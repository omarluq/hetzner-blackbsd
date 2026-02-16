package hcloud

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/samber/mo"
)

// EnsureSSHKey finds an SSH key by name, creating it if it doesn't exist.
func (c *Client) EnsureSSHKey(
	ctx context.Context,
	name string,
	publicKey string,
) (*hcloud.SSHKey, error) {
	existing, _, err := c.api.SSHKey.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("get ssh key %q: %w", name, err)
	}

	if existing != nil {
		slog.Info("ssh key found", "name", name, "id", existing.ID)
		return existing, nil
	}

	var createOpts hcloud.SSHKeyCreateOpts
	createOpts.Name = name
	createOpts.PublicKey = publicKey

	created, _, err := c.api.SSHKey.Create(ctx, createOpts)
	if err != nil {
		return nil, fmt.Errorf("create ssh key %q: %w", name, err)
	}

	slog.Info("ssh key created", "name", name, "id", created.ID)
	return created, nil
}

// FindSSHKeyByFingerprint finds an SSH key by its fingerprint.
func (c *Client) FindSSHKeyByFingerprint(
	ctx context.Context,
	fingerprint string,
) mo.Option[*hcloud.SSHKey] {
	sshKey, _, err := c.api.SSHKey.GetByFingerprint(ctx, fingerprint)
	if err != nil || sshKey == nil {
		return mo.None[*hcloud.SSHKey]()
	}

	return mo.Some(sshKey)
}
