<template>
  <Map :focusVehicle="vehicleId" :focusStop="stopId" />
</template>

<script>
import Map from '../components/MapGL.vue';

export default {
  name: 'MapView',
  components: {
    Map,
  },
  computed: {
    vehicleId() {
      return this.$route.params.vehicle || null;
    },
    stopId() {
      return this.$route.params.stop || null;
    },
  },
};
</script>
