<template>
  <div v-if="vehicle" class="flex flex-col min-h-0 flex-grow">
    <div class="flex flex-col pb-2 mb-2 border-b-1 dark:border-dark-100">
      <div class="flex space-x-2 items-center">
        <i-ic-outline-pedal-bike v-if="vehicle.type === 'bike'" />
        <h1 class="text-lg">{{ vehicle.name }}</h1>
      </div>
      <div class="flex mt-2 overflow-x-auto">
        <div class="inline-flex w-min gap-2 pb-4">
          <a
            :href="`https://nxtb.it/${vehicle.id.replace('nextbike-', '')}`"
            target="_blank"
            rel="noopener noreferrer"
            class="flex flex-shrink-0"
          >
            <button type="button" class="border border-blue-400 bg-blue-300 rounded-full px-4 py-1 dark:text-white">
              Rent this bike
            </button>
          </a>
          <a
            :href="`https://www.google.com/maps/place/${vehicle.location.latitude / 3600000},${vehicle.location.longitude / 3600000}`"
            target="_blank"
            rel="noopener noreferrer"
            class="flex flex-shrink-0"
          >
            <button type="button" class="border border-gray-400 rounded-full px-4 py-1 dark:text-white">
              Navigate to
            </button>
          </a>
          <button
            type="button"
            class="border border-gray-400 rounded-full px-4 py-1 dark:text-white flex flex-shrink-0"
          >
            Reserve bike
          </button>
        </div>
      </div>
    </div>

    <div class="flex flex-col overflow-y-auto gap-2">
      <div class="text-lg">
        <p class="underline">Price:</p>
        <p>1€ / 15 min</p>
        <p>12€ / 24 h</p>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, onUnmounted, toRef, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { subscribe, trips, unsubscribe, vehicles } from '~/api';
import { Marker, Vehicle } from '~/api/types';

const props = defineProps<{
  marker: Marker;
}>();

const { t } = useI18n();

const marker = toRef(props, 'marker');
let subject: string | null = null;

const vehicle = computed<Vehicle | undefined>(() => vehicles.value[marker.value.id]);

// const trip = computed(() => {
//   if (!trips.value || !vehicle.value) {
//     return null;
//   }
//   return trips.value[vehicle.value.tripId];
// });

watch(
  vehicle,
  async (newVehicle, oldVehicle) => {
    if (newVehicle?.tripId === oldVehicle?.tripId) {
      return;
    }
    if (subject !== null) {
      void unsubscribe(subject);
    }
    if (!newVehicle) {
      return;
    }
    subject = `data.map.trip.${newVehicle.tripId}`;
    await subscribe(subject, trips);
  },
  { immediate: true },
);

onUnmounted(() => {
  if (subject !== null) {
    void unsubscribe(subject);
  }
});
</script>
