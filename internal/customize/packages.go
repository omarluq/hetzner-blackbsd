// Package customize applies post-install customizations to a NetBSD host.
package customize

import (
	"context"
	"fmt"
	"strings"
)

// InstallPackages installs each package individually via pkg_add for error isolation.
func (c *Customizer) InstallPackages(ctx context.Context, packages []string) error {
	for _, packageName := range packages {
		command := fmt.Sprintf("pkg_add -v %s", packageName)

		result, execErr := c.runner.Exec(ctx, command)
		if execErr != nil {
			return fmt.Errorf("install package %s: %w", packageName, execErr)
		}

		if !result.Success() {
			return fmt.Errorf("install package %s: exited %d: %s",
				packageName, result.ExitCode, strings.TrimSpace(result.Stderr))
		}
	}

	return nil
}

// DefaultSecurityTools returns the default list of security packages to install.
func DefaultSecurityTools() []string {
	return []string{
		"nmap",
		"wireshark",
		"metasploit",
		"aircrack-ng",
		"snort",
		"hydra",
		"john",
		"tcpdump",
		"netcat",
		"socat",
	}
}
