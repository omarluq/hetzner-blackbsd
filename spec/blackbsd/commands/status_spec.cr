require "../../spec_helper"
require "webmock"

describe BlackBSD::Commands::Status do
  before_each { WebMock.reset }

  describe "#run" do
    it "lists all blackbsd servers" do
      WebMock.stub(:get, "https://api.hetzner.cloud/v1/servers")
        .with(headers: {"Authorization" => "Bearer test_token"}, query: {"label_selector" => "managed-by=blackbsd-builder"})
        .to_return(body: {
          servers: [
            {
              id:             1,
              name:           "blackbsd-builder-123",
              status:         "running",
              public_net:     {ipv4: {ip: "1.2.3.4"}},
              rescue_enabled: false,
            },
            {
              id:             2,
              name:           "blackbsd-builder-456",
              status:         "initializing",
              public_net:     {ipv4: {ip: "5.6.7.8"}},
              rescue_enabled: true,
            },
          ],
        }.to_json)

      config = test_config("test_token")
      io = IO::Memory.new
      BlackBSD::Commands::Status.new(config).run(io)
      output = io.to_s

      output.should contain("blackbsd-builder-123")
      output.should contain("running")
      output.should contain("1.2.3.4")
      output.should contain("blackbsd-builder-456")
      output.should contain("initializing")
      output.should contain("Rescue: yes")
    end

    it "shows message when no servers found" do
      WebMock.stub(:get, "https://api.hetzner.cloud/v1/servers")
        .with(headers: {"Authorization" => "Bearer test_token"}, query: {"label_selector" => "managed-by=blackbsd-builder"})
        .to_return(body: {servers: [] of String}.to_json)

      config = test_config("test_token")
      io = IO::Memory.new
      BlackBSD::Commands::Status.new(config).run(io)
      output = io.to_s

      output.should contain("No BlackBSD servers")
    end
  end
end
