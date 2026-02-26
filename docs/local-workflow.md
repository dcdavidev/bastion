# Local Development Workflow 🛠️

This guide walks you through the typical daily workflow for managing secrets using the Bastion CLI.

---

## 🔐 1. Authentication

Before performing any operation, you must authenticate with the Bastion server. The CLI will securely store your session token locally in `~/.bastion/token`.

### Login

```bash
./bastion login
```

- **Server URL**: Usually `http://localhost:8080` for local development.
- **Username**: Your admin username or collaborator name.
- **Password**: Your login password.

---

## 📂 2. Organizing Resources

Bastion uses a three-level hierarchy to keep your secrets organized:
**Client** → **Project** → **Secrets**.

### Create a Client

Think of a Client as a logical group (e.g., a customer, a department, or a large application suite).

```bash
./bastion create client
```

_Tip: You will be prompted for the client name. Use something descriptive like "AcmeCorp" or "InternalTools"._

### Create a Project

Projects belong to a Client and hold the actual secrets. Each project has its own unique, encrypted Data Key.

```bash
./bastion create project
```

_Tip: You will need the **Client ID** (a UUID) which you can find by running `bastion list clients`._

---

## 🔑 3. Managing Secrets

Once you have a project, you can start adding encrypted secrets to it.

### List Clients and Projects

To find the correct IDs for your commands:

```bash
# List all clients
./bastion list clients

# List projects for a specific client
./bastion list projects --client <CLIENT_ID>
```

### Store a Secret

You can store a secret using flags or positional arguments:

```bash
# Using flags
./bastion set --project <PROJECT_ID> --key DATABASE_URL --value "postgres://..."

# Using positional arguments (faster)
./bastion set DATABASE_URL "postgres://..." -p <PROJECT_ID>
```

_Note: You will be prompted for your **Master Password** to unlock the vault and encrypt the data._

---

## 🚀 4. Injecting Secrets (Runtime)

The most powerful feature of Bastion is the ability to inject secrets directly into your application's environment without ever writing them to a `.env` file.

### Basic Usage

```bash
./bastion run --project <PROJECT_ID> -- <your-command>
```

### Examples

```bash
# Run a Node.js application
./bastion run -p npm start < PROJECT_ID > --

# Run a Go binary
./bastion run -p ./my-app < PROJECT_ID > --

# Use a custom configuration file
./bastion run -p ./custom-config.yml -- python main.py < PROJECT_ID > -c
```

### Why use `run`?

1. **Security**: Secrets only exist in memory during the execution of the command.
2. **Convenience**: No need to manage multiple `.env` files across different environments.
3. **Audit**: The server logs whenever secrets are fetched for a project.

---

## ⏭️ Next Steps

- **[CLI Reference](cli-api.md)**: Explore all available flags and commands.
- **[Configuration](config.md)**: Learn how to customize your CLI settings.
