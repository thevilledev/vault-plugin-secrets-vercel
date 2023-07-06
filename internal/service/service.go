package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/thevilledev/vault-plugin-secrets-vercel/internal/client"
)

type Service struct {
	apiClient *client.Client
}

type Token struct {
	ID    string
	Token string
}

func New(apiKey string) *Service {
	c := &http.Client{}

	return &Service{
		apiClient: client.New(apiKey, c),
	}
}

func NewWithBaseURL(apiKey string, baseURL string) *Service {
	c := &http.Client{}

	return &Service{
		apiClient: client.NewWithBaseURL(apiKey, c, baseURL),
	}
}

func (s *Service) CreateAuthToken(ctx context.Context, name string, ttl int64) (string, string, error) {
	if ttl <= 0 {
		return "", "", fmt.Errorf("cannot create token with a ttl of 0")
	}

	expiresAt := time.Now().Add(time.Duration(ttl) * time.Second).UTC().UnixMilli()
	r, err := s.apiClient.CreateAuthToken(ctx, &client.CreateAuthTokenRequest{
		Name:      name,
		ExpiresAt: expiresAt,
	})

	if err != nil {
		return "", "", err
	}

	return r.Token.ID, r.BearerToken, nil
}

func (s *Service) DeleteAuthToken(ctx context.Context, id string) (string, error) {
	r, err := s.apiClient.DeleteAuthToken(ctx, &client.DeleteAuthTokenRequest{
		ID: id,
	})

	if err != nil {
		return "", err
	}

	return r.ID, err
}
