package customize

import (
	"context"
	"fmt"
	"strings"
)

// ConfigureNetworking enables DHCP on the first interface and writes resolv.conf.
func (c *Customizer) ConfigureNetworking(ctx context.Context) error {
	if err := c.enableDHCP(ctx); err != nil {
		return err
	}

	return c.writeResolvConf(ctx)
}

func (c *Customizer) enableDHCP(ctx context.Context) error {
	command := `echo "dhcpcd=YES" >> /etc/rc.conf`

	result, execErr := c.runner.Exec(ctx, command)
	if execErr != nil {
		return fmt.Errorf("enable dhcp: %w", execErr)
	}

	if !result.Success() {
		return fmt.Errorf("enable dhcp: exited %d: %s",
			result.ExitCode, strings.TrimSpace(result.Stderr))
	}

	return nil
}

func (c *Customizer) writeResolvConf(ctx context.Context) error {
	command := `cat > /etc/resolv.conf << 'RESOLVEOF'
nameserver 1.1.1.1
nameserver 8.8.8.8
RESOLVEOF`

	result, execErr := c.runner.Exec(ctx, command)
	if execErr != nil {
		return fmt.Errorf("write resolv.conf: %w", execErr)
	}

	if !result.Success() {
		return fmt.Errorf("write resolv.conf: exited %d: %s",
			result.ExitCode, strings.TrimSpace(result.Stderr))
	}

	return nil
}
