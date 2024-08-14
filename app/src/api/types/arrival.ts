export type ArrivalState = 'predicted' | 'stopping' | 'planned' | 'departed';

export type ArrivalType =
  | 'bus-stop'
  | 'parking-spot'
  | 'ferry-stop'
  | 'train-stop'
  | 'subway-stop'
  | 'bike-stop'
  | 'tram-stop';

export type StopArrival = {
  name: string;
  type: ArrivalType;
  vehicleId: string;
  tripId: string;
  routeId: string;
  routeName: string;
  direction: string;
  state: ArrivalState;
  planned: string;
  eta: number; // in seconds
  platform: string;
};

export type TripArrival = {
  id: string;
  name: string;
  state: ArrivalState;
  planned: string;
};
