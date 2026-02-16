require "../../spec_helper"
require "webmock"

describe BlackBSD::Hetzner::Client do
  before_each { WebMock.reset }

  describe "#create_server" do
    it "creates a server with correct payload" do
      WebMock.stub(:post, "https://api.hetzner.cloud/v1/servers")
        .with(headers: {"Authorization" => "Bearer test_token"}, body: {
          name:        "test-server",
          server_type: "cpx31",
          image:       "ubuntu-24.04",
          location:    "fsn1",
          ssh_keys:    [%w[key_id]],
          labels:      {"managed-by" => "blackbsd-builder"},
        }.to_json)
        .to_return(body: {
          server: {
            id:             42,
            name:           "test-server",
            status:         "running",
            public_net:     {ipv4: {ip: "1.2.3.4"}},
            rescue_enabled: false,
          },
        }.to_json)

      client = BlackBSD::Hetzner::Client.new("test_token")
      server = client.create_server("test-server", "cpx31", "ubuntu-24.04", "fsn1", ["key_id"])

      server.id.should eq 42
      server.name.should eq "test-server"
      server.status.should eq "running"
      server.ipv4.should eq "1.2.3.4"
      server.rescue_enabled?.should be_false
    end
  end

  describe "#get_server" do
    it "fetches server by id" do
      WebMock.stub(:get, "https://api.hetzner.cloud/v1/servers/42")
        .with(headers: {"Authorization" => "Bearer test_token"})
        .to_return(body: {
          server: {
            id:             42,
            name:           "test-server",
            status:         "running",
            public_net:     {ipv4: {ip: "1.2.3.4"}},
            rescue_enabled: false,
          },
        }.to_json)

      client = BlackBSD::Hetzner::Client.new("test_token")
      server = client.get_server(42)

      server.should_not be_nil
      server.as(BlackBSD::Hetzner::Server).id.should eq 42
    end

    it "returns nil for 404" do
      WebMock.stub(:get, "https://api.hetzner.cloud/v1/servers/42")
        .with(headers: {"Authorization" => "Bearer test_token"})
        .to_return(status: 404)

      client = BlackBSD::Hetzner::Client.new("test_token")
      server = client.get_server(42)

      server.should be_nil
    end
  end

  describe "#list_servers" do
    it "lists servers with label filter" do
      WebMock.stub(:get, "https://api.hetzner.cloud/v1/servers")
        .with(headers: {"Authorization" => "Bearer test_token"}, query: {"label_selector" => "managed-by=blackbsd-builder"})
        .to_return(body: {
          servers: [
            {id: 1, name: "server-1", status: "running", public_net: {ipv4: {ip: "1.2.3.4"}}, rescue_enabled: false},
            {id: 2, name: "server-2", status: "running", public_net: {ipv4: {ip: "5.6.7.8"}}, rescue_enabled: true},
          ],
        }.to_json)

      client = BlackBSD::Hetzner::Client.new("test_token")
      servers = client.list_servers("managed-by=blackbsd-builder")

      servers.size.should eq 2
      servers[0].id.should eq 1
      servers[1].id.should eq 2
      servers[1].rescue_enabled?.should be_true
    end
  end

  describe "#delete_server" do
    it "deletes server and returns true" do
      WebMock.stub(:delete, "https://api.hetzner.cloud/v1/servers/42")
        .with(headers: {"Authorization" => "Bearer test_token"})
        .to_return(status: 200)

      client = BlackBSD::Hetzner::Client.new("test_token")
      result = client.delete_server(42)

      result.should be_true
    end

    it "returns false for 404" do
      WebMock.stub(:delete, "https://api.hetzner.cloud/v1/servers/42")
        .with(headers: {"Authorization" => "Bearer test_token"})
        .to_return(status: 404)

      client = BlackBSD::Hetzner::Client.new("test_token")
      result = client.delete_server(42)

      result.should be_false
    end
  end

  describe "#enable_rescue" do
    it "enables rescue mode" do
      WebMock.stub(:post, "https://api.hetzner.cloud/v1/servers/42/actions/enable_rescue")
        .with(headers: {"Authorization" => "Bearer test_token"}, body: {
          type:       "linux64",
          ssh_key_id: "key_123",
        }.to_json)
        .to_return(body: {
          action: {
            id:      1,
            status:  "running",
            command: "enable_rescue",
            started: "2026-02-16T00:00:00Z",
          },
        }.to_json)

      client = BlackBSD::Hetzner::Client.new("test_token")
      rescue_info = client.enable_rescue(42, "key_123")

      rescue_info.id.should eq 1
      rescue_info.status.should eq "running"
    end
  end

  describe "#disable_rescue" do
    it "disables rescue mode" do
      WebMock.stub(:post, "https://api.hetzner.cloud/v1/servers/42/actions/disable_rescue")
        .with(headers: {"Authorization" => "Bearer test_token"})
        .to_return(status: 200)

      client = BlackBSD::Hetzner::Client.new("test_token")
      result = client.disable_rescue(42)

      result.should be_true
    end
  end

  describe "#get_server_status" do
    it "returns server status" do
      WebMock.stub(:get, "https://api.hetzner.cloud/v1/servers/42")
        .with(headers: {"Authorization" => "Bearer test_token"})
        .to_return(body: {
          server: {
            id:             42,
            status:         "running",
            public_net:     {ipv4: {ip: "1.2.3.4"}},
            rescue_enabled: false,
          },
        }.to_json)

      client = BlackBSD::Hetzner::Client.new("test_token")
      status = client.get_server_status(42)

      status.should eq "running"
    end
  end
end
