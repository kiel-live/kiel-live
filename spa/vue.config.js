module.exports = {
  configureWebpack: {
    devServer: {
      proxy: {
        '^/api/*': {
          target: 'http://localhost:8081/api/',
          secure: false,
        },
      },
    },
  },
};
