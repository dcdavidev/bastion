# Configuration Guide

Bastion relies exclusively on environment variables for configuration. This ensures that sensitive credentials are never accidentally committed to version control and simplifies containerized deployments.

---

## Environment Variables

These variables configure both the Bastion server and the CLI.

| Variable               | Description                                                                          | Default                   | Used By             |
| :--------------------- | :----------------------------------------------------------------------------------- | :------------------------ | :------------------ |
| `BASTION_HOST`         | The base URL of the Bastion server.                                                  | `http://localhost:8287`   | CLI                 |
| `BASTION_PORT`         | The port the server listens on.                                                      | `8287`                    | Server              |
| `BASTION_DATABASE_URL` | PostgreSQL connection string for an **empty database** (locally or managed service). | _(Required)_              | Server / CLI (init) |
| `BASTION_JWT_SECRET`   | 32-byte hex string used to sign session tokens.                                      | _(Required)_              | Server              |
| `BASTION_MASTER_KEY`   | (Optional) The master key for cryptographic operations.                              | _(Internal)_              | Server              |
| `BASTION_STORE_DIR`    | (Optional) Path to the password store directory for `pass` integration.              | `~/.config/bastion/store` | CLI (init)          |

---

## Getting Started with Environment Variables

### Local Development (Linux/macOS)

You can set these variables in your shell profile (e.g., `.bashrc` or `.zshrc`):

```bash
export BASTION_HOST="http://localhost:8287"
export BASTION_DATABASE_URL="postgres://user:pass@localhost:5432/bastion"
export BASTION_JWT_SECRET="your-32-byte-hex-secret"
```

### Docker Compose

In a `docker-compose.yaml` file, define them under the `environment` key:

```yaml
services:
  backend:
    image: bastion-server
    environment:
      - BASTION_DATABASE_URL=postgres://bastion:password@db:5432/bastion
      - BASTION_JWT_SECRET=your-secret
      - BASTION_PORT=8080
```
