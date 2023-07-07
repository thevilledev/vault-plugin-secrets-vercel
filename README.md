# vault-plugin-secrets-vercel

Vault Secrets Plugin for Vercel allows you to dynamically generate Vercel API tokens through Vault.
Useful for CI/CD as you can generate short-lived deployment tokens through Vault and manage developer access that way.
Vault revokes the tokens automatically upon expiration, which is user-configurable.

## Scope

Currently this project is scoped for personal (or "Hobby") Vercel accounts.

Any Pro/Enterprise plan features, such as these, are scoped out:

- Team-specific tokens
- Token-specific scope or permissions

Access to Pro/Enterprise plan or contributions are welcome, of course!

## Getting started

Get started by following the documentation:

- [Running the plugin locally](./docs/development.md)
- [Configuring and generating tokens with the plugin](./docs/configuration.md)
- [Installing the plugin to an existing Vault installation](./docs/install.md)