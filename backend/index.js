const path = require('path');
const express = require('express');
const socketIo = require('socket.io');
const http = require('http');

const Api = require('./src/api');
const cachedResponse = require('./src/cachedResponse');

const app = express();
const server = http.createServer(app);

const PORT = process.env.PORT || 8080;
const DIST_DIR = [__dirname, 'dist'];
let connectedClients = 0;

function start() {
  const io = socketIo(server, { path: '/api/socket' });

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

  app.get('/api/osm-tiles/:z/:x/:y.png', cachedResponse, async (req, resp) => {
    const { z, x, y } = req.params;
    const s = String.fromCharCode(97 + Math.floor(Math.random() * 3)); // select a, b or c

    const options = {
      hostname: `${s}.tile.osm.org`,
      port: 80,
      path: `/${z}/${x}/${y}.png`,
      method: 'GET',
      headers: {
        'User-Agent': 'osm-proxy-opnv-live',
        Accept: '*/*',
      },
    };

    resp.set({
      'Cache-Control': 'public, max-age=86400',
      Expires: new Date(Date.now() + 86400000).toUTCString(),
    });

    const proxy = http.request(options, (res) => {
      const data = [];

      res.on('data', (chunk) => {
        data.push(chunk);
      }).on('end', () => {
        const buffer = Buffer.concat(data);
        resp.send(buffer);
      });
    });

    proxy.end();
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
      const stops = await Api.lookupStops(query);
      if (stops) {
        socket.emit('stop:search', stops);
      }
    });

    socket.on('stop:nearby', async (opts) => {
      const stops = await Api.nearby(opts);
      if (stops) {
        socket.emit('stop:nearby', stops);
      }
    });

    socket.on('geo:stops', async () => {
      const stops = await Api.geoStops();
      if (stops) {
        socket.emit('geo:stops', stops);
      }
    });

    socket.on('geo:vehicles:join', () => {
      socket.join('geo:vehicles');
      Api.joinGeoVehicles({ clientId: socket.id }, (data) => io.to('geo:vehicles').emit('geo:vehicles', data));
    });

    socket.on('geo:vehicles:leave', () => {
      socket.leave('geo:vehicles');
      Api.leaveGeoVehicles({ clientId: socket.id });
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
