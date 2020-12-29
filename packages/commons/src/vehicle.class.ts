import { AbstractEntity } from './abstract-entity.class';
import { LocationHeading } from './location';

export enum VehicleType {
  BUS,
  BIKE,
  CAR,
  ESCOOTER,
  FERRY,
  TRAIN,
  SUBWAY,
}

export class Vehicle extends AbstractEntity {
  name!: string;
  type!: VehicleType;
  state?: string; // TODO define
  battery?: number;
  location!: LocationHeading;
}
