<template>
  <BusPopup v-if="marker.type === 'bus'" :marker="marker" />
  <BusStopPopup v-else-if="marker.type === 'bus-stop'" :marker="marker" />
  <FallbackStopPopup v-else :marker="marker" />
</template>

<script lang="ts" setup>
import { Marker } from '~/api/types';
import BusPopup from '~/components/popups/BusPopup.vue';
import BusStopPopup from '~/components/popups/BusStopPopup.vue';
import FallbackStopPopup from '~/components/popups/FallbackStopPopup.vue';

defineProps<{
  marker: Marker;
}>();
</script>
