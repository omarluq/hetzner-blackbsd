package ssh_test

import (
	"testing"

	"github.com/omarluq/hetzner-blackbsd/internal/ssh"
	"github.com/stretchr/testify/assert"
)

func TestEscapeShellArgBasic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string returns empty quotes", "", "''"},
		{"simple string wraps in single quotes", "hello", "'hello'"},
		{"string with spaces", "hello world", "'hello world'"},
		{"string with single quote", "it's", "'it'\\''s'"},
		{"string with multiple single quotes", "it's a 'test'", "'it'\\''s a '\\''test'\\'''"},
		{"string with double quotes", `"hello"`, `'"hello"'`},
		{"shell metacharacters", "foo; rm -rf /", "'foo; rm -rf /'"},
		{"dollar sign", "$HOME", "'$HOME'"},
		{"backticks", "`whoami`", "'`whoami`'"},
		{"newline", "line1\nline2", "'line1\nline2'"},
		{"pipe", "foo | bar", "'foo | bar'"},
		{"ampersand", "foo & bar", "'foo & bar'"},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, testCase.expected, ssh.EscapeShellArg(testCase.input))
		})
	}
}

func TestEscapeShellArgTypicalInputs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"typical hostname", "blackbsd", "'blackbsd'"},
		{"typical package name", "nmap", "'nmap'"},
		{"path with slashes", "/tmp/image.raw.xz", "'/tmp/image.raw.xz'"},
		{
			"url string",
			"https://cdn.netbsd.org/pub/NetBSD/NetBSD-10.1/amd64/installation/cdrom/boot-com.iso",
			"'https://cdn.netbsd.org/pub/NetBSD/NetBSD-10.1/amd64/installation/cdrom/boot-com.iso'",
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, testCase.expected, ssh.EscapeShellArg(testCase.input))
		})
	}
}
