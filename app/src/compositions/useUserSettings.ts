import { useStorage } from '@vueuse/core';

export const localStoragePrefix = 'kiel_live';

const userSettings = {
  liteMode: useStorage(`${localStoragePrefix}.lite`, false),
  geolocationAllowed: useStorage(`${localStoragePrefix}.geolocationAllowed`, false),
};

export function useUserSettings() {
  return userSettings;
}
