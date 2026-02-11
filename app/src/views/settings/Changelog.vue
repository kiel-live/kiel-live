<template>
  <SettingsContainer>
    <h1 class="mb-6 text-2xl font-bold">{{ t('changelog') }}</h1>

    <div class="flex flex-col gap-4">
      <div
        v-for="(entry, id) in changelog"
        :key="id"
        class="rounded-xl border border-gray-100 bg-white p-4 dark:border-neutral-800 dark:bg-neutral-950"
      >
        <h2 class="mb-2 font-semibold text-gray-900 dark:text-gray-100">
          {{ entry.date.toLocaleDateString() }}
        </h2>
        <span class="text-gray-600 dark:text-gray-400">{{ entry.text }}</span>
      </div>
      <div class="rounded-xl border border-gray-100 bg-white p-4 dark:border-neutral-800 dark:bg-neutral-950">
        <h2 class="mb-2 font-semibold text-gray-900 dark:text-gray-100">
          {{ new Date('2019-05-20').toLocaleDateString() }}
        </h2>
        <span class="text-gray-600 dark:text-gray-400"
          ><i-ph-cake-fill class="inline align-text-bottom" /> {{ t('changelog_project_started') }}</span
        >
      </div>
    </div>
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
