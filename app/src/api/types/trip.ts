import { TripArrival } from './arrival';

export type Trip = {
  id: string;
  provider: string;
  direction: string;
  arrivals?: TripArrival[];
};
