<template>
  <SettingsContainer>
    <h1 class="mb-4 text-xl font-bold">{{ t('kiel_live') }}</h1>

    <img src="../../assets/logo.png" alt="..." class="w-36 mx-auto mb-4" />

    <p class="mb-4 text-center">{{ t('where_is_my_bus') }}</p>

    <p class="mb-2 text-center">{{ t('support_us') }}</p>

    <Button href="https://www.instagram.com/kiel.live/" class="mx-auto mb-auto w-8/10">
      <i-mdi-instagram class="mr-2" />
      <span>{{ t('follow_instagram', { account: '@kiel.live' }) }}</span>
    </Button>

    <div class="flex flex-col mt-4">
      <router-link
        :to="{ name: 'settings-settings' }"
        class="flex items-center border-b border-t dark:border-gray-600 py-4 px-2 gap-2"
      >
        <i-ph-gear-fill />
        <span>{{ t('settings') }}</span>
      </router-link>
      <a :href="`mailto:${feedbackMail}`" class="flex items-center border-b dark:border-gray-600 py-4 px-2 gap-2">
        <i-ic-round-message />
        <span>{{ t('give_feedback') }}</span>
      </a>
      <a
        :href="analyticsUrl"
        target="_blank"
        rel="noopener noreferrer"
        class="flex items-center border-b dark:border-gray-600 py-4 px-2 gap-2"
      >
        <i-ion-md-analytics />
        <span>{{ t('analytics') }}</span>
      </a>
      <router-link
        :to="{ name: 'settings-changelog' }"
        class="flex items-center py-4 px-2 gap-2 border-b dark:border-gray-600"
      >
        <i-ic-outline-update />
        <span>{{ t('changelog') }}</span>
      </router-link>
      <a
        href="https://github.com/kiel-live/kiel-live"
        target="_blank"
        rel="noopener noreferrer"
        class="flex items-center py-4 px-2 gap-2"
      >
        <i-mdi-github />
        <span>{{ t('develop_on_github') }}</span>
      </a>
    </div>

    <p v-if="buildDate" class="flex dark:text-gray-500 text-sm gap-1">
      <span>{{ t('version_from') }}</span>
      <span>{{ buildDate }}</span>
    </p>
  </SettingsContainer>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import Button from '~/components/atomic/Button.vue';
import SettingsContainer from '~/components/layout/SettingsContainer.vue';
import { analyticsUrl, buildDate, feedbackMail } from '~/config';

const { t } = useI18n();
</script>
