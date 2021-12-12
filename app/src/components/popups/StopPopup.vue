<template>
  <div v-if="stop" class="flex flex-col min-h-0">
    <div class="mb-2 border-b-1 text-lg">{{ stop.name }}</div>
    <div class="overflow-y-auto">
      <router-link
        v-for="arrival in stop.arrivals"
        :key="arrival.tripId"
        class="flex p-4 w-full not-last:border-b-1"
        :to="{ name: 'map-marker', params: { markerType: 'vehicle', markerId: arrival.vehicleId } }"
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
import { StopArrival } from '~/api/types';
import { Marker } from '~/types';

export default defineComponent({
  name: 'StopPopup',

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
