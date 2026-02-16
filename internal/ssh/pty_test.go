package ssh_test

import (
	"context"
	"strings"
	"testing"

	"github.com/omarluq/hetzner-blackbsd/internal/ssh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecInteractiveDialError(t *testing.T) {
	t.Parallel()

	client, err := ssh.NewClient("192.0.2.1", createTempKey(t))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, execErr := client.ExecInteractive(ctx, "echo hello", strings.NewReader("input\n"))
	assert.Error(t, execErr)
	assert.Equal(t, "", result.Stdout)
	assert.Equal(t, "", result.Stderr)
	assert.Equal(t, 0, result.ExitCode)
}
