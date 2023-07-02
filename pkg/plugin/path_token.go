package plugin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/thevilledev/vault-plugin-secrets-vercel/pkg/service"
)

var (
	pathPatternToken        = "token"
	errBackendNotConfigured = errors.New("backend not configured")
)

func (b *backend) pathToken() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: pathPatternToken,

			Fields: map[string]*framework.FieldSchema{
				"role": {
					Type:        framework.TypeString,
					Description: "Name of the role to apply to the API key",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathTokenWrite,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathTokenWrite,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathTokenWrite,
				},
			},
		},
	}
}

func (b *backend) pathTokenWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var cfg backendConfig
	e, err := req.Storage.Get(ctx, pathPatternConfig)
	if err != nil {
		return nil, err
	}
	if err := e.DecodeJSON(&cfg); err != nil {
		return nil, err
	}

	if cfg.KeyToken == "" {
		return nil, errBackendNotConfigured
	}

	svc := service.New(cfg.KeyToken)
	ts := time.Now().UnixNano()
	name := fmt.Sprintf("vault-%d", ts)
	token, err := svc.CreateAuthToken(ctx, name)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"ID":          token.ID,
			"bearerToken": token.Token,
		},
	}, nil
}
