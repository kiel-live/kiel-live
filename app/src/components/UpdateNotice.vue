<template>
  <div
    v-if="show"
    class="flex fixed top-0 left-0 right-0 bottom-0 bg-gray-900 bg-opacity-80 z-1000 justify-center items-center"
    @click="close"
  >
    <div
      class="m-2 flex flex-col rounded-md p-4 bg-white border-1 border-gray-200 shadow-xl z-20 md:w-104 dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800"
    >
      <div class="flex flex-col text-center gap-4">
        <span class="text-xl mb-2">{{ t('update_title') }}</span>
        <span>{{ t('update_msg') }}</span>
        <span
          >{{ t('feedback') }}<a :href="`mailto:${feedbackMail}`" class="underline">{{ feedbackMail }}</a>
          {{ t('follow_us') }}
          <a href="https://www.instagram.com/kiel.live/" target="_blank" class="underline">Instagram</a>.</span
        >
      </div>

      <div class="flex flex-row w-full justify-center mt-8">
        <Button @click="close">Nice!</Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import confetti from 'canvas-confetti';
import { computed, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import Button from '~/components/atomic/Button.vue';

const LS_VERSION_KEY = 'kiel-live-version-v1';
const latestVersion = '2.0.0';
const feedbackMail = atob('YW5kcm9pZEBqdTYwLmRl'); // email as base64

const version = ref(localStorage.getItem(LS_VERSION_KEY));
const show = computed(() => version.value !== null && version.value !== latestVersion);
const { t } = useI18n();

onMounted(async () => {
  const oldVersion = localStorage.getItem('version');
  if (oldVersion !== null) {
    version.value = oldVersion;
    localStorage.removeItem('version');
  }
});

watch(
  show,
  () => {
    if (!show.value) {
      return;
    }

    const duration = 3 * 1000;
    const animationEnd = Date.now() + duration;

    const interval = setInterval(() => {
      const timeLeft = animationEnd - Date.now();

      if (timeLeft <= 0) {
        clearInterval(interval);
        return;
      }

      void confetti({ particleCount: 100, spread: 70, origin: { y: 1.1 }, startVelocity: 90, zIndex: 2000 });
    }, 250);
  },
  { immediate: true },
);

function close() {
  localStorage.setItem(LS_VERSION_KEY, latestVersion);
  version.value = latestVersion;
}
</script>
