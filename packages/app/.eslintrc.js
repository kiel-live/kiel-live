module.exports = {
  root: true,

  env: {
    node: true,
  },

  extends: [
    'plugin:vue/recommended',
    '@vue/airbnb',
  ],

  rules: {
    'no-console': 'error',
    'no-debugger': 'error',
    'no-param-reassign': ['error', { props: false }],
    'no-underscore-dangle': 'off',
    'vue/no-unused-properties': ['error', {
      groups: ['props', 'data', 'computed', 'methods', 'setup'],
    }],
    'vue/max-attributes-per-line': ['warn', {
      singleline: 10,
      multiline: {
        max: 1,
        allowFirstLine: false,
      },
    }],
  },

  parserOptions: {
    parser: 'babel-eslint',
  },

  overrides: [
    {
      files: ['*.vue'],
      rules: {
        'max-len': 'off',
      },
    },
  ],
};
