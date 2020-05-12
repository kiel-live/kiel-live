// get last commit id as revision for /env-config.js
const envConfigRevision = process.env.VERSION || null;
const manifestJSON = require('./public/manifest.json');

module.exports = {
  devServer: {
    proxy: {
      '^/api/*': {
        target: 'http://localhost:3000',
        secure: false,
      },
    },
    progress: false,
  },
  pwa: {
    workboxPluginMode: 'InjectManifest',
    workboxOptions: {
      swSrc: './src/sw.js',
      swDest: 'service-worker.js',
      additionalManifestEntries: [
        { url: '/env-config.js', revision: envConfigRevision },
      ],
    },
    iconPaths: {
      favicon32: 'img/icons/favicon-32x32.png',
      favicon16: 'img/icons/favicon-16x16.png',
      appleTouchIcon: 'img/icons/apple-touch-icon.png',
      maskIcon: 'img/icons/safari-pinned-tab.svg',
      msTileImage: 'img/icons/mstile-150x150.png',
    },
    themeColor: manifestJSON.theme_color,
    name: manifestJSON.short_name,
    msTileColor: manifestJSON.background_color,
    appleMobileWebAppCapable: 'yes',
    appleMobileWebAppStatusBarStyle: 'black',
  },
};
