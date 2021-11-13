import { connect, consumerOpts, createInbox, Events, StringCodec } from 'nats.ws';
import { ref } from 'vue';
import { Vehicle, Stop } from '~/api/types';

const sc = StringCodec();

export const vehicles = ref<Record<string, Vehicle>>({});
export const stops = ref<Record<string, Stop>>({});
export const isConnected = ref(false);

export const loadApi = async () => {
  const nc = await connect({
    servers: ['ws://localhost:4223', 'wss://api.kiel-live.ju60.de'],
  });
  isConnected.value = true;

  const js = nc.jetstream();

  const opts = consumerOpts();
  opts.deliverTo(createInbox());
  opts.deliverAll();
  opts.replayInstantly();
  const sub = await js.subscribe('data.>', opts);

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

  (async () => {
    for await (const m of sub) {
      // console.log(m.subject);
      const raw = sc.decode(m.data);
      if (raw === '---') {
        const id = m.subject;
        console.log('### remove', id);
        // delete vehicles.value[''];
        return;
      }

      const vehicle = JSON.parse(raw) as Vehicle;
      vehicle.location.latitude = vehicle.location.latitude / 3600000;
      vehicle.location.longitude = vehicle.location.longitude / 3600000;
      vehicles.value = Object.freeze({
        ...vehicles.value,
        [`${vehicle.provider}/${vehicle.id}`]: Object.freeze(vehicle),
      });
    }
  })();
};
