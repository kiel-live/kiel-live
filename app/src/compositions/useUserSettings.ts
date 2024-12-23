import type { LngLatLike } from 'maplibre-gl';
import { useStorage } from '@vueuse/core';

export const localStoragePrefix = 'kiel_live';

const userSettings = {
  liteMode: useStorage(`${localStoragePrefix}.lite`, false),
  lastLocation: useStorage<{ center: LngLatLike; zoom: number; pitch: number; bearing: number }>(
    `${localStoragePrefix}.last_location`,
    {
      center: [10.1283, 54.3166],
      zoom: 14,
      pitch: 0,
      bearing: 0,
    },
  ),
};

export function useUserSettings() {
  return userSettings;
}
