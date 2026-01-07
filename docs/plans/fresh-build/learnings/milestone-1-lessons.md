# Milestone 1: Lessons Learned

**Date:** 2026-01-07  
**Milestone:** M1 - Bootstrap & Config  
**Status:** Complete with refactoring needed

---

## Overview

This document captures key lessons from Milestone 1 implementation to prevent similar issues in future milestones.

---

## Critical Issues Identified

### 1. Package Structure Violations

**Issue:** Packages were created in wrong locations, not following project-structure.md

**Examples:**
- `internal/setup/` created instead of `internal/config/setup.go`
- `internal/bootstrap/` created instead of `internal/platform/bootstrap/`
- `internal/logging/` created instead of `internal/platform/logging/`

**Root Cause:** Did not verify package locations against project-structure.md before implementation

**Prevention:**
- ✅ Always consult project-structure.md BEFORE creating any package
- ✅ Verify domain classification (config, platform, domain, ui)
- ✅ Check if functionality belongs in existing package vs new package
- ✅ Use pre-implementation checklist in task-list.md

**Pattern to Follow:**
```
Domain Classification:
- Config domain → internal/config/
- Platform/Infrastructure → internal/platform/{subdomain}/
- Business logic → internal/{domain}/
- UI components → ui/{component}/
```

---

### 2. Global State Anti-Pattern

**Issue:** Global logger variable with mutex created instead of dependency injection

**Example:**
```go
// ❌ BAD - Global state
var (
    globalLogger *zap.Logger
    loggerMutex  sync.RWMutex
)

func GetLogger() (*zap.Logger, error) {
    loggerMutex.RLock()
    defer loggerMutex.RUnlock()
    return globalLogger, nil
}
```

**Root Cause:** Convenience over proper design patterns

**Prevention:**
- ✅ Never use global variables for dependencies
- ✅ Always pass dependencies through constructors
- ✅ Use dependency injection pattern consistently
- ✅ Check pre-implementation checklist for global state

**Pattern to Follow:**
```go
// ✅ GOOD - Dependency injection
func New() (*zap.Logger, error) {
    // ... setup code ...
    return logger, nil
}

// Pass logger explicitly
func main() {
    logger, err := logging.New()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to initialize logging: %v\n", err)
        os.Exit(1)
    }
    defer logger.Sync()
    
    // Pass to other components
    if err := bootstrap.Run(logger); err != nil {
        logger.Error("Bootstrap failed", zap.Error(err))
        os.Exit(1)
    }
}
```

---

### 3. Missing Package Documentation

**Issue:** No package-level comments in any file

**Root Cause:** Documentation not considered part of implementation

**Prevention:**
- ✅ Add package comment immediately when creating package
- ✅ Include package comment in pre-implementation checklist
- ✅ Verify package comments before task completion

**Pattern to Follow:**
```go
// Package config provides application configuration management including
// loading, validation, and persistence of user settings.
package config
```

---

### 4. Missing Tests

**Issue:** No test files created during implementation

**Root Cause:** Not following TDD strictly

**Prevention:**
- ✅ Write tests BEFORE implementation (TDD red-green-refactor)
- ✅ Create `{file}_test.go` alongside `{file}.go`
- ✅ Use table-driven tests from go-testing-guide.md
- ✅ Verify test coverage > 80% before task completion

**Pattern to Follow:**
```go
// internal/config/config_test.go
package config_test

import (
    "testing"
    "github.com/kyledavis/prompt-stack/internal/config"
)

func TestConfig_Validate(t *testing.T) {
    tests := []struct {
        name    string
        config  config.Config
        wantErr bool
    }{
        {
            name: "valid config",
            config: config.Config{
                ClaudeAPIKey: "sk-ant-test123",
                Model:        "claude-3-sonnet-20240229",
            },
            wantErr: false,
        },
        {
            name: "missing api key",
            config: config.Config{
                Model: "claude-3-sonnet-20240229",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.config.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

---

### 5. Hardcoded Magic Strings

**Issue:** Validation rules and constants embedded in code

**Example:**
```go
// ❌ BAD - Magic string
if !strings.HasPrefix(apiKey, "sk-ant-") {
    fmt.Println("Invalid API key format. API keys should start with 'sk-ant-'")
}
```

**Root Cause:** Not extracting constants during implementation

**Prevention:**
- ✅ Identify magic strings during planning phase
- ✅ Extract to constants immediately
- ✅ Include constant extraction in pre-implementation checklist

**Pattern to Follow:**
```go
// ✅ GOOD - Constants
const (
    APIKeyPrefix    = "sk-ant-"
    APIKeyMinLength = 20
)

func validateAPIKey(apiKey string) error {
    if !strings.HasPrefix(apiKey, APIKeyPrefix) {
        return fmt.Errorf("api key must start with %s", APIKeyPrefix)
    }
    if len(apiKey) < APIKeyMinLength {
        return fmt.Errorf("api key too short")
    }
    return nil
}
```

---

## Medium Priority Issues

### 6. Error Message Formatting

**Issue:** Inconsistent error message style (some with underscores, some capitalized)

**Example:**
```go
// ❌ BAD
return fmt.Errorf("claude_api_key is required")

// ✅ GOOD
return fmt.Errorf("claude api key is required")
```

**Prevention:**
- ✅ Follow go-style-guide.md: lowercase, no punctuation
- ✅ Review all error messages before task completion
- ✅ Use linter to catch style violations

---

### 7. Custom Error Types Not Used

**Issue:** Generic errors used instead of structured validation errors

**Prevention:**
- ✅ Consider custom error types for validation
- ✅ Use structured errors for better error handling
- ✅ Reference error-handling.md patterns

**Pattern to Follow:**
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}
```

---

## Process Improvements

### What Worked Well

1. **Clean separation of concerns** - Each file had clear responsibility
2. **Good error handling** - Consistent use of error wrapping with `%w`
3. **Proper resource management** - Used `defer` for cleanup
4. **Validation logic** - API key and config validation implemented

### What Needs Improvement

1. **Pre-implementation verification** - Need checklist before coding
2. **TDD discipline** - Tests must come first, not after
3. **Documentation as code** - Package comments written with code
4. **Constant extraction** - Identify and extract during planning

### New Process Requirements

1. **Pre-Implementation Checklist** (Added to milestone-execution-prompt.md)
   - Verify package structure against project-structure.md
   - Check for global state in design
   - Plan package documentation
   - Plan test files alongside implementation
   - Identify constants to extract

2. **During Implementation**
   - Write test first (TDD red phase)
   - Implement minimal code (TDD green phase)
   - Refactor and add documentation (TDD refactor phase)
   - Verify against checklist before completion

3. **Before Task Completion**
   - All tests passing
   - Package comments present
   - No global state
   - Constants extracted
   - Error messages follow style guide

---

## Refactoring Required

### High Priority

1. **Package Restructuring**
   - Move `internal/setup/` → `internal/config/`
   - Move `internal/bootstrap/` → `internal/platform/bootstrap/`
   - Move `internal/logging/` → `internal/platform/logging/`

2. **Remove Global State**
   - Change `logging.Initialize()` to `logging.New()`
   - Remove global logger variables
   - Pass logger through constructors

3. **Add Documentation**
   - Package comments for all packages
   - Function comments for exported functions

4. **Add Tests**
   - `internal/config/config_test.go`
   - `internal/config/setup_test.go`
   - `internal/platform/logging/logger_test.go`
   - `internal/platform/bootstrap/bootstrap_test.go`

### Medium Priority

5. **Extract Constants**
   - API key validation rules
   - Model names

6. **Standardize Error Messages**
   - Review all error messages
   - Ensure lowercase, no punctuation

---

## Key Takeaways

### For Future Milestones

1. **Always verify package structure BEFORE implementation**
   - Consult project-structure.md
   - Verify domain classification
   - Check if new package is needed

2. **Never use global state**
   - Pass dependencies explicitly
   - Use dependency injection
   - Makes testing easier

3. **Write tests first (TDD)**
   - Red: Write failing test
   - Green: Implement minimal code
   - Refactor: Improve code quality

4. **Document as you code**
   - Package comments immediately
   - Function comments for exports
   - Explain WHY, not WHAT

5. **Extract constants early**
   - Identify during planning
   - Extract during implementation
   - No magic strings in code

### Quality Metrics

**Target for Future Milestones:**
- Package structure: 100% compliance with project-structure.md
- Global state: 0 instances
- Package documentation: 100% coverage
- Test coverage: > 80% for new code
- Error message style: 100% compliance with go-style-guide.md
- Constants: 0 magic strings in code

---

## References

- [Code Review](../milestone-implementation-plans/M1/code-review.md)
- [Go Style Guide](../go-style-guide.md)
- [Project Structure](../project-structure.md)
- [Go Testing Guide](../go-testing-guide.md)
- [Error Handling Patterns](error-handling.md)

---

**Next Review:** After M1 refactoring is complete