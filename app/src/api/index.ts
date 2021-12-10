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

const sc = StringCodec();

export const DeletePayload = '---';

export const vehicles = ref<Record<string, Vehicle>>({});
export const stops = ref<Record<string, Stop>>({});
export const trips = ref<Record<string, Trip>>({});
export const isConnected = ref(false);

const subscriptions = ref<Record<string, JetStreamSubscription>>({});
let nc: NatsConnection | undefined;
let js: JetStreamClient | undefined;

export const subscribe = async (subject: string, state: Ref<Record<string, Models>>) => {
  if (subscriptions.value[subject]) {
    return;
  }

  if (!js) {
    // TODO: start retrying
    throw new Error('Connect before subscribing');
  }

  const opts = consumerOpts();
  opts.deliverTo(createInbox());
  opts.deliverAll();
  opts.ackNone();
  opts.replayInstantly();
  const sub = await js.subscribe(subject, opts);
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
  if (!subscriptions.value[subject]) {
    return;
  }
  subscriptions.value[subject].unsubscribe();
  delete subscriptions.value[subject];
};

export const loadApi = async () => {
  nc = await connect({
    servers: ['ws://localhost:4223'],
  });
  isConnected.value = true;

  js = nc.jetstream();

  void (async () => {
    // eslint-disable-next-line no-restricted-syntax
    for await (const s of nc.status()) {
      if (s.type === Events.Disconnect) {
        isConnected.value = false;
      }
      if (s.type === Events.Reconnect) {
        isConnected.value = true;
      }
    }
  })();
};
