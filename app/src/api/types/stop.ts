import type { Action } from './action';
import type { StopDeparture } from './departure';
import type { GpsLocation } from './location';
import type { Vehicle } from './vehicle';

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
