import type { VehicleType } from './vehicle';

export type ArrivalState = 'predicted' | 'stopping' | 'planned' | 'departed';

export interface StopArrival {
  name: string;
  type: VehicleType;
  vehicleId: string;
  tripId: string;
  routeId: string;
  routeName: string;
  direction: string;
  state: ArrivalState;
  planned: string;
  eta: number; // in seconds
  platform: string;
}

export interface TripArrival {
  id: string;
  name: string;
  state: ArrivalState;
  planned: string;
}
