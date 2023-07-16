package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thevilledev/vault-plugin-secrets-vercel/internal/client"
)

func TestService_New(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		apiKey  string
		baseURL string
		expURL  string
	}{
		"api client": {
			apiKey: "12345asbd",
			expURL: client.DefaultBaseURL,
		},
		"api client with custom base url": {
			apiKey:  "abbaacdc",
			baseURL: "http://somethingelse",
			expURL:  "http://somethingelse",
		},
		"mock client": {
			apiKey: "mock",
			expURL: "",
		},
		"mock client with custom base url": {
			apiKey:  "mock",
			baseURL: "http://somethingelse",
			expURL:  "",
		},
	}
	for name, tc := range cases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var s *Service
			if tc.baseURL != "" {
				s = NewWithBaseURL(tc.apiKey, tc.baseURL)
			} else {
				s = New(tc.apiKey)
			}
			require.Equal(t, s.client.GetBaseURL(), tc.expURL)
		})
	}
}

func TestService_CreateToken(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		name     string
		ttl      int64
		teamID   string
		expError string
	}{
		"success": {
			name:   "foo",
			ttl:    10,
			teamID: "asd",
		},
		"negative ttl": {
			name:     "foo",
			ttl:      -10,
			teamID:   "asd",
			expError: "invalid ttl",
		},
		"failed token creation": {
			name:     "",
			ttl:      3600,
			expError: "empty name for token",
		},
	}
	for name, tc := range cases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			s := New("mock")
			tid, bt, err := s.CreateAuthToken(ctx, tc.name, tc.ttl, tc.teamID)
			if tc.expError != "" {
				require.EqualError(t, err, tc.expError)
			} else {
				require.Equal(t, bt, "some-bearer-token")
				require.NotEmpty(t, tid)
			}
		})
	}
}

func TestService_DeleteToken(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		id          string
		createToken bool
		expError    string
	}{
		"success": {
			id:          "foo",
			createToken: true,
		},
		"failed token deletion": {
			id:          "",
			createToken: false,
			expError:    "empty id for token",
		},
	}
	for name, tc := range cases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			s := New("mock")

			var id string
			if tc.createToken {
				i, _, err := s.CreateAuthToken(ctx, tc.id+"foobar", 1, "foo")
				require.NoError(t, err)
				id = i
			}
			deletedID, err := s.DeleteAuthToken(ctx, id)
			if tc.expError != "" {
				require.EqualError(t, err, tc.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, id, deletedID)
			}
		})
	}
}
