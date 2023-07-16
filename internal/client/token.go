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
	ExpiresAt int64  `json:"expiresAt,omitempty"`
	TeamID    string `json:"-"`
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
}

type DeleteAuthTokenRequest struct {
	ID string `json:"id"`
}

type DeleteAuthTokenResponse struct {
	ID string `json:"tokenId"`
}

func (c *APIClient) CreateAuthToken(ctx context.Context,
	req *CreateAuthTokenRequest) (*CreateAuthTokenResponse, error) {
	resp := &CreateAuthTokenResponse{}

	if req == nil {
		return nil, errEmptyReq
	}

	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	p := make(map[string]string, 0)
	if req.TeamID != "" {
		p["teamId"] = req.TeamID
	}

	res, err := c.do(ctx, http.MethodPost, "/user/tokens", b, p)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	validStatusAbove := 200
	invalidStatusBelow := 300

	ok := res.StatusCode >= validStatusAbove && res.StatusCode < invalidStatusBelow
	if !ok {
		return nil, fmt.Errorf("http error %d with response body '%+v'", res.StatusCode, string(body))
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *APIClient) DeleteAuthToken(ctx context.Context,
	req *DeleteAuthTokenRequest) (*DeleteAuthTokenResponse, error) {
	resp := &DeleteAuthTokenResponse{}

	if req == nil {
		return nil, errEmptyReq
	}

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

	validStatusAbove := 200
	invalidStatusBelow := 300

	ok := res.StatusCode >= validStatusAbove && res.StatusCode < invalidStatusBelow
	if !ok {
		return nil, fmt.Errorf("http error %d with response body '%+v'", res.StatusCode, string(body))
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
