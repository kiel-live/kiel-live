<template>
  <div class="flex min-h-0 grow flex-col">
    <div class="mb-4 flex items-center space-x-2">
      <i-ph-star-fill class="text-xl" />
      <h1 class="text-xl font-semibold">{{ t('favorites') }}</h1>
    </div>
    <div v-if="favorites.length === 0" class="m-auto max-w-52 text-center text-lg">
      <p class="text-gray-500 dark:text-gray-400">{{ t('add_favorites') }}</p>
    </div>
    <div class="flex flex-col gap-1 overflow-y-auto">
      <router-link
        v-for="favorite in favorites"
        :key="favorite.id"
        :to="{ name: 'map-marker', params: { markerType: favorite.type, markerId: favorite.id } }"
        class="flex items-center gap-3 rounded-lg border border-gray-100 bg-white px-4 py-3 transition-colors hover:bg-gray-50 dark:border-neutral-950 dark:bg-neutral-800 dark:hover:bg-neutral-900"
      >
        <i-mdi-sign-real-estate v-if="favorite.type === 'bus-stop'" class="text-xl" />
        <div class="font-medium">
          {{ favorite.name }}
        </div>
      </router-link>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import { useFavorites } from '~/compositions/useFavorites';

const { t } = useI18n();
const { favorites } = useFavorites();
</script>
