<template>
  <div v-if="vehicle" class="flex flex-col">
    <span class="mb-2 border-b-1 text-lg">{{ vehicle.name }}</span>
    <template v-if="trip">
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
    </template>
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
