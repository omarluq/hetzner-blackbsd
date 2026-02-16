package ssh

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
)

const downloadFilePermissions = 0o600

// UploadFile transfers a local file to the remote host via SFTP.
func (c *Client) UploadFile(ctx context.Context, localPath, remotePath string) error {
	expandedLocal := filepath.Clean(expandPath(localPath))

	local, err := os.Open(expandedLocal)
	if err != nil {
		return &Error{Message: "open local file", Err: err}
	}
	defer closeQuietly(local)

	conn, err := c.dial(ctx)
	if err != nil {
		return err
	}
	defer closeQuietly(conn)

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return &Error{Message: "create sftp client", Err: err}
	}
	defer closeQuietly(sftpClient)

	remote, err := sftpClient.Create(remotePath)
	if err != nil {
		return &Error{Message: "create remote file", Err: err}
	}
	defer closeQuietly(remote)

	if _, err = io.Copy(remote, local); err != nil {
		return &Error{Message: "upload file", Err: err}
	}

	return nil
}

// DownloadFile transfers a remote file to the local filesystem via SFTP.
func (c *Client) DownloadFile(ctx context.Context, remotePath, localPath string) error {
	expandedLocal := filepath.Clean(expandPath(localPath))

	conn, err := c.dial(ctx)
	if err != nil {
		return err
	}
	defer closeQuietly(conn)

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return &Error{Message: "create sftp client", Err: err}
	}
	defer closeQuietly(sftpClient)

	remote, err := sftpClient.Open(remotePath)
	if err != nil {
		return &Error{Message: "open remote file", Err: err}
	}
	defer closeQuietly(remote)

	local, err := os.OpenFile(expandedLocal, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, downloadFilePermissions)
	if err != nil {
		return &Error{Message: "create local file", Err: err}
	}
	defer closeQuietly(local)

	if _, err = io.Copy(local, remote); err != nil {
		return &Error{Message: "download file", Err: err}
	}

	return nil
}
