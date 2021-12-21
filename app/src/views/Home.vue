<template>
  <div class="relative h-full w-full items-center justify-center overflow-hidden">
    <Map :selected-marker="selectedMarker" @marker-click="selectedMarker = $event" />
    <DetailsPopup :is-open="!!selectedMarker" @close="selectedMarker = undefined">
      <MarkerPopup v-if="selectedMarker" :marker="selectedMarker" />
    </DetailsPopup>
    <AppBar />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { Marker } from '~/api/types';
import AppBar from '~/components/AppBar.vue';
import DetailsPopup from '~/components/DetailsPopup.vue';
import Map from '~/components/Map.vue';
import MarkerPopup from '~/components/popups/MarkerPopup.vue';

export default defineComponent({
  // eslint-disable-next-line vue/multi-word-component-names
  name: 'Home',

  components: { Map, DetailsPopup, AppBar, MarkerPopup },

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

    return { selectedMarker };
  },
});
</script>
