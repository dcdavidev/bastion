# CLI API Reference

Bastion CLI provides a powerful interface for all vault operations.

## Authentication & Setup

- **`bastion init`**: The smart interactive wizard to initialize the database, migrations, and admin account.
- **`bastion login`**: Authenticate with the server and store the session token.
  - `--url, -u`: Server URL (default: `http://localhost:8287` or from `BASTION_HOST` env).
  - `--email, -e`: Email address (interactive prompt if omitted).
- **`bastion version`**: Print the version number of the Bastion CLI.
- **`bastion create superuser`**: Initialize the vault with admin credentials and master keys manually.
- **`bastion create jwtsecret`**: Generate and save a new random JWT secret.
- **`bastion create masterkey`**: Generate a new 32-byte random master key (hex).

## Resource Management

- **`bastion create client`**: Create a new tenant/client in the vault.
- **`bastion create project`**: Create a new project for a client with a dedicated encrypted data key.
- **`bastion create collaborator`**: Create a restricted user with specific access roles.
- **`bastion list clients`**: Display all clients in the vault.
- **`bastion list projects`**: List all projects belonging to a specific client.
  - `--client, -c`: Client ID (required).

## Secret Operations

- **`bastion set [KEY] [VALUE]`**: Encrypt and store a secret in a project.
  - `--project, -p`: Project ID (required).
  - `--key, -k`: Secret key name.
  - `--value, -v`: Secret value.
  - _Note: If KEY or VALUE are omitted, you will be prompted interactively with secure masking._
- **`bastion run --project <ID> -- <command>`**: Inject all decrypted secrets from a project as environment variables into the specified command.
  - `--project, -p`: Project ID to fetch secrets from (required).
  - _Note: Requires your Master Password to decrypt the data key locally._

## Global Flags

- `--url, -u`: Override the Bastion server URL for any command.
- `--help, -h`: Help for any command.

## Environment Variables

The CLI respects the following environment variables:

| Variable               | Description                                      | Default                   |
| :--------------------- | :----------------------------------------------- | :------------------------ |
| `BASTION_HOST`         | The base URL of the Bastion server.              | `http://localhost:8287`   |
| `BASTION_DATABASE_URL` | PostgreSQL connection string (used by `init`).   | -                         |
| `BASTION_STORE_DIR`    | Path to the password store directory for `pass`. | `~/.config/bastion/store` |
