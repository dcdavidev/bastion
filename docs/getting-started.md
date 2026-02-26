# Getting Started with Bastion 🚀

Welcome to Bastion! This guide will walk you through the initial steps to set up your own E2EE (End-to-End Encrypted) secrets vault for local development.

## 📋 Prerequisites

Before you begin, ensure your development environment meets these requirements:

### 🛠️ Required Tools

| Tool           | Version  | Purpose                                |
| :------------- | :------- | :------------------------------------- |
| **Go**         | `1.24+`  | Core backend logic and CLI.            |
| **Node.js**    | `24+`    | React dashboard and frontend tooling.  |
| **pnpm**       | `latest` | Fast, disk-efficient package manager.  |
| **PostgreSQL** | `15+`    | A running instance (local or managed). |

### 💡 Knowledge Base

- Basic familiarity with the command line.
- Knowledge of PostgreSQL connection strings.
- Experience managing a PostgreSQL instance.
- Basic Git usage.

---

## 🏗️ 1. Initial Setup

### Clone the Repository

Start by cloning the official repository and navigating into the project directory:

```bash
git clone https://github.com/dcdavidev/bastion.git
cd bastion
```

### Install Dependencies

Bastion is a monorepo managed by `pnpm`. Install all backend and frontend dependencies with a single command:

```bash
pnpm install
```

### Workspace Structure

Once installed, you'll see several key directories:

- `apps/server`: The Go backend.
- `apps/web`: The React dashboard.
- `apps/cli`: The Bastion command-line tool.
- `packages/core`: Shared logic, crypto, and database models.

---

## ✅ 2. Verification

To ensure your environment is correctly configured, run the build command for the entire workspace:

```bash
pnpm build
```

This will:

1. Compile the **CLI** tool.
2. Build the **Backend Server**.
3. Bundle the **Web Dashboard**.

If the build completes without errors, you are ready to proceed!

---

## ⏭️ Next Steps

Now that you have the source code and dependencies ready, it's time to configure and initialize your vault:

1. **[Configure your environment](config.md)**: Learn about the required environment variables.
2. **[Initialize the Vault](initial-steps.md)**: Set up your database and create your superuser account.
3. **[Manage Secrets](local-workflow.md)**: Start using the CLI to store and inject secrets.
