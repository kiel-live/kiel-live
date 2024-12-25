import { useFeatureFlags } from '~/compositions/useFeatureFlags';
import { NatsApi } from './nats';
import { RPCApi } from './rpc';
import type { Api } from './types';

const { checkFeatureFlag } = useFeatureFlags();

export const api: Api = checkFeatureFlag('new_api') ? new RPCApi() : new NatsApi();
