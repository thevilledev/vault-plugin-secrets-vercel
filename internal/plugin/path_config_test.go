package plugin

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
	"github.com/thevilledev/vault-plugin-secrets-vercel/internal/client"
)

func TestConfig_Get(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		input       []byte
		disabledOps []logical.Operation
		cfg         *backendConfig
		expError    string
	}{
		"empty": {},
		"default": {
			input: []byte(`{"api_key": "foo"}`),
			cfg:   &backendConfig{APIKey: "foo"},
		},
		"invalid config json": {
			input:    []byte(`lorem ipsum`),
			cfg:      &backendConfig{},
			expError: "failed to decode config",
		},
		"storage fail": {
			disabledOps: []logical.Operation{
				logical.ReadOperation,
			},
			expError: "failed to get config from storage",
		},
	}
	for name, tc := range cases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			b, storage := newTestBackend(t, tc.disabledOps)

			if tc.input != nil {
				require.NoError(t, storage.Put(ctx, &logical.StorageEntry{
					Key:   pathPatternConfig,
					Value: tc.input,
				}))
			}

			res, err := b.getConfig(ctx, storage)
			if tc.expError != "" {
				require.EqualError(t, err, tc.expError)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConfig_Read(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		disabledOps []logical.Operation
		inputConfig []byte
		expError    string
	}{
		"read configuration": {
			inputConfig: []byte(`{"api_key": "foo"}`),
			expError:    "unsupported operation",
		},
	}
	for name, tc := range cases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			b, storage := newTestBackend(t, tc.disabledOps)

			if tc.inputConfig != nil {
				require.NoError(t, storage.Put(ctx, &logical.StorageEntry{
					Key:   pathPatternConfig,
					Value: tc.inputConfig,
				}))
			}

			res, err := b.HandleRequest(ctx, &logical.Request{
				Storage:   storage,
				Operation: logical.ReadOperation,
				Path:      pathPatternConfig,
			})

			if tc.expError != "" {
				require.EqualError(t, err, tc.expError)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConfig_Write(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		disabledOps []logical.Operation
		data        map[string]any
		expError    string
		expConfig   *backendConfig
	}{
		"write configuration with empty data": {
			data:     map[string]any{},
			expError: "missing api key from configuration",
		},
		"write configuration with valid data": {
			data: map[string]any{
				"api_key": "foo",
			},
			expConfig: &backendConfig{
				APIKey:  "foo",
				BaseURL: client.DefaultBaseURL,
				MaxTTL:  defaultMaxTTL,
			},
		},
		"write configuration with valid team data": {
			data: map[string]any{
				"api_key":         "foo",
				"default_team_id": "bar",
			},
			expConfig: &backendConfig{
				APIKey:        "foo",
				BaseURL:       client.DefaultBaseURL,
				MaxTTL:        defaultMaxTTL,
				DefaultTeamID: "bar",
			},
		},
		"write configuration with custom url and ttl": {
			data: map[string]any{
				"api_key":         "foo",
				"base_url":        "http://baseurl",
				"max_ttl":         10,
				"default_team_id": "bar",
			},
			expConfig: &backendConfig{
				APIKey:        "foo",
				BaseURL:       "http://baseurl",
				MaxTTL:        10,
				DefaultTeamID: "bar",
			},
		},
		"write configuration with storage fail": {
			disabledOps: []logical.Operation{
				logical.CreateOperation,
			},
			data: map[string]any{
				"api_key": "foo",
			},
			expError: "failed to write config to storage",
		},
	}
	for name, tc := range cases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			b, storage := newTestBackend(t, tc.disabledOps)

			res, err := b.HandleRequest(ctx, &logical.Request{
				Storage:   storage,
				Operation: logical.CreateOperation,
				Path:      pathPatternConfig,
				Data:      tc.data,
			})
			if tc.expError != "" {
				require.EqualError(t, err, tc.expError)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				cfg, errg := b.getConfig(ctx, storage)
				require.NoError(t, errg)
				require.Equal(t, cfg, tc.expConfig)
			}
		})
	}
}

func TestConfig_Delete(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		input       []byte
		disabledOps []logical.Operation
		expError    string
	}{
		"delete without an init": {
			expError: "backend not configured",
		},
		"delete success": {
			input: []byte(`{"api_key": "foo"}`),
		},
		"delete with storage fail": {
			disabledOps: []logical.Operation{
				logical.ReadOperation,
				logical.DeleteOperation,
			},
			input:    []byte(`{"api_key": "foo"}`),
			expError: "failed to get config from storage",
		},
		"delete with storage delete fail": {
			disabledOps: []logical.Operation{
				logical.DeleteOperation,
			},
			input:    []byte(`{"api_key": "foo"}`),
			expError: "failed to delete config from storage",
		},
	}
	for name, tc := range cases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			b, storage := newTestBackend(t, tc.disabledOps)

			if tc.input != nil {
				require.NoError(t, storage.Put(ctx, &logical.StorageEntry{
					Key:   pathPatternConfig,
					Value: tc.input,
				}))
			}

			res, err := b.HandleRequest(ctx, &logical.Request{
				Storage:   storage,
				Operation: logical.DeleteOperation,
				Path:      pathPatternConfig,
			})
			if tc.expError != "" {
				require.EqualError(t, err, tc.expError)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				res, errg := b.getConfig(ctx, storage)
				require.NoError(t, errg)
				require.Nil(t, res)
			}
		})
	}
}
