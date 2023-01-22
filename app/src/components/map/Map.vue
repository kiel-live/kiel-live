<template>
  <div id="map" class="w-full h-full" />
</template>

<script lang="ts" setup>
// eslint-disable-next-line no-restricted-imports
import 'maplibre-gl/dist/maplibre-gl.css';

import type { Feature, FeatureCollection, Point, Position } from 'geojson';
import {
  AttributionControl,
  CircleLayerSpecification,
  GeoJSONSource,
  GeolocateControl,
  LineLayerSpecification,
  Map,
  NavigationControl,
  Source,
  SymbolLayerSpecification,
} from 'maplibre-gl';
import { computed, onMounted, Ref, toRef, watch } from 'vue';

import { stops, subscribe, trips, vehicles } from '~/api';
import { Marker } from '~/api/types';
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

const mapMovedManually = computed({
  get: () => props.mapMovedManually,
  set: (value) => emit('update:mapMovedManually', value),
});

const colorScheme = useColorMode();

const vehiclesGeoJson = computed<Feature[]>(() =>
  Object.values(vehicles.value).map((v) => ({
    type: 'Feature',
    properties: {
      type: v.type,
      name: v.name,
      id: v.id,
      number: v.name.split(' ')[0],
      to: v.name.split(' ').slice(1).join(' '),
      iconName: `${v.type}-unfocused-${v.name.split(' ')[0]}-${v.location.heading}`,
      iconNameFocused: `${v.type}-focused-${v.name.split(' ')[0]}-${v.location.heading}`,
    },

    geometry: {
      type: 'Point',
      coordinates: [v.location.longitude / 3600000, v.location.latitude / 3600000],
    },
  })),
);

const stopsGeoJson = computed<Feature[]>(() =>
  Object.values(stops.value).map((s) => ({
    type: 'Feature',
    properties: { type: s.type, name: s.name, id: s.id },
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

const tripsGeoJson = computed<Feature[]>(() => {
  if (selectedVehicle.value?.type === 'bus' && trip.value?.arrivals) {
    const coordinates: Position[] = [];
    const { arrivals } = trip.value;
    arrivals.forEach((arrival, i) => {
      if (
        selectedVehicle.value &&
        (arrival.state === 'predicted' || arrival.state === 'stopping') &&
        (i === 0 || arrivals[i - 1]?.state === 'departed')
      ) {
        coordinates.push([
          selectedVehicle.value.location.longitude / 3600000,
          selectedVehicle.value.location.latitude / 3600000,
        ]);
      }
      const stop = stops.value[arrival.id];
      if (stop) {
        coordinates.push([stop.location.longitude / 3600000, stop.location.latitude / 3600000]);
      }
    });

    return [
      {
        type: 'Feature',
        properties: {
          type: 'trip',
        },
        geometry: {
          type: 'LineString',
          coordinates,
        },
      },
    ];
  }
  return [];
});

const geojson = computed<FeatureCollection>(() => ({
  type: 'FeatureCollection',
  features: [...vehiclesGeoJson.value, ...stopsGeoJson.value, ...tripsGeoJson.value],
}));

const stopsLayer: Ref<CircleLayerSpecification> = computed(() => ({
  id: 'stops',
  type: 'circle',
  source: 'geojson',
  filter: ['==', 'type', 'bus-stop'],
  paint: {
    'circle-color': ['match', ['get', 'id'], selectedMarker.value.id || '', '#1673fc', '#4f96fc'],
    'circle-radius': ['match', ['get', 'id'], selectedMarker.value.id || '', 8, 5],
    'circle-opacity': selectedMarker.value.type === 'bus' ? 0.5 : 1,
    'circle-stroke-opacity': 0,
    'circle-stroke-width': 5,
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
  filter: ['==', 'type', 'bus'],
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

function flyTo(center: [number, number]) {
  if (!map) {
    return;
  }

  map.flyTo({
    center,
    padding: {
      bottom: 500, // TODO use 3/4 of screen height
    },
  });
}

onMounted(async () => {
  await subscribe('data.map.vehicle.>', vehicles);
  await subscribe('data.map.stop.>', stops);

  map = new Map({
    container: 'map',
    // style: 'https://demotiles.maplibre.org/style.json',
    style: colorScheme.value === 'dark' ? darkMapStyle : brightMapStyle,
    minZoom: 5,
    maxZoom: 18,
    center: [10.1283, 54.3166],
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

  map.on('styleimagemissing', (e) => {
    const [type, focus, route, heading] = e.id.split('-');
    if (type === 'bus') {
      map.addImage(e.id, new BusIcon(map, focus === 'focused', route, Number.parseInt(heading, 10)), {
        pixelRatio: 2,
      });
    }
  });

  map.on('load', () => {
    map.addSource('geojson', {
      type: 'geojson',
      data: Object.freeze(geojson.value),
    });

    map.addLayer(stopsLayer.value);

    map.addLayer(tripsLayer.value);

    map.addLayer(vehiclesLayer.value);

    initial = false;
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
    flyTo(feature.geometry.coordinates as [number, number]);
    emit('markerClick', { type: feature.properties.type, id: feature.properties.id });
  });

  map.on('drag', () => {
    mapMovedManually.value = true;
  });
});

watch(colorScheme, () => {
  if (colorScheme.value === 'dark') {
    map.setStyle(darkMapStyle);
  } else {
    map.setStyle(brightMapStyle);
  }
  // TODO: properly re-render custom layers
  // eslint-disable-next-line no-restricted-globals
  location.reload();
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
  return geojson.value.features.find((f) => (f.properties as Marker).id === marker.id);
});
watch(selectedMarkerItem, (_selectedMarkerItem) => {
  if (!map || !_selectedMarkerItem || mapMovedManually.value) {
    return;
  }

  flyTo((_selectedMarkerItem.geometry as Point)?.coordinates as [number, number]);
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
