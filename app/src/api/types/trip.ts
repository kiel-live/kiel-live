import type { TripArrival } from './arrival';
import type { GpsLocation } from './location';

export interface Trip {
  id: string;
  provider: string;
  direction: string;
  arrivals?: TripArrival[];
  path: GpsLocation[];
}
