<template>
  <div v-if="vehicle" class="flex flex-col min-h-0">
    <div class="flex pb-2 mb-2 border-b-1 dark:border-gray-500 space-x-2 items-center">
      <i-fa-bus v-if="vehicle.type === 'bus'" />
      <span class="text-lg">{{ vehicle.name }}</span>
    </div>
    <div v-if="trip" class="overflow-y-auto">
      <router-link
        v-for="(arrival, i) in trip.arrivals"
        :key="arrival.id"
        :to="{ name: 'map-marker', params: { markerType: 'bus-stop', markerId: arrival.id } }"
        class="flex w-full items-center"
      >
        <span class="w-12 mr-2">{{ arrival.planned }}</span>
        <div
          class="marker relative flex justify-center items-center mx-4 h-12 w-8 after:(absolute top-0 h-full bg-gray-800 dark:bg-gray-300)"
        >
          <div
            v-if="arrival.state !== 'departed' && trip.arrivals[i - 1]?.state === 'departed'"
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
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onUnmounted, PropType, toRef, watch } from 'vue';

import { subscribe, trips, unsubscribe, vehicles } from '~/api';
import { Marker } from '~/api/types';

export default defineComponent({
  name: 'BusPopup',

  props: {
    marker: {
      type: Object as PropType<Marker>,
      required: true,
    },
  },

  setup(props) {
    const marker = toRef(props, 'marker');
    let subject: string | null = null;

    const vehicle = computed(() => vehicles.value[marker.value.id]);

    const trip = computed(() => {
      if (!trips.value) {
        return null;
      }
      return trips.value[vehicle.value.tripId];
    });

    watch(
      vehicle,
      async () => {
        if (subject !== null) {
          await unsubscribe(subject);
        }
        if (!vehicle.value) {
          return;
        }
        subject = `data.map.trip.${vehicle.value.tripId}`;
        await subscribe(subject, trips);
      },
      { immediate: true },
    );

    onUnmounted(async () => {
      if (subject !== null) {
        await unsubscribe(subject);
      }
    });
    return { trip, vehicle };
  },
});
</script>

<style scoped>
.marker::after {
  left: calc(50% - 0.05rem);
  width: 0.1rem;
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
