<template>
  <BusPopup v-if="marker.type === 'bus'" :marker="marker" />
  <BikePopup v-if="marker.type === 'bike'" :marker="marker" />
  <StopPopup v-else-if="marker.type.endsWith('-stop')" :marker="marker" />
  <FallbackStopPopup v-else :marker="marker" />
</template>

<script lang="ts" setup>
import { Marker } from '~/api/types';
import BikePopup from '~/components/popups/BikePopup.vue';
import BusPopup from '~/components/popups/BusPopup.vue';
import FallbackStopPopup from '~/components/popups/FallbackStopPopup.vue';
import StopPopup from '~/components/popups/StopPopup.vue';

defineProps<{
  marker: Marker;
}>();
</script>
