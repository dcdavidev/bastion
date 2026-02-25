import type { Linter } from 'eslint';
import { defineConfig } from 'eslint/config';
import vuoto from 'eslint-plugin-vuoto';

export const configVuoto: Linter.Config[] = defineConfig([
  {
    plugins: {
      vuoto,
    },
  },
  vuoto.configs.all,
]);
