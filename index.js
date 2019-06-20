const path = require('path');
const express = require('express');
const socketIo = require('socket.io');
const hafas = require('nahsh-hafas');
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

    socket.on('info', () => {
      socket.emit('info', {
        connectedClients,
      });
    });

    socket.on('nearby', async ({latitude, longitude}) => {
      const result = hafas.nearby({
        type: 'location',
        latitude,
        longitude,
      }, {distance: 400});

      socket.emit('nearby', result);
    });

    socket.on('trip', async ({id, lineName}) => {
      const result = hafas.client.journeys(id, lineName, {results: 1});

      socket.emit('trip', result);
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
