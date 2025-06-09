import type { JetStreamClient, JetStreamSubscription, NatsConnection } from 'nats.ws';
import type { Ref } from 'vue';
import type { Api, Bounds, Models, Stop, Trip, Vehicle } from '~/api/types';
import Fuse from 'fuse.js';
import { connect, consumerOpts, createInbox, Events, StringCodec } from 'nats.ws';
import { computed, ref, watch } from 'vue';

import { natsServerUrl } from '~/config';

const sc = StringCodec();

export const DeletePayload = '---';

export class NatsApi implements Api {
  isConnected = ref(false);

  private vehicles = ref<Record<string, Vehicle>>({});

  private stops = ref<Record<string, Stop>>({});

  private trips = ref<Record<string, Trip>>({});

  private subscriptions = ref<Record<string, { subscription?: JetStreamSubscription; pending?: Promise<void> }>>({});

  private subscriptionsQueue: Record<string, Ref<Record<string, Models>>> = {};

  private nc: NatsConnection | undefined;

  js: Ref<JetStreamClient | undefined> = ref();

  constructor(autoLoad = true) {
    if (autoLoad) {
      void this.load();
    }
  }

  async load() {
    if (!natsServerUrl || typeof natsServerUrl !== 'string') {
      throw new Error('NATS_URL is invalid!');
    }

    this.nc = await connect({
      servers: [natsServerUrl],
      waitOnFirstConnect: true,
      maxReconnectAttempts: -1,
    });
    this.isConnected.value = true;
    this.js.value = this.nc.jetstream();

    await this.processSubscriptionsQueue();

    void (async () => {
      if (!this.nc) {
        throw new Error('NATS connection is not initialized');
      }

      for await (const s of this.nc.status()) {
        if (s.type === Events.Disconnect) {
          this.isConnected.value = false;
        }
        if (s.type === Events.Reconnect) {
          this.isConnected.value = true;

          await this.processSubscriptionsQueue();
        }
      }
    })();
  }

  async subscribe(subject: string, state: Ref<Record<string, Models>>) {
    if (this.subscriptions.value[subject]) {
      return;
    }

    if (!this.isConnected.value || !this.js.value) {
      this.subscriptionsQueue[subject] = state;
      return;
    }

    let resolvePendingSubscription: () => void = () => {};
    this.subscriptions.value[subject] = {
      pending: new Promise((resolve) => {
        resolvePendingSubscription = resolve;
      }),
    };

    const opts = consumerOpts();
    opts.deliverTo(createInbox());
    opts.deliverAll();
    opts.ackNone();
    opts.replayInstantly();
    const sub = await this.js.value.subscribe(subject, opts);

    this.subscriptions.value[subject].subscription = sub;
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
  }

  async unsubscribe(subject: string) {
    if (this.subscriptions.value[subject]) {
      const { pending } = this.subscriptions.value[subject];
      if (pending) {
        await pending;
      }
      this.subscriptions.value[subject]?.subscription?.unsubscribe();
      delete this.subscriptions.value[subject];
    }
    if (this.subscriptionsQueue[subject]) {
      delete this.subscriptionsQueue[subject];
    }
  }

  private async processSubscriptionsQueue() {
    await Promise.all(
      Object.keys(this.subscriptionsQueue).map(async (subject) => {
        await this.subscribe(subject, this.subscriptionsQueue[subject]);
        delete this.subscriptionsQueue[subject];
      }),
    );
  }

  useStops() {
    void this.subscribe(`data.map.stop.>`, this.stops);

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
    if (stopId.value) {
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
    if (vehicleId.value) {
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

  useSearch(query: Ref<string>, _bounds: Ref<Bounds>) {
    const { stops, loading } = this.useStops();

    const searchData = computed(() => [...Object.values(stops.value)]);
    const searchIndex = computed(
      () =>
        new Fuse(searchData.value, {
          includeScore: true,
          keys: ['name'],
          threshold: 0.4,
        }),
    );

    const results = computed(() => {
      if (query.value === '' || query.value.length < 3) {
        return [];
      }
      // limit to max 20 results
      return searchIndex.value
        .search(query.value)
        .slice(0, 20)
        .map((result) => result.item);
    });

    return { results, loading };
  }
}
