import { connect, consumerOpts, createInbox, Events, JetStreamClient, JetStreamSubscription, NatsConnection, StringCodec } from 'nats.ws';
import { Ref, ref } from 'vue';
import { Vehicle, Stop, Models, Trip } from '~/api/types';

const sc = StringCodec();

export const DeletePayload = '---'

export const vehicles = ref<Record<string, Vehicle>>({});
export const stops = ref<Record<string, Stop>>({});
export const trips = ref<Record<string, Trip>>({});
export const isConnected = ref(false);

const subscriptions = ref<Record<string, JetStreamSubscription>>({});
let nc: NatsConnection;
let js: JetStreamClient;

export const subscribe = async (subject: string, state: Ref<Record<string, Models>>) => {
  if (subscriptions.value[subject]) {
    return;
  }

  const opts = consumerOpts();
  opts.deliverTo(createInbox());
  opts.deliverAll();
  opts.ackNone();
  opts.replayInstantly();
  const sub = await js.subscribe(subject, opts);
  subscriptions.value[subject] = sub;

  (async () => {
    for await (const m of sub) {
      const raw = sc.decode(m.data);
      if (raw === DeletePayload) {
        const id = m.subject;
        console.log('### remove', id);
        // delete vehicles.value[''];
        continue;
      }

      const newModel = JSON.parse(raw);
      state.value = Object.freeze({
        ...state.value,
        [newModel.id]: Object.freeze(newModel),
      });
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

  (async () => {
    console.info(`connected ${nc.getServer()}`);
    for await (const s of nc.status()) {
      if (s.type === Events.Disconnect) {
        isConnected.value = false;
      }
      if (s.type === Events.Reconnect) {
        isConnected.value = true;
      }
      console.info(`${s.type}: ${s.data}`);
    }
  })();
};
