<template>
  <div class="flex flex-col m-auto w-full h-screen bg-white text-black dark:bg-dark-400 dark:text-gray-300">
    <main class="flex flex-grow min-h-0">
      <router-view />
    </main>
    <AppBarBottom class="mt-auto flex-shrink-0" />
    <UpdateNotice />
    <GeneralNotice />
  </div>
</template>

<script lang="ts" setup>
import { watch } from 'vue';
import { useI18n } from 'vue-i18n';

import AppBarBottom from '~/components/layout/AppBarBottom.vue';
import UpdateNotice from '~/components/UpdateNotice.vue';

import { useColorMode } from './compositions/useColorMode';

useColorMode();

const { locale } = useI18n();
watch(
  locale,
  () => {
    document.documentElement.setAttribute('lang', locale.value);
  },
  { immediate: true },
);
</script>

<!-- eslint-disable-next-line vue-scoped-css/enforce-style-type -->
<style>
html,
body,
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
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

.prose * {
  @apply dark:text-white;
}

.prose a {
  @apply text-blue-500 dark:text-white;
}
</style>
