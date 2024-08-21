<template>
  <div v-if="vehicle" class="flex flex-col min-h-0 flex-grow">
    <div class="flex pb-2 mb-2 border-b-1 dark:border-dark-100 space-x-2 items-center">
      <i-fa-bus v-if="vehicle.type === 'bus'" />
      <i-ic-outline-pedal-bike v-else-if="vehicle.type === 'bike'" />
      <i-ic-baseline-directions-car v-else-if="vehicle.type === 'car'" />
      <i-ic-twotone-electric-scooter v-else-if="vehicle.type === 'e-scooter'" />
      <i-ic-twotone-electric-scooter v-else-if="vehicle.type === 'ferry'" />
      <i-ic-baseline-train v-else-if="vehicle.type === 'train'" />
      <i-ic-outline-subway v-else-if="vehicle.type === 'subway'" />
      <i-ic-baseline-tram v-else-if="vehicle.type === 'tram'" />
      <i-ic-baseline-moped v-else-if="vehicle.type === 'moped'" />
      <i-ic-baseline-electric-moped v-else-if="vehicle.type === 'e-moped'" />
      <h1 class="text-lg">{{ vehicle.name }}</h1>
    </div>

    <Actions :actions="vehicle.actions ?? []" />

    <template v-if="trip">
      <div v-if="trip.arrivals?.length" class="overflow-y-auto">
        <router-link
          v-for="(arrival, i) in trip.arrivals"
          :key="arrival.id"
          :to="{ name: 'map-marker', params: { markerType: `${vehicle.type}-stop`, markerId: arrival.id } }"
          class="flex w-full items-center"
          :class="{
            'text-gray-500 dark:text-gray-400': arrival.state === 'departed',
            'mt-6': i === 0 && arrival.state === 'predicted',
          }"
        >
          <span class="w-14 min-w-12">{{ arrival.planned }}</span>
          <div
            class="marker relative flex justify-center items-center mx-4 h-12 w-8 min-w-4 after:(absolute top-0 h-full bg-gray-800 dark:bg-gray-300)"
            :class="{ 'after:(bg-gray-500 dark:bg-gray-400)': arrival.state === 'departed' }"
          >
            <div
              v-if="
                arrival.state !== 'departed' &&
                (trip.arrivals[i - 1] === undefined || trip.arrivals[i - 1].state === 'departed')
              "
              class="vehicle before:(h-4 w-4 bg-red-700 rounded-full)"
              :class="{ driving: arrival.state === 'predicted' }"
            >
              <div class="pulsating border-3 border-red-700 border-solid rounded-full" />
            </div>
            <div
              v-if="
                (arrival.state !== 'departed' && trip.arrivals[i - 1]?.state !== 'departed') ||
                arrival.state === 'predicted'
              "
              class="rounded-full h-4 w-4 flex items-center justify-center bg-gray-800 dark:bg-gray-300"
            />
          </div>
          <span class="w-full">{{ arrival.name }}</span>
        </router-link>
      </div>
      <NoData v-else>{{ t('trip_expired') }}</NoData>
    </template>
    <i-fa-solid-circle-notch v-else class="mx-auto mt-4 text-3xl animate-spin" />
  </div>
  <NoData v-else>{{ t('trip_does_not_exist') }}</NoData>
</template>

<script lang="ts" setup>
import { computed, onUnmounted, toRef, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { subscribe, trips, unsubscribe, vehicles } from '~/api';
import { Marker, Vehicle } from '~/api/types';
import NoData from '~/components/NoData.vue';
import Actions from '~/components/popups/Actions.vue';

const props = defineProps<{
  marker: Marker;
}>();

const { t } = useI18n();

const marker = toRef(props, 'marker');
let subject: string | null = null;

const vehicle = computed<Vehicle | undefined>(() => vehicles.value[marker.value.id]);

const trip = computed(() => {
  if (!trips.value || !vehicle.value || !vehicle.value.tripId) {
    return null;
  }
  return trips.value[vehicle.value.tripId];
});

watch(
  vehicle,
  async (newVehicle, oldVehicle) => {
    if (newVehicle?.tripId === oldVehicle?.tripId) {
      return;
    }
    if (subject !== null) {
      void unsubscribe(subject);
    }

    // don't subscribe if no vehicle was selected or it doesn't have a trip
    if (!newVehicle || !newVehicle.tripId) {
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

<style scoped>
.marker::after {
  left: calc(50% - 1px);
  width: 2px;
  content: '';
}

.vehicle {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 2;
}
.vehicle.driving {
  top: -0.25rem;
}

.vehicle::before {
  display: block;
  content: '';
}

.pulsating {
  position: absolute;
  width: 2rem;
  height: 2rem;
  left: calc(50% - 1rem);
  top: calc(50% - 1rem);
  transform: translate(-50%, -50%);
  animation: pulsate 1.5s ease-out;
  animation-iteration-count: infinite;
  opacity: 0;
}

@keyframes pulsate {
  0% {
    transform: scale(0.1, 0.1);
    opacity: 0;
  }
  50% {
    opacity: 1;
  }
  100% {
    transform: scale(1.2, 1.2);
    opacity: 0;
  }
}
</style>
