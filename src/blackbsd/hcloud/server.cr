require "json"

module BlackBSD
  module Hetzner
    struct Server
      include JSON::Serializable

      getter id : Int64
      getter name : String = ""
      getter status : String = ""
      @[JSON::Field(key: "rescue_enabled")]
      getter? rescue_enabled : Bool = false

      @[JSON::Field(key: "public_net")]
      getter public_net : PublicNet = PublicNet.new

      def ipv4 : String
        public_net.ipv4.ip
      end
    end

    struct PublicNet
      include JSON::Serializable

      getter ipv4 : IPv4 = IPv4.new

      def initialize
        @ipv4 = IPv4.new
      end
    end

    struct IPv4
      include JSON::Serializable

      getter ip : String = ""

      def initialize
        @ip = ""
      end
    end
  end
end
