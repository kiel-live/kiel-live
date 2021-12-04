import { TripArrival } from '.';

export type Trip = {
  id: string;
  provider: string;
  direction: string;
  arrivals: TripArrival[];
};
