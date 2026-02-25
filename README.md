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

## 🚀 Quick Start

### Prerequisites

- **Go** 1.24+
- **Node.js** 24+ (with `pnpm`)
- **PostgreSQL**

### Installation

Bastion is managed as a **Turborepo Monorepo**.

#### 🐧 Linux

- **Snap (Recommended)**: `sudo snap install bastion`
- **AUR (Arch Linux)**: `yay -S bastion-bin`
- **AppImage**: Download the `.AppImage` from [Releases](https://github.com/dcdavidev/bastion/releases)
- **Native Packages**: `.deb`, `.rpm`, and `.apk` available in [Releases](https://github.com/dcdavidev/bastion/releases)

#### 🍏 macOS

- **Homebrew**:
  ```bash
  brew install dcdavidev/tap/bastion
  ```

#### 🪟 Windows

- **WinGet**: `winget install dcdavidev.bastion`
- **Scoop**:
  ```bash
  scoop bucket add dcdavidev https://github.com/dcdavidev/scoop-bucket
  scoop install bastion
  ```
- **Chocolatey**: `choco install bastion`

#### 📦 Package Managers & Containers

- **NPM**: `npm install -g @dcdavidev/bastion`
- **Go Install**: `go install github.com/dcdavidev/bastion/apps/cli@latest`
- **Docker**:
  ```bash
  docker pull ghcr.io/dcdavidev/bastion:latest
  # or from Docker Hub
  docker pull dcdavidev/bastion:latest
  ```

#### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/dcdavidev/bastion.git
   cd bastion
   ```
2. Install dependencies:
   ```bash
   pnpm install
   ```
3. Build everything (CLI, Server, Web):
   ```bash
   pnpm build
   ```
4. Development mode (starts DB, Server, and Web):
   ```bash
   pnpm dev
   ```
5. Initialize the Vault (First Run):
   ```bash
   ./bastion create-superuser
   ```
   _Follow the instructions to update your `.env` and initialize the database._

## 🚀 CI/CD & Releases

Bastion uses **GitHub Actions** and **GoReleaser** for automated multi-platform builds and deployments.

### Release Process

To trigger a new release (including Binaries and Snap package):

1. Tag your commit: `git tag -a v0.1.0 -m "Release v0.1.0"`
2. Push the tag: `git push origin v0.1.0`

The `pipeline` workflow will:

- Run build checks on every push/PR.
- Automatically publish to GitHub Releases and Snap Store on valid version tags.

_Note: Ensure `SNAPCRAFT_STORE_CREDENTIALS` is configured in your repository secrets for Snap publishing._

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
