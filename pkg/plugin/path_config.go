package plugin

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/thevilledev/vault-plugin-secrets-vercel/pkg/client"
)

const (
	pathPatternConfig = "config"
	pathConfigAPIKey  = "api_key"
	pathConfigBaseURL = "base_url"
	pathConfigMaxTTL  = "max_ttl"
	defaultMaxTTL     = 600
)

var (
	errMissingAPIKey       = errors.New("missing API key from configuration")
	errTypeAssertionFailed = errors.New("type assertion failed")
)

type backendConfig struct {
	APIKey  string        `json:"api_key"`
	BaseURL string        `json:"base_url"`
	MaxTTL  time.Duration `json:"max_ttl"`
}

func (b *backend) pathConfig() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: pathPatternConfig,

			Fields: map[string]*framework.FieldSchema{
				pathConfigAPIKey: {
					Type:        framework.TypeString,
					Description: "API key for the Vercel account.",
					Required:    true,
				},
				pathConfigBaseURL: {
					Type:        framework.TypeString,
					Description: "Optional API base URL used by this backend.",
				},
				pathConfigMaxTTL: {
					Type:    framework.TypeDurationSecond,
					Default: defaultMaxTTL,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathConfigWrite,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathConfigWrite,
				},
			},
		},
	}
}

func (b *backend) getConfig(ctx context.Context, storage logical.Storage) (*backendConfig, error) {
	var config backendConfig

	e, err := storage.Get(ctx, pathPatternConfig)
	if err != nil {
		return nil, err
	}

	if e == nil || len(e.Value) == 0 {
		return &backendConfig{}, nil
	}

	if err = e.DecodeJSON(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request,
	data *framework.FieldData) (*logical.Response, error) {
	config := &backendConfig{}

	if v, ok := data.GetOk(pathConfigAPIKey); ok {
		config.APIKey, ok = v.(string)
		if !ok {
			b.Logger().Trace("type assertion failed: %+v", v)
			return nil, errTypeAssertionFailed
		}
	}

	if v, ok := data.GetOk(pathConfigBaseURL); ok {
		config.BaseURL, ok = v.(string)
		if !ok {
			b.Logger().Trace("type assertion failed: %+v", v)
			return nil, errTypeAssertionFailed
		}
	}

	if vr, ok := data.GetOk(pathConfigMaxTTL); ok {
		v, ta := vr.(int)
		if !ta {
			b.Logger().Trace("type assertion failed: %+v", v)
			return nil, errTypeAssertionFailed
		}

		ttl := time.Duration(v) * time.Second

		config.MaxTTL = time.Duration(ttl.Seconds())
	}

	if config.APIKey == "" {
		return nil, errMissingAPIKey
	}

	if config.BaseURL == "" {
		config.BaseURL = client.DefaultBaseURL
	}

	if config.MaxTTL == 0 {
		config.MaxTTL = defaultMaxTTL
	}

	e, err := logical.StorageEntryJSON(pathPatternConfig, config)
	if err != nil {
		return nil, err
	}

	if err = req.Storage.Put(ctx, e); err != nil {
		return nil, err
	}

	b.Logger().Info("config initialised")

	return &logical.Response{}, nil
}
