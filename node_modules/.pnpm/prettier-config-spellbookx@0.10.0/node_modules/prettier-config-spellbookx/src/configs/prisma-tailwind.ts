import type { Config } from 'prettier';

export const prismaTailwind: Config = {
  plugins: [
    'prettier-plugin-sh',
    'prettier-plugin-toml',
    'prettier-plugin-ini',
    'prettier-plugin-packagejson',
    'prettier-plugin-properties',
    'prettier-plugin-prisma',
    '@prettier/plugin-xml',
    'prettier-plugin-tailwindcss', // must be loaded last
  ],

  trailingComma: 'es5',
  tabWidth: 2,
  semi: true,
  singleQuote: true,
  printWidth: 80,
  endOfLine: 'lf',
  bracketSpacing: true,

  overrides: [
    {
      files: '**/.czrc',
      options: {
        parser: 'json',
      },
    },
    {
      files: ['*.toml'],
      options: {
        printWidth: 100,
      },
    },
    {
      files: ['*.astro'],
      options: {
        parser: 'astro',
      },
    },
    {
      files: [
        '**/*.xml',
        '**/*.svg',
        '**/*.xhtml',
        '**/*.xsd',
        '**/*.xsl',
        '**/*.xslt',
        '**/*.plist',
        '**/*.axml',
        '**/*.xaml',
        '**/*.config',
        '**/*.resx',
        '**/*.csproj',
        '**/*.vbproj',
      ],
      options: {
        parser: 'xml',
      },
    },
  ],
};
