import type { Ref } from 'vue';
import type { Api } from './types';
import type { Bounds } from './types/location';
import type { Stop } from './types/stop';
import type { Trip } from './types/trip';
import type { Vehicle } from './types/vehicle';
import { watchDebounced } from '@vueuse/core';
import { computed, ref, watch } from 'vue';
import { hubServerUrl } from '~/config';
import { Rpc } from './rpc';

interface Envelope {
  rid?: string;
  method: string;
  topic?: string;
  data?: unknown;
}

export class HubApi implements Api {
  isConnected = ref(false);

  private ws: WebSocket | undefined;
  private stops = ref<Record<string, Stop>>({});
  private vehicles = ref<Record<string, Vehicle>>({});
  private trips = ref<Record<string, Trip>>({});

  // Map kind subscriptions: "map.stop" / "map.vehicle" → active subscriber count.
  // The server only streams map data for kinds that have been explicitly subscribed.
  private mapKindSubCounts = new Map<string, number>();

  // Last sent viewport bounds — replayed on reconnect and when a map kind is
  // subscribed after a viewport has already been received.
  private currentBounds: Bounds | undefined;

  // Detail subscriptions: topic → count of active callers
  private detailSubCounts = new Map<string, number>();

  private rpc = new Rpc();

  constructor(autoLoad = true) {
    if (autoLoad) {
      void this.connect();
    }
  }

  private connect() {
    if (!hubServerUrl) {
      throw new Error('VITE_HUB_URL is not set');
    }
    const url = hubServerUrl.replace(/^http/, 'ws').replace(/\/$/, '');
    const ws = new WebSocket(`${url}/ws/client`);
    this.ws = ws;

    ws.addEventListener('open', () => {
      this.isConnected.value = true;
      // Clear stale map state — fresh snapshots arrive after re-subscribing.
      this.stops.value = {};
      this.vehicles.value = {};
      // Replay map kind subscriptions before viewport so the server is ready.
      for (const [topic] of this.mapKindSubCounts) {
        this.sendEnvelope({ method: 'subscribe', topic });
      }
      if (this.currentBounds) {
        this.sendViewport(this.currentBounds);
      }
      // Replay detail subscriptions.
      for (const [topic] of this.detailSubCounts) {
        this.sendEnvelope({ method: 'subscribe', topic });
      }
    });

    ws.addEventListener('close', () => {
      this.isConnected.value = false;
      setTimeout(() => this.connect(), 3000);
    });

    ws.addEventListener('message', (ev) => {
      try {
        const env = JSON.parse(ev.data as string) as Envelope;
        this.handleEnvelope(env);
      } catch (e) {
        console.error('[hub] parse error', e);
      }
    });
  }

  private sendEnvelope(env: Envelope) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(env));
    }
  }

  private sendRequest(topic: string, data: unknown, timeoutMs = 5000): Promise<unknown> {
    return this.rpc.request(this.sendEnvelope.bind(this), topic, data, timeoutMs);
  }

  private sendViewport(bounds: Bounds) {
    this.sendEnvelope({ method: 'viewport', data: bounds });
  }

  private handleEnvelope(env: Envelope) {
    switch (env.method) {
      case 'snapshot': {
        // Snapshots are per-cell and must be merged, not replaced.
        // Deletes from the server clean up entities that have left the viewport.
        const items = env.data as (Stop | Vehicle)[];
        const topic = env.topic ?? '';
        if (topic === 'map.stop') {
          const next = { ...this.stops.value };
          for (const item of items) {
            next[(item as Stop).id] = item as Stop;
          }
          this.stops.value = next;
        } else if (topic === 'map.vehicle') {
          const next = { ...this.vehicles.value };
          for (const item of items) {
            next[(item as Vehicle).id] = item as Vehicle;
          }
          this.vehicles.value = next;
        }
        break;
      }

      case 'update': {
        const topic = env.topic ?? '';
        if (topic.startsWith('map.stop.')) {
          const stop = env.data as Stop;
          this.stops.value = { ...this.stops.value, [stop.id]: stop };
        } else if (topic.startsWith('map.vehicle.')) {
          const vehicle = env.data as Vehicle;
          this.vehicles.value = { ...this.vehicles.value, [vehicle.id]: vehicle };
        } else if (topic.startsWith('stop.')) {
          const stop = env.data as Stop;
          this.stops.value = { ...this.stops.value, [stop.id]: stop };
        } else if (topic.startsWith('vehicle.')) {
          const vehicle = env.data as Vehicle;
          this.vehicles.value = { ...this.vehicles.value, [vehicle.id]: vehicle };
        } else if (topic.startsWith('trip.')) {
          const trip = env.data as Trip;
          this.trips.value = { ...this.trips.value, [trip.id]: trip };
        }
        break;
      }

      case 'delete': {
        const topic = env.topic ?? '';
        if (topic.startsWith('map.stop.')) {
          const id = topic.slice('map.stop.'.length);
          const next = { ...this.stops.value };
          delete next[id];
          this.stops.value = next;
        } else if (topic.startsWith('map.vehicle.')) {
          const id = topic.slice('map.vehicle.'.length);
          const next = { ...this.vehicles.value };
          delete next[id];
          this.vehicles.value = next;
        }
        break;
      }

      case 'result':
      case 'error': {
        if (env.rid) this.rpc.resolve(env.rid, env.data);
        break;
      }
    }
  }

  private subscribeMapKind(topic: string) {
    const count = this.mapKindSubCounts.get(topic) ?? 0;
    this.mapKindSubCounts.set(topic, count + 1);
    if (count === 0) {
      // Clear stale state before the server sends fresh snapshots.
      if (topic === 'map.stop') this.stops.value = {};
      if (topic === 'map.vehicle') this.vehicles.value = {};
      this.sendEnvelope({ method: 'subscribe', topic });
    }
  }

  private unsubscribeMapKind(topic: string) {
    const count = this.mapKindSubCounts.get(topic) ?? 0;
    if (count <= 1) {
      this.mapKindSubCounts.delete(topic);
      this.sendEnvelope({ method: 'unsubscribe', topic });
    } else {
      this.mapKindSubCounts.set(topic, count - 1);
    }
  }

  private subscribeDetail(topic: string) {
    const count = this.detailSubCounts.get(topic) ?? 0;
    this.detailSubCounts.set(topic, count + 1);
    if (count === 0) {
      this.sendEnvelope({ method: 'subscribe', topic });
    }
  }

  private unsubscribeDetail(topic: string) {
    const count = this.detailSubCounts.get(topic) ?? 0;
    if (count <= 1) {
      this.detailSubCounts.delete(topic);
      this.sendEnvelope({ method: 'unsubscribe', topic });
    } else {
      this.detailSubCounts.set(topic, count - 1);
    }
  }

  private sendBoundsAsViewport(b: Bounds | undefined) {
    if (!b) return;
    this.currentBounds = b;
    this.sendViewport(b);
  }

  useStops(bounds: Ref<Bounds | undefined>) {
    this.subscribeMapKind('map.stop');
    const stopWatcher = watch(bounds, (b) => this.sendBoundsAsViewport(b), { immediate: true });
    return {
      stops: computed(() => Object.values(this.stops.value)),
      loading: ref(false),
      unsubscribe: () => {
        stopWatcher();
        this.unsubscribeMapKind('map.stop');
      },
    };
  }

  useVehicles(bounds: Ref<Bounds | undefined>) {
    this.subscribeMapKind('map.vehicle');
    const boundsWatcher = watch(bounds, (b) => this.sendBoundsAsViewport(b), { immediate: true });
    return {
      vehicles: computed(() => Object.values(this.vehicles.value)),
      loading: ref(false),
      unsubscribe: () => {
        boundsWatcher();
        this.unsubscribeMapKind('map.vehicle');
      },
    };
  }

  useStop(stopId: Ref<string | undefined>) {
    let prevTopic: string | undefined;

    const subscribe = (id: string | undefined) => {
      if (prevTopic) {
        this.unsubscribeDetail(prevTopic);
        prevTopic = undefined;
      }
      if (id) {
        const topic = `stop.${id}`;
        this.subscribeDetail(topic);
        prevTopic = topic;
      }
    };

    subscribe(stopId.value);
    const stopWatcher = watch(stopId, subscribe);

    return {
      stop: computed(() => (stopId.value ? (this.stops.value[stopId.value] ?? null) : null)),
      loading: ref(false),
      unsubscribe: () => {
        stopWatcher();
        if (prevTopic) {
          this.unsubscribeDetail(prevTopic);
          prevTopic = undefined;
        }
      },
    };
  }

  useVehicle(vehicleId: Ref<string | undefined>) {
    let prevTopic: string | undefined;

    const subscribe = (id: string | undefined) => {
      if (prevTopic) {
        this.unsubscribeDetail(prevTopic);
        prevTopic = undefined;
      }
      if (id) {
        const topic = `vehicle.${id}`;
        this.subscribeDetail(topic);
        prevTopic = topic;
      }
    };

    subscribe(vehicleId.value);
    const stopWatcher = watch(vehicleId, subscribe);

    return {
      vehicle: computed(() => (vehicleId.value ? (this.vehicles.value[vehicleId.value] ?? null) : null)),
      loading: ref(false),
      unsubscribe: () => {
        stopWatcher();
        if (prevTopic) {
          this.unsubscribeDetail(prevTopic);
          prevTopic = undefined;
        }
      },
    };
  }

  useTrip(tripId: Ref<string | undefined>) {
    let prevTopic: string | undefined;

    const subscribe = (id: string | undefined) => {
      if (prevTopic) {
        this.unsubscribeDetail(prevTopic);
        prevTopic = undefined;
      }
      if (id) {
        const topic = `trip.${id}`;
        this.subscribeDetail(topic);
        prevTopic = topic;
      }
    };

    subscribe(tripId.value);
    const stopWatcher = watch(tripId, subscribe);

    return {
      trip: computed(() => (tripId.value ? (this.trips.value[tripId.value] ?? null) : null)),
      loading: ref(false),
      unsubscribe: () => {
        stopWatcher();
        if (prevTopic) {
          this.unsubscribeDetail(prevTopic);
          prevTopic = undefined;
        }
      },
    };
  }

  useSearch(query: Ref<string>, bounds: Ref<Bounds>) {
    const results = ref<(Stop | Vehicle)[]>([]);
    const loading = ref(false);
    let generation = 0;

    watchDebounced(
      [query, bounds] as const,
      async ([q, b]) => {
        const gen = ++generation;

        if (q.length < 3) {
          results.value = [];
          loading.value = false;
          return;
        }

        loading.value = true;
        try {
          const data = await this.sendRequest('search', { query: q, bounds: b, limit: 20 });
          if (gen === generation) {
            results.value = (data as (Stop | Vehicle)[]) ?? [];
          }
        } catch {
          if (gen === generation) results.value = [];
        } finally {
          if (gen === generation) loading.value = false;
        }
      },
      { immediate: true, debounce: 200 },
    );

    return { results, loading };
  }
}
