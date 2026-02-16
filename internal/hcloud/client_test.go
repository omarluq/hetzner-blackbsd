package hcloud_test

import (
	"testing"

	hcloudsdk "github.com/hetznercloud/hcloud-go/v2/hcloud"
	bsdhcloud "github.com/omarluq/hetzner-blackbsd/internal/hcloud"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	client := bsdhcloud.NewClient("test_token")

	assert.NotNil(t, client)
}

func TestNewClientWithOpts(t *testing.T) {
	t.Parallel()

	testServer := bsdhcloud.NewClientWithOpts(hcloudsdk.WithEndpoint("http://example.com"))

	assert.NotNil(t, testServer)
}
