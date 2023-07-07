package plugin

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/thevilledev/vault-plugin-secrets-vercel/internal/version"
)

const (
	pathPatternInfo = "info"
)

func (b *backend) pathInfo() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: pathPatternInfo,
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathInfoRead,
				},
			},
		},
	}
}

func (b *backend) pathInfoRead(
	_ context.Context,
	_ *logical.Request,
	_ *framework.FieldData,
) (*logical.Response, error) {
	var m map[string]any

	v := version.New()

	bs, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bs, &m)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: m,
	}, nil
}
