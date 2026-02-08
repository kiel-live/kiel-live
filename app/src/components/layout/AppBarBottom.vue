<template>
  <nav
    class="shadow-top flex w-full justify-center border-t border-gray-200 dark:border-gray-800 dark:bg-neutral-800 dark:text-gray-300"
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
