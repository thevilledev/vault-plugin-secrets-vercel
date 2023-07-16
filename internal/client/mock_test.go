package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMock_New(t *testing.T) {
	t.Parallel()

	m := NewMockClient()
	require.NotNil(t, m)
}

func TestMock_BaseURL(t *testing.T) {
	t.Parallel()

	m := NewMockClient()
	require.NotNil(t, m)

	u := m.GetBaseURL()
	require.Empty(t, u)
}

func TestMock_CreateToken(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		req      *CreateAuthTokenRequest
		expError string
	}{
		"empty token": {
			req:      &CreateAuthTokenRequest{},
			expError: "empty name for token",
		},
		"force fail": {
			req: &CreateAuthTokenRequest{
				Name:   "foo",
				TeamID: "force-fail",
			},
			expError: "force fail",
		},
		"success with just the name": {
			req: &CreateAuthTokenRequest{
				Name: "foo",
			},
		},
		"success with name + team id": {
			req: &CreateAuthTokenRequest{
				Name:   "foo",
				TeamID: "bar",
			},
		},
	}
	for name, tc := range cases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			m := NewMockClient()
			require.NotNil(t, m)

			r, err := m.CreateAuthToken(ctx, tc.req)
			if tc.expError != "" {
				require.EqualError(t, err, tc.expError)
			} else {
				require.NotNil(t, r)
			}
		})
	}
}

func TestMock_DeleteToken(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		req      *DeleteAuthTokenRequest
		expError string
	}{
		"empty id": {
			req:      &DeleteAuthTokenRequest{},
			expError: "empty id for token",
		},
		"success with just the name": {
			req: &DeleteAuthTokenRequest{
				ID: "foo",
			},
		},
	}
	for name, tc := range cases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			m := NewMockClient()
			require.NotNil(t, m)

			r, err := m.DeleteAuthToken(ctx, tc.req)
			if tc.expError != "" {
				require.EqualError(t, err, tc.expError)
			} else {
				require.NotNil(t, r)
			}
		})
	}
}
