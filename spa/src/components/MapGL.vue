<template>
  <MglMap
    :mapStyle="mapStyle"
    :center="center"
    :minZoom="minZoom"
    :maxZoom="maxZoom"
    :maxBounds="maxBounds"
  >
    <MglNavigationControl position="top-right"/>
    <MglGeolocateControl position="top-right" />
    <MglGeojsonLayer
      sourceId= "fakeID"
      :source="geoJsonSource"
      layerId="elf"
      :layer="geoJsonlayer"
    />
    <MglMarker :coordinates="bus" color="blue"/>
  </MglMap>
</template>

<script>
import Mapbox from 'mapbox-gl';
import {
  MglMap,
  MglNavigationControl,
  MglGeolocateControl,
  MglGeojsonLayer,
  MglMarker,
} from 'vue-mapbox';
import { mapState } from 'vuex';
import elf from '@/assets/11.json';

export default {
  components: {
    MglMap,
    MglNavigationControl,
    MglGeolocateControl,
    MglGeojsonLayer,
    MglMarker,
  },
  data() {
    return {
      mapStyle: 'https://maps.targomo.com/styles/dark-matter-gl-style.json',
      center: [10.1283, 54.3166],
      minZoom: 11,
      maxZoom: 18,
      // [west, south, east, north]
      maxBounds: [9.8, 54.21, 10.44, 54.52],
      geoJsonSource: {
        type: 'geojson',
        data: elf,
      },
      geoJsonlayer: {
        id: 'elf',
        type: 'line',
        paint: {
          'line-color': '#00ffff',
        },
      },
      bus: [10.1283, 54.3166],
    };
  },
  computed: {
    ...mapState({
      vehicles: (state) => state.map.vehicles,
      stops: (state) => state.map.stops,
      savedView: (state) => state.map.savedView,
    }),
  },
  created() {
    // We need to set mapbox-gl library here in order to use it in template
    this.mapbox = Mapbox;
  },
  mounted() {
    this.$store.dispatch('map/load');
  },
};
</script>

<style lang="scss" scoped>
  .map-container {
    position: relative;
  }
</style>
