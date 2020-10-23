module.exports = {
  root: true,
  env: {
    node: true,
  },
  extends: [
    'airbnb-base',
  ],

  rules: {
    'no-console': 'error',
    'no-debugger': 'error',
    'no-param-reassign': ['error', { props: false }],
    'no-underscore-dangle': 'off',
    'object-curly-newline': ['error', {
      ObjectPattern: { minProperties: 10 },
    }],
  },
};
