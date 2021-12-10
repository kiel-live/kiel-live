<template>
  <div v-if="vehicle" class="flex flex-col min-h-0">
    <div class="pb-2 mb-2 border-b-1 dark:border-gray-500 text-lg">{{ vehicle.name }}</div>
    <div v-if="trip" class="overflow-y-auto">
      <router-link
        v-for="arrival in trip.arrivals"
        :key="arrival.id"
        :to="{ name: 'map-marker', params: { markerType: 'stop', markerId: arrival.id } }"
        class="flex w-full"
      >
        <span class="mr-2">{{ arrival.planned }}</span>
        <TripMarker :marker="arrival.state === 'predicted' ? 'dot' : 'empty'" />
        <span>{{ arrival.name }}</span>
      </router-link>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onUnmounted, PropType, toRef, watch } from 'vue';

import { subscribe, trips, unsubscribe, vehicles } from '~/api';
import TripMarker from '~/components/TripMarker.vue';
import { Marker } from '~/types';

export default defineComponent({
  name: 'VehiclePopup',

  components: { TripMarker },

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
