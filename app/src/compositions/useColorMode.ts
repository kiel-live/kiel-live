import type { BasicColorSchema, UseColorModeOptions } from '@vueuse/core';
import { useColorMode as useColorModeVueUse } from '@vueuse/core';

import { localStoragePrefix } from '~/compositions/useUserSettings';
import { appStorage } from '~/storage';

export const useColorMode = (options?: UseColorModeOptions) =>
  useColorModeVueUse({ storageKey: `${localStoragePrefix}.theme`, storage: appStorage, ...options });

export type Theme = BasicColorSchema;
