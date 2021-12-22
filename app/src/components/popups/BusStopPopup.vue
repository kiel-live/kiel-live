<template>
  <div v-if="stop" class="flex flex-col min-h-0">
    <div class="flex flex-row pb-2 mb-2 border-b-1 dark:border-dark-800 items-center flex-grow">
      <i-mdi-sign-real-estate v-if="stop.type === 'bus-stop'" />
      <span class="text-lg ml-2">{{ stop.name }}</span>
      <div class="flex ml-auto items-center cursor-pointer select-none">
        <i-ph-star-fill v-if="isFavorite(stop.id)" class="text-yellow-300" @click="removeFavorite(stop)" />
        <i-ph-star-bold v-else @click="addFavorite(stop)" />
      </div>
    </div>
    <div v-if="stop.arrivals?.length > 0" class="overflow-y-auto">
      <router-link
        v-for="arrival in stop.arrivals"
        :key="arrival.tripId"
        class="flex py-2 w-full not-last:border-b-1 dark:border-dark-600"
        :to="{ name: 'map-marker', params: { markerType: 'bus', markerId: arrival.vehicleId } }"
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
    <i-fa-solid-circle-notch v-else class="mx-auto mt-4 text-3xl animate-spin" />
  </div>
  <NoData v-else>Diese Haltestelle gibt es wohl nicht.</NoData>
</template>

<script lang="ts">
import { computed, defineComponent, onUnmounted, PropType, toRef, watch } from 'vue';

import { stops, subscribe, unsubscribe } from '~/api';
import { Marker, StopArrival } from '~/api/types';
import NoData from '~/components/NoData.vue';
import { useFavorites } from '~/compositions/useFavorites';

export default defineComponent({
  name: 'BusStopPopup',

  components: { NoData },

  props: {
    marker: {
      type: Object as PropType<Marker>,
      required: true,
    },
  },

  setup(props) {
    const { addFavorite, removeFavorite, isFavorite } = useFavorites();

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

    return { stop, eta, addFavorite, removeFavorite, isFavorite };
  },
});
</script>
