module.exports = {
  extends: ['../../config/.eslintrc.base.js'],

  parserOptions: {
    tsconfigRootDir: __dirname,
    project: ['./tsconfig.json'],
  },
};
