import io from 'socket.io-client';
import config from '@/libs/config';

const url = config('backend_url', `${window.location.protocol}//${window.location.host}`);

const socket = io(url, {
  path: '/api/socket',
  transports: ['websocket'], // websocket only
});

export default socket;
