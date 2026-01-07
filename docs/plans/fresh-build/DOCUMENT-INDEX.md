# Document Index

**Purpose**: Quick reference index for all PromptStack documentation. Use this to quickly locate information across all documents.

---

## Quick Navigation

### By Document Type
- [Planning Documents](#planning-documents)
- [Implementation Guides](#implementation-guides)
- [Technical Specifications](#technical-specifications)
- [Supporting Documents](#supporting-documents)

### By Milestone Phase
- [Foundation (M1-M6)](#foundation-milestones-m1-m6)
- [Library Integration (M7-M10)](#library-integration-milestones-m7-m10)
- [Placeholders (M11-M14)](#placeholder-milestones-m11-m14)
- [History (M15-M17)](#history-milestones-m15-m17)
- [Commands & Files (M18-M22)](#commands--files-milestones-m18-m22)
- [Prompt Management (M23-M26)](#prompt-management-milestones-m23-m26)
- [AI Integration (M27-M33)](#ai-integration-milestones-m27-m33)
- [Vim Mode (M34-M35)](#vim-mode-milestones-m34-m35)
- [Polish (M36-M38)](#polish-milestones-m36-m38)

### By Domain
- [Editor Domain](#editor-domain)
- [Library Domain](#library-domain)
- [History Domain](#history-domain)
- [AI Domain](#ai-domain)
- [UI Domain](#ui-domain)
- [Platform Domain](#platform-domain)
- [Vim Domain](#vim-domain)
- [Commands Domain](#commands-domain)

---

## Planning Documents

### milestone-execution-prompt.md
**Purpose**: Main execution workflow for all milestones
**Key Sections**:
- Document checking workflow (CRITICAL)
- Step-by-step execution process
- Task list format
- Reference document format
- Checkpoint documentation format
- Testing guide format
- Summary format

**When to Use**: Always - before starting any milestone

### milestones.md
**Purpose**: Complete list of all 38 milestones with goals, deliverables, and test criteria
**Key Sections**:
- Milestone 1: Bootstrap & Config
- Milestone 2: Basic TUI Shell
- Milestone 3: File I/O Foundation
- ... (all 38 milestones)

**When to Use**: Always - to understand specific milestone requirements

### requirements.md
**Purpose**: Complete feature requirements for PromptStack
**Key Sections**:
- Core Concepts (Library, Composition, Placeholders)
- Features 1-11 (detailed feature specifications)
- Configuration (global config, initialization)
- User Interface (vim support, hotkeys, visual design)
- Technical Requirements (technology stack, file operations, search, AI context)
- Error Handling (by error type)
- Performance & Limits (file sizes, library scale, memory)

**When to Use**: Always - to understand feature requirements

### project-structure.md
**Purpose**: Domain-driven architecture and package structure
**Key Sections**:
- 8 Core Domains (Editor, Prompt, Library, History, AI, Config, UI, Platform)
- Recommended project structure (complete directory tree)
- Design principles (domain separation, dependency direction, package naming)
- File organization patterns
- Testing structure
- Interface design patterns
- Error handling patterns
- Concurrency patterns
- Configuration injection
- Theme system (Catppuccin Mocha palette)

**When to Use**: Always - to understand architecture and file organization

---

## Implementation Guides

### go-style-guide.md
**Purpose**: Go coding standards and best practices
**Key Sections**:
- Core principles (clarity, explicit, simple, testable)
- Package organization (naming, structure, comments)
- Type design (constructors, structs, method receivers)
- Error handling (creation, checking, custom errors)
- Interfaces (definition location, size)
- Dependency management (injection pattern, dependency direction)
- Concurrency (Bubble Tea pattern, synchronous domain logic)
- Testing (package names, table-driven tests, mocking)
- Code organization (file naming, function length, comments)
- Common patterns (options pattern, context usage)
- Anti-patterns to avoid
- Project-specific rules (UI components, theme usage, logging)

**When to Use**: Always - for all Go code

### go-testing-guide.md
**Purpose**: Testing patterns and TDD best practices
**Key Sections**:
- Core testing philosophy (effect-first, test pyramid)
- Bubble Tea testing patterns (model, command, message testing)
- User input simulation (key messages, sequences, complex interactions)
- Domain testing (editor, prompt, AI)
- Integration testing (component, E2E workflows)
- Test utilities (helpers, mock factories, fixtures)
- Testing anti-patterns (don't test view output, don't test private methods)
- Performance testing (benchmarks, memory testing)
- Advanced patterns (parallel tests, cleanup handlers, context/timeouts)
- Test organization (file structure, naming)
- Running tests (unit, integration, benchmarks)
- Quick reference (key patterns, common assertions, test helpers)
- Best practices summary

**When to Use**: Always - for all testing

### ENHANCED-TEST-CRITERIA-TEMPLATE.md
**Purpose**: Template for writing comprehensive acceptance criteria
**Key Sections**:
- 5-category testing framework (Functional, Integration, Edge Cases, Performance, UX)
- Functional Requirements (FR) format
- Integration Requirements (IR) format
- Edge Cases & Error Handling (EC) format
- Performance Requirements (PR) format
- User Experience Requirements (UX) format
- Success Criteria (Must Have, Should Have, Nice to Have)
- Test case examples with Go code
- Best practices for writing acceptance criteria

**When to Use**: When creating acceptance criteria for complex milestones

---

## Technical Specifications

### CONFIG-SCHEMA.md
**Purpose**: Complete configuration schema and validation rules
**Key Sections**:
- Config file location (~/.promptstack/config.yaml)
- Complete schema (all configuration options with defaults)
- Validation rules (required fields, field validation, error messages)
- Environment variable overrides (supported variables, precedence)
- Configuration migration (version tracking, migration process, examples)
- Configuration loading (load order, error handling, setup wizard)
- Configuration examples (minimal, development, production, custom paths)
- Best practices (security, performance, usability, debugging)
- Troubleshooting (common issues, resetting configuration)
- Future enhancements

**When to Use**: M1 (Bootstrap & Config), M36 (Settings Panel)

### DATABASE-SCHEMA.md
**Purpose**: SQLite database schema for history management
**Key Sections**:
- Database location (~/.promptstack/data/history.db)
- Tables (compositions, composition_tags)
- Full-Text Search (FTS5 virtual table)
- Triggers (auto-update timestamps, FTS5 sync)
- Query patterns (CRUD, listing, search, tag operations, statistics)
- Migration strategy (version management, migration process, rollback)
- Performance considerations (indexing, query optimization, FTS5)
- Data integrity (foreign keys, transactions, validation)
- Backup and recovery (backup strategy, recovery procedures)
- Maintenance (vacuum, analyze, cleanup)
- Testing (test data, test queries)
- Security considerations (file permissions, SQL injection, data privacy)
- Integration with markdown files (dual storage, sync verification, rebuild)

**When to Use**: M15 (SQLite Setup), M16 (History Sync), M17 (History Browser)

### DEPENDENCIES.md
**Purpose**: Complete list of all external dependencies
**Key Sections**:
- Core dependencies (TUI framework, database & storage, AI integration, fuzzy search, file system, text processing, terminal utilities, logging, UUID generation, clipboard, standard library extensions)
- Development dependencies (testing framework, mocking)
- Dependency categories (by layer, by license)
- License information (summary, requirements, compliance checklist)
- Adding dependencies (prerequisites, procedure, version pinning)
- Updating dependencies (update policy, procedure, rollback)
- Security considerations (scanning, vulnerability response, known issues)
- Dependency audit (regular audits, health metrics, removing dependencies)
- Appendix (dependency tree, license texts, resources)

**When to Use**: When adding or updating dependencies, or when milestone requires specific packages

### keybinding-system.md
**Purpose**: Vim mode keybinding specifications
**Key Sections**:
- (Read this document for specific keybinding details)

**When to Use**: M34 (Vim State Machine), M35 (Vim Keybindings)

---

## Testing Guides

### FOUNDATION-TESTING-GUIDE.md
**Purpose**: Comprehensive testing guide for Foundation milestones (M1-M6)
**Key Sections**:
- M1: Bootstrap & Config testing
- M2: Basic TUI Shell testing
- M3: File I/O Foundation testing
- M4: Basic Text Editor testing
- M5: Auto-save testing
- M6: Undo/Redo testing
- Integration tests for foundation components
- End-to-end scenarios
- Performance benchmarks
- Test data and fixtures

**When to Use**: M1-M6 (Foundation milestones)

### LIBRARY-INTEGRATION-TESTING-GUIDE.md
**Purpose**: Comprehensive testing guide for Library Integration milestones (M7-M10)
**Key Sections**:
- M7: Library Loader testing
- M8: Library Browser UI testing
- M9: Fuzzy Search in Library testing
- M10: Prompt Insertion testing
- Integration tests for library components
- End-to-end scenarios
- Performance benchmarks
- Test data and fixtures

**When to Use**: M7-M10 (Library Integration milestones)

### PLACEHOLDER-TESTING-GUIDE.md
**Purpose**: Comprehensive testing guide for Placeholder milestones (M11-M14)
**Key Sections**:
- M11: Placeholder Parser testing
- M12: Placeholder Navigation testing
- M13: Text Placeholder Editing testing
- M14: List Placeholder Editing testing
- Integration tests for placeholder components
- End-to-end scenarios
- Performance benchmarks
- Test data and fixtures

**When to Use**: M11-M14 (Placeholder milestones)

### HISTORY-TESTING-GUIDE.md
**Purpose**: Comprehensive testing guide for History milestones (M15-M17)
**Key Sections**:
- M15: SQLite Setup testing
- M16: History Sync testing
- M17: History Browser testing
- Integration tests for history components
- End-to-end scenarios
- Performance benchmarks
- Test data and fixtures

**When to Use**: M15-M17 (History milestones)

### COMMANDS-FILES-TESTING-GUIDE.md
**Purpose**: Comprehensive testing guide for Commands & Files milestones (M18-M22)
**Key Sections**:
- M18: Command Registry testing
- M19: Command Palette UI testing
- M20: File Finder testing
- M21: Title Extraction testing
- M22: Batch Title Editor & Link Insertion testing
- Integration tests for commands and files components
- End-to-end scenarios
- Performance benchmarks
- Test data and fixtures

**When to Use**: M18-M22 (Commands & Files milestones)

### PROMPT-MANAGEMENT-TESTING-GUIDE.md
**Purpose**: Comprehensive testing guide for Prompt Management milestones (M23-M26)
**Key Sections**:
- M23: Prompt Validation testing
- M24: Validation Results Display testing
- M25: Prompt Creator testing
- M26: Prompt Editor testing
- Integration tests for prompt management components
- End-to-end scenarios
- Performance benchmarks
- Test data and fixtures

**When to Use**: M23-M26 (Prompt Management milestones)

### AI-INTEGRATION-TESTING-GUIDE.md
**Purpose**: Comprehensive testing guide for AI Integration milestones (M27-M33)
**Key Sections**:
- M27: Claude API Client testing
- M28: Context Selection Algorithm testing
- M29: Token Estimation & Budget testing
- M30: Suggestion Parsing testing
- M31: Suggestions Panel testing
- M32: Diff Generation testing
- M33: Diff Application testing
- Integration tests for AI components
- End-to-end scenarios
- Performance benchmarks
- Test data and fixtures

**When to Use**: M27-M33 (AI Integration milestones)

### VIM-MODE-TESTING-GUIDE.md
**Purpose**: Comprehensive testing guide for Vim Mode milestones (M34-M35)
**Key Sections**:
- M34: Vim State Machine testing
- M35: Vim Keybindings testing
- Integration tests for vim components
- End-to-end scenarios
- Performance benchmarks
- Test data and fixtures

**When to Use**: M34-M35 (Vim Mode milestones)

### POLISH-TESTING-GUIDE.md
**Purpose**: Comprehensive testing guide for Polish milestones (M36-M38)
**Key Sections**:
- M36: Settings Panel testing
- M37: Responsive Layout testing
- M38: Error Handling & Log Viewer testing
- Integration tests for polish components
- End-to-end scenarios
- Performance benchmarks
- Test data and fixtures

**When to Use**: M36-M38 (Polish milestones)

## Acceptance Criteria Documents

### ACCEPTANCE-CRITERIA-M16-HISTORY-SYNC.md
**Purpose**: Detailed acceptance criteria for M16 - History Sync
**Key Sections**:
- Functional Requirements (FR-1 to FR-6)
- Integration Requirements (IR-1 to IR-3)
- Edge Cases & Error Handling (EC-1 to EC-5)
- Performance Requirements (PR-1 to PR-3)
- User Experience Requirements (UX-1 to UX-3)
- Success Criteria (Must Have, Should Have, Nice to Have)
- Test cases with Go code examples

**When to Use**: M16 (History Sync)

### ACCEPTANCE-CRITERIA-M28-CONTEXT-SELECTION.md
**Purpose**: Detailed acceptance criteria for M28 - Context Selection Algorithm
**Key Sections**:
- Functional Requirements (FR-1 to FR-6)
- Integration Requirements (IR-1 to IR-4)
- Edge Cases & Error Handling (EC-1 to EC-5)
- Performance Requirements (PR-1 to PR-3)
- User Experience Requirements (UX-1 to UX-3)
- Success Criteria (Must Have, Should Have, Nice to Have)
- Test cases with Go code examples

**When to Use**: M28 (Context Selection Algorithm)

### ACCEPTANCE-CRITERIA-M32-DIFF-GENERATION.md
**Purpose**: Detailed acceptance criteria for M32 - Diff Generation
**Key Sections**:
- Functional Requirements (FR-1 to FR-6)
- Integration Requirements (IR-1 to IR-4)
- Edge Cases & Error Handling (EC-1 to EC-5)
- Performance Requirements (PR-1 to PR-3)
- User Experience Requirements (UX-1 to UX-3)
- Success Criteria (Must Have, Should Have, Nice to Have)
- Test cases with Go code examples

**When to Use**: M32 (Diff Generation)

### ACCEPTANCE-CRITERIA-M33-DIFF-APPLICATION.md
**Purpose**: Detailed acceptance criteria for M33 - Diff Application
**Key Sections**:
- Functional Requirements (FR-1 to FR-6)
- Integration Requirements (IR-1 to IR-4)
- Edge Cases & Error Handling (EC-1 to EC-5)
- Performance Requirements (PR-1 to PR-3)
- User Experience Requirements (UX-1 to UX-3)
- Success Criteria (Must Have, Should Have, Nice to Have)
- Test cases with Go code examples

**When to Use**: M33 (Diff Application)

### ACCEPTANCE-CRITERIA-M35-VIM-KEYBINDINGS.md
**Purpose**: Detailed acceptance criteria for M35 - Vim Keybindings
**Key Sections**:
- Functional Requirements (FR-1 to FR-6)
- Integration Requirements (IR-1 to IR-4)
- Edge Cases & Error Handling (EC-1 to EC-5)
- Performance Requirements (PR-1 to PR-3)
- User Experience Requirements (UX-1 to UX-3)
- Success Criteria (Must Have, Should Have, Nice to Have)
- Test cases with Go code examples

**When to Use**: M35 (Vim Keybindings)

### ACCEPTANCE-CRITERIA-M37-RESPONSIVE-LAYOUT.md
**Purpose**: Detailed acceptance criteria for M37 - Responsive Layout
**Key Sections**:
- Functional Requirements (FR-1 to FR-6)
- Integration Requirements (IR-1 to IR-4)
- Edge Cases & Error Handling (EC-1 to EC-5)
- Performance Requirements (PR-1 to PR-3)
- User Experience Requirements (UX-1 to UX-3)
- Success Criteria (Must Have, Should Have, Nice to Have)
- Test cases with Go code examples

**When to Use**: M37 (Responsive Layout)

### ACCEPTANCE-CRITERIA-M38-ERROR-HANDLING.md
**Purpose**: Detailed acceptance criteria for M38 - Error Handling & Log Viewer
**Key Sections**:
- Functional Requirements (FR-1 to FR-6)
- Integration Requirements (IR-1 to IR-4)
- Edge Cases & Error Handling (EC-1 to EC-5)
- Performance Requirements (PR-1 to PR-3)
- User Experience Requirements (UX-1 to UX-3)
- Success Criteria (Must Have, Should Have, Nice to Have)
- Test cases with Go code examples

**When to Use**: M38 (Error Handling & Log Viewer)

## Supporting Documents

### BUILD.md
**Purpose**: Build and deployment guide
**Key Sections**:
- Prerequisites (Go version, Git, Make, platform tools)
- Local development build (quick build, debug build, clean build, tests)
- Production builds (build matrix, macOS, Linux, Windows)
- Build optimization (linker flags, build tags, UPX compression)
- Versioning (version information, semantic versioning)
- Installation methods (Homebrew, Linux packages, Windows installer, manual)
- Release process (pre-release checklist, create tag, build all platforms, create release assets, GitHub release, update Homebrew)
- Update mechanism (version checking, auto-update, migration procedures)
- CI/CD integration (GitHub Actions workflow, Makefile)
- Troubleshooting (build issues, cross-compilation, installation, runtime)

**When to Use**: When building, packaging, or deploying PromptStack

### SETUP.md
**Purpose**: Development environment setup guide
**Key Sections**:
- Prerequisites (required software, optional but recommended)
- Installation steps (install Go, clone repository, install dependencies, verify Go environment)
- IDE configuration (VS Code, GoLand, Vim/Neovim)
- Environment variables (required, optional for development, setting variables)
- Verification steps (verify Go, project structure, dependencies, build, run tests, run application, verify tools)
- Troubleshooting (common issues and solutions)
- Platform-specific notes (macOS, Linux, Windows)
- Next steps (read documentation, start development, configure IDE, join community)
- Additional resources (Go docs, Effective Go, Bubble Tea docs)

**When to Use**: When setting up development environment

### HOW-TO-USE.md
**Purpose**: Development workflow guide
**Key Sections**:
- (Read this document for development workflow details)

**When to Use**: When starting development or understanding workflow

---

## Foundation Milestones (M1-M6)

### M1: Bootstrap & Config
**Documents**: Core Planning + CONFIG-SCHEMA.md + SETUP.md + FOUNDATION-TESTING-GUIDE.md
**Key Information**:
- Config structure from CONFIG-SCHEMA.md
- Setup wizard from requirements.md
- Logging setup from DEPENDENCIES.md (zap)
- Version tracking from CONFIG-SCHEMA.md
- Testing guide from FOUNDATION-TESTING-GUIDE.md

### M2: Basic TUI Shell
**Documents**: Core Planning + project-structure.md (UI domain) + FOUNDATION-TESTING-GUIDE.md
**Key Information**:
- Bubble Tea model structure from project-structure.md
- Status bar from project-structure.md
- Theme system from project-structure.md
- Testing guide from FOUNDATION-TESTING-GUIDE.md

### M3: File I/O Foundation
**Documents**: Core Planning + project-structure.md (platform/files) + FOUNDATION-TESTING-GUIDE.md
**Key Information**:
- File operations from project-structure.md
- YAML parsing from DEPENDENCIES.md (yaml.v3)
- Markdown handling from requirements.md
- Testing guide from FOUNDATION-TESTING-GUIDE.md

### M4: Basic Text Editor
**Documents**: Core Planning + project-structure.md (editor domain) + FOUNDATION-TESTING-GUIDE.md
**Key Information**:
- Editor domain structure from project-structure.md
- Text buffer from project-structure.md
- Cursor management from project-structure.md
- Testing guide from FOUNDATION-TESTING-GUIDE.md

### M5: Auto-save
**Documents**: Core Planning + project-structure.md (editor domain) + FOUNDATION-TESTING-GUIDE.md
**Key Information**:
- Auto-save strategy from requirements.md
- Debouncing from go-style-guide.md
- Status bar updates from project-structure.md
- Testing guide from FOUNDATION-TESTING-GUIDE.md

### M6: Undo/Redo
**Documents**: Core Planning + project-structure.md (editor domain) + FOUNDATION-TESTING-GUIDE.md
**Key Information**:
- Undo stack from project-structure.md
- Smart batching from requirements.md
- Keyboard shortcuts from requirements.md
- Testing guide from FOUNDATION-TESTING-GUIDE.md

---

## Library Integration Milestones (M7-M10)

### M7: Library Loader
**Documents**: Core Planning + project-structure.md (library domain) + requirements.md (Library section) + LIBRARY-INTEGRATION-TESTING-GUIDE.md
**Key Information**:
- Library domain structure from project-structure.md
- Prompt models from project-structure.md
- YAML frontmatter from requirements.md
- Testing guide from LIBRARY-INTEGRATION-TESTING-GUIDE.md

### M8: Library Browser UI
**Documents**: Core Planning + project-structure.md (ui/browser) + LIBRARY-INTEGRATION-TESTING-GUIDE.md
**Key Information**:
- Browser UI structure from project-structure.md
- Modal patterns from go-testing-guide.md
- Fuzzy search from DEPENDENCIES.md (sahilm/fuzzy)
- Testing guide from LIBRARY-INTEGRATION-TESTING-GUIDE.md

### M9: Fuzzy Search in Library
**Documents**: Core Planning + DEPENDENCIES.md (sahilm/fuzzy) + LIBRARY-INTEGRATION-TESTING-GUIDE.md
**Key Information**:
- Fuzzy matching algorithm from DEPENDENCIES.md
- Search patterns from go-testing-guide.md
- Performance considerations from DEPENDENCIES.md
- Testing guide from LIBRARY-INTEGRATION-TESTING-GUIDE.md

### M10: Prompt Insertion
**Documents**: Core Planning + project-structure.md (library domain) + LIBRARY-INTEGRATION-TESTING-GUIDE.md
**Key Information**:
- Library integration from project-structure.md
- Editor integration from project-structure.md
- Insertion logic from requirements.md
- Testing guide from LIBRARY-INTEGRATION-TESTING-GUIDE.md

---

## Placeholder Milestones (M11-M14)

### M11: Placeholder Parser
**Documents**: Core Planning + requirements.md (Placeholder System) + PLACEHOLDER-TESTING-GUIDE.md
**Key Information**:
- Placeholder syntax from requirements.md
- Validation rules from requirements.md
- Regex patterns from go-style-guide.md
- Testing guide from PLACEHOLDER-TESTING-GUIDE.md

### M12: Placeholder Navigation
**Documents**: Core Planning + project-structure.md (editor domain) + PLACEHOLDER-TESTING-GUIDE.md
**Key Information**:
- Tab navigation from requirements.md
- Placeholder tracking from project-structure.md
- Visual highlighting from project-structure.md
- Testing guide from PLACEHOLDER-TESTING-GUIDE.md

### M13: Text Placeholder Editing
**Documents**: Core Planning + project-structure.md (editor domain) + PLACEHOLDER-TESTING-GUIDE.md
**Key Information**:
- Edit mode from requirements.md
- Placeholder replacement from project-structure.md
- Vim integration from keybinding-system.md
- Testing guide from PLACEHOLDER-TESTING-GUIDE.md

### M14: List Placeholder Editing
**Documents**: Core Planning + project-structure.md (editor domain) + PLACEHOLDER-TESTING-GUIDE.md
**Key Information**:
- List editing UI from requirements.md
- List management from project-structure.md
- Markdown conversion from requirements.md
- Testing guide from PLACEHOLDER-TESTING-GUIDE.md

---

## History Milestones (M15-M17)

### M15: SQLite Setup
**Documents**: Core Planning + DATABASE-SCHEMA.md + DEPENDENCIES.md (modernc.org/sqlite) + HISTORY-TESTING-GUIDE.md
**Key Information**:
- Database schema from DATABASE-SCHEMA.md
- Table creation from DATABASE-SCHEMA.md
- Index creation from DATABASE-SCHEMA.md
- SQLite package from DEPENDENCIES.md
- Testing guide from HISTORY-TESTING-GUIDE.md

### M16: History Sync
**Documents**: Core Planning + DATABASE-SCHEMA.md + requirements.md (History section) + HISTORY-TESTING-GUIDE.md + ACCEPTANCE-CRITERIA-M16-HISTORY-SYNC.md
**Key Information**:
- Sync strategy from DATABASE-SCHEMA.md
- Triggers from DATABASE-SCHEMA.md
- Dual storage from DATABASE-SCHEMA.md
- Markdown operations from requirements.md
- Testing guide from HISTORY-TESTING-GUIDE.md
- Acceptance criteria from ACCEPTANCE-CRITERIA-M16-HISTORY-SYNC.md

### M17: History Browser
**Documents**: Core Planning + project-structure.md (ui/history) + DATABASE-SCHEMA.md (Query Patterns) + HISTORY-TESTING-GUIDE.md
**Key Information**:
- History UI from project-structure.md
- Query patterns from DATABASE-SCHEMA.md
- FTS5 search from DATABASE-SCHEMA.md
- Modal patterns from go-testing-guide.md
- Testing guide from HISTORY-TESTING-GUIDE.md

---

## Commands & Files Milestones (M18-M22)

### M18: Command Registry
**Documents**: Core Planning + project-structure.md (commands domain) + COMMANDS-FILES-TESTING-GUIDE.md
**Key Information**:
- Command system from project-structure.md
- Registry pattern from go-style-guide.md
- Interface design from go-style-guide.md
- Testing guide from COMMANDS-FILES-TESTING-GUIDE.md

### M19: Command Palette UI
**Documents**: Core Planning + project-structure.md (ui/palette) + COMMANDS-FILES-TESTING-GUIDE.md
**Key Information**:
- Palette UI from project-structure.md
- Fuzzy search from DEPENDENCIES.md (sahilm/fuzzy)
- Modal patterns from go-testing-guide.md
- Testing guide from COMMANDS-FILES-TESTING-GUIDE.md

### M20: File Finder
**Documents**: Core Planning + project-structure.md (platform/files) + DEPENDENCIES.md (go-gitignore) + COMMANDS-FILES-TESTING-GUIDE.md
**Key Information**:
- File traversal from project-structure.md
- Gitignore parsing from DEPENDENCIES.md
- Fuzzy search from DEPENDENCIES.md (sahilm/fuzzy)
- Testing guide from COMMANDS-FILES-TESTING-GUIDE.md

### M21: Title Extraction
**Documents**: Core Planning + project-structure.md (platform/files) + COMMANDS-FILES-TESTING-GUIDE.md
**Key Information**:
- YAML parsing from DEPENDENCIES.md (yaml.v3)
- Frontmatter extraction from requirements.md
- Fallback logic from requirements.md
- Testing guide from COMMANDS-FILES-TESTING-GUIDE.md

### M22: Batch Title Editor & Link Insertion
**Documents**: Core Planning + project-structure.md (ui/filereference) + COMMANDS-FILES-TESTING-GUIDE.md
**Key Information**:
- File reference UI from project-structure.md
- Markdown link format from requirements.md
- Batch editing from requirements.md
- Testing guide from COMMANDS-FILES-TESTING-GUIDE.md

---

## Prompt Management Milestones (M23-M26)

### M23: Prompt Validation
**Documents**: Core Planning + project-structure.md (library domain) + requirements.md (Validation section) + PROMPT-MANAGEMENT-TESTING-GUIDE.md
**Key Information**:
- Validation checks from requirements.md
- Error handling from go-style-guide.md
- Library domain from project-structure.md
- Testing guide from PROMPT-MANAGEMENT-TESTING-GUIDE.md

### M24: Validation Results Display
**Documents**: Core Planning + project-structure.md (ui/validation) + PROMPT-MANAGEMENT-TESTING-GUIDE.md
**Key Information**:
- Validation UI from project-structure.md
- Modal patterns from go-testing-guide.md
- Error display from requirements.md
- Testing guide from PROMPT-MANAGEMENT-TESTING-GUIDE.md

### M25: Prompt Creator
**Documents**: Core Planning + project-structure.md (ui/promptcreator) + PROMPT-MANAGEMENT-TESTING-GUIDE.md
**Key Information**:
- Creator UI from project-structure.md
- Form patterns from go-testing-guide.md
- File operations from project-structure.md
- Testing guide from PROMPT-MANAGEMENT-TESTING-GUIDE.md

### M26: Prompt Editor
**Documents**: Core Planning + project-structure.md (ui/prompteditor) + DEPENDENCIES.md (glamour) + PROMPT-MANAGEMENT-TESTING-GUIDE.md
**Key Information**:
- Editor UI from project-structure.md
- Markdown rendering from DEPENDENCIES.md (glamour)
- Edit/preview toggle from requirements.md
- Testing guide from PROMPT-MANAGEMENT-TESTING-GUIDE.md

---

## AI Integration Milestones (M27-M33)

### M27: Claude API Client
**Documents**: Core Planning + project-structure.md (ai domain) + DEPENDENCIES.md (anthropic-sdk-go) + AI-INTEGRATION-TESTING-GUIDE.md
**Key Information**:
- AI domain from project-structure.md
- API client from DEPENDENCIES.md
- Error handling from go-style-guide.md
- Retry logic from requirements.md
- Testing guide from AI-INTEGRATION-TESTING-GUIDE.md

### M28: Context Selection Algorithm
**Documents**: Core Planning + project-structure.md (ai domain) + requirements.md (AI Context Window Management) + AI-INTEGRATION-TESTING-GUIDE.md + ACCEPTANCE-CRITERIA-M28-CONTEXT-SELECTION.md
**Key Information**:
- Context selection from requirements.md
- Scoring algorithm from requirements.md
- Library indexing from project-structure.md
- Testing guide from AI-INTEGRATION-TESTING-GUIDE.md
- Acceptance criteria from ACCEPTANCE-CRITERIA-M28-CONTEXT-SELECTION.md

### M29: Token Estimation & Budget
**Documents**: Core Planning + project-structure.md (ai domain) + requirements.md (AI Context Window Management) + AI-INTEGRATION-TESTING-GUIDE.md
**Key Information**:
- Token counting from requirements.md
- Budget management from requirements.md
- Warning thresholds from requirements.md
- Testing guide from AI-INTEGRATION-TESTING-GUIDE.md

### M30: Suggestion Parsing
**Documents**: Core Planning + project-structure.md (ai domain) + requirements.md (AI Suggestions section) + AI-INTEGRATION-TESTING-GUIDE.md
**Key Information**:
- Suggestion types from requirements.md
- JSON parsing from DEPENDENCIES.md (tidwall/gjson)
- Error handling from go-style-guide.md
- Testing guide from AI-INTEGRATION-TESTING-GUIDE.md

### M31: Suggestions Panel
**Documents**: Core Planning + project-structure.md (ui/suggestions) + AI-INTEGRATION-TESTING-GUIDE.md
**Key Information**:
- Suggestions UI from project-structure.md
- Panel layout from requirements.md
- Accept/dismiss actions from requirements.md
- Testing guide from AI-INTEGRATION-TESTING-GUIDE.md

### M32: Diff Generation
**Documents**: Core Planning + project-structure.md (ai domain) + DEPENDENCIES.md (sergi/go-diff) + AI-INTEGRATION-TESTING-GUIDE.md + ACCEPTANCE-CRITERIA-M32-DIFF-GENERATION.md
**Key Information**:
- Diff generation from DEPENDENCIES.md
- Unified diff format from requirements.md
- API integration from project-structure.md
- Testing guide from AI-INTEGRATION-TESTING-GUIDE.md
- Acceptance criteria from ACCEPTANCE-CRITERIA-M32-DIFF-GENERATION.md

### M33: Diff Application
**Documents**: Core Planning + project-structure.md (ai domain) + AI-INTEGRATION-TESTING-GUIDE.md + ACCEPTANCE-CRITERIA-M33-DIFF-APPLICATION.md
**Key Information**:
- Diff application from project-structure.md
- Undo integration from project-structure.md
- Editor locking from requirements.md
- Testing guide from AI-INTEGRATION-TESTING-GUIDE.md
- Acceptance criteria from ACCEPTANCE-CRITERIA-M33-DIFF-APPLICATION.md

---

## Vim Mode Milestones (M34-M35)

### M34: Vim State Machine
**Documents**: Core Planning + project-structure.md (vim domain) + keybinding-system.md + VIM-MODE-TESTING-GUIDE.md
**Key Information**:
- Vim domain from project-structure.md
- State machine from keybinding-system.md
- Mode transitions from keybinding-system.md
- Testing guide from VIM-MODE-TESTING-GUIDE.md

### M35: Vim Keybindings
**Documents**: Core Planning + project-structure.md (vim domain) + keybinding-system.md + CONFIG-SCHEMA.md (vim_mode) + VIM-MODE-TESTING-GUIDE.md + ACCEPTANCE-CRITERIA-M35-VIM-KEYBINDINGS.md
**Key Information**:
- Keybinding maps from keybinding-system.md
- Context-aware routing from keybinding-system.md
- Vim mode toggle from CONFIG-SCHEMA.md
- Testing guide from VIM-MODE-TESTING-GUIDE.md
- Acceptance criteria from ACCEPTANCE-CRITERIA-M35-VIM-KEYBINDINGS.md

---

## Polish Milestones (M36-M38)

### M36: Settings Panel
**Documents**: Core Planning + project-structure.md (ui/settings) + CONFIG-SCHEMA.md + POLISH-TESTING-GUIDE.md
**Key Information**:
- Settings UI from project-structure.md
- Config schema from CONFIG-SCHEMA.md
- Form patterns from go-testing-guide.md
- Testing guide from POLISH-TESTING-GUIDE.md

### M37: Responsive Layout
**Documents**: Core Planning + project-structure.md (ui/app) + requirements.md (Split-Pane Layout) + POLISH-TESTING-GUIDE.md + ACCEPTANCE-CRITERIA-M37-RESPONSIVE-LAYOUT.md
**Key Information**:
- Layout management from project-structure.md
- Responsive behavior from requirements.md
- Terminal resize handling from go-testing-guide.md
- Testing guide from POLISH-TESTING-GUIDE.md
- Acceptance criteria from ACCEPTANCE-CRITERIA-M37-RESPONSIVE-LAYOUT.md

### M38: Error Handling & Log Viewer
**Documents**: Core Planning + project-structure.md (platform/errors) + CONFIG-SCHEMA.md (logging section) + POLISH-TESTING-GUIDE.md + ACCEPTANCE-CRITERIA-M38-ERROR-HANDLING.md
**Key Information**:
- Error handling from project-structure.md
- Logging from CONFIG-SCHEMA.md
- Log viewer UI from project-structure.md
- Testing guide from POLISH-TESTING-GUIDE.md
- Acceptance criteria from ACCEPTANCE-CRITERIA-M38-ERROR-HANDLING.md

---

## Domain Reference

### Editor Domain
**Milestones**: M4, M5, M6, M11, M12, M13, M14
**Key Documents**:
- project-structure.md (editor section)
- go-style-guide.md (type design, error handling)
- go-testing-guide.md (domain testing)
- requirements.md (Editor features)

### Library Domain
**Milestones**: M7, M8, M9, M10, M23, M24, M25, M26
**Key Documents**:
- project-structure.md (library section)
- requirements.md (Library section)
- DEPENDENCIES.md (sahilm/fuzzy, yaml.v3)
- go-style-guide.md (package organization)

### History Domain
**Milestones**: M15, M16, M17
**Key Documents**:
- project-structure.md (history section)
- DATABASE-SCHEMA.md (all sections)
- requirements.md (History section)
- DEPENDENCIES.md (modernc.org/sqlite)

### AI Domain
**Milestones**: M27, M28, M29, M30, M31, M32, M33
**Key Documents**:
- project-structure.md (ai domain)
- requirements.md (AI sections)
- DEPENDENCIES.md (anthropic-sdk-go, tidwall/*, sergi/go-diff)
- go-style-guide.md (error handling, concurrency)

### UI Domain
**Milestones**: M2, M8, M19, M24, M25, M26, M31, M36, M37
**Key Documents**:
- project-structure.md (ui packages)
- go-testing-guide.md (Bubble Tea patterns)
- requirements.md (UI features)
- project-structure.md (theme system)

### Platform Domain
**Milestones**: M1, M3, M15, M20, M21, M38
**Key Documents**:
- project-structure.md (platform packages)
- CONFIG-SCHEMA.md (config, logging)
- DATABASE-SCHEMA.md (database)
- DEPENDENCIES.md (zap, go-homedir)

### Vim Domain
**Milestones**: M34, M35
**Key Documents**:
- project-structure.md (vim domain)
- keybinding-system.md (all sections)
- CONFIG-SCHEMA.md (vim_mode)
- requirements.md (Vim support)

### Commands Domain
**Milestones**: M18, M19
**Key Documents**:
- project-structure.md (commands domain)
- go-style-guide.md (interface design)
- requirements.md (Command palette)

---

## Quick Lookup Tables

### Find File Paths
| What You Need | Document | Section |
|---------------|----------|---------|
| Package structure | project-structure.md | Recommended Project Structure |
| Config schema | CONFIG-SCHEMA.md | Complete Schema |
| Database schema | DATABASE-SCHEMA.md | Tables |
| Dependencies | DEPENDENCIES.md | Core Dependencies |
| Style patterns | go-style-guide.md | Type Design, Error Handling |
| Test patterns | go-testing-guide.md | Bubble Tea Testing |
| Vim keybindings | keybinding-system.md | (read document) |
| Build process | BUILD.md | Production Builds |
| Setup process | SETUP.md | Installation Steps |

### Find Code Patterns
| Pattern Type | Document | Section |
|--------------|----------|---------|
| Constructors | go-style-guide.md | Type Design |
| Error handling | go-style-guide.md | Error Handling |
| Interfaces | go-style-guide.md | Interfaces |
| Dependency injection | go-style-guide.md | Dependency Management |
| Table-driven tests | go-testing-guide.md | Table-Driven Tests |
| Mocking | go-testing-guide.md | Mocking |
| Bubble Tea models | go-testing-guide.md | Model Testing |
| Bubble Tea commands | go-testing-guide.md | Command Testing |

### Find Feature Specifications
| Feature | Document | Section |
|---------|----------|---------|
| Placeholder system | requirements.md | Placeholder System |
| AI suggestions | requirements.md | AI Suggestions |
| History sync | DATABASE-SCHEMA.md | Integration with Markdown Files |
| Vim mode | requirements.md | Vim Support |
| Split-pane layout | requirements.md | Split-Pane Layout Behavior |
| Auto-save | requirements.md | Auto-Save Strategy |
| Undo/redo | requirements.md | Undo/Redo |

---

## Document Relationships

### Core Dependencies
```
milestone-execution-prompt.md
    ↓ references
milestones.md, requirements.md, project-structure.md, go-style-guide.md, go-testing-guide.md
    ↓ references
CONFIG-SCHEMA.md, DATABASE-SCHEMA.md, DEPENDENCIES.md, keybinding-system.md
```

### Milestone-Specific Dependencies
```
M1 (Bootstrap & Config)
    ↓ requires
CONFIG-SCHEMA.md, SETUP.md

M15 (SQLite Setup)
    ↓ requires
DATABASE-SCHEMA.md, DEPENDENCIES.md (modernc.org/sqlite)

M27 (Claude API Client)
    ↓ requires
DEPENDENCIES.md (anthropic-sdk-go), requirements.md (AI sections)

M34 (Vim State Machine)
    ↓ requires
keybinding-system.md, CONFIG-SCHEMA.md (vim_mode)
```

---

## Usage Tips

### For AI Creating Implementation Plans

1. **Start with DOCUMENT-REFERENCE-MATRIX.md**
   - Find your milestone in the matrix
   - List all required documents

2. **Follow DOCUMENT-CHECKING-WORKFLOW.md**
   - Read documents in batches (up to 5 at a time)
   - Extract relevant information
   - Organize by category

3. **Use this DOCUMENT-INDEX.md**
   - Quickly locate specific information
   - Find code patterns
   - Look up feature specifications

4. **Reference documents explicitly**
   - Cite document sections in your plan
   - Include file paths from project-structure.md
   - Apply patterns from go-style-guide.md and go-testing-guide.md

### For Human Reviewers

1. **Check document coverage**
   - Verify all required documents were read
   - Check that document sections are referenced
   - Ensure patterns are applied correctly

2. **Verify consistency**
   - Check that file paths match project-structure.md
   - Verify code follows go-style-guide.md
   - Ensure tests follow go-testing-guide.md

3. **Validate completeness**
   - Check that all deliverables are addressed
   - Verify that technical specifications are applied
   - Ensure dependencies are correctly identified

---

**Last Updated**: 2026-01-07  
**Status**: Active - Use this index for quick reference