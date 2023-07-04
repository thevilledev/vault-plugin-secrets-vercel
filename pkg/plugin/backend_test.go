package plugin

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		cfg     *logical.BackendConfig
		wantErr bool
	}{
		{
			name:    "Default",
			cfg:     &logical.BackendConfig{},
			wantErr: false,
		},
		{
			name:    "MissingConfig",
			cfg:     nil,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			b, err := Factory(ctx, tc.cfg)
			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, b)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func newTestBackend(t *testing.T) (*backend, logical.Storage) {
	t.Helper()

	config := logical.TestBackendConfig()
	config.StorageView = new(logical.InmemStorage)
	config.Logger = hclog.NewNullLogger()
	b, err := Factory(context.Background(), config)
	require.NoError(t, err)
	require.NotNil(t, b)

	return b.(*backend), config.StorageView
}

func TestBackend_Config(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		input   []byte
		cfg     *backendConfig
		wantErr bool
	}{
		{
			name:    "Default",
			input:   []byte(`{"api_key": "foo"}`),
			cfg:     &backendConfig{APIKey: "foo"},
			wantErr: false,
		},
		{
			name:    "InvalidJSON",
			input:   []byte(`lorem ipsum`),
			cfg:     &backendConfig{},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			b, storage := newTestBackend(t)

			if tc.input != nil {
				require.NoError(t, storage.Put(ctx, &logical.StorageEntry{
					Key:   pathPatternConfig,
					Value: tc.input,
				}))
			}

			_, err := b.getConfig(ctx, storage)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
