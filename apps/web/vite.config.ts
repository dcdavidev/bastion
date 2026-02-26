import { defineConfig } from 'vite';
import tsconfigPaths from 'vite-tsconfig-paths';
import { reactRouter } from '@react-router/dev/vite';

import pkg from './package.json';

export default defineConfig({
  plugins: [reactRouter(), tsconfigPaths()],
  define: {
    'import.meta.env.VITE_BASTION_VERSION': JSON.stringify(pkg.version),
  },
  server: {
    port: 3500,
    proxy: {
      '/api': {
        target: 'http://localhost:8287',
        changeOrigin: true,
      },
    },
  },
});
