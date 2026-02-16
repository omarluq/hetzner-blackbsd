// Package extract provides disk image extraction utilities for BlackBSD builds.
package extract

import (
	"github.com/omarluq/hetzner-blackbsd/internal/runner"
)

// Extractor handles disk image extraction from remote devices.
type Extractor struct {
	runner runner.Runner
	device string
}

// New creates a new Extractor.
func New(exec runner.Runner, device string) *Extractor {
	return &Extractor{
		runner: exec,
		device: device,
	}
}
