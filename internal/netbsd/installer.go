// Package netbsd automates NetBSD installation via QEMU inside Hetzner rescue mode.
package netbsd

import (
	"context"
	"fmt"
	"time"

	"github.com/omarluq/hetzner-blackbsd/internal/runner"
	"github.com/omarluq/hetzner-blackbsd/internal/ssh"
)

const qemuTimeout = 15 * time.Minute

// Installer automates NetBSD installation by downloading the ISO and running
// QEMU with KVM acceleration inside Hetzner rescue mode.
type Installer struct {
	runner  runner.Runner
	version string
	arch    string
}

// New creates an Installer for the given NetBSD version and architecture.
func New(exec runner.Runner, version, arch string) *Installer {
	return &Installer{
		runner:  exec,
		version: version,
		arch:    arch,
	}
}

// ISODownloadURL returns the CDN URL for the NetBSD serial-console boot ISO.
func (inst *Installer) ISODownloadURL() string {
	return fmt.Sprintf(
		"https://cdn.netbsd.org/pub/NetBSD/NetBSD-%s/%s/installation/cdrom/boot-com.iso",
		inst.version,
		inst.arch,
	)
}

// DownloadISO fetches the NetBSD boot ISO to the given directory on the remote host.
// It returns the remote path of the downloaded file.
func (inst *Installer) DownloadISO(ctx context.Context, destDir string) (string, error) {
	isoPath := fmt.Sprintf("%s/netbsd-%s-%s.iso", destDir, inst.version, inst.arch)
	cmd := fmt.Sprintf("wget -O %s %s",
		ssh.EscapeShellArg(isoPath),
		ssh.EscapeShellArg(inst.ISODownloadURL()))

	result, err := inst.runner.Exec(ctx, cmd)
	if err != nil {
		return "", fmt.Errorf("download iso: %w", err)
	}

	if !result.Success() {
		return "", fmt.Errorf("download iso: exit code %d: %s", result.ExitCode, result.Stderr)
	}

	return isoPath, nil
}

// InstallViaQEMU runs QEMU with KVM to boot the ISO and install NetBSD onto the target device.
func (inst *Installer) InstallViaQEMU(ctx context.Context, isoPath, device string) error {
	cmd := fmt.Sprintf(
		"qemu-system-x86_64 -enable-kvm -m 4G -smp 4 -cdrom %s -boot d "+
			"-drive file=%s,format=raw -nographic -serial mon:stdio",
		ssh.EscapeShellArg(isoPath),
		ssh.EscapeShellArg(device),
	)

	timeoutCtx, cancel := context.WithTimeout(ctx, qemuTimeout)
	defer cancel()

	result, err := inst.runner.Exec(timeoutCtx, cmd)
	if err != nil {
		return fmt.Errorf("run qemu install: %w", err)
	}

	if !result.Success() {
		return fmt.Errorf("run qemu install: exit code %d: %s", result.ExitCode, result.Stderr)
	}

	return nil
}
