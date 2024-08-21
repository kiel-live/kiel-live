<template>
  <div v-if="actions.length > 0" class="flex mt-2 overflow-x-auto">
    <div class="inline-flex w-min gap-2 pb-4">
      <template v-for="action in actions" :key="action.url">
        <a :href="action.url" target="_blank" rel="noopener noreferrer" class="flex flex-shrink-0">
          <button type="button" class="border border-blue-400 bg-blue-300 rounded-full px-4 py-1 dark:text-white">
            <template v-if="action.type === 'navigate-to'">
              <i-ic-baseline-directions />
              <span>{{ t('navigate_to') }}</span>
            </template>
            <template v-else-if="action.type === 'rent'">
              <i-ic-baseline-play-arrow />
              <span>{{ t('rent_vehicle') }}</span>
            </template>
            <template v-else>
              <span>{{ action.name }}</span>
            </template>
          </button>
        </a>
      </template>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import { Action } from '~/api/types/action';

const { t } = useI18n();

const actions = defineModel<Action[]>('actions', {
  required: true,
});
</script>
