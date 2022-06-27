<template>
  <div>
    <h1 class="mb-2 text-xl font-bold">{{ t('changelog') }}</h1>

    <div v-for="(entry, id) in changelog" :key="id">
      <h3 class="border-b-1 dark:border-gray-600 mt-2">{{ entry.date.toLocaleDateString() }}</h3>
      <span>{{ entry.text }}</span>
    </div>
    <h3 class="border-b-1 dark:border-gray-600 mt-2">{{ new Date('2019-05-20').toLocaleDateString() }}</h3>
    <span><i-fa-solid-birthday-cake class="align-text-bottom" /> {{ t('changelog_project_started') }}</span>
  </div>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

type ChangelogEntry = {
  date: Date;
  text: string;
};

const changelog: ChangelogEntry[] = [
  {
    date: new Date('2022-06-27'),
    text: t('changelog_rewrite_version_2'), // TODO: lite mode, trips path, ui updates (bottom app bar)
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
