# Installation Options 📦

Bastion can be installed via several package managers or compiled directly from the source.

## 🛠️ Build from Source (Manual)

If you have **Go 1.24+** installed, you can compile the CLI directly on your machine. This is the recommended way for development or for architectures not yet covered by binary releases.

```bash
# Clone the repository
git clone https://github.com/dcdavidev/bastion.git
cd bastion

# Build the CLI binary
go build -o bastion ./apps/cli/main.go

# (Optional) Move to your PATH
sudo mv bastion /usr/local/bin/
```

---

## 🐧 Linux

- **Snap Store**:
  ```bash
  sudo snap install bastion-cli
  ```
- **Native Packages**: `.deb`, `.rpm`, and `.apk` files are available for every release in the [GitHub Releases](https://github.com/dcdavidev/bastion/releases) page.

## 🍏 macOS

- **Homebrew**:
  ```bash
  brew tap dcdavidev/tap
  brew install bastion-cli
  ```

## 🪟 Windows

- **Scoop**:
  ```bash
  scoop bucket add dcdavidev https://github.com/dcdavidev/scoop-bucket
  scoop install bastion-cli
  ```

---

## 📦 Alternative Distribution

### NPM (Node Package Manager)

The CLI is also distributed as a global npm package for environments where Node.js is the primary runtime.

```bash
npm install -g @dcdavidev/bastion-cli
```

### Docker (Unified Server + UI)

For the server-side deployment, use the unified Docker image which includes both the Go backend and the React dashboard.

```bash
# Pull the latest unified image
docker pull ghcr.io/dcdavidev/bastion:latest

# Run with required environment variables
docker run -p 8287:8287 \
  -e BASTION_DATABASE_URL="postgres://..." \
  -e BASTION_JWT_SECRET="your-secret" \
  ghcr.io/dcdavidev/bastion:latest
```

---

## ✅ Post-Installation

After installing, verify the installation by checking the version:

```bash
bastion version
```

Next, follow the **[Getting Started](getting-started.md)** guide to initialize your vault.
