.PHONY: build test fmt clean lint help

build:
	go build -o dist/your-tool ./cmd/your-tool

test:
	go test -v ./...

fmt:
	go fmt ./...

lint:
	go vet ./...

clean:
	rm -rf dist/
	go clean

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build  - Build the binary"
	@echo "  test   - Run all tests"
	@echo "  fmt    - Format Go code"
	@echo "  lint   - Run go vet"
	@echo "  clean  - Clean build artifacts"
	@echo "  help   - Show this help message"
