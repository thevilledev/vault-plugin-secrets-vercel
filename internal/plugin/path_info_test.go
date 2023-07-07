package plugin

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
	"github.com/thevilledev/vault-plugin-secrets-vercel/internal/version"
)

func TestPathInfo(t *testing.T) {
	t.Parallel()

	t.Run("ValidInfo", func(t *testing.T) {
		t.Parallel()

		version.Branch = "main"
		version.BuildDate = time.Now().String()
		version.Commit = "xyz"
		version.CommitDate = time.Now().String()
		version.Dirty = "false"
		version.Tag = "v0.0.1"
		version.Version = "0.0.1"

		ctx := context.Background()
		b, storage := newTestBackend(t)

		res, err := b.HandleRequest(ctx, &logical.Request{
			Storage:   storage,
			Operation: logical.ReadOperation,
			Path:      pathPatternInfo,
		})
		require.NoError(t, err)

		var vi version.VersionInfo
		d, err := json.Marshal(res.Data)
		require.NoError(t, err)
		err = json.Unmarshal(d, &vi)
		require.NoError(t, err)

		vn := version.New()
		require.Equal(t, *vn, vi)
	})
}
