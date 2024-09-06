import { Ref } from 'vue';

import { Bounds } from './location';
import { Stop } from './stop';
import { Trip } from './trip';
import { Vehicle } from './vehicle';

export interface Api {
  useStops(bounds: Ref<Bounds>): { stops: Ref<Stop[]>; loading: Ref<boolean>; unsubscribe: () => void | Promise<void> };

  useVehicles(bounds: Ref<Bounds>): {
    vehicles: Ref<Vehicle[]>;
    loading: Ref<boolean>;
    unsubscribe: () => void | Promise<void>;
  };

  useStop(stopId: Ref<string | undefined>): {
    stop: Ref<Stop | null>;
    loading: Ref<boolean>;
    unsubscribe: () => void | Promise<void>;
  };

  useVehicle(vehicleId: Ref<string | undefined>): {
    vehicle: Ref<Vehicle | null>;
    loading: Ref<boolean>;
    unsubscribe: () => void | Promise<void>;
  };

  useTrip(tripId: Ref<string | undefined>): {
    trip: Ref<Trip | null>;
    loading: Ref<boolean>;
    unsubscribe: () => void | Promise<void>;
  };

  get isConnected(): Ref<boolean>;
}
