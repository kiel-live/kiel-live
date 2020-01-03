const path = require('path');
const express = require('express');
const socketIo = require('socket.io');

const app = express();
const server = require('http').createServer(app);
const Api = require('./src/api');

const PORT = process.env.PORT || 8080;
const DIST_DIR = [__dirname, 'dist'];
let connectedClients = 0;

function start() {
  const io = socketIo(server, { path: '/api' });

  app.use(express.static(path.join(...DIST_DIR)));
  app.use('/status', (req, resp) => {
    const status = {
      version: process.env.VERSION || null,
      clients: connectedClients,
      ...Api.status(),
    };

    // pretty json
    resp.header('Content-Type', 'application/json');
    resp.send(JSON.stringify(status, null, 4));
  });
  app.use('*', (req, resp) => {
    resp.sendFile(path.join(...DIST_DIR, 'index.html'));
  });

  io.sockets.on('connection', (socket) => {
    connectedClients += 1;

    socket.on('stop:join', (stopId) => {
      if (!stopId) { return; }
      socket.join(`stop:${stopId}`);
      Api.joinStop({ stopId, clientId: socket.id }, (data) => io.to(`stop:${stopId}`).emit('stop', data));
    });

    socket.on('stop:leave', (stopId) => {
      if (!stopId) { return; }
      socket.leave(`stop:${stopId}`);
      Api.leaveStop({ stopId, clientId: socket.id });
    });

    socket.on('trip:join', ({ tripId, vehicleId }) => {
      if (!tripId || !vehicleId) { return; }
      socket.join(`trip:${tripId}:${vehicleId}`);
      Api.joinTrip({ tripId, vehicleId, clientId: socket.id }, (data) => io.to(`trip:${tripId}:${vehicleId}`).emit('trip', data));
    });

    socket.on('trip:leave', ({ tripId, vehicleId }) => {
      if (!tripId || !vehicleId) { return; }
      socket.leave(`trip:${tripId}:${vehicleId}`);
      Api.leaveTrip({ tripId, vehicleId, clientId: socket.id });
    });

    socket.on('stop:search', async (query) => {
      socket.emit('stop:search', await Api.lookupStops(query));
    });

    socket.on('stop:nearby', async (opts) => {
      socket.emit('stop:nearby', await Api.nearby(opts));
    });

    socket.on('disconnect', () => {
      Api.leaveChannels(socket.id);
      connectedClients -= 1;
    });
  });

  server.listen(PORT, () => {
    console.log(`Server listening on port ${PORT}!`);
  });
}

start();
