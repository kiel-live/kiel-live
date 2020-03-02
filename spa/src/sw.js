/* eslint-disable no-restricted-globals */
/* global workbox */

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

workbox.core.clientsClaim(); // Vue CLI 4 and Workbox v4, else
// workbox.clientsClaim(); // Vue CLI 3 and Workbox v3.

// The precaching code provided by Workbox.
self.__precacheManifest = [].concat(self.__precacheManifest || []);
// workbox.precaching.suppressWarnings(); // Only used with Vue CLI 3 and Workbox v3.
workbox.precaching.precacheAndRoute(self.__precacheManifest, {});

// cache /env-config.js:
workbox.routing.registerRoute(
  new RegExp('/env-config.js'),
  workbox.strategies.cacheFirst({
    cacheName: 'config',
    plugins: [
      new workbox.expiration.Plugin({
        maxEntries: 30,
        maxAgeSeconds: 30 * 24 * 60 * 60, // 30 days
      }),
    ],
  }),
);
