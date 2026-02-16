// Package customize applies post-install customizations to a NetBSD host.
package customize

import (
	"github.com/omarluq/hetzner-blackbsd/internal/runner"
)

// Customizer applies branding, packages, and networking to a remote NetBSD host.
type Customizer struct {
	runner runner.Runner
}

// New creates a Customizer backed by the given Runner.
func New(exec runner.Runner) *Customizer {
	return &Customizer{runner: exec}
}
