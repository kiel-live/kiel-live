<template>
  <VehiclePopup v-if="marker.type === 'bus'" :marker="marker" />
  <StopPopup v-else-if="marker.type.endsWith('-stop')" :marker="marker" />
  <FallbackStopPopup v-else :marker="marker" />
</template>

<script lang="ts" setup>
import { Marker } from '~/api/types';
import FallbackStopPopup from '~/components/popups/FallbackStopPopup.vue';
import StopPopup from '~/components/popups/StopPopup.vue';
import VehiclePopup from '~/components/popups/VehiclePopup.vue';

defineProps<{
  marker: Marker;
}>();
</script>
