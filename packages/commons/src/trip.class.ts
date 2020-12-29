import { AbstractEntity } from './abstract-entity.class';

export enum TripState {
  PREDICTED,
  STOPPING,
  DEPARTED,
};

export type TripStop = {
  id: string;
  name: string;
  state: TripState;
  planned: Date;
  eta?: number;
};

export default class Trip extends AbstractEntity {
  id: string;
  vehicle: string;
  direction: string;
  stops: TripStop[];
}
