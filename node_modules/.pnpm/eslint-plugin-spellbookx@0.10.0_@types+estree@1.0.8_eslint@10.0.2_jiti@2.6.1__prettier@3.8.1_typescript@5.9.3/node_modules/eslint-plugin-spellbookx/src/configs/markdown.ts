import { Linter } from 'eslint';
import { defineConfig } from 'eslint/config';
import prettierRecommended from 'eslint-plugin-prettier/recommended';
import markdown from '@eslint/markdown';

export const configMarkdown: Linter.Config[] = defineConfig([
  {
    files: ['**/*.md'],
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    plugins: { markdown: markdown as any },
    language: 'markdown/gfm',
    extends: ['markdown/recommended'],
  },

  prettierRecommended,
]);
