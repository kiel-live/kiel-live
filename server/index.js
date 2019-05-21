const socketIo = require('socket.io');
const stops = require('./stops');

const PORT = 8084;

function start() {
  const io = socketIo.listen(PORT);

  stops.setIO(io);

  io.sockets.on('connection', (socket) => {
    console.log('client connected');

    socket.on('stop:join', (stopId) => {
      console.log('client joined stop', stopId);
      if (!stopId) { return; }
      socket.join(`stop:${stopId}`);
      stops.join(stopId);
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

  console.log(`Server listening on port: ${PORT}`);
}

start();
