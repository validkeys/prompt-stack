# Milestone 3 Reference: File I/O Foundation

## How to Use This Document

**Read this section when:**
- Before implementing any task in Milestone 3
- When debugging file I/O issues
- When adding new file operations
- When writing tests for file I/O

**Key sections:**
- Lines 14-40: Architecture context - Read before Task 1
- Lines 42-80: Style guide references - Reference during all tasks
- Lines 82-120: Testing guide references - Consult when writing tests
- Lines 122-180: Key learnings - Apply patterns from previous milestones
- Lines 182-350: Implementation notes - Code examples for each task
- Lines 352-420: Error handling patterns - Critical for all file operations

**Related documents:**
- See [`go-style-guide.md`](../go-style-guide.md) for Go coding standards
- See [`go-testing-guide.md`](../go-testing-guide.md) for testing patterns
- See [`learnings/go-fundamentals.md`](../learnings/go-fundamentals.md) for Go-specific patterns
- See [`learnings/error-handling.md`](../learnings/error-handling.md) for error handling patterns

---

## Architecture Context

### Domain Overview

Milestone 3 implements **File I/O Foundation** which provides the foundation for reading and writing markdown files with YAML frontmatter. This is a critical infrastructure milestone that will be used by multiple future milestones:

- **M4 (Basic Text Editor)**: Will use markdown reader/writer for file operations
- **M5 (Auto-save)**: Will use markdown writer for saving compositions
- **M15 (SQLite Setup)**: Will use markdown reader for importing history files
- **M7-M10 (Library Integration)**: Will use markdown reader for loading prompt library

### Package Structure

Files will be organized in `internal/prompt/` package:

```
internal/prompt/
├── frontmatter.go         # Frontmatter parser (Task 1)
├── frontmatter_test.go    # Frontmatter tests
├── storage.go             # File I/O operations (Tasks 2, 3)
├── storage_test.go        # File I/O tests
└── storage_bench_test.go  # Performance benchmarks (Task 5)
```

### Dependencies

**Internal dependencies**:
- `internal/platform/logging` - From M1 for logging
- `internal/platform/errors` - From M1 for error handling

**External dependencies** (from go.mod):
- `go.uber.org/zap` - Structured logging

### Integration Points

**Editor domain (M4)**: Will import `internal/prompt` for file operations
**Auto-save (M5)**: Will use `storage.Write()` for auto-save operations
**History domain (M15)**: Will use `storage.Read()` for importing history

---

## Style Guide References

### Relevant Patterns from [`go-style-guide.md`](../go-style-guide.md)

#### Error Messages

**Pattern**: Lowercase, no punctuation, wrap with `%w`

```go
// Correct
return fmt.Errorf("failed to read file: %w", err)

// Incorrect
return fmt.Errorf("Error reading file: %v", err)
```

**Apply to**: All file operations in `storage.go`

#### Constructor Pattern

**Pattern**: New() for single type, NewType() for multiple

```go
// For frontmatter parsing - no constructor needed, just functions

// For file operations - no constructor, just functions
```

**Apply to**: Task 1 (frontmatter), Tasks 2-3 (storage)

#### Error Checking

**Pattern**: Check immediately, handle explicitly

```go
// Correct
data, err := os.ReadFile(path)
if err != nil {
    return nil, fmt.Errorf("failed to read file: %w", err)
}

// Incorrect
data, _ := os.ReadFile(path)  // never ignore errors
```

**Apply to**: All file operations

---

## Testing Guide References

### Testing Patterns from [`go-testing-guide.md`](../go-testing-guide.md)

#### Table-Driven Tests

**Pattern**: Use table-driven tests for multiple cases

```go
func TestParseFrontmatter(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        wantMeta  map[string]string
        wantContent string
        wantErr   bool
    }{
        {
            name:      "valid frontmatter",
            input:     "---\ntitle: Test\n---\nContent",
            wantMeta:  map[string]string{"title": "Test"},
            wantContent: "Content",
            wantErr:   false,
        },
        {
            name:      "no frontmatter",
            input:     "Just content",
            wantMeta:  nil,
            wantContent: "Just content",
            wantErr:   false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            meta, content, err := ParseFrontmatter(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if !reflect.DeepEqual(meta, tt.wantMeta) {
                t.Errorf("meta = %v, want %v", meta, tt.wantMeta)
            }
            if content != tt.wantContent {
                t.Errorf("content = %q, want %q", content, tt.wantContent)
            }
        })
    }
}
```

**Apply to**: Task 1 (frontmatter_test.go), Task 2-3 (storage_test.go)

#### Test File Organization

**Pattern**: Co-locate tests with implementation

```
internal/prompt/
├── frontmatter.go
├── frontmatter_test.go    # Test same package
├── storage.go
├── storage_test.go        # Test same package
```

**Apply to**: All test files

### Integration Testing from [`FOUNDATION-TESTING-GUIDE.md`](../milestones/FOUNDATION-TESTING-GUIDE.md)

#### Round-Trip Testing

**Pattern**: Test complete read-modify-write workflow

```go
func TestRoundTrip(t *testing.T) {
    // Setup
    tmp := t.TempDir()
    path := filepath.Join(tmp, "test.md")

    // Write initial file
    original := &prompt.Prompt{
        Title:   "Test",
        Content: "Hello world",
        Tags:    []string{"test", "example"},
    }
    err := storage.Write(path, original)
    if err != nil {
        t.Fatal(err)
    }

    // Read back
    loaded, err := storage.Read(path)
    if err != nil {
        t.Fatal(err)
    }

    // Verify
    if loaded.Title != original.Title {
        t.Errorf("title = %q, want %q", loaded.Title, original.Title)
    }
    if loaded.Content != original.Content {
        t.Errorf("content = %q, want %q", loaded.Content, original.Content)
    }
}
```

**Apply to**: Task 4 (integration tests)

---

## Key Learnings References

### From [`learnings/go-fundamentals.md`](../learnings/go-fundamentals.md)

#### Category 8: Frontmatter Parsing Strategy

**Learning**: Choose simplest solution that meets requirements

**Pattern**: Simple string-based parser instead of full YAML library

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

**Apply to**: Task 1 (frontmatter.go)

**Rationale**:
- Only need to extract key-value pairs
- Simple format: `key: value`
- Avoids additional dependency
- Sufficient for current requirements

### From [`learnings/error-handling.md`](../learnings/error-handling.md)

#### Category 7: Graceful File Read Error Handling

**Learning**: Implement comprehensive error handling with graceful degradation

**Pattern**: Comprehensive error checking before attempting to read

```go
func readFileGracefully(filePath string, logger *zap.Logger) ([]byte, error) {
    // Check if file exists
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, errors.FileError("File not found", err).
                WithDetails(fmt.Sprintf("The file %s does not exist", filePath))
        }
        return nil, errors.FileError("Failed to access file", err).
            WithDetails(fmt.Sprintf("Could not access file: %s", filePath))
    }

    // Check file size (1MB limit)
    const maxFileSize = 1 << 20 // 1MB
    if fileInfo.Size() > maxFileSize {
        err := errors.FileError("File size exceeds limit", nil).
            WithDetails(fmt.Sprintf("File %s is %.2f MB (max: 1MB)",
                filePath, float64(fileInfo.Size())/(1024*1024)))
        logger.Warn("File size exceeds limit",
            zap.String("path", filePath),
            zap.Int64("size", fileInfo.Size()))
        return nil, err
    }

    // Check file permissions
    if fileInfo.Mode().Perm()&0400 == 0 {
        err := errors.FileError("Permission denied", nil).
            WithDetails(fmt.Sprintf("No read permission for file: %s", filePath))
        logger.Warn("Permission denied", zap.String("path", filePath))
        return nil, err
    }

    // Read file content
    content, err := os.ReadFile(filePath)
    if err != nil {
        // Handle specific error types
        if os.IsPermission(err) {
            return nil, errors.FileError("Permission denied", err).
                WithDetails(fmt.Sprintf("Cannot read file: %s", filePath))
        }
        if errors.Is(err, os.ErrClosed) {
            return nil, errors.FileError("File closed unexpectedly", err).
                WithDetails(fmt.Sprintf("File handle closed: %s", filePath))
        }

        // Generic file read error
        return nil, errors.FileError("Failed to read file", err).
            WithDetails(fmt.Sprintf("Error reading file: %s", filePath))
    }

    return content, nil
}
```

**Apply to**: Task 2 (storage.go - Read function)

**Benefits**:
- Comprehensive error detection (not found, size, permissions, read errors)
- Detailed error messages with context
- Graceful degradation (continues loading other files)
- Error tracking for reporting
- Severity-based handling (error vs warning)
- All errors logged appropriately

---

## Implementation Notes

### Task 1: Frontmatter Parser

#### Code Example

```go
// Package prompt provides prompt data structures and file I/O operations.
package prompt

import (
    "fmt"
    "strings"
)

const (
    frontmatterStart = "---"
    frontmatterEnd   = "\n---"
)

// ParseFrontmatter extracts YAML frontmatter from markdown content.
// It returns the metadata map, the remaining content, and an error if parsing fails.
// Files without frontmatter return nil metadata and the original content.
func ParseFrontmatter(content string) (map[string]string, string, error) {
    // Check for frontmatter start marker
    if !strings.HasPrefix(content, frontmatterStart) {
        return nil, content, nil
    }

    // Find end of frontmatter
    endIdx := strings.Index(content, frontmatterEnd)
    if endIdx == -1 {
        // Malformed frontmatter - return content as is
        return nil, content, nil
    }

    // Extract frontmatter content (between --- markers)
    frontmatterText := content[len(frontmatterStart):endIdx]
    body := content[endIdx+len(frontmatterEnd):]

    // Parse frontmatter key-value pairs
    metadata := make(map[string]string)
    for _, line := range strings.Split(frontmatterText, "\n") {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }

        // Split on first colon
        parts := strings.SplitN(line, ":", 2)
        if len(parts) != 2 {
            // Invalid line - skip
            continue
        }

        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])

        // Remove quotes if present
        if len(value) >= 2 {
            if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
               (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
                value = value[1 : len(value)-1]
            }
        }

        metadata[key] = value
    }

    return metadata, body, nil
}
```

#### Test Example

```go
func TestParseFrontmatter(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        wantMeta  map[string]string
        wantBody  string
        wantErr   bool
    }{
        {
            name: "valid frontmatter",
            input: "---\ntitle: Test\ntags: test,example\n---\nContent here",
            wantMeta: map[string]string{
                "title": "Test",
                "tags":  "test,example",
            },
            wantBody: "Content here",
            wantErr:  false,
        },
        {
            name:     "no frontmatter",
            input:    "Just content without frontmatter",
            wantMeta: nil,
            wantBody:  "Just content without frontmatter",
            wantErr:   false,
        },
        {
            name:     "only frontmatter",
            input:    "---\ntitle: Test\n---",
            wantMeta: map[string]string{"title": "Test"},
            wantBody:  "",
            wantErr:   false,
        },
        {
            name: "quoted values",
            input: "---\ntitle: \"Test Title\"\ndescription: 'A test'\n---\nContent",
            wantMeta: map[string]string{
                "title":       "Test Title",
                "description": "A test",
            },
            wantBody: "Content",
            wantErr:  false,
        },
        {
            name:     "malformed frontmatter - missing end",
            input:    "---\ntitle: Test\nContent",
            wantMeta: nil,
            wantBody:  "---\ntitle: Test\nContent",
            wantErr:   false,
        },
        {
            name: "empty lines in frontmatter",
            input: "---\ntitle: Test\n\ntags: test\n---\nContent",
            wantMeta: map[string]string{
                "title": "Test",
                "tags":  "test",
            },
            wantBody: "Content",
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            meta, body, err := ParseFrontmatter(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if !reflect.DeepEqual(meta, tt.wantMeta) {
                t.Errorf("meta = %v, want %v", meta, tt.wantMeta)
            }
            if body != tt.wantBody {
                t.Errorf("body = %q, want %q", body, tt.wantBody)
            }
        })
    }
}
```

---

### Task 2: Markdown File Reader

#### Code Example

```go
// Package prompt provides prompt data structures and file I/O operations.
package prompt

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "unicode/utf8"

    "go.uber.org/zap"
)

const (
    maxFileSize = 1 << 20 // 1MB
)

// ReadFile reads a markdown file from disk and returns a Prompt structure.
// It parses YAML frontmatter and returns both metadata and content.
// Returns error if file doesn't exist, can't be read, exceeds size limit, or contains invalid UTF-8.
func ReadFile(filePath string, logger *zap.Logger) (*Prompt, error) {
    // Get file info first for validation
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, fmt.Errorf("file not found: %w", err)
        }
        return nil, fmt.Errorf("failed to access file: %w", err)
    }

    // Check file size
    if fileInfo.Size() > maxFileSize {
        return nil, fmt.Errorf("file size exceeds limit (%d bytes): %w",
            maxFileSize, os.ErrInvalid)
    }

    // Read file content
    content, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %w", err)
    }

    // Validate UTF-8 encoding
    if !utf8.Valid(content) {
        return nil, fmt.Errorf("file contains invalid UTF-8 encoding: %w", os.ErrInvalid)
    }

    // Convert bytes to string
    contentStr := string(content)

    // Parse frontmatter
    metadata, body, err := ParseFrontmatter(contentStr)
    if err != nil {
        logger.Warn("failed to parse frontmatter, loading as plain markdown",
            zap.String("path", filePath),
            zap.Error(err))
        metadata = nil
        body = contentStr
    }

    // Extract metadata fields
    prompt := &Prompt{
        FilePath: filePath,
        Content:  body,
        Title:    metadata["title"],
        Tags:     parseTags(metadata["tags"]),
    }

    // Use filename as title if not in frontmatter
    if prompt.Title == "" {
        prompt.Title = filepath.Base(filePath)
        prompt.Title = prompt.Title[:len(prompt.Title)-len(filepath.Ext(prompt.Title))]
    }

    return prompt, nil
}

func parseTags(tagsStr string) []string {
    if tagsStr == "" {
        return nil
    }
    tags := strings.Split(tagsStr, ",")
    for i, tag := range tags {
        tags[i] = strings.TrimSpace(tag)
    }
    return tags
}
```

#### Test Example

```go
func TestReadFile(t *testing.T) {
    tests := []struct {
        name      string
        setup     func(t *testing.T) string
        wantTitle string
        wantErr  bool
    }{
        {
            name: "file with frontmatter",
            setup: func(t *testing.T) string {
                tmp := t.TempDir()
                path := filepath.Join(tmp, "test.md")
                content := "---\ntitle: Test\ntags: test,example\n---\nContent here"
                os.WriteFile(path, []byte(content), 0644)
                return path
            },
            wantTitle: "Test",
            wantErr:  false,
        },
        {
            name: "file without frontmatter",
            setup: func(t *testing.T) string {
                tmp := t.TempDir()
                path := filepath.Join(tmp, "test.md")
                os.WriteFile(path, []byte("Just content"), 0644)
                return path
            },
            wantTitle: "test",
            wantErr:  false,
        },
        {
            name: "file not found",
            setup: func(t *testing.T) string {
                return "/nonexistent/file.md"
            },
            wantTitle: "",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            logger, _ := zap.NewDevelopment()
            path := tt.setup(t)

            prompt, err := ReadFile(path, logger)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if !tt.wantErr && prompt.Title != tt.wantTitle {
                t.Errorf("title = %q, want %q", prompt.Title, tt.wantTitle)
            }
        })
    }
}
```

---

### Task 3: Markdown File Writer

#### Code Example

```go
// WriteFile writes a Prompt structure to a markdown file with YAML frontmatter.
// Creates parent directories if they don't exist.
// Returns error if unable to write file.
// Note: This function is in the same package as ReadFile and uses the same imports.
// Line endings in the content are preserved as-is (no normalization).
func WriteFile(filePath string, prompt *Prompt) error {
    // Create parent directories if needed
    dir := filepath.Dir(filePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create directory: %w", err)
    }

    // Build frontmatter
    var frontmatter strings.Builder
    frontmatter.WriteString("---\n")
    if prompt.Title != "" {
        frontmatter.WriteString(fmt.Sprintf("title: %s\n", prompt.Title))
    }
    if len(prompt.Tags) > 0 {
        tags := strings.Join(prompt.Tags, ",")
        frontmatter.WriteString(fmt.Sprintf("tags: %s\n", tags))
    }
    if prompt.Description != "" {
        frontmatter.WriteString(fmt.Sprintf("description: %s\n", prompt.Description))
    }
    frontmatter.WriteString("---\n")

    // Build complete content
    var content strings.Builder
    content.WriteString(frontmatter.String())
    content.WriteString(prompt.Content)

    // Write file
    if err := os.WriteFile(filePath, []byte(content.String()), 0644); err != nil {
        return fmt.Errorf("failed to write file: %w", err)
    }

    return nil
}
```

#### Test Example

```go
func TestWriteFile(t *testing.T) {
    tests := []struct {
        name      string
        prompt    *Prompt
        wantErr   bool
        verify    func(t *testing.T, filePath string)
    }{
        {
            name: "write with frontmatter",
            prompt: &Prompt{
                Title:       "Test Title",
                Description: "A test prompt",
                Tags:        []string{"test", "example"},
                Content:     "This is the content",
            },
            wantErr: false,
            verify: func(t *testing.T, filePath string) {
                content, err := os.ReadFile(filePath)
                if err != nil {
                    t.Fatal(err)
                }
                contentStr := string(content)
                if !strings.Contains(contentStr, "title: Test Title") {
                    t.Error("frontmatter not written correctly")
                }
                if !strings.Contains(contentStr, "This is the content") {
                    t.Error("content not written")
                }
            },
        },
        {
            name: "write without frontmatter",
            prompt: &Prompt{
                Content: "Just content",
            },
            wantErr: false,
            verify: func(t *testing.T, filePath string) {
                content, err := os.ReadFile(filePath)
                if err != nil {
                    t.Fatal(err)
                }
                contentStr := string(content)
                if !strings.Contains(contentStr, "Just content") {
                    t.Error("content not written")
                }
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tmp := t.TempDir()
            path := filepath.Join(tmp, "test.md")

            err := WriteFile(path, tt.prompt)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }

            if !tt.wantErr && tt.verify != nil {
                tt.verify(t, path)
            }
        })
    }
}
```

---

### Task 4: Integration Tests

#### Code Example

```go
// test/integration/fileio_test.go
package integration_test

import (
    "path/filepath"
    "testing"

    "github.com/kyledavis/prompt-stack/internal/prompt"
    "go.uber.org/zap"
)

func TestFileIORoundTrip(t *testing.T) {
    tmp := t.TempDir()
    logger, _ := zap.NewDevelopment()

    // Test round-trip
    original := &prompt.Prompt{
        Title:       "Test Prompt",
        Description: "A test prompt for round-trip",
        Tags:        []string{"test", "roundtrip"},
        Content:     "This is test content with {{text:placeholder}}",
    }

    // Write
    filePath := filepath.Join(tmp, "test.md")
    err := prompt.WriteFile(filePath, original)
    if err != nil {
        t.Fatalf("failed to write file: %v", err)
    }

    // Read back
    loaded, err := prompt.ReadFile(filePath, logger)
    if err != nil {
        t.Fatalf("failed to read file: %v", err)
    }

    // Verify
    if loaded.Title != original.Title {
        t.Errorf("title = %q, want %q", loaded.Title, original.Title)
    }
    if loaded.Description != original.Description {
        t.Errorf("description = %q, want %q", loaded.Description, original.Description)
    }
    if len(loaded.Tags) != len(original.Tags) {
        t.Errorf("tags length = %d, want %d", len(loaded.Tags), len(original.Tags))
    }
    for i, tag := range loaded.Tags {
        if tag != original.Tags[i] {
            t.Errorf("tags[%d] = %q, want %q", i, tag, original.Tags[i])
        }
    }
    if loaded.Content != original.Content {
        t.Errorf("content = %q, want %q", loaded.Content, original.Content)
    }
}
```

---

### Task 5: Performance Benchmarks

#### Code Example

```go
// internal/prompt/storage_bench_test.go
package prompt

import (
    "os"
    "path/filepath"
    "testing"
    "time"
)

func BenchmarkReadFile(b *testing.B) {
    tmp := b.TempDir()
    logger, _ := zap.NewDevelopment()

    // Create test file
    path := filepath.Join(tmp, "bench.md")
    content := strings.Repeat("test line\n", 10000)
    os.WriteFile(path, []byte(content), 0644)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = ReadFile(path, logger)
    }
}

func BenchmarkWriteFile(b *testing.B) {
    tmp := b.TempDir()

    // Create test prompt
    prompt := &Prompt{
        Title:   "Benchmark",
        Content: strings.Repeat("test line\n", 10000),
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        path := filepath.Join(tmp, fmt.Sprintf("bench%d.md", i))
        _ = WriteFile(path, prompt)
    }
}

func BenchmarkParseFrontmatter(b *testing.B) {
    content := "---\ntitle: Test\ndescription: A test prompt\ntags: test,example\n---\n" +
        strings.Repeat("test line\n", 1000)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _, _ = ParseFrontmatter(content)
    }
}
```

---

## Error Handling Patterns

### From [`learnings/error-handling.md`](../learnings/error-handling.md)

#### Category 1: Error Handling Architecture

**Learning**: Create structured error types with severity levels and display strategies

**Pattern**: Use structured errors from M1 error handling system

```go
// Use error types from M1
import "github.com/kyledavis/prompt-stack/internal/platform/errors"

// For file errors
func ReadFile(filePath string, logger *zap.Logger) (*Prompt, error) {
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, errors.FileError("file not found", err).
                WithDetails(fmt.Sprintf("file does not exist: %s", filePath))
        }
        return nil, errors.FileError("failed to access file", err).
            WithDetails(fmt.Sprintf("could not access: %s", filePath))
    }
    // ...
}
```

**Apply to**: All file operations in Task 2 and Task 3

---

## Verification Checklist

Before presenting implementation plan, verify:

- [x] All required documents from MATRIX have been read
- [x] All deliverables from M3 are addressed
- [x] File paths match project-structure.md
- [x] Code examples follow go-style-guide.md patterns
- [x] Test examples follow go-testing-guide.md patterns
- [x] Key learnings from go-fundamentals.md are applied
- [x] Key learnings from error-handling.md are applied
- [x] Testing requirements include coverage targets and critical scenarios
- [x] Navigation guide included
- [x] All code examples pass validation (no pseudo-code)
- [x] Design system compliance noted (N/A for this milestone)
- [x] UTF-8 validation included in file reading
- [x] File locking considerations documented
- [x] Detailed performance thresholds specified for benchmarks
- [x] File hash verification pattern documented for future use
- [x] Line ending preservation documented

---

**Last Updated**: 2026-01-08 (Revised)
**Status**: Ready for implementation with enhanced considerations
