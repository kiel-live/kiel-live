<template>
  <div class="flex h-full w-full min-h-0">
    <AppBar v-model:search-input="searchInput" />

    <div class="flex flex-col mt-16 mx-2 w-full overflow-y-auto">
      <router-link
        v-for="searchResult in searchResults"
        :key="searchResult.refIndex"
        :to="{ name: 'map-marker', params: { markerType: searchResult.item.type, markerId: searchResult.item.id } }"
        class="flex p-2 not-last:border-b-1 border-gray-600 dark:border-dark-800"
      >
        <i-mdi-sign-real-estate v-if="searchResult.item.type === 'bus-stop'" class="mr-2" />
        <i-fa-bus v-else-if="searchResult.item.type === 'bus'" class="mr-2" />
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

    const searchData = computed(() => [...Object.values(stops.value), ...Object.values(vehicles.value)]);
    const searchIndex = computed(
      () =>
        new Fuse(searchData.value, {
          includeScore: true,
          keys: ['name'],
          threshold: 0.4,
        }),
    );

    const searchResults = computed(() => {
      if (searchInput.value === '') {
        return [];
      }

      return searchIndex.value.search(searchInput.value);
    });

    onMounted(async () => {
      await subscribe('data.map.vehicle.>', vehicles);
      await subscribe('data.map.stop.>', stops);

      setTimeout(() => {
        searchInput.value = 'Bahn';
      }, 100);
    });

    return { searchInput, searchResults };
  },
});
</script>
