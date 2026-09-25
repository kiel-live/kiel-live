import type { JetStreamManager } from '@nats-io/jetstream';
import type { NatsConnection, Subscription } from '@nats-io/nats-core';
import type { Ref } from 'vue';
import type { Api, Bounds, Models, Stop, Trip, Vehicle } from '~/api/types';
import { AckPolicy, DeliverPolicy, jetstreamManager, ReplayPolicy } from '@nats-io/jetstream';
import { createInbox, wsconnect } from '@nats-io/nats-core';
import Fuse from 'fuse.js';
import { computed, ref, watch } from 'vue';

import { natsServerUrl } from '~/config';

export const DeletePayload = '---';

export class NatsApi implements Api {
  isConnected = ref(false);

  private vehicles = ref<Record<string, Vehicle>>({});

  private stops = ref<Record<string, Stop>>({});

  private trips = ref<Record<string, Trip>>({});

  private subscriptions: Record<
    string,
    { subscription?: Subscription; pending?: Promise<void>; state: Ref<Record<string, Models>> }
  > = {};

  private subscriptionsQueue: Record<string, Ref<Record<string, Models>>> = {};

  private nc: NatsConnection | undefined;

  private jsm: JetStreamManager | undefined;

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
    this.jsm = await jetstreamManager(this.nc, { checkAPI: false });

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

          void this.resubscribeAll();
        }
      }
    })();
  }

  async subscribe(topic: string, state: Ref<Record<string, Models>>) {
    if (this.subscriptions[topic]) {
      return;
    }

    if (!this.isConnected.value || !this.jsm || !this.nc) {
      this.subscriptionsQueue[topic] = state;
      return;
    }

    let resolvePendingSubscription: () => void = () => {};
    this.subscriptions[topic] = {
      pending: new Promise((resolve) => {
        resolvePendingSubscription = resolve;
      }),
      state,
    };

    const inbox = createInbox();
    const sub = this.nc.subscribe(inbox);

    try {
      const streamName = await this.jsm.streams.find(topic);
      await this.jsm.consumers.add(streamName, {
        deliver_subject: inbox,
        deliver_policy: DeliverPolicy.All,
        ack_policy: AckPolicy.None,
        replay_policy: ReplayPolicy.Instant,
        filter_subject: topic,
      });
    } catch (error) {
      // Most likely the connection dropped while setting up the consumer.
      // Clean up and queue the topic so it is retried on the next reconnect.
      sub.unsubscribe();
      delete this.subscriptions[topic];
      resolvePendingSubscription();
      this.subscriptionsQueue[topic] = state;
      console.error(`Failed to subscribe to ${topic}`, error);
      return;
    }

    this.subscriptions[topic].subscription = sub;
    resolvePendingSubscription();

    void (async () => {
      for await (const m of sub) {
        const raw = m.string();
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

  async unsubscribe(topic: string) {
    if (this.subscriptions[topic]) {
      const { pending } = this.subscriptions[topic];
      if (pending) {
        await pending;
      }
      this.subscriptions[topic]?.subscription?.unsubscribe();
      delete this.subscriptions[topic];
    }
    if (this.subscriptionsQueue[topic]) {
      delete this.subscriptionsQueue[topic];
    }
  }

  private async processSubscriptionsQueue() {
    await Promise.all(
      Object.keys(this.subscriptionsQueue).map(async (topic) => {
        const state = this.subscriptionsQueue[topic];
        delete this.subscriptionsQueue[topic];
        await this.subscribe(topic, state);
      }),
    );
  }

  /**
   * Drop all active subscriptions and re-create them (including their
   * JetStream consumers) on the current connection.
   */
  private async resubscribeAll() {
    await Promise.all(
      Object.keys(this.subscriptions).map(async (topic) => {
        const entry = this.subscriptions[topic];
        if (!entry) {
          return;
        }
        // remember the state ref before waiting, so a topic unsubscribed in
        // the meantime is not re-added by accident
        const { state } = entry;
        await entry.pending;
        if (!this.subscriptions[topic]) {
          return;
        }
        this.subscriptions[topic]?.subscription?.unsubscribe();
        delete this.subscriptions[topic];
        this.subscriptionsQueue[topic] = state;
      }),
    );

    await this.processSubscriptionsQueue();
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
