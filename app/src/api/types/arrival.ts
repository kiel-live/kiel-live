export type ArrivalState = 'predicted' | 'stopping' | 'planned';

export type StopArrival = {
  name: string;
  vehicleID: string;
  tripId: string;
  routeId: string;
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
