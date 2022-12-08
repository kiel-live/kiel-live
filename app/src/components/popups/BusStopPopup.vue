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

      <template v-if="stop.arrivals">
        <template v-if="stop.arrivals.length > 0">
          <router-link
            v-for="arrival in stop.arrivals"
            :key="arrival.tripId"
            class="flex py-2 w-full not-last:border-b-1 dark:border-dark-300"
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

import { stops, subscribe, unsubscribe } from '~/api';
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
</script>
