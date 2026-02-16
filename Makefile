BIN_NAME = hetzner-blackbsd

.PHONY: all build release spec test format format-check ameba ameba-fix clean

all: format-check ameba spec

build:
	shards build $(BIN_NAME) -d --error-trace

release:
	shards build $(BIN_NAME) --release --no-debug

spec test:
	crystal spec -v --error-trace

format:
	crystal tool format src spec

format-check:
	crystal tool format --check src spec

ameba:
	bin/ameba

ameba-fix:
	bin/ameba --fix

clean:
	rm -rf bin lib docs coverage .croupier
