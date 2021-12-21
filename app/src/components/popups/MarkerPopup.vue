<template>
  <BusPopup v-if="marker.type === 'bus'" :marker="marker" />
  <BusStopPopup v-else-if="marker.type === 'bus-stop'" :marker="marker" />
  <pre v-else>{{ marker }}</pre>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue';

import { Marker } from '~/api/types';
import BusPopup from '~/components/popups/BusPopup.vue';
import BusStopPopup from '~/components/popups/BusStopPopup.vue';

export default defineComponent({
  name: 'MarkerPopup',

  components: { BusPopup, BusStopPopup },

  props: {
    marker: {
      type: Object as PropType<Marker>,
      required: true,
    },
  },
});
</script>
