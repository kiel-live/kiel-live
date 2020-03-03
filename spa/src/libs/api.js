import io from 'socket.io-client';

const socket = io({
  path: '/api/socket',
  transports: ['websocket'], // websocket only
});

export default socket;
