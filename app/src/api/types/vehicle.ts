import { GpsLocation } from './location';

export type VehicleType = 'bus' | 'bike' | 'car' | 'e-scooter' | 'ferry' | 'train' | 'subway' | 'tram';

export type Vehicle = {
  id: string;
  provider: string;
  name: string;
  type: VehicleType;
  state: string;
  battery: string; // in percent
  location: GpsLocation;
  tripId: string;
};
