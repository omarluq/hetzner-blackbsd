require "crest"
require "./exceptions"
require "./server"
require "./action"
require "./responses"

module BlackBSD
  module Hetzner
    class Client
      BASE_URL = "https://api.hetzner.cloud/v1"

      def initialize(@token : String)
      end

      def create_server(name : String, server_type : String, image : String, location : String, ssh_keys : Array(String)) : Server
        body = {
          name:        name,
          server_type: server_type,
          image:       image,
          location:    location,
          ssh_keys:    [ssh_keys],
          labels:      {"managed-by" => "blackbsd-builder"},
        }

        response = post("/servers", body)
        ServerResponse.from_json(response.body).server
      end

      def get_server(id : Int64) : Server?
        response = get("/servers/#{id}")
        return nil if response.status_code == 404
        ServerResponse.from_json(response.body).server
      rescue NotFoundError
        nil
      end

      def list_servers(label_selector : String) : Array(Server)
        response = get("/servers", {"label_selector" => label_selector})
        ServersResponse.from_json(response.body).servers
      end

      def delete_server(id : Int64) : Bool
        response = delete("/servers/#{id}")
        response.status_code != 404
      rescue NotFoundError
        false
      end

      def enable_rescue(server_id : Int64, ssh_key_id : String) : ActionInfo
        body = {type: "linux64", ssh_key_id: ssh_key_id}
        response = post("/servers/#{server_id}/actions/enable_rescue", body)
        ActionResponse.from_json(response.body).action
      end

      def disable_rescue(server_id : Int64) : Bool
        post("/servers/#{server_id}/actions/disable_rescue", nil)
        true
      end

      def get_server_status(id : Int64) : String
        server = get_server(id)
        server.try(&.status) || "unknown"
      end

      private def headers : HTTP::Headers
        HTTP::Headers{
          "Authorization" => "Bearer #{@token}",
          "Content-Type"  => "application/json",
        }
      end

      private def get(path : String) : Crest::Response
        Crest.get("#{BASE_URL}#{path}", headers: headers)
      rescue ex : Crest::RequestFailed
        handle_error(ex)
      end

      private def get(path : String, params : Hash(String, String)) : Crest::Response
        Crest.get("#{BASE_URL}#{path}", headers: headers, params: params)
      rescue ex : Crest::RequestFailed
        handle_error(ex)
      end

      private def post(path : String, body) : Crest::Response
        if body
          Crest.post("#{BASE_URL}#{path}", body.to_json, headers: headers)
        else
          Crest.post("#{BASE_URL}#{path}", headers: headers)
        end
      rescue ex : Crest::RequestFailed
        handle_error(ex)
      end

      private def delete(path : String) : Crest::Response
        Crest.delete("#{BASE_URL}#{path}", headers: headers)
      rescue ex : Crest::RequestFailed
        handle_error(ex)
      end

      private def handle_error(ex : Crest::RequestFailed) : NoReturn
        case ex.http_code
        when 404
          raise NotFoundError.new(ex.message)
        when 429
          raise RateLimitError.new("Rate limited: #{ex.message}")
        else
          raise APIError.new("HTTP #{ex.http_code}: #{ex.message}")
        end
      end
    end
  end
end
