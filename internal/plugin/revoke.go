package plugin

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/thevilledev/vault-plugin-secrets-vercel/internal/service"
)

var (
	errRemoteTokenRevokeFailed = errors.New("failed to revoke token")
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

	if req.Secret == nil {
		return nil, errInternalDataMissing
	}

	k, ok := req.Secret.InternalData[pathTokenID]
	if !ok {
		return nil, errInternalDataMissing
	}

	ks, ok := k.(string)
	if !ok || ks == "" {
		return nil, errInternalDataMissing
	}

	_, err = svc.DeleteAuthToken(ctx, ks)
	if err != nil {
		b.Logger().Error("failed to revoke/delete token from Vercel", "error", err)

		return nil, errRemoteTokenRevokeFailed
	}

	return &logical.Response{}, nil
}
