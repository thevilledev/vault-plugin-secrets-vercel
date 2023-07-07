package plugin

import (
	"context"

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
	return &logical.Response{
		Data: map[string]any{
			"build_date":        version.BuildDate,
			"build_version":     version.Version,
			"build_commit":      version.Commit,
			"build_commit_date": version.CommitDate,
			"build_branch":      version.Branch,
			"build_tag":         version.Tag,
			"build_dirty":       version.Dirty,
		},
	}, nil
}
