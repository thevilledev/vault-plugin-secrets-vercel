package service

import (
	"context"
	"os"
	"testing"
)

func TestToken(t *testing.T) {
	token := os.Getenv("VERCEL_TOKEN")
	a := New(token)
	ctx := context.Background()
	tokenID, bearerToken, err := a.CreateAuthToken(ctx, "foobar")
	if err != nil {
		t.Fatal(err)
	}
	if tokenID == "" {
		t.Fatal("empty token")
	}
	if bearerToken == "" {
		t.Fatal("empty bearer token")
	}
}
