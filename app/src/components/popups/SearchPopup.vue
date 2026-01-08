<template>
  <div class="flex flex-col min-h-0 grow">
    <div class="flex pb-2 mb-2 border-b border-gray-200 dark:border-neutral-600 space-x-2 items-center">
      <i-ph-magnifying-glass-bold />
      <h1 class="text-lg">{{ t('search_result') }}</h1>
    </div>
    <div v-if="searchResults.length === 0 && searchInput.length < 3" class="m-auto max-w-52 text-center text-xl">
      <p>{{ t('search_stop_vehicle') }}</p>
    </div>
    <div v-else-if="searchResults.length === 0 && searchInput.length >= 3" class="m-auto max-w-52 text-center text-xl">
      <p>{{ t('no_entry') }}</p>
    </div>
    <div class="flex flex-col overflow-y-auto">
      <router-link
        v-for="searchResult in searchResults"
        :key="searchResult.id"
        :to="{ name: 'map-marker', params: { markerType: searchResult.type, markerId: searchResult.id } }"
        class="flex py-2 not-last:border-b border-gray-200 dark:border-neutral-700 max-w-full"
        @click="searchInput = ''"
      >
        <i-mdi-sign-real-estate v-if="searchResult.type === 'bus-stop'" class="mr-2" />
        <i-mdi-ferry v-else-if="searchResult.type === 'ferry-stop'" class="mr-2" />
        <!-- <i-fa-bus v-else-if="searchResult.type === 'bus'" class="mr-2" /> -->
        <div class="">
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
