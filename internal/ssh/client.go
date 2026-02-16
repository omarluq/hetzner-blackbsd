package ssh

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/cenkalti/backoff/v4"
	"golang.org/x/crypto/ssh"
)

const (
	defaultPort    = 22
	defaultTimeout = 30 * time.Second
)

// Client wraps golang.org/x/crypto/ssh with higher-level operations.
type Client struct {
	config  *ssh.ClientConfig
	host    string
	keyPath string
	port    int
}

// NewClient creates a new SSH client for the given host using key authentication.
func NewClient(host, keyPath string) (*Client, error) {
	expandedPath := expandPath(keyPath)

	key, err := readKeyFile(expandedPath)
	if err != nil {
		return nil, &Error{Message: "read private key", Err: err}
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, &Error{Message: "parse private key", Err: err}
	}

	var clientConfig ssh.ClientConfig
	clientConfig.User = "root"
	clientConfig.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	clientConfig.HostKeyCallback = insecureHostKeyFallback
	clientConfig.Timeout = defaultTimeout

	return &Client{
		host:    host,
		port:    defaultPort,
		keyPath: expandedPath,
		config:  &clientConfig,
	}, nil
}

// readKeyFile reads an SSH private key from a validated path.
func readKeyFile(keyPath string) ([]byte, error) {
	cleaned := filepath.Clean(keyPath)
	return os.ReadFile(cleaned)
}

// WithPort sets a custom SSH port.
func (c *Client) WithPort(port int) *Client {
	c.port = port
	return c
}

// Exec runs a command on the remote host and returns the result.
func (c *Client) Exec(ctx context.Context, command string) (CommandResult, error) {
	conn, err := c.dial(ctx)
	if err != nil {
		return CommandResult{Stdout: "", Stderr: "", ExitCode: 0}, err
	}
	defer closeQuietly(conn)

	session, err := conn.NewSession()
	if err != nil {
		return CommandResult{Stdout: "", Stderr: "", ExitCode: 0}, &Error{Message: "create session", Err: err}
	}
	defer closeQuietly(session)

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	exitCode := 0
	if runErr := session.Run(command); runErr != nil {
		var exitErr *ssh.ExitError
		if errors.As(runErr, &exitErr) {
			exitCode = exitErr.ExitStatus()
		} else {
			return CommandResult{Stdout: "", Stderr: "", ExitCode: 0}, &Error{Message: "run command", Err: runErr}
		}
	}

	return CommandResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}, nil
}

// WaitForReady waits until the host accepts SSH connections.
func (c *Client) WaitForReady(ctx context.Context) error {
	backoffPolicy := backoff.NewExponentialBackOff()
	backoffPolicy.MaxElapsedTime = 5 * time.Minute
	backoffPolicy.InitialInterval = 2 * time.Second

	retryOperation := func() error {
		conn, dialErr := c.dial(ctx)
		if dialErr != nil {
			slog.Debug("ssh not ready yet", "host", c.host, "err", dialErr)
			return dialErr
		}
		closeQuietly(conn)
		return nil
	}

	if err := backoff.Retry(retryOperation, backoff.WithContext(backoffPolicy, ctx)); err != nil {
		return &Error{
			Message: fmt.Sprintf("host %s not reachable after timeout", c.host),
			Err:     err,
		}
	}

	slog.Info("ssh connection ready", "host", c.host)
	return nil
}

func (c *Client) dial(ctx context.Context) (*ssh.Client, error) {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	var dialer net.Dialer
	netConn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, &Error{Message: "connect to " + addr, Err: err}
	}

	sshConn, chans, reqs, err := ssh.NewClientConn(netConn, addr, c.config)
	if err != nil {
		closeQuietly(netConn)
		return nil, &Error{Message: "ssh handshake", Err: err}
	}

	return ssh.NewClient(sshConn, chans, reqs), nil
}

// insecureHostKeyFallback accepts any host key. Rescue mode servers have
// ephemeral host keys that change on each provisioning cycle, making
// host-key verification impractical for this use case.
func insecureHostKeyFallback(_ string, _ net.Addr, _ ssh.PublicKey) error {
	return nil
}

type closer interface {
	Close() error
}

func closeQuietly(resource closer) {
	if err := resource.Close(); err != nil {
		slog.Debug("close error", "error", err)
	}
}

func expandPath(path string) string {
	if path != "" && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}
