package ssh

import (
	"bytes"
	"context"
	"errors"
	"io"

	"golang.org/x/crypto/ssh"
)

const (
	ptyCols = 80
	ptyRows = 40
)

// ExecInteractive runs a command with a PTY and pipes the given input.
// This is used for QEMU serial console interaction where install commands
// are piped through an interactive terminal session.
func (c *Client) ExecInteractive(ctx context.Context, command string, input io.Reader) (CommandResult, error) {
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

	if err = session.RequestPty("xterm", ptyRows, ptyCols, ssh.TerminalModes{}); err != nil {
		return CommandResult{Stdout: "", Stderr: "", ExitCode: 0}, &Error{Message: "request pty", Err: err}
	}

	session.Stdin = input

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	exitCode := 0
	if runErr := session.Run(command); runErr != nil {
		var exitErr *ssh.ExitError
		if errors.As(runErr, &exitErr) {
			exitCode = exitErr.ExitStatus()
		} else {
			return CommandResult{
				Stdout:   "",
				Stderr:   "",
				ExitCode: 0,
			}, &Error{Message: "run interactive command", Err: runErr}
		}
	}

	return CommandResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}, nil
}
