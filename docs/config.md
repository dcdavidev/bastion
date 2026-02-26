# Configuration Guide

Bastion relies exclusively on environment variables for configuration. This ensures that sensitive credentials are never accidentally committed to version control and simplifies containerized deployments.

---

## Environment Variables

These variables configure both the Bastion server and the CLI.

| Variable               | Description                                                                          | Default                   | Used By             |
| :--------------------- | :----------------------------------------------------------------------------------- | :------------------------ | :------------------ |
| `BASTION_HOST`         | The base URL of the Bastion server.                                                  | `http://localhost:8287`   | CLI                 |
| `BASTION_PORT`         | The port the server listens on.                                                      | `8287`                    | Server              |
| `BASTION_DATABASE_URL` | PostgreSQL connection string (fallback to `DATABASE_URL`).                           | _(Required)_              | Server / CLI (init) |
| `BASTION_JWT_SECRET`   | 32-byte hex string used to sign session tokens.                                      | _(Required)_              | Server              |
| `BASTION_UI_DIR`       | Path to the built frontend assets (`apps/web/build/client`).                         | _(Auto-detected)_         | Server              |
| `BASTION_STORE_DIR`    | Path to the password store directory for `pass` integration.                         | `~/.config/bastion/store` | CLI (init)          |

---

## Local Development

### .env File Support
The Bastion server and CLI automatically load variables from a `.env` file in the current directory if it exists.

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
    image: dcdavidev/bastion:latest
    ports:
      - "8287:8287"
    environment:
      - BASTION_DATABASE_URL=postgres://user:pass@db:5432/bastion
      - BASTION_JWT_SECRET=${BASTION_JWT_SECRET}
    restart: unless-stopped
```
