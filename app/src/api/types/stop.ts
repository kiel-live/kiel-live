import { StopArrival } from './arrival';
import { GpsLocation } from './location';

export type StopType = 'bus-stop' | 'parking-spot' | 'ferry-stop' | 'train-stop' | 'subway-stop' | 'bike-stop';

export type Stop = {
  id: string;
  provider: string;
  name: string;
  type: StopType;
  routes: string[]; // list of routes using this stop
  alerts: string[]; // general alerts for this stop
  arrivals: StopArrival[];
  location: GpsLocation;
};
