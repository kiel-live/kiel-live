<template>
  <nav
    class="flex w-full justify-center border-t-1 border-gray-200 dark:border-neutral-600 dark:bg-neutral-800 dark:text-gray-300 shadow-top"
  >
    <div class="flex w-full justify-around max-w-96">
      <router-link
        v-if="liteMode"
        :to="{ name: 'search' }"
        class="flex flex-col items-center w-1/3 p-2"
        :class="{ 'text-red-700 dark:text-red-500': activeArea === 'search' }"
        :aria-label="t('search')"
      >
        <i-ph-magnifying-glass-bold class="w-6 h-6 mb-1" />
        <span class="mt-auto text-xs">{{ t('search') }}</span>
      </router-link>
      <router-link
        v-else
        :to="{ name: 'home' }"
        class="flex flex-col items-center w-1/3 p-2"
        :class="{ 'text-red-700 dark:text-red-500': activeArea === 'map' }"
        :aria-label="t('map')"
      >
        <i-carbon-map class="w-6 h-6 mb-1" />
        <span class="mt-auto text-xs">{{ t('map') }}</span>
      </router-link>
      <router-link
        :to="{ name: 'favorites' }"
        class="flex flex-col items-center w-1/3 p-2"
        :class="{ 'text-red-700 dark:text-red-500': activeArea === 'favorites' }"
        :aria-label="t('favorites')"
      >
        <i-ph-star-fill class="w-6 h-6 mb-1" />
        <span class="mt-auto text-xs">{{ t('favorites') }}</span>
      </router-link>
      <router-link
        :to="{ name: 'settings-about' }"
        class="flex flex-col items-center w-1/3 p-2"
        :class="{ 'text-red-700 dark:text-red-500': activeArea === 'settings' }"
        :aria-label="t('settings')"
      >
        <i-ic-baseline-settings class="w-6 h-6 mb-1" />
        <span class="mt-auto text-xs">{{ t('settings') }}</span>
      </router-link>
    </div>
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

  if (route.name === 'search') {
    return 'search';
  }

  if (route.meta.settings) {
    return 'settings';
  }

  return 'map';
});
</script>
