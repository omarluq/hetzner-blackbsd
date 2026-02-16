package di

import (
	"fmt"
	"log/slog"

	hcloudsdk "github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"

	"github.com/omarluq/hetzner-blackbsd/internal/config"
	"github.com/omarluq/hetzner-blackbsd/internal/logger"
)

// ConfigService wraps the loaded configuration.
type ConfigService struct {
	Config *config.Config
	path   string
}

// NewConfigService loads config from the provided path.
func NewConfigService(i do.Injector) (*ConfigService, error) {
	path := do.MustInvokeNamed[string](i, ConfigPathKey)
	cfg, err := config.Load(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", path, err)
	}
	return &ConfigService{Config: cfg, path: path}, nil
}

// LoggerService provides the configured logger.
type LoggerService struct {
	Logger zerolog.Logger
}

// NewLoggerService creates the logger service.
func NewLoggerService(i do.Injector) (*LoggerService, error) {
	cfgSvc := do.MustInvoke[*ConfigService](i)

	zl := logger.New("info", nil)
	logger.SetupSlog(&zl)

	slog.Info("logger initialized", "path", cfgSvc.path)

	return &LoggerService{Logger: zl}, nil
}

// HetznerService wraps the Hetzner API client.
type HetznerService struct {
	Client *hcloudsdk.Client
}

// NewHetznerService creates the Hetzner API client.
func NewHetznerService(i do.Injector) (*HetznerService, error) {
	cfgSvc := do.MustInvoke[*ConfigService](i)

	client := hcloudsdk.NewClient(
		hcloudsdk.WithToken(cfgSvc.Config.HCloudToken),
	)

	slog.Info("hetzner client initialized", "location", cfgSvc.Config.Location)

	return &HetznerService{Client: client}, nil
}

// SSHService is a placeholder for future SSH operations.
type SSHService struct{}

// NewSSHService creates the SSH service.
func NewSSHService(_ do.Injector) (*SSHService, error) {
	return &SSHService{}, nil
}
