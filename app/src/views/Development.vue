<template>
  <SettingsContainer>
    <h1 class="mb-6 text-2xl font-bold">{{ t('kiel_live') }}</h1>

    <p v-if="buildDate" class="mb-6 flex justify-center gap-1 text-sm text-gray-500 dark:text-gray-400">
      <span>{{ t('version_from') }}</span>
      <router-link :to="{ name: 'dev' }" class="underline">{{ buildDate }}</router-link>
    </p>

    <div>
      <h2 class="mb-2 text-lg font-semibold">{{ t('feature_flags') }}</h2>

      <p class="mb-4 text-sm text-gray-500 dark:text-gray-400">{{ t('feature_flags_description') }}</p>

      <div class="flex flex-col gap-3">
        <div
          v-for="featureFlag in featureFlags"
          :key="featureFlag.id"
          class="flex items-center justify-between gap-4 rounded-xl border border-gray-100 bg-white p-4 transition-colors hover:bg-gray-50 dark:border-neutral-950 dark:bg-neutral-800 dark:hover:bg-neutral-900"
        >
          <label class="flex flex-col" for="lite-mode">
            <span class="font-medium">{{ featureFlag.name }}</span>
            <span v-if="featureFlag.description" class="text-sm text-gray-500 dark:text-gray-400">{{
              featureFlag.description
            }}</span>
          </label>
          <Checkbox
            id="lite-mode"
            :model-value="featureFlag.enabled.value"
            @update:model-value="
              featureFlag.enabled.value = $event;
              track('feature-flag', { enabled: $event, id: featureFlag.id });
            "
          />
        </div>
      </div>
    </div>
  </SettingsContainer>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import Checkbox from '~/components/atomic/Checkbox.vue';
import SettingsContainer from '~/components/layout/SettingsContainer.vue';
import { useFeatureFlags } from '~/compositions/useFeatureFlags';
import { useTrack } from '~/compositions/useTrack';
import { buildDate } from '~/config';

const { t } = useI18n();
const { track } = useTrack();

const { featureFlags } = useFeatureFlags();
</script>
