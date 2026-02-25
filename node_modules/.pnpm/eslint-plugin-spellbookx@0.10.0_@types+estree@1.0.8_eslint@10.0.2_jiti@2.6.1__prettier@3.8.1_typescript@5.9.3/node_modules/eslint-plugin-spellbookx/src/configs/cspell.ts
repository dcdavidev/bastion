import { Linter } from 'eslint';
import { defineConfig } from 'eslint/config';
import prettierRecommended from 'eslint-plugin-prettier/recommended';
import cspell from '@cspell/eslint-plugin/configs';

export const configCspell: Linter.Config[] = defineConfig([
  cspell.recommended,
  prettierRecommended,
]);
