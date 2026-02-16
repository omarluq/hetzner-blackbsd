package hcloud_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	hcloudsdk "github.com/hetznercloud/hcloud-go/v2/hcloud"
	bsdhcloud "github.com/omarluq/hetzner-blackbsd/internal/hcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListServers(t *testing.T) {
	t.Parallel()

	t.Run("returns servers with correct label", func(t *testing.T) {
		t.Parallel()

		serverJSON := `{
			"servers": [
				{"id": 1, "name": "server-1", "status": "running"},
				{"id": 2, "name": "server-2", "status": "running"}
			]
		}`

		testServer := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, _ *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusOK)
				writeJSON(t, writer, serverJSON)
			}))
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))
		servers, err := client.ListServers(context.Background())

		require.NoError(t, err)
		assert.Len(t, servers, 2)
		assert.Equal(t, int64(1), servers[0].ID)
		assert.Equal(t, int64(2), servers[1].ID)
	})
}

func TestDeleteServer(t *testing.T) {
	t.Parallel()

	newTestServer := func(statusCode int, responseBody string) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, _ *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(statusCode)
				writeJSON(t, writer, responseBody)
			}))
	}

	var srv hcloudsdk.Server
	srv.ID = 42

	t.Run("deletes existing server and returns true", func(t *testing.T) {
		t.Parallel()

		testServer := newTestServer(http.StatusOK, `{}`)
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))
		deleted, err := client.DeleteServer(context.Background(), &srv)

		require.NoError(t, err)
		assert.True(t, deleted)
	})

	t.Run("returns false for 404", func(t *testing.T) {
		t.Parallel()

		testServer := newTestServer(http.StatusNotFound, `{"error": {"code": "not_found"}}`)
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))
		deleted, err := client.DeleteServer(context.Background(), &srv)

		require.NoError(t, err)
		assert.False(t, deleted)
	})
}

func TestGetServer(t *testing.T) {
	t.Parallel()

	t.Run("fetches server by ID", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, r *http.Request) {
				if strings.Contains(r.URL.String(), "/servers/42") {
					writer.Header().Set("Content-Type", "application/json")
					writer.WriteHeader(http.StatusOK)
					writeJSON(t, writer, `{
						"server": {"id": 42, "name": "test-server", "status": "running"}
					}`)
				}
			}))
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))
		result := client.GetServer(context.Background(), 42)

		assert.True(t, result.IsPresent())
		srv, _ := result.Get()
		assert.Equal(t, int64(42), srv.ID)
	})

	t.Run("returns None for 404", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, _ *http.Request) {
				writer.WriteHeader(http.StatusNotFound)
			}))
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))
		result := client.GetServer(context.Background(), 42)

		assert.False(t, result.IsPresent())
	})
}

func TestServerStatus(t *testing.T) {
	t.Parallel()

	t.Run("returns status string for existing server", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, _ *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusOK)
				writeJSON(t, writer, `{
					"server": {"id": 42, "status": "running"}
				}`)
			}))
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))
		status := client.ServerStatus(context.Background(), 42)

		assert.Equal(t, "running", status)
	})

	t.Run("returns unknown for missing server", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, _ *http.Request) {
				writer.WriteHeader(http.StatusNotFound)
			}))
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))
		status := client.ServerStatus(context.Background(), 42)

		assert.Equal(t, "unknown", status)
	})
}

func TestCreateServer(t *testing.T) {
	t.Parallel()

	t.Run("creates server with correct options", func(t *testing.T) {
		t.Parallel()

		var requestBodies []string
		testServer := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.Contains(r.URL.String(), "/servers") {
					body, readErr := io.ReadAll(r.Body)
					require.NoError(t, readErr)
					requestBodies = append(requestBodies, string(body))

					writer.Header().Set("Content-Type", "application/json")
					writer.WriteHeader(http.StatusCreated)
					writeJSON(t, writer, `{
						"server": {
							"id": 42,
							"name": "test-server",
							"status": "running"
						},
						"action": {"id": 1, "status": "running"}
					}`)
				}
			}))
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))
		opts := &bsdhcloud.CreateOpts{
			Name:       "test-server",
			ServerType: "cpx31",
			Image:      "ubuntu-24.04",
			Location:   "fsn1",
			SSHKeyIDs:  []int64{123},
		}

		result, err := client.CreateServer(context.Background(), opts)

		require.NoError(t, err)
		assert.Equal(t, int64(42), result.ID)
		assert.Equal(t, "test-server", result.Name)

		require.Len(t, requestBodies, 1)
		assert.Contains(t, requestBodies[0], "test-server")
		assert.Contains(t, requestBodies[0], "cpx31")
		assert.Contains(t, requestBodies[0], "ubuntu-24.04")
		assert.Contains(t, requestBodies[0], "fsn1")
		assert.Contains(t, requestBodies[0], "managed-by")
		assert.Contains(t, requestBodies[0], "blackbsd-builder")
	})
}
