package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	DefaultBaseURL     = "https://api.vercel.com/v3"
	defaultHTTPTimeout = 60 * time.Second
)

var (
	errEmptyReq = errors.New("empty req")
)

type Client interface {
	GetBaseURL() string
	DeleteAuthToken(ctx context.Context, req *DeleteAuthTokenRequest) (*DeleteAuthTokenResponse, error)
	CreateAuthToken(ctx context.Context, req *CreateAuthTokenRequest) (*CreateAuthTokenResponse, error)
}

type APIClient struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

func NewAPIClient(apiKey string, client *http.Client) *APIClient {
	return &APIClient{
		baseURL:    DefaultBaseURL,
		httpClient: client,
		token:      apiKey,
	}
}

func NewAPIClientWithBaseURL(apiKey string, client *http.Client, baseURL string) *APIClient {
	return &APIClient{
		baseURL:    baseURL,
		httpClient: client,
		token:      apiKey,
	}
}

func (c *APIClient) GetBaseURL() string {
	return c.baseURL
}

func (c *APIClient) do(ctx context.Context, method, endpoint string, body []byte,
	params map[string]string) (*http.Response, error) {
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
