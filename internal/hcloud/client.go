// Package hcloud provides a high-level wrapper around the official Hetzner Cloud client.
package hcloud

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const (
	// LabelKey is the label key applied to all managed servers.
	LabelKey = "managed-by"

	// LabelValue is the label value for blackbsd-managed servers.
	LabelValue = "blackbsd-builder"

	// Label is the full label selector string.
	Label = LabelKey + "=" + LabelValue
)

// Client wraps the official hcloud.Client with domain-specific operations.
type Client struct {
	api *hcloud.Client
}

// NewClient creates a new Hetzner client with the given token.
func NewClient(token string) *Client {
	return &Client{
		api: hcloud.NewClient(hcloud.WithToken(token)),
	}
}

// NewClientWithOpts creates a new Hetzner client with custom options.
func NewClientWithOpts(opts ...hcloud.ClientOption) *Client {
	return &Client{
		api: hcloud.NewClient(opts...),
	}
}
