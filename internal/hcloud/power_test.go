package hcloud_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	hcloudsdk "github.com/hetznercloud/hcloud-go/v2/hcloud"
	bsdhcloud "github.com/omarluq/hetzner-blackbsd/internal/hcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerPowerActions(t *testing.T) {
	t.Parallel()

	type serverActionCase struct {
		name       string
		callAction func(*bsdhcloud.Client, context.Context, *hcloudsdk.Server) error
		errorText  string
	}

	cases := []serverActionCase{
		{
			name: "reset",
			callAction: func(client *bsdhcloud.Client, ctx context.Context, srv *hcloudsdk.Server) error {
				return client.ResetServer(ctx, srv)
			},
			errorText: "reset server 42",
		},
		{
			name: "power on",
			callAction: func(client *bsdhcloud.Client, ctx context.Context, srv *hcloudsdk.Server) error {
				return client.PowerOnServer(ctx, srv)
			},
			errorText: "power on server 42",
		},
		{
			name: "power off",
			callAction: func(client *bsdhcloud.Client, ctx context.Context, srv *hcloudsdk.Server) error {
				return client.PowerOffServer(ctx, srv)
			},
			errorText: "power off server 42",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name+" succeeds", func(t *testing.T) {
			t.Parallel()

			testServer := httptest.NewServer(http.HandlerFunc(
				func(writer http.ResponseWriter, _ *http.Request) {
					writer.Header().Set("Content-Type", "application/json")
					writeJSON(t, writer, `{"action": {"id": 1, "status": "running"}}`)
				}))
			defer testServer.Close()

			client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))

			var srv hcloudsdk.Server
			srv.ID = 42

			err := testCase.callAction(client, context.Background(), &srv)
			require.NoError(t, err)
		})

		t.Run(testCase.name+" returns error on failure", func(t *testing.T) {
			t.Parallel()

			testServer := httptest.NewServer(http.HandlerFunc(
				func(writer http.ResponseWriter, _ *http.Request) {
					writer.Header().Set("Content-Type", "application/json")
					writer.WriteHeader(http.StatusInternalServerError)
					writeJSON(t, writer, `{"error": {"code": "server_error", "message": "fail"}}`)
				}))
			defer testServer.Close()

			client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))

			var srv hcloudsdk.Server
			srv.ID = 42

			err := testCase.callAction(client, context.Background(), &srv)
			require.Error(t, err)
			assert.Contains(t, err.Error(), testCase.errorText)
		})
	}
}
