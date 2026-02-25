import { Linter } from 'eslint';
import { defineConfig } from 'eslint/config';
import prettierRecommended from 'eslint-plugin-prettier/recommended';

import { configCspell } from './cspell.js';
import configIgnores from './ignores.js';
import { configJavascript } from './javascript.js';
import { configJson } from './json.js';
import { configMarkdown } from './markdown.js';
import { configVuoto } from './vuoto.js';

export const configRecommended: Linter.Config[] = defineConfig([
  configIgnores,
  configJavascript,
  configJson,
  configMarkdown,
  configVuoto,
  configCspell,
  prettierRecommended,
]);
