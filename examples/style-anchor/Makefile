.PHONY: all build test install clean

BINARY=bin/mytool

all: build

build:
	@echo "==> building $(BINARY)"
	@mkdir -p bin
	@go build -o $(BINARY) ./cmd/mytool

test:
	@echo "==> running tests"
	@go test ./...

install: build
	@echo "==> installing to $(HOME)/.local/bin"
	@mkdir -p $(HOME)/.local/bin
	@cp $(BINARY) $(HOME)/.local/bin/

clean:
	@echo "==> cleaning"
	@rm -rf bin
