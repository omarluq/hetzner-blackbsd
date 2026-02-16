require "../../spec_helper"

describe BlackBSD::Config do
  describe ".from_file" do
    it "loads a valid config file" do
      with_temp_key do |key_path|
        config_yaml = <<-YAML
          hcloud_token: test_token
          ssh_key_path: #{key_path}
          location: fsn1
          server_type: cpx31
          netbsd_version: "10.1"
          netbsd_arch: "amd64"
          security_tools:
            - nmap
            - wireshark
          branding:
            hostname: blackbsd
            motd: "Welcome"
            default_user: hacker
          output_dir: ./output
          build_disk_image: true
          build_iso: true
          YAML

        with_tempfile(config_yaml) do |path|
          ENV.delete("HCLOUD_TOKEN")
          config = BlackBSD::Config.from_file(path)
          config.hcloud_token.should eq "test_token"
          config.ssh_key_path.should eq key_path
          config.location.should eq "fsn1"
          config.server_type.should eq "cpx31"
          config.netbsd_version.should eq "10.1"
          config.security_tools.should eq ["nmap", "wireshark"]
          config.branding.hostname.should eq "blackbsd"
          config.build_disk_image?.should be_true
        end
      end
    end

    it "uses HCLOUD_TOKEN env var as override" do
      with_temp_key do |key_path|
        config_yaml = <<-YAML
          hcloud_token: file_token
          ssh_key_path: #{key_path}
          branding:
            hostname: test
            motd: test
            default_user: test
          YAML

        with_tempfile(config_yaml) do |path|
          ENV["HCLOUD_TOKEN"] = "env_token"
          begin
            config = BlackBSD::Config.from_file(path)
            config.hcloud_token.should eq "env_token"
          ensure
            ENV.delete("HCLOUD_TOKEN")
          end
        end
      end
    end

    it "raises error for missing file" do
      expect_raises(BlackBSD::ConfigError, "Config file not found") do
        BlackBSD::Config.from_file("/nonexistent.yml")
      end
    end

    it "raises error for invalid YAML" do
      with_tempfile("invalid: yaml: [") do |path|
        expect_raises(BlackBSD::ConfigError, /Invalid YAML/) do
          BlackBSD::Config.from_file(path)
        end
      end
    end

    it "validates required fields" do
      with_temp_key do |key_path|
        config_yaml = <<-YAML
          hcloud_token: ""
          ssh_key_path: #{key_path}
          branding:
            hostname: test
            motd: test
            default_user: test
          YAML

        with_tempfile(config_yaml) do |path|
          expect_raises(BlackBSD::ConfigError, /hcloud_token is required/) do
            BlackBSD::Config.from_file(path)
          end
        end
      end
    end

    it "validates location is a Hetzner datacenter" do
      with_temp_key do |key_path|
        config_yaml = <<-YAML
          hcloud_token: test
          ssh_key_path: #{key_path}
          location: invalid_dc
          branding:
            hostname: test
            motd: test
            default_user: test
          YAML

        with_tempfile(config_yaml) do |path|
          expect_raises(BlackBSD::ConfigError, /Invalid location/) do
            BlackBSD::Config.from_file(path)
          end
        end
      end
    end

    it "validates at least one output format is enabled" do
      with_temp_key do |key_path|
        config_yaml = <<-YAML
          hcloud_token: test
          ssh_key_path: #{key_path}
          build_disk_image: false
          build_iso: false
          branding:
            hostname: test
            motd: test
            default_user: test
          YAML

        with_tempfile(config_yaml) do |path|
          expect_raises(BlackBSD::ConfigError, /At least one of/) do
            BlackBSD::Config.from_file(path)
          end
        end
      end
    end

    it "provides defaults for optional fields" do
      with_temp_key do |key_path|
        config_yaml = <<-YAML
          hcloud_token: test
          ssh_key_path: #{key_path}
          branding:
            hostname: test
            motd: test
            default_user: test
          YAML

        with_tempfile(config_yaml) do |path|
          config = BlackBSD::Config.from_file(path)
          config.location.should eq "fsn1"
          config.server_type.should eq "cpx31"
          config.netbsd_version.should eq "10.1"
          config.netbsd_arch.should eq "amd64"
          config.security_tools.should eq [] of String
          config.output_dir.should eq "./output"
          config.build_disk_image?.should be_true
          config.build_iso.should be_true
        end
      end
    end
  end
end
