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
  opts.ackNone()
  opts.replayInstantly();
  const sub = await js.subscribe('data.map.vehicle.>', opts);

  const optsStops = consumerOpts();
  optsStops.deliverTo(createInbox());
  optsStops.deliverAll();
  optsStops.ackNone()
  optsStops.replayInstantly();
  const subStops = await js.subscribe('data.map.stop.>', optsStops);

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
        continue;
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

  (async () => {
    for await (const m of subStops) {
      // console.log(m.subject);
      const raw = sc.decode(m.data);
      if (raw === '---') {
        const id = m.subject;
        console.log('### remove', id);
        // delete vehicles.value[''];
        continue;
      }

      const stop = JSON.parse(raw) as Stop;
      stop.location.latitude = stop.location.latitude / 3600000;
      stop.location.longitude = stop.location.longitude / 3600000;
      stops.value = Object.freeze({
        ...stops.value,
        [`${stop.provider}/${stop.id}`]: Object.freeze(stop),
      });
    }
  })();
};
