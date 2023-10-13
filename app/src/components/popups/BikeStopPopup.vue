<template>
  <div class="bg-red-400 w-full">
    <div>{{ stop.name }}</div>
  </div>
</template>

<script lang="ts" setup>
import { Stop } from '~/api/types';

const stop: Stop = {
  id: '123',
  alerts: [],
  arrivals: [],
  location: {
    latitude: 54.48855,
    longitude: 9.94689,
    heading: 0,
  },
  name: 'Reventloubrücke',
  provider: 'donkey',
  routes: [],
  type: 'bike-stop',
};
</script>
