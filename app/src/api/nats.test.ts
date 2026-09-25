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

function createStatusEmitter() {
  const queue: { type: string }[] = [];
  let notify: (() => void) | undefined;

  return {
    emit(type: string) {
      queue.push({ type });
      notify?.();
    },
    async *status() {
      while (true) {
        const next = queue.shift();
        if (next) {
          yield next;
        } else {
          await new Promise<void>((resolve) => {
            notify = resolve;
          });
        }
      }
    },
  };
}

async function createApi() {
  const unsubscribeFn = vi.fn();
  const addFn = vi.fn(async () => ({}));
  const statusEmitter = createStatusEmitter();
  vi.mocked(wsconnect).mockResolvedValue({
    subscribe: vi.fn(() => {
      // every inbox subscription delivers a single message and then ends
      let delivered = false;
      return {
        unsubscribe: unsubscribeFn,
        [Symbol.asyncIterator]() {
          return {
            next: async () => {
              if (delivered) {
                return { done: true as const, value: undefined };
              }
              delivered = true;
              return { done: false as const, value: { string: () => JSON.stringify({ id: 'model-1' }) } };
            },
          };
        },
      };
    }),
    status: () => statusEmitter.status(),
  } as unknown as NatsConnection);

  vi.mocked(jetstreamManager).mockResolvedValue({
    streams: { find: vi.fn(async () => 'test-stream') },
    consumers: { add: addFn },
  } as unknown as JetStreamManager);

  const api = new NatsApi(false);
  await api.load();

  return { api, addFn, unsubscribeFn, statusEmitter };
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

  it('should re-create consumers of active subscriptions after a reconnect', async () => {
    const { api, addFn, unsubscribeFn, statusEmitter } = await createApi();
    const state = ref({});
    await api.subscribe('test', state);
    expect(addFn).toHaveBeenCalledOnce();

    statusEmitter.emit('disconnect');
    await vi.waitFor(() => expect(api.isConnected.value).toBe(false));
    statusEmitter.emit('reconnect');

    await vi.waitFor(() => expect(addFn).toHaveBeenCalledTimes(2));
    expect(api.isConnected.value).toBe(true);
    // the stale inbox subscription of the old consumer is dropped
    expect(unsubscribeFn).toHaveBeenCalledOnce();
    // messages of the new consumer still end up in the state
    await vi.waitFor(() => expect(state.value).toHaveProperty('model-1'));

    // still a single active subscription that can be removed
    await api.unsubscribe('test');
    expect(unsubscribeFn).toHaveBeenCalledTimes(2);
    await api.subscribe('test', state);
    expect(addFn).toHaveBeenCalledTimes(3);
  });

  it('should not re-create subscriptions removed while disconnected', async () => {
    const { api, addFn, statusEmitter } = await createApi();
    const state = ref({});
    await api.subscribe('test', state);

    statusEmitter.emit('disconnect');
    await vi.waitFor(() => expect(api.isConnected.value).toBe(false));
    await api.unsubscribe('test');
    statusEmitter.emit('reconnect');

    await vi.waitFor(() => expect(api.isConnected.value).toBe(true));
    await new Promise((resolve) => setTimeout(resolve, 10));
    expect(addFn).toHaveBeenCalledOnce();
  });

  it('should subscribe to topics requested while disconnected after a reconnect', async () => {
    const { api, addFn, statusEmitter } = await createApi();
    const state = ref({});

    statusEmitter.emit('disconnect');
    await vi.waitFor(() => expect(api.isConnected.value).toBe(false));
    await api.subscribe('test', state);
    expect(addFn).not.toHaveBeenCalled();

    statusEmitter.emit('reconnect');
    await vi.waitFor(() => expect(addFn).toHaveBeenCalledOnce());
  });

  it('should retry a subscription on reconnect when creating the consumer failed', async () => {
    const { api, addFn, unsubscribeFn, statusEmitter } = await createApi();
    vi.spyOn(console, 'error').mockImplementation(() => {});
    addFn.mockRejectedValueOnce(new Error('connection lost'));
    const state = ref({});

    await api.subscribe('test', state);
    expect(addFn).toHaveBeenCalledOnce();
    expect(unsubscribeFn).toHaveBeenCalledOnce();

    statusEmitter.emit('disconnect');
    await vi.waitFor(() => expect(api.isConnected.value).toBe(false));
    statusEmitter.emit('reconnect');
    await vi.waitFor(() => expect(addFn).toHaveBeenCalledTimes(2));
  });
});
