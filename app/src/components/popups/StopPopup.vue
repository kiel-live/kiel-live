<template>
  <div class="flex flex-col">
    <span class="mb-2 border-b-1 text-lg">Stop</span>
    <span>{{ marker.id }}</span>
    <div v-if="stop" v-for="arrival in stop.arrivals">
      <span class="text-sm">{{ arrival.direction }}</span>
      <span class="text-sm">{{ dayjs.duration(arrival.eta, 'seconds').format('m [Min]') }}</span>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, onMounted, computed } from 'vue';
import { subscribe } from '~/api';
import { Marker } from '~/types';
import { stops } from '~/api';
import dayjs from '~/lib/dayjs';

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

    onMounted(async () => {
      await subscribe(`data.map.stop.${props.marker.id}`, stops);
    });

    return { stop, stops, dayjs };
  },
});
</script>
