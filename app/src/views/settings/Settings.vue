<template>
  <div class="flex flex-col min-h-0 flex-grow overflow-y-auto p-2">
    <h1 class="mb-2 text-xl font-bold">{{ t('settings') }}</h1>

    <div class="flex gap-4 items-center">
      <div class="flex flex-col">
        <span>{{ t('lite_mode') }}</span>
        <span class="text-sm dark:text-gray-400">{{ t('lite_mode_description') }}</span>
      </div>
      <Checkbox v-model="liteMode" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import Checkbox from '~/components/atomic/Checkbox.vue';
import { useUserSettings } from '~/compositions/useUserSettings';

const { liteMode } = useUserSettings();
const { t } = useI18n();
</script>
