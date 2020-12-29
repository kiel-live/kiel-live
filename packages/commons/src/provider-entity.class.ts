import { AbstractEntity } from './abstract-entity.class';

export class ProviderEntity extends AbstractEntity {
  _id!: string;
  provider!: string;

  get id(): string {
    return `${this.provider}-${this._id}`;
  }
}
