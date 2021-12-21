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
      autofocus
    />
    <div class="flex items-center cursor-pointer">
      <i-mdi-close v-if="$route.name === 'search' || $route.name === 'favorites'" @click="$router.back()" />
      <i-ph-star-fill v-else @click="$router.push({ name: 'favorites' })" />
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { isConnected } from '~/api';

export default defineComponent({
  name: 'AppBar',

  emits: {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    'search-input': (_searchInput: string) => true,
  },

  setup(_, { emit }) {
    const route = useRoute();
    const router = useRouter();
    const searchRaw = ref('');
    const searchInput = computed({
      get() {
        return searchRaw.value;
      },
      set(_searchInput: string) {
        searchRaw.value = _searchInput;

        emit('search-input', _searchInput);

        if (_searchInput.length > 0 && route.name !== 'search') {
          void router.push({ name: 'search' });
        }
      },
    });

    return { isConnected, searchInput };
  },
});
</script>
