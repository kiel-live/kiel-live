import { AbstractEntity, Ref } from './abstract-entity.class';
import Stop from './stop.class';

export enum RouteType {
  BUS,
  FERRY,
  TRAIN,
  SUBWAY,
}

export default class Route extends AbstractEntity {
  id: string;
  provider: string;
  name: string;
  type: RouteType;
  active: boolean;
  stops: Ref<Stop>[];
}
