import { Ref } from 'vue';

import { Bounds } from './location';
import { Stop } from './stop';
import { Trip } from './trip';
import { Vehicle } from './vehicle';

export interface Api {
  // getStops(bounds: Bounds): Stop[];
  // onStopsUpdated(bounds: Bounds, cb: (stop: Stop) => void): void;
  useStops(bounds: Ref<Bounds>): { stops: Ref<Stop[]>; loading: Ref<boolean>; unsubscribe: () => void | Promise<void> };

  // getVehicles(bounds: Bounds): Vehicle[];
  // onVehiclesUpdated(bounds: Bounds, cb: (vehicle: Vehicle) => void): void;
  useVehicles(bounds: Ref<Bounds>): {
    vehicles: Ref<Vehicle[]>;
    loading: Ref<boolean>;
    unsubscribe: () => void | Promise<void>;
  };

  // getStop(stopId: string): Stop | null;
  // onStopUpdated(stopId: string, cb: (stop: Stop) => void): void;
  useStop(stopId: Ref<string | undefined>): {
    stop: Ref<Stop | null>;
    loading: Ref<boolean>;
    unsubscribe: () => void | Promise<void>;
  };

  // getVehicle(vehicleId: string): Vehicle | null;
  // onVehicleUpdated(vehicleId: string, cb: (vehicle: Vehicle) => void): void;
  useVehicle(vehicleId: Ref<string | undefined>): {
    vehicle: Ref<Vehicle | null>;
    loading: Ref<boolean>;
    unsubscribe: () => void | Promise<void>;
  };

  // getTrip(tripId: string): Trip | null;
  // onTripUpdated(tripId: string, cb: (trip: Trip) => void): void;
  useTrip(tripId: Ref<string | undefined>): {
    trip: Ref<Trip | null>;
    loading: Ref<boolean>;
    unsubscribe: () => void | Promise<void>;
  };

  get isConnected(): Ref<boolean>;
}
