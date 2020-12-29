import { Stop } from './stop.class';
import { Vehicle } from './vehicle.class';

export type TopicListener<T> = (data: T) => void;

export type TopicName = 'stops' | 'vehicles';

export abstract class AbstractSubscription<T> {
  protected topic: string;
  protected listener: TopicListener<T>;

  constructor(topic: string, listener: TopicListener<T>) {
    this.topic = topic;
    this.listener = listener;
  }

  abstract stop(): AbstractSubscription<T>;
}

export abstract class AbstractTopic<T> {
  protected topic: string;

  constructor(topic: string) {
    this.topic = topic;
  }

  abstract sub(listener: TopicListener<T>): AbstractSubscription<T>;

  abstract pub(data: T): void;
}

export interface TopicFnc {
  (topic: 'stops'): AbstractTopic<Stop>;
  (topic: 'vehicles'): AbstractTopic<Vehicle>;
  <T>(topic: TopicName): AbstractTopic<T>;
}
