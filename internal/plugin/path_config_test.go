package plugin

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestBackend_PathConfigRead(t *testing.T) {
	t.Parallel()

	t.Run("ReadConfiguration", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		b, storage := newTestBackend(t, nil)

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
		b, storage := newTestBackend(t, nil)

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
		b, storage := newTestBackend(t, nil)

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

	t.Run("WriteConfigurationWithValidTeamData", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		b, storage := newTestBackend(t, nil)

		_, err := b.HandleRequest(ctx, &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternConfig,
			Data: map[string]any{
				"api_key":         "foo",
				"default_team_id": "bar",
			},
		})
		require.NoError(t, err)

		cfg, err := b.getConfig(ctx, storage)
		require.NoError(t, err)
		require.NotNil(t, cfg)
		require.Equal(t, cfg.APIKey, "foo")
		require.Equal(t, cfg.DefaultTeamID, "bar")
	})

	t.Run("WriteConfigurationWithStorageFail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		disabledOps := []logical.Operation{logical.CreateOperation}
		b, storage := newTestBackend(t, disabledOps)

		_, err := b.HandleRequest(ctx, &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternConfig,
			Data: map[string]any{
				"api_key": "foo",
			},
		})
		require.Error(t, err)
	})
}

func TestBackend_PathConfigDelete(t *testing.T) {
	t.Parallel()

	t.Run("DeleteWithoutInit", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		b, storage := newTestBackend(t, nil)

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
		b, storage := newTestBackend(t, nil)

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

	t.Run("DeleteWithStorageFail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		disabledOps := []logical.Operation{logical.ReadOperation, logical.DeleteOperation}
		b, storage := newTestBackend(t, disabledOps)

		res, err := b.HandleRequest(ctx, &logical.Request{
			Storage:   storage,
			Operation: logical.DeleteOperation,
			Path:      pathPatternConfig,
		})
		require.Equal(t, err, errGetConfig)
		require.Nil(t, res)
	})
	t.Run("DeleteWithStorageFailOnRead", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		disabledOps := []logical.Operation{logical.DeleteOperation}
		b, storage := newTestBackend(t, disabledOps)

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
		require.Error(t, err)
		require.Nil(t, res)
	})
}

func TestBackend_PathConfigGet(t *testing.T) {
	t.Parallel()

	t.Run("GetConfig", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		b, storage := newTestBackend(t, nil)

		cfg, err := b.getConfig(ctx, storage)
		require.Nil(t, err)
		require.Nil(t, cfg)
	})

	t.Run("GetConfigStorageFail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		disabledOps := []logical.Operation{logical.ReadOperation}
		b, storage := newTestBackend(t, disabledOps)

		_, err := b.getConfig(ctx, storage)
		require.Equal(t, err, errGetConfig)
	})
}
