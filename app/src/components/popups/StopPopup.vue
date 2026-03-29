<template>
  <div v-if="stop" class="flex min-h-0 grow flex-col">
    <div class="mb-2 flex flex-row items-center border-b border-gray-200 pb-2 dark:border-neutral-600">
      <i-mdi-ferry v-if="stop.type === 'ferry-stop'" />
      <i-mdi-sign-real-estate v-else />
      <h1 class="ml-2 text-lg">{{ stop.name }}</h1>
      <Button
        v-if="isFavorite(stop)"
        class="ml-auto border-0 text-yellow-300"
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

    <div class="flex grow flex-col overflow-y-auto">
      <div
        v-if="stop.alerts && stop.alerts.length >= 1"
        class="bg-opacity-50 dark:bg-opacity-50 mb-2 rounded-md bg-red-300 p-2 dark:bg-red-800"
      >
        <div class="mb-2 flex items-center border-b border-gray-500 dark:border-gray-300">
          <i-carbon-warning-alt-filled class="mr-2" /><span class="font-bold">{{ t('alerts') }}</span>
        </div>
        <ul>
          <li v-for="(alert, i) in stop.alerts" :key="i" class="ml-5 list-outside list-disc items-center">
            {{ alert }}
          </li>
        </ul>
      </div>

      <router-link
        v-for="departure in sortedDepartures"
        :key="departure.tripId"
        class="flex w-full flex-col border-gray-200 py-2 not-last:border-b dark:border-neutral-700"
        :to="{
          name: 'map-marker',
          params: { markerType: 'vehicle', markerId: departure.vehicleId ?? '-' },
        }"
      >
        <div class="flex flex-row items-center">
          <i-mdi-bus v-if="departure.type === 'bus'" class="mr-2 h-6 w-6" />
          <i-mdi-ferry v-else-if="departure.type === 'ferry'" class="mr-2" />
          <i-carbon-tram v-else-if="departure.type === 'tram'" class="mr-2 h-6 w-6" />
          <i-carbon-train-profile v-else-if="departure.type === 'train'" class="mr-2" />

          <span class="mr-2">{{ departure.routeName }}</span>
          <span class="grow">{{ departure.direction }}</span>
          <span>{{ getBoardTime(departure) }}</span>
          <div class="ml-2 flex items-center">
            <i-ph-clock-fill v-if="departure.state === 'planned'" />
            <i-ph-hand-fill v-if="departure.state === 'stopping'" />
            <i-ph-person-simple-run-bold v-if="departure.state === 'predicted'" />
          </div>
        </div>
      </router-link>
      <router-link
        v-for="vehicle in stop.vehicles"
        :key="vehicle.id"
        class="flex w-full flex-col border-gray-200 py-2 not-last:border-b dark:border-neutral-700"
        :to="{
          name: 'map-marker',
          params: { markerType: vehicle.type, markerId: vehicle.id },
        }"
      >
        <div class="flex flex-row">
          <i-carbon-bicycle v-if="vehicle.type === 'bike'" class="mr-2" />
          <i-mdi-ferry v-else-if="vehicle.type === 'ferry'" class="mr-2" />

          <span class="mr-2">{{ vehicle.name }}</span>
        </div>
      </router-link>
      <NoData v-if="sortedDepartures && sortedDepartures.length === 0">
        {{ t('no_bus_wants_to_stop_here_right_now') }}
      </NoData>
      <i-ph-circle-notch v-if="!sortedDepartures && !stop.vehicles" class="m-auto animate-spin text-3xl" />
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
import type { Marker, StopDeparture } from '~/api/types';
import { computed, onBeforeUnmount, toRef } from 'vue';

import { useI18n } from 'vue-i18n';
import { api } from '~/api';
import Button from '~/components/atomic/Button.vue';
import NoData from '~/components/NoData.vue';
import { get24hTime } from '~/compositions/date';
import { useFavorites } from '~/compositions/useFavorites';

const props = defineProps<{
  marker: Marker;
}>();

const { addFavorite, removeFavorite, isFavorite } = useFavorites();
const { t } = useI18n();

const marker = toRef(props, 'marker');

const { stop, unsubscribe: unsubscribeStop } = api.useStop(computed(() => props.marker.id));

function getBoardTime(departure: StopDeparture) {
  if (!departure.actual) {
    return get24hTime(departure.planned);
  }

  if (departure.state === 'stopping') {
    return t('stopping');
  }

  const actual = new Date(departure.actual);
  const diffMinutes = Math.round((actual.getTime() - new Date().getTime()) / 60_000);

  if (diffMinutes < 1) {
    return t('immediately');
  }

  if (diffMinutes < 60) {
    return t('minutes', { minutes: diffMinutes });
  }

  return get24hTime(departure.actual);
}

const sortedDepartures = computed(() => {
  const departures = stop.value?.departures;
  if (!departures) {
    return null;
  }

  return departures
    .toSorted((a, b) => {
      const aTime = new Date(a.actual ?? a.planned).getTime();
      const bTime = new Date(b.actual ?? b.planned).getTime();
      return aTime - bTime;
    })

    .map((a) => {
      return {
        ...a,
      };
    });
});

onBeforeUnmount(async () => {
  await unsubscribeStop();
});
</script>
