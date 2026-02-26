# Configuration Guide

Bastion uses environment variables for server configuration and a local configuration file for the CLI. This ensures that sensitive credentials are managed securely across different environments.

---

## CLI Configuration

The Bastion CLI stores its configuration, including profiles and session tokens, in a local YAML file.

- **Path**: `~/.bastion/config.yaml`
- **Managed by**: `bastion profile` and `bastion init` commands.

This file allows you to switch between different Bastion servers (e.g., development, staging, production) without re-entering credentials.

---

## Environment Variables

These variables configure the Bastion server and are used by the CLI during the `init` process.

| Variable               | Description                                                | Default                 | Used By             |
| :--------------------- | :--------------------------------------------------------- | :---------------------- | :------------------ |
| `BASTION_HOST`         | The base URL of the Bastion server.                        | `http://localhost:8287` | CLI (Fallback)      |
| `BASTION_PORT`         | The port the server listens on.                            | `8287`                  | Server              |
| `BASTION_DATABASE_URL` | PostgreSQL connection string (fallback to `DATABASE_URL`). | _(Required)_            | Server / CLI (init) |
| `BASTION_JWT_SECRET`   | 32-byte hex string used to sign session tokens.            | _(Required)_            | Server              |
| `BASTION_UI_DIR`       | Path to the built frontend assets.                         | `ui` (in Docker)        | Server              |

### Admin Fallback (Optional)

Bastion supports an environment-based admin fallback. This is useful for the first login before the database is initialized or as a recovery mechanism.

| Variable                      | Description                                |
| :---------------------------- | :----------------------------------------- |
| `BASTION_ADMIN_PASSWORD_HASH` | Argon2id hash of the admin password (hex). |
| `BASTION_ADMIN_PASSWORD_SALT` | 32-byte salt used for the hash (hex).      |

---

## Local Development

### .env File Support

Both the Bastion server and CLI automatically load variables from a `.env` file in the current directory if it exists.

### Example .env

```bash
BASTION_DATABASE_URL="postgres://user:pass@localhost:5432/bastion"
BASTION_JWT_SECRET="0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
BASTION_PORT=8287
```

---

## Docker Configuration

The unified Bastion container serves both the API and the UI on a single port.

### Docker Compose

```yaml
services:
  bastion:
    image: ghcr.io/dcdavidev/bastion:latest
    ports:
      - '8287:8287'
    environment:
      - BASTION_DATABASE_URL=postgres://user:pass@db:5432/bastion
      - BASTION_JWT_SECRET=${BASTION_JWT_SECRET}
      # Optional Admin Fallback
      - BASTION_ADMIN_PASSWORD_HASH=${BASTION_ADMIN_PASSWORD_HASH}
      - BASTION_ADMIN_PASSWORD_SALT=${BASTION_ADMIN_PASSWORD_SALT}
    restart: unless-stopped
```
