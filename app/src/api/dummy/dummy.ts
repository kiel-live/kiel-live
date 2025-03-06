import type { Ref } from 'vue';
import type { Api, Stop, Vehicle } from '~/api/types';
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
}
