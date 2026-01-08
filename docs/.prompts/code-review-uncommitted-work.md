# Code Review Prompt for Uncommitted Changes

Please perform a comprehensive code review of the uncommitted changes in the working directory. Evaluate the implementation against our established standards and guidelines before committing:

## Review Scope

Focus on all uncommitted changes including:
- Modified files (M)
- New untracked files (??)
- Staged changes

## Review Criteria

### 1. Go Code Quality & Idioms
Assess adherence to [`go-style-guide.md`](../plans/fresh-build/go-style-guide.md):
- **Package Organization**: Singular, lowercase names; proper package comments
- **Type Design**: Appropriate use of constructors (`New()` pattern), struct encapsulation
- **Method Receivers**: Consistent pointer vs value receivers
- **Error Handling**: Lowercase messages, proper wrapping with `%w`, contextual error information
- **Interfaces**: Defined at usage site, small and focused (3-5 methods max)
- **Naming Conventions**: `camelCase` for unexported, `PascalCase` for exported
- **Code Organization**: File naming, function length (<50 lines ideal), meaningful comments

### 2. Project Structure Compliance
Verify alignment with [`project-structure.md`](../plans/fresh-build/project-structure.md):
- **Domain Boundaries**: Clear separation between domains (editor, prompt, library, history, ai, config, platform, ui)
- **Dependency Direction**: Proper flow (ui → internal/domain → internal/platform)
- **File Organization**: Standard pattern (`entity.go`, `entity_test.go`, `operation.go`, etc.)
- **Interface Placement**: Interfaces defined where used, not where implemented
- **Theme System**: Consistent use of theme helpers from [`ui/theme/theme.go`](../../ui/theme/theme.go)

### 3. Testing Standards
Evaluate test coverage per [`go-testing-guide.md`](../plans/fresh-build/go-testing-guide.md):
- **Test Organization**: Co-located tests with `_test.go` suffix
- **Test Patterns**: Table-driven tests for multiple cases
- **Effect Testing**: Testing observable outcomes, not implementation details
- **Mock Usage**: Proper interface mocking for dependencies
- **Coverage**: Aim for >80% test coverage
- **Test Naming**: Descriptive names following `TestType_Method` pattern
- **Bubble Tea Testing**: For TUI components, follow [`bubble-tea-testing-best-practices.md`](../plans/fresh-build/bubble-tea-testing-best-practices.md) - test message handling effects, not implementation details

### 4. Requirements Adherence
Check implementation against [`requirements.md`](../plans/fresh-build/requirements.md):
- **Feature Completeness**: All specified features implemented correctly
- **Placeholder System**: Proper syntax validation (`{{type:name}}`)
- **Error Handling**: Appropriate error messages and recovery strategies
- **Performance**: Meeting specified performance targets
- **Security**: Input validation, API key handling

### 5. Architecture Patterns
Verify proper use of established patterns:
- **Factory Pattern**: For provider/repository instantiation
- **Repository Pattern**: For data access abstraction
- **Middleware Pattern**: For cross-cutting concerns
- **Event Pattern**: For component decoupling
- **Dependency Injection**: Explicit dependency passing

### 6. Bubble Tea UI Patterns
Assess adherence to [`bubble-tea-best-practices.md`](../../bubble-tea-best-practices.md):
- **Elm Architecture**: Clean separation of Model, Update, and View functions
- **Model Immutability**: Update returns modified model, no side effects in View
- **Message Handling**: All message types properly handled in Update with type assertions
- **Command Pattern**: Async I/O operations using Commands, not direct I/O in Update/View
- **View Purity**: No state updates or I/O in View function
- **Initialization**: Proper Init() implementation returning initial Commands
- **Testing**: Message simulation and effect verification per [`bubble-tea-testing-best-practices.md`](../plans/fresh-build/bubble-tea-testing-best-practices.md)

### 7. Code Quality Issues
Identify any:
- **Anti-patterns**: `init()` side effects, panic in library code, naked returns, generic variable names
- **Code Smells**: Long functions, tight coupling, hidden dependencies
- **Performance Issues**: Inefficient algorithms, unnecessary allocations
- **Security Concerns**: Input validation gaps, unsafe operations

### 8. Documentation
Assess:
- **Package Comments**: Clear package-level documentation
- **Function Comments**: Exported functions properly documented
- **Inline Comments**: Explain "why" not "what"
- **README Updates**: If applicable

### 9. Commit Readiness
Evaluate if changes are ready to commit:
- **Completeness**: Is this a logical, complete unit of work?
- **Build Status**: Does the code compile without errors?
- **Test Status**: Do all tests pass?
- **No Debug Code**: Remove console logs, debug statements, commented code
- **No Secrets**: Ensure no API keys, passwords, or sensitive data
- **Git Status**: Are all necessary files staged? Any accidental inclusions?

## Output Format

Please structure your review as follows:

### ✅ Strengths
- List what was done well
- Highlight good patterns and practices

### ⚠️ Issues Found
For each issue, provide:
- **Severity**: Critical / Major / Minor
- **Category**: Code Quality / Structure / Testing / Requirements / Performance / Security / Commit Readiness
- **Location**: File path and line number
- **Description**: What's wrong and why it matters
- **Recommendation**: Specific fix or improvement

### 🚫 Blockers for Commit
List any critical issues that MUST be fixed before committing:
- Build errors
- Failing tests
- Security vulnerabilities
- Accidental inclusion of secrets or debug code

### 📋 Suggestions
- Optional improvements
- Refactoring opportunities
- Performance optimizations

### 📊 Metrics
- Test coverage percentage
- Lines of code added/modified
- Number of files changed
- Complexity assessment

### ✅ Commit Recommendation
- **READY TO COMMIT**: All checks passed, safe to commit
- **NEEDS FIXES**: Address blockers before committing
- **NEEDS DISCUSSION**: Architectural or design decisions needed

## Example Issue Format

```
⚠️ **Major - Code Quality**
Location: `internal/library/loader.go:45`
Issue: Using global state for configuration access
```go
func (l *Loader) LoadPrompts() {
    dir := config.Global.DataDir // ❌ Global state
}
```
Recommendation: Pass config via dependency injection
```go
func (l *Loader) LoadPrompts(dir string) {
    // Use injected parameter
}
```
```

## Uncommitted Files to Review

Based on the current git status, review the following:

**Modified Files:**
- archive/code/ui/app/model.go
- archive/code/ui/palette/model.go
- archive/code/ui/workspace/model.go
- cmd/promptstack/main.go
- docs/plans/fresh-build/DOCUMENT-INDEX.md
- docs/plans/fresh-build/DOCUMENT-REFERENCE-MATRIX.md
- docs/plans/fresh-build/go-testing-guide.md
- go.mod
- go.sum

**New Files:**
- archive/code/internal/testutil/
- archive/code/ui/app/model_test.go
- archive/code/ui/palette/model_test.go
- archive/code/ui/suggestions/model_test.go
- archive/code/ui/workspace/model_test.go
- archive/code/ui/workspace/test content/
- bubble-tea-best-practices.md
- cmd/promptstack/main_test.go
- docs/.prompts/
- docs/plans/fresh-build/bubble-tea-testing-best-practices.md
- docs/plans/fresh-build/code-reviews/
- docs/plans/fresh-build/milestone-implementation-plans/M2/
- internal/editor/
- internal/testutil/
- issues/
- ui/

Please be thorough but constructive in your feedback, focusing on ensuring code quality and commit readiness. Flag any issues that would make this commit problematic or incomplete.
