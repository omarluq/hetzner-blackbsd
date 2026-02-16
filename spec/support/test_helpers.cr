# Test helper utilities for spec files.
module TestHelpers
  def with_tempfile(content : String, &)
    path = "/tmp/blackbsd_test_#{Random::Secure.hex(4)}.yml"
    File.write(path, content)
    begin
      yield path
    ensure
      File.delete(path) if File.exists?(path)
    end
  end

  def with_temp_key(&)
    path = "/tmp/blackbsd_test_key_#{Random::Secure.hex(4)}"
    File.write(path, "dummy key content")
    begin
      yield path
    ensure
      File.delete(path) if File.exists?(path)
    end
  end

  def test_config(token : String = "test_token") : BlackBSD::Config
    key_path = "/tmp/blackbsd_test_key_#{Random::Secure.hex(4)}"
    File.write(key_path, "dummy key")

    yaml_path = "/tmp/blackbsd_test_#{Random::Secure.hex(4)}.yml"
    yaml = <<-YAML
      hcloud_token: #{token}
      ssh_key_path: #{key_path}
      branding:
        hostname: test
        motd: test
        default_user: test
      YAML
    File.write(yaml_path, yaml)

    ENV.delete("HCLOUD_TOKEN")
    config = BlackBSD::Config.from_file(yaml_path)

    File.delete(yaml_path) if File.exists?(yaml_path)
    File.delete(key_path) if File.exists?(key_path)

    config
  end
end
