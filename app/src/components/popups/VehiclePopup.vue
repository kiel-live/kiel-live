<template>
  <div v-if="vehicle" class="flex flex-col min-h-0">
    <div class="flex pb-2 mb-2 border-b-1 dark:border-gray-500 space-x-2 items-center">
      <i-fa-bus v-if="vehicle.type === 'bus'" />
      <span class="text-lg">{{ vehicle.name }}</span>
    </div>
    <div v-if="trip" class="overflow-y-auto">
      <router-link
        v-for="arrival in trip.arrivals"
        :key="arrival.id"
        :to="{ name: 'map-marker', params: { markerType: 'bus-stop', markerId: arrival.id } }"
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
import { Marker } from '~/api/types';
import TripMarker from '~/components/TripMarker.vue';

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
