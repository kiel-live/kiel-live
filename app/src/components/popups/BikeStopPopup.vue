<template>
  <div v-if="stop" class="flex flex-col min-h-0 flex-grow">
    <div class="flex flex-col pb-2 mb-2 border-b-1 dark:border-dark-100">
      <div class="flex space-x-2 items-center">
        <i-fa-bus v-if="stop.type === 'bike-stop'" />
        <h1 class="text-lg">{{ stop.name }}</h1>
      </div>
      <div class="flex mt-2 gap-2 overflow-x-auto">
        <button type="button" class="border border-blue-600 bg-blue-600 rounded-full px-4 py-1 text-white">
          Rent this bike
        </button>
        <button type="button" class="border border-gray-400 rounded-full px-4 py-1 text-white">Navigate to</button>
        <button type="button" class="border border-gray-400 rounded-full px-4 py-1 text-white">Reserve bike</button>
      </div>
    </div>

    <div class="flex flex-col overflow-y-auto gap-2">
      <div class="text-lg">
        <p class="underline">Preis:</p>
        <p>0,49€ zum Entsperren</p>
        <p>0,32€ pro Minute</p>
      </div>

      <div class="flex gap-2">
        <i-mdi-map-marker />
        <div>{{ stop.location.latitude }}, {{ stop.location.longitude }}</div>
      </div>
    </div>
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
