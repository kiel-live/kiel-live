<template>
  <div
    class="absolute top-0 left-0 right-0 mx-2 mt-2 h-12 flex rounded-md py-4 items-center justify-center bg-white border-1 border-gray-200 shadow-xl z-20 md:transform md:-translate-x-1/2 md:right-auto md:left-1/2 md:w-96 dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800"
    :class="{ 'bg-red-300': !isConnected }"
  >
    <router-link :to="{ name: 'home' }" class="p-2">
      <img :alt="t('logo_alt')" src="../../assets/logo.png" class="w-6 h-6" />
    </router-link>
    <div class="flex flex-grow mr-2">
      <input
        :value="internalSearchInput"
        type="text"
        class="bg-transparent p-2 border border-transparent focus:outline-none focus-visible:(outline-none rounded-md border-gray-300 border-opacity-50) w-full h-full"
        :placeholder="t('search')"
        autofocus
        @input="(event) => (internalSearchInput = (event.currentTarget as HTMLInputElement).value)"
        @keydown.escape="$router.back()"
        @click="$router.push({ name: 'search' })"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, toRef } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';

import { isConnected } from '~/api';

export default defineComponent({
  name: 'AppBar',

  props: {
    searchInput: {
      type: String,
      required: true,
    },
  },

  emits: {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    'update:search-input': (_searchInput: string) => true,
  },

  setup(props, { emit }) {
    const { t } = useI18n();
    const route = useRoute();
    const router = useRouter();

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

    return { t, isConnected, internalSearchInput };
  },
});
</script>
