<template>
  <Story>
    <BusStopPopup :marker="marker" />
  </Story>
</template>

<script lang="ts" setup>
import { Marker } from '~/api/types';

import BusStopPopup from './BusStopPopup.vue';

const marker: Marker = {
  id: '123',
  type: 'bus-stop',
};
</script>
