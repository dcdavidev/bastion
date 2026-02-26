# Initial Steps 🚀

Once you have cloned the repository and installed the dependencies, follow these steps to initialize your vault and start the Bastion services.

---

## 🐘 1. Prepare your Database

Bastion requires a running PostgreSQL 15+ instance. You can use a local installation or a managed service (AWS RDS, Supabase, etc.).

1. **Create an empty database**: e.g., `bastion`.
2. **Obtain your Connection String**: It should follow this format:
   `postgres://<user>:<password>@<host>:<port>/<dbname>?sslmode=<mode>`

Once you have these, you can proceed with the automated setup.

---

## 🔐 2. Initialize the Vault (Interactive Wizard)

Bastion provides a comprehensive setup wizard that handles database connection, migrations, and superuser creation.

### Step 2.1: Build the CLI

Ensure you have the latest version of the CLI and Server built:

```bash
pnpm build
```

### Step 2.2: Run the Setup Wizard

Run the following command and follow the interactive prompts:

```bash
./bastion init
```

**What the wizard does:**

1. **Migrations**: Updates the database schema to the latest version.
2. **Master Key**: Generates a 256-bit Master Key, wraps it with your password, and stores it.
3. **Admin User**: Creates your initial administrator account.
4. **Environment**: Outputs the required variables (`BASTION_DATABASE_URL`, `BASTION_JWT_SECRET`, etc.).

_Tip: Save these variables in a `.env` file in the project root._

---

## 🏗️ 3. Run the Unified Service

Bastion serves both the API and the Web Dashboard from the same process. It relies on the frontend assets built in `apps/web/build/client`.

```bash
# Start the server (requires the frontend to be built first)
pnpm dev:server
```

### 📍 Access Points

| Service            | URL                                   | Description                     |
| :----------------- | :------------------------------------ | :------------------------------ |
| **Unified Portal** | `http://localhost:8287`               | Both the API and Web Dashboard. |
| **Status API**     | `http://localhost:8287/api/v1/status` | Real-time health check.         |

---

## ✅ 4. Verification

Confirm that everything is running correctly:

### Check Health Endpoint

```bash
curl http://localhost:8287/health
```

### Login via CLI

Authenticate using the credentials you set during the `init` wizard:

```bash
./bastion login --email your@email.com
```

---

## ⏭️ Next Steps

- **[Local Development Workflow](local-workflow.md)**: Start managing clients and projects.
- **[CLI Reference](cli-api.md)**: Explore the full power of the Bastion CLI.
