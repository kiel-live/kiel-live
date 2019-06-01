const path = require('path');
const express = require('express');
const socketIo = require('socket.io');
const app = express();
const server = require('http').createServer(app);
const stops = require('./stops');

const PORT = process.env.PORT || 8080;

function start() {
  const io = socketIo(server, { path: '/api' });

  app.use(express.static(path.join(__dirname, 'spa', 'dist')));
  app.use('*', (req, resp) => {
    resp.sendFile(path.join(__dirname, 'spa', 'dist', 'index.html'));
  });

  stops.setIO(io);

  io.sockets.on('connection', (socket) => {
    console.log('client connected');

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
  });

  server.listen(PORT, () => {
    console.log(`Server listening on port ${PORT}!`);
  });
}

start();
