# Milestone 2: Basic TUI Shell - Reference Document

**Milestone Number**: 2  
**Title**: Basic TUI Shell  
**Goal**: Render functional TUI with quit handling

---

## How to Use This Document

**Read this section when:**
- [Before implementing Task 1] - Understanding theme system requirements and Catppuccin Mocha palette
- [Before implementing Task 2] - Understanding status bar component patterns and Bubble Tea model implementation
- [Before implementing Task 3] - Understanding root app model structure and message handling
- [Before implementing Task 4] - Understanding main application integration and dependency wiring
- [When debugging TUI issues] - Referencing Bubble Tea patterns and common anti-patterns
- [When writing tests] - Referencing testing patterns and table-driven test structure
- [When encountering build errors] - Checking code examples for proper imports and syntax

**Key sections:**
- Lines 33-88: Architecture Context - Read before Task 1 to understand UI domain structure and package organization
- Lines 90-100: Document Navigation - Understand how this document is split across multiple parts
- See [`reference-part-2.md`](reference-part-2.md) for Style Guide References (lines 28-150) and Testing Guide References (lines 153-300)
- See [`reference-part-3.md`](reference-part-3.md) for Key Learnings References (lines 28-159) and Implementation Notes for Tasks 1-2 (lines 161-500)
- See [`reference-part-4.md`](reference-part-4.md) for Implementation Notes for Tasks 3-4 (lines 31-411), Common Patterns (lines 414-528), Performance Considerations (lines 532-560), and Testing Checklist (lines 563-584)

**Related documents:**
- See [`task-list.md`](task-list.md) for complete task breakdown, acceptance criteria, and integration contracts
- Cross-reference with [`go-style-guide.md`](../../go-style-guide.md) for coding standards and patterns
- Cross-reference with [`go-testing-guide.md`](../../go-testing-guide.md) for testing patterns and TDD approach
- Cross-reference with [`project-structure.md`](../../project-structure.md:227) for UI domain package structure
- Cross-reference with [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md) for integration test scenarios

---

## Architecture Context

### Domain Overview

Milestone 2 establishes the UI/TUI domain for PromptStack. This domain is responsible for all terminal user interface components and follows the Bubble Tea model-view-update pattern.

**UI Domain Structure** (from [`project-structure.md`](../../project-structure.md:227)):
```
ui/
├── app/           # Root application model
│   └── model.go   # Main Bubble Tea model
├── statusbar/     # Status bar component
│   └── model.go   # Status bar model
└── theme/         # Theme system
    └── theme.go   # Color constants and style helpers
```

### Package Structure

**ui/theme/** - Centralized theme system
- Provides Catppuccin Mocha color palette
- Exports style helper functions for consistent styling
- Foundation for all UI components

**ui/statusbar/** - Status bar component
- Displays application status information
- Implements Bubble Tea Model interface
- Integrates with root app model

**ui/app/** - Root application model
- Main Bubble Tea model for TUI
- Coordinates child components
- Handles keyboard input and quit logic

### Dependencies

**External Dependencies** (from [`DEPENDENCIES.md`](../../DEPENDENCIES.md)):
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling library

**Internal Dependencies**:
- Config system (M1) - For theme preferences (future use)
- Logging system (M1) - For TUI debugging

### Dependency Flow

```
cmd/promptstack/main.go
    ↓
ui/app/model.go
    ↓
ui/statusbar/model.go
    ↓
ui/theme/theme.go
```

---

## Document Navigation

This reference document is split into multiple parts to comply with the 600-line limit:

- **Part 1 (this file)**: Architecture Context, Navigation Guide
- **Part 2**: Style Guide References, Testing Guide References
- **Part 3**: Key Learnings References, Implementation Notes (Tasks 1-2)
- **Part 4**: Implementation Notes (Tasks 3-4), Common Patterns, Performance, Testing Checklist

**Continue reading**: See [`reference-part-2.md`](reference-part-2.md) for Style Guide and Testing Guide references.

---

**Last Updated**: 2026-01-07  
**Milestone Group**: Foundation (M1-M6)  
**Testing Guide**: [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md)