package hcloud

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const (
	// waitForStatusTimeout is the maximum time to wait for a server status change.
	waitForStatusTimeout = 10 * time.Minute
)

// WaitForAction waits until the given action completes.
func (c *Client) WaitForAction(ctx context.Context, action *hcloud.Action) error {
	if err := c.api.Action.WaitFor(ctx, action); err != nil {
		return fmt.Errorf("wait for action %d: %w", action.ID, err)
	}

	slog.Info("action completed", "id", action.ID)
	return nil
}

// WaitForServerStatus polls until the server reaches the target status.
func (c *Client) WaitForServerStatus(
	ctx context.Context,
	serverID int64,
	target hcloud.ServerStatus,
) error {
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = waitForStatusTimeout

	operation := func() error {
		server, _, getErr := c.api.Server.GetByID(ctx, serverID)
		if getErr != nil {
			return backoff.Permanent(fmt.Errorf("get server %d: %w", serverID, getErr))
		}

		if server == nil {
			return backoff.Permanent(fmt.Errorf("server %d: not found", serverID))
		}

		if server.Status != target {
			return fmt.Errorf(
				"server %d: status %s, want %s",
				serverID, server.Status, target,
			)
		}

		return nil
	}

	if err := backoff.Retry(operation, backoff.WithContext(expBackoff, ctx)); err != nil {
		return fmt.Errorf("wait for server %d status %s: %w", serverID, target, err)
	}

	slog.Info("server reached target status", "id", serverID, "status", target)
	return nil
}
