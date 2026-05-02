<template>
  <SettingsContainer>
    <h1 class="mb-2 text-xl font-bold">{{ t('changelog') }}</h1>

    <div v-for="(entry, id) in changelog" :key="id">
      <h2 class="mt-2 mb-1 border-b border-gray-200 font-bold dark:border-neutral-600">
        {{ entry.date.toLocaleDateString() }}
      </h2>
      <span>{{ entry.text }}</span>
    </div>
    <h2 class="mt-2 mb-1 border-b border-gray-200 font-bold dark:border-neutral-600">
      {{ new Date('2019-05-20').toLocaleDateString() }}
    </h2>
    <span><i-ph-cake-fill class="inline align-text-bottom" /> {{ t('changelog_project_started') }}</span>
  </SettingsContainer>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import SettingsContainer from '~/components/layout/SettingsContainer.vue';

const { t } = useI18n();

interface ChangelogEntry {
  date: Date;
  text: string;
}

const changelog: ChangelogEntry[] = [
  {
    date: new Date('2024-07-05'),
    text: t('changelog_nok_ferry'),
  },
  {
    date: new Date('2022-06-27'),
    text: t('changelog_trips_path_lite_mode'),
  },
  {
    date: new Date('2022-05-03'),
    text: t('changelog_rewrite_version_2'),
  },
  {
    date: new Date('2020-05-13'),
    text: t('changelog_live_locations'),
  },
  {
    date: new Date('2019-06-22'),
    text: t('changelog_alerts'),
  },
];
</script>
