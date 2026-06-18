interface Envelope {
  rid?: string;
  method: string;
  topic?: string;
  data?: unknown;
}

type SendFn = (env: Envelope) => void;

export class Rpc {
  private pending = new Map<string, (data: any) => void>();

  // Call when a 'result' or 'error' envelope arrives.
  resolve<T>(rid: string, data: T) {
    const cb = this.pending.get(rid);
    if (cb) {
      this.pending.delete(rid);
      cb(data);
    }
  }

  // Send a request envelope and return a promise for the result.
  // Rejects after timeoutMs if no response arrives.
  request<T>(send: SendFn, topic: string, data: T, timeoutMs = 5000): Promise<T> {
    return new Promise((resolve, reject) => {
      const rid = Math.random().toString(36).slice(2);
      const timer = setTimeout(() => {
        this.pending.delete(rid);
        reject(new Error(`rpc timeout: ${topic}`));
      }, timeoutMs);
      this.pending.set(rid, (result: T) => {
        clearTimeout(timer);
        resolve(result);
      });
      send({ method: 'request', topic, rid, data });
    });
  }
}
