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

## 🛠 Tech Stack
- **Backend:** Go (Golang)
- **Frontend:** React + TypeScript (Vite)
- **Database:** PostgreSQL
- **Security:** Argon2id, AES-256-GCM

## 🤝 Contributing
Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on our development process.
