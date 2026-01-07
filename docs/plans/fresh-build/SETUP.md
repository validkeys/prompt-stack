# Development Environment Setup

This guide will help you set up a complete development environment for building PromptStack from scratch.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation Steps](#installation-steps)
- [IDE Configuration](#ide-configuration)
- [Environment Variables](#environment-variables)
- [Verification Steps](#verification-steps)
- [Troubleshooting](#troubleshooting)
- [Platform-Specific Notes](#platform-specific-notes)

---

## Prerequisites

### Required Software

- **Go**: Version 1.21 or higher
  - Download from [https://go.dev/dl/](https://go.dev/dl/)
  - Verify installation: `go version`
  
- **Git**: For version control
  - macOS: Included with Xcode Command Line Tools
  - Linux: `sudo apt-get install git` or `sudo yum install git`
  - Windows: Download from [https://git-scm.com/](https://git-scm.com/)
  - Verify installation: `git --version`

- **Terminal Emulator**: With 256-color support
  - macOS: Terminal.app, iTerm2, or Warp
  - Linux: Most modern terminals support 256 colors
  - Windows: Windows Terminal (recommended) or Git Bash

### Optional but Recommended

- **Make**: For build automation (included with most Unix-like systems)
- **jq**: For JSON parsing (useful for debugging): `brew install jq` (macOS)

---

## Installation Steps

### Step 1: Install Go

#### macOS

```bash
# Using Homebrew (recommended)
brew install go

# Or download from https://go.dev/dl/ and follow the installer
```

#### Linux

```bash
# Download Go 1.21+ from https://go.dev/dl/
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz

# Extract to /usr/local
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin
```

#### Windows

1. Download the MSI installer from [https://go.dev/dl/](https://go.dev/dl/)
2. Run the installer and follow the prompts
3. Restart your terminal/command prompt

**Verify Go installation:**

```bash
go version
# Should output: go version go1.21.x darwin/amd64 (or similar)
```

### Step 2: Clone the Repository

```bash
# Clone the repository
git clone https://github.com/yourorg/promptstack.git
cd promptstack

# Verify you're in the correct directory
pwd
# Should show: /path/to/promptstack
```

### Step 3: Install Dependencies

```bash
# Download all Go module dependencies
go mod download

# Verify dependencies are downloaded
go mod verify

# Tidy up dependencies (removes unused ones)
go mod tidy
```

### Step 4: Verify Go Environment

```bash
# Check Go environment
go env

# Key values to verify:
# - GOPATH: Should be set (usually ~/go)
# - GOROOT: Should point to Go installation
# - GO111MODULE: Should be "on" or "auto"
```

---

## IDE Configuration

### Visual Studio Code (Recommended)

#### Required Extensions

1. **Go** (by Go Team at Google)
   - Install from VS Code Extensions marketplace
   - Provides IntelliSense, debugging, and testing support

2. **Makefile Support** (optional but recommended)
   - For syntax highlighting in Makefile

#### VS Code Settings

Create or update `.vscode/settings.json`:

```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "go.formatTool": "goimports",
  "go.testFlags": ["-v", "-race"],
  "go.coverOnSingleTest": true,
  "go.coverOnSingleTestFile": true,
  "go.coverageDecorator": {
    "type": "gutter"
  },
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  },
  "files.exclude": {
    "**/.git": true,
    "**/.DS_Store": true,
    "**/node_modules": true
  }
}
```

#### VS Code Launch Configuration

Create `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch PromptStack",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/promptstack",
      "env": {
        "PROMPTSTACK_LOG_LEVEL": "debug"
      },
      "args": []
    },
    {
      "name": "Debug Tests",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${workspaceFolder}",
      "args": [
        "-test.v",
        "-test.run",
        "TestFunctionName"
      ]
    }
  ]
}
```

### GoLand (JetBrains)

1. Open the project directory in GoLand
2. Go to **File → Settings → Go → GOROOT** and verify it's set correctly
3. Go to **File → Settings → Go → GOPATH** and verify it's set correctly
4. Enable **Go Modules** integration
5. Configure **Code Style** to match the project's [`go-style-guide.md`](go-style-guide.md)

### Vim/Neovim

#### Required Plugins

- **vim-go**: Go support for Vim/Neovim
- **coc.nvim** or **nvim-lspconfig**: LSP support

#### Configuration Example (Neovim with coc.nvim)

```vim
" Install vim-go
Plug 'fatih/vim-go'

" Install coc.nvim
Plug 'neoclide/coc.nvim'

" Configure Go
let g:go_def_mode='gopls'
let g:go_info_mode='gopls'

" Auto-format on save
autocmd BufWritePre *.go :GoImports
```

---

## Environment Variables

### Required Environment Variables

None are strictly required for development, but the following are useful:

```bash
# Set Go module proxy (optional, for faster downloads)
export GOPROXY=https://proxy.golang.org,direct

# Enable Go modules (default in Go 1.16+)
export GO111MODULE=on
```

### Optional Environment Variables for Development

```bash
# Enable debug logging for PromptStack
export PROMPTSTACK_LOG_LEVEL=debug

# Enable verbose debug output
export PROMPTSTACK_DEBUG=1

# Set custom data directory for testing
export PROMPTSTACK_DATA_DIR=/tmp/promptstack-test
```

### Setting Environment Variables

#### macOS/Linux (bash)

Add to `~/.bashrc` or `~/.bash_profile`:

```bash
export GOPROXY=https://proxy.golang.org,direct
export GO111MODULE=on
export PROMPTSTACK_LOG_LEVEL=debug
```

#### macOS/Linux (zsh)

Add to `~/.zshrc`:

```bash
export GOPROXY=https://proxy.golang.org,direct
export GO111MODULE=on
export PROMPTSTACK_LOG_LEVEL=debug
```

#### Windows

Add to System Environment Variables or create a `.env` file:

```
GOPROXY=https://proxy.golang.org,direct
GO111MODULE=on
PROMPTSTACK_LOG_LEVEL=debug
```

---

## Verification Steps

### Step 1: Verify Go Installation

```bash
go version
# Expected output: go version go1.21.x <os>/<arch>
```

### Step 2: Verify Project Structure

```bash
# List project structure
ls -la

# Expected directories:
# - cmd/
# - internal/
# - ui/
# - docs/
# - test/
# - go.mod
# - go.sum
```

### Step 3: Verify Dependencies

```bash
# Check go.mod exists
cat go.mod

# Verify dependencies are downloaded
go list -m all

# Should show all dependencies including:
# - github.com/charmbracelet/bubbletea
# - github.com/charmbracelet/lipgloss
# - modernc.org/sqlite
# - etc.
```

### Step 4: Build the Project

```bash
# Build the application
go build -o promptstack ./cmd/promptstack

# Verify binary was created
ls -lh promptstack

# Expected: promptstack binary file exists
```

### Step 5: Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...

# Expected: All tests pass
```

### Step 6: Run the Application

```bash
# Run the application (will trigger first-time setup)
./promptstack

# Expected: Application launches and prompts for initial setup
```

### Step 7: Verify Development Tools

```bash
# Check if gofmt is available
which gofmt

# Check if go vet is available
which go vet

# Check if go test is available
which go test

# All should be available in your PATH
```

---

## Troubleshooting

### Common Issues

#### Issue: "go: command not found"

**Solution:**
- Verify Go is installed: `which go`
- Add Go to PATH (see Installation Steps above)
- Restart your terminal

#### Issue: "go.mod file not found"

**Solution:**
- Ensure you're in the project root directory
- Run `ls -la` to verify go.mod exists
- If missing, run `go mod init github.com/yourorg/promptstack`

#### Issue: "module not found" or "cannot find package"

**Solution:**
```bash
# Download dependencies
go mod download

# Tidy up dependencies
go mod tidy

# Verify go.mod is correct
cat go.mod
```

#### Issue: Build fails with "undefined: ..."

**Solution:**
- Ensure all dependencies are downloaded: `go mod download`
- Check for typos in import paths
- Verify you're using Go 1.21 or higher: `go version`

#### Issue: Tests fail with "no test files"

**Solution:**
- Ensure you're in the correct directory
- Run tests from project root: `go test ./...`
- Check that test files exist: `find . -name "*_test.go"`

#### Issue: "permission denied" when running binary

**Solution:**
```bash
# Make binary executable
chmod +x promptstack

# Or build with correct permissions
go build -o promptstack ./cmd/promptstack
```

#### Issue: Terminal colors not displaying correctly

**Solution:**
- Verify your terminal supports 256 colors
- Try a different terminal emulator (iTerm2, Warp, Windows Terminal)
- Check TERM environment variable: `echo $TERM`
- Set TERM to xterm-256color: `export TERM=xterm-256color`

#### Issue: "cannot load package: package github.com/...: no Go files"

**Solution:**
- Verify you're in the correct directory
- Check that the package directory exists
- Ensure the directory contains .go files

#### Issue: "go: cannot find main module"

**Solution:**
- Ensure you're in the project root directory (where go.mod is)
- Run `pwd` to verify current directory
- Navigate to project root: `cd /path/to/promptstack`

---

## Platform-Specific Notes

### macOS

#### Xcode Command Line Tools

```bash
# Install Xcode Command Line Tools (required for some Go packages)
xcode-select --install

# Verify installation
xcode-select -p
```

#### Homebrew

```bash
# Install Homebrew if not already installed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Go via Homebrew
brew install go

# Install other useful tools
brew install make jq
```

#### File System

- macOS is case-insensitive by default
- Be careful with file naming (avoid case-only differences)
- Use `.DS_Store` in `.gitignore` to ignore macOS metadata files

### Linux

#### Package Managers

**Ubuntu/Debian:**
```bash
sudo apt-get update
sudo apt-get install golang git make jq
```

**Fedora/RHEL:**
```bash
sudo dnf install golang git make jq
```

**Arch Linux:**
```bash
sudo pacman -S go git make jq
```

#### Permissions

- You may need to use `sudo` for system-wide installations
- Consider using a user-local Go installation to avoid permission issues

### Windows

#### Git Bash

- Use Git Bash for a Unix-like environment
- Ensure Git Bash is in your PATH
- Configure Git Bash to use Windows-style paths if needed

#### Windows Terminal

- Install Windows Terminal for a better terminal experience
- Configure profiles for PowerShell, Git Bash, and WSL

#### WSL (Windows Subsystem for Linux)

- Consider using WSL2 for a full Linux development environment
- Install Go inside WSL: `sudo apt-get install golang`
- Access Windows files from WSL: `/mnt/c/`

#### Path Issues

- Ensure Go binary directory is in your PATH
- Use forward slashes in paths (even on Windows)
- Avoid spaces in directory names

---

## Next Steps

After completing the setup:

1. **Read the documentation:**
   - [`README.md`](README.md) - Project overview
   - [`HOW-TO-USE.md`](HOW-TO-USE.md) - Development workflow
   - [`go-style-guide.md`](go-style-guide.md) - Coding standards
   - [`go-testing-guide.md`](go-testing-guide.md) - Testing patterns

2. **Start development:**
   - Review [`milestones.md`](milestones.md) for the development roadmap
   - Begin with Milestone 1: Bootstrap & Config
   - Follow the TDD process strictly

3. **Configure your IDE:**
   - Set up linting and formatting
   - Configure test runners
   - Set up debugging

4. **Join the community:**
   - Read [`CONTRIBUTING.md`](../../CONTRIBUTING.md) (when available)
   - Report issues and suggest features

---

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Modules Reference](https://golang.org/ref/mod)
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss)

---

## Getting Help

If you encounter issues not covered in this guide:

1. Check the [troubleshooting section](#troubleshooting) above
2. Search existing issues in the project repository
3. Ask questions in the project's discussion forum
4. Create a new issue with detailed information about your problem

---

**Last Updated**: 2026-01-07  
**Status**: Ready for use  
**Maintainer**: PromptStack Development Team