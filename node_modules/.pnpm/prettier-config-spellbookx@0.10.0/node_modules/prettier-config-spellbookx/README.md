# prettier-config-spellbookx

Shared [Prettier](https://prettier.io/) configuration for the Spellbookx project. This package provides a consistent and opinionated code formatting setup for a wide range of file types.

## Description

This package offers a collection of Prettier configurations that are tailored for different development scenarios. It includes support for:

- **Base**: A solid foundation for formatting JavaScript, TypeScript, JSON, and more.
- **Astro**: Additional rules for Astro projects.
- **Prisma**: Formatting for Prisma schema files.
- **Tailwind CSS**: Sorting and organizing Tailwind CSS classes.
- A variety of other file types, including TOML, INI, XML, and shell scripts.

## Why use this?

Using this shared Prettier configuration ensures that:

- Code formatting is consistent across the entire monorepo.
- The formatting rules are easily maintainable and extendable.
- The need for manual formatting is eliminated, saving development time.

## Installation

Install the package and its peer dependencies using your favorite package manager:

**npm**

```bash
npm install --save-dev prettier-config-spellbookx prettier
```

**pnpm**

```bash
pnpm add -D prettier-config-spellbookx prettier @prettier/plugin-xml prettier-plugin-astro prettier-plugin-ini prettier-plugin-packagejson prettier-plugin-prisma prettier-plugin-properties prettier-plugin-sh prettier-plugin-tailwindcss prettier-plugin-toml
```

**yarn**

```bash
yarn add -D prettier-config-spellbookx prettier
```

**bun**

```bash
bun add -d prettier-config-spellbookx prettier
```

## Usage

To use this configuration, create a `prettier.config.mjs` file in your project's root and import the desired configuration.

### Base Configuration

The `base` configuration is a great starting point for most projects.

```javascript
import { base } from 'prettier-config-spellbookx';

export default base;
```

### Combined Configurations

You can easily combine configurations to suit your project's needs. For example, to use the Astro, Prisma, and Tailwind CSS configurations together:

```javascript
import { astroPrismaTailwind } from 'prettier-config-spellbookx';

export default astroPrismaTailwind;
```

### Available Configurations

Here is a list of all the available configurations:

- `base`
- `astro`
- `astro-prisma`
- `astro-tailwind`
- `astro-prisma-tailwind`
- `tailwind`
- `prisma`
- `prisma-tailwind`

### Troubleshooting

If node complains that some dependencies are missing, please install them:

```bash
npm install -D prettier @prettier/plugin-xml prettier-plugin-astro prettier-plugin-ini prettier-plugin-packagejson prettier-plugin-prisma prettier-plugin-properties prettier-plugin-sh prettier-plugin-tailwindcss prettier-plugin-toml
```

```bash
pnpm add -D prettier @prettier/plugin-xml prettier-plugin-astro prettier-plugin-ini prettier-plugin-packagejson prettier-plugin-prisma prettier-plugin-properties prettier-plugin-sh prettier-plugin-tailwindcss prettier-plugin-toml
```

```bash
yarn add -D prettier @prettier/plugin-xml prettier-plugin-astro prettier-plugin-ini prettier-plugin-packagejson prettier-plugin-prisma prettier-plugin-properties prettier-plugin-sh prettier-plugin-tailwindcss prettier-plugin-toml
```

```bash
bun add -d prettier @prettier/plugin-xml prettier-plugin-astro prettier-plugin-ini prettier-plugin-packagejson prettier-plugin-prisma prettier-plugin-properties prettier-plugin-sh prettier-plugin-tailwindcss prettier-plugin-toml
```

## License

This project is licensed under the MIT License.

**Copyright (c) 2026 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
