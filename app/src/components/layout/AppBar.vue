<template>
  <!-- TODO: remove id app-bar when android apps are updated to set --safe-area-top CSS variable -->
  <nav
    id="app-bar"
    class="absolute top-0 right-0 left-0 z-20 mx-2 mt-[calc(0.5rem+var(--safe-area-top,0px))] flex h-12 items-center justify-between gap-x-1 rounded-xl border border-gray-100 bg-white py-1 pr-1 shadow-lg md:right-auto md:left-1/2 md:w-96 md:-translate-x-1/2 md:transform dark:border-neutral-950 dark:bg-neutral-800 dark:text-gray-300"
  >
    <router-link :to="{ name: 'home' }" class="p-2">
      <img :alt="t('logo_alt')" src="../../assets/logo.png" class="h-6 w-6" />
    </router-link>
    <div v-if="isConnected" class="flex h-full grow">
      <input
        :value="internalSearchInput"
        type="text"
        class="focus-visible:border-opacity-50 h-full w-full border border-transparent bg-transparent p-2 focus:outline-none focus-visible:rounded-md focus-visible:border-gray-300 focus-visible:outline-none"
        :title="t('search')"
        :placeholder="`${t('search')} ...`"
        autofocus
        name="query"
        @input="(event) => (internalSearchInput = (event.currentTarget as HTMLInputElement).value)"
        @keydown.escape="$router.back()"
        @click="$router.push({ name: 'search' })"
      />
    </div>
    <div v-else class="mr-2 flex items-center gap-x-2">
      <span>{{ t('no_connection') }}</span>
      <i-carbon-cloud-offline class="text-red-600" />
    </div>
    <Button v-if="needRefresh" class="h-full gap-x-1" @click="updateServiceWorker(true)">
      <i-carbon-cloud-download />
      <span>{{ t('update') }}</span>
    </Button>
  </nav>
</template>

<script lang="ts" setup>
import { useRegisterSW } from 'virtual:pwa-register/vue';
import { computed, toRef } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';

import { api } from '~/api';
import Button from '~/components/atomic/Button.vue';

const props = defineProps<{
  searchInput: string;
}>();

const emit = defineEmits<{
  (e: 'update:search-input', searchInput: string): void;
}>();

const { isConnected } = api;

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const { needRefresh, updateServiceWorker } = useRegisterSW();

const searchInput = toRef(props, 'searchInput');
const internalSearchInput = computed({
  get() {
    return searchInput.value;
  },
  set(_searchInput: string) {
    emit('update:search-input', _searchInput);

    if (_searchInput.length > 0 && route.name !== 'search') {
      void router.push({ name: 'search' });
    }

    if (_searchInput.length === 0 && route.name === 'search') {
      void router.push({ name: 'home' });
    }
  },
});
</script>
