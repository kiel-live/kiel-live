import type { VehicleType } from './vehicle';

export type DepartureState = 'predicted' | 'stopping' | 'planned' | 'departed';

export interface StopDeparture {
  name: string;
  type: VehicleType;
  vehicleId: string;
  tripId: string;
  routeId: string;
  routeName: string;
  direction: string;
  state: DepartureState;
  planned: string; // iso datetime string
  actual: string; // iso datetime string
  platform: string;
}

export interface TripDeparture {
  id: string;
  name: string;
  state: DepartureState;
  planned: string;
  actual?: string;
}
