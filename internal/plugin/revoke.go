package plugin

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/thevilledev/vault-plugin-secrets-vercel/internal/service"
)

var (
	errRemoteTokenRevokeFailed = errors.New("failed to revoke token from Vercel")
)

func (b *backend) Revoke(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if cfg.APIKey == "" {
		return nil, errBackendNotConfigured
	}

	svc := service.NewWithBaseURL(cfg.APIKey, cfg.BaseURL)

	k, ok := req.Secret.InternalData[pathTokenID]
	if !ok {
		return nil, fmt.Errorf("token ID is missing from the secret")
	}

	ks, ok := k.(string)
	if !ok {
		b.Logger().Trace("type assertion failed: %+v", ks)

		return nil, errTypeAssertionFailed
	}

	_, err = svc.DeleteAuthToken(ctx, ks)
	if err != nil {
		b.Logger().Error("failed to revoke/delete the token from Vercel %s", err)

		return nil, errRemoteTokenRevokeFailed
	}

	return &logical.Response{}, nil
}
