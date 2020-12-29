module.exports = {
  root: true,

  parser: '@typescript-eslint/parser',

  plugins: ['@typescript-eslint', 'import', 'simple-import-sort'],

  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:@typescript-eslint/recommended-requiring-type-checking',
    'prettier/@typescript-eslint',
    'plugin:prettier/recommended',
    'plugin:import/errors',
    'plugin:import/warnings',
    'plugin:import/typescript',
    'plugin:promise/recommended',
  ],

  rules: {
    'no-unused-expressions': 'off',
    '@typescript-eslint/no-unused-expressions': ['error'],
    'no-unused-vars': 'off',
    '@typescript-eslint/no-explicit-any': 'error',
    '@typescript-eslint/no-unused-vars': 'error',
    'no-restricted-imports': [
      'error',
      {
        patterns: ['src', 'dist'],
      },
    ],
    'no-console': 'error',
    curly: ['warn', 'all'],
    'no-useless-concat': 'error',
    'promise/prefer-await-to-then': 'error',
    'promise/prefer-await-to-callbacks': 'error',
    // import rules
    'import/no-unresolved': 'off', // disable as this is handled by tsc itself
    'import/first': 'error',
    'import/newline-after-import': 'error',
    'import/no-cycle': 'error',
    'import/no-relative-parent-imports': 'error',
    'import/no-duplicates': 'error',
    'simple-import-sort/imports': 'error',
    'simple-import-sort/exports': 'error',
    // disable different sort rules because simple-import-sort takes care of this
    'sort-imports': 'off',
    'import/order': 'off',
    '@typescript-eslint/naming-convention': [
      'error',
      {
        selector: 'default',
        format: ['camelCase'],
        filter: {
          regex: '[-.]',
          match: false,
        },
      },
      {
        selector: 'variable',
        format: ['camelCase', 'UPPER_CASE'],
      },
      {
        selector: 'typeLike',
        format: ['PascalCase'],
      },
      {
        selector: 'memberLike',
        format: ['camelCase'],
        leadingUnderscore: 'allow',
        filter: {
          regex: '[-.]',
          match: false,
        },
      },
      {
        selector: 'property',
        format: ['camelCase', 'PascalCase'],
        leadingUnderscore: 'allow',
        // ignore properties that require quotes
        // this filter has to be applied to less specific selectors that include `property` (memberLike, default)
        filter: {
          regex: '[-.]',
          match: false,
        },
      },
      {
        selector: 'parameter',
        format: ['camelCase'],
        leadingUnderscore: 'allow',
      },
      {
        selector: 'enumMember',
        format: ['UPPER_CASE', 'PascalCase'], // TODO: should be only UPPER_CASE
      },
    ],
    // TODO: the following rules should be set to error at some point
    '@typescript-eslint/no-unnecessary-type-assertion': 'error',
    '@typescript-eslint/no-unsafe-assignment': 'warn',
    '@typescript-eslint/no-unsafe-member-access': 'warn',
    '@typescript-eslint/no-unsafe-call': 'warn',
    '@typescript-eslint/no-unsafe-return': 'warn',
    '@typescript-eslint/no-floating-promises': 'warn',
    '@typescript-eslint/restrict-template-expressions': 'warn',
  },

  // parserOptions: {
  //   sourceType: 'module',
  //   tsconfigRootDir: __dirname,
  //   project: ['./tsconfig.eslint.json'],
  // },

  // overrides: [
  //   {
  //     files: ['**/*.js'],
  //     rules: {
  //       '@typescript-eslint/no-unsafe-assignment': 'off',
  //       '@typescript-eslint/no-unsafe-member-access': 'off',
  //       '@typescript-eslint/no-unsafe-call': 'off',
  //       '@typescript-eslint/restrict-plus-operands': 'off',
  //       '@typescript-eslint/no-var-requires': 'off',
  //     },
  //   },
  // ],
};
