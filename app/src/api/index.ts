import type { Api } from './types';
import { useFeatureFlags } from '~/compositions/useFeatureFlags';
import { DummyApi } from './dummy/dummy';
import { HttpApi } from './http';
import { NatsApi } from './nats';

const { checkFeatureFlag } = useFeatureFlags();

function getApi() {
  if (import.meta.env.VITE_USE_DUMMY_API === 'true') {
    return new DummyApi();
  }

  if (checkFeatureFlag('new_api').value) {
    return new HttpApi();
  }

  return new NatsApi();
}

export const api: Api = getApi();
