<template>
  <StopPopup v-if="marker.type.endsWith('-stop')" :marker="marker" />
  <VehiclePopup v-else :marker="marker" />
</template>

<script lang="ts" setup>
import type { Marker } from '~/api/types';
import StopPopup from '~/components/popups/StopPopup.vue';
import VehiclePopup from '~/components/popups/VehiclePopup.vue';

defineProps<{
  marker: Marker;
}>();
</script>
