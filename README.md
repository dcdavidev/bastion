# Bastion 🏰

Bastion is a single-user, open-source E2EE secrets vault built with Go. It provides a secure, self-hosted fortress to manage multiple client secrets via a powerful CLI and dashboard, ensuring data stays private with blind-backend architecture.

## 🚀 Quick Start

### Prerequisites
- **Go** 1.24+
- **Node.js** 24+
- **PostgreSQL**
- **Docker** (optional, for DB)

### Setup
1. Clone the repository:
   ```bash
   git clone https://github.com/dcdavidev/bastion.git
   cd bastion
   ```
2. Install backend dependencies:
   ```bash
   go mod download
   ```
3. Install frontend dependencies:
   ```bash
   cd frontend
   npm install
   ```

## 🛡 Tech Stack & Architecture
- **Backend:** Go (Golang) 1.24+ with [chi](https://github.com/go-chi/chi) router.
- **Frontend:** React + TypeScript (Vite).
- **Database:** PostgreSQL with [pgx](https://github.com/jackc/pgx) driver.
- **Security:**
  - **Argon2id:** For secure admin password hashing and key derivation.
  - **AES-256-GCM:** For authenticated encryption (Key Wrapping).
  - **Blind Backend:** Secrets are never processed in plaintext by the server.

## 🏗 Project Structure
- `cmd/server/`: Backend entrypoint and server configuration.
- `internal/api/`: REST API handlers (incoming).
- `internal/auth/`: Admin authentication and authorization.
- `internal/crypto/`: Cryptographic primitives (encryption/decryption).
- `internal/db/`: Database connection and pooling.
- `internal/models/`: Data structures and domain models.
- `frontend/`: React application.

## 🤝 Contributing
Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on our development process.
