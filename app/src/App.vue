<template>
  <div class="app flex flex-col m-auto w-full h-full bg-gray-100 dark:bg-dark-400 dark:text-gray-300">
    <div class="flex flex-grow min-h-0">
      <router-view />
    </div>
    <AppBarBottom class="mt-auto flex-shrink-0" />
    <ReloadPrompt />
    <UpdateNotice />
  </div>
</template>

<script lang="ts" setup>
import AppBarBottom from '~/components/layout/AppBarBottom.vue';
import ReloadPrompt from '~/components/ReloadPrompt.vue';
import UpdateNotice from '~/components/UpdateNotice.vue';
import { usePrefersColorSchemeDark } from '~/compositions/usePrefersColorScheme';

usePrefersColorSchemeDark();
</script>

<!-- eslint-disable-next-line vue-scoped-css/enforce-style-type -->
<style>
html,
body,
#app {
  width: 100%;
  height: 100%;
}

body {
  /* disable android pull to refresh feature */
  overflow-y: hidden;
}
*::-webkit-scrollbar {
  @apply bg-transparent w-12px h-12px;
}

* {
  scrollbar-width: thin;
}

*::-webkit-scrollbar-thumb {
  transition: background 0.2s ease-in-out;
  border: 3px solid transparent;
  @apply bg-cool-gray-500 dark:bg-dark-200 rounded-full bg-clip-content;
}

*::-webkit-scrollbar-thumb:hover {
  @apply bg-cool-gray-600 dark:bg-dark-100;
}

*::-webkit-scrollbar-corner {
  @apply bg-transparent;
}
</style>

<style scoped>
.app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
