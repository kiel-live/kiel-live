import type { JetStreamClient, JetStreamManager, JetStreamSubscription } from '@nats-io/jetstream';
import type { NatsConnection } from '@nats-io/nats-core';
import type { Ref } from 'vue';
import type { Api, Models, Stop, Trip, Vehicle } from '~/api/types';
import { jetstream, jetstreamManager } from '@nats-io/jetstream';
import { consumerOpts, createInbox, wsconnect } from '@nats-io/nats-core';
import { computed, ref, watch } from 'vue';
import { natsServerUrl } from '~/config';

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
  jsm: Ref<JetStreamManager | undefined> = ref();

  constructor(autoLoad = true) {
    if (autoLoad) {
      void this.load();
    }
  }

  async load() {
    if (!natsServerUrl || typeof natsServerUrl !== 'string') {
      throw new Error('NATS_URL is invalid!');
    }

    this.nc = await wsconnect({
      servers: [natsServerUrl],
      waitOnFirstConnect: true,
      maxReconnectAttempts: -1,
    });
    this.isConnected.value = true;
    this.js.value = jetstream(this.nc);
    this.jsm.value = await jetstreamManager(this.nc);

    await this.processSubscriptionsQueue();

    void (async () => {
      if (!this.nc) {
        throw new Error('NATS connection is not initialized');
      }

      for await (const s of this.nc.status()) {
        if (s.type === 'disconnect') {
          this.isConnected.value = false;
        }
        if (s.type === 'reconnect') {
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
    // const sub = await this.js.value.subscribe(subject, opts);
    const consumer = await this.jsm.value?.consumers.add(subject, { ack_policy: 'none', durable_name: 'A' });

    const c2 = await this.js.value.consumers.get(subject, consumer?.name ?? '');
    const iter = await c2.consume();
    // this.subscriptions.value[subject].subscription = sub;
    resolvePendingSubscription();

    void (async () => {
      for await (const m of iter) {
        const raw = m.data.toString();
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
