<template>
  <div v-if="vehicle" class="flex flex-col">
    <span class="mb-2 border-b-1 text-lg">{{ vehicle.name }}</span>
    <template v-if="trip" v-for="arrival in trip.arrivals">
    <span>{{ arrival.name }}</span>
    </template>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, watch, toRef, onUnmounted, computed } from 'vue';
import { subscribe, trips, unsubscribe, vehicles } from '~/api';
import { Marker } from '~/types';

export default defineComponent({
  name: 'VehiclePopup',

  props: {
    marker: {
      type: Object as PropType<Marker>,
      required: true,
    },
  },

  setup(props, { emit }) {
    const marker = toRef(props, 'marker');
    let subject: string | null = null;

    const vehicle = computed(() => {
      return vehicles.value[marker.value.id];
    });

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
          unsubscribe(subject);
        }
        if(!vehicle.value) {
          return;
        }
        subject = `data.map.trip.${vehicle.value.tripId}`;
        await subscribe(subject, trips);
      },
      { immediate: true },
    );

    onUnmounted(() => {
      if (subject !== null) {
        unsubscribe(subject);
      }
    });
    return { trip, vehicle };
  },
});
</script>
