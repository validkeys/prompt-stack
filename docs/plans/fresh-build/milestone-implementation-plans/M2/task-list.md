# Milestone 2: Basic TUI Shell - Task List

**Milestone Number**: 2  
**Title**: Basic TUI Shell  
**Goal**: Render functional TUI with quit handling

---

## Overview

This milestone establishes the foundational TUI (Terminal User Interface) for PromptStack using the Bubble Tea framework. The focus is on creating a minimal but functional shell that can render a basic UI, handle keyboard input, and provide clean quit functionality.

### Deliverables
- Root Bubble Tea model at [`ui/app/model.go`](../../project-structure.md:227)
- Basic status bar component at [`ui/statusbar/model.go`](../../project-structure.md:247)
- Keyboard input handling (character input, special keys)
- Clean quit functionality (Ctrl+C, 'q' key)

### Dependencies
- Milestone 1: Bootstrap & Config (config loading, logging setup)

### Integration Points
- Config system from M1 for theme preferences
- Logging system from M1 for TUI debugging
- Theme system (Catppuccin Mocha colors)

---

## Pre-Implementation Checklist

Before writing any code, verify:

### Package Structure
- [ ] All file paths match [`project-structure.md`](../../project-structure.md)
- [ ] Packages are in correct domain (ui/app, ui/statusbar, ui/theme)
- [ ] No packages in wrong locations (e.g., setup as separate package)

### Dependency Injection
- [ ] No global state variables planned
- [ ] All dependencies passed through constructors
- [ ] Logger passed explicitly, not accessed globally

### Documentation
- [ ] Package comments planned for all new packages
- [ ] Exported function comments planned
- [ ] Error messages follow lowercase, no punctuation style

### Testing
- [ ] Test files planned alongside implementation files
- [ ] Table-driven test structure planned
- [ ] Mock interfaces identified for testing

### Style Compliance
- [ ] Constructor naming follows New() or NewType() pattern
- [ ] Error wrapping uses %w consistently
- [ ] Method receivers are consistent (all pointer or all value)
- [ ] No stuttering in names (e.g., app.AppModel → app.Model)

### Constants
- [ ] Magic strings identified for extraction to constants
- [ ] Validation rules defined as constants
- [ ] No hardcoded values in implementation

**If any item is unchecked, review and adjust plan before proceeding.**

---

## Tasks

### Task 1: Create Theme System

**Dependencies**: None  
**Files**: [`ui/theme/theme.go`](../../project-structure.md:267)  
**Integration Points**: None (foundation for other tasks)  
**Estimated Complexity**: Low

**Description**: Create a centralized theme package with Catppuccin Mocha color palette and style helper functions. This provides the visual foundation for all TUI components.

**Acceptance Criteria**:
- [x] Package MUST define all Catppuccin Mocha color constants (BackgroundPrimary, BackgroundSecondary, ForegroundPrimary, ForegroundMuted, AccentBlue, AccentGreen, AccentYellow, AccentRed, BorderPrimary)
- [x] Package MUST provide style helper functions (ModalStyle, StatusStyle, ActivePlaceholderStyle)
- [x] Style functions MUST return lipgloss.Style objects
- [x] Colors MUST match exact hex values from Catppuccin Mocha specification
- [x] Package MUST have package-level documentation comment

**Testing Requirements**:
- Coverage Target: > 90%
- Critical Test Scenarios:
  - Verify all color constants are defined with correct hex values
  - Verify style functions return non-nil lipgloss.Style objects
  - Verify style functions apply correct colors and properties
- Edge Cases:
  - Test that style functions can be chained
  - Test that styles render correctly with sample text

**Testing Guide Reference**: [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md) - Section: Performance Benchmarks (Benchmark 5: TUI Rendering)

**Acceptance Criteria Document**: None (use [`ENHANCED-TEST-CRITERIA-TEMPLATE.md`](../../milestones/ENHANCED-TEST-CRITERIA-TEMPLATE.md))

---

### Task 2: Create Status Bar Component

**Dependencies**: Task 1 (Theme System)  
**Files**: [`ui/statusbar/model.go`](../../project-structure.md:247), [`ui/statusbar/model_test.go`](../../go-testing-guide.md:682)  
**Integration Points**: Theme system from Task 1  
**Estimated Complexity**: Low

**Description**: Implement a basic status bar component that displays application status information. This component will be integrated into the root app model in Task 3.

**Acceptance Criteria**:
- [x] Component MUST implement Bubble Tea Model interface (Init, Update, View)
- [x] Model MUST track character count and line count
- [x] View() MUST render styled status bar using theme helpers
- [x] Update() MUST handle tea.WindowSizeMsg to adjust width
- [x] Component MUST have package-level documentation comment
- [x] Exported functions MUST have documentation comments

**Testing Requirements**:
- Coverage Target: > 85%
- Critical Test Scenarios:
  - Test Init() returns nil command
  - Test Update() handles WindowSizeMsg correctly
  - Test View() renders correct content
  - Test View() uses theme styles correctly
- Edge Cases:
  - Test with zero width/height
  - Test with very large width/height
  - Test rapid window resize events

**Testing Guide Reference**: [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md) - Section: Integration Tests (Test 5: TUI State Management)

**Acceptance Criteria Document**: None (use [`ENHANCED-TEST-CRITERIA-TEMPLATE.md`](../../milestones/ENHANCED-TEST-CRITERIA-TEMPLATE.md))

---

### Task 3: Create Root App Model

**Dependencies**: Task 1 (Theme System), Task 2 (Status Bar)  
**Files**: [`ui/app/model.go`](../../project-structure.md:227), [`ui/app/model_test.go`](../../go-testing-guide.md:682)  
**Integration Points**: Status bar component, config system (M1), logging system (M1)  
**Estimated Complexity**: Medium

**Description**: Implement the root Bubble Tea model that serves as the main application entry point. This model manages the overall TUI state and coordinates child components.

**Acceptance Criteria**:
- [x] Model MUST implement Bubble Tea Model interface (Init, Update, View)
- [x] Init() MUST return nil command
- [x] Update() MUST handle tea.KeyMsg for character input
- [x] Update() MUST handle tea.KeyMsg for quit (Ctrl+C, 'q')
- [x] Update() MUST handle tea.WindowSizeMsg
- [x] Update() MUST return tea.Quit command on quit
- [x] View() MUST render status bar at bottom of screen
- [x] View() MUST use theme styles for consistent appearance
- [x] Model MUST have package-level documentation comment
- [x] Exported functions MUST have documentation comments

**Testing Requirements**:
- Coverage Target: > 85%
- Critical Test Scenarios:
  - Test Init() returns nil command
  - Test Update() handles character input correctly
  - Test Update() handles Ctrl+C quit
  - Test Update() handles 'q' quit
  - Test Update() handles WindowSizeMsg
  - Test View() renders status bar
  - Test View() uses theme styles
- Edge Cases:
  - Test with rapid keyboard input
  - Test with multiple quit attempts
  - Test with window resize during input
  - Test with special characters and unicode

**Testing Guide Reference**: [`go-testing-guide.md`](../../go-testing-guide.md) - Section: Bubble Tea Testing Patterns (Model Testing, User Input Simulation)

**Acceptance Criteria Document**: None (use [`ENHANCED-TEST-CRITERIA-TEMPLATE.md`](../../milestones/ENHANCED-TEST-CRITERIA-TEMPLATE.md))

---

### Task 4: Integrate TUI with Main Application

**Dependencies**: Task 3 (Root App Model)  
**Files**: [`cmd/promptstack/main.go`](../../project-structure.md:227)  
**Integration Points**: Root app model, config system (M1), logging system (M1)  
**Estimated Complexity**: Medium

**Description**: Integrate the TUI model with the main application entry point. This involves creating the Bubble Tea program and wiring up dependencies.

**Acceptance Criteria**:
- [x] main() MUST create app model with dependencies
- [x] main() MUST create tea.Program with app model
- [x] main() MUST pass tea.WithAltScreen() option
- [x] main() MUST handle program errors with logging
- [x] main() MUST log TUI startup and shutdown
- [x] Application MUST launch TUI without errors
- [x] Application MUST exit cleanly on quit

**Testing Requirements**:
- Coverage Target: > 80%
- Critical Test Scenarios:
  - Test TUI launches successfully
  - Test TUI handles keyboard input
  - Test TUI quits on Ctrl+C
  - Test TUI quits on 'q'
  - Test TUI exits cleanly
- Edge Cases:
  - Test with invalid config
  - Test with missing config
  - Test with logging errors
  - Test with terminal size changes

**Testing Guide Reference**: [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md) - Section: End-to-End Scenarios (Scenario 1: First-Time User Setup)

**Acceptance Criteria Document**: None (use [`ENHANCED-TEST-CRITERIA-TEMPLATE.md`](../../milestones/ENHANCED-TEST-CRITERIA-TEMPLATE.md))

---

## Integration Contracts

### M1 Dependencies (from Milestone 1: Bootstrap & Config)

**Config System** ([`internal/config/config.go`](../../project-structure.md:183)):
```go
// Config represents application configuration
type Config struct {
    Path    string
    Version string
    // ... other fields
}

// Load loads configuration from default location
func Load() (*Config, error)
```

**Logging System** ([`internal/platform/logging/logger.go`](../../project-structure.md:197)):
```go
// Logger provides structured logging
type Logger interface {
    Info(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Debug(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Sync() error
}

// New creates a new logger instance
func New() (*Logger, error)
```

### Component Integration Contracts

**Task 3 → Task 4 Integration**:
```go
// ui/app/model.go - Expected interface
type Model struct {
    statusBar statusbar.Model
    width     int
    height    int
    quitting  bool
}

func New() Model
func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m Model) View() string
```

**Task 2 → Task 3 Integration**:
```go
// ui/statusbar/model.go - Expected interface
type Model struct {
    charCount int
    lineCount int
    width     int
}

func New() Model
func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m Model) View() string
func (m *Model) SetCharCount(count int)
func (m *Model) SetLineCount(count int)
```

**Task 1 → Task 2 & Task 3 Integration**:
```go
// ui/theme/theme.go - Expected interface
const (
    BackgroundPrimary   = "#1e1e2e"
    BackgroundSecondary = "#181825"
    ForegroundPrimary = "#cdd6f4"
    ForegroundMuted    = "#a6adc8"
    AccentBlue   = "#89b4fa"
    AccentGreen  = "#a6e3a1"
    AccentYellow = "#f9e2af"
    AccentRed    = "#f38ba8"
    BorderPrimary = "#45475a"
)

func StatusStyle() lipgloss.Style
func ModalStyle() lipgloss.Style
func ActivePlaceholderStyle() lipgloss.Style
```

## Task Dependencies

```
Task 1 (Theme System)
    ↓
Task 2 (Status Bar) ─────┐
    ↓                      │
Task 3 (Root App Model) ──┤
    ↓                      │
Task 4 (Main Integration) ┘
```

---

## Testing Strategy

### Unit Tests
- Each component tested independently
- Table-driven tests for multiple scenarios
- Mock dependencies where needed
- Focus on effects, not implementation

### Integration Tests
- Test component interactions
- Test message flow between components
- Test state management across components

### Manual Testing
- Launch TUI and verify rendering
- Test keyboard input handling
- Test quit functionality (Ctrl+C, 'q')
- Test window resize handling

---

## Success Criteria

### Functional Requirements
- [ ] TUI launches without errors
- [ ] Status bar displays at bottom of screen
- [ ] Keyboard input is captured
- [ ] Quit works with Ctrl+C
- [ ] Quit works with 'q' key
- [ ] Clean exit without errors

### Integration Requirements
- [ ] Theme system integrates with all components
- [ ] Status bar integrates with root app model
- [ ] Root app model integrates with main application
- [ ] Config system integrates with TUI (for future use)
- [ ] Logging system integrates with TUI

### Edge Cases & Error Handling
- [ ] Handle zero terminal size
- [ ] Handle very large terminal size
- [ ] Handle rapid keyboard input
- [ ] Handle window resize during input
- [ ] Handle special characters and unicode

### Performance Requirements
- [ ] TUI startup time < 100ms
- [ ] Input response time < 16ms (60 FPS)
- [ ] Render time < 16ms (60 FPS)
- [ ] No memory leaks during operation

### User Experience Requirements
- [ ] Status bar is clearly visible
- [ ] Colors are consistent (Catppuccin Mocha)
- [ ] Quit is intuitive (Ctrl+C, 'q')
- [ ] No visual glitches during resize
- [ ] Smooth keyboard input handling

---

## References

### Style Guide References
- [`go-style-guide.md`](../../go-style-guide.md) - Section: Type Design (lines 52-95)
- [`go-style-guide.md`](../../go-style-guide.md) - Section: Error Handling (lines 113-201)
- [`go-style-guide.md`](../../go-style-guide.md) - Section: Project-Specific Rules (lines 1062-1116)

### Testing Guide References
- [`go-testing-guide.md`](../../go-testing-guide.md) - Section: Bubble Tea Testing Patterns (lines 47-143)
- [`go-testing-guide.md`](../../go-testing-guide.md) - Section: User Input Simulation (lines 147-214)
- [`go-testing-guide.md`](../../go-testing-guide.md) - Section: Test Organization (lines 678-723)

### Key Learnings References
- [`learnings/ui-domain.md`](../learnings/ui-domain.md) - Category 1: Bubble Tea Model Implementation
- [`learnings/ui-domain.md`](../learnings/ui-domain.md) - Category 6: Centralized Theme System
- [`learnings/go-fundamentals.md`](../learnings/go-fundamentals.md) - Category 7: Error Handling Patterns

### Project Structure References
- [`project-structure.md`](../../project-structure.md) - UI Domain (lines 227-315)
- [`project-structure.md`](../../project-structure.md) - ui/app/ (lines 227-246)
- [`project-structure.md`](../../project-structure.md) - ui/statusbar/ (lines 247-266)
- [`project-structure.md`](../../project-structure.md) - ui/theme/ (lines 267-286)

---

**Last Updated**: 2026-01-08
**Milestone Group**: Foundation (M1-M6)
**Testing Guide**: [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md)

---

## Task Completion Summary

### Task 4: Integrate TUI with Main Application ✅ COMPLETED

**Files Created/Modified:**
- [`cmd/promptstack/main.go`](cmd/promptstack/main.go) - Updated to integrate TUI with main application
- [`cmd/promptstack/main_test.go`](cmd/promptstack/main_test.go) - Comprehensive test suite with 13 test functions

**Implementation Details:**
- Integrated Bubble Tea program creation with [`tea.NewProgram()`](cmd/promptstack/main.go:43)
- Added [`tea.WithAltScreen()`](cmd/promptstack/main.go:43) option for full-screen TUI experience
- Created app model using [`app.New()`](cmd/promptstack/main.go:38)
- Added TUI startup logging at line 40
- Added TUI shutdown logging at line 49
- Implemented error handling with logging for program errors (lines 44-47)
- Maintained existing bootstrap and logging initialization from Milestone 1

**Test Results:**
All 13 test suites passed successfully:
- TestMainFunction
- TestRunFunction (2 subtests)
- TestTUILaunch
- TestTUIWithProgram
- TestTUIHandlesKeyboardInput
- TestTUIHandlesQuit (3 subtests)
- TestTUIHandlesWindowSize
- TestTUILogging
- TestTUIErrorHandling
- TestTUIWithAltScreen
- TestTUIIntegration
- TestTUIPerformance
- TestTUIEdgeCases (4 subtests)
- TestTUIConcurrentAccess

**Milestone 2 Status:**
- Task 1: Create Theme System - ✅ Completed
- Task 2: Create Status Bar Component - ✅ Completed
- Task 3: Create Root App Model - ✅ Completed
- Task 4: Integrate TUI with Main Application - ✅ Completed

**Milestone 2 is now complete!** 🎉