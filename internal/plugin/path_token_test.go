package plugin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
	"github.com/thevilledev/vault-plugin-secrets-vercel/internal/client"
)

func TestToken_Create(t *testing.T) {
	t.Parallel()

	t.Run("CreateTokenWithoutBackend", func(t *testing.T) {
		t.Parallel()

		b, storage := newTestBackend(t, nil)

		r, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternToken,
			Data:      map[string]any{},
		})
		require.Equal(t, err, errBackendNotConfigured)
		require.Nil(t, r)
	})

	t.Run("CreateTokenWithValidBackend", func(t *testing.T) {
		t.Parallel()

		b, storage := newTestBackend(t, nil)

		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				t.Helper()

				body, _ := json.Marshal(&client.CreateAuthTokenResponse{
					Token: client.Token{
						ID:   "foo",
						Name: "bar",
					},
					BearerToken: "zyzz",
				})
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(body)
			}),
		)
		defer ts.Close()

		_, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternConfig,
			Data: map[string]interface{}{
				"api_key":  "foo",
				"base_url": ts.URL,
			},
		})
		require.NoError(t, err)

		r, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternToken,
			Data:      map[string]any{},
		})
		require.NoError(t, err)
		require.NotNil(t, r)
		require.Equal(t, r.Data["token_id"], "foo")
		require.Equal(t, r.Data["bearer_token"], "zyzz")
		require.Equal(t, r.Secret.LeaseOptions.TTL, defaultMaxTTL*time.Second)
		require.Equal(t, r.Secret.InternalData["token_id"], "foo")
	})

	t.Run("CreateTokenWithStorageFail", func(t *testing.T) {
		t.Parallel()

		disabledOps := []logical.Operation{logical.ReadOperation}
		b, storage := newTestBackend(t, disabledOps)

		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				t.Helper()

				body, _ := json.Marshal(&client.CreateAuthTokenResponse{
					Token: client.Token{
						ID:   "foo",
						Name: "bar",
					},
					BearerToken: "zyzz",
				})
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(body)
			}),
		)
		defer ts.Close()

		_, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternConfig,
			Data: map[string]interface{}{
				"api_key":  "foo",
				"base_url": ts.URL,
			},
		})
		require.NoError(t, err)

		r, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.ReadOperation,
			Path:      pathPatternToken,
			Data:      map[string]any{},
		})
		require.Error(t, err)
		require.Nil(t, r)
	})

	t.Run("CreateTokenWithUpstreamAPIError", func(t *testing.T) {
		t.Parallel()

		b, storage := newTestBackend(t, nil)

		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				t.Helper()

				body := []byte("not authorized")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write(body)
			}),
		)
		defer ts.Close()

		_, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternConfig,
			Data: map[string]interface{}{
				"api_key":  "foo",
				"base_url": ts.URL,
			},
		})
		require.NoError(t, err)

		r, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternToken,
			Data:      map[string]any{},
		})
		require.Error(t, err)
		require.Nil(t, r)
	})

	t.Run("CreateTokenWithConflictingTTLs", func(t *testing.T) {
		t.Parallel()

		b, storage := newTestBackend(t, nil)

		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				t.Helper()

				body, _ := json.Marshal(&client.CreateAuthTokenResponse{
					Token: client.Token{
						ID:   "foo",
						Name: "bar",
					},
					BearerToken: "zyzz",
				})
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(body)
			}),
		)
		defer ts.Close()

		_, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternConfig,
			Data: map[string]any{
				"api_key":  "foo",
				"base_url": ts.URL,
				"max_ttl":  10,
			},
		})
		require.NoError(t, err)

		_, err = b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternToken,
			Data: map[string]any{
				"ttl": 11,
			},
		})
		require.ErrorIs(t, err, errTokenMaxTTLExceeded)
	})

	t.Run("CreateTokenWithDefaultTeamID", func(t *testing.T) {
		t.Parallel()

		b, storage := newTestBackend(t, nil)

		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				t.Helper()

				body, _ := json.Marshal(&client.CreateAuthTokenResponse{
					Token: client.Token{
						ID:   "foo",
						Name: "bar",
					},
					BearerToken: "zyzz",
				})
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(body)
			}),
		)
		defer ts.Close()

		_, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternConfig,
			Data: map[string]any{
				"api_key":         "foo",
				"base_url":        ts.URL,
				"default_team_id": "foo",
			},
		})
		require.NoError(t, err)

		r, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternToken,
			Data:      map[string]any{},
		})
		require.NoError(t, err)
		require.Equal(t, r.Data["team_id"], "foo")
	})

	t.Run("CreateTokenWithTeamID", func(t *testing.T) {
		t.Parallel()

		b, storage := newTestBackend(t, nil)

		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				t.Helper()

				body, _ := json.Marshal(&client.CreateAuthTokenResponse{
					Token: client.Token{
						ID:   "foo",
						Name: "bar",
					},
					BearerToken: "zyzz",
				})
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(body)
			}),
		)
		defer ts.Close()

		_, err := b.HandleRequest(context.Background(), &logical.Request{
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
			Operation: logical.CreateOperation,
			Path:      pathPatternToken,
			Data: map[string]any{
				"team_id": "bar",
			},
		})
		require.NoError(t, err)
		require.Equal(t, r.Data["team_id"], "bar")
	})

	t.Run("CreateTokenWithConflictingTeamIDs", func(t *testing.T) {
		t.Parallel()

		b, storage := newTestBackend(t, nil)

		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				t.Helper()

				body, _ := json.Marshal(&client.CreateAuthTokenResponse{
					Token: client.Token{
						ID:   "foo",
						Name: "bar",
					},
					BearerToken: "zyzz",
				})
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(body)
			}),
		)
		defer ts.Close()

		_, err := b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternConfig,
			Data: map[string]any{
				"api_key":         "foo",
				"base_url":        ts.URL,
				"default_team_id": "foo",
			},
		})
		require.NoError(t, err)

		_, err = b.HandleRequest(context.Background(), &logical.Request{
			Storage:   storage,
			Operation: logical.CreateOperation,
			Path:      pathPatternToken,
			Data: map[string]any{
				"team_id": "bar",
			},
		})
		require.ErrorIs(t, err, errCannotOverrideDefaultTeamID)
	})
}
