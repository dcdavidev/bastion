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
3. VS Code users: Accept the recommended extensions when opening the project.

## 📜 Coding Standards
- **Go**: Follow standard Go Layout and idiomatic patterns.
- **TypeScript**: Use ES Modules, strict typing (no `any`).
- **Commits**: Use [Conventional Commits](https://www.conventionalcommits.org/).

## 🧪 Testing
- Backend: `go test ./...`
- Frontend: `npm test` (when implemented)
