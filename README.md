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

## 📚 Documentation

For detailed information on how to use Bastion, please refer to the following guides:

- **[Getting Started](docs/getting-started.md)**: Prerequisites and initial project setup.
- **[Configuration](docs/config.md)**: Detailed explanation of environment variables and configuration settings.
- **[Initial Steps](docs/initial-steps.md)**: How to start the database, initialize the vault, and run the services.
- **[Local Development Workflow](docs/local-workflow.md)**: A practical guide to using the CLI for daily secret management.
- **[CLI Reference](docs/cli-api.md)**: Full command reference for the Bastion CLI.
- **[Installation Options](docs/cli-install.md)**: How to install Bastion on Linux, macOS, and Windows.

## 🤝 Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for detailed development instructions.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Built with ❤️ by [dcdavidev](https://github.com/dcdavidev)
