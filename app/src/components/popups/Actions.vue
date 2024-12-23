<template>
  <div v-if="actions.length > 0 && checkFeatureFlag('vehicle_stop_actions').value" class="flex mt-2 overflow-x-auto">
    <div class="inline-flex w-min gap-2 pb-4">
      <template v-for="action in actions" :key="action.url">
        <Button :href="action.url" rounded class="gap-2 flex-shrink-0 px-4 py-2" @click="doAction(action)">
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
        </Button>
      </template>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import type { Action } from '~/api/types/action';
import Button from '~/components/atomic/Button.vue';
import { useFeatureFlags } from '~/compositions/useFeatureFlags';

const { t } = useI18n();

const actions = defineModel<Action[]>('actions', {
  required: true,
});

const { checkFeatureFlag } = useFeatureFlags();

async function doAction(action: Action) {
  await navigator?.share({
    title: `Kiel Live - ${action.name}`,
    url: window?.location?.href,
  });
}
</script>
