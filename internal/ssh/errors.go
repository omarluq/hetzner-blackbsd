// Package ssh provides an SSH client wrapper using golang.org/x/crypto/ssh.
package ssh

import "fmt"

// Error is the base error for SSH operations.
type Error struct {
	Err     error
	Message string
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("ssh: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("ssh: %s", e.Message)
}

func (e *Error) Unwrap() error { return e.Err }

// CommandFailedError indicates a remote command exited with a non-zero code.
type CommandFailedError struct {
	Command  string
	Stderr   string
	ExitCode int
}

func (e *CommandFailedError) Error() string {
	return fmt.Sprintf("command %q failed with exit code %d: %s", e.Command, e.ExitCode, e.Stderr)
}
