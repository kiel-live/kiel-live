<template>
  <div
    class="absolute top-0 left-0 right-0 mx-2 mt-2 h-12 flex rounded-full p-4 items-center justify-center bg-white border-1 border-gray-200 shadow-xl z-20 md:transform md:-translate-x-1/2 md:right-auto md:left-1/2 md:w-64 dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800"
    :class="{ 'bg-red-300': !isConnected }"
  >
    <img src="../assets/logo.png" class="w-6 h-6 mr-auto" />
    <input
      type="text"
      class="bg-transparent focus:outline-transparent mx-2 w-full"
      placeholder="Suchen ..."
      :autofocus="$route.name === 'search'"
      :value="searchInput"
      @input="$emit('update:search-input', ($event.target as HTMLInputElement).value)"
      @click="$router.push({ name: 'search' })"
    />
    <i-ph-star-fill />
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue';

import { isConnected } from '~/api';

export default defineComponent({
  name: 'AppBar',

  props: {
    searchInput: {
      type: String as PropType<string>,
      default: '',
    },
  },

  emits: {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    'update:search-input': (_searchInput: string) => true,
  },

  setup() {
    return { isConnected };
  },
});
</script>
