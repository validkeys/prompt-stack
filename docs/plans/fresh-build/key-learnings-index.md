# Key Learnings Index

**Purpose**: Index of domain-specific key learnings from previous PromptStack implementation.

**Status**: ✅ Reorganized - Learnings now organized by domain in [`learnings/`](learnings/) directory

**Last Updated**: 2026-01-07

---

## Overview

The original [`key-learnings.md`](key-learnings.md) file (3,118 lines) has been reorganized into domain-specific documents for easier reference during milestone planning and execution.

**New Organization**: [`learnings/`](learnings/) directory with 8 domain-specific documents plus index

---

## Domain-Specific Learning Documents

### Go Fundamentals
**File**: [`learnings/go-fundamentals.md`](learnings/go-fundamentals.md)

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

**Priority**: High - These are foundational learnings applicable to all milestones

---

### Editor Domain
**File**: [`learnings/editor-domain.md`](learnings/editor-domain.md)

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

**Priority**: High - Directly applicable to editor milestones

---

### UI/TUI Domain
**File**: [`learnings/ui-domain.md`](learnings/ui-domain.md)

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

**Priority**: High - Directly applicable to UI milestones

---

### Error Handling
**File**: [`learnings/error-handling.md`](learnings/error-handling.md)

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

**Priority**: High - Applicable to all milestones

---

### AI Domain
**File**: [`learnings/ai-domain.md`](learnings/ai-domain.md)

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

**Priority**: High - Directly applicable to AI milestones

---

### Vim Domain
**File**: [`learnings/vim-domain.md`](learnings/vim-domain.md)

**Description**: Vim mode patterns

**Related Milestones**: M34, M35

**Key Topics**:
- Vim Mode Transition Logic
- Read-Only Mode During Async Operations

**Priority**: Medium - Applicable to vim milestones

---

### History Domain
**File**: [`learnings/history-domain.md`](learnings/history-domain.md)

**Description**: History management patterns

**Related Milestones**: M15, M16, M17

**Key Topics**:
- History Browser Integration
- History Manager Initialization in Bootstrap

**Priority**: Medium - Applicable to history milestones

---

### Architecture Patterns
**File**: [`learnings/architecture-patterns.md`](learnings/architecture-patterns.md)

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

**Priority**: Medium - Applicable to multiple milestones

---

## Learnings Directory Index

**File**: [`learnings/README.md`](learnings/README.md)

**Description**: Index and navigation guide for all learning documents

**Contents**:
- Overview of all learning documents
- How to use these documents during milestone planning and implementation
- Document template reference
- Maintenance guidelines
- Related documentation links

---

## Original Document

**File**: [`key-learnings.md`](key-learnings.md)

**Status**: Archived - Content has been reorganized into domain-specific documents

**Description**: Original comprehensive learnings document (3,118 lines) organized by topic

**Note**: This document is preserved for reference but should not be used for new implementations. Use the domain-specific documents in the [`learnings/`](learnings/) directory instead.

---

## Organization Plan

**File**: [`key-learnings-organization-plan.md`](key-learnings-organization-plan.md)

**Status**: In Progress - Phase 1 Complete

**Description**: Detailed plan for reorganizing key learnings into domain-specific documents

**Completed Phases**:
- ✅ Phase 1: Create Learning Documents (All 9 documents created)

**Remaining Phases**:
- ⏳ Phase 2: Update DOCUMENT-REFERENCE-MATRIX.md
- ⏳ Phase 3: Update milestone-execution-prompt.md
- ⏳ Phase 4: Update DOCUMENT-INDEX.md
- ⏳ Phase 5: Archive Original Document

---

## Quick Reference by Milestone

### Foundation (M1-M6)
- [`learnings/go-fundamentals.md`](learnings/go-fundamentals.md) - Go-specific patterns
- [`learnings/ui-domain.md`](learnings/ui-domain.md) - Bubble Tea patterns
- [`learnings/error-handling.md`](learnings/error-handling.md) - Error handling
- [`learnings/architecture-patterns.md`](learnings/architecture-patterns.md) - Architecture patterns

### Editor (M4, M5, M6, M11, M12, M13, M14)
- [`learnings/editor-domain.md`](learnings/editor-domain.md) - Editor implementation
- [`learnings/ui-domain.md`](learnings/ui-domain.md) - UI patterns
- [`learnings/error-handling.md`](learnings/error-handling.md) - Error handling

### History (M15, M16, M17)
- [`learnings/go-fundamentals.md`](learnings/go-fundamentals.md) - SQLite driver
- [`learnings/history-domain.md`](learnings/history-domain.md) - History patterns
- [`learnings/error-handling.md`](learnings/error-handling.md) - Error handling
- [`learnings/architecture-patterns.md`](learnings/architecture-patterns.md) - Database schema

### AI (M27-M33)
- [`learnings/ai-domain.md`](learnings/ai-domain.md) - AI integration
- [`learnings/ui-domain.md`](learnings/ui-domain.md) - Diff viewer
- [`learnings/error-handling.md`](learnings/error-handling.md) - Error handling

### Vim (M34, M35)
- [`learnings/vim-domain.md`](learnings/vim-domain.md) - Vim mode
- [`learnings/ui-domain.md`](learnings/ui-domain.md) - Keyboard navigation

### All Milestones
- [`learnings/error-handling.md`](learnings/error-handling.md) - Error handling (applicable to all)
- [`learnings/architecture-patterns.md`](learnings/architecture-patterns.md) - Architecture patterns

---

## Related Documentation

- [`key-learnings-organization-plan.md`](key-learnings-organization-plan.md) - Organization plan and rationale
- [`DOCUMENT-REFERENCE-MATRIX.md`](DOCUMENT-REFERENCE-MATRIX.md) - Document reference matrix (to be updated)
- [`milestone-execution-prompt.md`](milestone-execution-prompt.md) - Execution prompt (to be updated)
- [`DOCUMENT-INDEX.md`](DOCUMENT-INDEX.md) - Document index (to be updated)
- [`project-structure.md`](project-structure.md) - Project structure by domain
- [`go-style-guide.md`](go-style-guide.md) - Go coding standards
- [`go-testing-guide.md`](go-testing-guide.md) - Go testing patterns

---

**Total Learning Documents**: 8 domain-specific + 1 index = 9 documents
**Total Learnings Preserved**: 60+ learning sections from original document
**Last Updated**: 2026-01-07
**Status**: ✅ Phase 1 Complete - All learning documents created and indexed