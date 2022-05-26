export type ArrivalState = 'predicted' | 'stopping' | 'planned' | 'departed';

export type StopArrival = {
  name: string;
  vehicleId: string;
  tripId: string;
  routeId: string;
  routeName: string;
  direction: string;
  state: ArrivalState;
  planned: string;
  eta: number; // in seconds
};

export type TripArrival = {
  id: string;
  name: string;
  state: ArrivalState;
  planned: string;
};
