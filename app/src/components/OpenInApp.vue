<template>
  <div
    v-if="show"
    class="fixed left-0 right-0 bottom-0 mx-2 mb-2 flex flex-col rounded-md p-4 items-center justify-center bg-white border-1 border-gray-200 shadow-xl z-20 md:transform md:-translate-x-1/2 md:right-auto md:left-1/2 md:w-96 dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800"
    role="alert"
  >
    Hey Android!
    <i-maki-cross @click="close" />
  </div>
</template>

<script setup lang="ts">
import UAParser from 'ua-parser-js';
import { ref } from 'vue';

const LS_HIDE_APP_HINT = 'kiel-live-hide-app-hint';

const uaParser = new UAParser();
const os = uaParser.getOS().name?.toLowerCase();
const browser = uaParser.getBrowser().name?.toLowerCase();
const isApp = !browser?.includes('webview');
const hideAppHint = localStorage.getItem(LS_HIDE_APP_HINT) === 'true';
const show = ref(!hideAppHint && os === 'android' && !isApp);

function close() {
  localStorage.setItem(LS_HIDE_APP_HINT, 'true');
  show.value = false;
}
</script>
