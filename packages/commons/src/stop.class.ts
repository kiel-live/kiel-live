import { AbstractEntity } from './abstract-entity.class';

export enum StopType {
  BUS_STOP,
  PARKING_SPOT,
  FERRY_STOP,
  TRAIN_STOP,
  SUBWAY_STOP,
}

export enum StopArrivalState {
  PREDICTED,
  STOPPING,
  DEPARTED,
}

export type StopArrival = {
  name: string;
  vehicle: string; // TODO ref
  trip: string; // TODO ref
  route: string; // TODO ref
  direction: string;
  state: StopArrivalState;
  planned: Date;
  eta?: number;
};

export default class Stop extends AbstractEntity {
  id: string;
  name: string;
  provider: string;
  type: StopType;
  location: {
    longitude: string;
    latitude: string;
  };
  alerts: string[];
  arrivals: StopArrival[]; // changing frequently
}
