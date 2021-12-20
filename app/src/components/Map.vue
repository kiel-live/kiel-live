<template>
  <div id="map" class="w-full h-full" />
</template>

<script lang="ts">
import { Feature, FeatureCollection, Point } from 'geojson';
import { CircleLayer, GeoJSONSource, Map, SymbolLayer } from 'maplibre-gl';
import { computed, defineComponent, onMounted, PropType, Ref, toRef, watch } from 'vue';

import { stops, subscribe, vehicles } from '~/api';
import BusIcon from '~/components/busIcon';
import { usePrefersColorSchemeDark } from '~/compositions/usePrefersColorScheme';
import { Marker } from '~/types';

export default defineComponent({
  // eslint-disable-next-line vue/multi-word-component-names
  name: 'Map',

  props: {
    selectedMarker: {
      type: Object as PropType<Marker>,
      default: () => ({}),
    },
  },

  emits: {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    markerClick: (_marker?: Marker) => true,
  },

  setup(props, { emit }) {
    let map: Map;
    let initial = true;

    const vehiclesGeoJson = computed<Feature[]>(() =>
      Object.values(vehicles.value).map((v) => ({
        type: 'Feature',
        properties: {
          type: 'vehicle',
          name: v.name,
          id: v.id,
          number: v.name.split(' ')[0],
          to: v.name.split(' ').slice(1).join(' '),
          iconName: `busIcon-unfocused-${v.name.split(' ')[0]}-${v.location.heading}`,
          iconNameFocused: `busIcon-focused-${v.name.split(' ')[0]}-${v.location.heading}`,
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
        properties: { type: 'stop', name: s.name, id: s.id },
        geometry: {
          type: 'Point',
          coordinates: [s.location.longitude / 3600000, s.location.latitude / 3600000],
        },
      })),
    );

    const geojson = computed<FeatureCollection>(() => ({
      type: 'FeatureCollection',
      features: [...vehiclesGeoJson.value, ...stopsGeoJson.value],
    }));

    const selectedMarker = toRef(props, 'selectedMarker');

    const stopsLayer: Ref<CircleLayer> = computed(() => ({
      id: 'stops',
      type: 'circle',
      source: 'geojson',
      filter: ['==', 'type', 'stop'],
      paint: {
        'circle-color': ['match', ['get', 'id'], selectedMarker.value.id || '', '#1673fc', '#4f96fc'],
        'circle-radius': ['match', ['get', 'id'], selectedMarker.value.id || '', 8, 5],
        'circle-opacity': selectedMarker.value.type === 'vehicle' ? 0.5 : 1,
      },
    }));

    const vehiclesLayer: Ref<SymbolLayer> = computed(() => ({
      id: 'vehicles',
      type: 'symbol',
      source: 'geojson',
      paint: {
        'icon-opacity': [
          'match',
          ['get', 'number'],
          selectedMarker.value.id && vehicles.value[selectedMarker.value.id]
            ? vehicles.value[selectedMarker.value.id].name.split(' ')[0]
            : '',
          1,
          selectedMarker.value.type === 'vehicle' ? 0.3 : 1,
        ],
      },
      filter: ['==', 'type', 'vehicle'],
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
        'symbol-sort-key': [
          'match',
          ['get', 'number'],
          selectedMarker.value.id && vehicles.value[selectedMarker.value.id]
            ? vehicles.value[selectedMarker.value.id].name.split(' ')[0]
            : '',
          2,
          1,
        ],
      },
    }));

    function flyTo(center: [number, number]) {
      if (!map) {
        return;
      }

      map.flyTo({
        center,
        // TODO: fix upstream type
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
        style: 'https://tiles.slucky.de/styles/gray-matter/style.json',
        accessToken: '',
        minZoom: 11,
        maxZoom: 18,
        center: [10.1283, 54.3166],
        zoom: 14,
        // [west, south, east, north]
        maxBounds: [9.8, 54.21, 10.44, 54.52],
      });

      // var nav = new MapLibre.NavigationControl();
      // map.addControl(nav, 'bottom-right');

      map.on('styleimagemissing', (e: { id: string; type: 'styleimagemissing' }) => {
        const [, focus, route, heading] = e.id.split('-');
        map.addImage(e.id, new BusIcon(map, focus === 'focused', route, Number.parseInt(heading, 10)), {
          pixelRatio: 2,
        });
      });

      map.on('load', () => {
        map.addSource('geojson', {
          type: 'geojson',
          data: Object.freeze(geojson.value),
        });

        map.addLayer(stopsLayer.value);

        map.addLayer(vehiclesLayer.value);

        initial = false;
      });

      map.on('click', 'vehicles', (e) => {
        if (!e.features || e.features.length === 0) {
          return;
        }
        const feature = e.features[0] as unknown as {
          geometry: Point;
          properties: Marker;
        };
        flyTo(feature.geometry.coordinates as [number, number]);
        emit('markerClick', { type: feature.properties.type, id: feature.properties.id });
      });

      // Change the cursor to a pointer when the it enters a feature in the 'symbols' layer.
      map.on('mouseenter', 'vehicles', () => {
        map.getCanvas().style.cursor = 'pointer';
      });

      // Change it back to a pointer when it leaves.
      map.on('mouseleave', 'vehicles', () => {
        map.getCanvas().style.cursor = '';
      });

      map.on('click', 'stops', (e) => {
        if (!e.features || e.features.length === 0) {
          return;
        }
        const feature = e.features[0] as unknown as {
          geometry: Point;
          properties: Marker;
        };
        map.flyTo({
          center: feature.geometry.coordinates as [number, number],
        });
        emit('markerClick', { type: feature.properties.type, id: feature.properties.id });
      });

      // Change the cursor to a pointer when the it enters a feature in the 'symbols' layer.
      map.on('mouseenter', 'stops', () => {
        map.getCanvas().style.cursor = 'pointer';
      });

      // Change it back to a pointer when it leaves.
      map.on('mouseleave', 'stops', () => {
        map.getCanvas().style.cursor = '';
      });

      // Deselect marker when the map is clicked.
      map.on('click', (e) => {
        const features = map.queryRenderedFeatures(e.point, {
          layers: ['stops', 'vehicles'],
        });

        if (features.length === 0) {
          emit('markerClick');
        }
      });
    });

    watch(usePrefersColorSchemeDark(), (prefersColorSchemeDark) => {
      if (prefersColorSchemeDark) {
        // TODO configurable tiles server
        map.setStyle('https://tiles.slucky.de/styles/gray-matter/style.json');
      } else {
        map.setStyle('https://tiles.slucky.de/styles/bright-matter/style.json');
      }
    });

    watch(geojson, () => {
      if (!map) {
        return;
      }

      const source = map.getSource('geojson') as GeoJSONSource | undefined;
      source?.setData(Object.freeze(geojson.value));
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
            map.setPaintProperty(
              'vehicles',
              key,
              vehiclesLayer.value.paint[key as keyof typeof vehiclesLayer.value.paint],
            );
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
      if (!map || !_selectedMarkerItem) {
        return;
      }

      flyTo((_selectedMarkerItem.geometry as Point)?.coordinates as [number, number]);
    });
  },
});
</script>
