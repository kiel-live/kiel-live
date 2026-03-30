import type { Action } from './action';
import type { DepartureState } from './departure';
import type { GpsLocation } from './location';
import type { Vehicle, VehicleType } from './vehicle';

export type StopType =
  | 'bus-stop'
  | 'parking-spot'
  | 'ferry-stop'
  | 'train-stop'
  | 'subway-stop'
  | 'bike-stop'
  | 'tram-stop';

export interface Stop {
  id: string;
  provider: string;
  name: string;
  /**
   * @deprecated use departures[].type or vehicles[].type instead
   */
  type: StopType;
  routes?: string[] | null; // list of routes using this stop
  alerts?: string[] | null; // general alerts for this stop
  departures?: StopDeparture[];
  location: GpsLocation;
  vehicles?: Vehicle[];
  actions?: Action[];
}

export interface StopDeparture {
  name: string;
  type: VehicleType;
  vehicleId: string;
  tripId: string;
  routeId: string;
  routeName: string;
  direction: string;
  state: DepartureState;
  planned: string; // iso datetime string
  actual: string; // iso datetime string
  platform: string;
}
