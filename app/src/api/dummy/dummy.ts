import type { Ref } from 'vue';
import type { Api, Bounds, Stop, Vehicle } from '~/api/types';
import Fuse from 'fuse.js';
import { computed, ref } from 'vue';
import { DUMMY_STOPS, DUMMY_TRIPS, DUMMY_VEHICLES } from './data';

export class DummyApi implements Api {
  isConnected = ref(true);

  useStops() {
    const stops = ref<Stop[]>(DUMMY_STOPS);

    return {
      stops,
      loading: ref(false),
      unsubscribe: async () => {},
    };
  }

  useVehicles() {
    const vehicles = ref<Vehicle[]>(DUMMY_VEHICLES);

    return {
      vehicles,
      loading: ref(false),
      unsubscribe: async () => {},
    };
  }

  useStop(stopId: Ref<string | undefined>) {
    const stop = computed(() => DUMMY_STOPS.find((s) => s.id === stopId.value) || null);

    return {
      stop,
      loading: ref(false),
      unsubscribe: async () => {},
    };
  }

  useVehicle(vehicleId: Ref<string | undefined>) {
    const vehicle = computed(() => DUMMY_VEHICLES.find((v) => v.id === vehicleId.value) || null);

    return {
      vehicle,
      loading: ref(false),
      unsubscribe: async () => {},
    };
  }

  useTrip(tripId: Ref<string | undefined>) {
    const trip = computed(() => {
      return DUMMY_TRIPS.find((t) => t.id === tripId.value) || null;
    });

    return {
      trip,
      loading: ref(false),
      unsubscribe: async () => {},
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
