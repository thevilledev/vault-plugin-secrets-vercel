package plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/thevilledev/vault-plugin-secrets-vercel/pkg/service"
)

func (b *backend) Revoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("backend is missing API key")
	}

	svc := service.New(cfg.APIKey)

	k, ok := req.Secret.InternalData[pathTokenID]
	if !ok {
		return nil, fmt.Errorf("API key is missing from the secret")
	}
	ks := k.(string)

	err = svc.DeleteAuthToken(ctx, ks)
	if err != nil {
		return nil, fmt.Errorf("failed to delete token: %s", err)
	}

	return &logical.Response{}, nil
}
