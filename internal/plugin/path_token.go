package plugin

import (
	"context"
	"errors"
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
	pathTokenTeamID      = "team_id"
	//nolint:gosec
	pathTokenTTLDescription = `
(Optional) TTL for the generated API key ("bearer token"). Less than or equal to the maximum TTL set in configuration.
Defaults to maximum TTL.`
	//nolint:gosec
	pathTokenTeamIDDescription = `
(Optional) Team ID used for generating the API key.
This acts as a scope for the key. It only has access to the given team.`
	pathTokenDescription = `
Supports only read operations. Token ID for the generated key is stored in the plugin backend for revocation purposes.
Generated bearer token is NOT stored in the plugin backend.
Tokens are automatically revoked & deleted by Vault once TTL hits zero.
Tokens cannot be renewed. Generate a new token if you need one.`
	//nolint:gosec
	pathTokenSynopsis = `
Generate a Vercel API token with the given TTL.`
)

var (
	errTokenMaxTTLExceeded         = errors.New("given TTL exceeds the maximum allowed value")
	errCannotOverrideDefaultTeamID = errors.New("cannot override team_id different than set on backend config")
)

func (b *backend) pathToken() []*framework.Path {
	return []*framework.Path{
		{
			Pattern:         pathPatternToken,
			HelpDescription: pathTokenDescription,
			HelpSynopsis:    pathTokenSynopsis,
			Fields: map[string]*framework.FieldSchema{
				pathTokenTTL: {
					Type:        framework.TypeDurationSecond,
					Description: pathTokenTTLDescription,
				},
				pathTokenTeamID: {
					Type:        framework.TypeString,
					Description: pathTokenTeamIDDescription,
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

	if cfg == nil {
		return nil, errBackendNotConfigured
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
		return nil, errTokenMaxTTLExceeded
	}

	teamID := cfg.DefaultTeamID

	if vr, ok := data.GetOk(pathTokenTeamID); ok {
		v, ta := vr.(string)
		if !ta {
			b.Logger().Trace("type assertion failed: %+v", v)

			return nil, errTypeAssertionFailed
		}

		if teamID != "" && v != "" && teamID != v {
			return nil, errCannotOverrideDefaultTeamID
		}

		teamID = v
	}

	svc := service.NewWithBaseURL(cfg.APIKey, cfg.BaseURL)
	ts := time.Now().UnixNano()
	name := fmt.Sprintf("%s-%d", keyPrefix, ts)

	b.Logger().Info(fmt.Sprintf("creating token with %s and with TTL of %d", name, ttl))

	tokenID, bearerToken, err := svc.CreateAuthToken(ctx, name, ttl, teamID)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]any{
			pathTokenID:          tokenID,
			pathTokenBearerToken: bearerToken,
			pathTokenTeamID:      teamID,
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
