package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultBaseURL          = "https://api.vercel.com/v3"
	defaultHTTPTimeout      = 60 * time.Second
	maxHTTPErrorBodyLength  = 1024
	truncatedHTTPBodyMarker = "...(truncated)"
)

var (
	errEmptyReq                       = errors.New("empty req")
	errInvalidCreateAuthTokenResponse = errors.New("invalid create auth token response")
	errInvalidDeleteAuthTokenResponse = errors.New("invalid delete auth token response")
	errMissingTokenID                 = errors.New("missing token id")
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

type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http error %d with response body %q", e.StatusCode, e.Body)
}

func NewAPIClient(apiKey string, client *http.Client) *APIClient {
	return &APIClient{
		baseURL:    DefaultBaseURL,
		httpClient: configuredHTTPClient(client),
		token:      apiKey,
	}
}

func NewAPIClientWithBaseURL(apiKey string, client *http.Client, baseURL string) *APIClient {
	return &APIClient{
		baseURL:    baseURL,
		httpClient: configuredHTTPClient(client),
		token:      apiKey,
	}
}

func (c *APIClient) GetBaseURL() string {
	return c.baseURL
}

func (c *APIClient) do(ctx context.Context, method, endpoint string, body []byte,
	params map[string]string) (*http.Response, error) {
	u, err := c.requestURL(endpoint, params)
	if err != nil {
		return nil, err
	}

	bearer := fmt.Sprintf("Bearer %s", c.token)

	req, err := http.NewRequestWithContext(ctx, method, u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer)

	httpClient := c.httpClient
	if httpClient == nil {
		httpClient = configuredHTTPClient(nil)
	}

	return httpClient.Do(req)
}

func (c *APIClient) requestURL(endpoint string, params map[string]string) (string, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}

	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("invalid base url %q", c.baseURL)
	}

	basePath := strings.TrimRight(u.EscapedPath(), "/")
	endpointPath := strings.TrimLeft(endpoint, "/")

	rawPath := endpointPath
	if basePath != "" {
		rawPath = fmt.Sprintf("%s/%s", basePath, endpointPath)
	}

	if rawPath == "" {
		rawPath = "/"
	}

	decodedPath, err := url.PathUnescape(rawPath)
	if err != nil {
		return "", err
	}

	u.Path = decodedPath
	u.RawPath = rawPath

	q := u.Query()
	for key, val := range params {
		q.Set(key, val)
	}

	u.RawQuery = q.Encode()

	return u.String(), nil
}

func configuredHTTPClient(client *http.Client) *http.Client {
	if client == nil {
		return &http.Client{Timeout: defaultHTTPTimeout}
	}

	if client.Timeout != 0 {
		return client
	}

	clone := *client
	clone.Timeout = defaultHTTPTimeout

	return &clone
}

func newHTTPError(statusCode int, body []byte) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Body:       sanitizeHTTPErrorBody(body),
	}
}

func sanitizeHTTPErrorBody(body []byte) string {
	var payload any
	if err := json.Unmarshal(body, &payload); err == nil {
		redacted := redactJSON(payload)

		b, marshalErr := json.Marshal(redacted)
		if marshalErr == nil {
			body = b
		}
	}

	if len(body) <= maxHTTPErrorBodyLength {
		return string(body)
	}

	return string(body[:maxHTTPErrorBodyLength]) + truncatedHTTPBodyMarker
}

func redactJSON(value any) any {
	switch v := value.(type) {
	case map[string]any:
		for key, child := range v {
			if isSensitiveKey(key) {
				v[key] = "[REDACTED]"

				continue
			}

			v[key] = redactJSON(child)
		}

		return v
	case []any:
		for i, child := range v {
			v[i] = redactJSON(child)
		}

		return v
	default:
		return value
	}
}

func isSensitiveKey(key string) bool {
	switch strings.ToLower(key) {
	case "authorization", "apikey", "api_key", "bearertoken", "secret", "token":
		return true
	default:
		return false
	}
}
