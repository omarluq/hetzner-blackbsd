package customize_test

import (
	"context"
	"testing"

	"github.com/omarluq/hetzner-blackbsd/internal/config"
	"github.com/omarluq/hetzner-blackbsd/internal/customize"
	"github.com/omarluq/hetzner-blackbsd/internal/ssh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRunner struct {
	err      error
	results  map[string]ssh.CommandResult
	commands []string
}

func (runner *mockRunner) Exec(_ context.Context, command string) (ssh.CommandResult, error) {
	runner.commands = append(runner.commands, command)

	if runner.err != nil {
		return ssh.CommandResult{Stdout: "", Stderr: "", ExitCode: 0}, runner.err
	}

	if result, found := runner.results[command]; found {
		return result, nil
	}

	return ssh.CommandResult{Stdout: "", Stderr: "", ExitCode: 0}, nil
}

func successResult() ssh.CommandResult {
	return ssh.CommandResult{Stdout: "", Stderr: "", ExitCode: 0}
}

func failureResult(stderr string, exitCode int) ssh.CommandResult {
	return ssh.CommandResult{Stdout: "", Stderr: stderr, ExitCode: exitCode}
}

func testBranding() config.Branding {
	return config.Branding{
		Hostname:    "blackbsd",
		MOTD:        "Welcome to BlackBSD",
		DefaultUser: "security",
	}
}

func TestInstallPackages(t *testing.T) {
	t.Parallel()

	t.Run("installs all packages successfully", func(t *testing.T) {
		t.Parallel()

		runner := &mockRunner{
			err: nil,
			results: map[string]ssh.CommandResult{
				"pkg_add -v 'nmap'":    successResult(),
				"pkg_add -v 'tcpdump'": successResult(),
			},
			commands: nil,
		}
		customizer := customize.New(runner)

		installErr := customizer.InstallPackages(context.Background(), []string{"nmap", "tcpdump"})

		require.NoError(t, installErr)
		assert.Len(t, runner.commands, 2)
		assert.Equal(t, "pkg_add -v 'nmap'", runner.commands[0])
		assert.Equal(t, "pkg_add -v 'tcpdump'", runner.commands[1])
	})

	t.Run("returns error with package name on failure", func(t *testing.T) {
		t.Parallel()

		runner := &mockRunner{
			err: nil,
			results: map[string]ssh.CommandResult{
				"pkg_add -v 'nmap'":    successResult(),
				"pkg_add -v 'badpkg'":  failureResult("package not found", 1),
				"pkg_add -v 'tcpdump'": successResult(),
			},
			commands: nil,
		}
		customizer := customize.New(runner)

		installErr := customizer.InstallPackages(context.Background(), []string{"nmap", "badpkg", "tcpdump"})

		require.Error(t, installErr)
		assert.Contains(t, installErr.Error(), "badpkg")
		assert.Len(t, runner.commands, 2)
	})

	t.Run("returns error on exec failure", func(t *testing.T) {
		t.Parallel()

		runner := &mockRunner{
			err:      assert.AnError,
			results:  nil,
			commands: nil,
		}
		customizer := customize.New(runner)

		installErr := customizer.InstallPackages(context.Background(), []string{"nmap"})

		require.Error(t, installErr)
		assert.Contains(t, installErr.Error(), "nmap")
	})

	t.Run("handles empty package list", func(t *testing.T) {
		t.Parallel()

		runner := &mockRunner{
			err:      nil,
			results:  nil,
			commands: nil,
		}
		customizer := customize.New(runner)

		installErr := customizer.InstallPackages(context.Background(), []string{})

		require.NoError(t, installErr)
		assert.Empty(t, runner.commands)
	})
}

func TestApplyBrandingCommands(t *testing.T) {
	t.Parallel()

	t.Run("executes hostname motd and user commands", func(t *testing.T) {
		t.Parallel()

		runner := &mockRunner{err: nil, results: nil, commands: nil}
		customizer := customize.New(runner)

		brandingErr := customizer.ApplyBranding(context.Background(), testBranding())

		require.NoError(t, brandingErr)
		require.Len(t, runner.commands, 3)
		assert.Contains(t, runner.commands[0], "hostname='blackbsd'")
		assert.Contains(t, runner.commands[0], "/etc/rc.conf")
		assert.Contains(t, runner.commands[1], "Welcome to BlackBSD")
		assert.Contains(t, runner.commands[1], "/etc/motd")
		assert.Contains(t, runner.commands[2], "useradd")
		assert.Contains(t, runner.commands[2], "'security'")
	})

	t.Run("returns error when hostname fails", func(t *testing.T) {
		t.Parallel()

		runner := &mockRunner{err: assert.AnError, results: nil, commands: nil}
		customizer := customize.New(runner)

		brandingErr := customizer.ApplyBranding(context.Background(), testBranding())

		require.Error(t, brandingErr)
		assert.Contains(t, brandingErr.Error(), "hostname")
	})
}

func TestApplyBrandingUserCreation(t *testing.T) {
	t.Parallel()

	t.Run("tolerates user already exists exit code", func(t *testing.T) {
		t.Parallel()

		runner := &mockRunner{
			err: nil,
			results: map[string]ssh.CommandResult{
				"useradd -m -G wheel 'security'": failureResult("user already exists", 9),
			},
			commands: nil,
		}
		customizer := customize.New(runner)

		brandingErr := customizer.ApplyBranding(context.Background(), testBranding())

		require.NoError(t, brandingErr)
	})

	t.Run("returns error on unexpected user creation failure", func(t *testing.T) {
		t.Parallel()

		runner := &mockRunner{
			err: nil,
			results: map[string]ssh.CommandResult{
				"useradd -m -G wheel 'security'": failureResult("permission denied", 1),
			},
			commands: nil,
		}
		customizer := customize.New(runner)

		brandingErr := customizer.ApplyBranding(context.Background(), testBranding())

		require.Error(t, brandingErr)
		assert.Contains(t, brandingErr.Error(), "security")
	})
}

func TestConfigureNetworking(t *testing.T) {
	t.Parallel()

	t.Run("enables dhcp and writes resolv conf", func(t *testing.T) {
		t.Parallel()

		runner := &mockRunner{err: nil, results: nil, commands: nil}
		customizer := customize.New(runner)

		networkErr := customizer.ConfigureNetworking(context.Background())

		require.NoError(t, networkErr)
		require.Len(t, runner.commands, 2)
		assert.Contains(t, runner.commands[0], "dhcpcd=YES")
		assert.Contains(t, runner.commands[0], "/etc/rc.conf")
		assert.Contains(t, runner.commands[1], "1.1.1.1")
		assert.Contains(t, runner.commands[1], "8.8.8.8")
		assert.Contains(t, runner.commands[1], "/etc/resolv.conf")
	})

	t.Run("returns error when dhcp setup fails", func(t *testing.T) {
		t.Parallel()

		runner := &mockRunner{err: assert.AnError, results: nil, commands: nil}
		customizer := customize.New(runner)

		networkErr := customizer.ConfigureNetworking(context.Background())

		require.Error(t, networkErr)
		assert.Contains(t, networkErr.Error(), "dhcp")
	})
}

func TestDefaultSecurityTools(t *testing.T) {
	t.Parallel()

	tools := customize.DefaultSecurityTools()

	expectedTools := []string{
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

	assert.Equal(t, expectedTools, tools)
	assert.Len(t, tools, 10)
}
