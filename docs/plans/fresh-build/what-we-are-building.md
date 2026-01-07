# What We're Building: PromptStack

## Elevator Pitch

**PromptStack** is a sophisticated CLI tool that revolutionizes AI prompt engineering by providing a modern, terminal-based workspace for composing complex prompts. It combines a reusable prompt library, intelligent AI-powered suggestions, and developer-friendly features like vim keybindings and auto-save to help developers build better prompts faster. Think of it as an IDE for AI prompts—complete with templates, placeholders, file references, and Claude integration—all in a beautiful, minimal TUI interface.

## Key Features

### Core Composition
- Split-pane workspace with live editor and AI suggestions panel
- Auto-save with debouncing and undo/redo support (100 levels)
- Optional vim keybindings (Normal/Insert/Visual modes)
- Placeholder system for template variables (`{{text:name}}` and `{{list:name}}`)

### Prompt Library
- Global library of reusable prompt templates (workflows, commands, decorations, rules)
- Fuzzy search and browser with color-coded categories
- Prompt creation and editing with markdown preview
- YAML frontmatter for metadata (title, description, tags)

### AI Integration
- Claude API integration for intelligent suggestions
- Context-aware prompt recommendations
- Diff-based suggestion application with preview
- Conservative token usage with smart context selection

### Developer Tools
- Command palette (Space key) for quick actions
- File reference system with fuzzy finder and gitignore support
- History browser with full-text search (SQLite + FTS5)
- Library validation with error/warning reporting

### Polish
- Responsive layout (adapts to terminal size)
- Settings panel for configuration
- Comprehensive error handling and logging
- Modern, minimal design with Catppuccin Mocha theme

### Technical Excellence
- Built with Go and Bubble Tea
- Test-driven development with 80%+ coverage target
- 38 granular milestones with comprehensive testing guides
- Domain-driven architecture with 8 core domains

---

## Development Philosophy

The documentation shows an incredibly well-planned system with meticulous attention to quality, testing, and developer experience. It's essentially a professional-grade development system for building a professional-grade tool.

**Key Principles:**
- Quality over speed
- Every line of code tested
- Every milestone verified
- Every checkpoint documented

## Quick Reference

- **Primary Documentation**: [`README.md`](README.md)
- **Requirements**: [`requirements.md`](requirements.md)
- **Architecture**: [`project-structure.md`](project-structure.md)
- **Milestones**: [`milestones.md`](milestones.md)
- **How to Use**: [`HOW-TO-USE.md`](HOW-TO-USE.md)
- **Document Index**: [`DOCUMENT-INDEX.md`](DOCUMENT-INDEX.md)

---

**Last Updated**: 2026-01-07