package hcloud

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/samber/lo"
	"github.com/samber/mo"
)

// EnableRescue enables rescue mode on the given server.
func (c *Client) EnableRescue(
	ctx context.Context,
	server *hcloud.Server,
	sshKeyIDs []int64,
) (mo.Result[hcloud.ServerEnableRescueResult], error) {
	sshKeys := lo.Map(sshKeyIDs, func(id int64, _ int) *hcloud.SSHKey {
		var sshKey hcloud.SSHKey
		sshKey.ID = id
		return &sshKey
	})

	var rescueOpts hcloud.ServerEnableRescueOpts
	rescueOpts.Type = hcloud.ServerRescueTypeLinux64
	rescueOpts.SSHKeys = sshKeys

	result, _, err := c.api.Server.EnableRescue(ctx, server, rescueOpts)
	if err != nil {
		return mo.Err[hcloud.ServerEnableRescueResult](err), err
	}

	slog.Info("rescue enabled", "server_id", server.ID, "action_id", result.Action.ID)
	return mo.Ok(result), nil
}

// DisableRescue disables rescue mode on the given server.
func (c *Client) DisableRescue(ctx context.Context, server *hcloud.Server) error {
	_, _, err := c.api.Server.DisableRescue(ctx, server)
	if err != nil {
		return fmt.Errorf("disable rescue for server %d: %w", server.ID, err)
	}
	return nil
}
