package plugin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
	"github.com/thevilledev/vault-plugin-secrets-vercel/pkg/client"
)

func TestRevokeToken(t *testing.T) {
	t.Parallel()

	t.Run("RevokeToken", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		b, storage := newTestBackend(t)

		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				t.Helper()

				body, _ := json.Marshal(&client.DeleteAuthTokenResponse{
					ID: "zyzz",
				})
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(body)
			}),
		)
		defer ts.Close()

		_, err := b.HandleRequest(ctx, &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternConfig,
			Data: map[string]any{
				"api_key":  "foo",
				"base_url": ts.URL,
			},
		})
		require.NoError(t, err)

		r, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.RevokeOperation,
			Path:      pathPatternToken,
			Data:      map[string]any{},
			Secret: &logical.Secret{
				InternalData: map[string]any{
					"secret_type": backendSecretType,
					"token_id":    "zyzz",
				},
			},
		})
		require.NoError(t, err)
		require.Equal(t, r, &logical.Response{})
	})

	t.Run("RevokeTokenFail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		b, storage := newTestBackend(t)

		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				t.Helper()

				w.WriteHeader(http.StatusForbidden)
			}),
		)
		defer ts.Close()

		_, err := b.HandleRequest(ctx, &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternConfig,
			Data: map[string]any{
				"api_key":  "foo",
				"base_url": ts.URL,
			},
		})
		require.NoError(t, err)

		_, err = b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.RevokeOperation,
			Path:      pathPatternToken,
			Data:      map[string]any{},
			Secret: &logical.Secret{
				InternalData: map[string]any{
					"secret_type": backendSecretType,
					"token_id":    "zyzz",
				},
			},
		})
		require.Error(t, err)
	})
}
