import { Ref } from './abstract-entity.class';
import { ProviderEntity } from './provider-entity.class';
import { Stop } from './stop.class';

export enum RouteType {
  BUS,
  FERRY,
  TRAIN,
  SUBWAY,
}

export class Route extends ProviderEntity {
  name!: string;
  type!: RouteType;
  active!: boolean;
  stops!: Ref<Stop>[];
}
