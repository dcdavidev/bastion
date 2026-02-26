# Getting Started with Bastion 🚀

Welcome to Bastion! This guide will walk you through the initial steps to set up your own E2EE (End-to-End Encrypted) secrets vault.

## 📋 Prerequisites

### 🛠️ Required Tools

| Tool           | Version  | Purpose                                   |
| :------------- | :------- | :---------------------------------------- |
| **Go**         | `1.24+`  | Core backend logic and CLI.               |
| **Node.js**    | `24+`    | React dashboard build and frontend tools. |
| **pnpm**       | `latest` | Fast, disk-efficient package manager.     |
| **PostgreSQL** | `15+`    | A running instance (local or managed).    |

---

## 🏗️ 1. Initial Setup

### Clone and Install

```bash
git clone https://github.com/dcdavidev/bastion.git
cd bastion
pnpm install
```

### Workspace Structure

- **`apps/server`**: Unified Go server (API + Web UI).
- **`apps/web`**: React-based Dashboard source code.
- **`apps/cli`**: The `bastion` command-line management tool.
- **`packages/core`**: Shared logic, crypto, and database migrations.

---

## ✅ 2. Build for Production

Bastion uses a **Unified Architecture**. The Go server is capable of directly serving the compiled Dashboard. To prepare the environment:

```bash
pnpm build
```

This command will generate:

1. The `bastion` executable (CLI).
2. The `bastion-server` executable (Backend).
3. The optimized web assets in `apps/web/build/client`.

---

## 🚀 3. Quick Launch

To immediately test the entire stack (API + Dashboard) after building:

```bash
pnpm dev:server
```

You can now access the unified portal at **`http://localhost:8287`**.

---

## ⏭️ Next Steps

1. **[Configure Environment](config.md)**: Set up your database connection and JWT secrets.
2. **[Initialize the Vault](initial-steps.md)**: Use the `bastion init` wizard to configure the database and create your admin account.
3. **[CLI Reference](cli-api.md)**: Explore the available commands to manage your secrets.
