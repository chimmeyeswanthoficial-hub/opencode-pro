# OpenCode Pro Makefile

.PHONY: all build test clean run

all: build

build:
	go build -ldflags="-s -w" -o opencode .

test:
	go test -v ./internal/optimizer/...

clean:
	rm -f opencode

run: build
	./opencode
