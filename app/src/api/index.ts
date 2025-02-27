import type { Api } from './types';
import { NatsApi } from './nats';

export const api: Api = new NatsApi();
