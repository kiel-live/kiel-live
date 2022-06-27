<template>
  <div class="flex flex-col min-h-0 flex-grow overflow-y-auto p-2">
    <h1 class="mb-2 text-xl font-bold">{{ t('kiel_live') }}</h1>

    <img src="../../assets/logo.png" alt="..." class="w-36 mx-auto my-4" />

    <p class="mb-2 text-center">{{ t('where_is_my_bus') }}</p>

    <p class="mb-2 text-center">{{ t('support_us') }}</p>

    <Button href="https://www.instagram.com/kiel.live/" class="mx-auto w-8/10">
      <i-mdi-instagram class="mr-2" />
      <span>{{ t('follow_instagram', { account: '@kiel.live' }) }}</span>
    </Button>

    <div class="flex flex-col mt-4">
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
      <router-link :to="{ name: 'settings-changelog' }" class="flex items-center py-4 px-2 gap-2">
        <i-mdi-github />
        <span>{{ t('develop_on_github') }}</span>
      </router-link>
    </div>

    <p v-if="buildDate" class="flex dark:text-gray-500 text-sm gap-1">
      <span>{{ t('version_from') }}</span>
      <span>{{ buildDate }}</span>
    </p>
  </div>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import Button from '~/components/atomic/Button.vue';
import { analyticsUrl, buildDate, feedbackMail } from '~/config';

const { t } = useI18n();
</script>
