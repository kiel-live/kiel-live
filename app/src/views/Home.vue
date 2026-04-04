<template>
  <div class="relative h-full w-full items-center justify-center overflow-hidden">
    <AppBar v-model:search-input="searchInput" />

    <DetailsPopup :is-open="isPopupOpen" v-model:snap-point="snapPoint" @close="$router.replace({ name: 'home' })">
      <MarkerPopup v-if="selectedMarker" :marker="selectedMarker" />
      <SearchPopup v-if="route.name === 'search'" v-model:search-input="searchInput" />
      <FavoritesPopup v-if="route.name === 'favorites'" />
    </DetailsPopup>

    <Map
      v-if="!liteMode"
      v-model:map-moved-manually="mapMovedManually"
      :selected-marker="selectedMarker"
      @marker-click="onMarkerClick"
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

const snapPoint = ref<'header' | 'half' | 'full'>('half');

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

function onMarkerClick(marker: Marker) {
  snapPoint.value = 'half';
  selectedMarker.value = marker;
}

const searchInput = ref('');
const mapMovedManually = ref(false);

// Shrink to header snap when the user manually pans the map
watch(mapMovedManually, (moved) => {
  if (moved && isPopupOpen.value) {
    snapPoint.value = 'header';
  }
});

const isPopupOpen = computed(() => {
  return selectedMarker.value !== undefined || route.name === 'search' || route.name === 'favorites';
});
</script>
