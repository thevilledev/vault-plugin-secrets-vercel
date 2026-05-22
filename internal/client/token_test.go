package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

		var httpErr *HTTPError
		require.ErrorAs(t, err, &httpErr)
		require.Equal(t, http.StatusForbidden, httpErr.StatusCode)
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

		var httpErr *HTTPError
		require.ErrorAs(t, err, &httpErr)
		require.Equal(t, http.StatusForbidden, httpErr.StatusCode)
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

func TestTokenResponseHardening(t *testing.T) {
	t.Parallel()

	t.Run("create token redacts http error body", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		srv := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				t.Helper()

				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte(`{"bearerToken":"secret","message":"denied"}`))
			}),
		)
		defer srv.Close()

		c := NewAPIClientWithBaseURL("foo", nil, srv.URL)
		res, err := c.CreateAuthToken(ctx, &CreateAuthTokenRequest{Name: "foo"})
		require.Nil(t, res)

		var httpErr *HTTPError
		require.ErrorAs(t, err, &httpErr)
		require.Equal(t, http.StatusForbidden, httpErr.StatusCode)
		require.Contains(t, httpErr.Body, "[REDACTED]")
		require.NotContains(t, httpErr.Body, "secret")
	})

	t.Run("create token truncates http error body", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		srv := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				t.Helper()

				w.WriteHeader(http.StatusBadGateway)
				_, _ = w.Write([]byte(strings.Repeat("x", maxHTTPErrorBodyLength+1)))
			}),
		)
		defer srv.Close()

		c := NewAPIClientWithBaseURL("foo", nil, srv.URL)
		res, err := c.CreateAuthToken(ctx, &CreateAuthTokenRequest{Name: "foo"})
		require.Nil(t, res)

		var httpErr *HTTPError
		require.ErrorAs(t, err, &httpErr)
		require.Equal(t, http.StatusBadGateway, httpErr.StatusCode)
		require.Len(t, httpErr.Body, maxHTTPErrorBodyLength+len(truncatedHTTPBodyMarker))
		require.True(t, strings.HasSuffix(httpErr.Body, truncatedHTTPBodyMarker))
	})

	t.Run("create token validates success payload", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		srv := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				t.Helper()

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{}`))
			}),
		)
		defer srv.Close()

		c := NewAPIClientWithBaseURL("foo", nil, srv.URL)
		res, err := c.CreateAuthToken(ctx, &CreateAuthTokenRequest{Name: "foo"})
		require.Nil(t, res)
		require.ErrorIs(t, err, errInvalidCreateAuthTokenResponse)
	})

	t.Run("delete token escapes id path segment", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		tokenID := "foo/bar baz"
		srv := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				t.Helper()

				require.Equal(t, "/v3/user/tokens/foo%2Fbar%20baz", r.URL.EscapedPath())
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"tokenId":"foo/bar baz"}`))
			}),
		)
		defer srv.Close()

		c := NewAPIClientWithBaseURL("foo", nil, srv.URL+"/v3/")
		res, err := c.DeleteAuthToken(ctx, &DeleteAuthTokenRequest{ID: tokenID})
		require.NoError(t, err)
		require.Equal(t, tokenID, res.ID)
	})

	t.Run("delete token validates request id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := NewAPIClientWithBaseURL("foo", nil, "https://example.com")
		res, err := c.DeleteAuthToken(ctx, &DeleteAuthTokenRequest{})
		require.Nil(t, res)
		require.ErrorIs(t, err, errMissingTokenID)
	})

	t.Run("delete token validates success payload", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		srv := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				t.Helper()

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{}`))
			}),
		)
		defer srv.Close()

		c := NewAPIClientWithBaseURL("foo", nil, srv.URL)
		res, err := c.DeleteAuthToken(ctx, &DeleteAuthTokenRequest{ID: "foo"})
		require.Nil(t, res)
		require.ErrorIs(t, err, errInvalidDeleteAuthTokenResponse)
	})

	t.Run("http error supports errors as", func(t *testing.T) {
		t.Parallel()

		err := error(newHTTPError(http.StatusTooManyRequests, []byte("rate limited")))
		var httpErr *HTTPError
		require.True(t, errors.As(err, &httpErr))
		require.Equal(t, http.StatusTooManyRequests, httpErr.StatusCode)
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
