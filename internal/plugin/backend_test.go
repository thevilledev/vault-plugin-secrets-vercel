package plugin

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		name     string
		cfg      *logical.BackendConfig
		expError string
	}{
		"default": {
			cfg: &logical.BackendConfig{},
		},
		"missing config": {
			cfg:      nil,
			expError: "configuration passed into backend is nil",
		},
	}

	for name, tc := range cases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			b, err := Factory(ctx, tc.cfg)
			if tc.expError != "" {
				require.EqualError(t, err, tc.expError)
				require.Nil(t, b)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
