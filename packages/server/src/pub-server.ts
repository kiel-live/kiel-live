import { AbstractEntity, TopicName } from '@kiel-live/commons';
import http from 'http';
import { Server, Socket } from 'socket.io';

import Database from './database';

export default (port: number): Promise<void> => {
  const database = Database();

  const httpServer = http.createServer();

  const io = new Server(httpServer, {
    serveClient: false,
  });

  io.on('connection', (socket: Socket) => {
    socket.on('publish-items', (topic: TopicName, items: AbstractEntity[]) => {
      database.setItems(topic, items);
    });

    socket.on('publish-item', (topic: TopicName, id: string, item: AbstractEntity) => {
      database.setItem(topic, id, item);
    });
  });

  return new Promise((resolve) => {
    httpServer.listen(port, resolve);
  });
};
