import type { JetStreamClient } from 'nats.ws';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ref } from 'vue';

vi.mock('nats.ws', async () => {
  const original = await vi.importActual<typeof import('nats.ws')>('nats.ws');
  return {
    ...original,
    consumerOpts: vi.fn(() => ({
      deliverTo: vi.fn(),
      deliverAll: vi.fn(),
      ackNone: vi.fn(),
      replayInstantly: vi.fn(),
    })),
  };
});

describe('api', () => {
  beforeEach(() => {
    vi.resetModules();
  });

  it('should only subscribe once when called multiple times', async () => {
    const { NatsApi } = await import('./nats');
    const api = new NatsApi(false);
    const state = ref({});
    const subscribeMock = vi.fn(() => []);
    api.js.value = { subscribe: subscribeMock } as unknown as JetStreamClient;
    api.isConnected.value = true;

    await Promise.all([api.subscribe('test', state), api.subscribe('test', state), api.subscribe('test', state)]);

    expect(subscribeMock).toBeCalledTimes(1);
  });

  it('should unsubscribe immediately after subscribing', async () => {
    const { NatsApi } = await import('./nats');
    const api = new NatsApi(false);
    const unsubscribeMock = vi.fn();
    const subscribeMock = vi.fn(() => ({
      unsubscribe: unsubscribeMock,
      [Symbol.asyncIterator]() {
        return {
          next: async () =>
            new Promise((resolve) => {
              resolve({ done: true });
            }),
        };
      },
    }));
    const state = ref({});
    api.js.value = {
      subscribe: subscribeMock,
    } as unknown as JetStreamClient;
    api.isConnected.value = true;

    await Promise.all([api.subscribe('test', state), api.unsubscribe('test')]);

    expect(subscribeMock).toBeCalledTimes(1);
    expect(unsubscribeMock).toBeCalledTimes(1);
  });

  it('should subscribe after unsubscribing', async () => {
    const { NatsApi } = await import('./nats');
    const api = new NatsApi(false);
    const unsubscribeMock = vi.fn();
    const subscribeMock = vi.fn(() => ({
      unsubscribe: unsubscribeMock,
      [Symbol.asyncIterator]() {
        return {
          next: async () =>
            new Promise((resolve) => {
              resolve({ done: true });
            }),
        };
      },
    }));
    const state = ref({});
    api.js.value = {
      subscribe: subscribeMock,
    } as unknown as JetStreamClient;
    api.isConnected.value = true;

    await api.subscribe('test', state);
    await api.unsubscribe('test');
    await api.subscribe('test', state);
    await api.unsubscribe('test');

    expect(subscribeMock).toBeCalledTimes(2);
    expect(unsubscribeMock).toBeCalledTimes(2);
  });
});
