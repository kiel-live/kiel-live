import type { Api } from './types';
import { useFeatureFlags } from '~/compositions/useFeatureFlags';
import { DummyApi } from './dummy/dummy';
import { NatsApi } from './nats';
import { RPCApi } from './rpc';

const { checkFeatureFlag } = useFeatureFlags();

export const api: Api =
  import.meta.env.VITE_USE_DUMMY_API === 'true'
    ? new DummyApi()
    : checkFeatureFlag('new_api')
      ? new RPCApi()
      : new NatsApi();
