import io from 'socket.io-client';
import store from '@/store';

const socket = io('', { path: '/api' });

socket.on('connect', () => {
  store.commit('connect');
});

socket.on('disconnect', () => {
  store.commit('disconnect');
});

export default socket;
