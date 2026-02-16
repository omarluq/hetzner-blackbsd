// Package hcloud provides a high-level wrapper around the official Hetzner Cloud client.
package hcloud

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cenkalti/backoff/v4"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/samber/lo"
	"github.com/samber/mo"
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

// ListServers returns all servers matching the blackbsd label.
func (c *Client) ListServers(ctx context.Context) ([]*hcloud.Server, error) {
	var listOpts hcloud.ListOpts
	listOpts.LabelSelector = Label

	var serverListOpts hcloud.ServerListOpts
	serverListOpts.ListOpts = listOpts

	servers, err := c.api.Server.AllWithOpts(ctx, serverListOpts)
	if err != nil {
		return nil, fmt.Errorf("list servers: %w", err)
	}
	return servers, nil
}

// GetServer fetches a server by ID, returning None if not found.
func (c *Client) GetServer(ctx context.Context, id int64) mo.Option[*hcloud.Server] {
	server, _, err := c.api.Server.GetByID(ctx, id)
	if err != nil || server == nil {
		return mo.None[*hcloud.Server]()
	}
	return mo.Some(server)
}

// DeleteServer deletes a server. Returns true if deleted, false if not found.
func (c *Client) DeleteServer(
	ctx context.Context,
	server *hcloud.Server,
) (bool, error) {
	_, _, err := c.api.Server.DeleteWithResult(ctx, server)
	if err != nil {
		if hcloud.IsError(err, hcloud.ErrorCodeNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("delete server %d: %w", server.ID, err)
	}

	slog.Info("server deleted", "id", server.ID, "name", server.Name)
	return true, nil
}

// CreateServer provisions a new build server with the blackbsd label.
func (c *Client) CreateServer(
	ctx context.Context,
	opts *CreateOpts,
) (*hcloud.Server, error) {
	var result hcloud.ServerCreateResult

	var serverType hcloud.ServerType
	serverType.Name = opts.ServerType

	var image hcloud.Image
	image.Name = opts.Image

	var location hcloud.Location
	location.Name = opts.Location

	sshKeys := lo.Map(opts.SSHKeyIDs, func(id int64, _ int) *hcloud.SSHKey {
		var sshKey hcloud.SSHKey
		sshKey.ID = id
		return &sshKey
	})

	var createOpts hcloud.ServerCreateOpts
	createOpts.Name = opts.Name
	createOpts.ServerType = &serverType
	createOpts.Image = &image
	createOpts.Location = &location
	createOpts.SSHKeys = sshKeys
	createOpts.Labels = map[string]string{LabelKey: LabelValue}

	retryOperation := func() error {
		var err error
		result, _, err = c.api.Server.Create(ctx, createOpts)
		return err
	}

	if err := backoff.Retry(retryOperation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)); err != nil {
		return nil, fmt.Errorf("create server %s: %w", opts.Name, err)
	}

	slog.Info("server created", "id", result.Server.ID, "name", result.Server.Name)
	return result.Server, nil
}

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

// ServerStatus returns the status of a server by ID.
func (c *Client) ServerStatus(ctx context.Context, id int64) string {
	opt := c.GetServer(ctx, id)
	if server, ok := opt.Get(); ok {
		return string(server.Status)
	}
	return "unknown"
}

// CreateOpts defines options for creating a build server.
type CreateOpts struct {
	Name       string
	ServerType string
	Image      string
	Location   string
	SSHKeyIDs  []int64
}
