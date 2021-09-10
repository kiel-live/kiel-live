<template>
  <div id="map" class="w-full h-full" />
</template>

<script lang="ts">
import { defineComponent, onMounted, PropType, ref, toRef, watch } from 'vue';
import MapLibre, { GeoJSONSource, GeoJSONSourceRaw } from 'maplibre-gl';
import 'maplibre-gl/dist/maplibre-gl.css';
import { Marker } from '~/types';

export default defineComponent({
  name: 'Map',

  props: {
    geojson: {
      type: Object as PropType<GeoJSONSourceRaw>,
      required: true,
    },
  },

  emits: {
    markerClick: (marker: Marker) => true,
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

      map.on('load', () => {
        map.addSource('geojson', geojson.value);

        map.addLayer({
          id: 'vehicles',
          source: 'geojson',
          type: 'circle',
          filter: ['has', 'point_count'],
          paint: {
            'circle-color': ['step', ['get', 'point_count'], '#bfbc00', 5, '#bf7c00', 10, '#bf3c00'],
            'circle-radius': ['step', ['get', 'point_count'], 10, 5, 15, 10, 20],
          },
        });

        map.addLayer({
          id: 'vehicles-unclustered',
          source: 'geojson',
          type: 'circle',
          filter: ['!', ['has', 'point_count']],
          paint: {
            'circle-color': '#007cbf',
            'circle-radius': 7,
          },
        });
      });

      map.on('click', 'vehicles-unclustered', (e) => {
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
      map.on('mouseenter', 'vehicles-unclustered', function () {
        map.getCanvas().style.cursor = 'pointer';
      });

      // Change it back to a pointer when it leaves.
      map.on('mouseleave', 'vehicles-unclustered', function () {
        map.getCanvas().style.cursor = '';
      });
    });

    watch(geojson, () => {
      if (!map) {
        return;
      }

      const source = map.getSource('geojson') as GeoJSONSource | undefined;
      source?.setData(Object.freeze(geojson.value.data));
    });
  },
});
</script>
