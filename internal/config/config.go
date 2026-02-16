// Package config defines the configuration model for blackbsd.
package config

// Config is the root configuration for blackbsd.
type Config struct {
	Branding       Branding `yaml:"branding"`
	HCloudToken    string   `yaml:"hcloud_token"`
	SSHKeyPath     string   `yaml:"ssh_key_path"`
	ServerType     string   `yaml:"server_type"`
	Location       string   `yaml:"location"`
	Image          string   `yaml:"image"`
	OutputISO      bool     `yaml:"output_iso"`
	OutputRaw      bool     `yaml:"output_raw"`
	BuildDiskImage bool     `yaml:"build_disk_image"`
}

// Branding holds the customization settings for the built image.
type Branding struct {
	Hostname    string `yaml:"hostname"`
	MOTD        string `yaml:"motd"`
	DefaultUser string `yaml:"default_user"`
}

// Defaults returns a Config populated with sensible default values.
func Defaults() Config {
	return Config{
		HCloudToken:    "",
		SSHKeyPath:     "",
		ServerType:     "cpx31",
		Location:       "fsn1",
		Image:          "ubuntu-24.04",
		OutputISO:      true,
		OutputRaw:      false,
		BuildDiskImage: true,
		Branding: Branding{
			Hostname:    "blackbsd",
			MOTD:        "Welcome to BlackBSD",
			DefaultUser: "security",
		},
	}
}
