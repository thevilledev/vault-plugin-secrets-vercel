package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			t.Helper()

			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write(nil)
		}),
	)
	defer ts.Close()

	t.Run("create client", func(t *testing.T) {
		hc := &http.Client{}
		k := NewAPIClient("api key", hc)
		require.Equal(t, k.baseURL, DefaultBaseURL)
		require.Equal(t, k.token, "api key")
		require.Equal(t, defaultHTTPTimeout, k.httpClient.Timeout)
		require.Zero(t, hc.Timeout)
	})

	t.Run("create client with nil http client", func(t *testing.T) {
		k := NewAPIClient("api key", nil)
		require.NotNil(t, k.httpClient)
		require.Equal(t, defaultHTTPTimeout, k.httpClient.Timeout)
	})

	t.Run("create client preserves custom timeout", func(t *testing.T) {
		hc := &http.Client{Timeout: time.Second}
		k := NewAPIClient("api key", hc)
		require.Same(t, hc, k.httpClient)
		require.Equal(t, time.Second, k.httpClient.Timeout)
	})

	t.Run("client do", func(t *testing.T) {
		ctx := context.Background()
		hc := &http.Client{}

		k := NewAPIClientWithBaseURL("foo", hc, "http://doesnotexist")
		u := k.GetBaseURL()
		require.Equal(t, u, "http://doesnotexist")
		_, err := k.do(ctx, http.MethodGet, "/", nil, nil)
		require.Error(t, err)
	})

	t.Run("client do preserves base path and query", func(t *testing.T) {
		ctx := context.Background()

		var gotAuth string
		var gotPath string
		var gotQuery string

		srv := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				t.Helper()

				gotAuth = r.Header.Get("Authorization")
				gotPath = r.URL.EscapedPath()
				gotQuery = r.URL.RawQuery
				w.WriteHeader(http.StatusNoContent)
			}),
		)
		defer srv.Close()

		k := NewAPIClientWithBaseURL("foo", nil, srv.URL+"/v3/")
		_, err := k.do(ctx, http.MethodGet, "/user/tokens", nil, map[string]string{
			"teamId": "team id",
		})
		require.NoError(t, err)
		require.Equal(t, "Bearer foo", gotAuth)
		require.Equal(t, "/v3/user/tokens", gotPath)
		require.Equal(t, "teamId=team+id", gotQuery)
	})

	t.Run("client without context", func(t *testing.T) {
		hc := &http.Client{}

		k := NewAPIClientWithBaseURL("foo", hc, ts.URL)
		//nolint:staticcheck
		_, err := k.do(nil, http.MethodGet, "/", nil, nil)
		require.Error(t, err)
	})
}
