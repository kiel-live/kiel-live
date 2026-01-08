<template>
  <nav
    id="app-bar"
    class="absolute top-0 left-0 right-0 mx-2 mt-2 h-12 flex rounded-md py-1 pr-1 gap-x-1 items-center justify-between bg-white border border-gray-200 shadow-xl z-20 md:transform md:-translate-x-1/2 md:right-auto md:left-1/2 md:w-96 dark:bg-neutral-800 dark:text-gray-300 dark:border-neutral-950"
  >
    <router-link :to="{ name: 'home' }" class="p-2">
      <img :alt="t('logo_alt')" src="../../assets/logo.png" class="w-6 h-6" />
    </router-link>
    <div v-if="isConnected" class="flex flex-grow h-full">
      <input
        :value="internalSearchInput"
        type="text"
        class="bg-transparent p-2 border border-transparent focus:outline-none focus-visible:outline-none focus-visible:rounded-md focus-visible:border-gray-300 focus-visible:border-opacity-50 w-full h-full"
        :title="t('search')"
        :placeholder="`${t('search')} ...`"
        autofocus
        name="query"
        @input="(event) => (internalSearchInput = (event.currentTarget as HTMLInputElement).value)"
        @keydown.escape="$router.back()"
        @click="$router.push({ name: 'search' })"
      />
    </div>
    <div v-else class="flex gap-x-2 mr-2 items-center">
      <span>{{ t('no_connection') }}</span>
      <i-ic-baseline-cloud-off class="text-red-600" />
    </div>
    <Button v-if="needRefresh" class="h-full gap-x-1" @click="updateServiceWorker(true)">
      <i-ic-outline-cloud-download />
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
