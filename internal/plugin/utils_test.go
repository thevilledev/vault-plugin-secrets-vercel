package plugin

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func newTestBackend(t *testing.T, disabledOps []logical.Operation) (*backend, logical.Storage) {
	t.Helper()

	config := logical.TestBackendConfig()
	sw := new(logical.InmemStorage)

	for _, v := range disabledOps {
		switch v {
		case logical.ReadOperation:
			sw.Underlying().FailGet(true)
		case logical.UpdateOperation:
			sw.Underlying().FailPut(true)
		case logical.DeleteOperation:
			sw.Underlying().FailDelete(true)
		case logical.CreateOperation:
			sw.Underlying().FailPut(true)
		case logical.RevokeOperation:
			sw.Underlying().FailDelete(true)
		}
	}

	config.StorageView = sw
	config.Logger = hclog.NewNullLogger()
	br, err := Factory(context.Background(), config)
	require.NoError(t, err)
	require.NotNil(t, br)

	b, ok := br.(*backend)
	require.Equal(t, ok, true)

	return b, config.StorageView
}
