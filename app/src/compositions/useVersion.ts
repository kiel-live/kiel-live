import { useRegisterSW } from 'virtual:pwa-register/vue';
import { ref } from 'vue';

import { useTrack } from './useTrack';

export function useVersion() {
  const { needRefresh, updateServiceWorker } = useRegisterSW();
  const version = ref(import.meta.env.VITE_BUILD_DATE);

  async function updateApp() {
    useTrack().track('app-update', {
      from: version.value,
    });
    await updateServiceWorker(true);
  }

  return {
    needRefresh,
    version,
    updateApp,
  };
}
