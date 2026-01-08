<template>
  <div class="app flex flex-col m-auto w-full h-full bg-white text-black dark:bg-neutral-800 dark:text-gray-300">
    <main class="flex grow min-h-0">
      <router-view />
    </main>
    <AppBarBottom class="mt-auto shrink-0" />
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

<style scoped>
.app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
