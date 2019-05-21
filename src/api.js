import io from 'socket.io-client';

const socket = io('http://localhost:8084');

socket.on('connect', () => {
  console.log('connected');
});

socket.on('disconnect', () => {
  console.log('disconnected');
});

export default socket;
