<template>
  <SettingsContainer>
    <h1 class="mb-4 text-xl font-bold">{{ t('kiel_live') }}</h1>

    <img src="../../assets/logo.png" :alt="t('logo_alt')" class="w-36 mx-auto mb-4" />

    <p class="mb-4 text-center">{{ t('contact_title') }}</p>

    <form class="flex flex-col gap-4 items-center mx-auto mb-auto w-8/10" @submit.prevent="sendEmail">
      <textarea
        v-model="message"
        rows="10"
        class="w-full p-2 rounded-md border border-gray-300 bg-white dark:border-gray-400 dark:bg-neutral-800"
        type="text"
      />

      <Button type="submit">
        <i-mdi-email class="mr-2" />
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
