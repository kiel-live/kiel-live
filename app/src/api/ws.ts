export class ReconnectingWebSocket {
  ws?: WebSocket;
  url: string;
  eventListeners: Record<string, ((e: Event) => void)[]> = {};

  constructor(url: string, autoConnect = true) {
    this.url = url.replace('http', 'ws');
    if (autoConnect) {
      this.connect();
    }
  }

  async connect() {
    this.ws = new WebSocket(this.url);

    this.ws.onclose = (e) => {
      console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
      setTimeout(() => {
        this.connect();
      }, 1000);
    };

    this.ws.onerror = (err: Event) => {
      console.error('Socket encountered error: ', err.message, 'Closing socket');
      this.ws?.close();
    };

    for (const [event, listeners] of Object.entries(this.eventListeners)) {
      for (const listener of listeners) {
        this.ws.addEventListener(event, listener);
      }
    }

    return new Promise<void>((resolve) => {
      this.ws!.onopen = () => {
        console.log('Socket is open');
        resolve();
      };
    });
  }

  send(data: string) {
    this.ws?.send(data);
  }

  on(event: string, cb: (e: Event) => void) {
    if (!this.eventListeners[event]) {
      this.eventListeners[event] = [];
    }
    this.eventListeners[event].push(cb);
    this.ws?.addEventListener(event, cb);
  }
}
