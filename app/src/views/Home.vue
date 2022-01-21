<template>
  <div class="relative h-full w-full items-center justify-center overflow-hidden">
    <Map :selected-marker="selectedMarker" @marker-click="selectedMarker = $event" />
    <DetailsPopup :is-open="!!selectedMarker" size="3/4" @close="selectedMarker = undefined">
      <MarkerPopup v-if="selectedMarker" :marker="selectedMarker" />
    </DetailsPopup>

    <DetailsPopup :is-open="$route.name === 'search'" size="1/2" @close="$router.back()">
      <SearchPopup v-model:search-input="searchInput" />
    </DetailsPopup>

    <DetailsPopup :is-open="$route.name === 'favorites'" size="1/2" @close="$router.back()">
      <FavoritesPopup />
    </DetailsPopup>

    <AppBar v-model:search-input="searchInput" />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { Marker } from '~/api/types';
import DetailsPopup from '~/components/DetailsPopup.vue';
import AppBar from '~/components/layout/AppBar.vue';
import Map from '~/components/map/Map.vue';
import FavoritesPopup from '~/components/popups/FavoritesPopup.vue';
import MarkerPopup from '~/components/popups/MarkerPopup.vue';
import SearchPopup from '~/components/popups/SearchPopup.vue';

export default defineComponent({
  // eslint-disable-next-line vue/multi-word-component-names
  name: 'Home',

  components: { Map, DetailsPopup, AppBar, MarkerPopup, SearchPopup, FavoritesPopup },

  setup() {
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
      async set(marker) {
        if (!marker) {
          await router.replace({ name: 'home' });
          return;
        }
        await router.replace({ name: 'map-marker', params: { markerType: marker.type, markerId: marker.id } });
      },
    });

    const searchInput = ref('');

    return { selectedMarker, searchInput };
  },
});
</script>
