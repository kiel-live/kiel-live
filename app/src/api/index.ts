import {
  connect,
  consumerOpts,
  createInbox,
  Events,
  JetStreamClient,
  JetStreamSubscription,
  NatsConnection,
  StringCodec,
} from 'nats.ws';
import { Ref, ref } from 'vue';

import { Models, Stop, Trip, Vehicle } from '~/api/types';
import { natsServerUrl } from '~/config';

const sc = StringCodec();

export const DeletePayload = '---';

export const vehicles = ref<Record<string, Vehicle>>({});
export const stops = ref<Record<string, Stop>>({});
export const trips = ref<Record<string, Trip>>({});
export const isConnected = ref(false);

const subscriptions = ref<Record<string, JetStreamSubscription | 'pending'>>({});
const subscriptionsQueue: Record<string, Ref<Record<string, Models>>> = {};

let nc: NatsConnection | undefined;
export const js: Ref<JetStreamClient | undefined> = ref();

export const subscribe = async (subject: string, state: Ref<Record<string, Models>>) => {
  if (subscriptions.value[subject]) {
    return;
  }
  subscriptions.value[subject] = 'pending';

  if (!isConnected.value || !js.value) {
    subscriptionsQueue[subject] = state;
    return;
  }

  const opts = consumerOpts();
  opts.deliverTo(createInbox());
  opts.deliverAll();
  opts.ackNone();
  opts.replayInstantly();
  const sub = await js.value.subscribe(subject, opts);
  if (subscriptions.value[subject] !== 'pending') {
    sub.unsubscribe();
    return;
  }
  subscriptions.value[subject] = sub;

  void (async () => {
    // eslint-disable-next-line no-restricted-syntax
    for await (const m of sub) {
      const raw = sc.decode(m.data);
      if (raw === DeletePayload) {
        // TODO
        // delete vehicles.value[''];
      } else {
        const newModel = JSON.parse(raw) as Models;
        // eslint-disable-next-line no-param-reassign
        state.value = Object.freeze({
          ...state.value,
          [newModel.id]: Object.freeze(newModel),
        });
      }
    }
  })();
};

export const unsubscribe = async (subject: string) => {
  if (subscriptions.value[subject]) {
    const subscription = subscriptions.value[subject];
    if (subscription !== 'pending') {
      subscription.unsubscribe();
    }
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
    // eslint-disable-next-line no-restricted-syntax
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
