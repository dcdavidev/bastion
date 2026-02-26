# Contributing to Bastion 🛡️

Thank you for your interest in contributing to Bastion!

## 🛠 Development Tools

To ensure a smooth development experience, please make sure you have the following tools installed:

### Required Tools

- **Go 1.24+**: The core backend language.
- **Node.js 24+**: For the React frontend.
- **pnpm** (preferred) or **npm**: Package manager for the frontend.
- **PostgreSQL 15+**: A running instance (local or managed) with an empty database created.
- **VS Code**: Recommended editor with extensions.

### Backend Tech Details

- **Router**: [chi v5](https://github.com/go-chi/chi)
- **Configuration**: Environment variable based (`BASTION_HOST`, `BASTION_PORT`, etc.)
- **Database Driver**: [pgx v5](https://github.com/jackc/pgx)

### Setup Environment

1. Follow the installation steps in the [README.md](README.md).
2. **Database Setup**: Ensure you have a running PostgreSQL instance. Create a new, empty database (e.g., `bastion`).
3. **Environment Variables**: Bastion uses environment variables for configuration.
   Example environment setup:
   ```bash
   export BASTION_DATABASE_URL='postgres://user:password@localhost:5432/bastion?sslmode=disable'
   export BASTION_JWT_SECRET='your-secure-jwt-secret'
   export BASTION_PORT='8287'
   ```
4. **Initialize Schema**: Apply the migrations to your database:
   ```bash
   # Using psql (replace with your connection string if necessary)
   psql $BASTION_DATABASE_URL -f packages/core/db/migrations/000001_initial_schema.up.sql
   ```
   *Note: Alternatively, the `bastion init` command can also run migrations for you.*
4. Generate required security keys (example using openssl):

   ```bash
   # Generate a 32-byte Master Key or Salt (hex)
   openssl rand -hex 32
   ```

   _Note: For the admin password hash, you should use a small Go utility with `packages/core/crypto` to ensure compatibility with the Argon2id parameters._

5. VS Code users: Accept the recommended extensions when opening the project.

## 🛡️ Bastion CLI

The CLI is located in `apps/cli/`. You can build it using:

```bash
go build -o bastion ./apps/cli/main.go
```

### Basic Commands

- **Login**: `bastion login --url http://localhost:8287` (Admin) or `bastion login -n <username>` (Collaborator).
  Authenticates and stores the JWT locally in `~/.bastion/token`.
- **Create Project**: `bastion create project -n MyProject -c <CLIENT_UUID>`
  Generates a new data key and creates a project for a client.
- **List Clients**: `bastion list clients`
  Displays a table of all clients in the vault.
- **List Projects**: `bastion list projects -c <CLIENT_UUID>`
  Displays all projects for a specific client.
- **Run**: `bastion run -p <PROJECT_ID> -- <command>`
  Injects secrets into the environment of the specified command.
- **Set**: `bastion set -p <PROJECT_UUID> -k MY_KEY -v my_value`
  Encrypts and stores a secret in the specified project.

## 📜 Coding Standards

- **Go**: Follow standard Go Layout and idiomatic patterns.
- **TypeScript**: Use ES Modules, strict typing (no `any`).
- **Commits**: Use [Conventional Commits](https://www.conventionalcommits.org/).

## 🧪 Testing

- Backend: `go test ./...`
- Frontend: `npm test` (when implemented)
