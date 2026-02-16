require "../../spec_helper"
require "webmock"

describe BlackBSD::Hetzner::Server do
  describe "#ipv4" do
    it "extracts ipv4 from public_net" do
      json = {
        id:             42,
        name:           "test",
        status:         "running",
        rescue_enabled: false,
        public_net:     {ipv4: {ip: "1.2.3.4"}},
      }.to_json

      server = BlackBSD::Hetzner::Server.from_json(json)
      server.ipv4.should eq "1.2.3.4"
    end
  end

  describe "#rescue_enabled?" do
    it "returns rescue enabled status" do
      json = {
        id:             42,
        name:           "test",
        status:         "running",
        rescue_enabled: true,
        public_net:     {ipv4: {ip: "1.2.3.4"}},
      }.to_json

      server = BlackBSD::Hetzner::Server.from_json(json)
      server.rescue_enabled?.should be_true
    end
  end
end
