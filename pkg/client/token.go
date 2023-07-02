package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type CreateAuthTokenRequest struct {
	Name      string `json:"name"`
	ExpiresAt int    `json:"expiresAt,omitempty"`
}

type CreateAuthTokenResponse struct {
	Token       Token  `json:"token"`
	BearerToken string `json:"bearerToken"`
}

type Token struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Origin string `json:"origin"`
	// TODO scopes
}

func (c *Client) CreateAuthToken(ctx context.Context, req *CreateAuthTokenRequest) (*CreateAuthTokenResponse, error) {
	resp := &CreateAuthTokenResponse{}
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	res, err := c.do(ctx, http.MethodPost, "/user/tokens", b, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}
