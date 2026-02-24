# Contributing to Bastion 🛡️

Thank you for your interest in contributing to Bastion!

## 🛠 Development Tools

To ensure a smooth development experience, please make sure you have the following tools installed:

### Required Tools
- **Go 1.24+**: The core backend language.
- **Node.js 24+**: For the React frontend.
- **pnpm** (preferred) or **npm**: Package manager for the frontend.
- **Docker**: Used to run PostgreSQL locally for development.
- **VS Code**: Recommended editor with extensions.

### Backend Tech Details
- **Router**: [chi v5](https://github.com/go-chi/chi)
- **Environment Management**: [godotenv](https://github.com/joho/godotenv)
- **Database Driver**: [pgx v5](https://github.com/jackc/pgx)

### Setup Environment
1. Follow the installation steps in the [README.md](README.md).
2. Copy `.env.example` to `.env` and fill in the required variables:
   ```bash
   cp .env.example .env
   ```
3. Initialize the database schema:
   ```bash
   # You can run this directly if you have psql installed
   psql $DATABASE_URL -f internal/db/migrations/000001_initial_schema.up.sql
   ```
4. Generate required security keys (example using openssl):
   ```bash
   # Generate a 32-byte Master Key or Salt (hex)
   openssl rand -hex 32
   ```
   *Note: For the admin password hash, you should use a small Go utility with `internal/crypto` to ensure compatibility with the Argon2id parameters.*

5. VS Code users: Accept the recommended extensions when opening the project.

## 🛡️ Bastion CLI
The CLI is located in `cmd/bastion/`. You can build it using:
```bash
go build -o bastion ./cmd/bastion
```

### Basic Commands
- **Login**: `bastion login --url http://localhost:8080`
  Authenticates and stores the JWT locally in `~/.bastion/token`.
- **Run**: `bastion run -p <PROJECT_UUID> -- <command>`
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
