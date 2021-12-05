<template>
  <div v-if="stop" class="flex flex-col">
    <span class="mb-2 border-b-1 text-lg">{{ stop.name }}</span>
    <router-link
      v-for="arrival in stop.arrivals"
      :key="arrival.tripId"
      class="flex p-4 w-full border-b-1 cursor-pointer"
      :to="{ name: 'map-marker', params: { markerType: 'vehicle', markerId: arrival.vehicleId } }"
    >
      <IconBus class="mr-2" />
      <span class="mr-2">{{ arrival.routeName }}</span>
      <span class="flex-grow">{{ arrival.direction }}</span>
      <span>{{ eta(arrival) }}</span>
      <div class="ml-2">
        <IconClock v-if="arrival.state === 'planned'" />
        <IconHandPaper v-if="arrival.state === 'stopping'" />
        <IconRunning v-if="arrival.state === 'predicted'" />
      </div>
    </router-link>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onUnmounted, PropType, toRef, watch } from 'vue';

import { stops, subscribe, unsubscribe } from '~/api';
import { StopArrival } from '~/api/types';
import { Marker } from '~/types';
import IconBus from '~icons/fa/bus';
import IconClock from '~icons/fa-solid/clock';
import IconHandPaper from '~icons/fa-solid/hand-paper';
import IconRunning from '~icons/fa-solid/running';

export default defineComponent({
  name: 'StopPopup',

  components: {
    IconClock,
    IconHandPaper,
    IconRunning,
    IconBus,
  },

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
          unsubscribe(subject);
        }
        subject = `data.map.stop.${props.marker.id}`;
        await subscribe(subject, stops);
      },
      { immediate: true },
    );

    onUnmounted(() => {
      if (subject !== null) {
        unsubscribe(subject);
      }
    });

    return { stop, eta };
  },
});
</script>
