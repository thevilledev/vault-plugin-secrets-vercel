package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

func TestTokenMock(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			t.Helper()

			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write(nil)
		}),
	)
	defer ts.Close()

	t.Run("create token forbidden", func(t *testing.T) {
		ctx := context.Background()
		hc := &http.Client{}

		k := NewAPIClientWithBaseURL("foo", hc, ts.URL)
		r, err := k.CreateAuthToken(ctx, &CreateAuthTokenRequest{Name: "foo"})
		require.Nil(t, r)
		require.Error(t, err)
	})

	t.Run("create token bogus url", func(t *testing.T) {
		ctx := context.Background()
		hc := &http.Client{}
		k := NewAPIClientWithBaseURL("foo", hc, "http://localhost:69696")
		r, err := k.CreateAuthToken(ctx, &CreateAuthTokenRequest{Name: "foo"})
		require.Nil(t, r)
		require.Error(t, err)
	})

	t.Run("delete token forbidden", func(t *testing.T) {
		ctx := context.Background()
		hc := &http.Client{}

		k := NewAPIClientWithBaseURL("foo", hc, ts.URL)
		r, err := k.DeleteAuthToken(ctx, &DeleteAuthTokenRequest{ID: "foo"})
		require.Nil(t, r)
		require.Error(t, err)
	})

	t.Run("delete token bogus url", func(t *testing.T) {
		ctx := context.Background()
		hc := &http.Client{}
		k := NewAPIClientWithBaseURL("foo", hc, "http://localhost:69696")
		r, err := k.DeleteAuthToken(ctx, &DeleteAuthTokenRequest{ID: "foo"})
		require.Nil(t, r)
		require.Error(t, err)
	})
}

func TestCreateDeleteToken(t *testing.T) {
	t.Parallel()

	recordHelper(t, "auth_token", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, c *APIClient) {
		t.Helper()

		require.NotNil(t, c.httpClient)
		pfx := "vault-plugin-secrets-vercel-fixtures-token"
		ts := time.Now().UnixNano()
		name := fmt.Sprintf("%s-%d", pfx, ts)

		res, err := c.CreateAuthToken(ctx, nil)
		require.Equal(t, err, errEmptyReq)
		require.Nil(t, res)

		res, err = c.CreateAuthToken(ctx, &CreateAuthTokenRequest{
			Name: name,
		})
		require.NoError(t, err)

		resd, err := c.DeleteAuthToken(ctx, nil)
		require.Equal(t, err, errEmptyReq)
		require.Nil(t, resd)

		resd, err = c.DeleteAuthToken(ctx, &DeleteAuthTokenRequest{
			ID: res.Token.ID,
		})

		require.NoError(t, err)

		require.Equal(t, resd.ID, res.Token.ID)
	})
}

func TestCreateDeleteTokenTeam(t *testing.T) {
	t.Parallel()
	recordHelper(t, "auth_token_team", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, c *APIClient) {
		t.Helper()

		require.NotNil(t, c.httpClient)
		pfx := "vault-plugin-secrets-vercel-fixtures-token-team"
		ts := time.Now().UnixNano()
		name := fmt.Sprintf("%s-%d", pfx, ts)
		res, err := c.CreateAuthToken(ctx, &CreateAuthTokenRequest{
			Name:   name,
			TeamID: "thevilledev-team-1",
		})
		require.NoError(t, err)

		res2, err := c.DeleteAuthToken(ctx, &DeleteAuthTokenRequest{
			ID: res.Token.ID,
		})

		require.NoError(t, err)

		require.Equal(t, res2.ID, res.Token.ID)
	})
}
