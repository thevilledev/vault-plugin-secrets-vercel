package plugin

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// #nosec G101
	backendSecretType = "vercel_token"
	backendPathHelp   = `
Vercel Secrets backend is a secrets backend for dynamically managing Vercel tokens.`
	// #nosec G101
	secretTokenIDDescription = `
Token ID of the generated API key is stored in the plugin backend.
This ID is used for revocation purposes. It can only be used to identify a key,
and cannot be used to do API operations.`
)

var (
	errBackendEmptyConfig = errors.New("configuration passed into backend is nil")
)

type backend struct {
	*framework.Backend
}

var _ logical.Factory = Factory

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := newBackend()

	if conf == nil {
		return nil, errBackendEmptyConfig
	}

	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	return b, nil
}

func newBackend() *backend {
	b := &backend{}

	b.Backend = &framework.Backend{
		Help:        backendPathHelp,
		BackendType: logical.TypeLogical,
		Paths: framework.PathAppend(
			b.pathConfig(),
			b.pathToken(),
			b.pathInfo(),
		),
		Secrets: []*framework.Secret{
			{
				Type: backendSecretType,
				Fields: map[string]*framework.FieldSchema{
					pathTokenID: {
						Type:        framework.TypeString,
						Description: secretTokenIDDescription,
					},
				},
				Revoke: b.Revoke,
			},
		},
	}

	return b
}
