import { AbstractEntity, TopicName } from '@kiel-live/commons';

type DatabaseListener = (topic: string, item: AbstractEntity) => void;

type DatabaseTopic<T extends AbstractEntity> = {
  [key: string]: T;
};

type DatabaseItems = {
  [key in TopicName]: DatabaseTopic<AbstractEntity>;
};

class Database {
  items: DatabaseItems = {
    stops: {},
    vehicles: {},
  };

  listener?: DatabaseListener;

  getItems(topic: TopicName): AbstractEntity[] | null {
    if (!this.items[topic]) {
      return null;
    }

    return Object.values(this.items[topic]);
  }

  getItem(topic: TopicName, id: string): AbstractEntity | null {
    if (!this.items[topic][id]) {
      return null;
    }

    return this.items[topic][id];
  }

  setItem(topic: TopicName, id: string, item: AbstractEntity) {
    this.items[topic][id] = item;

    if (this.listener) {
      this.listener(topic, item);
    }
  }

  setItems(topic: TopicName, items: AbstractEntity[]) {
    for (const item of items) {
      this.setItem(topic, item.id, item);
    }
  }

  on(listener: DatabaseListener) {
    this.listener = listener;
  }
}

let database: Database;

export default (): Database => {
  if (!database) {
    database = new Database();
  }

  return database;
};
