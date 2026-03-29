import type { TripDeparture } from './departure';
import type { GpsLocation } from './location';

export interface Trip {
  id: string;
  provider: string;
  direction: string;
  departures?: TripDeparture[];
  path?: GpsLocation[];
}
