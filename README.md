# vault-plugin-secrets-vercel

Vault Secrets Plugin for Vercel allows you to dynamically generate Vercel API tokens through Vault.

It is useful for more advanced CI/CD use cases where the common
[Vercel git integration](https://vercel.com/docs/concepts/deployments/git/vercel-for-github) is not being utilised.

With this plugin, the CI/CD pipeline should:

- Authenticate to Vault through a number of means. See [hashicorp/vault-action docs](https://github.com/hashicorp/vault-action#authentication-methods) for full list of available methods, such as:
    - AppRole
    - JWT OIDC
    - A pre-defined token
- Call the plugin to generate a short-lived Vercel token. TTL for the toke is user-configurable.
- Run the actual deployment pipeline, such as [Github Actions for Vercel](https://vercel.com/guides/how-can-i-use-github-actions-with-vercel)
- After token lifetime runs out, Vault revokes the token automatically.

Here's a full example of a Github Actions pipeline utilising this plugin:

```
name: Vercel Preview Deployment
env:
  VERCEL_ORG_ID: ${{ secrets.VERCEL_ORG_ID }}
  VERCEL_PROJECT_ID: ${{ secrets.VERCEL_PROJECT_ID }}
on:
  push:
    branches-ignore:
      - main

jobs:
  Deploy-Preview:
    runs-on: ubuntu-latest
    steps:
      - name: Import Secrets
        id: secrets
        uses: hashicorp/vault-action@65d7a12a8098b0aa7fcfdf22ad850c051f8b3ccb # v2.7.2
        with:
          url: ${{ secrets.VAULT_ADDR }}
          method: approle
          roleId: ${{ secrets.VAULT_ROLE_ID }}
          secretId: ${{ secrets.VAULT_SECRET_ID }}
          secrets: |
            vercel-secrets/token bearer_token | VERCEL_TOKEN

      - uses: actions/checkout@c85c95e3d7251135ab7dc9ce3241c5835cc595a9 # v3.5.3

      - name: Install Vercel CLI
        run: npm install --global vercel@latest

      - name: Pull Vercel Environment Information
        run: vercel pull --yes --environment=preview --token=${{ steps.secrets.outputs.VERCEL_TOKEN }}

      - name: Build Project Artifacts
        run: vercel build --token=${{ steps.secrets.outputs.VERCEL_TOKEN }}

      - name: Deploy Project Artifacts to Vercel
        run: vercel deploy --prebuilt --token=${{ steps.secrets.outputs.VERCEL_TOKEN }}
```
## Scope

Currently this project is scoped for personal (or "Hobby") Vercel accounts.

Any Pro/Enterprise plan features, such as these, are scoped out:

- Team-specific tokens
- Token-specific scope or permissions

Reason is, I don't have a Pro/Enterprise plan at hand. Contributions are welcome, of course!

## Getting started

Get started by following the documentation:

- [Running the plugin locally](./docs/development.md)
- [Configuring and generating tokens with the plugin](./docs/configuration.md)
- [Installing the plugin to an existing Vault installation](./docs/install.md)
