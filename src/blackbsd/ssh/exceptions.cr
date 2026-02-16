module BlackBSD
  module SSH
    class SSHError < Exception; end

    class CommandFailedError < SSHError
      getter exit_code : Int32

      def initialize(command : String, @exit_code : Int32, stderr : String = "")
        msg = "Command failed (exit #{@exit_code}): #{command}"
        msg += "\n#{stderr}" unless stderr.empty?
        super(msg)
      end
    end

    class ConnectionError < SSHError; end

    class AuthenticationError < SSHError; end

    class TimeoutError < SSHError; end
  end
end
