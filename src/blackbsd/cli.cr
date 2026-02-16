require "admiral"

module BlackBSD
  class CLI < Admiral::Command
    define_help description: "BlackBSD ISO build pipeline on Hetzner Cloud"
    define_version BlackBSD::VERSION

    class StatusCmd < Admiral::Command
      define_help description: "Show BlackBSD build servers"

      define_flag config : String,
        description: "Path to config file",
        default: "blackbsd.yml",
        short: c

      def run
        config = BlackBSD::Config.from_file(flags.config)
        Commands::Status.new(config).run
      rescue ex : ConfigError
        STDERR.puts "Error: #{ex.message}"
        exit 1
      end
    end

    register_sub_command status, StatusCmd, description: "Show BlackBSD build servers"

    def run
      puts help
    end
  end
end
