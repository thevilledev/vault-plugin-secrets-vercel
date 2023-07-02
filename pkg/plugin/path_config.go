package plugin

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	pathPatternConfig  = "config"
	pathConfigKeyToken = "key_token"

	errMissingMainAPIToken = errors.New("missing main API token configuration")
)

type backendConfig struct {
	KeyToken string `json:"key_token"`
}

func (b *backend) pathConfig() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: pathPatternConfig,

			Fields: map[string]*framework.FieldSchema{
				pathConfigKeyToken: {
					Type:        framework.TypeString,
					Description: "Main API key for the Vercel account.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleConfigRead,
					Summary:  "Retrieve configuration from Vercel secrets plugin.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleConfigWrite,
					Summary:  "Update configuration for an existing Vercel secrets plugin.",
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.handleConfigWrite,
					Summary:  "Create configuration for Vercel secrets plugin.",
				},
			},
			ExistenceCheck: b.handleConfigExistenceCheck,
		},
	}
}

func (b *backend) handleConfigExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %s", err)
	}

	return out != nil, nil
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
	if err := e.DecodeJSON(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (b *backend) handleConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	} else if config == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			pathConfigKeyToken: config.KeyToken,
		},
	}, nil
}

func (b *backend) handleConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config := &backendConfig{}

	if v, ok := data.GetOk(pathConfigKeyToken); ok {
		config.KeyToken = v.(string)
	}

	if config.KeyToken == "" {
		return nil, errMissingMainAPIToken
	}

	e, err := logical.StorageEntryJSON(pathPatternConfig, config)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, e); err != nil {
		return nil, err
	}

	res := &logical.Response{}
	res.AddWarning("The 'read' endpoint returns sensitive information. Please set an ACL.")

	return res, nil
}
