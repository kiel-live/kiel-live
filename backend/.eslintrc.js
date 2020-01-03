module.exports = {
  root: true,
  env: {
    node: true,
  },
  extends: [
    'airbnb-base',
  ],

  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-param-reassign': ['error', { props: false }],
    'no-underscore-dangle': 'off',
    'object-curly-newline': ['error', {
      'ObjectPattern': { 'minProperties': 10 },
    }],
  },
};
