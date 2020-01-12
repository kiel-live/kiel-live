<template>
  <Map :showVehicle="vehicleId" :showStop="stopId" />
</template>

<script>
import Map from '../components/Map.vue';

export default {
  name: 'MapView',
  components: {
    Map,
  },
  computed: {
    vehicleId() {
      return this.$route.params.vehicle;
    },
    stopId() {
      return this.$route.params.stop;
    },
  },
};
</script>
