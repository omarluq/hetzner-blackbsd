package extract

import (
	"context"
	"fmt"

	"github.com/omarluq/hetzner-blackbsd/internal/ssh"
)

// ExtractRawImage creates a compressed raw disk image using dd and xz.
func (e *Extractor) ExtractRawImage(ctx context.Context, outputPath string) error {
	if err := ValidatePath(outputPath); err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}

	cmd := fmt.Sprintf("dd if=%s bs=4M status=progress | xz -T0 -9 > %s",
		ssh.EscapeShellArg(e.device),
		ssh.EscapeShellArg(outputPath))

	result, err := e.runner.Exec(ctx, cmd)
	if err != nil {
		return fmt.Errorf("extract raw image: %w", err)
	}

	if !result.Success() {
		return fmt.Errorf("extract raw image: %s", result.Stderr)
	}

	return nil
}
