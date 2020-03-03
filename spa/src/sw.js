/* eslint-disable no-restricted-globals */
/* global workbox, importScripts */

// Import workbox-sw, which defines the global `workbox` object.
importScripts('https://storage.googleapis.com/workbox-cdn/releases/5.0.0/workbox-sw.js');

// This is the code piece that GenerateSW mode can't provide for us.
// This code listens for the user's confirmation to update the app.
self.addEventListener('message', (e) => {
  if (!e.data) {
    return;
  }

  switch (e.data) {
    case 'skipWaiting':
      self.skipWaiting();
      break;
    default:
      // NOOP
      break;
  }
});

workbox.core.clientsClaim();

workbox.precaching.precacheAndRoute(self.__WB_MANIFEST);
