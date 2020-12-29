import { Route } from '@kiel-live/commons';
import socketIo, { Socket } from 'socket.io';

const PORT = process.env.BACKEND_PORT ? parseInt(process.env.BACKEND_PORT) : 3000;

interface PremiumSocket extends Socket {
  on(event: 'route', listener: (route: Route) => void): this;
  on(event: 'disconnect', listener: (socket: Socket) => void): this;
}

function start() {
  const io = socketIo({ path: '/api/socket' });

  io.sockets.on('connection', (_socket) => {
    const socket = _socket as PremiumSocket;

    socket.on('route', (r) => {
      console.log(r.name);
    });

    socket.on('disconnect', () => {
      Api.leaveChannels(socket.id);
    });
  });

  io.listen(PORT, () => {
    // eslint-disable-next-line no-console
    console.log(`Server listening on port ${PORT}!`);
  });
}

start();
