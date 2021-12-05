<template>
  <div
    class="absolute top-0 left-0 right-0 mx-2 mt-2 h-12 bg-white flex rounded-md p-4 items-center justify-center border-1 border-gray-200 shadow-md md:transform md:-translate-x-1/2 md:right-auto md:left-1/2 md:w-64"
    :class="{ 'bg-red-300': !isConnected }"
  >
    <span class="text-xl font-bold">Kiel-Live</span>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

import { isConnected } from '~/api';

export default defineComponent({
  name: 'AppBar',

  setup() {
    return { isConnected };
  },
});
</script>
