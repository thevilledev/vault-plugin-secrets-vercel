# Development setup

**NOTE: This setup is purely for local development, demonstration or testing.**

## Run it locally

The project is built and tested with Go 1.25.7. The CI lint job uses golangci-lint v2.11, and release configuration is validated with GoReleaser v2.

Run a local development setup with the following commands:

```
$ make build
$ make start
$ make enable
```

This creates a local Vault installation in development mode and enables the plugin.

## Configuring the plugin

Follow [the configuration guide](configuration.md).

## Mock usage

By setting the API key to `mock` the plugin is forced to use the mock API client, which does not communicate
with the Vercel API at all. Useful for development purposes and refactoring. The returned `bearer_token` is hard coded to `some-bearer-token`.

```
$ vault write vercel-secrets/config api_key=mock
$ vault read vercel-secrets/token
Key                Value
---                -----
lease_id           vercel-secrets/token/BIxRweNgXNSQsnbeBBmiea8X
lease_duration     10m
lease_renewable    false
bearer_token       some-bearer-token
team_id            n/a
token_id           vault-plugin-secrets-vercel-1689595722412039000-1689595722412067000
```

## Live integration tests

Live tests create and delete real Vercel tokens. They are not part of default CI.

```
$ ACC_TEST=yes VERCEL_TOKEN=<vercel-token> go test -race -count=1 -v ./internal/service -run TestIntegration
```

To include the team-scoped integration test, also set `VERCEL_TEAM_ID=<vercel-team-id>`.

The same live tests can be run manually from the `live integration` GitHub Actions workflow. It requires a `VERCEL_TOKEN` repository secret; `VERCEL_TEAM_ID` is optional.
