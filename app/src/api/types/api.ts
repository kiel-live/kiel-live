import type { Ref } from 'vue';

import type { Bounds } from './location';
import type { Stop } from './stop';
import type { Trip } from './trip';
import type { Vehicle } from './vehicle';

export interface Api {
  useStops: (bounds: Ref<Bounds | undefined>) => {
    stops: Ref<Stop[]>;
    loading: Ref<boolean>;
    unsubscribe: () => void | Promise<void>;
  };

  useVehicles: (bounds: Ref<Bounds | undefined>) => {
    vehicles: Ref<Vehicle[]>;
    loading: Ref<boolean>;
    unsubscribe: () => void | Promise<void>;
  };

  useStop: (stopId: Ref<string | undefined>) => {
    stop: Ref<Stop | null>;
    loading: Ref<boolean>;
    unsubscribe: () => void | Promise<void>;
  };

  useVehicle: (vehicleId: Ref<string | undefined>) => {
    vehicle: Ref<Vehicle | null>;
    loading: Ref<boolean>;
    unsubscribe: () => void | Promise<void>;
  };

  useTrip: (tripId: Ref<string | undefined>) => {
    trip: Ref<Trip | null>;
    loading: Ref<boolean>;
    unsubscribe: () => void | Promise<void>;
  };

  useSearch: (
    query: Ref<string>,
    bounds: Ref<Bounds>,
  ) => {
    results: Ref<(Stop | Vehicle)[]>;
    loading: Ref<boolean>;
  };

  get isConnected(): Ref<boolean>;
}
