package plugin

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// backend wraps the backend framework and adds a map for storing key value pairs
type backend struct {
	*framework.Backend
}

var _ logical.Factory = Factory

// Factory configures and returns Mock backends
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
		),
	}

	return b
}

const vercelHelp = `
Vercel Secrets backend is a secrets backend for dynamically managing Vercel tokens.
`
