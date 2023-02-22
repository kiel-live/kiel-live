import { JetStreamClient } from 'nats.ws';
import { describe, expect, it, vi } from 'vitest';
import { ref } from 'vue';

import { isConnected, js, subscribe } from '.';

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
  it('should only subscribe once when called multiple times', async () => {
    const state = ref({});
    js.value = { subscribe: vi.fn(() => []) } as unknown as JetStreamClient;
    isConnected.value = true;

    await Promise.all([subscribe('test', state), subscribe('test', state), subscribe('test', state)]);

    // eslint-disable-next-line @typescript-eslint/unbound-method
    expect(js.value.subscribe).toBeCalledTimes(1);
  });
});
