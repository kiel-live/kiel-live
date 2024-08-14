import { ArrivalType, StopArrival } from './arrival';
import { GpsLocation } from './location';
import { Vehicle } from './vehicle';

export type Stop = {
  id: string;
  provider: string;
  name: string;
  type: ArrivalType; // deprecated: use arrivals[].type or vehicles[].type instead
  routes: string[] | null; // list of routes using this stop
  alerts: string[] | null; // general alerts for this stop
  arrivals?: StopArrival[];
  location: GpsLocation;
  vehicles?: Vehicle[];
};
