# @dcdavidev/bastion-web

Static assets for the **Bastion E2EE Secrets Vault** Dashboard.

This package contains the compiled frontend of Bastion, a secure self-hosted end-to-end encrypted vault. It is designed to be served by the [Bastion Server](https://www.npmjs.com/package/@dcdavidev/bastion-server) or any static file server.

## Features

- **E2EE Dashboard**: Encrypt and decrypt secrets directly in your browser.
- **Blind Backend Integration**: Zero-knowledge architecture.
- **Multi-Tenant UI**: Manage clients, projects, and collaborators.
- **Built with React Router 7 & Pittorica UI**.

## Installation

This package is intended to be used as a dependency for the Bastion server or for manual deployments:

```bash
npm install @dcdavidev/bastion-web
```

The assets will be available in `node_modules/@dcdavidev/bastion-web/build/client`.

## Manual Usage

If you want to serve the dashboard using Nginx or another web server, copy the contents of the `build/client` directory to your web root and ensure you handle SPA routing (redirecting all non-file requests to `index.html`).

## Links

- **Main Repository**: [github.com/dcdavidev/bastion](https://github.com/dcdavidev/bastion)
- **CLI Package**: [@dcdavidev/bastion-cli](https://www.npmjs.com/package/@dcdavidev/bastion-cli)

## License

MIT © [Davide Di Criscito](https://github.com/dcdavidev)
