package service

import (
	"context"

	"github.com/thevilledev/vault-plugin-secrets-vercel/pkg/client"
)

type Service struct {
	apiClient *client.Client
}

type Token struct {
	ID    string
	Token string
}

func New(apiKey string) *Service {
	return &Service{
		apiClient: client.New(apiKey),
	}
}

func (s *Service) CreateAuthToken(ctx context.Context, name string) (*Token, error) {
	r, err := s.apiClient.CreateAuthToken(ctx, &client.CreateAuthTokenRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return &Token{
		ID:    r.Token.ID,
		Token: r.BearerToken,
	}, nil
}
