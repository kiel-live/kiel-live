import type { JetStreamManager } from '@nats-io/jetstream';
import type { NatsConnection, Subscription } from '@nats-io/nats-core';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ref } from 'vue';

function makeSubscription(unsubscribeFn = vi.fn()): Subscription {
  return {
    unsubscribe: unsubscribeFn,
    [Symbol.asyncIterator]() {
      return {
        next: async () =>
          new Promise((resolve) => {
            resolve({ done: true, value: undefined });
          }),
      };
    },
  } as unknown as Subscription;
}

function makeJsm(addFn = vi.fn(async () => ({}))): JetStreamManager {
  return {
    streams: { find: vi.fn(async () => 'test-stream') },
    consumers: { add: addFn },
  } as unknown as JetStreamManager;
}

function makeNc(sub: Subscription): NatsConnection {
  return { subscribe: vi.fn(() => sub) } as unknown as NatsConnection;
}

describe('api', () => {
  beforeEach(() => {
    vi.resetModules();
  });

  it('should only subscribe once when called multiple times', async () => {
    const { NatsApi } = await import('./nats');
    const api = new NatsApi(false);
    const state = ref({});
    const addMock = vi.fn(async () => ({}));
    api.jsm = makeJsm(addMock);
    api.nc = makeNc(makeSubscription());
    api.isConnected.value = true;

    await Promise.all([api.subscribe('test', state), api.subscribe('test', state), api.subscribe('test', state)]);

    expect(addMock).toBeCalledTimes(1);
  });

  it('should unsubscribe immediately after subscribing', async () => {
    const { NatsApi } = await import('./nats');
    const api = new NatsApi(false);
    const unsubscribeMock = vi.fn();
    api.jsm = makeJsm();
    api.nc = makeNc(makeSubscription(unsubscribeMock));
    api.isConnected.value = true;

    const state = ref({});
    await Promise.all([api.subscribe('test', state), api.unsubscribe('test')]);

    expect(unsubscribeMock).toBeCalledTimes(1);
  });

  it('should subscribe after unsubscribing', async () => {
    const { NatsApi } = await import('./nats');
    const api = new NatsApi(false);
    const unsubscribeMock = vi.fn();
    const addMock = vi.fn(async () => ({}));
    api.jsm = makeJsm(addMock);
    api.nc = makeNc(makeSubscription(unsubscribeMock));
    api.isConnected.value = true;

    const state = ref({});
    await api.subscribe('test', state);
    await api.unsubscribe('test');
    await api.subscribe('test', state);
    await api.unsubscribe('test');

    expect(addMock).toBeCalledTimes(2);
    expect(unsubscribeMock).toBeCalledTimes(2);
  });
});
