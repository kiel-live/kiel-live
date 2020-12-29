import { Location } from './location';
import { ProviderEntity } from './provider-entity.class';

export enum StopType {
  BUS_STOP,
  PARKING_SPOT,
  FERRY_STOP,
  TRAIN_STOP,
  SUBWAY_STOP,
}

export class Stop extends ProviderEntity {
  name!: string;
  type!: StopType;
  location!: Location;
  alerts!: string[];
}
