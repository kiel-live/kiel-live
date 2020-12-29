import { AbstractEntity, Ref } from './abstract-entity.class';
import { Stop } from './stop.class';

export enum RouteType {
  BUS,
  FERRY,
  TRAIN,
  SUBWAY,
}

export class Route extends AbstractEntity {
  name!: string;
  type!: RouteType;
  active!: boolean;
  stops!: Ref<Stop>[];
}
