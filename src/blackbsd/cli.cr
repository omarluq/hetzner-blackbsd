module BlackBSD
  module CLI
    def self.run(args = ARGV)
      if args.includes?("--version") || args.includes?("-v")
        puts "hetzner-blackbsd #{VERSION}"
        return
      end

      if args.includes?("--help") || args.includes?("-h") || args.empty?
        puts help_text
        return
      end
    end

    private def self.help_text
      <<-HELP
      hetzner-blackbsd v#{VERSION} - BlackBSD ISO build pipeline on Hetzner Cloud

      Usage:
        hetzner-blackbsd [command] [options]

      Options:
        -h, --help     Show this help
        -v, --version  Show version

      HELP
    end
  end
end
