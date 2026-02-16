package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/omarluq/hetzner-blackbsd/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testToken = "test_token"

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	configPath := filepath.Join(dir, "blackbsd.yml")
	require.NoError(t, os.WriteFile(configPath, []byte(content), 0o600))
	return configPath
}

func writeSSHKey(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "id_test")
	require.NoError(t, os.WriteFile(keyPath, []byte("fake-key"), 0o600))
	return keyPath
}

func validConfigYAML(keyPath string) string {
	return `hcloud_token: test_token
ssh_key_path: ` + keyPath + `
location: fsn1
server_type: cpx31
branding:
  hostname: blackbsd
  motd: "Welcome to BlackBSD"
  default_user: security
output_iso: true
output_raw: false
build_disk_image: true
`
}

func configYAMLWithoutToken(keyPath string) string {
	return `ssh_key_path: ` + keyPath + `
location: fsn1
server_type: cpx31
branding:
  hostname: blackbsd
  motd: "Welcome to BlackBSD"
  default_user: security
output_iso: true
output_raw: false
build_disk_image: true
`
}

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("loads valid config", func(t *testing.T) {
		t.Parallel()

		keyPath := writeSSHKey(t)
		configPath := writeConfigFile(t, validConfigYAML(keyPath))

		cfg, err := config.Load(configPath)

		require.NoError(t, err)
		assert.Equal(t, keyPath, cfg.SSHKeyPath)
		assert.Equal(t, "fsn1", cfg.Location)
		assert.Equal(t, "cpx31", cfg.ServerType)
		assert.Equal(t, "blackbsd", cfg.Branding.Hostname)
		assert.True(t, cfg.OutputISO)
	})

	t.Run("error for missing file", func(t *testing.T) {
		t.Parallel()

		_, err := config.Load("/nonexistent/blackbsd.yml")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read config")
	})

	t.Run("error for invalid YAML", func(t *testing.T) {
		t.Parallel()

		configPath := writeConfigFile(t, "invalid: yaml: [")

		_, err := config.Load(configPath)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse config")
	})
}

func TestLoadEnvOverride(t *testing.T) {
	t.Setenv("HCLOUD_TOKEN", "env_override_value")

	keyPath := writeSSHKey(t)
	configPath := writeConfigFile(t, configYAMLWithoutToken(keyPath))

	cfg, err := config.Load(configPath)

	require.NoError(t, err)
	assert.Equal(t, "env_override_value", cfg.HCloudToken)
}

func TestDefaults(t *testing.T) {
	t.Parallel()

	cfg := config.Defaults()

	assert.Equal(t, "cpx31", cfg.ServerType)
	assert.Equal(t, "fsn1", cfg.Location)
	assert.Equal(t, "ubuntu-24.04", cfg.Image)
	assert.True(t, cfg.OutputISO)
	assert.False(t, cfg.OutputRaw)
	assert.True(t, cfg.BuildDiskImage)
	assert.Equal(t, "blackbsd", cfg.Branding.Hostname)
	assert.Equal(t, "Welcome to BlackBSD", cfg.Branding.MOTD)
	assert.Equal(t, "security", cfg.Branding.DefaultUser)
}

func TestValidateRequiredFields(t *testing.T) {
	t.Parallel()

	keyPath := writeSSHKey(t)

	t.Run("valid config passes", func(t *testing.T) {
		t.Parallel()

		cfg := config.Defaults()
		cfg.HCloudToken = testToken
		cfg.SSHKeyPath = keyPath
		assert.NoError(t, config.Validate(&cfg))
	})

	t.Run("empty token fails", func(t *testing.T) {
		t.Parallel()

		cfg := config.Defaults()
		cfg.SSHKeyPath = keyPath
		err := config.Validate(&cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "hcloud_token")
	})

	t.Run("empty ssh_key_path fails", func(t *testing.T) {
		t.Parallel()

		cfg := config.Defaults()
		cfg.HCloudToken = testToken
		err := config.Validate(&cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ssh_key_path")
	})

	t.Run("nonexistent ssh_key_path fails", func(t *testing.T) {
		t.Parallel()

		cfg := config.Defaults()
		cfg.HCloudToken = testToken
		cfg.SSHKeyPath = "/nonexistent/key"
		err := config.Validate(&cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})
}

func TestValidateConstraints(t *testing.T) {
	t.Parallel()

	keyPath := writeSSHKey(t)

	t.Run("invalid location fails", func(t *testing.T) {
		t.Parallel()

		cfg := config.Defaults()
		cfg.HCloudToken = testToken
		cfg.SSHKeyPath = keyPath
		cfg.Location = "invalid_dc"
		err := config.Validate(&cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "location")
	})

	t.Run("no output format fails", func(t *testing.T) {
		t.Parallel()

		cfg := config.Defaults()
		cfg.HCloudToken = testToken
		cfg.SSHKeyPath = keyPath
		cfg.OutputISO = false
		cfg.OutputRaw = false
		err := config.Validate(&cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "output")
	})

	t.Run("all valid locations accepted", func(t *testing.T) {
		t.Parallel()

		for _, location := range config.ValidLocations {
			cfg := config.Defaults()
			cfg.HCloudToken = testToken
			cfg.SSHKeyPath = keyPath
			cfg.Location = location
			assert.NoError(t, config.Validate(&cfg), "location %s should be valid", location)
		}
	})
}
