package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Load reads and parses a YAML config file, applying defaults
// and environment variable overrides.
func Load(path string) (*Config, error) {
	cleanPath := filepath.Clean(path)
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config %s: %w", path, err)
	}

	cfg := Defaults()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config %s: %w", path, err)
	}

	// Environment variable override for token
	if envToken := os.Getenv("HCLOUD_TOKEN"); envToken != "" {
		cfg.HCloudToken = envToken
	}

	if err := Validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
