<template>
  <div id="map" class="w-full h-full" />
</template>

<script lang="ts">
import { CircleLayer, GeoJSONSource, GeoJSONSourceRaw, Map, SymbolLayer } from 'maplibre-gl';
import { computed, defineComponent, onMounted, PropType, Ref, toRef, watch } from 'vue';

import BusIcon from '~/components/busIcon';
import { usePrefersColorSchemeDark } from '~/compositions/usePrefersColorScheme';
import { Marker } from '~/types';

export default defineComponent({
  // eslint-disable-next-line vue/multi-word-component-names
  name: 'Map',

  props: {
    geojson: {
      type: Object as PropType<GeoJSONSourceRaw['data']>,
      required: true,
    },

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

    const geojson = toRef(props, 'geojson');
    const selectedMarker = toRef(props, 'selectedMarker');

    const stopsLayer: Ref<CircleLayer> = computed(() => ({
      id: 'stops',
      type: 'circle',
      source: 'geojson',
      filter: ['==', 'type', 'stop'],
      paint: {
        'circle-color': ['match', ['get', 'id'], selectedMarker.value.id || '', '#1673fc', '#4f96fc'],
        'circle-radius': ['match', ['get', 'id'], selectedMarker.value.id || '', 8, 5],
      },
    }));

    const vehiclesLayer: Ref<SymbolLayer> = computed(() => ({
      id: 'vehicles',
      type: 'symbol',
      source: 'geojson',
      paint: {
        'icon-opacity': ['match', ['get', 'number'], '', 1, 1],
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
        'symbol-sort-key': ['match', ['get', 'number'], '', 2, 1],
      },
    }));

    onMounted(async () => {
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

      map.on('styleimagemissing', (e) => {
        const [, focus, route, heading] = e.id.split('-');
        map.addImage(e.id, new BusIcon(map, focus === 'focused', route, heading), { pixelRatio: 2 });
      });

      map.on('load', () => {
        map.addSource('geojson', {
          type: 'geojson',
          data: Object.freeze(geojson.value),
        });

        map.addLayer(stopsLayer.value);

        map.addLayer(vehiclesLayer.value);
      });

      map.on('click', 'vehicles', (e) => {
        if (!e.features || e.features.length === 0) {
          return;
        }
        const feature = e.features[0];
        map.flyTo({
          center: feature.geometry.coordinates,
        });
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
        const feature = e.features[0];
        map.flyTo({
          center: feature.geometry.coordinates,
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
      if (!map) {
        return;
      }

      map.removeLayer('stops');
      map.addLayer(stopsLayer.value);
    });

    watch(vehiclesLayer, () => {
      if (!map) {
        return;
      }

      map.removeLayer('vehicles');
      map.addLayer(vehiclesLayer.value);
    });
  },
});
</script>
