<template>
  <div v-if="stop" class="flex flex-col min-h-0 flex-grow">
    <div class="flex flex-row pb-2 mb-2 border-b-1 dark:border-dark-100 items-center">
      <i-mdi-sign-real-estate v-if="stop.type === 'bus-stop'" />
      <span class="text-lg ml-2">{{ stop.name }}</span>
      <div class="flex ml-auto items-center cursor-pointer select-none">
        <i-ph-star-fill v-if="isFavorite(stop)" class="text-yellow-300" @click="removeFavorite(stop)" />
        <i-ph-star-bold v-else @click="addFavorite(stop)" />
      </div>
    </div>

    <div class="flex flex-col flex-grow overflow-y-auto">
      <div v-if="stop.alerts && stop.alerts.length >= 1" class="bg-red-600 bg-opacity-50 p-2 mb-2 rounded-md">
        <div class="flex items-center border-b-1 mb-4">
          <i-mdi-alert class="mr-2" /><span class="font-bold">Hinweise</span>
        </div>
        <div v-for="(alert, i) in stop.alerts" :key="i" class="flex items-center">{{ alert }}</div>
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
        <NoData v-else>{{ t('no_man_wants_to_stop_here_right_now') }}</NoData>
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
      <i-ph-star-fill class="mr-2 text-yellow-300" /><span>Favorit löschen</span>
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
</script>
