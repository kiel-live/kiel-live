import type { JetStreamManager } from '@nats-io/jetstream';
import type { NatsConnection } from '@nats-io/nats-core';
import { jetstreamManager } from '@nats-io/jetstream';
import { wsconnect } from '@nats-io/nats-core';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ref } from 'vue';
import { NatsApi } from './nats';

vi.mock('@nats-io/nats-core', () => ({
  wsconnect: vi.fn(),
  createInbox: vi.fn(() => '_INBOX.test'),
}));

vi.mock('@nats-io/jetstream', () => ({
  jetstreamManager: vi.fn(),
  DeliverPolicy: { All: 'all' },
  AckPolicy: { None: 'none' },
  ReplayPolicy: { Instant: 'instant' },
}));

vi.mock('~/config', () => ({
  natsServerUrl: 'ws://test',
}));

async function createApi() {
  const unsubscribeFn = vi.fn();
  const addFn = vi.fn(async () => ({}));
  vi.mocked(wsconnect).mockResolvedValue({
    subscribe: vi.fn(() => ({
      unsubscribe: unsubscribeFn,
      [Symbol.asyncIterator]() {
        return { next: async () => ({ done: true as const, value: undefined }) };
      },
    })),
    async *status() {},
  } as unknown as NatsConnection);

  vi.mocked(jetstreamManager).mockResolvedValue({
    streams: { find: vi.fn(async () => 'test-stream') },
    consumers: { add: addFn },
  } as unknown as JetStreamManager);

  const api = new NatsApi(false);
  await api.load();

  return { api, addFn, unsubscribeFn };
}

describe('api', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should only subscribe once when called multiple times', async () => {
    const { api, addFn } = await createApi();
    const state = ref({});
    await Promise.all([api.subscribe('test', state), api.subscribe('test', state), api.subscribe('test', state)]);
    expect(addFn).toHaveBeenCalledOnce();
  });

  it('should unsubscribe immediately after subscribing', async () => {
    const { api, unsubscribeFn } = await createApi();
    const state = ref({});
    await Promise.all([api.subscribe('test', state), api.unsubscribe('test')]);
    expect(unsubscribeFn).toHaveBeenCalledOnce();
  });

  it('should subscribe after unsubscribing', async () => {
    const { api, unsubscribeFn, addFn } = await createApi();
    const state = ref({});
    await api.subscribe('test', state);
    await api.unsubscribe('test');
    await api.subscribe('test', state);
    await api.unsubscribe('test');
    expect(addFn).toHaveBeenCalledTimes(2);
    expect(unsubscribeFn).toHaveBeenCalledTimes(2);
  });
});
