import { ref } from 'vue';

const localStoragePrefix = 'kiel_live.';

export function useUserSettings() {
  return {
    useLiteMode: ref(true),
  };
}
