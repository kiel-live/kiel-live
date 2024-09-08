import { NatsApi } from './nats';
import type { Api } from './types';

export const api: Api = new NatsApi();
