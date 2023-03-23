import { useStorage } from '@vueuse/core';
import { LngLatLike } from 'maplibre-gl';

export const localStoragePrefix = 'kiel_live';

const userSettings = {
  liteMode: useStorage(`${localStoragePrefix}.lite`, false),
  lastLocation: useStorage<LngLatLike>(`${localStoragePrefix}.last_location`, [10.1283, 54.3166]),
};

export function useUserSettings() {
  return userSettings;
}
