package plugin

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/thevilledev/vault-plugin-secrets-vercel/internal/service"
)

var (
	errRemoteTokenRevokeFailed = errors.New("failed to revoke token from Vercel")
	errInternalDataMissing     = errors.New("missing internal data from secret")
)

func (b *backend) Revoke(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		return nil, errBackendNotConfigured
	}

	svc := service.NewWithBaseURL(cfg.APIKey, cfg.BaseURL)

	k, ok := req.Secret.InternalData[pathTokenID]
	if !ok {
		return nil, errInternalDataMissing
	}

	ks, _ := k.(string)

	_, err = svc.DeleteAuthToken(ctx, ks)
	if err != nil {
		b.Logger().Error("failed to revoke/delete the token from Vercel %s", err)

		return nil, errRemoteTokenRevokeFailed
	}

	return &logical.Response{}, nil
}
