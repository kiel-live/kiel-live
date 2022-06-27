<template>
  <div class="flex flex-col min-h-0 flex-grow">
    <h1 class="flex p-2 mb-2 border-b-1 dark:border-dark-100 space-x-2 items-center">
      <i-ph-info />
      <span class="text-lg">{{ t('kiel_live') }}</span>
    </h1>
    <div class="flex items-start flex-col gap-2 p-2 overflow-y-auto">
      <p class="mb-2">
        {{ t('where_is_my_bus') }}
      </p>
      <p class="mb-2">{{ t('support_us') }}</p>
      <Button href="https://www.instagram.com/kiel.live/">
        <i-mdi-instagram class="mr-2" /><span>{{ t('follow_instagram', { account: '@kiel.live' }) }}</span>
      </Button>
      <Button href="https://github.com/kiel-live/kiel-live">
        <i-mdi-github class="mr-2" /><span>{{ t('develop_on_github') }}</span>
      </Button>
      <Button href="https://boomerang.ju60.de/share/z8KAHmGY/Kiel%20Live%202.0">
        <i-ion-md-analytics class="mr-2" /><span>{{ t('analytics') }}</span>
      </Button>

      <h2 class="border-b-2 w-full dark:border-dark-100 mt-4 font-bold">{{ t('changelog') }}</h2>
      <h3 class="border-b-1 dark:border-gray-600 mt-2">{{ new Date('2022-05-03').toLocaleDateString() }}</h3>
      <span>{{ t('changelog_rewrite_version_2') }}</span>
      <h3 class="border-b-1 dark:border-gray-600 mt-2">{{ new Date('2020-05-13').toLocaleDateString() }}</h3>
      <span>{{ t('changelog_live_locations') }}</span>
      <h3 class="border-b-1 dark:border-gray-600 mt-2">{{ new Date('2019-06-22').toLocaleDateString() }}</h3>
      <span>{{ t('changelog_alerts') }}</span>
      <h3 class="border-b-1 dark:border-gray-600 mt-2">{{ new Date('2019-05-20').toLocaleDateString() }}</h3>
      <span><i-fa-solid-birthday-cake class="align-text-bottom" /> {{ t('changelog_project_started') }}</span>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import Button from '~/components/atomic/Button.vue';

const { t } = useI18n();
</script>
