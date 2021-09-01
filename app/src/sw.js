/* eslint-disable no-restricted-globals */
import { skipWaiting, clientsClaim } from 'workbox-core';
import { precacheAndRoute } from 'workbox-precaching';

// This is the code piece that GenerateSW mode can't provide for us.
// This code listens for the user's confirmation to update the app.
self.addEventListener('message', (e) => {
  if (!e.data) {
    return;
  }

  switch (e.data) {
    case 'skipWaiting':
      skipWaiting();
      break;
    default:
      // NOOP
      break;
  }
});

clientsClaim();

precacheAndRoute(self.__WB_MANIFEST);
