package hcloud_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	hcloudsdk "github.com/hetznercloud/hcloud-go/v2/hcloud"
	bsdhcloud "github.com/omarluq/hetzner-blackbsd/internal/hcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWaitForAction(t *testing.T) {
	t.Parallel()

	t.Run("completes for finished action", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set("Content-Type", "application/json")

				if strings.Contains(request.URL.Path, "/actions") {
					writeJSON(t, writer, `{
						"actions": [{"id": 1, "status": "success", "progress": 100}]
					}`)
				}
			}))
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))

		var action hcloudsdk.Action
		action.ID = 1
		action.Status = hcloudsdk.ActionStatusSuccess

		err := client.WaitForAction(context.Background(), &action)
		require.NoError(t, err)
	})

	t.Run("returns error for failed action", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, _ *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				writeJSON(t, writer, `{
					"actions": [{"id": 1, "status": "error", "error": {"code": "action_failed", "message": "boom"}}]
				}`)
			}))
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))

		var action hcloudsdk.Action
		action.ID = 1
		action.Status = hcloudsdk.ActionStatusError
		action.ErrorCode = "action_failed"
		action.ErrorMessage = "boom"

		err := client.WaitForAction(context.Background(), &action)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "wait for action 1")
	})
}

func TestWaitForServerStatus(t *testing.T) {
	t.Parallel()

	t.Run("returns when status matches", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, _ *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				writeJSON(t, writer, `{"server": {"id": 42, "status": "running"}}`)
			}))
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))

		err := client.WaitForServerStatus(
			context.Background(), 42, hcloudsdk.ServerStatusRunning,
		)
		require.NoError(t, err)
	})

	t.Run("returns error for missing server", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, _ *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusNotFound)
				writeJSON(t, writer, `{"error": {"code": "not_found", "message": "not found"}}`)
			}))
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))

		err := client.WaitForServerStatus(
			context.Background(), 99, hcloudsdk.ServerStatusRunning,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "99")
	})
}
