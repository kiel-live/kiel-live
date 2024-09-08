import type { Action } from './action';
import type { GpsLocation } from './location';

export type VehicleType =
  | 'bus'
  | 'bike'
  | 'car'
  | 'e-scooter'
  | 'ferry'
  | 'train'
  | 'subway'
  | 'tram'
  | 'moped'
  | 'e-moped';

export interface Vehicle {
  id: string;
  provider: string;
  name: string;
  type: VehicleType;
  state: string;
  battery?: string; // in percent
  location: GpsLocation;
  tripId?: string;
  actions?: Action[];
  description?: string;
}
