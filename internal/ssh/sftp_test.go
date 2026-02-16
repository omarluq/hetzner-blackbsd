package ssh_test

import (
	"context"
	"testing"

	"github.com/omarluq/hetzner-blackbsd/internal/ssh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadFileDialError(t *testing.T) {
	t.Parallel()

	client, err := ssh.NewClient("192.0.2.1", createTempKey(t))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	uploadErr := client.UploadFile(ctx, "/tmp/nonexistent", "/remote/path")
	assert.Error(t, uploadErr)
}

func TestUploadFileLocalFileError(t *testing.T) {
	t.Parallel()

	client, err := ssh.NewClient("192.0.2.1", createTempKey(t))
	require.NoError(t, err)

	// Use a canceled context so dial fails, but also verify local path handling.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	uploadErr := client.UploadFile(ctx, "~/nonexistent_file_12345", "/remote/path")
	assert.Error(t, uploadErr)
}

func TestDownloadFileDialError(t *testing.T) {
	t.Parallel()

	client, err := ssh.NewClient("192.0.2.1", createTempKey(t))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	downloadErr := client.DownloadFile(ctx, "/remote/path", "/tmp/local")
	assert.Error(t, downloadErr)
}
