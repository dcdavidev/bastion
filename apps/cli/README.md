# @dcdavidev/bastion-cli 🛡️

The official Node.js wrapper for **Bastion**, a secure E2EE (End-to-End Encrypted) secrets vault.

## Features

- **Zero-Config Secrets**: Inject secrets directly into your application's environment.
- **E2EE**: All secrets are encrypted client-side.
- **Profiles**: Manage multiple environments (dev, staging, prod) with `bastion profile`.
- **CI/CD Ready**: All commands support flags to bypass interactive prompts.

## Installation

```bash
npm install -g @dcdavidev/bastion-cli
```

## Quick Start

### 1. Initialize and Configure

```bash
bastion init
```

### 2. Login

```bash
bastion login --email your@email.com --password your-pass
```

### 3. Run your application with secrets

```bash
bastion run --project PROJECT_ID --password vault-pass -- node main.js
```

## CLI Reference 📖

### Setup & Profiles

- `bastion init`: Interactive setup wizard for server and profile configuration.
- `bastion profile list`: List all configured environments.
- `bastion profile add <name> <url>`: Add a new server environment.
- `bastion profile use <name>`: Switch the active environment.
- `bastion login`: Authenticate and store session token for the active profile.

### Management

- `bastion create client --name <name>`: Register a new tenant.
- `bastion create project --client <id> --name <name>`: Add an E2EE project.
- `bastion list clients`: Show all registered clients.
- `bastion list projects --client <id>`: Show projects for a client.

### Secrets

- `bastion set <KEY> <VALUE> --project <id> --password <vault-pass>`: Encrypt and store a secret.
- `bastion run --project <id> --password <vault-pass> -- <command>`: Inject secrets and execute a command.

## Global Flags

- `--profile, -P <name>`: Use a specific profile for the command.
- `--url, -u <url>`: Override the Bastion server URL for any command.
- `--help, -h`: Help for any command.

## Documentation

For the full documentation and CLI reference, visit the [official repository](https://github.com/dcdavidev/bastion).

## License

MIT © [Davide Di Criscito](https://github.com/dcdavidev)
