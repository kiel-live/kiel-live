<template>
  <div
    class="flex w-full justify-around border-t-1 border-gray-200 dark:border-gray-700 bg-gray-100 dark:bg-dark-400 dark:text-gray-300 drop-shadow"
  >
    <router-link
      :to="{ name: 'home' }"
      class="flex flex-col items-center w-1/3 p-2"
      :class="{ 'text-blue-400': $route.name !== 'favorites' && $route.name !== 'settings' }"
      :aria-label="t('map')"
    >
      <i-carbon-map class="w-6 h-6 mb-1" />
      <span class="mt-auto text-xs">{{ t('map') }}</span>
    </router-link>
    <router-link
      :to="{ name: 'favorites' }"
      class="flex flex-col items-center w-1/3 p-2"
      :class="{ 'text-blue-400': $route.name === 'favorites' }"
      :aria-label="t('favorites')"
    >
      <i-ph-star-fill class="w-6 h-6 mb-1" />
      <span class="mt-auto text-xs">{{ t('favorites') }}</span>
    </router-link>
    <router-link
      :to="{ name: 'settings' }"
      class="flex flex-col items-center w-1/3 p-2"
      :class="{ 'text-blue-400': $route.name === 'settings' }"
      :aria-label="t('settings')"
    >
      <i-ic-baseline-settings class="w-6 h-6 mb-1" />
      <span class="mt-auto text-xs">{{ t('settings') }}</span>
    </router-link>
  </div>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>
