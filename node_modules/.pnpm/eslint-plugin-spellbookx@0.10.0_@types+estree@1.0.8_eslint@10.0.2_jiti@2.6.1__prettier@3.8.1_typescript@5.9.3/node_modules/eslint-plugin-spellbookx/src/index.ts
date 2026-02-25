import type { ESLint, Linter } from 'eslint';

import { configCspell } from './configs/cspell.js';
import configIgnores from './configs/ignores.js';
import { configJavascript } from './configs/javascript.js';
import { configJson } from './configs/json.js';
import { configMarkdown } from './configs/markdown.js';
import { configRecommended } from './configs/recommended.js';
import { configVuoto } from './configs/vuoto.js';
import { packageJson } from './helpers.js';
import { SbxESLintConfig } from './types.js';

const configs: Record<SbxESLintConfig, Linter.Config[]> = {
  cspell: configCspell,
  ignores: configIgnores,
  javascript: configJavascript,
  json: configJson,
  markdown: configMarkdown,
  recommended: configRecommended,
  vuoto: configVuoto,
};

type Plugin = ESLint.Plugin & {
  configs: typeof configs;
};

const plugin: Plugin = {
  meta: {
    name: packageJson.name,
    version: packageJson.version,
    namespace: 'spellbookx',
  },
  configs,
  rules: {},
  processors: {},
};

export default plugin;
