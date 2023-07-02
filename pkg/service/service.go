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

func (s *Service) CreateAuthToken(ctx context.Context, name string) (string, string, error) {
	r, err := s.apiClient.CreateAuthToken(ctx, &client.CreateAuthTokenRequest{
		Name: name,
	})
	if err != nil {
		return "", "", err
	}
	return r.Token.ID, r.BearerToken, nil
}

func (s *Service) DeleteAuthToken(ctx context.Context, id string) error {
	_, err := s.apiClient.DeleteAuthToken(ctx, &client.DeleteAuthTokenRequest{
		ID: id,
	})
	return err
}
