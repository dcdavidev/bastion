import type { Linter } from 'eslint';
import { defineConfig, globalIgnores } from 'eslint/config';

const configIgnores: Linter.Config[] = defineConfig([
  globalIgnores([
    // Node / JS / TS build & lock files
    '**/node_modules/**',
    '**/dist/**',
    '**/build/**',
    '**/out/**',
    '**/out-tsc/**',
    '**/coverage/**',
    '**/.astro/**',
    '**/.next/**',
    '**/.nuxt/**',
    '**/.svelte-kit/**',
    '**/.turbo/**',
    '**/.pnpm/**',
    '**/.npm/**',
    '**/.yarn/**',
    '**/.pnp/**',
    '**/package-lock.json',
    '**/pnpm-lock.yaml',
    '**/yarn.lock',
    '**/bun.lockb',

    // Python environments and cache
    '**/__pycache__/**',
    '**/*.pyc',
    '**/*.pyo',
    '**/.mypy_cache/**',
    '**/.pytest_cache/**',
    '**/.tox/**',
    '**/.ruff_cache/**',
    '**/.venv/**',
    '**/venv/**',
    '**/env/**',
    '**/.ipynb_checkpoints/**',

    // Rust build artifacts
    '**/target/**',
    '**/.cargo/**',
    '**/Cargo.lock',

    // Go dependencies and cache
    '**/go.sum',
    '**/go.work',
    '**/go.work.sum',
    '**/vendor/**',
    '**/.gopath/**',
    '**/.cache/go-build/**',

    // Miscellaneous project artifacts
    '**/*LICENSE*',
    '**/*.log',
    '**/.cache/**',
    '**/.temp/**',
    '**/.tmp/**',
    '**/.DS_Store',
    '**/.idea/**',
    '**/.vscode/**',
    '**/.nx/**',
    '**/.cspell/**',
    '**/.cursor/**',
    '**/.history/**',
    '**/.terraform/**',
    '**/.devcontainer/**',
    '**/.direnv/**',
    '**/.editorconfig',
    '**/.eslintcache',
    '**/.babelrc',
    '**/.prettier*',
    '**/.sass-cache/**',
    '**/.gradle/**',
    '**/.docker/**',
    '**/.kube/**',
    '**/.git/**',
    '**/.svn/**',
    '**/.hg/**',
    '**/.gemini/**',
    '**/.codacy/**',
    '.github/instructions/**',

    // Windows system artifacts
    '**/*:Zone.Identifier',
    '**/Thumbs.db',
    '**/desktop.ini',
    '**/$RECYCLE.BIN/**',
    '**/System Volume Information/**',
    '**/pagefile.sys',
    '**/swapfile.sys',
    '**/hiberfil.sys',

    // Other unsupported files
    '**/*.hbs',
  ]),
]);

export default configIgnores;
