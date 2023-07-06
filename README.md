# vault-plugin-secrets-vercel

Vault Secrets Plugin for Vercel allows you to dynamically generate Vercel API tokens through Vault.
Useful for CI/CD as you can generate short-lived deployment tokens and let them revoke once you are done.

Currently supports personal Vercel accounts. Additional features, such as token-specific fine-grained permissions
are not supported by the plugin - simply because I have no Pro/Enterprise plan to develop it against with.

## Run it locally

Run a local development setup with the following commands:

```
$ make build
$ make start
$ make enable
```

This sets up a local Vault installation in development mode and enables the plugin.

Go to the [Vercel tokens page](https://vercel.com/account/tokens) and generate an admin token. Then configure the plugin:

```
$ vault write vercel-secrets/config api_key=<your-api-key-here>
```

You can also define a maximum TTL for the secrets by defining an additional parameter `max_ttl=<seconds>`. By default it is 10 minutes. TTLs can be defined on a per-token basis, but they will need to be lower than or equal to the max.

## Generate tokens

Generate a new Vault plugin managed Vercel token:

```
$ vault read vercel-secrets/token
Key                Value
---                -----
lease_id           vercel-secrets/token/<lease-id>
lease_duration     10m
lease_renewable    false
bearer_token       xyzabbacdc
token_id           bababababa
```

You can set a custom lease duration with the parameter `ttl=<seconds>`.

Vault will automatically revoke & delete the API key after the lease duration.
The generated token also has an expiration time equal to the lease duration on Vercel side.

## Running it on production

Please don't - yet!

If you do, please refer to the documentation on deployment model for plugins [from the official HashiCorp Vault tutorials](https://developer.hashicorp.com/vault/tutorials/app-integration/plugin-backends#setup-vault).