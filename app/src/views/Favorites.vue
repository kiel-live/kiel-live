<template>
  <div class="flex h-full w-full min-h-0">
    <AppBar />

    <div class="flex flex-col mt-16 mx-2 w-full overflow-y-auto">
      <div v-if="favorites.length === 0" class="m-auto max-w-52 text-center text-xl">
        <p>Füge neue Haltestellen als Favoriten hinzu, indem du beim Öffnen der Haltestelle auf den Stern klickst.</p>
      </div>
      <router-link
        v-for="favorite in favorites"
        :key="favorite.id"
        :to="{ name: 'map-marker', params: { markerType: favorite.type, markerId: favorite.id } }"
        class="flex p-2 not-last:border-b-1 border-gray-600 dark:border-dark-800"
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

import AppBar from '~/components/AppBar.vue';
import { useFavorites } from '~/compositions/useFavorites';

export default defineComponent({
  // eslint-disable-next-line vue/multi-word-component-names
  name: 'Favorites',

  components: { AppBar },

  setup() {
    const { favorites } = useFavorites();

    return { favorites };
  },
});
</script>
