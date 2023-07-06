package client

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

func recordHelper(t *testing.T, fixture string, f func(context.Context, *testing.T, *recorder.Recorder, *Client)) {
	t.Helper()

	r, err := recorder.New("fixtures/" + fixture)
	require.NoError(t, err)

	hook := func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")

		if strings.Contains(i.Request.URL, "/user/tokens") && i.Request.Method == http.MethodPost {
			var c map[string]any
			if e := json.Unmarshal([]byte(i.Response.Body), &c); e != nil {
				return e
			}

			_, ok := c["bearerToken"]
			if ok {
				c["bearerToken"] = "REDACTED"
			}

			res, e := json.Marshal(c)
			if e != nil {
				return e
			}

			i.Response.Body = string(res)
		}

		return nil
	}

	r.AddHook(hook, recorder.BeforeSaveHook)

	defer func() {
		err = r.Stop()
		require.NoError(t, err)
	}()

	// Required when updating fixtures
	apiKey := os.Getenv("VERCEL_TOKEN")

	httpClient := r.GetDefaultClient()

	ctx, cancel := context.WithTimeout(context.Background(), defaultHTTPTimeout)
	defer cancel()

	f(ctx, t, r, NewWithBaseURL(apiKey, httpClient, DefaultBaseURL))
}
