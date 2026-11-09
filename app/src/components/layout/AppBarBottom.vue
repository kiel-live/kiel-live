<template>
  <!-- Mobile Bottom Navigation -->
  <nav
    class="flex w-full justify-center border-t border-gray-100 bg-white md:hidden dark:border-neutral-950 dark:bg-neutral-800 dark:text-gray-300"
  >
    <div class="flex w-full max-w-96 justify-around">
      <router-link
        v-if="liteMode"
        :to="{ name: 'search' }"
        class="flex w-1/3 flex-col items-center p-2"
        :class="{ 'dark:text-red-high-contrast text-red-700': activeArea === 'map-or-search' }"
        :aria-label="t('search')"
      >
        <i-ph-magnifying-glass-bold class="mb-1 h-6 w-6" />
        <span class="mt-auto text-xs">{{ t('search') }}</span>
      </router-link>
      <router-link
        v-else
        :to="{ name: 'home' }"
        class="flex w-1/3 flex-col items-center p-2"
        :class="{ 'dark:text-red-high-contrast text-red-700': activeArea === 'map-or-search' }"
        :aria-label="t('map')"
      >
        <i-carbon-map class="mb-1 h-6 w-6" />
        <span class="mt-auto text-xs">{{ t('map') }}</span>
      </router-link>
      <router-link
        :to="{ name: 'favorites' }"
        class="flex w-1/3 flex-col items-center p-2"
        :class="{ 'dark:text-red-high-contrast text-red-700': activeArea === 'favorites' }"
        :aria-label="t('favorites')"
      >
        <i-ph-star-fill class="mb-1 h-6 w-6" />
        <span class="mt-auto text-xs">{{ t('favorites') }}</span>
      </router-link>
      <router-link
        :to="{ name: 'settings-about' }"
        class="flex w-1/3 flex-col items-center p-2"
        :class="{ 'dark:text-red-high-contrast text-red-700': activeArea === 'settings' }"
        :aria-label="t('settings')"
      >
        <i-ph-gear-fill class="mb-1 h-6 w-6" />
        <span class="mt-auto text-xs">{{ t('settings') }}</span>
      </router-link>
    </div>
  </nav>

  <!-- Desktop Left Sidebar Navigation -->
  <nav
    class="fixed top-0 bottom-0 left-0 z-30 hidden w-20 flex-col items-center gap-6 border-r border-gray-100 bg-white py-6 md:flex dark:border-neutral-950 dark:bg-neutral-800 dark:text-gray-300"
  >
    <router-link :to="{ name: 'home' }" class="mb-4" :aria-label="t('logo_alt')">
      <img :alt="t('logo_alt')" src="../../assets/logo.png" class="h-10 w-10" />
    </router-link>

    <router-link
      v-if="!liteMode"
      :to="{ name: 'home' }"
      class="flex flex-col items-center rounded-xl p-3 transition-all hover:bg-gray-100 dark:hover:bg-neutral-900"
      :class="{
        'dark:text-red-high-contrast bg-red-50 text-red-700 dark:bg-red-950/30': activeArea === 'map-or-search',
      }"
      :aria-label="t('map')"
    >
      <i-carbon-map class="h-7 w-7" />
    </router-link>

    <router-link
      :to="{ name: 'favorites' }"
      class="flex flex-col items-center rounded-xl p-3 transition-all hover:bg-gray-100 dark:hover:bg-neutral-900"
      :class="{ 'dark:text-red-high-contrast bg-red-50 text-red-700 dark:bg-red-950/30': activeArea === 'favorites' }"
      :aria-label="t('favorites')"
    >
      <i-ph-star-fill class="h-7 w-7" />
    </router-link>

    <router-link
      :to="{ name: 'settings-about' }"
      class="mt-auto flex flex-col items-center rounded-xl p-3 transition-all hover:bg-gray-100 dark:hover:bg-neutral-900"
      :class="{ 'dark:text-red-high-contrast bg-red-50 text-red-700 dark:bg-red-950/30': activeArea === 'settings' }"
      :aria-label="t('settings')"
    >
      <i-ph-gear-fill class="h-7 w-7" />
    </router-link>
  </nav>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute } from 'vue-router';

import { useUserSettings } from '~/compositions/useUserSettings';

const { t } = useI18n();
const { liteMode } = useUserSettings();

const route = useRoute();

const activeArea = computed(() => {
  if (route.name === 'favorites') {
    return 'favorites';
  }

  if (route.meta.settings) {
    return 'settings';
  }

  return 'map-or-search';
});
</script>
