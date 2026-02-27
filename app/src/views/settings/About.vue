<template>
  <SettingsContainer>
    <h1 class="mb-6 text-2xl font-bold">{{ t('kiel_live') }}</h1>

    <img src="../../assets/logo.png" :alt="t('logo_alt')" class="mx-auto mb-6 w-32" />

    <p class="mb-4 text-center text-gray-600 dark:text-gray-400">{{ t('where_is_my_bus') }}</p>

    <p class="mb-3 text-center text-sm text-gray-600 dark:text-gray-400">{{ t('support_us') }}</p>

    <Button
      href="https://www.instagram.com/kiel.live/"
      class="mx-auto mb-8 w-full max-w-sm"
      @click="track('click:instagram')"
    >
      <i-ph-instagram-logo class="mr-2" />
      <span>{{ t('follow_instagram', { account: '@kiel.live' }) }}</span>
    </Button>

    <div
      class="mb-6 flex flex-col overflow-hidden rounded-xl border border-gray-100 bg-white dark:border-neutral-950 dark:bg-neutral-800"
    >
      <router-link
        :to="{ name: 'settings-settings' }"
        class="flex items-center gap-3 border-b border-gray-100 px-4 py-4 transition-colors hover:bg-gray-50 dark:border-neutral-950 dark:hover:bg-neutral-900"
      >
        <i-ph-gear class="text-xl" />
        <span class="font-medium">{{ t('settings') }}</span>
      </router-link>
      <router-link
        :to="{ name: 'settings-contact' }"
        class="flex items-center gap-3 border-b border-gray-100 px-4 py-4 transition-colors hover:bg-gray-50 dark:border-neutral-950 dark:hover:bg-neutral-900"
      >
        <i-ph-chat-text class="text-xl" />
        <span class="font-medium">{{ t('contact_us') }}</span>
      </router-link>
      <a
        :href="analyticsUrl"
        target="_blank"
        rel="noopener noreferrer"
        class="flex items-center gap-3 border-b border-gray-100 px-4 py-4 transition-colors hover:bg-gray-50 dark:border-neutral-950 dark:hover:bg-neutral-900"
        @click="track('click:analytics')"
      >
        <i-ph-chart-line class="text-xl" />
        <span class="font-medium">{{ t('analytics') }}</span>
      </a>
      <router-link
        :to="{ name: 'settings-changelog' }"
        class="flex items-center gap-3 border-b border-gray-100 px-4 py-4 transition-colors hover:bg-gray-50 dark:border-neutral-950 dark:hover:bg-neutral-900"
      >
        <i-ph-clock-clockwise class="text-xl" />
        <span class="font-medium">{{ t('changelog') }}</span>
      </router-link>
      <a
        href="https://github.com/kiel-live/kiel-live"
        target="_blank"
        rel="noopener noreferrer"
        class="flex items-center gap-3 px-4 py-4 transition-colors hover:bg-gray-50 dark:hover:bg-neutral-900"
        @click="track('click:github')"
      >
        <i-ph-github-logo class="text-xl" />
        <span class="font-medium">{{ t('develop_on_github') }}</span>
      </a>
    </div>

    <p v-if="buildDate" class="flex justify-center gap-1 text-sm text-gray-500 dark:text-gray-400">
      <span>{{ t('version_from') }}</span>
      <router-link :to="{ name: 'dev' }" class="underline">{{ buildDate }}</router-link>
    </p>
  </SettingsContainer>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import Button from '~/components/atomic/Button.vue';
import SettingsContainer from '~/components/layout/SettingsContainer.vue';
import { useTrack } from '~/compositions/useTrack';
import { analyticsUrl, buildDate } from '~/config';

const { t } = useI18n();
const { track } = useTrack();
</script>
