<template>
  <div class="relative h-full w-full items-center justify-center overflow-hidden">
    <AppBar v-model:search-input="searchInput" />

    <DetailsPopup
      v-model:size="actualPopupSize"
      :is-open="!!selectedMarker"
      :disable-resize="liteMode"
      @close="selectedMarker = undefined"
    >
      <MarkerPopup v-if="selectedMarker" :marker="selectedMarker" />
    </DetailsPopup>

    <DetailsPopup
      v-model:size="actualPopupSize"
      :is-open="$route.name === 'search'"
      :disable-resize="liteMode"
      @close="$router.replace({ name: 'home' })"
    >
      <SearchPopup v-model:search-input="searchInput" />
    </DetailsPopup>

    <DetailsPopup
      v-model:size="actualPopupSize"
      :is-open="$route.name === 'favorites'"
      :disable-resize="liteMode"
      @close="$router.replace({ name: 'home' })"
    >
      <FavoritesPopup />
    </DetailsPopup>

    <Map
      v-if="!liteMode"
      :selected-marker="selectedMarker"
      @marker-click="selectedMarker = $event"
      @map-moved="popupSize = 'minimized'"
    />
  </div>
</template>

<script lang="ts" setup>
import type { Marker } from '~/api/types';
import { computed, ref, watch } from 'vue';

import { useRoute, useRouter } from 'vue-router';
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
    if (route.name === 'map-marker') {
      void router.replace({ name: 'map-marker', params: { markerType: marker.type, markerId: marker.id } });
    } else {
      void router.push({ name: 'map-marker', params: { markerType: marker.type, markerId: marker.id } });
    }
  },
});

const searchInput = ref('');

type PopupSize = 'full' | 'half' | 'minimized';

const popupSize = ref<PopupSize>('half');

const actualPopupSize = computed<PopupSize>({
  get() {
    if (liteMode.value) {
      return 'full';
    }
    return popupSize.value;
  },
  set(size) {
    popupSize.value = size;
  },
});

watch(selectedMarker, (marker) => {
  if (marker) {
    popupSize.value = 'half';
  }
});

watch(route, () => {
  if ((route.name === 'search' || route.name === 'favorites') && popupSize.value === 'minimized') {
    popupSize.value = 'half';
  }
});
</script>
