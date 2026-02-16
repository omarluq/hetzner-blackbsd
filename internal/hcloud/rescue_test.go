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

func TestEnableRescue(t *testing.T) {
	t.Parallel()

	t.Run("enables rescue mode successfully", func(t *testing.T) {
		t.Parallel()

		body := `{"action": {"id": 1, "status": "running"}, "root_password": "test"}`

		testServer := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, _ *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusOK)
				writeJSON(t, writer, body)
			}))
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))
		var srv hcloudsdk.Server
		srv.ID = 42

		result, err := client.EnableRescue(context.Background(), &srv, []int64{123})

		require.NoError(t, err)
		assert.True(t, result.IsOk())
		rescueResult, getErr := result.Get()
		require.NoError(t, getErr)
		assert.Equal(t, int64(1), rescueResult.Action.ID)
	})
}

func TestDisableRescue(t *testing.T) {
	t.Parallel()

	t.Run("disables rescue mode successfully", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, _ *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusOK)
				writeJSON(t, writer, `{"action": {"id": 1, "status": "running", "command": "disable_rescue"}}`)
			}))
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))
		var srv hcloudsdk.Server
		srv.ID = 42

		err := client.DisableRescue(context.Background(), &srv)
		require.NoError(t, err)
	})
}
