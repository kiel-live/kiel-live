import { AbstractSubscription, AbstractTopic, TopicFnc, TopicListener } from '@kiel-live/commons';
import io from 'socket.io-client';

export const socket = io('ws://example.com/my-namespace');

class Subscription<T> extends AbstractSubscription<T> {
  constructor(topic: string, listener: TopicListener<T>) {
    super(topic, listener);

    socket.emit('subscribe', this.topic);
    socket.on(this.topic, listener);
  }

  stop() {
    socket.emit('unsubscribe', this.topic);
    socket.off(this.topic, this.listener);
    return this;
  }
}

class Topic<T> extends AbstractTopic<T> {
  sub(listener: TopicListener<T>): Subscription<T> {
    return new Subscription<T>(this.topic, listener);
  }

  pub(data: T): void {
    socket.emit(this.topic, data);
  }
}

export const topic: TopicFnc = <T>(topic: string) => new Topic<T>(topic);
