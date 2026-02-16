package extract

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// ImageSize returns the size of a remote image file in bytes.
func (e *Extractor) ImageSize(ctx context.Context, imagePath string) (int64, error) {
	cmd := fmt.Sprintf("stat -c %%s %s", imagePath)

	result, err := e.runner.Exec(ctx, cmd)
	if err != nil {
		return 0, fmt.Errorf("get image size: %w", err)
	}

	if !result.Success() {
		return 0, fmt.Errorf("get image size: %s", result.Stderr)
	}

	sizeStr := strings.TrimSpace(result.Stdout)

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse image size: %w", err)
	}

	return size, nil
}

// Checksum returns the SHA256 checksum of a remote image file.
func (e *Extractor) Checksum(ctx context.Context, imagePath string) (string, error) {
	cmd := fmt.Sprintf("sha256sum %s", imagePath)

	result, err := e.runner.Exec(ctx, cmd)
	if err != nil {
		return "", fmt.Errorf("compute checksum: %w", err)
	}

	if !result.Success() {
		return "", fmt.Errorf("compute checksum: %s", result.Stderr)
	}

	output := strings.TrimSpace(result.Stdout)
	parts := strings.Fields(output)

	if len(parts) == 0 {
		return "", fmt.Errorf("parse checksum: empty output")
	}

	return parts[0], nil
}
