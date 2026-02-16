# Test helper utilities for spec files.
module TestHelpers
  def with_tempfile(content : String, &)
    path = "/tmp/blackbsd_test_#{Random::Secure.hex(4)}.yml"
    File.write(path, content)
    begin
      yield path
    ensure
      File.delete(path) if File.exists?(path)
    end
  end

  def with_temp_key(&)
    path = "/tmp/blackbsd_test_key_#{Random::Secure.hex(4)}"
    File.write(path, "dummy key content")
    begin
      yield path
    ensure
      File.delete(path) if File.exists?(path)
    end
  end
end
