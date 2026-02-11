<template>
  <div v-if="actions.length > 0 && checkFeatureFlag('vehicle_stop_actions').value" class="mt-2 flex overflow-x-auto">
    <div class="inline-flex w-min gap-2 pb-4">
      <template v-for="action in actions" :key="action.url">
        <Button :href="action.url" rounded class="shrink-0 gap-2 px-4 py-2">
          <template v-if="action.type === 'navigate-to'">
            <i-mdi-directions />
            <span>{{ t('navigate_to') }}</span>
          </template>
          <template v-else-if="action.type === 'rent'">
            <i-ph-play-fill />
            <span>{{ t('rent_vehicle') }}</span>
          </template>
          <template v-else>
            <span>{{ action.name }}</span>
          </template>
        </Button>
      </template>
    </div>
  </div>
</template>

<script lang="ts" setup>
import type { Action } from '~/api/types/action';

import { useI18n } from 'vue-i18n';
import Button from '~/components/atomic/Button.vue';
import { useFeatureFlags } from '~/compositions/useFeatureFlags';

const { t } = useI18n();

const actions = defineModel<Action[]>('actions', {
  required: true,
});

const { checkFeatureFlag } = useFeatureFlags();
</script>
