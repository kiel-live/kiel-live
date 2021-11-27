<template>
  <div class="flex flex-col">
    <span class="mb-2 border-b-1 text-lg">Stop</span>
    <span>{{ marker.id }}</span>
    <div v-if="stop" v-for="arrival in stop.arrivals">
      <span class="text-sm">{{ arrival.direction }}</span>
      <span class="text-sm">{{ eta(arrival) }}</span>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, onMounted, computed } from 'vue';
import { subscribe } from '~/api';
import { Marker } from '~/types';
import { stops } from '~/api';
import { StopArrival } from '~/api/types';

export default defineComponent({
  name: 'StopPopup',

  props: {
    marker: {
      type: Object as PropType<Marker>,
      required: true,
    },
  },

  setup(props, { emit }) {
    const stop = computed(() => stops.value[props.marker.id]);

    const eta = (arrival: StopArrival) => {
      const minutes = Math.round(arrival.eta / 60);

      if (arrival.state === 'stopping') {
        return 'hält';
      }
      if (arrival.state === 'predicted') {
        return arrival.planned;
      }
      if (minutes < 1) {
        return 'sofort';
      }

      return `${minutes} Min`;
    };

    onMounted(async () => {
      await subscribe(`data.map.stop.${props.marker.id}`, stops);
    });

    return { stop, stops, eta };
  },
});
</script>
