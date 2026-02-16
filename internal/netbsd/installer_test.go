// Package netbsd tests NetBSD installation automation.
package netbsd_test

import (
	"context"
	"testing"

	"github.com/omarluq/hetzner-blackbsd/internal/netbsd"
	"github.com/omarluq/hetzner-blackbsd/internal/ssh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRunner is a test double for the Runner interface.
type mockRunner struct {
	err         error
	results     map[string]ssh.CommandResult
	lastCommand string
	commands    []string
}

func (mock *mockRunner) Exec(_ context.Context, command string) (ssh.CommandResult, error) {
	mock.commands = append(mock.commands, command)
	mock.lastCommand = command

	if mock.err != nil {
		return ssh.CommandResult{Stdout: "", Stderr: "", ExitCode: 0}, mock.err
	}

	if result, found := mock.results[command]; found {
		return result, nil
	}

	return ssh.CommandResult{Stdout: "", Stderr: "", ExitCode: 0}, nil
}

func okResult() ssh.CommandResult {
	return ssh.CommandResult{Stdout: "", Stderr: "", ExitCode: 0}
}

func errResult(stderr string) ssh.CommandResult {
	return ssh.CommandResult{Stdout: "", Stderr: stderr, ExitCode: 1}
}

func newMock(results map[string]ssh.CommandResult) *mockRunner {
	return &mockRunner{err: nil, results: results, lastCommand: "", commands: []string{}}
}

func newErrMock(mockErr error) *mockRunner {
	return &mockRunner{err: mockErr, results: nil, lastCommand: "", commands: nil}
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("creates installer with runner and version", func(t *testing.T) {
		t.Parallel()

		installer := netbsd.New(newMock(nil), "10.1", "amd64")
		assert.NotNil(t, installer)
	})

	t.Run("creates installer for different arch", func(t *testing.T) {
		t.Parallel()

		installer := netbsd.New(newMock(nil), "10.0", "earmv7hf")
		assert.NotNil(t, installer)
	})
}

func TestISODownloadURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		version  string
		arch     string
		expected string
	}{
		{
			name:     "10.1 amd64",
			version:  "10.1",
			arch:     "amd64",
			expected: "https://cdn.netbsd.org/pub/NetBSD/NetBSD-10.1/amd64/installation/cdrom/boot-com.iso",
		},
		{
			name:     "10.0 earmv7hf",
			version:  "10.0",
			arch:     "earmv7hf",
			expected: "https://cdn.netbsd.org/pub/NetBSD/NetBSD-10.0/earmv7hf/installation/cdrom/boot-com.iso",
		},
		{
			name:     "9.3 amd64",
			version:  "9.3",
			arch:     "amd64",
			expected: "https://cdn.netbsd.org/pub/NetBSD/NetBSD-9.3/amd64/installation/cdrom/boot-com.iso",
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			installer := netbsd.New(newMock(nil), testCase.version, testCase.arch)

			assert.Equal(t, testCase.expected, installer.ISODownloadURL())
		})
	}
}

func TestDownloadISO(t *testing.T) {
	t.Parallel()

	t.Run("generates correct wget command", func(t *testing.T) {
		t.Parallel()

		wgetCmd := "wget -O /tmp/netbsd-10.1-amd64.iso " +
			"https://cdn.netbsd.org/pub/NetBSD/NetBSD-10.1/amd64/installation/cdrom/boot-com.iso"
		mock := newMock(map[string]ssh.CommandResult{wgetCmd: okResult()})
		installer := netbsd.New(mock, "10.1", "amd64")

		path, err := installer.DownloadISO(context.Background(), "/tmp")

		require.NoError(t, err)
		assert.Equal(t, "/tmp/netbsd-10.1-amd64.iso", path)
		assert.Equal(t, wgetCmd, mock.lastCommand)
	})

	t.Run("uses different dest directory", func(t *testing.T) {
		t.Parallel()

		wgetCmd := "wget -O /var/tmp/netbsd-10.0-amd64.iso " +
			"https://cdn.netbsd.org/pub/NetBSD/NetBSD-10.0/amd64/installation/cdrom/boot-com.iso"
		mock := newMock(map[string]ssh.CommandResult{wgetCmd: okResult()})
		installer := netbsd.New(mock, "10.0", "amd64")

		path, err := installer.DownloadISO(context.Background(), "/var/tmp")

		require.NoError(t, err)
		assert.Equal(t, "/var/tmp/netbsd-10.0-amd64.iso", path)
	})

	t.Run("returns error on wget failure", func(t *testing.T) {
		t.Parallel()

		wgetCmd := "wget -O /tmp/netbsd-10.1-amd64.iso " +
			"https://cdn.netbsd.org/pub/NetBSD/NetBSD-10.1/amd64/installation/cdrom/boot-com.iso"
		mock := newMock(map[string]ssh.CommandResult{wgetCmd: errResult("404 Not Found")})
		installer := netbsd.New(mock, "10.1", "amd64")

		path, err := installer.DownloadISO(context.Background(), "/tmp")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "download iso")
		assert.Equal(t, "", path)
	})

	t.Run("returns error on exec failure", func(t *testing.T) {
		t.Parallel()

		installer := netbsd.New(newErrMock(assert.AnError), "10.1", "amd64")

		path, err := installer.DownloadISO(context.Background(), "/tmp")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "download iso")
		assert.Equal(t, "", path)
	})
}

func TestInstallViaQEMU(t *testing.T) {
	t.Parallel()

	t.Run("generates correct qemu command", func(t *testing.T) {
		t.Parallel()

		qemuCmd := "qemu-system-x86_64 -enable-kvm -m 4G -smp 4" +
			" -cdrom /tmp/iso.iso -boot d -drive file=/dev/sda,format=raw" +
			" -nographic -serial mon:stdio"
		mock := newMock(map[string]ssh.CommandResult{qemuCmd: okResult()})
		installer := netbsd.New(mock, "10.1", "amd64")

		err := installer.InstallViaQEMU(context.Background(), "/tmp/iso.iso", "/dev/sda")

		require.NoError(t, err)
		assert.Contains(t, mock.lastCommand, "-enable-kvm")
		assert.Contains(t, mock.lastCommand, "-cdrom /tmp/iso.iso")
		assert.Contains(t, mock.lastCommand, "file=/dev/sda,format=raw")
	})

	t.Run("uses different device path", func(t *testing.T) {
		t.Parallel()

		qemuCmd := "qemu-system-x86_64 -enable-kvm -m 4G -smp 4" +
			" -cdrom /tmp/iso.iso -boot d -drive file=/dev/nvme0n1,format=raw" +
			" -nographic -serial mon:stdio"
		mock := newMock(map[string]ssh.CommandResult{qemuCmd: okResult()})
		installer := netbsd.New(mock, "10.1", "amd64")

		err := installer.InstallViaQEMU(context.Background(), "/tmp/iso.iso", "/dev/nvme0n1")

		require.NoError(t, err)
		assert.Contains(t, mock.lastCommand, "file=/dev/nvme0n1,format=raw")
	})

	t.Run("returns error on qemu failure", func(t *testing.T) {
		t.Parallel()

		qemuCmd := "qemu-system-x86_64 -enable-kvm -m 4G -smp 4" +
			" -cdrom /tmp/iso.iso -boot d -drive file=/dev/sda,format=raw" +
			" -nographic -serial mon:stdio"
		mock := newMock(map[string]ssh.CommandResult{qemuCmd: errResult("device busy")})
		installer := netbsd.New(mock, "10.1", "amd64")

		err := installer.InstallViaQEMU(context.Background(), "/tmp/iso.iso", "/dev/sda")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "run qemu install")
	})

	t.Run("returns error on exec failure", func(t *testing.T) {
		t.Parallel()

		installer := netbsd.New(newErrMock(assert.AnError), "10.1", "amd64")

		err := installer.InstallViaQEMU(context.Background(), "/tmp/iso.iso", "/dev/sda")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "run qemu install")
	})
}
