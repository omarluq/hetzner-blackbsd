// Package extract tests disk image extraction utilities.
package extract_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/omarluq/hetzner-blackbsd/internal/extract"
	"github.com/omarluq/hetzner-blackbsd/internal/ssh"
	"github.com/stretchr/testify/assert"
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

	t.Run("creates extractor with runner and device", func(t *testing.T) {
		t.Parallel()
		assert.NotNil(t, extract.New(newMock(nil), "/dev/sda"))
	})

	t.Run("creates extractor with different device", func(t *testing.T) {
		t.Parallel()
		assert.NotNil(t, extract.New(newMock(nil), "/dev/nvme0n1"))
	})
}

func TestExtractRawImageSuccess(t *testing.T) {
	t.Parallel()

	t.Run("generates correct dd and xz command", func(t *testing.T) {
		t.Parallel()

		runner := newMock(map[string]ssh.CommandResult{
			"dd if=/dev/sda bs=4M status=progress | xz -T0 -9 > /tmp/image.raw.xz": okResult(),
		})

		extractor := extract.New(runner, "/dev/sda")
		extractErr := extractor.ExtractRawImage(context.Background(), "/tmp/image.raw.xz")

		assert.NoError(t, extractErr)
		assert.Contains(t, runner.lastCommand, "dd if=/dev/sda")
		assert.Contains(t, runner.lastCommand, "xz -T0 -9")
	})

	t.Run("uses different device path", func(t *testing.T) {
		t.Parallel()

		runner := newMock(map[string]ssh.CommandResult{
			"dd if=/dev/nvme0n1 bs=4M status=progress | xz -T0 -9 > /tmp/out.img.xz": okResult(),
		})

		extractor := extract.New(runner, "/dev/nvme0n1")
		extractErr := extractor.ExtractRawImage(context.Background(), "/tmp/out.img.xz")

		assert.NoError(t, extractErr)
		assert.Contains(t, runner.lastCommand, "dd if=/dev/nvme0n1")
	})
}

func TestExtractRawImageErrors(t *testing.T) {
	t.Parallel()

	t.Run("returns error on command failure", func(t *testing.T) {
		t.Parallel()

		runner := newMock(map[string]ssh.CommandResult{
			"dd if=/dev/sda bs=4M status=progress | xz -T0 -9 > /tmp/image.raw.xz": errResult("device not found"),
		})

		extractErr := extract.New(runner, "/dev/sda").ExtractRawImage(context.Background(), "/tmp/image.raw.xz")

		assert.Error(t, extractErr)
		assert.Contains(t, extractErr.Error(), "extract raw image")
	})

	t.Run("returns error on runner error", func(t *testing.T) {
		t.Parallel()

		extractErr := extract.New(newErrMock(fmt.Errorf("connection lost")), "/dev/sda").
			ExtractRawImage(context.Background(), "/tmp/image.raw.xz")

		assert.Error(t, extractErr)
	})
}

func isoSuccessResults() map[string]ssh.CommandResult {
	return map[string]ssh.CommandResult{
		"mount -r /dev/sda1 /mnt/iso": okResult(),
		"xorriso -as mkisofs -o /tmp/blackbsd.iso -b boot/cdboot -no-emul-boot /mnt/iso": okResult(),
		"umount /mnt/iso": okResult(),
	}
}

func TestExtractISOSuccess(t *testing.T) {
	t.Parallel()

	t.Run("mounts device creates ISO unmounts", func(t *testing.T) {
		t.Parallel()

		runner := newMock(isoSuccessResults())
		extractor := extract.New(runner, "/dev/sda")
		isoErr := extractor.ExtractISO(context.Background(), "/mnt/iso", "/tmp/blackbsd.iso")

		assert.NoError(t, isoErr)
		assert.Len(t, runner.commands, 3)
		assert.Contains(t, runner.commands[0], "mount -r /dev/sda1")
		assert.Contains(t, runner.commands[1], "xorriso")
		assert.Contains(t, runner.commands[2], "umount")
	})

	t.Run("uses p separator for nvme partition", func(t *testing.T) {
		t.Parallel()

		runner := newMock(map[string]ssh.CommandResult{
			"mount -r /dev/nvme0n1p1 /mnt/build":                                         okResult(),
			"xorriso -as mkisofs -o /output.iso -b boot/cdboot -no-emul-boot /mnt/build": okResult(),
			"umount /mnt/build": okResult(),
		})

		isoErr := extract.New(runner, "/dev/nvme0n1").ExtractISO(context.Background(), "/mnt/build", "/output.iso")

		assert.NoError(t, isoErr)
		assert.Contains(t, runner.commands[0], "mount -r /dev/nvme0n1p1")
	})
}

func TestExtractISOMountError(t *testing.T) {
	t.Parallel()

	runner := newMock(map[string]ssh.CommandResult{
		"mount -r /dev/sda1 /mnt/iso": errResult("mount point does not exist"),
	})

	isoErr := extract.New(runner, "/dev/sda").ExtractISO(context.Background(), "/mnt/iso", "/tmp/blackbsd.iso")

	assert.Error(t, isoErr)
	assert.Contains(t, isoErr.Error(), "mount device")
}

func TestExtractISOXorrisoError(t *testing.T) {
	t.Parallel()

	xorrisoCmd := "xorriso -as mkisofs -o /tmp/blackbsd.iso -b boot/cdboot -no-emul-boot /mnt/iso"
	runner := newMock(map[string]ssh.CommandResult{
		"mount -r /dev/sda1 /mnt/iso": okResult(),
		xorrisoCmd:                     errResult("bootloader not found"),
	})

	isoErr := extract.New(runner, "/dev/sda").ExtractISO(context.Background(), "/mnt/iso", "/tmp/blackbsd.iso")

	assert.Error(t, isoErr)
	assert.Contains(t, isoErr.Error(), "create ISO")
}

func TestExtractISOUmountError(t *testing.T) {
	t.Parallel()

	runner := newMock(map[string]ssh.CommandResult{
		"mount -r /dev/sda1 /mnt/iso": okResult(),
		"xorriso -as mkisofs -o /tmp/blackbsd.iso -b boot/cdboot -no-emul-boot /mnt/iso": okResult(),
		"umount /mnt/iso": errResult("device is busy"),
	})

	isoErr := extract.New(runner, "/dev/sda").ExtractISO(context.Background(), "/mnt/iso", "/tmp/blackbsd.iso")

	assert.Error(t, isoErr)
	assert.Contains(t, isoErr.Error(), "unmount device")
}

func TestImageSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		imagePath   string
		stdout      string
		errContains string
		exitCode    int
		expected    int64
		expectErr   bool
	}{
		{
			name:        "returns parsed size from stat output",
			imagePath:   "/tmp/image.raw.xz",
			stdout:      "1073741824",
			exitCode:    0,
			expected:    1073741824,
			expectErr:   false,
			errContains: "",
		},
		{
			name:        "handles different file path",
			imagePath:   "/var/tmp/blackbsd.img.xz",
			stdout:      "2147483648",
			exitCode:    0,
			expected:    2147483648,
			expectErr:   false,
			errContains: "",
		},
		{
			name:        "returns error on stat failure",
			imagePath:   "/tmp/image.raw.xz",
			stdout:      "",
			exitCode:    1,
			expected:    0,
			expectErr:   true,
			errContains: "get image size",
		},
		{
			name:        "returns error on parse failure",
			imagePath:   "/tmp/image.raw.xz",
			stdout:      "not a number",
			exitCode:    0,
			expected:    0,
			expectErr:   true,
			errContains: "parse image size",
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			command := fmt.Sprintf("stat -c %%s %s", testCase.imagePath)
			runner := newMock(map[string]ssh.CommandResult{
				command: {Stdout: testCase.stdout, Stderr: "", ExitCode: testCase.exitCode},
			})

			size, sizeErr := extract.New(runner, "/dev/sda").ImageSize(context.Background(), testCase.imagePath)

			if testCase.expectErr {
				assert.Error(t, sizeErr)
				assert.Contains(t, sizeErr.Error(), testCase.errContains)
			} else {
				assert.NoError(t, sizeErr)
			}
			assert.Equal(t, testCase.expected, size)
		})
	}
}

func TestChecksum(t *testing.T) {
	t.Parallel()

	longHash := "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6a7b8c9d0e1f2"

	tests := []struct {
		name        string
		imagePath   string
		stdout      string
		expected    string
		errContains string
		exitCode    int
		expectErr   bool
	}{
		{
			name:        "returns checksum from sha256sum output",
			imagePath:   "/tmp/image.raw.xz",
			stdout:      longHash + "  /tmp/image.raw.xz",
			exitCode:    0,
			expected:    longHash,
			expectErr:   false,
			errContains: "",
		},
		{
			name:        "handles different file path",
			imagePath:   "/var/tmp/blackbsd.img.xz",
			stdout:      "f1e2d3c4b5a6978800  /var/tmp/blackbsd.img.xz",
			exitCode:    0,
			expected:    "f1e2d3c4b5a6978800",
			expectErr:   false,
			errContains: "",
		},
		{
			name:        "returns error on sha256sum failure",
			imagePath:   "/tmp/image.raw.xz",
			stdout:      "",
			exitCode:    1,
			expected:    "",
			expectErr:   true,
			errContains: "compute checksum",
		},
		{
			name:        "returns error on empty output",
			imagePath:   "/tmp/image.raw.xz",
			stdout:      "",
			exitCode:    0,
			expected:    "",
			expectErr:   true,
			errContains: "parse checksum",
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			command := fmt.Sprintf("sha256sum %s", testCase.imagePath)
			stderr := ""
			if testCase.exitCode != 0 {
				stderr = "file not found"
			}
			runner := newMock(map[string]ssh.CommandResult{
				command: {Stdout: testCase.stdout, Stderr: stderr, ExitCode: testCase.exitCode},
			})

			result, checksumErr := extract.New(runner, "/dev/sda").Checksum(context.Background(), testCase.imagePath)

			if testCase.expectErr {
				assert.Error(t, checksumErr)
				assert.Contains(t, checksumErr.Error(), testCase.errContains)
			} else {
				assert.NoError(t, checksumErr)
			}
			assert.Equal(t, testCase.expected, result)
		})
	}
}
