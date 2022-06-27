import { useStorage } from '@vueuse/core';

const localStoragePrefix = 'kiel_live';

const userSettings = {
  liteMode: useStorage(`${localStoragePrefix}.lite`, false),
};

export function useUserSettings() {
  return userSettings;
}
