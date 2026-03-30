import type { DepartureState } from './departure';
import type { GpsLocation } from './location';

export interface Trip {
  id: string;
  provider: string;
  direction: string;
  departures?: TripDeparture[];
  path?: GpsLocation[];
}

export interface TripDeparture {
  id: string;
  name: string;
  state: DepartureState;
  planned: string;
  actual?: string;
}
