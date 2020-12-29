import { LocationHeading } from './location';
import { ProviderEntity } from './provider-entity.class';

export enum VehicleType {
  BUS,
  BIKE,
  CAR,
  ESCOOTER,
  FERRY,
  TRAIN,
  SUBWAY,
}

export class Vehicle extends ProviderEntity {
  name!: string;
  type!: VehicleType;
  state?: string; // TODO define
  battery?: number;
  location!: LocationHeading;
}
