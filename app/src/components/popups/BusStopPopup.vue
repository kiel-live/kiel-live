<template>
  <div v-if="stop" class="flex flex-col min-h-0 flex-grow">
    <div class="flex flex-row pb-2 mb-2 border-b-1 dark:border-dark-100 items-center">
      <i-mdi-sign-real-estate v-if="stop.type === 'bus-stop'" />
      <h1 class="text-lg ml-2">{{ stop.name }}</h1>
      <Button
        v-if="isFavorite(stop)"
        class="text-yellow-300 ml-auto border-0"
        :title="t('remove_favorite')"
        @click="removeFavorite(stop)"
      >
        <i-ph-star-fill />
      </Button>
      <Button v-else class="ml-auto border-0" :title="t('add_favorite')" @click="addFavorite(stop)">
        <i-ph-star-bold />
      </Button>
    </div>

    <div class="flex flex-col flex-grow overflow-y-auto">
      <div
        v-if="stop.alerts && stop.alerts.length >= 1"
        class="bg-red-300 dark:bg-red-800 bg-opacity-50 dark:bg-opacity-50 p-2 mb-2 rounded-md"
      >
        <div class="flex items-center border-b-1 border-gray-500 dark:border-gray-300 mb-2">
          <i-mdi-alert class="mr-2" /><span class="font-bold">{{ t('alerts') }}</span>
        </div>
        <ul>
          <li v-for="(alert, i) in stop.alerts" :key="i" class="items-center ml-5 list-outside list-disc">
            {{ alert }}
          </li>
        </ul>
      </div>

      <template v-if="augmentedArrivals">
        <template v-if="augmentedArrivals.length > 0">
          <router-link
            v-for="arrival in augmentedArrivals"
            :key="arrival.tripId"
            class="flex flex-col py-2 w-full not-last:border-b-1 dark:border-dark-300"
            :to="{ name: 'map-marker', params: { markerType: 'bus', markerId: arrival.vehicleId } }"
          >
            <div class="flex flex-row">
              <i-fa-bus class="mr-2" />
              <span class="mr-2">{{ arrival.routeName }}</span>
              <span class="flex-grow">{{ arrival.direction }}</span>
              <span>{{ eta(arrival) }}</span>
              <div class="ml-2">
                <i-fa-solid-clock v-if="arrival.state === 'planned'" />
                <i-fa-solid-hand-paper v-if="arrival.state === 'stopping'" />
                <i-fa-solid-running v-if="arrival.state === 'predicted'" />
              </div>
            </div>
            <div>
              <span class="text-gray-600 text-sm">Next Stop: {{ arrival.nextStopName }}</span>
            </div>
          </router-link>
        </template>
        <NoData v-else>{{ t('no_bus_wants_to_stop_here_right_now') }}</NoData>
      </template>
      <i-fa-solid-circle-notch v-else class="m-auto text-3xl animate-spin" />
    </div>
  </div>
  <NoData v-else>
    {{ t('this_stop_probably_does_not_exist') }}
    <Button
      v-if="isFavorite(marker)"
      class="mt-2"
      @click="
        () => {
          removeFavorite(marker);
          $router.replace({ name: 'home' });
        }
      "
    >
      <i-ph-star-fill class="mr-2 text-yellow-300" /><span>{{ t('remove_favorite') }}</span>
    </Button>
  </NoData>
</template>

<script setup lang="ts">
import { computed, onUnmounted, toRef, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { stops, subscribe, trips, unsubscribe } from '~/api';
import { Marker, StopArrival } from '~/api/types';
import Button from '~/components/atomic/Button.vue';
import NoData from '~/components/NoData.vue';
import { useFavorites } from '~/compositions/useFavorites';

const props = defineProps<{
  marker: Marker;
}>();

const { addFavorite, removeFavorite, isFavorite } = useFavorites();
const { t } = useI18n();

const marker = toRef(props, 'marker');
const stop = computed(() => stops.value[props.marker.id]);
let subject: string | null = null;
const eta = (arrival: StopArrival) => {
  const minutes = Math.round(arrival.eta / 60);

  if (arrival.state === 'stopping') {
    return t('stopping');
  }
  if (minutes < 1) {
    return t('immediately');
  }

  return t('minutes', { minutes });
};

const augmentedArrivals = computed<(StopArrival & { nextStopName?: string })[]>(() => {
  if (stop.value === undefined || !stop.value.arrivals) {
    return [];
  }

  return stop.value.arrivals.map((a) => {
    const trip = trips.value[a.tripId];
    console.log(trip);

    if (trip === undefined || trip.arrivals === undefined) {
      return a;
    }

    const nextStopIndex = trip.arrivals.findIndex((s) => s.id === props.marker.id);
    if (nextStopIndex === -1) {
      return a;
    }

    return {
      ...a,
      nextStopName: trip.arrivals[nextStopIndex + 1]?.name,
    };
  });
});

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

const tripSubscriptions = new Set<string>();

// watch arrivals and subscribe to trips
watch(
  stop,
  async (newStop, oldStop) => {
    console.log('stop changed', newStop, oldStop);

    if (newStop === null || newStop.arrivals === null || newStop.arrivals === oldStop?.arrivals) {
      return;
    }

    // eslint-disable-next-line no-restricted-syntax
    for await (const arrival of newStop.arrivals) {
      if (!tripSubscriptions.has(arrival.tripId)) {
        tripSubscriptions.add(arrival.tripId);
        await subscribe(`data.map.trip.${arrival.tripId}`, trips);
      }
    }
  },
  { immediate: true },
);

onUnmounted(async () => {
  if (subject !== null) {
    await unsubscribe(subject);
  }
  // eslint-disable-next-line no-restricted-syntax
  for await (const tripId of tripSubscriptions) {
    await unsubscribe(`data.map.trip.${tripId}`);
  }
});
</script>
