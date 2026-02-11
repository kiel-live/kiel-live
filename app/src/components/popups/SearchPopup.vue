<template>
  <div class="flex min-h-0 grow flex-col">
    <div class="mb-4 flex items-center space-x-2">
      <i-ph-magnifying-glass-bold class="text-xl" />
      <h1 class="text-xl font-semibold">{{ t('search_result') }}</h1>
    </div>
    <div v-if="searchResults.length === 0 && searchInput.length < 3" class="m-auto max-w-52 text-center text-lg">
      <p class="text-gray-500 dark:text-gray-400">{{ t('search_stop_vehicle') }}</p>
    </div>
    <div v-else-if="searchResults.length === 0 && searchInput.length >= 3" class="m-auto max-w-52 text-center text-lg">
      <p class="text-gray-500 dark:text-gray-400">{{ t('no_entry') }}</p>
    </div>
    <div class="flex flex-col gap-1 overflow-y-auto">
      <router-link
        v-for="searchResult in searchResults"
        :key="searchResult.id"
        :to="{ name: 'map-marker', params: { markerType: searchResult.type, markerId: searchResult.id } }"
        class="flex max-w-full items-center gap-3 rounded-lg border border-gray-100 bg-white px-4 py-3 transition-colors hover:bg-gray-50 dark:border-neutral-800 dark:bg-neutral-950 dark:hover:bg-neutral-900"
        @click="searchInput = ''"
      >
        <i-mdi-sign-real-estate v-if="searchResult.type === 'bus-stop'" class="text-xl" />
        <i-mdi-ferry v-else-if="searchResult.type === 'ferry-stop'" class="text-xl" />
        <div class="font-medium">
          {{ searchResult.name }}
        </div>
      </router-link>
    </div>
  </div>
</template>

<script lang="ts" setup>
import type { Bounds } from '~/api/types';
import { ref } from 'vue';

import { useI18n } from 'vue-i18n';
import { api } from '~/api';

const searchInput = defineModel<string>('searchInput', {
  default: '',
});

const { t } = useI18n();

// TODO: use proper bounds / server search
const bounds = ref<Bounds>({
  east: 0,
  west: 0,
  north: 0,
  south: 0,
});

const { results: searchResults } = api.useSearch(searchInput, bounds);
</script>
