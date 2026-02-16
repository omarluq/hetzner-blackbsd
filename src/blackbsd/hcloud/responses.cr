require "json"
require "./server"
require "./action"

module BlackBSD
  module Hetzner
    private struct ServerResponse
      include JSON::Serializable
      getter server : Server
    end

    private struct ServersResponse
      include JSON::Serializable
      getter servers : Array(Server)
    end

    private struct ActionResponse
      include JSON::Serializable
      getter action : ActionInfo
    end
  end
end
