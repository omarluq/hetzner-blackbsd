package hcloud

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/samber/mo"
)

// RescueOption configures rescue mode enablement.
type RescueOption func(*hcloud.ServerEnableRescueOpts)

// WithRescueType sets the rescue type (defaults to Linux64).
func WithRescueType(t hcloud.ServerRescueType) RescueOption {
	return func(opts *hcloud.ServerEnableRescueOpts) {
		opts.Type = t
	}
}

// WithRescueSSHKeys adds SSH keys to the rescue environment.
func WithRescueSSHKeys(keys []int64) RescueOption {
	return func(opts *hcloud.ServerEnableRescueOpts) {
		sshKeys := make([]*hcloud.SSHKey, len(keys))
		for i, id := range keys {
			var k hcloud.SSHKey
			k.ID = id
			sshKeys[i] = &k
		}
		opts.SSHKeys = sshKeys
	}
}

// EnableRescue enables rescue mode on the given server.
func (c *Client) EnableRescue(
	ctx context.Context,
	server *hcloud.Server,
	opts ...RescueOption,
) mo.Result[hcloud.ServerEnableRescueResult] {
	rescueOpts := &hcloud.ServerEnableRescueOpts{
		Type:   hcloud.ServerRescueTypeLinux64,
		SSHKeys: nil,
	}

	for _, opt := range opts {
		opt(rescueOpts)
	}

	result, _, err := c.api.Server.EnableRescue(ctx, server, *rescueOpts)
	if err != nil {
		slog.Error("enable rescue failed", "server_id", server.ID, "error", err)
		return mo.Err[hcloud.ServerEnableRescueResult](err)
	}

	slog.Info("rescue enabled", "server_id", server.ID, "action_id", result.Action.ID)
	return mo.Ok(result)
}

// DisableRescue disables rescue mode on the given server.
func (c *Client) DisableRescue(ctx context.Context, server *hcloud.Server) error {
	_, _, err := c.api.Server.DisableRescue(ctx, server)
	if err != nil {
		return fmt.Errorf("disable rescue for server %d: %w", server.ID, err)
	}
	return nil
}
