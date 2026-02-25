<p align="center">
  <img src="https://github.com/dcdavidev/bastion/blob/main/bastion-banner.png?raw=true" alt="bastion banner">
</p>

# Bastion

Bastion is a single-user controlled, multi-tenant E2EE secrets vault built with Go and React. It provides a secure, self-hosted fortress to manage multiple client secrets via a powerful CLI and a modern dashboard, ensuring data stays private with a robust **blind-backend architecture**.

## 🛡️ Security Model

- **End-to-End Encryption (E2EE):** All secrets are encrypted client-side (CLI or Dashboard) before reaching the server.
- **Key Wrapping:** Uses a multi-layered key hierarchy (Master Key -> Project Data Key -> Secret).
- **Blind Backend:** The server never processes or stores plaintext secrets or raw keys.
- **Multi-User Access:** Admin can delegate project-specific access to Collaborators using secure re-wrapping techniques.
- **Audit Logging:** Every sensitive operation is cryptographically linked and logged.

## 🚀 Getting Started (Step-by-Step)

Follow these steps to set up and start using Bastion locally.

### 1. Prerequisites

Ensure you have the following installed:

- **Go** (1.24+)
- **Node.js** (24+) and **pnpm**
- **Docker** (to run the database)

### 2. Initial Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/dcdavidev/bastion.git
   cd bastion
   ```

2. **Install dependencies:**

   ```bash
   pnpm install
   ```

3. **Configure Environment:**
   Create your local `.env` file from the example:
   ```bash
   cp .env.example .env
   ```

### 3. Start the Database

Bastion uses PostgreSQL. You can start a local instance using Docker:

```bash
pnpm docker:db:up
```

### 4. Initialize the Vault (First Time Only)

Before running the server, you must create a Superuser and initialize the encryption keys:

1. **Build the CLI:**

   ```bash
   pnpm build
   ```

2. **Run Initialization:**

   ```bash
   ./bastion create-superuser
   ```

   Follow the prompts to set your **Username** and **Master Password**.

   > Save your Master Password securely. It is the only way to decrypt your secrets.

3. **Verify `.env`:**
   The command will generate a `VAULT_ID`. Ensure it is correctly reflected in your `.env` file if prompted.

### 5. Run the Services (Server & Dashboard)

You can start both the backend server and the web dashboard in development mode with a single command:

```bash
pnpm dev
```

- **Backend Server:** Runs on `http://localhost:8080`
- **Web Dashboard:** Runs on `http://localhost:5173`

### 6. Using the CLI

Open a new terminal to interact with Bastion via the CLI.

1. **Login:**

   ```bash
   ./bastion login
   ```

   Enter the server URL (`http://localhost:8080`) and your superuser credentials.

2. **Organize Secrets:**
   Secrets are organized by **Clients** and **Projects**.

   ```bash
   # Create a Client
   ./bastion create-client --name "MyClient"

   # Create a Project for that Client
   ./bastion create-project --name "ProductionApp" --client-id <CLIENT_ID>
   ```

3. **Store a Secret:**

   ```bash
   ./bastion set -p DATABASE_URL -v "postgres://user:pass@host:5432/db" < PROJECT_ID > -k
   ```

4. **Inject Secrets (Runtime):**
   Execute any command with your project's secrets injected as environment variables:
   ```bash
   ./bastion run -p npm run start < PROJECT_ID > --
   ```

## 📦 Installation Options

Bastion can be installed via several package managers for production use.

### 🐧 Linux

- **Snap**: `sudo snap install bastion-cli`
- **Native Packages**: `.deb`, `.rpm`, and `.apk` available in [Releases](https://github.com/dcdavidev/bastion/releases)

### 🍏 macOS

- **Homebrew**:
  ```bash
  brew install dcdavidev/tap/bastion-cli
  ```

### 🪟 Windows

- **Scoop**:
  ```bash
  scoop bucket add dcdavidev https://github.com/dcdavidev/scoop-bucket
  scoop install bastion-cli
  ```

### 📦 NPM & Docker

- **NPM**: `npm install -g @dcdavidev/bastion-cli`
- **Docker**: `docker pull ghcr.io/dcdavidev/bastion-server:latest`

## 🚀 CI/CD & Releases

Bastion uses **GitHub Actions** and **GoReleaser** for automated multi-platform builds. To trigger a new release:

1. Tag your commit: `git tag -a v0.1.0 -m "Release v0.1.0"`
2. Push the tag: `git push origin v0.1.0`

## 🤝 Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for detailed development instructions.

---

Built with ❤️ by [dcdavidev](https://github.com/dcdavidev)
