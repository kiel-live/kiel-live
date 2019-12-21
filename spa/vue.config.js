module.exports = {
  configureWebpack: {
    devServer: {
      proxy: {
        '^/api/*': {
          target: 'http://localhost:8080/api/',
          secure: false,
        },
      },
    },
  },
  pwa: {
    workboxPluginMode: 'InjectManifest',
    workboxOptions: {
      swSrc: './src/sw.js',
      swDest: 'service-worker.js',
    },
    iconPaths: {
      favicon32: 'img/icons/favicon-32x32.png',
      favicon16: 'img/icons/favicon-16x16.png',
      appleTouchIcon: 'img/icons/apple-touch-icon.png',
      maskIcon: 'img/icons/safari-pinned-tab.svg',
      msTileImage: 'img/icons/mstile-150x150.png'
    },
    themeColor: '#2c3e50',
    msTileColor: '#ffffff',
    name: "Kiel Live"
  },
};
