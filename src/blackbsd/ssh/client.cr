require "ssh2"
require "./exceptions"
require "./result"

module BlackBSD
  module SSH
    class Client
      getter host : String
      getter user : String
      getter key_path : String
      getter port : Int32

      def initialize(@host : String, @user : String, *, key_path : String, @port : Int32 = 22)
        @key_path = File.expand_path(key_path, home: true)
      end

      def exec(command : String) : String
        result = exec_capture(command)
        unless result.success?
          raise CommandFailedError.new(command, result.exit_code, result.stderr)
        end
        result.stdout
      end

      def exec_capture(command : String) : CommandResult
        stdout = IO::Memory.new
        stderr = IO::Memory.new
        exit_code = 0

        SSH2::Session.open(@host, @port) do |session|
          session.login_with_pubkey(@user, @key_path)
          session.open_session do |channel|
            channel.command(command)
            IO.copy(channel, stdout)
            exit_code = channel.exit_status
          end
        end

        CommandResult.new(stdout.to_s.chomp, stderr.to_s.chomp, exit_code)
      rescue ex : Socket::ConnectError
        raise ConnectionError.new("Cannot connect to #{@host}:#{@port}: #{ex.message}")
      rescue ex : SSH2::SessionError
        raise AuthenticationError.new("Auth failed for #{@user}@#{@host}: #{ex.message}")
      end

      def upload(local_path : String, remote_path : String) : Nil
        content = File.read(local_path)
        SSH2::Session.open(@host, @port) do |session|
          session.login_with_pubkey(@user, @key_path)
          session.sftp_session do |sftp|
            sftp.open_file(remote_path, "w") do |file|
              file.print(content)
            end
          end
        end
      rescue ex : Socket::ConnectError
        raise ConnectionError.new("Cannot connect to #{@host}:#{@port}: #{ex.message}")
      end

      def download(remote_path : String, local_path : String) : Nil
        SSH2::Session.open(@host, @port) do |session|
          session.login_with_pubkey(@user, @key_path)
          session.sftp_session do |sftp|
            sftp.open_file(remote_path, "r") do |file|
              File.write(local_path, file.gets_to_end)
            end
          end
        end
      rescue ex : Socket::ConnectError
        raise ConnectionError.new("Cannot connect to #{@host}:#{@port}: #{ex.message}")
      end

      def wait_for_ready(timeout : Int32 = 120, interval : Int32 = 5) : Nil
        deadline = Time.utc + timeout.seconds
        last_error : Exception? = nil

        while Time.utc < deadline
          begin
            exec("echo ready")
            return
          rescue ex : SSHError
            last_error = ex
            sleep interval.seconds
          end
        end

        raise TimeoutError.new("SSH not ready after #{timeout}s: #{last_error.try(&.message)}")
      end
    end
  end
end
