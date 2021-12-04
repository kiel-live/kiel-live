<template>
  <div v-if="vehicle" class="flex flex-col">
    <span class="mb-2 border-b-1 text-lg">{{ vehicle.name }}</span>
    <div v-if="trip" v-for="arrival in trip.arrivals" class="flex w-full cursor-pointer">
      <span class="mr-2">{{ arrival.planned }}</span>
      <TripMarker :marker="arrival.state === 'predicted' ? 'dot' : 'empty'" />
      <span>{{ arrival.name }}</span>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, watch, toRef, onUnmounted, computed } from 'vue';
import { subscribe, trips, unsubscribe, vehicles } from '~/api';
import { Marker } from '~/types';
import TripMarker from '../TripMarker.vue';

export default defineComponent({
  name: 'VehiclePopup',

  components: { TripMarker },

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
        if (!vehicle.value) {
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
