package service

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_Token(t *testing.T) {
	if os.Getenv("ACC_TEST") == "" {
		t.Skip("test skipped as ACC_TEST environment variable is not set")
	}

	token := os.Getenv("VERCEL_TOKEN")
	a := New(token)
	ctx := context.Background()

	tokenID, bearerToken, err := a.CreateAuthToken(ctx, "foobar")
	require.NoError(t, err)
	require.NotEmpty(t, tokenID)
	require.NotEmpty(t, bearerToken)

	s, err := a.DeleteAuthToken(ctx, tokenID)
	require.NoError(t, err)
	require.Equal(t, s, tokenID)
}
