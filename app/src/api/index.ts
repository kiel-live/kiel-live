import { useFeatureFlags } from '~/compositions/useFeatureFlags';
import { RPCApi } from './nats';
import type { Api } from './types';

const { checkFeatureFlag } = useFeatureFlags();

export const api: Api = checkFeatureFlag('new_api') ? new RPCApi() : new RPCApi();
