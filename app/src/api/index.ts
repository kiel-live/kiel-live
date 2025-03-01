import type { Api } from './types';
import { DummyApi } from './dummy';
import { NatsApi } from './nats';

export const api: Api = import.meta.env.VITE_USE_DUMMY_API === 'true' ? new DummyApi() : new NatsApi();
