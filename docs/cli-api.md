# CLI API Reference

Bastion CLI provides a powerful interface for all vault operations.

## Authentication & Setup

- **`bastion login`**: Authenticate with the server.
  - `--url, -u`: Server URL (default: from BASTION_HOST env).
  - `--username, -n`: Username (interactive prompt if omitted).
- **`bastion create superuser`**: Initialize the vault with admin credentials and master keys.
- **`bastion create jwtsecret`**: Generate and save a new random JWT secret.
- **`bastion create masterkey`**: Generate a new 32-byte random master key (hex).

## Resource Management

- **`bastion create client`**: Create a new tenant/client in the vault.
- **`bastion create project`**: Create a new project for a client with a dedicated encrypted data key.
- **`bastion create collaborator`**: Create a restricted user with specific access roles.
- **`bastion list clients`**: Display all clients in the vault.
- **`bastion list projects --client <ID>`**: List all projects belonging to a specific client.

## Secret Operations

- **`bastion set --project <ID> --key <KEY> --value <VALUE>`**: Encrypt and store a secret.
  - Can be used positionally: `bastion set <KEY> <VALUE> --project <ID>`.
- **`bastion run --project <ID> -- <command>`**: Inject all decrypted secrets from a project as environment variables into the specified command.
  - `--project, -p`: Project ID to fetch secrets from.

## Global Flags

- `--url, -u`: Override the Bastion server URL for any command.
