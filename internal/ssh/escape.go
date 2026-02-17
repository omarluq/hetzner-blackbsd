package ssh

import "strings"

// EscapeShellArg escapes a string for safe use as a shell argument.
// Uses single-quote escaping: wraps in single quotes and escapes embedded quotes.
func EscapeShellArg(arg string) string {
	if arg == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(arg, "'", "'\\''") + "'"
}
