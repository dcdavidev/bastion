import { Linter } from 'eslint';
import { defineConfig } from 'eslint/config';
import packageJson from 'eslint-plugin-package-json';
import prettierRecommended from 'eslint-plugin-prettier/recommended';
import * as jsoncParser from 'jsonc-eslint-parser';
import json from '@eslint/json';

export const configJson: Linter.Config[] = defineConfig([
  {
    files: ['**/*.json'],
    plugins: { json: json },
    language: 'json/json',
    extends: ['json/recommended'],
    languageOptions: { parser: jsoncParser },
    rules: {
      'prettier/prettier': 'off',
    },
  },
  {
    files: ['**/*.json5'],
    plugins: { json: json },
    language: 'json/json5',
    extends: ['json/recommended'],
    languageOptions: { parser: jsoncParser },
  },
  {
    files: ['**/*.jsonc', '**/tsconfig*.json', '**/.vscode/**/*.json'],
    plugins: { json: json },
    language: 'json/jsonc',
    extends: ['json/recommended'],
    languageOptions: { parser: jsoncParser },
  },
  {
    files: ['**/package.json'],
    plugins: { json: json, 'package-json': packageJson },
    language: 'json/json',
    extends: [packageJson.configs.recommended],
    languageOptions: { parser: jsoncParser },
    rules: {
      'package-json/order-properties': 'error',
      'package-json/sort-collections': 'error',
      'package-json/require-description': 'error',
      'package-json/require-bugs': 'error',
      'package-json/require-keywords': 'error',
      'package-json/require-name': 'error',
      'package-json/require-version': 'error',
      'package-json/valid-description': 'error',
      'package-json/valid-license': 'error',
      'package-json/valid-name': 'error',
      'package-json/valid-package-definition': 'error',
      'package-json/valid-version': 'error',
      'prettier/prettier': 'off',
    },
  },

  prettierRecommended,
]);
