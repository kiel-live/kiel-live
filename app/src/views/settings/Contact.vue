<template>
  <SettingsContainer>
    <h1 class="mb-6 text-2xl font-bold">{{ t('kiel_live') }}</h1>

    <img src="../../assets/logo.png" :alt="t('logo_alt')" class="mx-auto mb-6 w-32" />

    <p class="mb-6 text-center text-gray-600 dark:text-gray-400">{{ t('contact_title') }}</p>

    <form class="mx-auto flex w-full max-w-lg flex-col items-center gap-4" @submit.prevent="sendEmail">
      <textarea
        v-model="message"
        rows="10"
        class="w-full rounded-xl border border-gray-100 bg-white p-4 transition-colors focus:border-gray-300 focus:outline-none dark:border-neutral-800 dark:bg-neutral-950 dark:focus:border-neutral-700"
        type="text"
      />

      <Button type="submit" class="w-full">
        <i-ph-envelope class="mr-2" />
        <span>{{ t('send_email') }}</span>
      </Button>
    </form>
  </SettingsContainer>
</template>

<script lang="ts" setup>
import { useStorage } from '@vueuse/core';
import { useI18n } from 'vue-i18n';

import Button from '~/components/atomic/Button.vue';
import SettingsContainer from '~/components/layout/SettingsContainer.vue';
import { useTrack } from '~/compositions/useTrack';
import { localStoragePrefix } from '~/compositions/useUserSettings';
import { buildDate, feedbackMail } from '~/config';

const { t } = useI18n();
const { track } = useTrack();

const message = useStorage(`${localStoragePrefix}.contact_message`, t('contact_email_body'));

async function sendEmail() {
  const subject = encodeURIComponent(t('feedback_subject'));
  const additionalData = {
    version: buildDate,
  };
  const body = encodeURIComponent(
    `${message.value}\n\n---\n${Object.entries(additionalData)
      .map(([key, value]) => `${key}: ${value}`)
      .join('\n')}\n`,
  );
  track('contact:send-email');
  window.open(`mailto:${feedbackMail}?subject=${subject}&body=${body}`);
  message.value = t('contact_email_body');
}
</script>
