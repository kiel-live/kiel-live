import { AbstractEntity } from './abstract-entity.class';
import { Location } from './location';

export enum StopType {
  BUS_STOP,
  PARKING_SPOT,
  FERRY_STOP,
  TRAIN_STOP,
  SUBWAY_STOP,
}

export class Stop extends AbstractEntity {
  name!: string;
  type!: StopType;
  location!: Location;
  alerts!: string[];
}
