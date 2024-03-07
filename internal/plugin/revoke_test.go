package plugin

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestToken_Revoke(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		disabledOps  []logical.Operation
		cfgData      map[string]any
		tokenData    map[string]any
		internalData map[string]any
		expError     string
	}{
		"token revocation without backend": {
			expError: "backend not configured",
		},
		"token revocation with storage fail": {
			disabledOps: []logical.Operation{
				logical.ReadOperation,
			},
			expError: "failed to get config from storage",
		},
		"token revocation success": {
			cfgData: map[string]any{
				"api_key": "mock",
			},
			tokenData: map[string]any{
				"name": "foo",
			},
		},
		"token revocation backend fail": {
			cfgData: map[string]any{
				"api_key": "mock",
			},
			tokenData: map[string]any{
				"name": "foo",
			},
			internalData: map[string]any{
				"secret_type": backendSecretType,
				"token_id":    "", // force backend mock fail
			},
			expError: "failed to revoke token",
		},
		"token revocation internal data fail": {
			cfgData: map[string]any{
				"api_key": "mock",
			},
			tokenData: map[string]any{
				"name": "foo",
			},
			internalData: map[string]any{
				"secret_type": backendSecretType,
			},
			expError: "missing internal data from secret",
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

			var tokenID string

			if tc.tokenData != nil {
				r, err := b.HandleRequest(ctx, &logical.Request{
					Storage:   storage,
					Operation: logical.ReadOperation,
					Path:      pathPatternToken,
					Data:      tc.tokenData,
				})
				require.NoError(t, err)

				tokenID, _ = r.Data["token_id"].(string)
			}

			id := map[string]any{
				"secret_type": backendSecretType,
				"token_id":    tokenID,
			}
			if tc.internalData != nil {
				id = tc.internalData
			}

			r, err := b.HandleRequest(ctx, &logical.Request{
				Storage:   storage,
				Operation: logical.RevokeOperation,
				Path:      pathPatternToken,
				Data:      map[string]any{},
				Secret: &logical.Secret{
					InternalData: id,
				},
			})

			if tc.expError != "" {
				require.EqualError(t, err, tc.expError)
				require.Nil(t, r)
			}
		})
	}
}
