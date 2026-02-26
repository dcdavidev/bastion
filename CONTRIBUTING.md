# Contributing to Bastion 🛡️

Thank you for your interest in contributing to Bastion!

## 🛠 Development Tools

To ensure a smooth development experience, please make sure you have the following tools installed:

### Required Tools

- **Go 1.24+**: The core backend language.
- **Node.js 24+**: For the React frontend.
- **pnpm** (preferred): Package manager for the monorepo.
- **PostgreSQL 15+**: A running instance (local or managed) with an empty database created.
- **VS Code**: Recommended editor with extensions.

### Backend Tech Details

- **Router**: [chi v5](https://github.com/go-chi/chi)
- **Configuration**: Environment variable based (`BASTION_HOST`, `BASTION_PORT`, etc.)
- **Database Driver**: [pgx v5](https://github.com/jackc/pgx)

### Setup Environment

1. Follow the installation steps in the [README.md](README.md).
2. **Interactive Setup**: The easiest way to get started is by building the CLI and running the wizard:
   ```bash
   pnpm build
   ./bastion init
   ```
3. **Environment Variables**: Alternatively, set them manually in a `.env` file:
   ```bash
   BASTION_DATABASE_URL='postgres://user:password@localhost:5432/bastion?sslmode=disable'
   BASTION_JWT_SECRET='your-secure-jwt-secret'
   BASTION_PORT='8287'
   ```
4. **Run Services**:
   ```bash
   # Start the unified server (Backend + Frontend)
   pnpm dev:server
   ```

## 🛡️ Bastion CLI

The CLI is located in `apps/cli/`. You can build it using:

```bash
go build -o bastion ./apps/cli/main.go
```

### Basic Commands

- **Login**: `bastion login --email your@email.com`
  Authenticates and stores the JWT locally in `~/.bastion/token`.
- **Create Project**: `bastion create project`
  Generates a new data key and creates a project for a client.
- **List Clients**: `bastion list clients`
- **Run**: `bastion run -p <PROJECT_ID> -- <command>`
  Injects secrets into the environment of the specified command.
- **Set**: `bastion set <KEY> <VALUE> -p <PROJECT_UUID>`
  Encrypts and stores a secret in the specified project.

## 📜 Coding Standards

- **Go**: Follow standard Go Layout and idiomatic patterns.
- **TypeScript**: Use ES Modules, strict typing (no `any`).
- **Commits**: Use [Conventional Commits](https://www.conventionalcommits.org/).

## 🧪 Testing

- **Full Workspace**: `pnpm test`
- **Backend**: `go test ./...`
- **Frontend**: `pnpm run typecheck`
