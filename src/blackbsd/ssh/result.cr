module BlackBSD
  module SSH
    struct CommandResult
      getter stdout : String
      getter stderr : String
      getter exit_code : Int32

      def initialize(@stdout : String, @stderr : String, @exit_code : Int32)
      end

      def success? : Bool
        @exit_code == 0
      end
    end
  end
end
