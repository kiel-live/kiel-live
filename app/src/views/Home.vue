<template>
  <div class="relative h-full w-full items-center justify-center overflow-hidden">
    <AppBar v-model:search-input="searchInput" />

    <DetailsPopup
      :is-open="!!selectedMarker"
      :disable-resize="liteMode"
      :size="popupSize"
      @close="selectedMarker = undefined"
    >
      <MarkerPopup v-if="selectedMarker" :marker="selectedMarker" />
    </DetailsPopup>

    <DetailsPopup
      :is-open="$route.name === 'search'"
      :disable-resize="liteMode"
      :size="popupSize"
      @close="$router.replace({ name: 'home' })"
    >
      <SearchPopup v-model:search-input="searchInput" />
    </DetailsPopup>

    <DetailsPopup
      :is-open="$route.name === 'favorites'"
      :disable-resize="liteMode"
      :size="popupSize"
      @close="$router.replace({ name: 'home' })"
    >
      <FavoritesPopup />
    </DetailsPopup>

    <Map
      v-if="!liteMode"
      v-model:map-moved-manually="mapMovedManually"
      :selected-marker="selectedMarker"
      @marker-click="selectedMarker = $event"
    />
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import type { Marker } from '~/api/types';
import DetailsPopup from '~/components/DetailsPopup.vue';
import AppBar from '~/components/layout/AppBar.vue';
import Map from '~/components/map/Map.vue';
import FavoritesPopup from '~/components/popups/FavoritesPopup.vue';
import MarkerPopup from '~/components/popups/MarkerPopup.vue';
import SearchPopup from '~/components/popups/SearchPopup.vue';
import { useUserSettings } from '~/compositions/useUserSettings';

const { liteMode } = useUserSettings();
const route = useRoute();
const router = useRouter();
const selectedMarker = computed<Marker | undefined>({
  get() {
    if (route.name !== 'map-marker') {
      return undefined;
    }
    return {
      type: route.params.markerType,
      id: route.params.markerId,
    } as Marker;
  },
  set(marker) {
    if (!marker) {
      void router.replace({ name: 'home' });
      return;
    }
    void router.replace({ name: 'map-marker', params: { markerType: marker.type, markerId: marker.id } });
  },
});

const searchInput = ref('');

const mapMovedManually = ref(false);
const popupSize = computed(() => {
  if (liteMode.value) {
    return '1';
  }
  if (route.name === 'search' || route.name === 'favorites' || mapMovedManually.value) {
    return '1/2';
  }
  return '3/4';
});
</script>
