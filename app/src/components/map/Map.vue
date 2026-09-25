<template>
  <div id="map" ref="mapElement" class="h-full w-full" />
</template>

<script lang="ts" setup>
import type {
  GeoJsonProperties as _GeoJsonProperties,
  Feature,
  FeatureCollection,
  Geometry,
  LineString,
  Point,
} from 'geojson';
import type {
  AllLayoutProperties,
  AllPaintProperties,
  CircleLayerSpecification,
  GeoJSONSource,
  LineLayerSpecification,
  Source,
  SymbolLayerSpecification,
} from 'maplibre-gl';
import type { Ref } from 'vue';
import type { Bounds, Marker, StopType, VehicleType } from '~/api/types';
import { refThrottled, useElementSize } from '@vueuse/core';

import { AttributionControl, GeolocateControl, Map, NavigationControl, setWorkerUrl } from 'maplibre-gl';
import maplibreWorkerUrl from 'maplibre-gl/dist/maplibre-gl-worker.mjs?worker&url';
import { computed, onBeforeUnmount, onMounted, ref, toRef, useTemplateRef, watch } from 'vue';
import { api } from '~/api';
import { labeledVehicleTypes, stopColor, vehicleColors } from '~/components/map/markerColors';
import { createVehicleBadgeIcon, createVehicleNoseIcon } from '~/components/map/vehicleIcon';
import { useColorMode } from '~/compositions/useColorMode';
import { useUserSettings } from '~/compositions/useUserSettings';
import { brightMapStyle, darkMapStyle } from '~/config';

import 'maplibre-gl/dist/maplibre-gl.css';

const props = withDefaults(
  defineProps<{
    selectedMarker?: Partial<Marker>;
    mapMovedManually: boolean;
  }>(),
  {
    selectedMarker: () => ({}),
  },
);

const emit = defineEmits<{
  (e: 'markerClick', marker?: Marker): void;
  (e: 'update:mapMovedManually', mapMovedManually: boolean): void;
}>();

let map: Map;
let initial = true;

type GeoJsonProperties = _GeoJsonProperties & {
  type: StopType | VehicleType | 'trip';
  id?: string;
};

const mapMovedManually = computed({
  get: () => props.mapMovedManually,
  set: (value) => emit('update:mapMovedManually', value),
});

const colorScheme = useColorMode();

const bounds = ref<Bounds>();
const throttledBounds = refThrottled(bounds, 500);
const { stops, unsubscribe: unsubscribeStops } = api.useStops(throttledBounds);
const { vehicles, unsubscribe: unsubscribeVehicles } = api.useVehicles(throttledBounds);

const vehiclesGeoJson = computed<Feature<Point, GeoJsonProperties>[]>(() =>
  Object.values(vehicles.value).map((v) => {
    const heading = v.location.heading;
    const hasHeading = heading !== undefined && heading !== null;

    return {
      type: 'Feature',
      properties: {
        kind: 'vehicle',
        type: v.type,
        name: v.name,
        id: v.id,
        number: labeledVehicleTypes.has(v.type) ? v.name.split(' ')[0] : '',
        to: v.name.split(' ').slice(1).join(' '),
        iconName: v.type,
        iconNameFocused: `${v.type}-selected`,
        noseIcon: hasHeading ? `${v.type}-nose` : '',
        noseIconFocused: hasHeading ? `${v.type}-nose-selected` : '',
        heading: heading ?? 0,
      },

      geometry: {
        type: 'Point',
        coordinates: [v.location.longitude / 3600000, v.location.latitude / 3600000],
      },
    };
  }),
);

const stopsGeoJson = computed<Feature<Point, GeoJsonProperties>[]>(() =>
  Object.values(stops.value).map((s) => ({
    type: 'Feature',
    properties: {
      kind: 'stop',
      type: s.type,
      name: s.name,
      id: s.id,
    },
    geometry: {
      type: 'Point',
      coordinates: [s.location.longitude / 3600000, s.location.latitude / 3600000],
    },
  })),
);

const selectedMarker = toRef(props, 'selectedMarker');

const { vehicle: selectedVehicle, unsubscribe: unsubscribeSelectedVehicle } = api.useVehicle(
  computed(() => selectedMarker.value.id),
);

const { trip, unsubscribe: unsubscribeTrip } = api.useTrip(computed(() => selectedVehicle.value?.tripId));

const tripsGeoJson = computed<Feature<LineString, GeoJsonProperties>[]>(() => {
  if (selectedVehicle.value?.type === 'bus' && trip.value?.path) {
    return [
      {
        type: 'Feature',
        properties: {
          type: 'trip',
        },
        geometry: {
          type: 'LineString',
          coordinates: trip.value.path.map((p) => [p.longitude / 3600000, p.latitude / 3600000]),
        },
      },
    ];
  }
  return [];
});

const geojson = computed<FeatureCollection<Geometry, GeoJsonProperties>>(() => ({
  type: 'FeatureCollection',
  features: [...vehiclesGeoJson.value, ...stopsGeoJson.value, ...tripsGeoJson.value],
}));

// not reactive: a color scheme change already triggers a full page reload (see
// the `colorScheme` watcher below)
const labelTextColor = colorScheme.value === 'light' ? '#1a1a1a' : '#f5f5f5';
const labelHaloColor = colorScheme.value === 'light' ? '#ffffff' : '#1a1a1a';
// same as the popup's dark-mode text color (DetailsPopup.vue's dark:text-gray-300)
const markerOutlineColor = colorScheme.value === 'light' ? '#ffffff' : '#d1d5db';

const stopsLayer: Ref<CircleLayerSpecification> = computed(() => ({
  id: 'stops',
  type: 'circle',
  source: 'geojson',
  minzoom: 10,
  filter: ['==', 'kind', 'stop'],
  paint: {
    'circle-color': stopColor,
    'circle-radius': ['match', ['get', 'id'], selectedMarker.value.id || '', 10, 7],
    'circle-stroke-width': ['match', ['get', 'id'], selectedMarker.value.id || '', 3, 2],
    'circle-stroke-color': markerOutlineColor,
    'circle-opacity': selectedMarker.value.type === 'bus' ? 0.4 : 1,
    'circle-stroke-opacity': selectedMarker.value.type === 'bus' ? 0.4 : 1,
  },
}));

// stop names, only shown once zoomed in far enough to avoid cluttering the map
const stopsLabelLayer: Ref<SymbolLayerSpecification> = computed(() => ({
  id: 'stops-labels',
  type: 'symbol',
  source: 'geojson',
  minzoom: 14,
  filter: ['==', 'kind', 'stop'],
  layout: {
    'text-field': ['get', 'name'],
    'text-size': 12,
    'text-offset': [0, 1],
    'text-anchor': 'top',
    'text-allow-overlap': false,
    'text-optional': true,
  },
  paint: {
    'text-color': labelTextColor,
    'text-halo-color': labelHaloColor,
    'text-halo-width': 1.2,
  },
}));

// two layers so the nose can rotate with the heading while the badge (icon +
// label) stays upright; see vehicleIcon.ts for how they stay aligned
const vehicleIconSize = 1.2;

const vehiclesNoseLayer: Ref<SymbolLayerSpecification> = computed(() => ({
  id: 'vehicles-nose',
  type: 'symbol',
  source: 'geojson',
  minzoom: 14,
  filter: ['==', 'kind', 'vehicle'],
  paint: {
    'icon-opacity': [
      'match',
      ['get', 'number'],
      selectedVehicle.value?.name.split(' ')[0] ?? '',
      1,
      selectedMarker.value.type === 'bus' ? 0.3 : 1,
    ],
  },
  layout: {
    'icon-image': [
      'match',
      ['get', 'id'],
      selectedMarker.value.id || '',
      ['get', 'noseIconFocused'],
      ['get', 'noseIcon'],
    ],
    'icon-size': vehicleIconSize,
    'icon-rotate': ['get', 'heading'],
    'icon-rotation-alignment': 'map',
    'icon-allow-overlap': true,
    'symbol-sort-key': ['match', ['get', 'number'], selectedVehicle.value?.name.split(' ')[0] ?? '', 2, 1],
  },
}));

const vehiclesLayer: Ref<SymbolLayerSpecification> = computed(() => ({
  id: 'vehicles',
  type: 'symbol',
  source: 'geojson',
  minzoom: 6,
  paint: {
    'icon-opacity': [
      'match',
      ['get', 'number'],
      selectedVehicle.value?.name.split(' ')[0] ?? '',
      1,
      selectedMarker.value.type === 'bus' ? 0.3 : 1,
    ],
  },
  filter: ['==', 'kind', 'vehicle'],
  layout: {
    'icon-image': [
      'match',
      ['get', 'id'],
      selectedMarker.value.id || '',
      ['get', 'iconNameFocused'],
      ['get', 'iconName'],
    ],
    'icon-size': vehicleIconSize,
    'icon-allow-overlap': true,
    'symbol-sort-key': ['match', ['get', 'number'], selectedVehicle.value?.name.split(' ')[0] ?? '', 2, 1],
  },
}));

// route/line number, only shown once zoomed in far enough to avoid cluttering the map
const vehiclesLabelLayer: Ref<SymbolLayerSpecification> = computed(() => ({
  id: 'vehicles-labels',
  type: 'symbol',
  source: 'geojson',
  minzoom: 14,
  filter: ['==', 'kind', 'vehicle'],
  layout: {
    'text-field': ['get', 'number'],
    'text-size': 13,
    'text-offset': [0, 1.3],
    'text-anchor': 'top',
    'text-allow-overlap': false,
    'text-optional': true,
  },
  paint: {
    'text-color': labelTextColor,
    'text-halo-color': labelHaloColor,
    'text-halo-width': 1.2,
  },
}));

const tripsLayer: Ref<LineLayerSpecification> = computed(() => ({
  id: 'trips',
  type: 'line',
  source: 'geojson',
  filter: ['==', 'type', 'trip'],
  paint: {
    'line-width': 3,
    'line-color': 'rgb(170, 0, 0)',
  },
}));

const mapElement = useTemplateRef('mapElement');
const { width, height } = useElementSize(mapElement);

function flyTo(center: [number, number]) {
  if (!map) {
    return;
  }

  map.flyTo({
    center,
    padding: {
      // 768: md breakpoint
      // 480: left nav (80px) + sidebar width (400px)
      left: width.value >= 768 ? 480 : 0,
      bottom: width.value >= 768 ? 0 : height.value * (2 / 3),
    },
  });
}

onMounted(async () => {
  const { lastLocation } = useUserSettings();

  setWorkerUrl(maplibreWorkerUrl);

  map = new Map({
    container: 'map',
    // style: 'https://demotiles.maplibre.org/style.json',
    style: colorScheme.value === 'dark' ? darkMapStyle : brightMapStyle,
    minZoom: 5,
    maxZoom: 18,
    center: lastLocation.value.center,
    zoom: lastLocation.value.zoom,
    pitch: lastLocation.value.pitch,
    bearing: lastLocation.value.bearing,
    // [west, south, east, north]
    maxBounds: [5.0, 46.0, 15.0, 57.0],
    attributionControl: false,
  });

  const attributionControl = new AttributionControl({ compact: true });
  map.addControl(attributionControl, 'bottom-left');

  const geolocateControl = new GeolocateControl({
    positionOptions: {
      enableHighAccuracy: true,
    },
    trackUserLocation: true,
  });

  map.addControl(geolocateControl, 'bottom-right');

  map.addControl(new NavigationControl({}), 'bottom-right');

  // Trigger geolocation if permission has been granted
  if (navigator.permissions) {
    const permissionStatus = await navigator.permissions.query({ name: 'geolocation' });
    if (permissionStatus.state === 'granted') {
      geolocateControl.trigger();
    }
  }

  async function loadVehicleIcons() {
    await Promise.all(
      (Object.keys(vehicleColors) as VehicleType[]).flatMap((type) => {
        const color = vehicleColors[type];
        return [
          createVehicleBadgeIcon({ type, color, outlineColor: markerOutlineColor }).then((icon) =>
            map.addImage(type, icon, { pixelRatio: 2 }),
          ),
          createVehicleBadgeIcon({ type, color, outlineColor: markerOutlineColor, selected: true }).then((icon) =>
            map.addImage(`${type}-selected`, icon, { pixelRatio: 2 }),
          ),
          map.addImage(`${type}-nose`, createVehicleNoseIcon(color, markerOutlineColor), { pixelRatio: 2 }),
          map.addImage(`${type}-nose-selected`, createVehicleNoseIcon(color, markerOutlineColor, true), {
            pixelRatio: 2,
          }),
        ];
      }),
    );
  }

  map.on('load', () => {
    void loadVehicleIcons();

    map.addSource('geojson', {
      type: 'geojson',
      data: Object.freeze(geojson.value),
    });

    map.addLayer(stopsLayer.value);
    map.addLayer(stopsLabelLayer.value);
    map.addLayer(tripsLayer.value);
    map.addLayer(vehiclesNoseLayer.value);
    map.addLayer(vehiclesLayer.value);
    map.addLayer(vehiclesLabelLayer.value);
    applyLayerOrder();

    bounds.value = {
      north: map.getBounds().getNorth(),
      east: map.getBounds().getEast(),
      south: map.getBounds().getSouth(),
      west: map.getBounds().getWest(),
    };

    initial = false;
  });

  function addPointerOnHover(layerName: string) {
    map.on('mouseenter', layerName, () => {
      map.getCanvas().style.cursor = 'pointer';
    });
    map.on('mouseleave', layerName, () => {
      map.getCanvas().style.cursor = '';
    });
  }

  addPointerOnHover('vehicles');
  addPointerOnHover('vehicles-labels');
  addPointerOnHover('stops');
  addPointerOnHover('stops-labels');

  map.on('click', (e) => {
    const features = map.queryRenderedFeatures(e.point, {
      layers: ['stops', 'stops-labels', 'vehicles', 'vehicles-labels'],
    });

    // Deselect marker when the map is clicked
    if (features.length === 0) {
      emit('markerClick');
      return;
    }

    const feature = features[0] as unknown as {
      geometry: Point;
      properties: Marker;
    };

    // Prevent reloading the same marker
    if (feature.properties.id === selectedMarker.value.id) {
      return;
    }

    mapMovedManually.value = false;
    emit('markerClick', { type: feature.properties.type, id: feature.properties.id });
  });

  map.on('drag', () => {
    mapMovedManually.value = true;
  });

  map.on('move', () => {
    lastLocation.value = {
      center: map.getCenter(),
      zoom: map.getZoom(),
      pitch: map.getPitch(),
      bearing: map.getBearing(),
    };
    bounds.value = {
      north: map.getBounds().getNorth(),
      east: map.getBounds().getEast(),
      south: map.getBounds().getSouth(),
      west: map.getBounds().getWest(),
    };
  });

  map.on('idle', () => {
    if (mapElement.value) {
      mapElement.value.setAttribute('data-idle', 'true');
    }
  });
});

onBeforeUnmount(async () => {
  await unsubscribeStops();
  await unsubscribeVehicles();
  await unsubscribeSelectedVehicle();
  await unsubscribeTrip();
});

watch(colorScheme, () => {
  if (colorScheme.value === 'dark') {
    map.setStyle(darkMapStyle);
  } else {
    map.setStyle(brightMapStyle);
  }

  // TODO: properly re-render custom layers
  window.location.reload();
});

watch(geojson, () => {
  if (!map) {
    return;
  }

  const geoJSONSource = map.getSource('geojson');
  const isGeoJsonSource = (source?: Source): source is GeoJSONSource => source?.type === 'geojson';
  if (isGeoJsonSource(geoJSONSource)) {
    geoJSONSource.setData(Object.freeze(geojson.value));
  }
});

function syncMapLayer(layerId: string, layer: { layout?: Record<string, unknown>; paint?: Record<string, unknown> }) {
  if (!map || initial) {
    return;
  }

  if (layer.layout) {
    Object.entries(layer.layout).forEach(([key, value]) => {
      map.setLayoutProperty(layerId, key as keyof AllLayoutProperties, value as never);
    });
  }

  if (layer.paint) {
    Object.entries(layer.paint).forEach(([key, value]) => {
      map.setPaintProperty(layerId, key as keyof AllPaintProperties, value as never);
    });
  }
}

watch(stopsLayer, () => syncMapLayer('stops', stopsLayer.value));
watch(vehiclesLayer, () => syncMapLayer('vehicles', vehiclesLayer.value));
watch(vehiclesNoseLayer, () => syncMapLayer('vehicles-nose', vehiclesNoseLayer.value));

const selectedMarkerItem = computed(() => {
  const marker = selectedMarker.value;
  if (!marker) {
    return undefined;
  }
  return geojson.value.features.find((f) => f.properties.id === marker.id);
});
watch(selectedMarkerItem, (newSelectedMarkerItem, oldSelectedMarkerItem) => {
  if (!map || !newSelectedMarkerItem || newSelectedMarkerItem.properties.id === oldSelectedMarkerItem?.properties.id) {
    return;
  }

  flyTo((newSelectedMarkerItem.geometry as Point)?.coordinates as [number, number]);
});

// stops normally sit behind vehicles, but a selected stop is brought to the
// front so it isn't hidden by vehicle markers on top of it
const defaultLayerOrder = ['stops', 'stops-labels', 'trips', 'vehicles-nose', 'vehicles', 'vehicles-labels'];
const stopSelectedLayerOrder = ['trips', 'vehicles-nose', 'vehicles', 'vehicles-labels', 'stops', 'stops-labels'];

const layerOrder = computed(() =>
  selectedMarkerItem.value?.properties.kind === 'stop' ? stopSelectedLayerOrder : defaultLayerOrder,
);

function applyLayerOrder() {
  if (!map) {
    return;
  }
  layerOrder.value.forEach((id) => map.moveLayer(id));
}

watch(layerOrder, applyLayerOrder);
</script>

<style scoped>
@reference "tailwindcss";

#map :deep(.maplibregl-ctrl-attrib) {
  @apply dark:bg-neutral-800;
}

#map :deep(.maplibregl-ctrl-attrib a) {
  @apply dark:text-gray-300;
}

#map :deep(.maplibregl-ctrl-attrib-button) {
  @apply dark:invert dark:filter;
}

#map :deep(.maplibregl-ctrl-group) {
  @apply dark:bg-neutral-800;
}

.dark #map :deep(.maplibregl-ctrl-group:not(:empty)) {
  box-shadow: 0 0 0 2px rgb(60 60 60);
}

#map :deep(.maplibregl-ctrl-group button + button) {
  @apply dark:border-t-neutral-700;
}

#map :deep(.maplibregl-ctrl button .maplibregl-ctrl-icon) {
  @apply dark:invert dark:filter;
}
</style>
