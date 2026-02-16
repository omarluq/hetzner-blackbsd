require "yaml"

module BlackBSD
  class ConfigError < Exception; end

  struct Branding
    include YAML::Serializable

    property hostname : String
    property motd : String
    property default_user : String
  end

  struct Config
    include YAML::Serializable

    property hcloud_token : String
    property ssh_key_path : String
    property location : String = "fsn1"
    property server_type : String = "cpx31"

    property netbsd_version : String = "10.1"
    property netbsd_arch : String = "amd64"

    property security_tools : Array(String) = [] of String

    property branding : Branding

    property output_dir : String = "./output"

    @[YAML::Field(key: "build_disk_image")]
    property? build_disk_image : Bool = true

    @[YAML::Field(key: "build_iso")]
    property? build_iso : Bool = true

    @[YAML::Field(key: "upload_to_github")]
    property? upload_to_github : Bool = false

    @[YAML::Field(key: "deploy_test_vm")]
    property? deploy_test_vm : Bool = false

    def self.from_file(path : String) : self
      unless File.exists?(path)
        raise ConfigError.new("Config file not found: #{path}")
      end

      env_token = ENV["HCLOUD_TOKEN"]?

      content = File.read(path)
      config = self.from_yaml(content)

      if token = env_token
        config = config.with_token(token)
      end

      config.validate!
      config
    rescue ex : YAML::ParseException
      raise ConfigError.new("Invalid YAML in #{path}: #{ex.message}")
    end

    protected def with_token(token : String) : self
      copy = self.dup
      copy.hcloud_token = token
      copy
    end

    def validate! : Nil
      errors = [] of String

      errors << "hcloud_token is required" if @hcloud_token.empty? || @hcloud_token == "your_token_here"
      errors << "ssh_key_path is required" if @ssh_key_path.empty?
      errors << "server_type is required" if @server_type.empty?
      errors << "netbsd_version is required" if @netbsd_version.empty?

      unless File.exists?(Path[@ssh_key_path].expand(home: true))
        errors << "ssh_key_path does not exist: #{@ssh_key_path}"
      end

      valid_locations = ["fsn1", "nbg1", "hel1", "ash", "hil"]
      unless valid_locations.includes?(@location)
        errors << "Invalid location: #{@location} (valid: #{valid_locations.join(", ")})"
      end

      unless build_disk_image? || build_iso?
        errors << "At least one of build_disk_image or build_iso must be true"
      end

      unless errors.empty?
        raise ConfigError.new("Config validation failed:\n  - #{errors.join("\n  - ")}")
      end
    end
  end
end
