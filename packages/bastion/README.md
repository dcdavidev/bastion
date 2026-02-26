# @dcdavidev/bastion-cli 🛡️

The official Node.js wrapper for **Bastion**, a secure E2EE (End-to-End Encrypted) secrets vault.

## Features
- **Zero-Config Secrets**: Inject secrets directly into your application's environment.
- **E2EE**: All secrets are encrypted client-side; the server never sees your plaintext data.
- **Blind Backend**: Secure, self-hosted architecture.
- **Cross-Platform**: Built-in binaries for Linux, macOS, and Windows.

## Installation

```bash
npm install -g @dcdavidev/bastion-cli
```

## Quick Start

### 1. Initialize your vault
```bash
bastion init
```

### 2. Login
```bash
bastion login --email your@email.com
```

### 3. Run your application with secrets
```bash
bastion run -p <PROJECT_ID> -- npm start
```

## Documentation
For the full documentation and CLI reference, visit the [official repository](https://github.com/dcdavidev/bastion).

## License
MIT © [Davide Di Criscito](https://github.com/dcdavidev)
