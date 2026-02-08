<template>
  <div class="flex min-h-0 grow flex-col">
    <div class="mb-2 flex items-center space-x-2 border-b border-gray-200 pb-2 dark:border-neutral-600">
      <i-ph-star-fill />
      <h1 class="text-lg">{{ t('favorites') }}</h1>
    </div>
    <div v-if="favorites.length === 0" class="m-auto max-w-52 text-center text-xl">
      <p>{{ t('add_favorites') }}</p>
    </div>
    <div class="flex flex-col overflow-y-auto">
      <router-link
        v-for="favorite in favorites"
        :key="favorite.id"
        :to="{ name: 'map-marker', params: { markerType: favorite.type, markerId: favorite.id } }"
        class="flex border-gray-200 py-2 not-last:border-b dark:border-neutral-700"
      >
        <i-mdi-sign-real-estate v-if="favorite.type === 'bus-stop'" class="mr-2" />
        <div class="">
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
