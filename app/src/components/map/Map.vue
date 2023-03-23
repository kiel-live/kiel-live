<template>
  <div id="map" ref="mapElement" class="w-full h-full" />
</template>

<script lang="ts" setup>
// eslint-disable-next-line no-restricted-imports
import 'maplibre-gl/dist/maplibre-gl.css';

import { useElementSize } from '@vueuse/core';
import type {
  Feature,
  FeatureCollection,
  GeoJsonProperties as _GeoJsonProperties,
  Geometry,
  LineString,
  Point,
} from 'geojson';
import {
  AttributionControl,
  GeoJSONSource,
  GeolocateControl,
  LineLayerSpecification,
  LngLatLike,
  Map,
  NavigationControl,
  Source,
  SymbolLayerSpecification,
} from 'maplibre-gl';
import { computed, onMounted, Ref, ref, toRef, watch } from 'vue';

import { stops, subscribe, trips, vehicles } from '~/api';
import { Marker, StopType, VehicleType } from '~/api/types';
import BusIcon from '~/components/map/busIcon';
import { useColorMode } from '~/compositions/useColorMode';
import { brightMapStyle, darkMapStyle } from '~/config';

const props = withDefaults(
  defineProps<{
    selectedMarker: Partial<Marker>;
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

const vehiclesGeoJson = computed<Feature<Point, GeoJsonProperties>[]>(() =>
  Object.values(vehicles.value).map((v) => {
    const iconData = {
      kind: 'vehicle',
      type: v.type,
      name: v.name.split(' ')[0],
      focused: false,
      heading: v.location.heading,
    };

    return {
      type: 'Feature',
      properties: {
        kind: 'vehicle',
        type: v.type,
        name: v.name,
        id: v.id,
        number: v.name.split(' ')[0],
        to: v.name.split(' ').slice(1).join(' '),
        iconName:
          // eslint-disable-next-line no-nested-ternary
          v.type === 'bus' ? JSON.stringify(iconData) : colorScheme.value === 'dark' ? `dark-${v.type}` : v.type,
        iconNameFocused: v.type === 'bus' ? JSON.stringify({ ...iconData, focused: true }) : v.type,
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
      iconName: colorScheme.value === 'dark' ? `dark-${s.type}` : s.type,
      iconNameFocused: s.type,
    },
    geometry: {
      type: 'Point',
      coordinates: [s.location.longitude / 3600000, s.location.latitude / 3600000],
    },
  })),
);

const selectedMarker = toRef(props, 'selectedMarker');

const selectedVehicle = computed(() => {
  if (!selectedMarker.value.id) {
    return null;
  }
  return vehicles.value[selectedMarker.value.id];
});

const trip = computed(() => {
  if (!trips.value || !selectedVehicle.value) {
    return null;
  }
  return trips.value[selectedVehicle.value.tripId];
});

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

const stopsLayer: Ref<SymbolLayerSpecification> = computed(() => ({
  id: 'stops',
  type: 'symbol',
  source: 'geojson',
  filter: ['==', 'kind', 'stop'],
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
      ['get', 'iconNameFocused'],
      ['get', 'iconName'],
    ],
    'icon-size': 0.6,
    'icon-rotation-alignment': 'map',
    'icon-allow-overlap': true,
    'symbol-sort-key': ['match', ['get', 'number'], selectedVehicle.value?.name.split(' ')[0] ?? '', 2, 1],
  },
}));

const vehiclesLayer: Ref<SymbolLayerSpecification> = computed(() => ({
  id: 'vehicles',
  type: 'symbol',
  source: 'geojson',
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
    'icon-rotation-alignment': 'map',
    'icon-allow-overlap': true,
    'symbol-sort-key': ['match', ['get', 'number'], selectedVehicle.value?.name.split(' ')[0] ?? '', 2, 1],
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

const mapElement = ref(null);
const { width, height } = useElementSize(mapElement);

function flyTo(center: [number, number]) {
  if (!map) {
    return;
  }

  map.flyTo({
    center,
    padding: {
      // 768: md breakpoint
      // 320: sidebar width w-80
      left: width.value >= 768 ? 320 : 0,
      bottom: width.value >= 768 ? 0 : height.value * (2 / 3),
    },
  });
}

onMounted(async () => {
  void subscribe('data.map.vehicle.>', vehicles);
  void subscribe('data.map.stop.>', stops);

  let center: [number, number] = [10.1283, 54.3166];
  const { state: geolocationPermission } = await navigator.permissions.query({ name: 'geolocation' });
  if (geolocationPermission === 'granted') {
    const position = await new Promise<GeolocationPosition>((resolve, reject) => {
      navigator.geolocation.getCurrentPosition(resolve, reject);
    });
    center = [position.coords.longitude, position.coords.latitude];
  }

  const center: LngLatLike = [10.1283, 54.3166];

  map = new Map({
    container: 'map',
    // style: 'https://demotiles.maplibre.org/style.json',
    style: colorScheme.value === 'dark' ? darkMapStyle : brightMapStyle,
    minZoom: 5,
    maxZoom: 18,
    center,
    zoom: 14,
    // [west, south, east, north]
    maxBounds: [5.0, 46.0, 15.0, 57.0],
    attributionControl: false,
  });

  const attributionControl = new AttributionControl({ compact: true });
  map.addControl(attributionControl, 'bottom-left');

  map.addControl(
    new GeolocateControl({
      positionOptions: {
        enableHighAccuracy: true,
      },
      trackUserLocation: true,
    }),
    'bottom-right',
  );

  map.addControl(new NavigationControl({}), 'bottom-right');

  type IconData =
    | { kind: 'vehicle'; type: string; name: string; focused: boolean; heading: number }
    | { kind: 'stop'; type: string; name: string; focused: boolean };

  map.on('styleimagemissing', (e) => {
    if (e.id[0] !== '{') {
      return;
    }

    const iconData = JSON.parse(e.id) as IconData;
    if (iconData.kind === 'vehicle' && iconData.type === 'bus') {
      map.addImage(e.id, new BusIcon(map, iconData.focused, iconData.name, iconData.heading), {
        pixelRatio: 2,
      });
    } else {
      console.log(iconData);
    }
  });

  const loadImage = (name: string, url: string) =>
    new Promise<void>((resolve, reject) => {
      // eslint-disable-next-line promise/prefer-await-to-callbacks
      map.loadImage(url, (error, image) => {
        if (error) {
          reject(error);
        } else if (image) {
          map.addImage(name, image, { pixelRatio: 2 });
          resolve();
        }
      });
    });

  async function loadImages() {
    await loadImage('bus-stop', '/icons/bus.png');
    await loadImage('dark-bus-stop', '/icons/dark-bus.png');
    await loadImage('bike-stop', '/icons/bike.png');
    await loadImage('dark-bike-stop', '/icons/dark-bike.png');
  }

  map.on('load', () => {
    map.addSource('geojson', {
      type: 'geojson',
      data: Object.freeze(geojson.value),
    });

    map.addLayer(stopsLayer.value);

    map.addLayer(tripsLayer.value);

    map.addLayer(vehiclesLayer.value);

    initial = false;

    void loadImages();
  });

  // Change the cursor to a pointer when the it enters a feature in the 'symbols' layer.
  map.on('mouseenter', 'vehicles', () => {
    map.getCanvas().style.cursor = 'pointer';
  });

  // Change it back to a pointer when it leaves.
  map.on('mouseleave', 'vehicles', () => {
    map.getCanvas().style.cursor = '';
  });

  // Change the cursor to a pointer when the it enters a feature in the 'symbols' layer.
  map.on('mouseenter', 'stops', () => {
    map.getCanvas().style.cursor = 'pointer';
  });

  // Change it back to a pointer when it leaves.
  map.on('mouseleave', 'stops', () => {
    map.getCanvas().style.cursor = '';
  });

  map.on('click', (e) => {
    const features = map.queryRenderedFeatures(e.point, {
      layers: ['stops', 'vehicles'],
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
});

watch(colorScheme, () => {
  if (colorScheme.value === 'dark') {
    map?.setStyle(darkMapStyle);
  } else {
    map?.setStyle(brightMapStyle);
  }
  // TODO: properly re-render custom layers
  // location.reload();
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

watch(stopsLayer, () => {
  if (!map || initial) {
    return;
  }

  if (stopsLayer.value.layout) {
    Object.keys(stopsLayer.value.layout).forEach((key) => {
      if (stopsLayer.value.layout) {
        map.setLayoutProperty('stops', key, stopsLayer.value.layout[key as keyof typeof stopsLayer.value.layout]);
      }
    });
  }

  if (stopsLayer.value.paint) {
    Object.keys(stopsLayer.value.paint).forEach((key) => {
      if (stopsLayer.value.paint) {
        map.setPaintProperty('stops', key, stopsLayer.value.paint[key as keyof typeof stopsLayer.value.paint]);
      }
    });
  }
});

watch(vehiclesLayer, () => {
  if (!map || initial) {
    return;
  }

  if (vehiclesLayer.value.layout) {
    Object.keys(vehiclesLayer.value.layout).forEach((key) => {
      if (vehiclesLayer.value.layout) {
        map.setLayoutProperty(
          'vehicles',
          key,
          vehiclesLayer.value.layout[key as keyof typeof vehiclesLayer.value.layout],
        );
      }
    });
  }

  if (vehiclesLayer.value.paint) {
    Object.keys(vehiclesLayer.value.paint).forEach((key) => {
      if (vehiclesLayer.value.paint) {
        map.setPaintProperty('vehicles', key, vehiclesLayer.value.paint[key as keyof typeof vehiclesLayer.value.paint]);
      }
    });
  }
});

const selectedMarkerItem = computed(() => {
  const marker = selectedMarker.value;
  if (!marker) {
    return undefined;
  }
  return geojson.value.features.find((f) => f.properties.id === marker.id);
});
watch(selectedMarkerItem, (newSelectedMarkerItem, oldSelectedMarkerItem) => {
  if (
    !map ||
    !newSelectedMarkerItem ||
    newSelectedMarkerItem.properties.id === oldSelectedMarkerItem?.properties.id ||
    mapMovedManually.value
  ) {
    return;
  }

  flyTo((newSelectedMarkerItem.geometry as Point)?.coordinates as [number, number]);
});
</script>

<style scoped>
#map :deep(.maplibregl-ctrl-attrib) {
  @apply dark:bg-dark-400;
}

#map :deep(.maplibregl-ctrl-attrib a) {
  @apply dark:text-gray-300;
}

#map :deep(.maplibregl-ctrl-attrib-button) {
  @apply dark:(filter invert);
}

#map :deep(.maplibregl-ctrl-group) {
  @apply dark:bg-dark-400;
}

.dark #map :deep(.maplibregl-ctrl-group:not(:empty)) {
  box-shadow: 0 0 0 2px rgb(60 60 60);
}

#map :deep(.maplibregl-ctrl-group button + button) {
  @apply dark:border-t-dark-100;
}

#map :deep(.maplibregl-ctrl button .maplibregl-ctrl-icon) {
  @apply dark:(filter invert);
}
</style>
