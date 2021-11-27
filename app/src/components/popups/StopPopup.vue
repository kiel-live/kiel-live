<template>
  <div v-if="stop" class="flex flex-col">
    <span class="mb-2 border-b-1 text-lg">{{ stop.name }}</span>
    <div class="flex p-4 w-full border-b-1 cursor-pointer" v-for="arrival in stop.arrivals">
      <IconBus class="mr-2" />
      <span class="flex-grow">{{ arrival.direction }}</span>
      <span>{{ eta(arrival) }}</span>
      <div class="ml-2">
        <IconClock v-if="arrival.state === 'planned'" />
        <IconHandPaper v-if="arrival.state === 'stopping'" />
        <IconRunning v-if="arrival.state === 'predicted'" />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, onMounted, computed, watch, toRef } from 'vue';
import { subscribe } from '~/api';
import { Marker } from '~/types';
import { stops } from '~/api';
import { StopArrival } from '~/api/types';
import IconClock from '~icons/fa-solid/clock';
import IconHandPaper from '~icons/fa-solid/hand-paper';
import IconRunning from '~icons/fa-solid/running';
import IconBus from '~icons/fa/bus';

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

  setup(props, { emit }) {
    const marker = toRef(props, 'marker');
    const stop = computed(() => stops.value[props.marker.id]);

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
        await subscribe(`data.map.stop.${props.marker.id}`, stops);
      },
      { immediate: true },
    );

    return { stop, stops, eta };
  },
});
</script>
