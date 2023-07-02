package plugin

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	pathPatternConfig = "config"
	pathConfigAPIKey  = "api_key"
)

var (
	errMissingAPIKey       = errors.New("missing API key from configuration")
	errTypeAssertionFailed = errors.New("type assertion failed")
)

type backendConfig struct {
	APIKey string `json:"api_key"`
}

func (b *backend) pathConfig() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: pathPatternConfig,

			Fields: map[string]*framework.FieldSchema{
				pathConfigAPIKey: {
					Type:        framework.TypeString,
					Description: "API key for the Vercel account.",
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

	if e == nil {
		return nil, nil
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

	if config.APIKey == "" {
		return nil, errMissingAPIKey
	}

	e, err := logical.StorageEntryJSON(pathPatternConfig, config)
	if err != nil {
		return nil, err
	}

	if err = req.Storage.Put(ctx, e); err != nil {
		return nil, err
	}

	return &logical.Response{}, nil
}
