require "json"

module BlackBSD
  module Hetzner
    struct ActionInfo
      include JSON::Serializable

      getter id : Int64
      getter status : String = ""
      getter command : String = ""
      getter started : String = ""
    end
  end
end
