// Package di provides dependency injection using samber/do v2.
package di

import (
	"context"
	"fmt"

	"github.com/samber/do/v2"
)

// ConfigPathKey is the DI key for the config file path.
const ConfigPathKey = "config.path"

// Container wraps the do.Injector with blackbsd-specific configuration.
type Container struct {
	injector *do.RootScope
}

// NewContainer creates and configures the DI container.
func NewContainer(configPath string) (*Container, error) {
	injector := do.New()

	do.ProvideNamedValue(injector, ConfigPathKey, configPath)
	RegisterProviders(injector)

	if _, err := do.Invoke[*ConfigService](injector); err != nil {
		return nil, err
	}

	return &Container{injector: injector}, nil
}

// Injector returns the underlying do.Injector.
func (c *Container) Injector() *do.RootScope {
	return c.injector
}

// Invoke resolves a service from the container.
func Invoke[T any](c *Container) (T, error) {
	return do.Invoke[T](c.injector)
}

// MustInvoke resolves a service or panics.
func MustInvoke[T any](c *Container) T {
	return do.MustInvoke[T](c.injector)
}

// Shutdown gracefully shuts down all services.
func (c *Container) Shutdown() error {
	report := c.injector.Shutdown()
	if report != nil && !report.Succeed {
		return fmt.Errorf("shutdown failed: %s", report.Error())
	}
	return nil
}

// ShutdownWithContext shuts down with context.
func (c *Container) ShutdownWithContext(ctx context.Context) error {
	done := make(chan *do.ShutdownReport, 1)
	go func() {
		done <- c.injector.ShutdownWithContext(ctx)
	}()

	select {
	case report := <-done:
		if report != nil && !report.Succeed {
			return fmt.Errorf("shutdown failed: %s", report.Error())
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("shutdown timed out: %w", ctx.Err())
	}
}
