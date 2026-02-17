package customize

import (
	"context"
	"fmt"
	"strings"

	"github.com/omarluq/hetzner-blackbsd/internal/config"
	"github.com/omarluq/hetzner-blackbsd/internal/ssh"
)

// ApplyBranding sets hostname, MOTD, and creates the default user.
func (c *Customizer) ApplyBranding(ctx context.Context, branding config.Branding) error {
	if err := c.setHostname(ctx, branding.Hostname); err != nil {
		return err
	}

	if err := c.writeMOTD(ctx, branding.MOTD); err != nil {
		return err
	}

	return c.createUser(ctx, branding.DefaultUser)
}

func (c *Customizer) setHostname(ctx context.Context, hostname string) error {
	command := fmt.Sprintf("echo hostname=%s >> /etc/rc.conf", ssh.EscapeShellArg(hostname))

	result, execErr := c.runner.Exec(ctx, command)
	if execErr != nil {
		return fmt.Errorf("set hostname: %w", execErr)
	}

	if !result.Success() {
		return fmt.Errorf("set hostname: exited %d: %s",
			result.ExitCode, strings.TrimSpace(result.Stderr))
	}

	return nil
}

func (c *Customizer) writeMOTD(ctx context.Context, motd string) error {
	// Escape for printf - percent signs and backslashes
	escaped := strings.ReplaceAll(motd, "%", "%%")
	escaped = strings.ReplaceAll(escaped, "\\", "\\\\")

	command := fmt.Sprintf("printf %%s %s > /etc/motd", ssh.EscapeShellArg(escaped))

	result, execErr := c.runner.Exec(ctx, command)
	if execErr != nil {
		return fmt.Errorf("write motd: %w", execErr)
	}

	if !result.Success() {
		return fmt.Errorf("write motd: exited %d: %s",
			result.ExitCode, strings.TrimSpace(result.Stderr))
	}

	return nil
}

func (c *Customizer) createUser(ctx context.Context, username string) error {
	command := fmt.Sprintf("useradd -m -G wheel %s", ssh.EscapeShellArg(username))

	result, execErr := c.runner.Exec(ctx, command)
	if execErr != nil {
		return fmt.Errorf("create user %s: %w", username, execErr)
	}

	// Exit code 0 means created; code 9 means "already exists" which is acceptable.
	const userExistsCode = 9
	if !result.Success() && result.ExitCode != userExistsCode {
		return fmt.Errorf("create user %s: exited %d: %s",
			username, result.ExitCode, strings.TrimSpace(result.Stderr))
	}

	return nil
}
