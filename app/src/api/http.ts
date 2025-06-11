import type { Ref } from 'vue';
import type { Api, Bounds, Models as Model, Stop, Trip, Vehicle } from '~/api/types';
import { s2 } from 's2js';

import { computed, ref, watch } from 'vue';

interface WebsocketMessage<T = unknown> {
  topic: string;
  action?: 'updated' | 'deleted';
  data?: T;
  // TODO: we could extend websocket message and models with a timestamp and only update using newer messages
}

type Store<T extends Model> = Ref<Map<string, T>>;

function isModel<T extends Model>(data: unknown): data is T {
  return data !== null && typeof data === 'object' && 'id' in data;
}

function getBoundsCellIds(bounds: Bounds): string[] {
  const cellIds: string[] = [];
  // const swCell = s2.cellid.latLngToCellId(bounds.sw.lat, bounds.sw.lng, 15);
  // const neCell = s2.latLngToCellId(bounds.ne.lat, bounds.ne.lng, 15);

  // // Iterate over the cells in the bounding box
  // for (let lat = swCell.lat; lat <= neCell.lat; lat++) {
  //   for (let lng = swCell.lng; lng <= neCell.lng; lng++) {
  //     cellIds.push(s2.cellIdToString(s2.latLngToCellId(lat, lng, 15)));
  //   }
  // }

  return cellIds;
}

export class HttpApi implements Api {
  isConnected = ref(false);
  ws: WebSocket;

  private topics = new Map<string, Store<Model>>();
  private vehicles = ref<Map<string, Vehicle>>(new Map());
  private stops = ref<Map<string, Stop>>(new Map());
  private trips = ref<Map<string, Trip>>(new Map());

  constructor() {
    this.ws = new WebSocket('ws://localhost:4568/api/ws'); // TODO: replace with proper reconnecting websocket

    this.ws.addEventListener('message', (event) => {
      const message: WebsocketMessage = JSON.parse(event.data);
      const topic = message.topic;

      if (isModel(message.data)) {
        const store = this.topics.get(topic);
        if (!store) {
          console.warn(`No store found for topic: ${topic}. Probably failed to unsubscribe from topic.`);
          return;
        }

        if (message.action === 'updated') {
          store.value.set(message.data.id, message.data);
        }

        if (message.action === 'deleted') {
          store.value.delete(message.data.id);
        }
      }
    });
  }

  private async fetch<T>(url: string, options?: globalThis.RequestInit): Promise<T> {
    const response = await fetch(url, options);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json() as Promise<T>;
  }

  private async subscribe<T extends Model>(topic: string, store: Ref<Map<string, T>>) {
    this.ws.send(JSON.stringify({ action: 'subscribe', topic }));
    this.topics.set(topic, store);
  }

  private async unsubscribe(topic: string) {
    this.ws.send(JSON.stringify({ action: 'unsubscribe', topic }));
    this.topics.delete(topic);
  }

  private useMapItems<T extends Model>(itemType: string, bounds: Ref<Bounds>, store: Ref<Map<string, T>>) {
    // TODO: in case we have the same item-type + id combination, we should proxy the existing query
    const loading = ref(false);

    async function loadItems(this: HttpApi, cellIds: string[]) {
      if (loading.value === true) {
        return;
      }

      loading.value = true;
      const items = await this.fetch<T[]>(`/api/${itemType}?cells=${cellIds.join(',')}`);
      store.value.clear();
      items.forEach((item) => {
        store.value.set(item.id, item);
      });
      loading.value = false;
    }

    watch(
      bounds,
      async (newBounds, oldBounds) => {
        const newCellIds = getBoundsCellIds(newBounds);
        const oldCellIds = oldBounds ? getBoundsCellIds(oldBounds) : [];

        // Load the complete set of items for the new bounds
        await loadItems.call(this, newCellIds);

        const addedCellIds = newCellIds.filter((id) => !oldCellIds.includes(id));
        for (const cellId of addedCellIds) {
          await this.subscribe(`map.${itemType}:${cellId}`, store);
        }

        const removedCellIds = oldCellIds.filter((id) => !newCellIds.includes(id));
        for (const cellId of removedCellIds) {
          await this.unsubscribe(`map.${itemType}:${cellId}`);
        }
      },
      { immediate: true },
    );

    return {
      items: computed(() => Array.from(store.value.values())),
      loading,
      unsubscribe: async () => {
        const cellIds = getBoundsCellIds(bounds.value);
        for (const cellId of cellIds) {
          await this.unsubscribe(`map.${itemType}:${cellId}`);
        }
      },
    };
  }

  private useItem<T extends Model>(itemType: string, id: Ref<string | undefined>, store: Ref<Map<string, T>>) {
    // TODO: in case we have the same item-type + id combination, we should proxy the existing query
    const loading = ref(false);

    async function loadItem(this: HttpApi, itemId: string | undefined) {
      if (!itemId || loading.value === true) {
        return;
      }

      loading.value = true;
      const item = await this.fetch<T>(`/api/${itemType}/${itemId}`);
      store.value.set(itemId, item);
      loading.value = false;
    }

    watch(
      id,
      async (newId, oldId) => {
        if (oldId) {
          await this.unsubscribe(`${itemType}:${oldId}`);
        }
        if (newId) {
          await this.subscribe(`${itemType}:${newId}`, store);
          await loadItem.call(this, newId);
        }
      },
      { immediate: true },
    );

    return {
      item: computed(() => (id.value ? (store.value.get(id.value) ?? null) : null)),
      loading,
      unsubscribe: async () => {
        await this.unsubscribe(`${itemType}.${id.value}`);
      },
    };
  }

  useStops(bounds: Ref<Bounds>) {
    const store = ref<Map<string, Stop>>(new Map());
    const { items, loading, unsubscribe } = this.useMapItems<Stop>('stops', bounds, store);
    return {
      stops: items,
      loading,
      unsubscribe,
    };
  }

  useVehicles(bounds: Ref<Bounds>) {
    const store = ref<Map<string, Vehicle>>(new Map());
    const { items, loading, unsubscribe } = this.useMapItems<Vehicle>('vehicles', bounds, store);
    return {
      vehicles: items,
      loading,
      unsubscribe,
    };
  }

  useStop(stopId: Ref<string | undefined>) {
    const { item, loading, unsubscribe } = this.useItem<Stop>('stops', stopId, this.stops);
    return {
      stop: item,
      loading,
      unsubscribe,
    };
  }

  useVehicle(vehicleId: Ref<string | undefined>) {
    const { item, loading, unsubscribe } = this.useItem<Vehicle>('vehicles', vehicleId, this.vehicles);
    return {
      vehicle: item,
      loading,
      unsubscribe,
    };
  }

  useTrip(tripId: Ref<string | undefined>) {
    const { item, loading, unsubscribe } = this.useItem<Trip>('trips', tripId, this.trips);
    return {
      trip: item,
      loading,
      unsubscribe,
    };
  }

  useSearch(_query: Ref<string>, _bounds: Ref<Bounds>) {
    const results = ref<(Stop | Vehicle)[]>([]);
    const loading = ref(false);

    // TODO: implement search logic

    return {
      results,
      loading,
    };
  }
}
