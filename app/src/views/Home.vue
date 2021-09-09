<template>
  <div class="flex flex-col h-full w-full items-center justify-center">
    <div id="map" class="w-full h-full" />
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted } from 'vue';
import MapLibre from 'maplibre-gl';
import 'maplibre-gl/dist/maplibre-gl.css';

export default defineComponent({
  name: 'Home',

  setup() {
    let map: MapLibre.Map;

    onMounted(() => {
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
    });
  },
});
</script>
