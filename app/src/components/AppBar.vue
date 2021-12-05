<template>
  <div
    class="absolute top-0 left-0 right-0 mx-2 mt-2 h-12 flex rounded-full p-4 items-center justify-center bg-white border-1 border-gray-200 shadow-xl z-20 md:transform md:-translate-x-1/2 md:right-auto md:left-1/2 md:w-64 dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800"
    :class="{ 'bg-red-300': !isConnected }"
  >
    <img src="../assets/logo.png" class="w-6 h-6 mr-auto" />
    <input
      v-model="searchInput"
      type="text"
      class="bg-transparent focus:outline-transparent mx-2 w-full"
      placeholder="Suchen ..."
      @click="$router.push({ name: 'search' })"
    />
    <i-ph-star-fill />
  </div>
  <div v-if="searchResults.length > 0" class="bg-white w-full">
    <div v-for="result in searchResults" :key="result.refIndex">
      {{ result }}
    </div>
  </div>
</template>

<script lang="ts">
import Fuse from 'fuse.js';
import { computed, defineComponent, ref } from 'vue';

import { isConnected, stops } from '~/api';

export default defineComponent({
  name: 'AppBar',

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
      if (searchInput.value === '') {
        return [];
      }

      return searchIndex.value.search(searchInput.value);
    });

    return { isConnected, searchResults, searchInput };
  },
});
</script>
