<template>
  <div v-if="stop" class="flex flex-col min-h-0 flex-grow">
    <div class="flex flex-row pb-2 mb-2 border-b-1 border-gray-200 dark:border-neutral-700 items-center">
      <i-mdi-ferry v-if="stop.type === 'ferry-stop'" />
      <i-mdi-sign-real-estate v-else />
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

    <Actions :actions="stop.actions ?? []" />

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

      <router-link
        v-for="arrival in augmentedArrivals"
        :key="arrival.tripId"
        class="flex flex-col py-2 w-full not-last:border-b-1 dark:border-neutral-800"
        :to="{
          name: 'map-marker',
          params: { markerType: stop.type.replace('-stop', ''), markerId: arrival.vehicleId },
        }"
      >
        <div class="flex flex-row items-center">
          <i-mdi-bus v-if="arrival.type === 'bus'" class="mr-2 w-6 h-6" />
          <i-mdi-ferry v-else-if="arrival.type === 'ferry'" class="mr-2" />
          <i-mdi-tram v-else-if="arrival.type === 'tram'" class="mr-2 w-6 h-6" />
          <i-carbon-train-profile v-else-if="arrival.type === 'train'" class="mr-2" />

          <span class="mr-2">{{ arrival.routeName }}</span>
          <span class="flex-grow">{{ arrival.direction }}</span>
          <span>{{ arrival.eta ?? arrival.planned }}</span>
          <div class="ml-2 flex items-center">
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
      </router-link>
      <router-link
        v-for="vehicle in stop.vehicles"
        :key="vehicle.id"
        class="flex flex-col py-2 w-full not-last:border-b-1 dark:border-neutral-800"
        :to="{
          name: 'map-marker',
          params: { markerType: vehicle.type, markerId: vehicle.id },
        }"
      >
        <div class="flex flex-row">
          <i-ic-outline-pedal-bike v-if="vehicle.type === 'bike'" class="mr-2" />
          <i-mdi-ferry v-else-if="vehicle.type === 'ferry'" class="mr-2" />

          <span class="mr-2">{{ vehicle.name }}</span>
        </div>
      </router-link>
      <NoData v-if="augmentedArrivals && augmentedArrivals.length === 0">
        {{ t('no_bus_wants_to_stop_here_right_now') }}
      </NoData>
      <i-fa-solid-circle-notch v-if="!augmentedArrivals && !stop.vehicles" class="m-auto text-3xl animate-spin" />
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
import type { Marker, StopArrival } from '~/api/types';
import { computed, onBeforeUnmount, toRef } from 'vue';

import { useI18n } from 'vue-i18n';
import { api } from '~/api';
import Button from '~/components/atomic/Button.vue';
import NoData from '~/components/NoData.vue';
import { useFavorites } from '~/compositions/useFavorites';

const props = defineProps<{
  marker: Marker;
}>();

const { addFavorite, removeFavorite, isFavorite } = useFavorites();
const { t } = useI18n();

const marker = toRef(props, 'marker');

const { stop, unsubscribe: unsubscribeStop } = api.useStop(computed(() => props.marker.id));

const eta = (arrival: StopArrival) => {
  if (arrival.eta === 0) {
    return null;
  }

  const minutes = Math.round(arrival.eta / 60);

  if (arrival.state === 'stopping') {
    return t('stopping');
  }
  if (minutes < 1) {
    return t('immediately');
  }

  return t('minutes', { minutes });
};

const augmentedArrivals = computed<(Omit<StopArrival, 'eta'> & { nextStopName?: string; eta: string | null })[] | null>(
  () => {
    if (!stop.value?.arrivals) {
      return null;
    }

    return stop.value.arrivals
      .toSorted((a, b) => {
        if (a.eta === 0 || b.eta === 0) {
          return a.planned.localeCompare(b.planned);
        }
        return a.eta - b.eta;
      })

      .map((a) => {
        // const trip = trips.value[a.tripId];

        // let nextStopName: string | undefined;
        // if (trip !== undefined && trip.arrivals !== undefined) {
        //   const nextStopIndex = trip.arrivals.findIndex((s) => s.id === props.marker.id);
        //   if (nextStopIndex !== -1) {
        //     nextStopName = trip.arrivals[nextStopIndex + 1]?.name;
        //   }
        // }

        return {
          ...a,
          nextStopName: undefined,
          eta: eta(a),
        };
      });
  },
);

// const tripSubscriptions = new Set<string>();

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

//     newStop.arrivals?.forEach((arrival) => {
//       if (!tripSubscriptions.has(arrival.tripId)) {
//         tripSubscriptions.add(arrival.tripId);
//         void subscribe(`data.map.trip.${arrival.tripId}`, trips);
//       }
//     });
//   },
//   { immediate: true },
// );

onBeforeUnmount(async () => {
  await unsubscribeStop();
});
</script>
