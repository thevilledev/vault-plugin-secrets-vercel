package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

	p := make(map[string]string, 1)
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

	if !successStatus(res.StatusCode) {
		return nil, newHTTPError(res.StatusCode, body)
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Token.ID == "" || resp.BearerToken == "" {
		return nil, errInvalidCreateAuthTokenResponse
	}

	return resp, nil
}

func (c *APIClient) DeleteAuthToken(ctx context.Context,
	req *DeleteAuthTokenRequest) (*DeleteAuthTokenResponse, error) {
	resp := &DeleteAuthTokenResponse{}

	if req == nil {
		return nil, errEmptyReq
	}

	if req.ID == "" {
		return nil, errMissingTokenID
	}

	path := fmt.Sprintf("%s/%s", "/user/tokens", url.PathEscape(req.ID))

	res, err := c.do(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if !successStatus(res.StatusCode) {
		return nil, newHTTPError(res.StatusCode, body)
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.ID == "" {
		return nil, errInvalidDeleteAuthTokenResponse
	}

	return resp, nil
}

func successStatus(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices
}
