import type { Api } from './types';
import { DummyApi } from './dummy/dummy';
import { NatsApi } from './nats';

function getApi() {
  if (import.meta.env.VITE_USE_DUMMY_API === 'true') {
    return new DummyApi();
  }

  return new NatsApi();
}

export const api: Api = getApi();
