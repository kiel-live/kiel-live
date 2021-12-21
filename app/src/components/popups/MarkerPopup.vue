<template>
  <VehiclePopup v-if="marker.type === 'bus'" :marker="marker" />
  <StopPopup v-else-if="marker.type === 'bus-stop'" :marker="marker" />
  <pre v-else>{{ marker }}</pre>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue';

import { Marker } from '~/api/types';
import StopPopup from '~/components/popups/StopPopup.vue';
import VehiclePopup from '~/components/popups/VehiclePopup.vue';

export default defineComponent({
  name: 'MarkerPopup',

  components: { VehiclePopup, StopPopup },

  props: {
    marker: {
      type: Object as PropType<Marker>,
      required: true,
    },
  },
});
</script>
