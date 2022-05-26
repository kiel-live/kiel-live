<template>
  <div class="flex flex-col min-h-0 flex-grow">
    <div class="flex pb-2 mb-2 border-b-1 dark:border-dark-100 space-x-2 items-center">
      <i-ph-star-fill />
      <span class="text-lg">Favoriten</span>
    </div>
    <div v-if="favorites.length === 0" class="m-auto max-w-52 text-center text-xl">
      <p>Füge neue Haltestellen als Favoriten hinzu, indem du beim Öffnen der Haltestelle auf den Stern klickst.</p>
    </div>
    <div class="flex flex-col overflow-y-auto">
      <router-link
        v-for="favorite in favorites"
        :key="favorite.id"
        :to="{ name: 'map-marker', params: { markerType: favorite.type, markerId: favorite.id } }"
        class="flex py-2 not-last:border-b-1 dark:border-dark-300"
      >
        <i-mdi-sign-real-estate v-if="favorite.type === 'bus-stop'" class="mr-2" />
        <!-- <i-fa-bus v-else-if="searchResult.item.type === 'bus'" class="mr-2" /> -->
        <div class="">
          {{ favorite.name }}
        </div>
      </router-link>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

import { useFavorites } from '~/compositions/useFavorites';

export default defineComponent({
  name: 'FavoritesPopup',

  setup() {
    const { favorites } = useFavorites();

    return { favorites };
  },
});
</script>
