module BlackBSD
  module Commands
    class Status
      LABEL = "managed-by=blackbsd-builder"

      def initialize(@config : Config)
      end

      def run(io : IO = STDOUT) : Nil
        client = Hetzner::Client.new(@config.hcloud_token)
        servers = client.list_servers(LABEL)

        if servers.empty?
          io.puts "No BlackBSD servers found."
          return
        end

        io.puts "BlackBSD servers:"
        io.puts ""
        servers.each do |server|
          rescue_str = server.rescue_enabled? ? "yes" : "no"
          io.puts "  #{server.name}"
          io.puts "    ID: #{server.id} | Status: #{server.status} | IP: #{server.ipv4}"
          io.puts "    Rescue: #{rescue_str}"
          io.puts ""
        end
      end
    end
  end
end
