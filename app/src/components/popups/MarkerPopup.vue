<template>
  <BusPopup v-if="marker.type === 'bus'" :marker="marker" />
  <StopPopup v-else-if="marker.type.endsWith('-stop')" :marker="marker" />
  <FallbackPopup v-else :marker="marker" />
</template>

<script lang="ts" setup>
import { Marker } from '~/api/types';
import BusPopup from '~/components/popups/BusPopup.vue';
import FallbackPopup from '~/components/popups/FallbackPopup.vue';
import StopPopup from '~/components/popups/StopPopup.vue';

defineProps<{
  marker: Marker;
}>();
</script>
