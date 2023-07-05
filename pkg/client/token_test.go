package client

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

func TestCreateDeleteToken(t *testing.T) {
	recordHelper(t, "auth_token", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, c *Client) {
		require.NotNil(t, c.httpClient)
		pfx := "vault-plugin-secrets-vercel-fixtures-token"
		ts := time.Now().UnixNano()
		name := fmt.Sprintf("%s-%d", pfx, ts)
		res, err := c.CreateAuthToken(ctx, &CreateAuthTokenRequest{
			Name: name,
		})
		require.NoError(t, err)

		res2, err := c.DeleteAuthToken(ctx, &DeleteAuthTokenRequest{
			ID: res.Token.ID,
		})

		require.NoError(t, err)

		require.Equal(t, res2.ID, res.Token.ID)
	})
}
