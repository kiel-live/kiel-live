<template>
  <div class="flex h-full w-full min-h-0">
    <AppBar v-model:search-input="searchInput" />

    <div class="flex flex-col mt-16 mx-2 w-full overflow-y-auto">
      <div v-if="searchResults.length === 0" class="m-auto max-w-52 text-center text-xl">
        <p>Suche nach einer Haltestelle oder einem Fahrzeug</p>
      </div>
      <router-link
        v-for="searchResult in searchResults"
        :key="searchResult.refIndex"
        :to="{ name: 'map-marker', params: { markerType: searchResult.item.type, markerId: searchResult.item.id } }"
        class="flex p-2 not-last:border-b-1 border-gray-600 dark:border-dark-800"
      >
        <i-mdi-sign-real-estate v-if="searchResult.item.type === 'bus-stop'" class="mr-2" />
        <!-- <i-fa-bus v-else-if="searchResult.item.type === 'bus'" class="mr-2" /> -->
        <div class="">
          {{ searchResult.item.name }}
        </div>
      </router-link>
    </div>
  </div>
</template>

<script lang="ts">
import Fuse from 'fuse.js';
import { computed, defineComponent, onMounted, ref } from 'vue';

import { stops, subscribe, vehicles } from '~/api';
import AppBar from '~/components/AppBar.vue';

export default defineComponent({
  // eslint-disable-next-line vue/multi-word-component-names
  name: 'Search',

  components: { AppBar },

  setup() {
    const searchInput = ref('');

    const searchData = computed(() => [...Object.values(stops.value)]);
    const searchIndex = computed(
      () =>
        new Fuse(searchData.value, {
          includeScore: true,
          keys: ['name'],
          threshold: 0.4,
        }),
    );

    const searchResults = computed(() => {
      if (searchInput.value === '' || searchInput.value.length < 3) {
        return [];
      }

      // limit to max 20 results
      return searchIndex.value.search(searchInput.value).slice(0, 20);
    });

    onMounted(async () => {
      await subscribe('data.map.vehicle.>', vehicles);
      await subscribe('data.map.stop.>', stops);
    });

    return { searchInput, searchResults };
  },
});
</script>
