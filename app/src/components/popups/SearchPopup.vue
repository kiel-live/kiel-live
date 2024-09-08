<template>
  <div class="flex flex-col min-h-0 flex-grow">
    <div class="flex pb-2 mb-2 border-b-1 dark:border-dark-100 space-x-2 items-center">
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
        :key="searchResult.refIndex"
        :to="{ name: 'map-marker', params: { markerType: searchResult.item.type, markerId: searchResult.item.id } }"
        class="flex py-2 not-last:border-b-1 dark:border-dark-300 max-w-full"
        @click="$emit('update:search-input', '')"
      >
        <i-mdi-sign-real-estate v-if="searchResult.item.type === 'bus-stop'" class="mr-2" />
        <i-mdi-ferry v-else-if="searchResult.item.type === 'ferry-stop'" class="mr-2" />
        <!-- <i-fa-bus v-else-if="searchResult.item.type === 'bus'" class="mr-2" /> -->
        <div class="">
          {{ searchResult.item.name }}
        </div>
      </router-link>
    </div>
  </div>
</template>

<script lang="ts">
import Fuse from 'fuse.js';
import { computed, defineComponent, onMounted, toRef } from 'vue';
import { useI18n } from 'vue-i18n';

import { stops, subscribe, vehicles } from '~/api';

export default defineComponent({
  name: 'SearchPopup',

  props: {
    searchInput: {
      type: String,
      default: '',
    },
  },

  emits: {
    'update:search-input': (_searchInput: string) => true,
  },

  setup(props) {
    const { t } = useI18n();

    const searchInput = toRef(props, 'searchInput');
    const searchData = computed(() => [...Object.values(stops.value)]);
    const searchIndex = computed(
      () =>
        new Fuse(searchData.value, {
          includeScore: true,
          keys: ['name'],
          threshold: 0.4,
        }),
    );

    const searchResults = computed(() => {
      if (searchInput.value === '' || searchInput.value.length < 3) {
        return [];
      }
      // limit to max 20 results
      return searchIndex.value.search(searchInput.value).slice(0, 20);
    });

    onMounted(async () => {
      void subscribe('data.map.vehicle.>', vehicles);
      void subscribe('data.map.stop.>', stops);
    });

    return { t, searchResults };
  },
});
</script>
