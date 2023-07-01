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
	newToken, err := a.CreateAuthToken(ctx, "foobar")
	if err != nil {
		t.Fatal(err)
	}
	if newToken.Token == "" {
		t.Fatal("empty token")
	}
}
