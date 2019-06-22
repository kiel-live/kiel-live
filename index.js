const path = require('path');
const express = require('express');
const socketIo = require('socket.io');
const app = express();
const server = require('http').createServer(app);
const stops = require('./stops');

const PORT = process.env.PORT || 8080;
let connectedClients = 0;

function start() {
  const io = socketIo(server, { path: '/api' });

  app.use(express.static(path.join(__dirname, 'spa', 'dist')));
  app.use('*', (req, resp) => {
    resp.sendFile(path.join(__dirname, 'spa', 'dist', 'index.html'));
  });

  stops.setIO(io);

  io.sockets.on('connection', (socket) => {
    connectedClients += 1;

    socket.on('stop:join', (stopId) => {
      if (!stopId) { return; }
      console.log('client joined stop', stopId);
      socket.join(`stop:${stopId}`);
      stops.join(stopId, socket);
    });

    socket.on('stop:leave', (stopId) => {
      if (!stopId) { return; }
      console.log('client left stop', stopId);
      stops.leave(stopId);
      socket.leave(`stop:${stopId}`);
    });

    socket.on('stop:search', async (query) => {
      socket.emit('stop:search', await stops.lookupStops(query));
    });

    socket.on('stop:nearby', async (opts) => {
      const res = {};
      const tmp = await stops.nearby(opts);

      for (let i = 0; i < tmp.length; i++) {
        const item = tmp[i];
        const name = item.name.replace(/Kiel\s/, '');
        const lookup = await stops.lookupStops(name);
        console.log(lookup);
        if (lookup && lookup.length === 1) {
          const found = {
            ...item,
            ...lookup[0],
          };
          found.gps = true;
          res[found.id] = found;          
        }
      }

      socket.emit('stop:nearby', res);
    });

    socket.on('trip', async (opts) => {
      socket.emit('trip', await stops.trip(opts));
    });

    socket.on('info', () => {
      socket.emit('info', {
        connectedClients,
      });
    });

    socket.on('disconnect', () => {
      connectedClients -= 1;
    });
  });

  server.listen(PORT, () => {
    console.log(`Server listening on port ${PORT}!`);
  });
}

start();
