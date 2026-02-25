# eslint-plugin-spellbookx

An ESLint plugin containing shared configurations for the Spellbookx project. This package bundles and pre-configures a variety of essential ESLint plugins for a consistent and high-quality codebase.

## Description

This ESLint plugin provides a set of shareable configurations that are fine-tuned for modern JavaScript/TypeScript development. It includes rules for:

- **JavaScript**: Best practices, error prevention, and code style.
- **TypeScript**: Static typing, and modern language features.
- **JSON**: Linting for JSON files.
- **Markdown**: Linting for Markdown files.
- **Imports**: Sorting and organizing imports.
- **JSDoc**: Enforcing JSDoc comments.
- **Prettier**: Integrating Prettier for code formatting.
- **CSpell**: Spell-checking within your code.

## Why use this?

This plugin saves you the hassle of manually installing and configuring multiple ESLint plugins. By providing a centralized and opinionated configuration, it helps to:

- Enforce a consistent code style across the entire monorepo.
- Catch potential errors and bugs early in the development process.
- Improve code readability and maintainability.

## Installation

Install the package and its peer dependencies using your favorite package manager:

**npm**

```bash
npm install --save-dev eslint-plugin-spellbookx eslint
```

**pnpm**

```bash
pnpm add -D eslint-plugin-spellbookx eslint
```

**yarn**

```bash
yarn add -D eslint-plugin-spellbookx eslint
```

**bun**

```bash
bun add -d eslint-plugin-spellbookx eslint
```

## Usage

To use this plugin, create an `eslint.config.mjs` file in your project's root and extend the desired configurations.

### Recommended Configuration

The `recommended` configuration includes all the essential rules for a typical JavaScript/TypeScript project.

```javascript
import spellbookx from 'eslint-plugin-spellbookx';

export default [
  ...spellbookx.configs.recommended,
  // Add your custom rules here
];
```

### Advanced Configuration

You can also create a more tailored configuration by picking and choosing the configs you need:

```javascript
import spellbookx from 'eslint-plugin-spellbookx';

export default [
  ...spellbookx.configs.ignores,
  ...spellbookx.configs.javascript,
  ...spellbookx.configs.typescript,
  ...spellbookx.configs.json,
  ...spellbookx.configs.markdown,
  ...spellbookx.configs.vuoto,
  ...spellbookx.configs.cspell,
  // any other configuration
];
```

## License

This project is licensed under the MIT License.

**Copyright (c) 2026 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
