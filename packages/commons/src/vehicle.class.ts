import { AbstractEntity } from './abstract-entity.class';

export enum VehicleType {
  BUS,
  BIKE,
  CAR,
  ESCOOTER,
  FERRY,
  TRAIN,
  SUBWAY,
};

export default class Vehicle extends AbstractEntity {
  id: string;
  name: string;
  provider: string;
  type: VehicleType;
  state?: string;
  battery?: number;
  location: {
    heading: number;
    longitude: number;
    latitude: number;
  };
}