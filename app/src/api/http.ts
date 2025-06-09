import type { Ref } from 'vue';
import type { Api, Bounds, Models, Stop, Trip, Vehicle } from '~/api/types';

import { computed, ref, watch } from 'vue';

export class HttpApi implements Api {
  isConnected = ref(false);

  private vehicles = ref<Map<string, Vehicle>>(new Map());

  private stops = ref<Map<string, Stop>>(new Map());

  private trips = ref<Map<string, Trip>>(new Map());

  constructor(autoLoad = true) {
    if (autoLoad) {
      void this.load();
    }
  }

  async load() {
    // Validate http endpoint
    this.isConnected.value = true;
  }

  private async fetch<T>(url: string): Promise<T> {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json() as Promise<T>;
  }

  private async subscribe<T extends Models>(subject: string, store: Ref<Map<string, T>>) {
    // TODO: Implement the logic to subscribe to the subject
  }

  private async unsubscribe(subject: string) {
    // TODO: Implement the logic to unsubscribe from the subject
  }

  useStops() {
    const loading = ref(true);

    async function loadStops(this: HttpApi) {
      loading.value = true;
      const stops = await this.fetch<Stop[]>(`/api/stops`);
      stops.forEach((stop) => {
        this.stops.value.set(stop.id, stop);
      });
      loading.value = false;
    }

    void this.subscribe(`data.map.stop.>`, this.stops);
    void loadStops.call(this);

    // TODO: watch bounds

    return {
      stops: computed(() => Array.from(this.stops.value.values())),
      loading,
      unsubscribe: async () => {
        await this.unsubscribe(`data.map.stop.>`);
      },
    };
  }

  useVehicles() {
    const loading = ref(true);

    async function loadVehicles(this: HttpApi) {
      loading.value = true;
      const vehicles = await this.fetch<Vehicle[]>(`/api/vehicles`);
      vehicles.forEach((vehicle) => {
        this.vehicles.value.set(vehicle.id, vehicle);
      });
      loading.value = false;
    }

    void this.subscribe(`data.map.vehicle.>`, this.vehicles);
    void loadVehicles.call(this);

    // TODO: watch bounds

    return {
      vehicles: computed(() => Array.from(this.vehicles.value.values())),
      loading,
      unsubscribe: async () => {
        await this.unsubscribe(`data.map.vehicle.>`);
      },
    };
  }

  useStop(stopId: Ref<string | undefined>) {
    const loading = ref(true);

    async function loadStop(this: HttpApi, id: string) {
      loading.value = true;
      if (!id) {
        loading.value = false;
        return;
      }

      const stop = await this.fetch<Stop>(`/api/stops/${id}`);
      this.stops.value.set(stop.id, stop);
      loading.value = false;
    }

    watch(
      stopId,
      async (newId, oldId) => {
        if (oldId) {
          await this.unsubscribe(`data.map.stop.${oldId}`);
        }
        if (newId) {
          await this.subscribe(`data.map.stop.${newId}`, this.stops);
          await loadStop.call(this, newId);
        }
      },
      { immediate: true },
    );

    return {
      stop: computed(() => (stopId.value ? (this.stops.value.get(stopId.value) ?? null) : null)),
      loading,
      unsubscribe: async () => {
        await this.unsubscribe(`data.map.stop.${stopId.value}`);
      },
    };
  }

  useVehicle(vehicleId: Ref<string | undefined>) {
    const loading = ref(true);

    async function loadVehicle(this: HttpApi, id: string) {
      loading.value = true;
      if (!id) {
        loading.value = false;
        return;
      }

      const vehicle = await this.fetch<Vehicle>(`/api/vehicles/${id}`);
      this.vehicles.value.set(vehicle.id, vehicle);
      loading.value = false;
    }

    watch(
      vehicleId,
      async (newId, oldId) => {
        if (oldId) {
          await this.unsubscribe(`data.map.vehicle.${oldId}`);
        }
        if (newId) {
          await this.subscribe(`data.map.vehicle.${newId}`, this.vehicles);
          await loadVehicle.call(this, newId);
        }
      },
      { immediate: true },
    );

    return {
      vehicle: computed(() => (vehicleId.value ? (this.vehicles.value.get(vehicleId.value) ?? null) : null)),
      loading,
      unsubscribe: async () => {
        await this.unsubscribe(`data.map.vehicle.${vehicleId.value}`);
      },
    };
  }

  useTrip(tripId: Ref<string | undefined>) {
    const loading = ref(true);

    async function loadTrip(this: HttpApi, id: string) {
      loading.value = true;
      if (!id) {
        loading.value = false;
        return;
      }

      const trip = await this.fetch<Trip>(`/api/trips/${id}`);
      this.trips.value.set(trip.id, trip);
      loading.value = false;
    }

    watch(
      tripId,
      async (newId, oldId) => {
        if (oldId) {
          await this.unsubscribe(`data.map.trip.${oldId}`);
        }
        if (newId) {
          await this.subscribe(`data.map.trip.${newId}`, this.trips);
          await loadTrip.call(this, newId);
        }
      },
      { immediate: true },
    );

    return {
      trip: computed(() => (tripId.value ? (this.trips.value.get(tripId.value) ?? null) : null)),
      loading,
      unsubscribe: async () => {
        await this.unsubscribe(`data.map.trip.${tripId.value}`);
      },
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
