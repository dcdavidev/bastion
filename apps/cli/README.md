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
bastion run -p npm start < PROJECT_ID > --
```

## CLI Reference 📖

### Setup & Auth

- `bastion init`: Interactive setup wizard for database and admin.
- `bastion login --email <email>`: Authenticate and store session token.
- `bastion version`: Show current version and check for updates.

### Management

- `bastion create client`: Register a new tenant.
- `bastion create project`: Add an E2EE environment to a client.
- `bastion list clients`: Show all registered clients.
- `bastion list projects -c <CLIENT_ID>`: Show projects for a client.

### Secrets

- `bastion set <KEY> <VALUE> -p <PROJECT_ID>`: Encrypt and store a secret.
- `bastion run -p <PROJECT_ID> -- <command>`: Inject secrets and execute a command.

## Documentation

For the full documentation and CLI reference, visit the [official repository](https://github.com/dcdavidev/bastion).

## License

MIT © [Davide Di Criscito](https://github.com/dcdavidev)
