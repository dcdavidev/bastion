import type { Linter } from 'eslint';
import { defineConfig } from 'eslint/config';
import importPlugin from 'eslint-plugin-import';
import jsdoc from 'eslint-plugin-jsdoc';
import nodeDependencies from 'eslint-plugin-node-dependencies';
import prettierPlugin from 'eslint-plugin-prettier';
import prettierRecommended from 'eslint-plugin-prettier/recommended';
import simpleImportSort from 'eslint-plugin-simple-import-sort';
import unicorn from 'eslint-plugin-unicorn';
import globals from 'globals';
import tseslint from 'typescript-eslint';
import js from '@eslint/js';
import eslintReact from '@eslint-react/eslint-plugin';

export const configJavascript: Linter.Config[] = defineConfig([
  // --- Common for JS/TS/JSX/TSX ---
  {
    files: ['**/*.{js,jsx,mjs,cjs,ts,tsx,mts,cts}'],
    plugins: {
      jsdoc,
      import: importPlugin,
      prettier: prettierPlugin,
      'simple-import-sort': simpleImportSort,
    },
    extends: [
      js.configs.recommended,
      ...tseslint.configs.recommended,
      unicorn.configs.recommended,
    ],
    rules: {
      'prettier/prettier': 'error',
      'no-unused-vars': 'off',
      '@typescript-eslint/no-unused-vars': [
        'error',
        { argsIgnorePattern: '^_', varsIgnorePattern: '^_' },
      ],

      // Enforce structured and grouped imports using simple-import-sort
      'sort-imports': 'off',
      'simple-import-sort/exports': 'error',
      'simple-import-sort/imports': [
        'error',
        {
          groups: [
            // dotenv & dotenvx packages
            ['^@dotenvx/dotenvx', '^dotenv'],

            // Side-effect imports (e.g. polyfills)
            [String.raw`^\u0000`],

            // Env imports for React Native, Expo, etc.
            ['^@env'],

            // Node.js built-in modules
            [
              '^assert',
              '^buffer',
              '^child_process',
              '^cluster',
              '^console',
              '^constants',
              '^crypto',
              '^dgram',
              '^dns',
              '^domain',
              '^events',
              '^fs',
              '^http',
              '^https',
              '^inspector',
              '^module',
              '^net',
              '^os',
              '^path',
              '^perf_hooks',
              '^process',
              '^punycode',
              '^querystring',
              '^readline',
              '^repl',
              '^stream',
              '^string_decoder',
              '^timers',
              '^tls',
              '^tty',
              '^url',
              '^util',
              '^v8',
              '^vm',
              '^zlib',
            ],

            // Node.js built-in modules #2
            ['^node:'],

            // Node.js backend frameworks
            [
              '^@fastify',
              '^@hapi/hapi',
              '^@koa/koa',
              '^@nestjs',
              '^express',
              '^fastify',
              '^hapi',
              '^koa',
              '^loopback',
              '^nest',
              '^sails',
            ],

            // Node.js backend middlewares & utilities
            [
              '^body-parser',
              '^connect-redis',
              '^cookie-parser',
              '^cors',
              '^express-rate-limit',
              '^express-session',
              '^helmet',
              '^morgan',
              '^passport',
              String.raw`^pino-`,
              '^redis',
              '^winston',
            ],

            // UI Frameworks (React, Vue, Svelte, etc.)
            ['^@angular', '^react', '^solid-js', '^svelte', '^vue'],

            // React specific packages
            [
              '^@tanstack/router',
              '^react-dom',
              '^react-helmet',
              '^react-intl',
              '^react-router',
              '^react-router-dom',
            ],

            // Full-stack/SSR frameworks (Next.js, Remix, etc.)
            ['^@nuxt/kit', '^@remix-run', '^@sveltejs/kit', '^gatsby', '^next'],

            // React Native & Expo
            [
              '^@expo',
              String.raw`^@expo\/`,
              '^@react-native',
              '^expo',
              String.raw`^expo-`,
              '^react-native',
              '^react-navigation',
            ],

            // State management libraries
            [
              '^@reduxjs/toolkit',
              '^jotai',
              '^mobx',
              '^recoil',
              '^redux',
              '^valtio',
              '^zustand',
            ],

            // Data-fetching libraries
            [
              '^@apollo/client',
              '^@tanstack/react-query',
              '^axios',
              '^graphql',
              '^swr',
            ],

            // UI libraries & design systems
            [
              '^@chakra-ui',
              '^@headlessui/react',
              '^@lottiefiles',
              '^@material-ui',
              '^@mui',
              '^@nextui-org/react',
              '^@radix-ui',
              '^@pittorica/pitto',
              '^radix-ui',
              '^antd',
              '^framer-motion',
              '^native-base',
              '^react-native-paper',
              '^shadcn-ui',
              '^tailwindcss',
            ],

            // CSS-in-JS & utility libraries
            [
              '^@emotion',
              '^class-variance-authority',
              '^clsx',
              '^lucide-react',
              '^styled-components',
              '^tailwind-merge',
              '^twin.macro',
              '^tw-animate-css',
            ],

            // Common icon packages
            [
              String.raw`^@expo\/vector-icons`,
              '^@fortawesome',
              '^@tabler/icons-react',
              '^lucide',
              '^react-feather',
              '^react-icons',
              '^react-native-feather',
              '^react-native-vector-icons',
            ],

            // Testing libraries and utilities
            [
              '^@testing-library',
              '^cypress',
              '^jest',
              '^playwright',
              '^vitest',
            ],

            // Generic third-party packages (npm scope and plain)
            ['^[a-z]', String.raw`^@\w`],

            // Monorepo/workspace scoped packages
            ['@org/', '^@my-org/', '^@workspace/', '^@repo/'],

            // Asset imports (images, fonts, etc.)
            [
              String.raw`^.+\.(avi|mkv|mov|mp4|webm)$`,
              String.raw`^.+\.(mp3|ogg|wav|weba)$`,
              String.raw`^.+\.(gif|jpe?g|png|svg|webp)$`,
              String.raw`^.+\.lottie$`,
              String.raw`^.+\.(eot|otf|ttf|woff|woff2)$`,
            ],

            // JSON files
            [String.raw`^.+\.json$`],

            // Stylesheets (css, scss, less, etc.)
            [String.raw`^.+\.less$`, String.raw`^.+\.s?css$`],

            // Relative imports (parent, sibling, current)
            [String.raw`^\.?\.\/`],
          ],
        },
      ],

      'unicorn/no-null': 'off',
      'unicorn/prevent-abbreviations': 'off',
    },
  },

  // --- Node.js Specifics (TS/JS files) ---
  {
    files: ['**/*.{ts,js,mts,mjs,cts,cjs}'],
    ignores: ['**/*.{tsx,jsx}'],
    extends: [
      nodeDependencies.configs['flat/recommended'],
      jsdoc.configs['flat/contents-typescript-flavor'],
      jsdoc.configs['flat/logical-typescript-flavor'],
    ],
    languageOptions: {
      globals: { ...globals.node },
    },
    rules: {
      'unicorn/filename-case': ['error', { case: 'kebabCase' }],
      'jsdoc/require-jsdoc': [
        'warn',
        {
          require: {
            FunctionDeclaration: true,
            MethodDefinition: true,
            ClassDeclaration: true,
          },
        },
      ],
      'jsdoc/require-description': 'warn',
    },
  },

  // --- React Specifics (TSX/JSX files) ---
  {
    files: ['**/*.{tsx,jsx}'],
    extends: [eslintReact.configs['recommended-typescript']],
    languageOptions: {
      sourceType: 'module',
      globals: { ...globals.browser, ...globals.es2022 },
      parser: tseslint.parser,
      parserOptions: {
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
      },
    },
    rules: {
      // Filename case OFF for React components
      'unicorn/filename-case': 'off',
      // JSDoc OFF for React components
      'jsdoc/require-jsdoc': 'off',
      'jsdoc/require-description': 'off',
      'jsdoc/require-param': 'off',
      'jsdoc/require-returns': 'off',
      'jsdoc/check-tag-names': 'off',
    },
  },

  // --- CommonJS Overrides ---
  {
    files: ['**/*.{cjs,cts}'],
    languageOptions: { globals: { ...globals.commonjs } },
    rules: {
      '@typescript-eslint/no-require-imports': 'off',
      'unicorn/prefer-module': 'off',
    },
  },

  prettierRecommended,
]);
