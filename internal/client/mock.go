package client

import (
	"context"
	"fmt"
	"time"
)

type MockClient struct {
	tokens map[string]any
}

func NewMockClient() *MockClient {
	return &MockClient{
		tokens: make(map[string]any, 0),
	}
}

func (m *MockClient) CreateAuthToken(_ context.Context,
	req *CreateAuthTokenRequest) (*CreateAuthTokenResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("empty name for token")
	}

	if req.TeamID == "force-fail" {
		return nil, fmt.Errorf("force fail")
	}

	r := &CreateAuthTokenResponse{
		Token: Token{
			ID:   fmt.Sprintf("%s-%d", req.Name, time.Now().UnixNano()),
			Name: req.Name,
		},
		BearerToken: "some-bearer-token",
	}

	m.tokens[r.Token.ID] = req.Name

	return r, nil
}

func (m *MockClient) DeleteAuthToken(_ context.Context,
	req *DeleteAuthTokenRequest) (*DeleteAuthTokenResponse, error) {
	if req.ID == "" {
		return nil, fmt.Errorf("empty id for token")
	}

	delete(m.tokens, req.ID)

	return &DeleteAuthTokenResponse{
		ID: req.ID,
	}, nil
}

func (m *MockClient) GetBaseURL() string {
	return ""
}
