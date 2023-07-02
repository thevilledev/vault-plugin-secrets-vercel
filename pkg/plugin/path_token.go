package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/thevilledev/vault-plugin-secrets-vercel/pkg/service"
)

const (
	keyPrefix = "vault-plugin-secrets-vercel"
)

var (
	pathPatternToken     = "token"
	pathTokenID          = "token_id"
	pathTokenBearerToken = "bearer_token"
)

func (b *backend) pathToken() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: pathPatternToken,
			Fields: map[string]*framework.FieldSchema{
				pathTokenID: {
					Type:        framework.TypeString,
					Description: "Token ID for the generated API key.",
				},
				pathTokenBearerToken: {
					Type:        framework.TypeString,
					Description: "Generated API key.",
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
	cfg, err := b.getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("backend is missing api key")
	}

	svc := service.New(cfg.APIKey)
	ts := time.Now().UnixNano()
	name := fmt.Sprintf("%s-%d", keyPrefix, ts)

	tokenID, bearerToken, err := svc.CreateAuthToken(ctx, name)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			pathTokenID:          tokenID,
			pathTokenBearerToken: bearerToken,
		},
		Secret: &logical.Secret{
			InternalData: map[string]interface{}{
				"secret_type": backendSecretType,
				pathTokenID:   tokenID,
			},
			LeaseOptions: logical.LeaseOptions{
				// TODO: add user-configurable TTL
				TTL: time.Until(time.Now().Add(10 * time.Second)),
			},
		},
	}, nil
}
