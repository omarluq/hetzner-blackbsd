package ssh

// CommandResult holds the output of a remote command.
type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Success returns true if the command exited with code 0.
func (r CommandResult) Success() bool {
	return r.ExitCode == 0
}
