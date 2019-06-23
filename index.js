const path = require('path');
const express = require('express');
const socketIo = require('socket.io');
const app = express();
const server = require('http').createServer(app);
const Api = require('./src/api');

const PORT = process.env.PORT || 8080;
let connectedClients = 0;

function start() {
  const io = socketIo(server, { path: '/api' });

  app.use(express.static(path.join(__dirname, 'spa', 'dist')));
  app.use('*', (req, resp) => {
    resp.sendFile(path.join(__dirname, 'spa', 'dist', 'index.html'));
  });

  io.sockets.on('connection', (socket) => {
    connectedClients += 1;

    socket.on('stop:join', (stopId) => {
      if (!stopId) { return; }
      console.log('client joined stop', stopId);
      socket.join(`stop:${stopId}`);
      Api.joinStop(stopId, (data) => io.to(`stop:${stopId}`).emit('stop', data));
    });

    socket.on('stop:leave', (stopId) => {
      if (!stopId) { return; }
      console.log('client left stop', stopId);
      Api.leaveStop(stopId);
      socket.leave(`stop:${stopId}`);
    });

    socket.on('stop:search', async (query) => {
      socket.emit('stop:search', await Api.lookupStops(query));
    });

    socket.on('stop:nearby', async (opts) => {
      socket.emit('stop:nearby', await Api.nearby(opts));
    });

    socket.on('trip:join', ({ tripId, vehicleId }) => {
      if (!tripId || !vehicleId) { return; }
      console.log('client joined trip', tripId);
      socket.join(`trip:${tripId}:${vehicleId}`);
      Api.joinTrip({ tripId, vehicleId }, (data) => io.to(`trip:${tripId}:${vehicleId}`).emit('trip', data));
    });

    socket.on('trip:leave', ({ tripId, vehicleId }) => {
      if (!tripId || !vehicleId) { return; }
      console.log('client left trip', tripId);
      Api.leaveTrip({ tripId, vehicleId });
      socket.leave(`trip:${tripId}:${vehicleId}`);
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
