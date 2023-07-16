package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

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
	})

	t.Run("client do", func(t *testing.T) {
		ctx := context.Background()
		hc := &http.Client{}

		k := NewAPIClientWithBaseURL("foo", hc, "http://doesnotexist")
		_, err := k.do(ctx, http.MethodGet, "/", nil, nil)
		require.Error(t, err)
	})

	t.Run("client without context", func(t *testing.T) {
		hc := &http.Client{}

		k := NewAPIClientWithBaseURL("foo", hc, ts.URL)
		//nolint:staticcheck
		_, err := k.do(nil, http.MethodGet, "/", nil, nil)
		require.Error(t, err)
	})
}
