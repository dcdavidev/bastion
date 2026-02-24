# Bastion 🏰

Bastion is a single-user controlled, multi-tenant E2EE secrets vault built with Go and React. It provides a secure, self-hosted fortress to manage multiple client secrets via a powerful CLI and a modern dashboard, ensuring data stays private with a robust **blind-backend architecture**.

## 🛡️ Security Model
- **End-to-End Encryption (E2EE):** All secrets are encrypted client-side (CLI or Dashboard) before reaching the server.
- **Key Wrapping:** Uses a multi-layered key hierarchy (Master Key -> Project Data Key -> Secret).
- **Blind Backend:** The server never processes or stores plaintext secrets or raw keys.
- **Multi-User Access:** Admin can delegate project-specific access to Collaborators using secure re-wrapping techniques.
- **Audit Logging:** Every sensitive operation is cryptographically linked and logged.

## 🚀 Quick Start

### Prerequisites
- **Go** 1.24+
- **Node.js** 24+ (with `pnpm`)
- **PostgreSQL**

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/dcdavidev/bastion.git
   cd bastion
   ```
2. Build the CLI:
   ```bash
   go build -o bastion ./cmd/bastion
   ```
3. Initialize the Vault (First Run):
   ```bash
   ./bastion create-superuser
   ```
   *Follow the instructions to update your `.env` and initialize the database.*

## 🛠️ CLI Usage
Bastion comes with a powerful CLI for secret injection and management.

- `bastion login`: Authenticate with the server.
- `bastion run -p <PROJECT_ID> -- <command>`: Inject secrets into a process (dotenvx style).
- `bastion set -p <PROJECT_ID> -k KEY -v value`: Encrypt and store a secret.
- `bastion create-collaborator`: Grant restricted access to a team member.

## 🌐 Dashboard
The web interface (built with React + Pittorica) allows for easy management of:
- **Clients & Projects:** Organize secrets by entity.
- **Audit Logs:** Monitor access and changes in real-time.
- **Collaborators:** Manage user permissions and project access.

## 🤝 Contributing
Please see [CONTRIBUTING.md](CONTRIBUTING.md) for detailed development instructions.

---
Built with ❤️ by [dcdavidev](https://github.com/dcdavidev)
