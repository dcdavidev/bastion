import { defineConfig } from 'rolldown';

const esmBanner = `
import { fileURLToPath as _fileURLToPath } from 'node:url';
import { dirname } from 'node:path';
import { createRequire as _createRequire } from 'node:module';

const __filename = _fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const require = _createRequire(import.meta.url);
`;

export default defineConfig([
  {
    input: ['src/index.ts', 'src/types.ts'],
    output: {
      format: 'esm',
      dir: 'dist',
      entryFileNames: '[name].js',
      chunkFileNames: '[name]-[hash].js',
      assetFileNames: '[name]-[hash][extname]',
      banner: esmBanner,
    },
    platform: 'node',
    external: (id) => {
      if (id.startsWith('node:')) return true;
      if (id.startsWith('.') || id.startsWith('/')) return false;
      return true;
    },
  },
]);
