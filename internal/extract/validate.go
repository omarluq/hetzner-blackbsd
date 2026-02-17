package extract

import (
	"errors"
	"path/filepath"
	"strings"
)

// ErrInvalidPath is returned when a path contains traversal sequences or shell metacharacters.
var ErrInvalidPath = errors.New("invalid path: contains path traversal or shell metacharacters")

// ValidatePath checks that a path is safe for use in shell commands.
func ValidatePath(path string) error {
	if path == "" {
		return ErrInvalidPath
	}

	// Reject path traversal in original or cleaned form
	if strings.Contains(path, "..") {
		return ErrInvalidPath
	}

	cleaned := filepath.Clean(path)
	if strings.Contains(cleaned, "..") {
		return ErrInvalidPath
	}

	// Reject shell metacharacters and null bytes
	forbidden := []string{"\x00", ";", "|", "&", "$", "`", "\n", "\r", "(", ")", "<", ">"}
	for _, char := range forbidden {
		if strings.Contains(path, char) {
			return ErrInvalidPath
		}
	}

	return nil
}
