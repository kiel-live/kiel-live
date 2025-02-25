import { rpcClient } from 'typed-rpc';
import { websocketTransport } from 'typed-rpc/ws';
import { computed, ref, watch } from 'vue';
import type { Ref } from 'vue';

import type { Api, Models, Stop, Trip, Vehicle } from '~/api/types';
import { ReconnectingWebSocket } from './ws';

type Client = ReturnType<typeof rpcClient>;

class RPCClient {
  client: Client;

  constructor(url: string) {
    this.client = rpcClient<{
      subscribe: (topic: string) => void;
      unsubscribe: (topic: string) => void;
    }>({
      transport: websocketTransport({
        url,
      }),
    });
  }

  async subscribe(topic: string, cb: (data: T) => void) {
    return this.subscribe(topic);
  }

  async unsubscribe(topic: string) {
    return this.send('unsubscribe', { topic });
  }
}

export class RPCApi implements Api {
  isConnected = ref(false);

  socket: ReconnectingWebSocket;
  client: RPCClient;

  constructor(url: string) {
    this.socket = new ReconnectingWebSocket(url);
    this.socket.on('connected', () => {
      this.isConnected.value = true;
    });
    this.socket.on('disconnected', () => {
      this.isConnected.value = false;
    });
    this.client = new RPCClient((payload, _clientParams) => {
      this.socket.send(JSON.stringify(payload));
    });
    this.socket.on('message', (event) => {
      this.client.receive(JSON.parse(event.data.toString()));
    });
    this.socket.on('close', (event) => {
      this.client.rejectAllPendingRequests(`Connection is closed (${event.reason}).`);
    });

    const c = rpcClient({
      transport: websocketTransport({
        url,
      }),
    });
  }

  async load() {}

  useStops() {
    return {
      stops: computed(() => Object.values(this.stops.value)),
      loading: ref(false),
      unsubscribe: async () => {
        await this.unsubscribe(`data.map.stop.>`);
      },
    };
  }

  useVehicles() {
    void this.subscribe(`data.map.vehicle.>`, this.vehicles);

    return {
      vehicles: computed(() => Object.values(this.vehicles.value)),
      loading: ref(false),
      unsubscribe: async () => {
        await this.unsubscribe(`data.map.vehicle.>`);
      },
    };
  }

  useStop(stopId: Ref<string | undefined>) {
    if (stopId) {
      void this.subscribe(`data.map.stop.${stopId.value}`, this.stops);
    }

    watch(stopId, async (newId, oldId) => {
      if (oldId) {
        await this.unsubscribe(`data.map.stop.${oldId}`);
      }
      if (newId) {
        await this.subscribe(`data.map.stop.${newId}`, this.stops);
      }
    });

    return {
      stop: computed(() => (stopId.value ? (this.stops.value[stopId.value] ?? null) : null)),
      loading: ref(false),
      unsubscribe: async () => {
        await this.unsubscribe(`data.map.stop.${stopId.value}`);
      },
    };
  }

  useVehicle(vehicleId: Ref<string | undefined>) {
    if (vehicleId) {
      void this.subscribe(`data.map.vehicle.${vehicleId.value}`, this.vehicles);
    }

    watch(vehicleId, async (newId, oldId) => {
      if (oldId) {
        await this.unsubscribe(`data.map.vehicle.${oldId}`);
      }
      if (newId) {
        await this.subscribe(`data.map.vehicle.${newId}`, this.vehicles);
      }
    });

    return {
      vehicle: computed(() => (vehicleId.value ? (this.vehicles.value[vehicleId.value] ?? null) : null)),
      loading: ref(false),
      unsubscribe: async () => {
        await this.unsubscribe(`data.map.vehicle.${vehicleId.value}`);
      },
    };
  }

  useTrip(tripId: Ref<string | undefined>) {
    if (tripId.value) {
      void this.subscribe(`data.map.trip.${tripId.value}`, this.trips);
    }

    watch(tripId, async (newId, oldId) => {
      if (oldId) {
        await this.unsubscribe(`data.map.trip.${oldId}`);
      }
      if (newId) {
        await this.subscribe(`data.map.trip.${newId}`, this.trips);
      }
    });

    return {
      trip: computed(() => (tripId.value ? (this.trips.value[tripId.value] ?? null) : null)),
      loading: ref(false),
      unsubscribe: async () => {
        await this.unsubscribe(`data.map.trip.${tripId.value}`);
      },
    };
  }
}
