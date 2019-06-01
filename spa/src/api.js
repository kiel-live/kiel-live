import io from 'socket.io-client';

const socket = io('', { path: '/api' });

socket.on('connect', () => {
  console.log('connected');
});

socket.on('disconnect', () => {
  console.log('disconnected');
});

export default socket;
