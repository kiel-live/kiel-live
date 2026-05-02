<template>
  <SettingsContainer>
    <h1 class="mb-4 text-xl font-bold">{{ t('kiel_live') }}</h1>

    <img src="../../assets/logo.png" :alt="t('logo_alt')" class="mx-auto mb-4 w-36" />

    <p class="mb-4 text-center">{{ t('where_is_my_bus') }}</p>

    <p class="mb-2 text-center">{{ t('support_us') }}</p>

    <Button
      href="https://www.instagram.com/kiel.live/"
      class="mx-auto mb-auto w-8/10"
      @click="track('click:instagram')"
    >
      <i-ph-instagram-logo class="mr-2" />
      <span>{{ t('follow_instagram', { account: '@kiel.live' }) }}</span>
    </Button>

    <div class="mt-4 flex flex-col">
      <router-link
        :to="{ name: 'settings-settings' }"
        class="flex items-center gap-2 border-t border-b border-gray-200 px-2 py-4 dark:border-neutral-600"
      >
        <i-ph-gear />
        <span>{{ t('settings') }}</span>
      </router-link>
      <router-link
        :to="{ name: 'settings-contact' }"
        class="flex items-center gap-2 border-b border-gray-200 px-2 py-4 dark:border-neutral-600"
      >
        <i-ph-chat-text />
        <span>{{ t('contact_us') }}</span>
      </router-link>
      <a
        :href="analyticsUrl"
        target="_blank"
        rel="noopener noreferrer"
        class="flex items-center gap-2 border-b border-gray-200 px-2 py-4 dark:border-neutral-600"
        @click="track('click:analytics')"
      >
        <i-ph-chart-line />
        <span>{{ t('analytics') }}</span>
      </a>
      <router-link
        :to="{ name: 'settings-changelog' }"
        class="flex items-center gap-2 border-b border-gray-200 px-2 py-4 dark:border-neutral-600"
      >
        <i-ph-clock-clockwise />
        <span>{{ t('changelog') }}</span>
      </router-link>
      <a
        href="https://github.com/kiel-live/kiel-live"
        target="_blank"
        rel="noopener noreferrer"
        class="flex items-center gap-2 px-2 py-4"
        @click="track('click:github')"
      >
        <i-ph-github-logo />
        <span>{{ t('develop_on_github') }}</span>
      </a>
    </div>

    <p v-if="buildDate" class="flex gap-1 text-sm text-gray-500 dark:text-gray-400">
      <span>{{ t('version_from') }}</span>
      <router-link :to="{ name: 'dev' }">{{ buildDate }}</router-link>
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
