require "../../spec_helper"

describe BlackBSD::SSH::Client do
  describe "#initialize" do
    it "stores connection details" do
      client = BlackBSD::SSH::Client.new("example.com", "testuser", key_path: "/path/to/key")
      client.host.should eq "example.com"
      client.user.should eq "testuser"
      client.key_path.should eq "/path/to/key"
      client.port.should eq 22
    end

    it "expands ~ in key path" do
      client = BlackBSD::SSH::Client.new("example.com", "testuser", key_path: "~/test_key")
      client.key_path.should_not start_with("~")
      client.key_path.should contain("/test_key")
    end

    it "accepts custom port" do
      client = BlackBSD::SSH::Client.new("example.com", "testuser", key_path: "/key", port: 2222)
      client.port.should eq 2222
    end
  end
end

describe BlackBSD::SSH::CommandResult do
  describe "#success?" do
    it "returns true when exit_code is 0" do
      result = BlackBSD::SSH::CommandResult.new("output", "", 0)
      result.success?.should be_true
    end

    it "returns false when exit_code is non-zero" do
      result = BlackBSD::SSH::CommandResult.new("", "error", 1)
      result.success?.should be_false
    end
  end
end

describe BlackBSD::SSH::CommandFailedError do
  it "formats message with exit code" do
    error = BlackBSD::SSH::CommandFailedError.new("ls", 1, "No such file")
    error.message.to_s.should contain("exit 1")
    error.message.to_s.should contain("ls")
    error.message.to_s.should contain("No such file")
  end

  it "exposes exit code" do
    error = BlackBSD::SSH::CommandFailedError.new("test", 42)
    error.exit_code.should eq 42
  end
end
