import { BasicColorSchema, useColorMode as useColorModeVueUse, UseColorModeOptions } from '@vueuse/core';

import { localStoragePrefix } from '~/compositions/useUserSettings';

export const useColorMode = (options?: UseColorModeOptions) =>
  useColorModeVueUse({ storageKey: `${localStoragePrefix}.theme`, ...options });

export type Theme = BasicColorSchema;
