package ssh_test

import (
	"os"
	"strings"
	"testing"

	"github.com/omarluq/hetzner-blackbsd/internal/ssh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandResultSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		result   ssh.CommandResult
		expected bool
	}{
		{
			name: "zero exit code is success",
			result: ssh.CommandResult{
				Stdout:   "output",
				Stderr:   "",
				ExitCode: 0,
			},
			expected: true,
		},
		{
			name: "non-zero exit code is failure",
			result: ssh.CommandResult{
				Stdout:   "",
				Stderr:   "error",
				ExitCode: 1,
			},
			expected: false,
		},
		{
			name: "exit code 127 is failure",
			result: ssh.CommandResult{
				Stdout:   "",
				Stderr:   "command not found",
				ExitCode: 127,
			},
			expected: false,
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expected, testCase.result.Success())
		})
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	t.Run("Error formats with underlying error", func(t *testing.T) {
		t.Parallel()

		var sshErr ssh.Error
		sshErr.Message = "connect failed"
		sshErr.Err = assert.AnError

		msg := sshErr.Error()
		assert.Contains(t, msg, "ssh")
		assert.Contains(t, msg, "connect failed")
	})

	t.Run("Error formats without underlying error", func(t *testing.T) {
		t.Parallel()

		var sshErr ssh.Error
		sshErr.Message = "connection lost"

		msg := sshErr.Error()
		assert.Contains(t, msg, "ssh")
		assert.Contains(t, msg, "connection lost")
	})

	t.Run("Error unwraps", func(t *testing.T) {
		t.Parallel()

		underlying := assert.AnError

		var sshErr ssh.Error
		sshErr.Message = "test"
		sshErr.Err = underlying

		assert.Equal(t, underlying, sshErr.Unwrap())
	})
}

func TestCommandFailedError(t *testing.T) {
	t.Parallel()

	t.Run("formats message with all details", func(t *testing.T) {
		t.Parallel()

		var cmdErr ssh.CommandFailedError
		cmdErr.Command = "ls /nonexistent"
		cmdErr.Stderr = "No such file or directory"
		cmdErr.ExitCode = 1

		msg := cmdErr.Error()
		assert.Contains(t, msg, "ls /nonexistent")
		assert.Contains(t, msg, "exit code 1")
		assert.Contains(t, msg, "No such file or directory")
	})

	t.Run("exits exit code", func(t *testing.T) {
		t.Parallel()

		var cmdErr ssh.CommandFailedError
		cmdErr.ExitCode = 42

		assert.Equal(t, 42, cmdErr.ExitCode)
	})
}

func TestExpandPath(t *testing.T) {
	t.Parallel()

	t.Run("expands tilde to home directory", func(t *testing.T) {
		t.Parallel()

		home := "/home/user"
		require.NoError(t, os.Setenv("HOME", home))
		t.Cleanup(func() { require.NoError(t, os.Unsetenv("HOME")) })

		result := ssh.ExpandPath("~/test_key")
		assert.False(t, strings.HasPrefix(result, "~"))
		assert.Contains(t, result, "/test_key")
	})

	t.Run("handles empty string", func(t *testing.T) {
		t.Parallel()

		result := ssh.ExpandPath("")
		assert.Equal(t, "", result)
	})

	t.Run("leaves absolute path unchanged", func(t *testing.T) {
		t.Parallel()

		result := ssh.ExpandPath("/absolute/path/key")
		assert.Equal(t, "/absolute/path/key", result)
	})

	t.Run("leaves relative path without tilde unchanged", func(t *testing.T) {
		t.Parallel()

		result := ssh.ExpandPath("relative/path/key")
		assert.Equal(t, "relative/path/key", result)
	})
}
