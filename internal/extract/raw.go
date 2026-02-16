package extract

import (
	"context"
	"fmt"
)

// ExtractRawImage creates a compressed raw disk image using dd and xz.
func (e *Extractor) ExtractRawImage(ctx context.Context, outputPath string) error {
	cmd := fmt.Sprintf("dd if=%s bs=4M status=progress | xz -T0 -9 > %s", e.device, outputPath)

	result, err := e.runner.Exec(ctx, cmd)
	if err != nil {
		return fmt.Errorf("extract raw image: %w", err)
	}

	if !result.Success() {
		return fmt.Errorf("extract raw image: %s", result.Stderr)
	}

	return nil
}
