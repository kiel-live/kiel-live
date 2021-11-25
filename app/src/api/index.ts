import { connect, consumerOpts, createInbox, Events, JetStreamClient, NatsConnection, StringCodec } from 'nats.ws';
import { Ref, ref } from 'vue';
import { Vehicle, Stop, Models } from '~/api/types';

const sc = StringCodec();

export const vehicles = ref<Record<string, Vehicle>>({});
export const stops = ref<Record<string, Stop>>({});
export const isConnected = ref(false);

const subscriptions = ref<string[]>([]);
let nc: NatsConnection;
let js: JetStreamClient;

export const subscribe = async (subject: string, state: Ref<Record<string, Models>>) => {
  if (subscriptions.value.includes(subject)) {
    return;
  }

  subscriptions.value.push(subject);

  const opts = consumerOpts();
  opts.deliverTo(createInbox());
  opts.deliverAll();
  opts.ackNone();
  opts.replayInstantly();
  const sub = await js.subscribe(subject, opts);

  (async () => {
    for await (const m of sub) {
      const raw = sc.decode(m.data);
      if (raw === '---') {
        const id = m.subject;
        console.log('### remove', id);
        // delete vehicles.value[''];
        continue;
      }

      const newModel = JSON.parse(raw) as Vehicle;
      state.value = Object.freeze({
        ...state.value,
        [`${newModel.provider}/${newModel.id}`]: Object.freeze(newModel),
      });
    }
  })();
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
