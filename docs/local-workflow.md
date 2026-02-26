# Local Development Workflow 🛠️

This guide walks you through the daily workflow for managing and using secrets with the Bastion CLI.

---

## 🔐 1. Authentication

The CLI securely stores your session token in `~/.bastion/token` after a successful login.

### Login

```bash
./bastion login
```

- **Server**: `http://localhost:8287` (default).
- **Email**: Your admin or collaborator email.
- **Password**: Your login password.

---

## 📂 2. Organizing Resources

Bastion uses a three-level hierarchy: **Client** → **Project** → **Secrets**.

### Create a Client

Logical groups for your environments (e.g., "Internal Tools" or "Client Alpha").

```bash
./bastion create client
```

### Create a Project

Specific environments belonging to a client (e.g., "Production", "Staging").

```bash
./bastion create project
```

_Tip: You will need the **Client ID**, which you can find via `bastion list clients`._

---

## 🔑 3. Managing Secrets

Secrets are encrypted locally before transmission.

### Find IDs

```bash
# List all clients
./bastion list clients

# List projects for a specific client
./bastion list projects --client <CLIENT_ID>
```

### Store a Secret

You can use positional arguments or flags. If you omit values, the CLI will prompt you securely.

```bash
# Positional (Fastest)
./bastion set DATABASE_URL "postgres://..." -p <PROJECT_ID>

# Using flags
./bastion set --project <PROJECT_ID> --key DATABASE_URL --value "postgres://..."

# Interactive (Safest, hides input)
./bastion set -p <PROJECT_ID>
```

_Note: You must enter your **Master Password** to derive the local encryption key._

---

## 🚀 4. Injecting Secrets (Runtime)

The `run` command fetches decrypted secrets and injects them as environment variables directly into your application's process.

### Usage

```bash
./bastion run --project <PROJECT_ID> -- <your-command>
```

### Examples

```bash
# Node.js
./bastion run -p npm start < PROJECT_ID > --
```

```bash
# Go / Binaries
./bastion run -p ./my-app < PROJECT_ID > --
```

**Benefits:**

1. **Memory Only**: Secrets never touch the disk in plaintext.
2. **Zero Config**: No `.env` files to leak or manage.
3. **Audit**: Every secret fetch is logged on the server.

---

## ⏭️ Next Steps

- **[CLI Reference](cli-api.md)**: Full list of flags and commands.
- **[Configuration](config.md)**: Customize your host and storage settings.
