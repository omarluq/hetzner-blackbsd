package config

import (
	"os"
	"strings"
)

// ValidLocations lists the Hetzner datacenter locations.
var ValidLocations = []string{"fsn1", "nbg1", "hel1", "ash", "hil", "sin"}

// Validate checks the configuration for required fields and valid values.
func Validate(cfg *Config) error {
	if cfg.HCloudToken == "" {
		return &Error{Field: "hcloud_token", Message: "required (set in config or HCLOUD_TOKEN env)"}
	}

	if cfg.SSHKeyPath == "" {
		return &Error{Field: "ssh_key_path", Message: "required"}
	}

	if _, err := os.Stat(cfg.SSHKeyPath); os.IsNotExist(err) {
		return &Error{Field: "ssh_key_path", Message: "file does not exist: " + cfg.SSHKeyPath}
	}

	if !contains(ValidLocations, cfg.Location) {
		return &Error{
			Field:   "location",
			Message: "must be a valid Hetzner datacenter: " + strings.Join(ValidLocations, ", "),
		}
	}

	if !cfg.OutputISO && !cfg.OutputRaw {
		return &Error{Field: "output_iso/output_raw", Message: "at least one output format must be enabled"}
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
