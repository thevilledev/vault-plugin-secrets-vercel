package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/thevilledev/vault-plugin-secrets-vercel/internal/service"
)

const (
	keyPrefix            = "vault-plugin-secrets-vercel"
	pathPatternToken     = "token"
	pathTokenID          = "token_id"
	pathTokenBearerToken = "bearer_token"
	pathTokenTTL         = "ttl"
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
				pathTokenTTL: {
					Type:        framework.TypeDurationSecond,
					Description: "TTL for the generated API key. Less than or equal to the maximum TTL set in configuration.",
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

func (b *backend) pathTokenWrite(ctx context.Context, req *logical.Request,
	data *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if cfg.APIKey == "" {
		return nil, errMissingAPIKey
	}

	ttl := int64(0)

	if vr, ok := data.GetOk(pathTokenTTL); ok {
		v, ta := vr.(int)
		if !ta {
			b.Logger().Trace("type assertion failed: %+v", v)

			return nil, errTypeAssertionFailed
		}

		ttl = int64(v)
	}

	if ttl == 0 {
		ttl = int64(cfg.MaxTTL)
	}

	if ttl > int64(cfg.MaxTTL) {
		return nil, fmt.Errorf("TTL %d exceeds maximum of %d", ttl, int64(cfg.MaxTTL))
	}

	svc := service.NewWithBaseURL(cfg.APIKey, cfg.BaseURL)
	ts := time.Now().UnixNano()
	name := fmt.Sprintf("%s-%d", keyPrefix, ts)

	b.Logger().Info(fmt.Sprintf("creating token with %s and with TTL of %d", name, ttl))

	tokenID, bearerToken, err := svc.CreateAuthToken(ctx, name, ttl)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]any{
			pathTokenID:          tokenID,
			pathTokenBearerToken: bearerToken,
		},
		Secret: &logical.Secret{
			InternalData: map[string]any{
				"secret_type": backendSecretType,
				pathTokenID:   tokenID,
			},
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Duration(ttl) * time.Second,
			},
		},
	}, nil
}
