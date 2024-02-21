<template>
  <PopupNotice :show="show" @close="close">
    <div class="flex flex-col text-center gap-4">
      <span class="text-xl mb-2">{{ t('update_title') }}</span>
      <span>{{ t('update_msg') }}</span>
      <i18n-t keypath="feedback" tag="span">
        <template #email>
          <a :href="`mailto:${feedbackMail}`" class="underline">{{ feedbackMail }}</a>
        </template>
        <template #instagram>
          <a href="https://www.instagram.com/kiel.live/" target="_blank" class="underline">{{ t('instagram') }}</a>
        </template>
      </i18n-t>
    </div>

    <div class="flex flex-row w-full justify-center mt-8">
      <Button @click="close">{{ t('nice') }}</Button>
    </div>
  </PopupNotice>
</template>

<script setup lang="ts">
import { useStorage } from '@vueuse/core';
import confetti from 'canvas-confetti';
import { computed, onMounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import Button from '~/components/atomic/Button.vue';
import PopupNotice from '~/components/PopupNotice.vue';
import { localStoragePrefix } from '~/compositions/useUserSettings';
import { feedbackMail } from '~/config';

const latestVersion = '2.0.0';
const version = useStorage(`${localStoragePrefix}.version`, '2.0.0');
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
  version.value = latestVersion;
}
</script>
