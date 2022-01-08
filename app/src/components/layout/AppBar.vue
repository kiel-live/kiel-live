<template>
  <div
    class="absolute top-0 left-0 right-0 mx-2 mt-2 h-12 flex rounded-md py-4 items-center justify-center bg-white border-1 border-gray-200 shadow-xl z-20 md:transform md:-translate-x-1/2 md:right-auto md:left-1/2 md:w-96 dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800"
    :class="{ 'bg-red-300': !isConnected }"
  >
    <div class="flex items-center cursor-pointer select-none p-2">
      <i-ic-baseline-arrow-back v-if="$route.name === 'search'" class="w-6 h-6" @click="$router.back()" />
      <router-link v-else :to="{ name: 'home' }">
        <img src="../../assets/logo.png" class="w-6 h-6" />
      </router-link>
    </div>
    <div class="flex-grow">
      <span v-if="$route.name === 'favorites'">Favorites</span>
      <input
        v-else
        v-model="internalSearchInput"
        type="text"
        class="bg-transparent focus:outline-transparent w-full h-full"
        placeholder="Suchen ..."
        autofocus
        @click="$router.push({ name: 'search' })"
      />
    </div>
    <div class="flex items-center cursor-pointer select-none w-10 p-2 items-center justify-center">
      <i-ph-star-fill
        v-if="$route.name !== 'favorites' && $route.name !== 'search'"
        @click="$router.push({ name: 'favorites' })"
      />
      <i-uil-times v-else-if="$route.name === 'favorites'" class="w-6 h-6" @click="$router.back()" />
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, toRef } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { isConnected } from '~/api';

export default defineComponent({
  name: 'AppBar',

  props: {
    searchInput: {
      type: String,
      required: true,
    },
  },

  emits: {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    'update:search-input': (_searchInput: string) => true,
  },

  setup(props, { emit }) {
    const route = useRoute();
    const router = useRouter();

    const searchInput = toRef(props, 'searchInput');
    const internalSearchInput = computed({
      get() {
        return searchInput.value;
      },
      set(_searchInput: string) {
        searchInput.value = _searchInput;

        emit('update:search-input', _searchInput);

        if (_searchInput.length > 0 && route.name !== 'search') {
          void router.push({ name: 'search' });
        }

        if (_searchInput.length === 0 && route.name === 'search') {
          void router.push({ name: 'home' });
        }
      },
    });

    return { isConnected, internalSearchInput };
  },
});
</script>
