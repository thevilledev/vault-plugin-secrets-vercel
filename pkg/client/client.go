package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL     = "https://api.vercel.com/v3"
	httpTimeout = 60 * time.Second
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

type CreateAuthTokenRequest struct {
	Name      string `json:"name"`
	ExpiresAt int    `json:"expiresAt,omitempty"`
}

type CreateAuthTokenResponse struct {
	Token       TokenResponse `json:"token"`
	BearerToken string        `json:"bearerToken"`
}

type TokenResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Origin string `json:"origin"`
	// TODO scopes
}

func New(apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: httpTimeout,
		},
		token: apiKey,
	}
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

func (c *Client) do(ctx context.Context, method, endpoint string, body []byte, params map[string]string) (*http.Response, error) {
	u := fmt.Sprintf("%s%s", c.baseURL, endpoint)
	bearer := fmt.Sprintf("Bearer %s", c.token)
	req, err := http.NewRequestWithContext(ctx, method, u, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer)

	q := req.URL.Query()
	for key, val := range params {
		q.Set(key, val)
	}
	req.URL.RawQuery = q.Encode()

	return c.httpClient.Do(req)
}
