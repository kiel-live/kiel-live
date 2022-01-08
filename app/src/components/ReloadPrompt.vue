<template>
  <div
    v-if="offlineReady || needRefresh"
    class="fixed left-0 right-0 bottom-0 mx-2 mb-2 flex flex-col rounded-md p-4 items-center justify-center bg-white border-1 border-gray-200 shadow-xl z-20 dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800"
    role="alert"
  >
    <div class="mb-2 text-xl">
      <span v-if="offlineReady">Die App wurde fertig installiert.</span>
      <span v-if="needRefresh">Eine neue Version ist verfügbar.</span>
    </div>

    <div class="flex flex-row w-full gap-x-4 justify-center">
      <template v-if="needRefresh">
        <Button type="button" @click="updateServiceWorker(true)">Installieren</Button>
        <Button type="button" @click="close">Abbrechen</Button>
      </template>
      <Button v-else type="button" @click="close">Schließen</Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRegisterSW } from 'virtual:pwa-register/vue';

import Button from '~/components/atomic/Button.vue';

const { offlineReady, needRefresh, updateServiceWorker } = useRegisterSW();

const close = async () => {
  offlineReady.value = false;
  needRefresh.value = false;
};
</script>
