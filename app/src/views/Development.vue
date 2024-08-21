<template>
  <SettingsContainer>
    <h1 class="mb-2 text-xl font-bold">{{ t('kiel_live') }}</h1>

    <p v-if="buildDate" class="flex mb-4 text-gray-500 dark:text-gray-400 text-sm gap-1">
      <span>{{ t('version_from') }}</span>
      <router-link :to="{ name: 'dev' }">{{ buildDate }}</router-link>
    </p>

    <div>
      <h2 class="text-lg font-bold">{{ t('feature_flags') }}</h2>

      <p class="mb-2 text-gray-500 dark:text-gray-400 text-sm">{{ t('feature_flags_description') }}</p>

      <div class="flex flex-col gap-2">
        <div v-for="featureFlag in featureFlags" :key="featureFlag.id" class="flex gap-4 items-center justify-between">
          <label class="flex flex-col" for="lite-mode">
            <span>{{ featureFlag.name }}</span>
            <span v-if="featureFlag.description" class="text-sm text-gray-500 dark:text-gray-400">{{
              featureFlag.description
            }}</span>
          </label>
          <Checkbox
            id="lite-mode"
            :model-value="featureFlag.enabled.value"
            @update:model-value="featureFlag.enabled.value = $event"
          />
        </div>
      </div>
    </div>
  </SettingsContainer>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import SettingsContainer from '~/components/layout/SettingsContainer.vue';
import { useFeatureFlags } from '~/compositions/useFeatureFlags';
import { buildDate } from '~/config';

const { t } = useI18n();

const featureFlags = useFeatureFlags();
</script>
