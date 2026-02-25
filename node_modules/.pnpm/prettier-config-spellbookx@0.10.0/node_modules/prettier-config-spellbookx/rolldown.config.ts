import { defineConfig } from 'rolldown';

export default defineConfig([
  {
    input: ['src/index.ts', 'src/types.ts'],
    output: {
      format: 'esm',
      dir: 'dist',
      entryFileNames: '[name].js',
      chunkFileNames: '[name]-[hash].js',
      assetFileNames: '[name]-[hash][extname]',
    },
    platform: 'node',
  },
]);
