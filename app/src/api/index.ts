import { connect, consumerOpts, createInbox, Events, StringCodec } from 'nats.ws';
import { ref, type Ref } from 'vue';
import type { JetStreamClient, JetStreamSubscription, NatsConnection } from 'nats.ws';

import type { Models, Stop, Trip, Vehicle } from '~/api/types';
import { natsServerUrl } from '~/config';

const sc = StringCodec();

export const DeletePayload = '---';

export const vehicles = ref<Record<string, Vehicle>>({});
export const stops = ref<Record<string, Stop>>({});
export const trips = ref<Record<string, Trip>>({});
export const isConnected = ref(false);

const subscriptions = ref<Record<string, { subscription?: JetStreamSubscription; pending?: Promise<void> }>>({});
const subscriptionsQueue: Record<string, Ref<Record<string, Models>>> = {};

let nc: NatsConnection | undefined;
export const js: Ref<JetStreamClient | undefined> = ref();

export const subscribe = async (subject: string, state: Ref<Record<string, Models>>) => {
  if (subscriptions.value[subject]) {
    return;
  }

  if (!isConnected.value || !js.value) {
    subscriptionsQueue[subject] = state;
    return;
  }

  let resolvePendingSubscription: () => void = () => {};
  subscriptions.value[subject] = {
    pending: new Promise((resolve) => {
      resolvePendingSubscription = resolve;
    }),
  };

  const opts = consumerOpts();
  opts.deliverTo(createInbox());
  opts.deliverAll();
  opts.ackNone();
  opts.replayInstantly();
  const sub = await js.value.subscribe(subject, opts);

  subscriptions.value[subject].subscription = sub;
  resolvePendingSubscription();

  void (async () => {
    for await (const m of sub) {
      const raw = sc.decode(m.data);
      if (raw === DeletePayload) {
        // TODO
        // delete vehicles.value[''];
      } else {
        const newModel = JSON.parse(raw) as Models;
        if (raw !== JSON.stringify(state.value[newModel.id])) {
          state.value = Object.freeze({
            ...state.value,
            [newModel.id]: Object.freeze(newModel),
          });
        }
      }
    }
  })();
};

export const unsubscribe = async (subject: string) => {
  if (subscriptions.value[subject]) {
    const { pending } = subscriptions.value[subject];
    if (pending) {
      await pending;
    }
    subscriptions.value[subject]?.subscription?.unsubscribe();
    delete subscriptions.value[subject];
  }
  if (subscriptionsQueue[subject]) {
    delete subscriptionsQueue[subject];
  }
};

const processSubscriptionsQueue = async () => {
  await Promise.all(
    Object.keys(subscriptionsQueue).map(async (subject) => {
      await subscribe(subject, subscriptionsQueue[subject]);
      delete subscriptionsQueue[subject];
    }),
  );
};

export const loadApi = async () => {
  if (!natsServerUrl || typeof natsServerUrl !== 'string') {
    throw new Error('NATS_URL is invalid!');
  }

  nc = await connect({
    servers: [natsServerUrl],
    waitOnFirstConnect: true,
    maxReconnectAttempts: -1,
  });
  isConnected.value = true;
  js.value = nc.jetstream();

  await processSubscriptionsQueue();

  void (async () => {
    for await (const s of nc.status()) {
      if (s.type === Events.Disconnect) {
        isConnected.value = false;
      }
      if (s.type === Events.Reconnect) {
        isConnected.value = true;

        await processSubscriptionsQueue();
      }
    }
  })();
};
