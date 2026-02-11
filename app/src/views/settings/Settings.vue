<template>
  <SettingsContainer>
    <h1 class="mb-6 text-2xl font-bold">{{ t('settings') }}</h1>

    <div class="flex flex-col gap-3">
      <div
        class="flex items-center justify-between gap-4 rounded-xl border border-gray-100 bg-white p-4 transition-colors hover:bg-gray-50 dark:border-neutral-950 dark:bg-neutral-800 dark:hover:bg-neutral-900"
      >
        <label class="flex flex-col" for="lite-mode">
          <span class="font-medium">{{ t('lite_mode') }}</span>
          <span class="text-sm text-gray-500 dark:text-gray-400">{{ t('lite_mode_description') }}</span>
        </label>
        <Checkbox
          id="lite-mode"
          v-model="liteMode"
          @update:model-value="track('setting:lite-mode', { enabled: $event })"
        />
      </div>

      <div
        class="flex items-center justify-between gap-4 rounded-xl border border-gray-100 bg-white p-4 transition-colors hover:bg-gray-50 dark:border-neutral-950 dark:bg-neutral-800 dark:hover:bg-neutral-900"
      >
        <label class="flex flex-col" for="theme">
          <span class="font-medium">{{ t('theme') }}</span>
          <span class="text-sm text-gray-500 dark:text-gray-400">{{ t('theme_description') }}</span>
        </label>
        <Select id="theme" v-model="theme" :options="options" />
      </div>
    </div>
  </SettingsContainer>
</template>

<script lang="ts" setup>
import type { Theme } from '~/compositions/useColorMode';
import { computed } from 'vue';

import { useI18n } from 'vue-i18n';
import Checkbox from '~/components/atomic/Checkbox.vue';
import Select from '~/components/atomic/Select.vue';
import SettingsContainer from '~/components/layout/SettingsContainer.vue';
import { useColorMode } from '~/compositions/useColorMode';
import { useTrack } from '~/compositions/useTrack';
import { useUserSettings } from '~/compositions/useUserSettings';

const { liteMode } = useUserSettings();
const { t } = useI18n();
const { track } = useTrack();

const theme = useColorMode({ emitAuto: true });

const options = computed<{ value: Theme; label: string }[]>(() => [
  { value: 'auto', label: t('auto') },
  { value: 'light', label: t('light') },
  { value: 'dark', label: t('dark') },
]);
</script>
