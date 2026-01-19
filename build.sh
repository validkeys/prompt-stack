#!/bin/bash

set -e

VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
COMMIT=${COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")}
BUILD_DATE=${BUILD_DATE:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")}

echo "Building prompt-stack..."
echo "Version: $VERSION"
echo "Commit: $COMMIT"
echo "Build Date: $BUILD_DATE"

mkdir -p dist

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.Date=${BUILD_DATE} -s -w" \
  -o dist/prompt-stack-linux-amd64 ./cmd/prompt-stack

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build \
  -ldflags="-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.Date=${BUILD_DATE} -s -w" \
  -o dist/prompt-stack-darwin-amd64 ./cmd/prompt-stack

CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build \
  -ldflags="-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.Date=${BUILD_DATE} -s -w" \
  -o dist/prompt-stack-darwin-arm64 ./cmd/prompt-stack

echo "Build complete!"
echo "Binaries created:"
ls -lh dist/prompt-stack-*
