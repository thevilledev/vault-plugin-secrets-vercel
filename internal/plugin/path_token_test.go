package plugin

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestToken_Create(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		disabledOps   []logical.Operation
		cfgData       map[string]any
		tokenData     map[string]any
		expError      string
		expDataFields map[string]any
	}{
		"token without backend": {
			expError: "backend not configured",
		},
		"token success": {
			cfgData: map[string]any{
				"api_key": "mock",
			},
			tokenData: map[string]any{
				"name": "foo",
			},
			expDataFields: map[string]any{
				"bearer_token": "some-bearer-token",
			},
		},
		"token with storage fail": {
			disabledOps: []logical.Operation{
				logical.ReadOperation,
			},
			cfgData: map[string]any{
				"api_key": "mock",
			},
			tokenData: map[string]any{
				"name": "foo",
			},
			expError: "failed to get config from storage",
		},
		"token with backend fail": {
			cfgData: map[string]any{
				"api_key": "mock",
			},
			tokenData: map[string]any{
				"team_id": "force-fail",
			},
			expError: "failed to create token",
		},
		"token with conflicting ttl": {
			cfgData: map[string]any{
				"api_key": "mock",
				"max_ttl": 10,
			},
			tokenData: map[string]any{
				"name": "foo",
				"ttl":  11,
			},
			expError: "TTL exceeds the maximum value",
		},
		"token with default team id": {
			cfgData: map[string]any{
				"api_key":         "mock",
				"default_team_id": "default-team-id",
			},
			tokenData: map[string]any{
				"name": "foo",
			},
			expDataFields: map[string]any{
				"bearer_token": "some-bearer-token",
				"team_id":      "default-team-id",
			},
		},
		"token with custom team id": {
			cfgData: map[string]any{
				"api_key": "mock",
			},
			tokenData: map[string]any{
				"name":    "foo",
				"team_id": "custom-team-id",
			},
			expDataFields: map[string]any{
				"bearer_token": "some-bearer-token",
				"team_id":      "custom-team-id",
			},
		},
		"token with conflicting team ids": {
			cfgData: map[string]any{
				"api_key":         "mock",
				"default_team_id": "default-team-id",
			},
			tokenData: map[string]any{
				"name":    "foo",
				"team_id": "custom-team-id",
			},
			expError: "cannot override default_team_id",
		},
	}
	for name, tc := range cases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			b, storage := newTestBackend(t, tc.disabledOps)

			if tc.cfgData != nil {
				_, err := b.HandleRequest(ctx, &logical.Request{
					Storage:   storage,
					Operation: logical.CreateOperation,
					Path:      pathPatternConfig,
					Data:      tc.cfgData,
				})
				require.NoError(t, err)
			}

			r, err := b.HandleRequest(ctx, &logical.Request{
				Storage:   storage,
				Operation: logical.ReadOperation,
				Path:      pathPatternToken,
				Data:      tc.tokenData,
			})

			if tc.expError != "" {
				require.EqualError(t, err, tc.expError)
				require.Nil(t, r)
			} else {
				require.Equal(t, r.Secret.LeaseOptions.TTL, defaultMaxTTL*time.Second)

				for k, v := range tc.expDataFields {
					require.Equal(t, r.Data[k], v)
				}
			}
		})
	}
}
