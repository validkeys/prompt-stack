# Build and Deployment Guide

**Last Updated**: 2026-01-07  
**Go Version**: 1.24.0

This guide provides comprehensive instructions for building, packaging, and deploying PromptStack across different platforms.

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Local Development Build](#local-development-build)
- [Production Builds](#production-builds)
  - [macOS](#macos)
  - [Linux](#linux)
  - [Windows](#windows)
- [Build Optimization](#build-optimization)
- [Versioning](#versioning)
- [Installation Methods](#installation-methods)
  - [Homebrew](#homebrew)
  - [Linux Packages](#linux-packages)
  - [Windows Installer](#windows-installer)
  - [Manual Installation](#manual-installation)
- [Release Process](#release-process)
- [Update Mechanism](#update-mechanism)
- [CI/CD Integration](#cicd-integration)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

Before building PromptStack, ensure you have:

- **Go 1.24.0 or later** - [Install Go](https://golang.org/dl/)
- **Git** - For version control
- **Make** - For build automation (optional, but recommended)
- **Platform-specific tools**:
  - **macOS**: Xcode Command Line Tools (`xcode-select --install`)
  - **Linux**: Build tools (`sudo apt-get install build-essential` on Ubuntu)
  - **Windows**: MinGW-w64 or similar

### Verify Go Installation

```bash
go version
# Should output: go version go1.24.0 darwin/amd64 (or similar)

go env GOPATH
go env GOROOT
```

---

## Local Development Build

For development and testing, build the binary for your current platform:

### Quick Build

```bash
# Build for current platform
go build -o promptstack ./cmd/promptstack

# Run the binary
./promptstack
```

### Development Build with Debug Info

```bash
# Build with debug symbols and optimizations disabled
go build -gcflags="all=-N -l" -o promptstack-debug ./cmd/promptstack

# Build with race detector support
go build -race -o promptstack-race ./cmd/promptstack
```

### Clean Build

```bash
# Clean build artifacts
go clean

# Remove cached modules
go clean -modcache

# Full clean and rebuild
go clean && go build -o promptstack ./cmd/promptstack
```

### Run Tests Before Building

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...
```

---

## Production Builds

For production releases, build optimized binaries for all target platforms.

### Build Matrix

| Platform | Architecture | Binary Name |
|----------|--------------|-------------|
| macOS | amd64 (Intel) | `promptstack-darwin-amd64` |
| macOS | arm64 (Apple Silicon) | `promptstack-darwin-arm64` |
| Linux | amd64 | `promptstack-linux-amd64` |
| Linux | arm64 | `promptstack-linux-arm64` |
| Windows | amd64 | `promptstack-windows-amd64.exe` |

### macOS

#### Intel (amd64)

```bash
# Set environment variables
export GOOS=darwin
export GOARCH=amd64
export CGO_ENABLED=0

# Build optimized binary
go build -ldflags="-s -w" -o promptstack-darwin-amd64 ./cmd/promptstack

# Verify binary
file promptstack-darwin-amd64
# Output: promptstack-darwin-amd64: Mach-O 64-bit executable x86_64
```

#### Apple Silicon (arm64)

```bash
# Set environment variables
export GOOS=darwin
export GOARCH=arm64
export CGO_ENABLED=0

# Build optimized binary
go build -ldflags="-s -w" -o promptstack-darwin-arm64 ./cmd/promptstack

# Verify binary
file promptstack-darwin-arm64
# Output: promptstack-darwin-arm64: Mach-O 64-bit executable arm64
```

#### Universal Binary (Intel + Apple Silicon)

```bash
# Build both architectures first
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o promptstack-darwin-amd64 ./cmd/promptstack
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o promptstack-darwin-arm64 ./cmd/promptstack

# Create universal binary
lipo -create -output promptstack-darwin-universal promptstack-darwin-amd64 promptstack-darwin-arm64

# Verify universal binary
file promptstack-darwin-universal
# Output: promptstack-darwin-universal: Mach-O universal binary with 2 architectures
```

### Linux

#### amd64

```bash
# Set environment variables
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

# Build optimized binary
go build -ldflags="-s -w" -o promptstack-linux-amd64 ./cmd/promptstack

# Verify binary
file promptstack-linux-amd64
# Output: promptstack-linux-amd64: ELF 64-bit LSB executable, x86-64
```

#### arm64

```bash
# Set environment variables
export GOOS=linux
export GOARCH=arm64
export CGO_ENABLED=0

# Build optimized binary
go build -ldflags="-s -w" -o promptstack-linux-arm64 ./cmd/promptstack

# Verify binary
file promptstack-linux-arm64
# Output: promptstack-linux-arm64: ELF 64-bit LSB executable, ARM aarch64
```

### Windows

#### amd64

```bash
# Set environment variables (PowerShell)
$env:GOOS="windows"
$env:GOARCH="amd64"
$env:CGO_ENABLED="0"

# Build optimized binary
go build -ldflags="-s -w" -o promptstack-windows-amd64.exe ./cmd/promptstack

# Verify binary (PowerShell)
Get-FileItem promptstack-windows-amd64.exe
```

Or using Command Prompt:

```cmd
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0

go build -ldflags="-s -w" -o promptstack-windows-amd64.exe ./cmd/promptstack
```

---

## Build Optimization

### Linker Flags

Use these linker flags to reduce binary size:

```bash
# -s: Omit symbol table and debug information
# -w: Omit DWARF symbol table
go build -ldflags="-s -w" -o promptstack ./cmd/promptstack
```

### Build Tags

Use build tags to conditionally compile code:

```bash
# Build with specific tags
go build -tags="production" -o promptstack ./cmd/promptstack

# Build with debug features
go build -tags="debug" -o promptstack ./cmd/promptstack
```

### Upx Compression (Optional)

Further reduce binary size using UPX compression:

```bash
# Install UPX
brew install upx  # macOS
sudo apt-get install upx  # Linux

# Compress binary
upx --best --lzma promptstack-linux-amd64

# Verify compressed binary
upx -t promptstack-linux-amd64
```

**Note**: UPX compression can increase startup time slightly. Test performance before using in production.

---

## Versioning

### Version Information

Embed version information into the binary using linker flags:

```bash
# Set version variables
VERSION=1.0.0
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(git rev-parse --short HEAD)

# Build with version info
go build -ldflags="-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT" -o promptstack ./cmd/promptstack
```

### Display Version

Add version display to `main.go`:

```go
package main

import (
    "fmt"
    "os"
)

var (
    Version   = "dev"
    BuildTime = "unknown"
    GitCommit = "unknown"
)

func main() {
    if len(os.Args) > 1 && os.Args[1] == "--version" {
        fmt.Printf("PromptStack %s\n", Version)
        fmt.Printf("Built: %s\n", BuildTime)
        fmt.Printf("Commit: %s\n", GitCommit)
        os.Exit(0)
    }
    
    // Rest of application
}
```

### Semantic Versioning

Follow semantic versioning (MAJOR.MINOR.PATCH):

- **MAJOR**: Incompatible API changes
- **MINOR**: Backwards-compatible functionality additions
- **PATCH**: Backwards-compatible bug fixes

Examples:
- `1.0.0` - Initial stable release
- `1.1.0` - New feature addition
- `1.1.1` - Bug fix
- `2.0.0` - Breaking changes

---

## Installation Methods

### Homebrew

#### Create Homebrew Formula

Create `Formula/promptstack.rb`:

```ruby
class Promptstack < Formula
  desc "AI-powered prompt management CLI"
  homepage "https://github.com/yourorg/promptstack"
  url "https://github.com/yourorg/promptstack/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "sha256_checksum_here"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/promptstack"
  end

  test do
    system bin/"promptstack", "--version"
  end
end
```

#### Install from Tap

```bash
# Add tap
brew tap yourorg/promptstack

# Install
brew install promptstack

# Update
brew upgrade promptstack

# Uninstall
brew uninstall promptstack
```

#### Install from Local Formula

```bash
# Install from local formula file
brew install Formula/promptstack.rb
```

### Linux Packages

#### Debian/Ubuntu (.deb)

Create `debian/control`:

```
Package: promptstack
Version: 1.0.0
Section: utils
Priority: optional
Architecture: amd64
Maintainer: Your Name <you@example.com>
Description: AI-powered prompt management CLI
 PromptStack is a terminal-based tool for managing AI prompts
 with intelligent suggestions and library integration.
Depends: libc6
```

Build package:

```bash
# Install packaging tools
sudo apt-get install dpkg-dev

# Create package directory structure
mkdir -p promptstack_1.0.0_amd64/DEBIAN
mkdir -p promptstack_1.0.0_amd64/usr/local/bin

# Copy binary
cp promptstack-linux-amd64 promptstack_1.0.0_amd64/usr/local/bin/promptstack
chmod +x promptstack_1.0.0_amd64/usr/local/bin/promptstack

# Create control file
cat > promptstack_1.0.0_amd64/DEBIAN/control << EOF
Package: promptstack
Version: 1.0.0
Section: utils
Priority: optional
Architecture: amd64
Maintainer: Your Name <you@example.com>
Description: AI-powered prompt management CLI
 PromptStack is a terminal-based tool for managing AI prompts
 with intelligent suggestions and library integration.
EOF

# Build package
dpkg-deb --build promptstack_1.0.0_amd64

# Install package
sudo dpkg -i promptstack_1.0.0_amd64.deb
```

#### Red Hat/CentOS/Fedora (.rpm)

Create `.spec` file:

```spec
Name:           promptstack
Version:        1.0.0
Release:        1%{?dist}
Summary:        AI-powered prompt management CLI
License:        MIT
URL:            https://github.com/yourorg/promptstack
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang

%description
PromptStack is a terminal-based tool for managing AI prompts
with intelligent suggestions and library integration.

%prep
%setup -q

%build
go build -ldflags="-s -w" -o promptstack ./cmd/promptstack

%install
install -D -m 0755 promptstack %{buildroot}%{_bindir}/promptstack

%files
%{_bindir}/promptstack

%changelog
* Mon Jan 07 2026 Your Name <you@example.com> - 1.0.0-1
- Initial release
```

Build package:

```bash
# Install packaging tools
sudo dnf install rpm-build

# Build package
rpmbuild -ba promptstack.spec

# Install package
sudo dnf install promptstack-1.0.0-1.x86_64.rpm
```

#### Arch Linux (AUR)

Create `PKGBUILD`:

```bash
pkgname=promptstack
pkgver=1.0.0
pkgrel=1
pkgdesc="AI-powered prompt management CLI"
arch=('x86_64')
url="https://github.com/yourorg/promptstack"
license=('MIT')
depends=()
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::https://github.com/yourorg/promptstack/archive/v$pkgver.tar.gz")
sha256sums=('sha256_checksum_here')

build() {
  cd "$pkgname-$pkgver"
  go build -ldflags="-s -w" -o promptstack ./cmd/promptstack
}

package() {
  cd "$pkgname-$pkgver"
  install -Dm755 promptstack "$pkgdir/usr/bin/promptstack"
}
```

Build package:

```bash
# Build package
makepkg -si

# Or build without installing
makepkg
```

### Windows Installer

#### Using NSIS

Create `installer.nsi`:

```nsis
!define APPNAME "PromptStack"
!define COMPANYNAME "Your Company"
!define DESCRIPTION "AI-powered prompt management CLI"
!define VERSIONMAJOR 1
!define VERSIONMINOR 0
!define VERSIONBUILD 0

RequestExecutionLevel admin

InstallDir "$PROGRAMFILES\${APPNAME}"

Page directory
Page instfiles

Section "install"
    SetOutPath $INSTDIR
    File "promptstack-windows-amd64.exe"
    File /r "starter-prompts"
    
    CreateDirectory "$SMPROGRAMS\${APPNAME}"
    CreateShortcut "$SMPROGRAMS\${APPNAME}\${APPNAME}.lnk" "$INSTDIR\promptstack-windows-amd64.exe"
    
    WriteUninstaller "$INSTDIR\uninstall.exe"
    CreateShortcut "$SMPROGRAMS\${APPNAME}\Uninstall.lnk" "$INSTDIR\uninstall.exe"
SectionEnd

Section "uninstall"
    Delete $INSTDIR\promptstack-windows-amd64.exe
    Delete $INSTDIR\uninstall.exe
    RMDir /r "$INSTDIR"
    RMDir /r "$SMPROGRAMS\${APPNAME}"
SectionEnd
```

Build installer:

```bash
# Install NSIS
# Download from https://nsis.sourceforge.io/

# Build installer
makensis installer.nsi
```

#### Using Inno Setup

Create `installer.iss`:

```iss
[Setup]
AppName=PromptStack
AppVersion=1.0.0
DefaultDirName={commonpf}\PromptStack
DefaultGroupName=PromptStack
OutputBaseFilename=promptstack-setup
Compression=lzma
SolidCompression=yes

[Files]
Source: "promptstack-windows-amd64.exe"; DestDir: "{app}"
Source: "starter-prompts\*"; DestDir: "{app}\starter-prompts"; Flags: ignoreversion recursesubdirs createallsubdirs

[Icons]
Name: "{group}\PromptStack"; Filename: "{app}\promptstack-windows-amd64.exe"
Name: "{group}\Uninstall PromptStack"; Filename: "{uninstallexe}"

[UninstallDelete]
Type: filesandordirs; Name: "{app}"
```

Build installer:

```bash
# Install Inno Setup
# Download from https://jrsoftware.org/isdl.php

# Build installer
iscc installer.iss
```

### Manual Installation

#### macOS/Linux

```bash
# Download binary
wget https://github.com/yourorg/promptstack/releases/download/v1.0.0/promptstack-darwin-amd64

# Make executable
chmod +x promptstack-darwin-amd64

# Move to PATH
sudo mv promptstack-darwin-amd64 /usr/local/bin/promptstack

# Verify installation
promptstack --version
```

#### Windows

```powershell
# Download binary
Invoke-WebRequest -Uri "https://github.com/yourorg/promptstack/releases/download/v1.0.0/promptstack-windows-amd64.exe" -OutFile "promptstack.exe"

# Add to PATH (PowerShell)
$env:Path += ";C:\path\to\promptstack"

# Or add permanently
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\path\to\promptstack", "Machine")

# Verify installation
promptstack.exe --version
```

---

## Release Process

### Pre-Release Checklist

- [ ] All tests passing (`go test ./...`)
- [ ] Code coverage >80%
- [ ] No linting errors (`golangci-lint run`)
- [ ] No security vulnerabilities (`govulncheck ./...`)
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version number updated
- [ ] Release notes prepared

### Create Release Tag

```bash
# Update version in code
# Update CHANGELOG.md

# Commit changes
git add .
git commit -m "Release v1.0.0"

# Create tag
git tag -a v1.0.0 -m "Release v1.0.0"

# Push tag
git push origin v1.0.0
```

### Build All Platforms

```bash
# Create build script
cat > build.sh << 'EOF'
#!/bin/bash
set -e

VERSION=${1:-1.0.0}
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(git rev-parse --short HEAD)

LDFLAGS="-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT"

echo "Building PromptStack $VERSION..."

# macOS amd64
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o promptstack-darwin-amd64 ./cmd/promptstack

# macOS arm64
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o promptstack-darwin-arm64 ./cmd/promptstack

# macOS universal
lipo -create -output promptstack-darwin-universal promptstack-darwin-amd64 promptstack-darwin-arm64

# Linux amd64
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o promptstack-linux-amd64 ./cmd/promptstack

# Linux arm64
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o promptstack-linux-arm64 ./cmd/promptstack

# Windows amd64
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o promptstack-windows-amd64.exe ./cmd/promptstack

echo "Build complete!"
EOF

chmod +x build.sh
./build.sh 1.0.0
```

### Create Release Assets

```bash
# Create checksums
shasum -a 256 promptstack-* > checksums.txt

# Create archives
tar -czf promptstack-darwin-amd64.tar.gz promptstack-darwin-amd64
tar -czf promptstack-darwin-arm64.tar.gz promptstack-darwin-arm64
tar -czf promptstack-darwin-universal.tar.gz promptstack-darwin-universal
tar -czf promptstack-linux-amd64.tar.gz promptstack-linux-amd64
tar -czf promptstack-linux-arm64.tar.gz promptstack-linux-arm64
zip promptstack-windows-amd64.zip promptstack-windows-amd64.exe
```

### GitHub Release

1. Go to GitHub repository
2. Click "Releases" → "Create a new release"
3. Select tag `v1.0.0`
4. Add release title and description
5. Upload binaries and checksums
6. Click "Publish release"

### Update Homebrew Formula

```bash
# Update version and SHA256 in Formula/promptstack.rb
# Commit and push
git add Formula/promptstack.rb
git commit -m "Update Homebrew formula to v1.0.0"
git push origin main
```

---

## Update Mechanism

### Version Checking

Implement version checking in the application:

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "runtime"
)

type ReleaseInfo struct {
    TagName string `json:"tag_name"`
    HtmlUrl string `json:"html_url"`
}

func CheckForUpdates() error {
    resp, err := http.Get("https://api.github.com/repos/yourorg/promptstack/releases/latest")
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    var release ReleaseInfo
    if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
        return err
    }

    if release.TagName > Version {
        fmt.Printf("New version available: %s (current: %s)\n", release.TagName, Version)
        fmt.Printf("Download: %s\n", release.HtmlUrl)
    } else {
        fmt.Println("You're using the latest version!")
    }

    return nil
}
```

### Auto-Update (Optional)

Implement auto-update functionality:

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
)

func SelfUpdate() error {
    // Download new binary
    url := fmt.Sprintf("https://github.com/yourorg/promptstack/releases/download/v%s/promptstack-%s-%s", 
        LatestVersion, runtime.GOOS, runtime.GOARCH)
    
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Create temp file
    tmpFile := filepath.Join(os.TempDir(), "promptstack-new")
    out, err := os.Create(tmpFile)
    if err != nil {
        return err
    }
    defer out.Close()

    // Download binary
    if _, err := io.Copy(out, resp.Body); err != nil {
        return err
    }

    // Make executable
    if err := os.Chmod(tmpFile, 0755); err != nil {
        return err
    }

    // Replace current binary
    exePath, err := os.Executable()
    if err != nil {
        return err
    }

    // Move new binary to current location
    if err := os.Rename(tmpFile, exePath); err != nil {
        return err
    }

    fmt.Println("Update successful! Restarting...")
    
    // Restart application
    return exec.Command(exePath).Start()
}
```

### Migration Procedures

When upgrading between versions:

1. **Backup Configuration**:
   ```bash
   cp ~/.promptstack/config.yaml ~/.promptstack/config.yaml.backup
   ```

2. **Backup Database**:
   ```bash
   cp ~/.promptstack/data/history.db ~/.promptstack/data/history.db.backup
   ```

3. **Install New Version**:
   ```bash
   # Using Homebrew
   brew upgrade promptstack
   
   # Or manual
   wget https://github.com/yourorg/promptstack/releases/download/v1.0.0/promptstack-darwin-amd64
   chmod +x promptstack-darwin-amd64
   sudo mv promptstack-darwin-amd64 /usr/local/bin/promptstack
   ```

4. **Run Migration** (if needed):
   ```bash
   promptstack migrate
   ```

5. **Verify Installation**:
   ```bash
   promptstack --version
   promptstack doctor
   ```

---

## CI/CD Integration

### GitHub Actions Workflow

Create `.github/workflows/build.yml`:

```yaml
name: Build and Release

on:
  push:
    tags:
      - 'v*'
  workflow_dispatch:

jobs:
  build:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        include:
          - os: macos-latest
            goos: darwin
            goarch: amd64
          - os: macos-latest
            goos: darwin
            goarch: arm64
          - os: ubuntu-latest
            goos: linux
            goarch: amd64
          - os: ubuntu-latest
            goos: linux
            goarch: arm64
          - os: windows-latest
            goos: windows
            goarch: amd64

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24.0'

      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out

      - name: Build binary
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
          CGO_ENABLED: 0
          VERSION: ${{ github.ref_name }}
        run: |
          go build -ldflags="-s -w -X main.Version=$VERSION" -o promptstack-${{ matrix.goos }}-${{ matrix.goarch }}${{ matrix.goos == 'windows' && '.exe' || '' }} ./cmd/promptstack

      - name: Upload artifact
        uses: actions/upload-artifact@v3
        with:
          name: promptstack-${{ matrix.goos }}-${{ matrix.goarch }}
          path: promptstack-${{ matrix.goos }}-${{ matrix.goarch }}${{ matrix.goos == 'windows' && '.exe' || '' }}

  release:
    needs: build
    runs-on: ubuntu-latest
    if: startsWith(github.ref, 'refs/tags/')
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Download all artifacts
        uses: actions/download-artifact@v3

      - name: Create checksums
        run: |
          shasum -a 256 promptstack-* > checksums.txt

      - name: Create release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            promptstack-*
            checksums.txt
          draft: false
          prerelease: false
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Makefile

Create `Makefile` for build automation:

```makefile
.PHONY: build test clean install release

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)

build:
	go build -ldflags="$(LDFLAGS)" -o promptstack ./cmd/promptstack

build-all:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o promptstack-darwin-amd64 ./cmd/promptstack
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o promptstack-darwin-arm64 ./cmd/promptstack
	lipo -create -output promptstack-darwin-universal promptstack-darwin-amd64 promptstack-darwin-arm64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o promptstack-linux-amd64 ./cmd/promptstack
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o promptstack-linux-arm64 ./cmd/promptstack
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o promptstack-windows-amd64.exe ./cmd/promptstack

test:
	go test -v -race -coverprofile=coverage.out ./...

coverage:
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

clean:
	go clean
	rm -f promptstack-*
	rm -f coverage.out coverage.html

install: build
	install -m 755 promptstack /usr/local/bin/promptstack

release: clean test lint build-all
	shasum -a 256 promptstack-* > checksums.txt
	tar -czf promptstack-darwin-amd64.tar.gz promptstack-darwin-amd64
	tar -czf promptstack-darwin-arm64.tar.gz promptstack-darwin-arm64
	tar -czf promptstack-darwin-universal.tar.gz promptstack-darwin-universal
	tar -czf promptstack-linux-amd64.tar.gz promptstack-linux-amd64
	tar -czf promptstack-linux-arm64.tar.gz promptstack-linux-arm64
	zip promptstack-windows-amd64.zip promptstack-windows-amd64.exe
```

---

## Troubleshooting

### Build Issues

**Problem**: `go: cannot find main module`

**Solution**:
```bash
# Initialize module if needed
go mod init github.com/yourorg/promptstack

# Download dependencies
go mod download
```

**Problem**: `CGO_ENABLED=0` build fails

**Solution**:
```bash
# Check if package requires CGO
go list -f '{{.CgoFiles}}' ./...

# If CGO is required, install C compiler
# macOS: xcode-select --install
# Linux: sudo apt-get install build-essential
# Windows: Install MinGW-w64
```

### Cross-Compilation Issues

**Problem**: Binary doesn't run on target platform

**Solution**:
```bash
# Verify binary architecture
file promptstack-linux-amd64

# Rebuild with correct GOOS/GOARCH
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o promptstack-linux-amd64 ./cmd/promptstack
```

**Problem**: macOS binary notarization required

**Solution**:
```bash
# Sign binary
codesign --force --deep --sign "Developer ID Application: Your Name" promptstack-darwin-amd64

# Notarize binary
xcrun notarytool submit promptstack-darwin-amd64 --apple-id "your@email.com" --password "app-specific-password" --team-id "TEAMID" --wait
```

### Installation Issues

**Problem**: Permission denied when installing

**Solution**:
```bash
# Use sudo for system-wide installation
sudo install -m 755 promptstack /usr/local/bin/promptstack

# Or install to user directory
mkdir -p ~/.local/bin
install -m 755 promptstack ~/.local/bin/promptstack
```

**Problem**: Binary not found in PATH

**Solution**:
```bash
# Add to PATH (bash)
echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
source ~/.bashrc

# Add to PATH (zsh)
echo 'export PATH=$PATH:/usr/local/bin' >> ~/.zshrc
source ~/.zshrc
```

### Runtime Issues

**Problem**: `command not found: promptstack`

**Solution**:
```bash
# Verify installation
which promptstack

# Check PATH
echo $PATH

# Reinstall if needed
sudo install -m 755 promptstack /usr/local/bin/promptstack
```

**Problem**: Database migration fails after update

**Solution**:
```bash
# Backup current database
cp ~/.promptstack/data/history.db ~/.promptstack/data/history.db.backup

# Run migration
promptstack migrate

# If migration fails, restore backup
cp ~/.promptstack/data/history.db.backup ~/.promptstack/data/history.db
```

---

## Additional Resources

- [Go Build Documentation](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies)
- [Cross-Compilation Guide](https://go.dev/doc/install/source#environment)
- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Semantic Versioning](https://semver.org/)

---

**Document Maintenance**: Update this document when build processes change, new platforms are added, or deployment procedures are modified.

**Questions or Issues**: Contact the maintainers or open an issue in the repository.