<template>
  <div class="">
    <VehiclePopup v-if="marker.type === 'vehicle'" :marker="marker" />
    <StopPopup v-else-if="marker.type === 'stop'" :marker="marker" />
    <TripPopup v-else-if="marker.type === 'trip'" :marker="marker" />
    <RoutePopup v-else-if="marker.type === 'route'" :marker="marker" />
    <pre v-else>{{ marker }}</pre>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType, ref, toRef } from 'vue';
import { Marker } from '~/types';
import VehiclePopup from '~/components/popups/VehiclePopup.vue';
import StopPopup from '~/components/popups/StopPopup.vue';
import TripPopup from '~/components/popups/TripPopup.vue';
import RoutePopup from '~/components/popups/RoutePopup.vue';

export default defineComponent({
  name: 'MarkerPopup',

  components: { VehiclePopup, StopPopup, TripPopup, RoutePopup },

  props: {
    marker: {
      type: Object as PropType<Marker>,
      required: true,
    },
  },

  setup(props, { emit }) {},
});
</script>
