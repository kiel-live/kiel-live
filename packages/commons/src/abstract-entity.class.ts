export abstract class AbstractEntity {
  readonly _id!: string;
}

export type Ref<T extends AbstractEntity> = T['_id'];
