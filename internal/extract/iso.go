package extract

import (
	"context"
	"fmt"
	"unicode"
)

// ExtractISO creates a bootable ISO from a mounted device partition.
func (e *Extractor) ExtractISO(ctx context.Context, mountPoint, outputPath string) error {
	partition := partitionPath(e.device, 1)

	if err := e.mountDevice(ctx, partition, mountPoint); err != nil {
		return err
	}

	if err := e.createISO(ctx, mountPoint, outputPath); err != nil {
		return err
	}

	return e.unmountDevice(ctx, mountPoint)
}

func (e *Extractor) mountDevice(ctx context.Context, partition, mountPoint string) error {
	cmd := fmt.Sprintf("mount -r %s %s", partition, mountPoint)

	result, err := e.runner.Exec(ctx, cmd)
	if err != nil {
		return fmt.Errorf("mount device: %w", err)
	}

	if !result.Success() {
		return fmt.Errorf("mount device: %s", result.Stderr)
	}

	return nil
}

func (e *Extractor) createISO(ctx context.Context, mountPoint, outputPath string) error {
	cmd := fmt.Sprintf("xorriso -as mkisofs -o %s -b boot/cdboot -no-emul-boot %s", outputPath, mountPoint)

	result, err := e.runner.Exec(ctx, cmd)
	if err != nil {
		return fmt.Errorf("create ISO: %w", err)
	}

	if !result.Success() {
		return fmt.Errorf("create ISO: %s", result.Stderr)
	}

	return nil
}

func (e *Extractor) unmountDevice(ctx context.Context, mountPoint string) error {
	cmd := fmt.Sprintf("umount %s", mountPoint)

	result, err := e.runner.Exec(ctx, cmd)
	if err != nil {
		return fmt.Errorf("unmount device: %w", err)
	}

	if !result.Success() {
		return fmt.Errorf("unmount device: %s", result.Stderr)
	}

	return nil
}

// partitionPath returns the first partition path for a device.
// NVMe devices (ending in digit) use "p" separator: /dev/nvme0n1 -> /dev/nvme0n1p1.
// Traditional devices just append the number: /dev/sda -> /dev/sda1.
func partitionPath(device string, partNum int) string {
	runes := []rune(device)
	if len(runes) > 0 && unicode.IsDigit(runes[len(runes)-1]) {
		return fmt.Sprintf("%sp%d", device, partNum)
	}

	return fmt.Sprintf("%s%d", device, partNum)
}
