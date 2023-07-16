package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/thevilledev/vault-plugin-secrets-vercel/internal/client"
)

var (
	errInvalidTTL = errors.New("invalid ttl")
)

type Service struct {
	client client.Client
}

type Token struct {
	ID    string
	Token string
}

func New(apiKey string) *Service {
	c := &http.Client{}

	var ac client.Client

	if apiKey == "mock" {
		ac = client.NewMockClient()
	} else {
		ac = client.NewAPIClient(apiKey, c)
	}

	return &Service{
		client: ac,
	}
}

func NewWithBaseURL(apiKey string, baseURL string) *Service {
	c := &http.Client{}

	var ac client.Client

	if apiKey == "mock" {
		ac = client.NewMockClient()
	} else {
		ac = client.NewAPIClientWithBaseURL(apiKey, c, baseURL)
	}

	return &Service{
		client: ac,
	}
}

func (s *Service) CreateAuthToken(ctx context.Context, name string, ttl int64, teamID string) (string, string, error) {
	if ttl <= 0 {
		return "", "", errInvalidTTL
	}

	expiresAt := time.Now().Add(time.Duration(ttl) * time.Second).UTC().UnixMilli()
	r, err := s.client.CreateAuthToken(ctx, &client.CreateAuthTokenRequest{
		Name:      name,
		ExpiresAt: expiresAt,
		TeamID:    teamID,
	})

	if err != nil {
		return "", "", err
	}

	return r.Token.ID, r.BearerToken, nil
}

func (s *Service) DeleteAuthToken(ctx context.Context, id string) (string, error) {
	r, err := s.client.DeleteAuthToken(ctx, &client.DeleteAuthTokenRequest{
		ID: id,
	})

	if err != nil {
		return "", err
	}

	return r.ID, err
}
