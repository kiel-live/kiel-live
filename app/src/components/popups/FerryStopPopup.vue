<template>
  <div v-if="stop" class="flex flex-col min-h-0 flex-grow">
    <div class="flex flex-row pb-2 mb-2 border-b-1 dark:border-dark-100 items-center">
      <i-mdi-ferry />
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
          <div
            v-for="arrival in augmentedArrivals"
            :key="arrival.tripId"
            class="flex flex-col py-2 w-full not-last:border-b-1 dark:border-dark-300"
          >
            <div class="flex flex-row">
              <i-mdi-ferry class="mr-2" />
              <span class="mr-2">{{ arrival.routeName }}</span>
              <span class="flex-grow">{{ arrival.direction }}</span>
              <span>{{ arrival.planned }}</span>
              <div class="ml-2">
                <i-fa-solid-clock v-if="arrival.state === 'planned'" />
                <i-fa-solid-hand-paper v-if="arrival.state === 'stopping'" />
                <i-fa-solid-running v-if="arrival.state === 'predicted'" />
              </div>
            </div>
            <div class="flex flex-row gap-1 text-gray-500 dark:text-gray-400 text-xs">
              <template v-if="arrival.nextStopName">
                <span>{{ t('next_stop') }}</span>
                <span>{{ arrival.nextStopName }}</span>
              </template>
              <!-- <div v-else class="w-1/3 mb-1 bg-gray-500 dark:bg-gray-400 rounded-lg animate-pulse opacity-10" /> -->
              <span class="ml-auto">{{ arrival.platform }}</span>
            </div>
          </div>
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

const augmentedArrivals = computed<(Omit<StopArrival, 'eta'> & { nextStopName?: string; eta: string })[] | null>(() => {
  if (stop.value === undefined || !stop.value.arrivals) {
    return null;
  }

  return stop.value.arrivals.map((a) => {
    const trip = trips.value[a.tripId];

    let nextStopName: string | undefined;
    if (trip !== undefined && trip.arrivals !== undefined) {
      const nextStopIndex = trip.arrivals.findIndex((s) => s.id === props.marker.id);
      if (nextStopIndex !== -1) {
        nextStopName = trip.arrivals[nextStopIndex + 1]?.name;
      }
    }

    return {
      ...a,
      nextStopName,
      eta: eta(a),
    };
  });
});

watch(
  marker,
  async (newMarker, oldMarker) => {
    if (newMarker.id === oldMarker?.id) {
      return;
    }
    if (subject !== null) {
      void unsubscribe(subject);
    }
    subject = `data.map.stop.${newMarker.id}`;
    await subscribe(subject, stops);
  },
  { immediate: true },
);

const tripSubscriptions = new Set<string>();

// watch arrivals and subscribe to trips
// watch(
//   stop,
//   async (newStop, oldStop) => {
//     if (!newStop || newStop.arrivals === null || newStop.arrivals === oldStop?.arrivals) {
//       return;
//     }

//     oldStop?.arrivals?.forEach((arrival) => {
//       if (!newStop.arrivals?.some((a) => a.tripId === arrival.tripId)) {
//         tripSubscriptions.delete(arrival.tripId);
//         void unsubscribe(`data.map.trip.${arrival.tripId}`);
//       }
//     });

//     newStop.arrivals.forEach((arrival) => {
//       if (!tripSubscriptions.has(arrival.tripId)) {
//         tripSubscriptions.add(arrival.tripId);
//         void subscribe(`data.map.trip.${arrival.tripId}`, trips);
//       }
//     });
//   },
//   { immediate: true },
// );

onUnmounted(() => {
  if (subject !== null) {
    void unsubscribe(subject);
  }
  tripSubscriptions.forEach((tripId) => {
    void unsubscribe(`data.map.trip.${tripId}`);
  });
});
</script>
