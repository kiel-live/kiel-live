import { TopicName } from '@kiel-live/commons';
import http from 'http';
import { Server, Socket } from 'socket.io';

import Database from './database';

export default (port: number): Promise<void> => {
  const database = Database();

  const httpServer = http.createServer();

  const io = new Server(httpServer, {
    serveClient: false,
  });

  database.on((topic, item) => {
    io.to(topic).emit(topic, item);
  });

  io.on('connection', (socket: Socket) => {
    socket.on('subscribe', (topic: TopicName) => {
      void socket.join(topic);

      // send stored data from topic
      (database.getItems(topic) || []).forEach((item) => {
        socket.emit(topic, item);
      });
    });

    socket.on('unsubscribe', (topic) => {
      void socket.leave(topic);
    });
  });

  return new Promise((resolve) => {
    httpServer.listen(port, resolve);
  });
};
