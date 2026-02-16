package hcloud

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// serverActionFunc is a function that performs a server power action.
type serverActionFunc func(context.Context, *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)

// runServerAction executes a server power action, logs success, and wraps errors.
func (c *Client) runServerAction(
	ctx context.Context,
	server *hcloud.Server,
	actionFn serverActionFunc,
	actionName string,
) error {
	_, _, err := actionFn(ctx, server)
	if err != nil {
		return fmt.Errorf("%s server %d: %w", actionName, server.ID, err)
	}

	slog.Info(actionName+" completed", "id", server.ID, "name", server.Name)
	return nil
}

// ResetServer performs a hardware reset on the given server.
func (c *Client) ResetServer(ctx context.Context, server *hcloud.Server) error {
	return c.runServerAction(ctx, server, c.api.Server.Reset, "reset")
}

// PowerOnServer powers on the given server.
func (c *Client) PowerOnServer(ctx context.Context, server *hcloud.Server) error {
	return c.runServerAction(ctx, server, c.api.Server.Poweron, "power on")
}

// PowerOffServer powers off the given server.
func (c *Client) PowerOffServer(ctx context.Context, server *hcloud.Server) error {
	return c.runServerAction(ctx, server, c.api.Server.Poweroff, "power off")
}
