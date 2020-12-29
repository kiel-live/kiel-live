export abstract class AbstractEntity {
  protected readonly _id!: string;
  provider!: string;

  get id(): string {
    return `${this.provider}-${this._id}`;
  }
}

export type Ref<T extends AbstractEntity> = T['id'];
