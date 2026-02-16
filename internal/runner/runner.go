// Package runner provides a shared interface for command execution on remote hosts.
package runner

import (
	"context"

	"github.com/omarluq/hetzner-blackbsd/internal/ssh"
)

// Runner executes commands on a remote host.
type Runner interface {
	Exec(ctx context.Context, command string) (ssh.CommandResult, error)
}
