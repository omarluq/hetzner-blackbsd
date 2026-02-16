package di

import "github.com/samber/do/v2"

// RegisterProviders registers all service providers as singletons.
func RegisterProviders(i do.Injector) {
	do.Provide(i, NewConfigService)
	do.Provide(i, NewLoggerService)
	do.Provide(i, NewHetznerService)
	do.Provide(i, NewSSHService)
}
