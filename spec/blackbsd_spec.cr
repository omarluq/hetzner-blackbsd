require "./spec_helper"

describe BlackBSD do
  it "has a version" do
    BlackBSD::VERSION.should_not be_nil
  end
end
