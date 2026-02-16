package hcloud_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	hcloudsdk "github.com/hetznercloud/hcloud-go/v2/hcloud"
	bsdhcloud "github.com/omarluq/hetzner-blackbsd/internal/hcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureSSHKey(t *testing.T) {
	t.Parallel()

	t.Run("returns existing key", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set("Content-Type", "application/json")

				if request.Method == http.MethodGet && strings.Contains(request.URL.Path, "/ssh_keys") {
					writeJSON(t, writer, `{
						"ssh_keys": [{"id": 1, "name": "test-key", "public_key": "ssh-rsa AAA"}]
					}`)
				}
			}))
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))
		sshKey, err := client.EnsureSSHKey(context.Background(), "test-key", "ssh-rsa AAA")

		require.NoError(t, err)
		require.NotNil(t, sshKey)
		assert.Equal(t, int64(1), sshKey.ID)
		assert.Equal(t, "test-key", sshKey.Name)
	})

	t.Run("creates key when not found", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set("Content-Type", "application/json")

				switch request.Method {
				case http.MethodGet:
					writeJSON(t, writer, `{"ssh_keys": []}`)
				case http.MethodPost:
					writer.WriteHeader(http.StatusCreated)
					writeJSON(t, writer, `{
						"ssh_key": {"id": 2, "name": "new-key", "public_key": "ssh-rsa BBB"}
					}`)
				}
			}))
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))
		sshKey, err := client.EnsureSSHKey(context.Background(), "new-key", "ssh-rsa BBB")

		require.NoError(t, err)
		require.NotNil(t, sshKey)
		assert.Equal(t, int64(2), sshKey.ID)
		assert.Equal(t, "new-key", sshKey.Name)
	})
}

func TestFindSSHKeyByFingerprint(t *testing.T) {
	t.Parallel()

	t.Run("finds key by fingerprint", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set("Content-Type", "application/json")

				params, parseErr := url.ParseQuery(request.URL.RawQuery)
				require.NoError(t, parseErr)

				if params.Get("fingerprint") != "" {
					writeJSON(t, writer, `{
						"ssh_keys": [{"id": 1, "name": "my-key", "fingerprint": "ab:cd:ef"}]
					}`)
				}
			}))
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))
		result := client.FindSSHKeyByFingerprint(context.Background(), "ab:cd:ef")

		assert.True(t, result.IsPresent())
		sshKey, present := result.Get()
		assert.True(t, present)
		assert.Equal(t, int64(1), sshKey.ID)
	})

	t.Run("returns None when not found", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, _ *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				writeJSON(t, writer, `{"ssh_keys": []}`)
			}))
		defer testServer.Close()

		client := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint(testServer.URL))
		result := client.FindSSHKeyByFingerprint(context.Background(), "no:such:fp")

		assert.False(t, result.IsPresent())
	})
}
