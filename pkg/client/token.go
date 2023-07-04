package client

import (
	"context"
	"encoding/json"
	"fmt"
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

type DeleteAuthTokenRequest struct {
	ID string `json:"id"`
}

type DeleteAuthTokenResponse struct {
	ID string `json:"tokenId"`
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

	ok := res.StatusCode >= 200 && res.StatusCode < 300
	if !ok {
		return nil, fmt.Errorf("http error %d with response body '%+v'", res.StatusCode, string(body))
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *Client) DeleteAuthToken(ctx context.Context, req *DeleteAuthTokenRequest) (*DeleteAuthTokenResponse, error) {
	resp := &DeleteAuthTokenResponse{}
	path := fmt.Sprintf("%s/%s", "/user/tokens", req.ID)

	res, err := c.do(ctx, http.MethodDelete, path, nil, nil)
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
