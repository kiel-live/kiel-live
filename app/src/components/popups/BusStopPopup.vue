<template>
  <div v-if="stop" class="flex flex-col min-h-0">
    <div class="flex pb-2 mb-2 border-b-1 dark:border-dark-800 space-x-2 items-center">
      <i-mdi-sign-real-estate v-if="stop.type === 'bus-stop'" />
      <span class="text-lg">{{ stop.name }}</span>
    </div>
    <div class="overflow-y-auto">
      <router-link
        v-for="arrival in stop.arrivals"
        :key="arrival.tripId"
        class="flex py-2 w-full not-last:border-b-1 dark:border-dark-600"
        :to="{ name: 'map-marker', params: { markerType: 'bus', markerId: arrival.vehicleId } }"
      >
        <i-fa-bus class="mr-2" />
        <span class="mr-2">{{ arrival.routeName }}</span>
        <span class="flex-grow">{{ arrival.direction }}</span>
        <span>{{ eta(arrival) }}</span>
        <div class="ml-2">
          <i-fa-solid-clock v-if="arrival.state === 'planned'" />
          <i-fa-solid-hand-paper v-if="arrival.state === 'stopping'" />
          <i-fa-solid-running v-if="arrival.state === 'predicted'" />
        </div>
      </router-link>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onUnmounted, PropType, toRef, watch } from 'vue';

import { stops, subscribe, unsubscribe } from '~/api';
import { Marker, StopArrival } from '~/api/types';

export default defineComponent({
  name: 'BusStopPopup',

  props: {
    marker: {
      type: Object as PropType<Marker>,
      required: true,
    },
  },

  setup(props) {
    const marker = toRef(props, 'marker');
    const stop = computed(() => stops.value[props.marker.id]);
    let subject: string | null = null;
    const eta = (arrival: StopArrival) => {
      const minutes = Math.round(arrival.eta / 60);

      if (arrival.state === 'stopping') {
        return 'hält';
      }
      if (minutes < 1) {
        return 'sofort';
      }

      return `${minutes} Min`;
    };

    watch(
      marker,
      async () => {
        if (subject !== null) {
          await unsubscribe(subject);
        }
        subject = `data.map.stop.${props.marker.id}`;
        await subscribe(subject, stops);
      },
      { immediate: true },
    );

    onUnmounted(async () => {
      if (subject !== null) {
        await unsubscribe(subject);
      }
    });

    return { stop, eta };
  },
});
</script>
