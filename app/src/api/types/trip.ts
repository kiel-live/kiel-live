import { TripArrival } from './arrival';
import { GpsLocation } from './location';

export type Trip = {
  id: string;
  provider: string;
  direction: string;
  arrivals?: TripArrival[];
  path: GpsLocation[];
};
