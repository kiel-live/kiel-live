import { Action } from './action';
import { StopArrival } from './arrival';
import { GpsLocation } from './location';
import { Vehicle } from './vehicle';

export type StopType =
  | 'bus-stop'
  | 'parking-spot'
  | 'ferry-stop'
  | 'train-stop'
  | 'subway-stop'
  | 'bike-stop'
  | 'tram-stop';

export type Stop = {
  id: string;
  provider: string;
  name: string;
  /**
   * @deprecated use arrivals[].type or vehicles[].type instead
   */
  type: StopType;
  routes: string[] | null; // list of routes using this stop
  alerts: string[] | null; // general alerts for this stop
  arrivals?: StopArrival[];
  location: GpsLocation;
  vehicles?: Vehicle[];
  actions?: Action[];
};
