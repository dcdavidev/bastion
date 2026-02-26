# Initial Steps 🚀

Once you have cloned the repository and installed the dependencies, follow these steps to initialize your vault and start the Bastion services.

---

## 🐘 1. Prepare your Database

Bastion does not provide a built-in database container for development. You must have a running PostgreSQL 15+ instance (local or a managed service like AWS RDS, Supabase, or Azure Database for PostgreSQL).

1. **Create an empty database**: For example, `bastion`.
2. **Obtain your Connection String**: It should follow this format:
   `postgres://<user>:<password>@<host>:<port>/<dbname>?sslmode=<mode>`

Once you have these, you can proceed with the automated setup.

---

## 🔐 2. Initialize the Vault (Interactive Wizard)

Bastion provides a comprehensive setup wizard that handles database connection, migrations, and superuser creation.

### Step 2.1: Build the CLI
Ensure you have the latest version of the CLI built:
```bash
pnpm build
```

### Step 2.2: Run the Setup Wizard
Run the following command and follow the interactive prompts:
```bash
./bastion init
```

The wizard will ask for:
- **Admin Email & Password**: Used for your initial login.
- **Database Details**: Host, Port, User, Password, and Database Name.
- **Bastion URL**: The public address where your server will be reachable.
- **Password Store (Optional)**: If you have `pass` installed, you can choose to save your credentials in an encrypted, git-backed repository at `~/.config/bastion/store`.

**What happens during this step?**
1. **Migrations**: The database schema is automatically updated to the latest version.
2. **Master Key Generation**: A 256-bit random Master Key is generated.
3. **Key Wrapping**: The Master Key is encrypted with your Admin Password and stored in the database.
4. **Admin Creation**: An administrator account is created in the `users` table.
5. **Environment Variables**: Bastion will output the necessary environment variables (`BASTION_DATABASE_URL`, `BASTION_JWT_SECRET`, `BASTION_PORT`) for your server. Make sure to export them in your shell or use an `.env` file.

---

## 🏗️ 3. Run the Unified Service

Starting the backend server also starts the web dashboard, as they are now served by the same process:

```bash
# From the monorepo root
pnpm dev:server
```

### 📍 Access Points
| Service | URL | Description |
| :--- | :--- | :--- |
| **Unified Portal** | `http://localhost:8287` | API and Web Dashboard combined. |
| **Status Check** | `http://localhost:8287/api/v1/status` | Real-time health and migration check. |

---

## ✅ 4. Verification

Confirm the system state:

### Check Health Endpoint
```bash
curl http://localhost:8287/health
```

### Login via CLI
Try to authenticate using the email and password you set during `init`:
```bash
./bastion login --email your@email.com
```

---

## ⏭️ Next Steps
- **[Local Development Workflow](local-workflow.md)**: Learn how to manage client projects.
- **[CLI Reference](cli-api.md)**: Explore the full power of the Bastion CLI.
