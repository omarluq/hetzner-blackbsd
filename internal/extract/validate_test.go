package extract_test

import (
	"testing"

	"github.com/omarluq/hetzner-blackbsd/internal/extract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatePathAcceptsValidPaths(t *testing.T) {
	t.Parallel()

	t.Run("accepts valid absolute path", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, extract.ValidatePath("/tmp/image.raw.xz"))
	})

	t.Run("accepts valid relative path", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, extract.ValidatePath("output/image.raw.xz"))
	})

	t.Run("accepts device path", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, extract.ValidatePath("/dev/sda"))
	})

	t.Run("accepts nvme device path", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, extract.ValidatePath("/dev/nvme0n1"))
	})

	t.Run("accepts mount point path", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, extract.ValidatePath("/mnt/iso"))
	})
}

func TestValidatePathRejectsEmptyString(t *testing.T) {
	t.Parallel()

	err := extract.ValidatePath("")
	require.Error(t, err)
	assert.ErrorIs(t, err, extract.ErrInvalidPath)
}

func TestValidatePathRejectsPathTraversal(t *testing.T) {
	t.Parallel()

	t.Run("rejects path traversal with double dots", func(t *testing.T) {
		t.Parallel()
		err := extract.ValidatePath("/tmp/../etc/passwd")
		require.Error(t, err)
		assert.ErrorIs(t, err, extract.ErrInvalidPath)
	})

	t.Run("rejects path starting with double dots", func(t *testing.T) {
		t.Parallel()
		err := extract.ValidatePath("../../etc/shadow")
		require.Error(t, err)
		assert.ErrorIs(t, err, extract.ErrInvalidPath)
	})
}

func TestValidatePathRejectsShellMetacharacters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{"semicolon injection", "/tmp/foo; rm -rf /"},
		{"pipe injection", "/tmp/foo | cat /etc/passwd"},
		{"ampersand injection", "/tmp/foo & whoami"},
		{"dollar sign injection", "/tmp/$HOME"},
		{"backtick injection", "/tmp/`whoami`"},
		{"newline injection", "/tmp/foo\nrm -rf /"},
		{"carriage return injection", "/tmp/foo\rbar"},
		{"parenthesis injection", "/tmp/$(whoami)"},
		{"redirect less than", "/tmp/foo < /etc/passwd"},
		{"redirect greater than", "/tmp/foo > /etc/passwd"},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			err := extract.ValidatePath(testCase.input)
			require.Error(t, err)
			assert.ErrorIs(t, err, extract.ErrInvalidPath)
		})
	}
}
