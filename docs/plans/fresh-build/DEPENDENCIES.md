# Dependencies

**Last Updated**: 2026-01-07  
**Go Version**: 1.24.0

This document provides a complete overview of all dependencies used in PromptStack, including their purposes, versions, licensing information, and management procedures.

---

## Table of Contents

- [Core Dependencies](#core-dependencies)
- [Development Dependencies](#development-dependencies)
- [Dependency Categories](#dependency-categories)
- [License Information](#license-information)
- [Adding Dependencies](#adding-dependencies)
- [Updating Dependencies](#updating-dependencies)
- [Security Considerations](#security-considerations)
- [Dependency Audit](#dependency-audit)

---

## Core Dependencies

These dependencies are required for the application to run.

### TUI Framework

| Package | Version | Purpose | License |
|---------|---------|---------|---------|
| `github.com/charmbracelet/bubbletea` | v1.3.10 | TUI framework for terminal interfaces | MIT |
| `github.com/charmbracelet/lipgloss` | v1.1.1 | Styling library for Bubble Tea | MIT |
| `github.com/charmbracelet/bubbles` | v0.21.0 | Pre-built Bubble Tea components (viewport, textinput, spinner) | MIT |
| `github.com/charmbracelet/glamour` | v0.10.0 | Markdown rendering for terminal | MIT |
| `github.com/charmbracelet/x/ansi` | v0.10.1 | ANSI escape sequence handling | MIT |
| `github.com/charmbracelet/x/cellbuf` | v0.0.13 | Cell buffer utilities | MIT |
| `github.com/charmbracelet/x/term` | v0.2.1 | Terminal utilities | MIT |
| `github.com/charmbracelet/x/exp/slice` | v0.0.0-20250327172914-2fdc97757edf | Slice utilities | MIT |
| `github.com/charmbracelet/colorprofile` | v0.2.3-0.20250311203215-f60798e515dc | Color profile management | MIT |

**Rationale**: Charmbracelet provides a complete, well-maintained ecosystem for building modern terminal UIs. All components work together seamlessly and follow consistent design patterns.

### Database & Storage

| Package | Version | Purpose | License |
|---------|---------|---------|---------|
| `modernc.org/sqlite` | v1.42.2 | Pure Go SQLite implementation (no CGO) | BSD-3-Clause |
| `gopkg.in/yaml.v3` | v3.0.1 | YAML parsing for frontmatter | MIT |

**Rationale**: 
- `modernc.org/sqlite` chosen over `github.com/mattn/go-sqlite3` to avoid CGO dependency, making cross-compilation simpler and binary distribution easier
- `yaml.v3` is the de facto standard for YAML parsing in Go

### AI Integration

| Package | Version | Purpose | License |
|---------|---------|---------|---------|
| `github.com/anthropics/anthropic-sdk-go` | v1.19.0 | Official Claude API SDK | MIT |

**Rationale**: Official SDK provides best-in-class support for Claude API features, including streaming, error handling, and context management.

### Fuzzy Search & Text Processing

| Package | Version | Purpose | License |
|---------|---------|---------|---------|
| `github.com/sahilm/fuzzy` | v0.1.1 | Lightweight fuzzy matching algorithm | MIT |
| `github.com/tidwall/gjson` | v1.18.0 | Fast JSON parsing and querying | MIT |
| `github.com/tidwall/sjson` | v1.2.5 | Fast JSON modification | MIT |
| `github.com/tidwall/pretty` | v1.2.1 | JSON pretty-printing | MIT |
| `github.com/tidwall/match` | v1.1.1 | String matching utilities | MIT |

**Rationale**: 
- `sahilm/fuzzy` is the standard for fuzzy matching in Go, used by many popular CLI tools
- `tidwall` libraries provide excellent performance for JSON operations

### File System & Git Operations

| Package | Version | Purpose | License |
|---------|---------|---------|---------|
| `github.com/mitchellh/go-homedir` | v1.1.0 | Cross-platform home directory detection | MIT |
| `github.com/erikgeiser/coninput` | v0.0.0-20211004153227-1c3628e74d0f | Console input handling for Windows | MIT |

**Rationale**: 
- `go-homedir` provides reliable cross-platform home directory resolution
- `coninput` enables proper terminal input handling on Windows (future-proofing)

### Text Processing & Syntax Highlighting

| Package | Version | Purpose | License |
|---------|---------|---------|---------|
| `github.com/alecthomas/chroma/v2` | v2.14.0 | Syntax highlighting for code blocks | MIT |
| `github.com/yuin/goldmark` | v1.7.8 | CommonMark-compliant markdown parser | MIT |
| `github.com/yuin/goldmark-emoji` | v1.0.5 | Emoji support for goldmark | MIT |
| `github.com/lucasb-eyer/go-colorful` | v1.2.0 | Color manipulation utilities | BSD-3-Clause |
| `github.com/gorilla/css` | v1.0.1 | CSS parsing for markdown rendering | BSD-2-Clause |
| `github.com/microcosm-cc/bluemonday` | v1.0.27 | HTML sanitization | BSD-2-Clause |

**Rationale**: 
- `chroma` provides excellent syntax highlighting with many language support
- `goldmark` is fast, extensible, and CommonMark compliant
- `bluemonday` ensures safe HTML rendering from markdown

### Terminal & UI Utilities

| Package | Version | Purpose | License |
|---------|---------|---------|---------|
| `github.com/mattn/go-isatty` | v0.0.20 | Terminal detection | MIT |
| `github.com/mattn/go-runewidth` | v0.0.16 | Unicode-aware string width calculation | MIT |
| `github.com/mattn/go-localereader` | v0.0.1 | Locale-aware reader | MIT |
| `github.com/muesli/termenv` | v0.16.0 | Terminal environment detection | MIT |
| `github.com/muesli/ansi` | v0.0.0-20230316100256-276c6243b2f6 | ANSI escape sequence handling | MIT |
| `github.com/muesli/reflow` | v0.3.0 | Text reflow utilities | MIT |
| `github.com/muesli/cancelreader` | v0.2.2 | Cancelable reader for async operations | MIT |
| `github.com/xo/terminfo` | v0.0.0-20220910002029-abceb7e1c41e | Terminal capability database | MIT |
| `github.com/dlclark/regexp2` | v1.11.0 | Regex engine with Unicode support | BSD-3-Clause |
| `github.com/aymanbagabas/go-osc52/v2` | v2.0.1 | OSC 52 clipboard support | MIT |
| `github.com/aymerick/douceur` | v0.2.0 | CSS parser and serializer | BSD-2-Clause |
| `github.com/rivo/uniseg` | v0.4.7 | Unicode text segmentation | MIT |
| `github.com/dustin/go-humanize` | v1.0.1 | Human-readable formatting | MIT |
| `github.com/ncruces/go-strftime` | v0.1.9 | Time formatting utilities | MIT |

**Rationale**: These utilities provide cross-platform terminal compatibility, proper Unicode handling, and text formatting capabilities essential for a polished TUI experience.

### Logging

| Package | Version | Purpose | License |
|---------|---------|---------|---------|
| `go.uber.org/zap` | v1.27.1 | Structured logging | MIT |
| `gopkg.in/natefinch/lumberjack.v2` | v2.2.1 | Log rotation | MIT |
| `go.uber.org/multierr` | v1.10.0 | Multi-error handling | MIT |

**Rationale**: 
- `zap` is the industry standard for structured logging in Go, offering excellent performance
- `lumberjack` provides automatic log rotation

### UUID Generation

| Package | Version | Purpose | License |
|---------|---------|---------|---------|
| `github.com/google/uuid` | v1.6.0 | UUID generation | BSD-3-Clause |

**Rationale**: Google's UUID library is the de facto standard for UUID generation in Go.

### Clipboard Support

| Package | Version | Purpose | License |
|---------|---------|---------|---------|
| `github.com/atotto/clipboard` | v0.1.4 | Cross-platform clipboard access | BSD-3-Clause |

**Rationale**: Simple, cross-platform clipboard access for copy-to-clipboard functionality.

### Standard Library Extensions

| Package | Version | Purpose | License |
|---------|---------|---------|---------|
| `golang.org/x/sys` | v0.36.0 | System call wrappers | BSD-3-Clause |
| `golang.org/x/term` | v0.32.0 | Terminal utilities | BSD-3-Clause |
| `golang.org/x/text` | v0.27.0 | Text processing extensions | BSD-3-Clause |
| `golang.org/x/net` | v0.41.0 | Network extensions | BSD-3-Clause |
| `golang.org/x/exp` | v0.0.0-20250620022241-b7579e27df2b | Experimental features | BSD-3-Clause |

**Rationale**: These are official Go extensions that provide additional functionality beyond the standard library.

### Math & Performance

| Package | Version | Purpose | License |
|---------|---------|---------|---------|
| `modernc.org/libc` | v1.66.10 | C library implementation | BSD-3-Clause |
| `modernc.org/mathutil` | v1.7.1 | Math utilities | BSD-3-Clause |
| `modernc.org/memory` | v1.11.0 | Memory utilities | BSD-3-Clause |
| `github.com/remyoudompheng/bigfft` | v0.0.0-20230129092748-24d4a6f8daec | Fast FFT for big integers | BSD-3-Clause |

**Rationale**: These are indirect dependencies required by `modernc.org/sqlite` and other packages.

---

## Development Dependencies

These dependencies are only required for testing and development.

| Package | Version | Purpose | License |
|---------|---------|---------|---------|
| `github.com/stretchr/testify` | v1.8.4 | Testing framework and assertions | MIT |
| `github.com/golang/mock` | v1.6.0 | Mock generation for testing | Apache-2.0 |

**Rationale**: 
- `testify` provides comprehensive testing utilities including assertions, test suites, and mocking
- `golang/mock` is the official Go mocking framework for generating mock interfaces

---

## Dependency Categories

### By Layer

**UI Layer**:
- All `charmbracelet/*` packages
- `chroma`, `goldmark`, `bluemonday`
- Terminal utilities (`mattn/*`, `muesli/*`, `xo/terminfo`)

**Domain Layer**:
- `anthropic-sdk-go` (AI integration)
- `sahilm/fuzzy` (search)
- `yaml.v3` (parsing)

**Platform Layer**:
- `modernc.org/sqlite` (database)
- `zap`, `lumberjack` (logging)
- `go-homedir`, `coninput` (filesystem)

### By License

**MIT License** (permissive, no restrictions):
- All `charmbracelet/*` packages
- `anthropic-sdk-go`
- `sahilm/fuzzy`
- `tidwall/*` packages
- `yaml.v3`
- `zap`, `lumberjack`, `multierr`
- `google/uuid`
- `atotto/clipboard`
- `testify`
- `chroma`
- `goldmark`, `goldmark-emoji`
- `mattn/*` packages
- `muesli/*` packages
- `rivo/uniseg`
- `dustin/go-humanize`
- `ncruces/go-strftime`
- `aymanbagabas/go-osc52`
- `golang.org/x/*` packages

**BSD-3-Clause** (permissive, requires attribution):
- `modernc.org/*` packages
- `google/uuid`
- `lucasb-eyer/go-colorful`
- `dlclark/regexp2`

**BSD-2-Clause** (permissive, requires attribution):
- `gorilla/css`
- `microcosm-cc/bluemonday`
- `aymerick/douceur`

**Apache-2.0** (permissive, requires attribution):
- `golang/mock`

**License Compatibility**: All dependencies use permissive licenses (MIT, BSD, Apache) that are compatible with the project's intended open-source distribution.

---

## License Information

### Summary

All dependencies use permissive open-source licenses that allow:
- ✅ Commercial use
- ✅ Modification
- ✅ Distribution
- ✅ Private use

### License Requirements

**MIT License**:
- Include license copy
- Preserve copyright notice

**BSD Licenses**:
- Include license copy
- Preserve copyright notice
- Acknowledge use in documentation

**Apache-2.0 License**:
- Include license copy
- Preserve copyright notice
- State changes made
- Include NOTICE file if present

### Compliance Checklist

- [ ] All dependency licenses are documented in this file
- [ ] Third-party licenses will be included in distribution
- [ ] Attribution notices will be included in documentation
- [ ] No GPL or copyleft licenses that would restrict distribution

---

## Adding Dependencies

### Prerequisites

Before adding a new dependency, consider:

1. **Is it necessary?** Can the functionality be implemented with existing dependencies or standard library?
2. **License compatibility** - Must use a permissive license (MIT, BSD, Apache)
3. **Maintenance status** - Actively maintained with recent commits
4. **Community adoption** - Widely used and tested
5. **Dependency tree** - Minimal transitive dependencies
6. **Performance** - Efficient and well-optimized

### Procedure

1. **Research the dependency**:
   ```bash
   # Check package information
   go list -m -json github.com/package/name
   ```

2. **Add the dependency**:
   ```bash
   # Add specific version
   go get github.com/package/name@v1.2.3
   
   # Add to go.mod
   go mod tidy
   ```

3. **Update documentation**:
   - Add entry to appropriate section in this file
   - Include version, purpose, license, and rationale
   - Update dependency categories if needed

4. **Test thoroughly**:
   ```bash
   # Run all tests
   go test ./...
   
   # Run with race detector
   go test -race ./...
   ```

5. **Commit with clear message**:
   ```bash
   git add go.mod go.sum DEPENDENCIES.md
   git commit -m "Add github.com/package/name@v1.2.3 for [purpose]"
   ```

### Version Pinning Strategy

- **Pin exact versions** for all dependencies (no `@latest`)
- **Use semantic versioning** - prefer stable releases over pre-releases
- **Document rationale** for each dependency choice
- **Review updates** before applying (see [Updating Dependencies](#updating-dependencies))

---

## Updating Dependencies

### Update Policy

- **Security updates**: Apply immediately
- **Minor updates**: Review and test before applying
- **Major updates**: Evaluate breaking changes, plan migration
- **Frequency**: Review monthly, update quarterly

### Procedure

1. **Check for updates**:
   ```bash
   # List available updates
   go list -u -m all
   
   # Check specific package
   go list -u -m github.com/package/name
   ```

2. **Review changes**:
   - Read release notes and changelog
   - Check for breaking changes
   - Review commit history
   - Check for security advisories

3. **Update dependency**:
   ```bash
   # Update specific package
   go get github.com/package/name@v1.2.4
   
   # Update all dependencies
   go get -u ./...
   
   # Tidy go.mod
   go mod tidy
   ```

4. **Test thoroughly**:
   ```bash
   # Run all tests
   go test ./...
   
   # Run with race detector
   go test -race ./...
   
   # Build for all platforms
   make build-all
   ```

5. **Update documentation**:
   - Update version in this file
   - Document any breaking changes
   - Update rationale if needed

6. **Commit**:
   ```bash
   git add go.mod go.sum DEPENDENCIES.md
   git commit -m "Update github.com/package/name from v1.2.3 to v1.2.4"
   ```

### Rollback Procedure

If an update causes issues:

```bash
# Rollback to previous version
go get github.com/package/name@v1.2.3
go mod tidy

# Test and commit
go test ./...
git add go.mod go.sum
git commit -m "Rollback github.com/package/name to v1.2.3"
```

---

## Security Considerations

### Security Scanning

Regular security scans should be performed:

```bash
# Use govulncheck (official Go vulnerability scanner)
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Use gosec for additional security checks
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...
```

### Vulnerability Response

1. **Identify affected packages**:
   ```bash
   govulncheck ./...
   ```

2. **Check for available updates**:
   ```bash
   go list -u -m all
   ```

3. **Apply security update**:
   ```bash
   go get github.com/vulnerable/package@v1.2.4
   go mod tidy
   ```

4. **Test thoroughly**:
   ```bash
   go test ./...
   go test -race ./...
   ```

5. **Document the fix**:
   - Update version in this file
   - Add note about security fix
   - Update CHANGELOG

### Known Security Considerations

**API Key Storage**:
- API keys stored in plain text in `~/.promptstack/config.yaml`
- Config file created with restrictive permissions (0600)
- Future enhancement: Consider using system keychain or encrypted storage

**Input Validation**:
- All file paths validated to prevent directory traversal
- User input sanitized before file operations
- Markdown rendering uses `bluemonday` for HTML sanitization

**Prompt Injection**:
- AI suggestions treated as untrusted content
- User must explicitly accept suggestions (no auto-apply)
- Suggestions displayed in read-only view before acceptance

**Dependency Supply Chain**:
- Use `go.sum` for dependency verification
- Regular security scans with `govulncheck`
- Review new dependencies before adding

---

## Dependency Audit

### Regular Audits

Perform dependency audits quarterly:

1. **Check for vulnerabilities**:
   ```bash
   govulncheck ./...
   ```

2. **Review dependency tree**:
   ```bash
   go mod graph
   ```

3. **Check for unused dependencies**:
   ```bash
   go mod tidy
   ```

4. **Review license compliance**:
   - Verify all licenses are documented
   - Check for new dependencies with restrictive licenses
   - Ensure attribution requirements are met

5. **Update documentation**:
   - Update versions in this file
   - Remove deprecated dependencies
   - Add new dependencies with rationale

### Dependency Health Metrics

Track these metrics for each dependency:

- **Last commit date** (should be recent)
- **Release frequency** (active maintenance)
- **Open issues** (manageable backlog)
- **Stars/forks** (community adoption)
- **Security advisories** (CVE history)

### Removing Dependencies

If a dependency is no longer needed:

1. **Identify unused dependencies**:
   ```bash
   go mod tidy
   ```

2. **Remove from code**:
   - Remove imports
   - Replace functionality with alternatives
   - Update tests

3. **Update go.mod**:
   ```bash
   go mod tidy
   ```

4. **Update documentation**:
   - Remove entry from this file
   - Update dependency categories
   - Document removal reason

5. **Test thoroughly**:
   ```bash
   go test ./...
   go test -race ./...
   ```

6. **Commit**:
   ```bash
   git add go.mod go.sum DEPENDENCIES.md
   git commit -m "Remove github.com/unused/package - replaced with [alternative]"
   ```

---

## Appendix

### Complete Dependency Tree

To view the complete dependency tree:

```bash
go mod graph
```

To view dependency tree with versions:

```bash
go mod graph | grep github.com/package/name
```

### License Texts

Full license texts for all dependencies are available in:

- Package repositories (GitHub, GitLab, etc.)
- `go mod download -json github.com/package/name` for license location
- Third-party license notices will be included in distribution

### Resources

- [Go Module Reference](https://golang.org/ref/mod)
- [govulncheck Documentation](https://golang.org/x/vuln/cmd/govulncheck)
- [Go Security Policy](https://go.dev/security/policy)
- [Open Source Initiative](https://opensource.org/licenses)

---

**Document Maintenance**: This document should be updated whenever dependencies are added, updated, or removed. Review quarterly as part of dependency audit process.

**Questions or Issues**: Contact the maintainers or open an issue in the repository.