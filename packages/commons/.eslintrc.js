module.exports = {
  extends: ['../../config/.eslintrc.base.js'],

  env: {
    es6: true,
    browser: true,
    node: true,
  },

  parserOptions: {
    tsconfigRootDir: __dirname,
    project: ['./tsconfig.eslint.json'],
  },
};
