import type { StorageLike } from '@vueuse/core';

// Injected by the Android app (SharedPreferences backed) so data survives
// a change of the WebView domain.
declare global {
  interface Window {
    NativeStorage?: StorageLike;
  }
}

const nativeStorage = typeof window !== 'undefined' ? window.NativeStorage : undefined;

// Copy existing localStorage data into native storage once, so it is kept
// even if the WebView origin (and with it localStorage) changes later.
if (nativeStorage) {
  for (let i = 0; i < localStorage.length; i++) {
    const key = localStorage.key(i);
    if (key === null || nativeStorage.getItem(key) !== null) {
      continue;
    }
    const value = localStorage.getItem(key);
    if (value !== null) {
      nativeStorage.setItem(key, value);
    }
  }
}

/**
 * Prefers native storage (Android app) and falls back to localStorage for reads.
 * Writes go to both, so localStorage stays a valid fallback.
 */
export const appStorage: StorageLike = {
  getItem(key) {
    return nativeStorage?.getItem(key) ?? localStorage.getItem(key);
  },
  setItem(key, value) {
    nativeStorage?.setItem(key, value);
    localStorage.setItem(key, value);
  },
  removeItem(key) {
    nativeStorage?.removeItem(key);
    localStorage.removeItem(key);
  },
};
