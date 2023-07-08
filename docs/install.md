# Installing vault-plugin-secrets-vercel

**NOTE: This project is under active development and has not been production battle-tested yet - treat it as such and contributions welcome.**

Install the plugin with the following instructions to an existing Vault installation. For a more detailed explanation, please refer to [the official HashiCorp Vault documentation on plugin backends](https://developer.hashicorp.com/vault/tutorials/app-integration/plugin-backends#setup-vault).

For local development or testing, please refer to [the development setup guide](development.md).

## Downloading the plugin

Grab the latest release from [the releases page](https://github.com/thevilledev/vault-plugin-secrets-vercel/releases) that matches
your operating system and architecture.

### Verifying the plugin build

All builds have a matching SHA256 checksum. These are signed and the matching public key can be found from [signing-key-public.asc](https://raw.githubusercontent.com/thevilledev/vault-plugin-secrets-vercel/main/signing-key-public.asc).

You can validate these builds by doing the following:

- Download the release files `vault-plugin-secrets-vercel_<version>_SHA256SUMS` and `vault-plugin-secrets-vercel_<version>_SHA256SUMS.sig` to the same directory.
- Download the signing key [signing-key-public.asc](https://raw.githubusercontent.com/thevilledev/vault-plugin-secrets-vercel/main/signing-key-public.asc).
- Import the signing key with `gpg --import signing-key-public.asc`
- Verify checksums with `gpg --verify vault-plugin-secrets-vercel_<version>_SHA256SUMS.sig`

You have now verified that the checksums are signed by the project keys.

Next, generate a SHA256 checksum for the build you downloaded, for example:

```
$ sha256sum vault-plugin-secrets-vercel_v0.2.3_Darwin_arm64.tar.gz
cdadd3ec3115995625541dd4222ae4986bfdc2531b2f442ead3f8d385f87df88  vault-plugin-secrets-vercel_v0.2.3_Darwin_arm64.tar.gz
```

And finally, check that it matches the SHA256SUMS file:

```
$ grep 'vault-plugin-secrets-vercel_v0.2.3_Darwin_arm64.tar.gz' vault-plugin-secrets-vercel_v0.2.3_SHA256SUMS
cdadd3ec3115995625541dd4222ae4986bfdc2531b2f442ead3f8d385f87df88  vault-plugin-secrets-vercel_v0.2.3_Darwin_arm64.tar.gz
```

All good!

## Registering the plugin

Un-tar the installation package and move the plugin binary to the `plugin_directory` of your Vault installation. For example:

```
$ mv vault-plugin-secrets-vercel /opt/vault/plugins/
$ chown vault:vault /opt/vault/plugins/vault-plugin-secrets-vercel
```

Calculate the SHA256 checksum of the plugin binary.

```
$ SHA256=$(sha256sum /opt/vault/plugins/vault-plugin-secrets-vercel | cut -d ' ' -f1)
```

And register the plugin into the Vault plugin catalog. This assumes you have Vault CLI access properly configured.

```
$ vault plugin register -sha256=$SHA256 secret vault-plugin-secrets-vercel
Success! Registered plugin: vault-plugin-secrets-vercel
```

## Enabling the plugin

To enable the plugin, run:

```
$ vault secrets enable -path=vercel-secrets vault-plugin-secrets-vercel
```

Validate that it works:

```
$ vault read vercel-secrets/info
Key                    Value
---                    -----
build_commit           78f02477a582126fb68e195b0c1de3df63f335bf
build_commit_branch    HEAD
build_commit_date      2023-07-07T19:10:53Z
build_date             2023-07-07T19:13:32Z
build_dirty            false
build_tag              v0.2.4
build_version          0.2.4
```

## Configuring the plugin

Follow [the configuration guide](configuration.md).