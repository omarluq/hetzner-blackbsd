module BlackBSD
  module Hetzner
    class APIError < Exception; end

    class RateLimitError < APIError; end

    class NotFoundError < APIError; end
  end
end
