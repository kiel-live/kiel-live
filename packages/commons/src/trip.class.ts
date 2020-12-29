import { Ref } from './abstract-entity.class';
import { ProviderEntity } from './provider-entity.class';
import { Route } from './route.class';

export enum TripState {
  PREDICTED,
  STOPPING,
  DEPARTED,
}

export type TripStop = {
  id: string;
  name: string;
  state: TripState;
  planned: Date;
  eta?: number;
};

export class Trip extends ProviderEntity {
  vehicle!: string;
  direction!: string;
  route!: Ref<Route>;
  stops!: TripStop[];
}
