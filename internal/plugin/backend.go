package plugin

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// #nosec G101
const (
	backendSecretType = "vercel_token"
	vercelHelp        = `
Vercel Secrets backend is a secrets backend for dynamically managing Vercel tokens.
`
)

type backend struct {
	*framework.Backend
}

var _ logical.Factory = Factory

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := newBackend()

	if conf == nil {
		return nil, fmt.Errorf("configuration passed into backend is nil")
	}

	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	return b, nil
}

func newBackend() *backend {
	b := &backend{}

	b.Backend = &framework.Backend{
		Help:        strings.TrimSpace(vercelHelp),
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
						Description: "Vercel API token ID.",
					},
					pathTokenBearerToken: {
						Type:        framework.TypeString,
						Description: "Vercel API token.",
					},
				},
				Revoke: b.Revoke,
			},
		},
	}

	return b
}