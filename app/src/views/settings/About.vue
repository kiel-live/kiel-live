<template>
  <SettingsContainer>
    <h1 class="mb-4 text-xl font-bold">{{ t('kiel_live') }}</h1>

    <img src="../../assets/logo.png" :alt="t('logo_alt')" class="w-36 mx-auto mb-4" />

    <p class="mb-4 text-center">{{ t('where_is_my_bus') }}</p>

    <p class="mb-2 text-center">{{ t('support_us') }}</p>

    <Button
      href="https://www.instagram.com/kiel.live/"
      class="mx-auto mb-auto w-8/10"
      @click="track('click:instagram')"
    >
      <i-mdi-instagram class="mr-2" />
      <span>{{ t('follow_instagram', { account: '@kiel.live' }) }}</span>
    </Button>

    <div class="flex flex-col mt-4">
      <router-link
        :to="{ name: 'settings-settings' }"
        class="flex items-center border-b border-t border-gray-200 dark:border-neutral-600 py-4 px-2 gap-2 rounded-t-md"
      >
        <i-ph-gear-fill />
        <span>{{ t('settings') }}</span>
      </router-link>
      <router-link
        :to="{ name: 'settings-contact' }"
        class="flex items-center border-b border-gray-200 dark:border-neutral-600 py-4 px-2 gap-2"
      >
        <i-ic-round-message />
        <span>{{ t('contact_us') }}</span>
      </router-link>
      <a
        :href="analyticsUrl"
        target="_blank"
        rel="noopener noreferrer"
        class="flex items-center border-b border-gray-200 dark:border-neutral-600 py-4 px-2 gap-2"
        @click="track('click:analytics')"
      >
        <i-ion-md-analytics />
        <span>{{ t('analytics') }}</span>
      </a>
      <router-link
        :to="{ name: 'settings-changelog' }"
        class="flex items-center py-4 px-2 gap-2 border-b border-gray-200 dark:border-neutral-600"
      >
        <i-ic-outline-update />
        <span>{{ t('changelog') }}</span>
      </router-link>
      <a
        href="https://github.com/kiel-live/kiel-live"
        target="_blank"
        rel="noopener noreferrer"
        class="flex items-center py-4 px-2 gap-2 rounded-b-md"
        @click="track('click:github')"
      >
        <i-mdi-github />
        <span>{{ t('develop_on_github') }}</span>
      </a>
    </div>

    <p v-if="buildDate" class="flex text-gray-500 dark:text-gray-400 text-sm gap-1">
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
