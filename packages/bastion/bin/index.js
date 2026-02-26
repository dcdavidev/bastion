#!/usr/bin/env node

import { spawn } from 'node:child_process';
import path from 'node:path';
import os from 'node:os';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Mapping OS and architecture to the folders created by GoReleaser
const platforms = {
  darwin: {
    x64: 'darwin_amd64/bastion-cli',
    arm64: 'darwin_arm64/bastion-cli',
  },
  linux: {
    x64: 'linux_amd64/bastion-cli',
    arm64: 'linux_arm64/bastion-cli',
  },
  win32: {
    x64: 'windows_amd64/bastion-cli.exe',
    arm64: 'windows_arm64/bastion-cli.exe',
  },
};

const platform = os.platform();
const arch = os.arch();

if (!platforms[platform] || !platforms[platform][arch]) {
  console.error(`Unsupported platform/architecture: ${platform}/${arch}`);
  process.exit(1);
}

const binaryPath = path.join(__dirname, platforms[platform][arch]);
const args = process.argv.slice(2);

const child = spawn(binaryPath, args, { stdio: 'inherit' });

child.on('close', (code) => {
  process.exit(code);
});
