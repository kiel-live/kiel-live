import { DEBUG } from '~/config';

export class ReconnectingWebSocket {
  ws?: WebSocket;
  url: string;
  emitter: Emitter;

  constructor(url: string, autoConnect = true) {
    this.url = url.replace('http', 'ws');
    this.emitter = new Emitter();
    if (autoConnect) {
      this.connect();
    }
  }

  async connect() {
    this.ws = new WebSocket(this.url);

    this.ws.addEventListener('close', (event) => {
      this.log('connection closed. Reconnect will be attempted in 1 second.', event.reason);
      this.emitter.emit('close', event);
      setTimeout(() => {
        this.connect();
      }, 1000);
    });

    this.ws.addEventListener('error', (event) => {
      this.log('error encountered: ', (event as unknown as Error)?.message ?? event, 'Closing socket');
      this.emitter.emit('error', event);
      this.ws?.close();
    });

    this.ws.addEventListener('message', (event) => {
      this.emitter.emit('message', event);
    });

    return new Promise<void>((resolve) => {
      this.ws?.addEventListener('open', (event) => {
        this.log('connected');
        this.emitter.emit('open', event);
        resolve();
      });
    });
  }

  send(data: string) {
    this.ws?.send(data);
  }

  on(event: 'message', cb: (e: MessageEvent) => void): void;
  on(event: 'open' | 'close' | 'error', cb: (e: Event) => void): void;
  on(event: string, cb: (e: any) => void) {
    this.emitter.on(event, cb);
  }

  private log(...args: unknown[]) {
    if (DEBUG) {
      // eslint-disable-next-line no-console
      console.log('Reconnecting Websocket', ...args);
    }
  }
}

type Callback = (...args: unknown[]) => void;

class Emitter {
  eventMap: Map<string, Callback[]>;

  constructor() {
    this.eventMap = new Map();
  }

  on(event: string, callback: Callback) {
    if (!this.eventMap.has(event)) {
      this.eventMap.set(event, []);
    }
    this.eventMap.get(event)!.push(callback);
  }

  off(event: string, callback: Callback) {
    if (this.eventMap.has(event)) {
      const callbacks = this.eventMap.get(event)!.filter((cb) => cb !== callback);
      this.eventMap.set(event, callbacks);
    }
  }

  emit(event: string, ...data: unknown[]) {
    if (this.eventMap.has(event)) {
      this.eventMap.get(event)!.forEach((callback) => {
        setTimeout(() => callback(...data), 0);
      });
    }
  }
}
