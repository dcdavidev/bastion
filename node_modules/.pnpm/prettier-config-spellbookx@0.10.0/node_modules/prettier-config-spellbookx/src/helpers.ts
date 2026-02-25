import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';

const packagePath = fileURLToPath(new URL('../package.json', import.meta.url));
export const packageJson = JSON.parse(readFileSync(packagePath, 'utf8')) as {
  name: string;
  version: string;
  peerDependencies: Record<string, string>;
};
