package client

import (
	"bytes"
	"context"
	"fmt"
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

func New(apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: httpTimeout,
		},
		token: apiKey,
	}
}

func (c *Client) do(ctx context.Context, method, endpoint string, body []byte,
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
