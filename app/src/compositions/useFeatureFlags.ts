import type { Ref } from 'vue';
import { useStorage } from '@vueuse/core';
import { computed, getCurrentInstance } from 'vue';
import { useI18n } from 'vue-i18n';

import { localStoragePrefix } from './useUserSettings';

const enabledFeatureFlags = useStorage<string[]>(`${localStoragePrefix}.feature_flags`, []);

interface FeatureFlag {
  id: string;
  name: string;
  description?: string;
}

export function useFeatureFlags() {
  // we can only use i18n inside components, as it is only needed for the settings mock i18n otherwise
  let t: ReturnType<typeof useI18n>['t'] = (key) => key;
  if (getCurrentInstance()) {
    t = useI18n().t;
  }

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

  const featureFlagWithEnabled = featureFlags.map((flag) => ({
    ...flag,
    enabled: computed({
      get: () => enabledFeatureFlags.value.includes(flag.id),
      set: (value) => {
        if (value) {
          enabledFeatureFlags.value = [...enabledFeatureFlags.value, flag.id];
        } else {
          enabledFeatureFlags.value = enabledFeatureFlags.value.filter((f) => f !== flag.id);
        }
      },
    }),
  }));

  function checkFeatureFlag(flagId: string): Ref<boolean> {
    const featureFlag = featureFlagWithEnabled.find((flag) => flag.id === flagId);

    if (!featureFlag) {
      throw new Error(`Unknown feature flag: ${flagId}`);
    }

    return featureFlag.enabled;
  }

  return {
    featureFlags: featureFlagWithEnabled,
    checkFeatureFlag,
  };
}
