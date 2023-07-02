# vault-plugin-secrets-vercel

Vault Secrets Plugin for Vercel allows you to dynamically generate Vercel API tokens through Vault.

Currently supports personal Vercel accounts. Additional features, such as token-specific fine-grained permissions
are not supported by the plugin. I do not have a Pro/Enterprise plan to develop this further, but contributions welcome.

## Getting started

Run a local development setup with the following commands:

```
$ make build
$ make start
$ make enable
```

Go to the [Vercel tokens page](https://vercel.com/account/tokens) and generate an admin token. Then configure the plugin:

```
$ vault write vercel-secrets/config api_key=<your-api-key-here>
```

Generate a new Vault plugin managed token:

```
$ vault read vercel-secrets/token
Key                Value
---                -----
lease_id           vercel-secrets/token/GtbmIK80YfqX3hOwn1A23Lro
lease_duration     10s
lease_renewable    false
bearer_token       xyzabbacdc
token_id           bababababa
```

Vault will automatically revoke the API key after the lease duration.

## Running it on production

Please don't - yet!

If you do, please refer to the deployment model for plugins [from the official HashiCorp Vault tutorials](https://developer.hashicorp.com/vault/tutorials/app-integration/plugin-backends#setup-vault). 