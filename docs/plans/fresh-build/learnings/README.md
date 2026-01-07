# Key Learnings Index

**Purpose**: Domain-specific key learnings from previous PromptStack implementation, organized for easy reference during milestone planning and execution.

**Last Updated**: 2026-01-07

---

## Overview

This directory contains key learnings and implementation patterns extracted from the previous PromptStack implementation. Each document focuses on a specific domain, making it easy to find relevant learnings when planning or implementing milestones.

**Source Document**: [`../../../archive/key-learnings.md`](../../../archive/key-learnings.md) (3,118 lines of learnings) - **Archived**

---

## Learning Documents

### Go Fundamentals
**File**: [`go-fundamentals.md`](go-fundamentals.md)

**Description**: Go-specific patterns and pitfalls

**Related Milestones**: M1-M6 (Foundation), M15 (SQLite Setup)

**Key Topics**:
- Go Embed Limitations
- Zap Logger Structured Fields
- Regex Matching in Go
- SQLite Driver Selection
- Go Version Requirements
- Project Structure Organization
- Error Handling Patterns
- Frontmatter Parsing Strategy

---

### Editor Domain
**File**: [`editor-domain.md`](editor-domain.md)

**Description**: Editor implementation patterns

**Related Milestones**: M4, M5, M6, M11, M12, M13, M14

**Key Topics**:
- Placeholder Parsing
- Index Scoring Algorithm
- Validation Strategy
- Placeholder System Implementation
- Placeholder Validation Strategy
- Cursor Position Management for Placeholders
- Placeholder Highlighting in TUI
- Theme Integration for Placeholders
- Text Placeholder Editing Mode
- Text Editor Cursor Positioning
- File Path Management for History
- Lipgloss Styling

---

### UI/TUI Domain
**File**: [`ui-domain.md`](ui-domain.md)

**Description**: UI/TUI implementation patterns

**Related Milestones**: M2, M8, M19, M24, M25, M26, M31, M36, M37

**Key Topics**:
- Bubble Tea Model Implementation
- Cursor and Viewport Management
- Auto-save Debouncing with Bubble Tea
- Custom Message Types
- Status Bar State Management
- Centralized Theme System
- Library Browser Implementation
- Modal Overlay Pattern
- Fuzzy Matching Integration
- Split-Pane Layout with Lipgloss
- Keyboard Navigation with Vim Mode
- Message-Based Command Execution
- Command Palette Implementation
- Command Categorization
- Placeholder Command Handlers

---

### Error Handling
**File**: [`error-handling.md`](error-handling.md)

**Description**: Error handling patterns

**Related Milestones**: M1-M38 (All milestones)

**Key Topics**:
- Error Handling Architecture
- Status Bar Component Design
- Modal Component Pattern
- Error Handler Integration
- Import Cycle Prevention
- Error Recovery Strategies
- Graceful File Read Error Handling
- Error Logging Integration

---

### AI Domain
**File**: [`ai-domain.md`](ai-domain.md)

**Description**: AI integration patterns

**Related Milestones**: M27, M28, M29, M30, M31, M32, M33

**Key Topics**:
- AI Applying Indicator and Read-Only Mode
- Diff Viewer Modal Implementation
- AI Message-Based Workflow
- Token Budget Enforcement
- Context Selection Algorithm
- Command Palette Integration for AI Features
- Read-Only Mode During Async Operations

---

### Vim Domain
**File**: [`vim-domain.md`](vim-domain.md)

**Description**: Vim mode patterns

**Related Milestones**: M34, M35

**Key Topics**:
- Vim Mode Transition Logic
- Read-Only Mode During Async Operations

---

### History Domain
**File**: [`history-domain.md`](history-domain.md)

**Description**: History management patterns

**Related Milestones**: M15, M16, M17

**Key Topics**:
- History Browser Integration
- History Manager Initialization in Bootstrap

---

### Architecture Patterns
**File**: [`architecture-patterns.md`](architecture-patterns.md)

**Description**: Architecture and design patterns

**Related Milestones**: M1-M38 (All milestones)

**Key Topics**:
- Database Schema Design
- Configuration Management
- Starter Prompt Extraction
- Logging Strategy
- Command Registry Pattern
- Confirmation Dialog Integration for Destructive Operations
- Type Assertion for Bubble Tea Model Updates

---

## How to Use These Documents

### During Milestone Planning

1. **Identify the milestone domain** (e.g., M4 is editor domain)
2. **Read the relevant learning document** (e.g., [`editor-domain.md`](editor-domain.md))
3. **Review related learnings** for the milestone
4. **Apply learnings to implementation plan**
5. **Note any deviations** with justification

### During Implementation

1. **Reference the learning document** when encountering challenges
2. **Apply the documented patterns** to solve problems
3. **Follow the implementation examples** provided
4. **Avoid repeating past mistakes** documented in learnings

### Integration with Documentation Workflow

These learning documents are integrated into the documentation workflow:

- **DOCUMENT-REFERENCE-MATRIX.md**: Lists learning documents as context-specific references
- **milestone-execution-prompt.md**: Includes key learnings in document checking workflow
- **Implementation plans**: Should reference relevant learnings

---

## Document Template

Each learning document follows this structure:

```markdown
# {Domain} Key Learnings

**Purpose**: Key learnings and implementation patterns for {domain} from previous PromptStack implementation.

**Related Milestones**: M{N}, M{N+1}, M{N+2}

**Related Documents**: 
- [`project-structure.md`](../project-structure.md) - {domain} section
- [`go-style-guide.md`](../go-style-guide.md) - Relevant patterns
- [`go-testing-guide.md`](../go-testing-guide.md) - Testing patterns

---

## Learning Categories

### Category 1: {Category Name}

**Learning**: {Brief description}

**Problem**: {What problem was encountered}

**Solution**: {How it was solved}

**Implementation Pattern**:
```go
// Code example showing the pattern
```

**Lesson**: {Key takeaway}

**Related Milestones**: M{N}, M{N+1}

**When to Apply**: {Specific scenarios where this applies}

---

## Quick Reference

| Learning | Milestone | Priority |
|----------|-----------|----------|
| {Learning 1} | M{N} | High |
| {Learning 2} | M{N+1} | Medium |

---

**Last Updated**: {Date}
**Source**: [`../../../archive/key-learnings.md`](../../../archive/key-learnings.md) (Archived)
```

---

## Maintenance

### Adding New Learnings

When new learnings are discovered during implementation:

1. **Identify the domain** of the learning
2. **Add to the appropriate document** in this directory
3. **Update the Quick Reference table** at the end of the document
4. **Update this README** if adding a new category or document
5. **Consider updating the archived** [`key-learnings.md`](../../../archive/key-learnings.md) for completeness

### Updating Existing Learnings

When existing learnings need refinement:

1. **Update the learning document** with new information
2. **Update the Quick Reference table** if priorities change
3. **Document the reason** for the update in the learning section

---

## Related Documentation

- [`../../../archive/key-learnings.md`](../../../archive/key-learnings.md) - Original comprehensive learnings document (Archived)
- [`key-learnings-organization-plan.md`](../key-learnings-organization-plan.md) - Organization plan and rationale
- [`DOCUMENT-REFERENCE-MATRIX.md`](../DOCUMENT-REFERENCE-MATRIX.md) - Document reference matrix with learning documents
- [`milestone-execution-prompt.md`](../milestone-execution-prompt.md) - Execution prompt with learning integration
- [`project-structure.md`](../project-structure.md) - Project structure by domain
- [`go-style-guide.md`](../go-style-guide.md) - Go coding standards
- [`go-testing-guide.md`](../go-testing-guide.md) - Go testing patterns

---

**Last Updated**: 2026-01-07
**Status**: ✅ All Phases Complete - Key learnings organization fully implemented
**Next Steps**: Reference domain-specific learning documents during milestone planning and execution