# CLI API Reference

Bastion CLI provides a powerful interface for all vault operations.

## Setup & Configuration

- **`bastion init`**: The smart interactive wizard to initialize the database, migrations, and admin account. Now also configures the local client profile.
- **`bastion profile`**: Manage multiple Bastion environments.
  - `list`: Show all configured profiles.
  - `add [NAME] [URL]`: Add a new server environment.
  - `use [NAME]`: Set the default profile for subsequent commands.
- **`bastion login`**: Authenticate with the server and store the session token.
  - `--url, -u`: Server URL.
  - `--email, -e`: Email address.
  - `--password, -p`: Password (avoids interactive prompt).
- **`bastion version`**: Print the version number and check for updates.

## Resource Management

- **`bastion create client`**: Create a new tenant/client in the vault.
  - `--name, -n`: Client name.
- **`bastion create project`**: Create a new project for a client.
  - `--client, -c`: Client ID (UUID).
  - `--name, -n`: Project name.
- **`bastion list clients`**: Display all clients in the vault.
- **`bastion list projects`**: List all projects for a specific client.
  - `--client, -c`: Client ID (optional, interactive prompt if omitted).

## Secret Operations

- **`bastion set [KEY] [VALUE]`**: Encrypt and store a secret in a project.
  - `--project, -p`: Project ID (required).
  - `--key, -k`: Secret key name.
  - `--value, -v`: Secret value.
  - `--password`: Admin password to unlock the vault (avoids interactive prompt).
- **`bastion run --project <ID> -- <command>`**: Inject all decrypted secrets from a project as environment variables.
  - `--project, -p`: Project ID (required).
  - `--password`: Password to unlock the vault.

## Global Flags

- `--profile, -P`: Use a specific profile for the command.
- `--url, -u`: Override the server URL for any command.
- `--version, -v`: Print the version number.
- `--help, -h`: Show help.

## Environment Variables

| Variable               | Description                                      | Default                   |
| :--------------------- | :----------------------------------------------- | :------------------------ |
| `BASTION_HOST`         | The base URL of the Bastion server.              | `http://localhost:8287`   |
| `BASTION_DATABASE_URL` | PostgreSQL connection string (used by `init`).   | -                         |

## Config File

The CLI stores its configuration (profiles and tokens) in `~/.bastion/config.yaml`.
