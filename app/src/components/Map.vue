<template>
  <div id="map" class="w-full h-full" />
</template>

<script lang="ts">
import { defineComponent, onMounted, PropType, ref, toRef, watch } from 'vue';
import MapLibre, { GeoJSONSource, GeoJSONSourceRaw } from 'maplibre-gl';
import 'maplibre-gl/dist/maplibre-gl.css';
import { Marker } from '~/types';
import BusIcon from '~/components/busIcon';
import { usePrefersColorSchemeDark } from '~/compositions/usePrefersColorScheme';
import { log } from 'console';

export default defineComponent({
  name: 'Map',

  props: {
    geojson: {
      type: Object as PropType<GeoJSONSourceRaw['data']>,
      required: true,
    },
  },

  emits: {
    markerClick: (_marker?: Marker) => true,
  },

  setup(props, { emit }) {
    let map: MapLibre.Map;

    const geojson = toRef(props, 'geojson');

    onMounted(async () => {
      map = new MapLibre.Map({
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

        map.addLayer({
          id: 'vehicles',
          type: 'symbol',
          source: 'geojson',
          paint: {
            'icon-opacity': ['match', ['get', 'number'], '', 1, 1],
          },
          filter: ['==', 'type', 'vehicle'],
          layout: {
            'icon-image': ['match', ['get', 'id'], '', ['get', 'iconNameFocused'], ['get', 'iconName']],
            'icon-rotation-alignment': 'map',
            'icon-allow-overlap': true,
            'symbol-sort-key': ['match', ['get', 'number'], '', 2, 1],
          },
          // filter: ['!', ['has', 'point_count']],
          // paint: {
          //   'circle-color': '#007cbf',
          //   'circle-radius': 7,
          // },
        });

        map.addLayer({
          id: 'stops',
          type: 'circle',
          source: 'geojson',
          filter: ['==', 'type', 'stop'],
          paint: {
            'circle-color': '#4f96fc',
            'circle-radius': 5,
          },
        });
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
      map.on('mouseenter', 'vehicles', function () {
        map.getCanvas().style.cursor = 'pointer';
      });

      // Change it back to a pointer when it leaves.
      map.on('mouseleave', 'vehicles', function () {
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
      map.on('mouseenter', 'stops', function () {
        map.getCanvas().style.cursor = 'pointer';
      });

      // Change it back to a pointer when it leaves.
      map.on('mouseleave', 'stops', function () {
        map.getCanvas().style.cursor = '';
      });

      // Deselect marker when the map is clicked.
      map.on('click', function (e) {
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
  },
});
</script>
