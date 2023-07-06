package plugin

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestBackend_PathConfigRead(t *testing.T) {
	t.Parallel()

	t.Run("ReadConfigurationValid", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		b, storage := newTestBackend(t)

		_, err := b.HandleRequest(ctx, &logical.Request{
			Storage:   storage,
			Operation: logical.ReadOperation,
			Path:      pathPatternConfig,
		})
		require.Error(t, err)
	})
}

func TestBackend_PathConfigWrite(t *testing.T) {
	t.Parallel()

	t.Run("WriteConfigurationWithEmptyData", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		b, storage := newTestBackend(t)

		_, err := b.HandleRequest(ctx, &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternConfig,
			Data:      map[string]any{},
		})
		require.Error(t, err)
		require.Equal(t, err, errMissingAPIKey)
	})

	t.Run("WriteConfigurationWithValidData", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		b, storage := newTestBackend(t)

		_, err := b.HandleRequest(ctx, &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternConfig,
			Data: map[string]any{
				"api_key": "foo",
			},
		})
		require.NoError(t, err)

		cfg, err := b.getConfig(ctx, storage)
		require.NoError(t, err)
		require.NotNil(t, cfg)
		require.Equal(t, cfg.APIKey, "foo")
	})
}

func TestBackend_PathConfigDelete(t *testing.T) {
	t.Parallel()

	t.Run("DeleteWithoutInit", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		b, storage := newTestBackend(t)

		res, err := b.HandleRequest(ctx, &logical.Request{
			Storage:   storage,
			Operation: logical.DeleteOperation,
			Path:      pathPatternConfig,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("DeleteSuccess", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		b, storage := newTestBackend(t)

		_, err := b.HandleRequest(ctx, &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternConfig,
			Data: map[string]any{
				"api_key": "foo",
			},
		})
		require.NoError(t, err)

		res, err := b.HandleRequest(ctx, &logical.Request{
			Storage:   storage,
			Operation: logical.DeleteOperation,
			Path:      pathPatternConfig,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})
}
