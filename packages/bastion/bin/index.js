#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');

// Mapping OS and architecture to the folders created by GoReleaser
const platforms = {
  'darwin': {
    'x64': 'darwin_amd64/bastion',
    'arm64': 'darwin_arm64/bastion'
  },
  'linux': {
    'x64': 'linux_amd64/bastion',
    'arm64': 'linux_arm64/bastion'
  },
  'win32': {
    'x64': 'windows_amd64/bastion.exe',
    'arm64': 'windows_arm64/bastion.exe'
  }
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
