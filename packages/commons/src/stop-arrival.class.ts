import { Ref } from './abstract-entity.class';
import { ProviderEntity } from './provider-entity.class';
import { Stop } from './stop.class';
import { Trip } from './trip.class';
import { Vehicle } from './vehicle.class';

export enum StopArrivalState {
  PREDICTED,
  STOPPING,
  DEPARTED,
}

export class StopArrival extends ProviderEntity {
  stop!: Ref<Stop>;
  name!: string;
  vehicle!: Ref<Vehicle>;
  trip!: Ref<Trip>;
  direction!: string;
  state!: StopArrivalState;
  planned!: Date;
  eta?: number;
}
