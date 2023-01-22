<template>
  <div
    v-if="needRefresh"
    class="fixed left-0 right-0 bottom-0 mx-2 mb-2 flex flex-col rounded-md p-4 items-center justify-center bg-white border-1 border-gray-200 shadow-xl z-20 md:transform md:-translate-x-1/2 md:right-auto md:left-1/2 md:w-96 dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800"
    role="alert"
  >
    <div class="mb-2 text-xl">
      <span>{{ t('new_version') }}</span>
    </div>

    <div class="flex flex-row w-full gap-x-4 justify-center">
      <Button @click="updateServiceWorker(true)">{{ t('update') }}</Button>
      <Button @click="close">{{ t('cancel') }}</Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRegisterSW } from 'virtual:pwa-register/vue';
import { useI18n } from 'vue-i18n';

import Button from '~/components/atomic/Button.vue';

const { needRefresh, updateServiceWorker } = useRegisterSW();
const { t } = useI18n();

const close = async () => {
  needRefresh.value = false;
};
</script>
