import type { LngLatLike } from 'maplibre-gl';
import { useStorage } from '@vueuse/core';

import { appStorage } from '~/storage';

export const localStoragePrefix = 'kiel_live';

const userSettings = {
  liteMode: useStorage(`${localStoragePrefix}.lite`, false, appStorage),
  lastLocation: useStorage<{ center: LngLatLike; zoom: number; pitch: number; bearing: number }>(
    `${localStoragePrefix}.last_location`,
    {
      center: [10.1283, 54.3166],
      zoom: 14,
      pitch: 0,
      bearing: 0,
    },
    appStorage,
  ),
};

export function useUserSettings() {
  return userSettings;
}
