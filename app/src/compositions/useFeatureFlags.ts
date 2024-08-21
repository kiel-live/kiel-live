import { useStorage } from '@vueuse/core';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { localStoragePrefix } from './useUserSettings';

const enabledFeatureFlags = useStorage<string[]>(`${localStoragePrefix}.feature_flags`, []);

const { t } = useI18n();

type FeatureFlag = {
  id: string;
  name: string;
  description?: string;
};

const featureFlags = [
  {
    id: 'new_api',
    name: t('feature_flag_new_api'),
  },
  {
    id: 'vehicle_stop_actions',
    name: t('feature_flag_actions'),
    description: t('feature_flag_actions_description'),
  },
] satisfies FeatureFlag[];

export function useFeatureFlags() {
  return featureFlags.map((flag) => ({
    ...flag,
    enabled: computed({
      get: () => enabledFeatureFlags.value.includes(flag.name),
      set: (value) => {
        if (value) {
          enabledFeatureFlags.value = [...enabledFeatureFlags.value, flag.name];
        } else {
          enabledFeatureFlags.value = enabledFeatureFlags.value.filter((f) => f !== flag.name);
        }
      },
    }),
  }));
}
