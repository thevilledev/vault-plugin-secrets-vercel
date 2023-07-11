package service

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

//nolint:paralleltest
func TestIntegration_Token(t *testing.T) {
	if os.Getenv("ACC_TEST") == "" {
		t.Skip("test skipped as ACC_TEST environment variable is not set")
	}

	token := os.Getenv("VERCEL_TOKEN")
	a := New(token)
	ctx := context.Background()

	ttl := int64(10)
	teamID := ""
	name := fmt.Sprintf("%s-%d", "vault-plugin-secrets-vercel-service-test", time.Now().UTC().UnixMilli())
	tokenID, bearerToken, err := a.CreateAuthToken(ctx, name, ttl, teamID)
	require.NoError(t, err)
	require.NotEmpty(t, tokenID)
	require.NotEmpty(t, bearerToken)

	s, err := a.DeleteAuthToken(ctx, tokenID)
	require.NoError(t, err)
	require.Equal(t, s, tokenID)
}

//nolint:paralleltest
func TestIntegration_Token_Team(t *testing.T) {
	if os.Getenv("ACC_TEST") == "" {
		t.Skip("test skipped as ACC_TEST environment variable is not set")
	}

	token := os.Getenv("VERCEL_TOKEN")
	teamID := os.Getenv("VERCEL_TEAM_ID")
	a := New(token)
	ctx := context.Background()

	ttl := int64(10)
	name := fmt.Sprintf("%s-%d", "vault-plugin-secrets-vercel-service-test", time.Now().UTC().UnixMilli())
	tokenID, bearerToken, err := a.CreateAuthToken(ctx, name, ttl, teamID)
	require.NoError(t, err)
	require.NotEmpty(t, tokenID)
	require.NotEmpty(t, bearerToken)

	s, err := a.DeleteAuthToken(ctx, tokenID)
	require.NoError(t, err)
	require.Equal(t, s, tokenID)
}
