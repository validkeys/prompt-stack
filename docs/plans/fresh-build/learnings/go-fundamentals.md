# Go Fundamentals Key Learnings

**Purpose**: Key learnings and implementation patterns for Go-specific patterns and pitfalls from previous PromptStack implementation.

**Related Milestones**: M1-M6 (Foundation), M15 (SQLite Setup)

**Related Documents**: 
- [`project-structure.md`](../project-structure.md) - Go project structure
- [`go-style-guide.md`](../go-style-guide.md) - Go coding standards
- [`go-testing-guide.md`](../go-testing-guide.md) - Go testing patterns

---

## Learning Categories

### Category 1: Go Embed Limitations

**Learning**: `go:embed` does not support parent directory references (`..`)

**Problem**: Initially tried to embed `starter-prompts` from `internal/bootstrap/starter.go` using `//go:embed ../../starter-prompts`, which resulted in "invalid pattern syntax" error.

**Solution**: Created the embed file at the root level (`starter.go`) where it can directly reference `starter-prompts` without parent directory traversal.

**Implementation Pattern**:
```go
// Correct: Place embed directive at same level or higher than target
//go:embed starter-prompts
var starterFS embed.FS

// Incorrect: Never use parent directory references
//go:embed ../../starter-prompts // This will fail
```

**Lesson**: Always place `go:embed` directives in files at the same directory level or higher than the target directory. Never use `..` in embed patterns.

**Related Milestones**: M1, M2

**When to Apply**: When embedding static files or resources in Go applications

---

### Category 2: Zap Logger Structured Fields

**Learning**: Zap requires structured field objects, not string literals

**Problem**: Initially used `logger.Info("message", "key", value)` which caused compilation errors about untyped string constants.

**Solution**: Use zap field constructors: `logger.Info("message", zap.String("key", value))`

**Implementation Pattern**:
```go
import "go.uber.org/zap"

// Correct: Use zap field constructors
logger.Info("User logged in",
    zap.String("username", "john"),
    zap.Int("userID", 123),
    zap.Error(err),
)

// Incorrect: String literals cause compilation errors
logger.Info("User logged in", "username", "john") // Error!
```

**Lesson**: Always use zap's field constructors (`zap.String()`, `zap.Int()`, `zap.Error()`, etc.) for structured logging. This provides type safety and better performance.

**Related Milestones**: M1, M2

**When to Apply**: When using zap logger for structured logging throughout the application

---

### Category 3: Regex Matching in Go

**Learning**: Different regex methods return different data types

**Problem**: Used `FindAllStringSubmatchIndex()` which returns `[][]int` (positions), but needed actual string matches.

**Solution**: Switched to `FindAllStringSubmatch()` which returns `[][]string` with the actual matched text.

**Implementation Pattern**:
```go
import "regexp"

// For placeholder parsing - need actual text matches
re := regexp.MustCompile(`\{\{(\w+):(\w+)\}\}`)
matches := re.FindAllStringSubmatch(content, -1) // Returns [][]string

for _, match := range matches {
    // match[0] = full match
    // match[1] = first capture group
    // match[2] = second capture group
}

// If you need positions instead
positions := re.FindAllStringSubmatchIndex(content, -1) // Returns [][]int
```

**Lesson**: Carefully choose the right regex method based on whether you need positions or actual matches. For placeholder parsing, we needed the actual text.

**Related Milestones**: M4, M5

**When to Apply**: When parsing structured text with regex patterns

---

### Category 4: SQLite Driver Selection

**Learning**: Choose pure Go implementation over CGO for build simplicity

**Problem**: Need SQLite database for history management. Multiple drivers available with different trade-offs.

**Solution**: Chose `modernc.org/sqlite` over `github.com/mattn/go-sqlite3`

**Implementation Pattern**:
```go
import (
    "database/sql"
    _ "modernc.org/sqlite" // Pure Go SQLite driver
)

// Open database connection
db, err := sql.Open("sqlite", dbPath)
if err != nil {
    return fmt.Errorf("failed to open database: %w", err)
}
```

**Rationale**:
- Pure Go implementation (no CGO dependency)
- Simplifies cross-platform builds (macOS Intel/ARM)
- Adequate performance for personal-scale usage
- FTS5 support for full-text search

**Trade-off**: Slightly slower than CGO-based driver, but build simplicity outweighs performance for this use case.

**Lesson**: Consider build complexity vs. performance trade-offs when choosing dependencies. For CLI tools distributed as binaries, pure Go implementations often win.

**Related Milestones**: M15, M16

**When to Apply**: When selecting database drivers for Go applications, especially CLI tools

---

### Category 5: Go Version Requirements

**Learning**: Some packages require newer Go versions

**Problem**: `modernc.org/sqlite` required Go 1.24+, which was newer than the installed version (1.23.2).

**Solution**: Running `go get` automatically upgraded the Go toolchain to 1.24.11.

**Implementation Pattern**:
```bash
# Go toolchain automatically manages version requirements
go get modernc.org/sqlite

# Check current Go version
go version

# The toolchain downloads and uses required version automatically
```

**Lesson**: Be aware of Go version requirements in dependencies. The Go toolchain can manage multiple versions, but this may surprise users with older Go installations.

**Related Milestones**: M1, M15

**When to Apply**: When adding new dependencies with specific Go version requirements

---

### Category 6: Project Structure Organization

**Learning**: Organize internal packages by domain/feature rather than technical layer

**Problem**: Need clear, maintainable project structure for a CLI application with multiple features.

**Solution**: Standard Go project layout with feature-based internal packages

**Implementation Pattern**:
```
cmd/promptstack/    # Main application entry point
internal/            # Private packages organized by feature
  ├── config/      # Configuration management
  ├── setup/        # First-run setup
  ├── bootstrap/    # Application initialization
  ├── library/      # Prompt library management
  ├── history/      # History and database
  ├── prompt/       # Prompt data models
  └── logging/      # Logging setup
ui/                 # TUI components
starter.go           # Embedded resources at root
```

**Benefits**:
- Clear separation of concerns
- Easy to locate code by feature
- Standard Go conventions
- Internal packages are truly private

**Lesson**: Organize internal packages by domain/feature rather than technical layer. This makes the codebase more navigable and maintainable.

**Related Milestones**: M1-M6

**When to Apply**: When structuring new Go projects or refactoring existing ones

---

### Category 7: Error Handling Patterns

**Learning**: Use `fmt.Errorf` with `%w` for error wrapping

**Problem**: Need to add context to errors while preserving the original error for debugging.

**Solution**: Use error wrapping with `%w` verb

**Implementation Pattern**:
```go
import (
    "errors"
    "fmt"
)

// Wrap errors with context
if err != nil {
    return nil, fmt.Errorf("failed to load config: %w", err)
}

// Check for specific errors
if errors.Is(err, os.ErrNotExist) {
    // Handle file not found
}

// Extract wrapped error
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    // Handle path error
}
```

**Benefits**:
- Preserves original error for unwrapping
- Adds context at each layer
- Enables `errors.Is()` and `errors.As()` checks
- Clear error messages for debugging

**Lesson**: Always wrap errors with context using `%w`. Never discard the original error. This makes debugging and error handling much easier.

**Related Milestones**: M1-M38 (All milestones)

**When to Apply**: When handling errors throughout the application

---

### Category 8: Frontmatter Parsing Strategy

**Learning**: Choose the simplest solution that meets requirements

**Problem**: Need to parse frontmatter from markdown files (key: value pairs).

**Solution**: Simple string-based parser instead of full YAML library

**Implementation Pattern**:
```go
func parseFrontmatter(content string) (map[string]string, string, error) {
    // Check for --- markers
    if !strings.HasPrefix(content, "---") {
        return nil, content, nil
    }
    
    // Find end of frontmatter
    endIdx := strings.Index(content, "\n---")
    if endIdx == -1 {
        return nil, content, nil
    }
    
    // Parse frontmatter lines
    frontmatter := content[4:endIdx]
    metadata := make(map[string]string)
    
    for _, line := range strings.Split(frontmatter, "\n") {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }
        
        parts := strings.SplitN(line, ":", 2)
        if len(parts) == 2 {
            key := strings.TrimSpace(parts[0])
            value := strings.TrimSpace(parts[1])
            metadata[key] = value
        }
    }
    
    // Return metadata and remaining content
    return metadata, content[endIdx+4:], nil
}
```

**Rationale**:
- Only need to extract key-value pairs
- Simple format: `key: value`
- Avoids additional dependency
- Sufficient for current requirements

**Trade-off**: Less robust than full YAML parser, but adequate for simple frontmatter.

**Lesson**: Choose the simplest solution that meets requirements. Don't over-engineer for future needs that may never materialize.

**Related Milestones**: M4, M5

**When to Apply**: When parsing simple structured data formats

---

## Quick Reference

| Learning | Milestone | Priority |
|----------|-----------|----------|
| Go Embed Limitations | M1, M2 | High |
| Zap Logger Structured Fields | M1, M2 | High |
| Regex Matching in Go | M4, M5 | High |
| SQLite Driver Selection | M15, M16 | High |
| Go Version Requirements | M1, M15 | Medium |
| Project Structure Organization | M1-M6 | High |
| Error Handling Patterns | M1-M38 | High |
| Frontmatter Parsing Strategy | M4, M5 | Medium |

---

**Last Updated**: 2026-01-07
**Source**: [`key-learnings.md`](../key-learnings.md)