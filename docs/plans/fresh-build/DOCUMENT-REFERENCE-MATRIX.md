# Document Reference Matrix

**Purpose**: This matrix maps each milestone to the relevant documents that must be consulted when creating implementation plans.

**Usage**: Before starting any milestone, AI should read all documents listed in the "Required Documents" column for that milestone.

---

## Document Categories

### Core Planning Documents (Always Required)
- [`milestone-execution-prompt.md`](milestone-execution-prompt.md) - Main execution workflow
- [`milestones.md`](milestones.md) - All milestone definitions
- [`requirements.md`](requirements.md) - Complete feature requirements
- [`project-structure.md`](project-structure.md) - Domain architecture and package structure

### Implementation Guides (Always Required)
- [`go-style-guide.md`](go-style-guide.md) - Go coding standards
- [`go-testing-guide.md`](go-testing-guide.md) - Testing patterns and TDD
- [`ENHANCED-TEST-CRITERIA-TEMPLATE.md`](milestones/ENHANCED-TEST-CRITERIA-TEMPLATE.md) - Acceptance criteria template

### Technical Specifications (Context-Specific)
- [`CONFIG-SCHEMA.md`](CONFIG-SCHEMA.md) - Configuration structure
- [`DATABASE-SCHEMA.md`](DATABASE-SCHEMA.md) - SQLite schema
- [`DEPENDENCIES.md`](DEPENDENCIES.md) - External dependencies
- [`keybinding-system.md`](keybinding-system.md) - Vim mode details

### Key Learnings (Context-Specific)
- [`learnings/go-fundamentals.md`](learnings/go-fundamentals.md) - Go-specific patterns and pitfalls
- [`learnings/editor-domain.md`](learnings/editor-domain.md) - Editor implementation patterns
- [`learnings/ui-domain.md`](learnings/ui-domain.md) - UI/TUI implementation patterns
- [`learnings/error-handling.md`](learnings/error-handling.md) - Error handling patterns
- [`learnings/ai-domain.md`](learnings/ai-domain.md) - AI integration patterns
- [`learnings/vim-domain.md`](learnings/vim-domain.md) - Vim mode patterns
- [`learnings/history-domain.md`](learnings/history-domain.md) - History management patterns
- [`learnings/architecture-patterns.md`](learnings/architecture-patterns.md) - Architecture and design patterns

### Supporting Documents (As Needed)
- [`BUILD.md`](BUILD.md) - Build processes
- [`SETUP.md`](SETUP.md) - Development setup
- [`HOW-TO-USE.md`](HOW-TO-USE.md) - Development workflow

---

## Milestone Document Mapping

### Foundation Milestones (1-6)

| Milestone | Title | Required Documents | Additional Context |
|-----------|-------|-------------------|-------------------|
| **M1** | Bootstrap & Config | Core Planning, Implementation Guides, CONFIG-SCHEMA.md | SETUP.md for environment setup, FOUNDATION-TESTING-GUIDE.md |
| **M2** | Basic TUI Shell | Core Planning, Implementation Guides, project-structure.md (UI domain) | FOUNDATION-TESTING-GUIDE.md |
| **M3** | File I/O Foundation | Core Planning, Implementation Guides, project-structure.md (platform/files) | FOUNDATION-TESTING-GUIDE.md |
| **M4** | Basic Text Editor | Core Planning, Implementation Guides, project-structure.md (editor domain) | FOUNDATION-TESTING-GUIDE.md |
| **M5** | Auto-save | Core Planning, Implementation Guides, project-structure.md (editor domain) | FOUNDATION-TESTING-GUIDE.md |
| **M6** | Undo/Redo | Core Planning, Implementation Guides, project-structure.md (editor domain) | FOUNDATION-TESTING-GUIDE.md |

### Library Integration Milestones (7-10)

| Milestone | Title | Required Documents | Additional Context |
|-----------|-------|-------------------|-------------------|
| **M7** | Library Loader | Core Planning, Implementation Guides, project-structure.md (library domain) | requirements.md (Library section), LIBRARY-INTEGRATION-TESTING-GUIDE.md |
| **M8** | Library Browser UI | Core Planning, Implementation Guides, project-structure.md (ui/browser) | LIBRARY-INTEGRATION-TESTING-GUIDE.md |
| **M9** | Fuzzy Search in Library | Core Planning, Implementation Guides, DEPENDENCIES.md (sahilm/fuzzy) | LIBRARY-INTEGRATION-TESTING-GUIDE.md |
| **M10** | Prompt Insertion | Core Planning, Implementation Guides, project-structure.md (library domain) | LIBRARY-INTEGRATION-TESTING-GUIDE.md |

### Placeholder Milestones (11-14)

| Milestone | Title | Required Documents | Additional Context |
|-----------|-------|-------------------|-------------------|
| **M11** | Placeholder Parser | Core Planning, Implementation Guides, requirements.md (Placeholder System) | PLACEHOLDER-TESTING-GUIDE.md |
| **M12** | Placeholder Navigation | Core Planning, Implementation Guides, project-structure.md (editor domain) | PLACEHOLDER-TESTING-GUIDE.md |
| **M13** | Text Placeholder Editing | Core Planning, Implementation Guides, project-structure.md (editor domain) | PLACEHOLDER-TESTING-GUIDE.md |
| **M14** | List Placeholder Editing | Core Planning, Implementation Guides, project-structure.md (editor domain) | PLACEHOLDER-TESTING-GUIDE.md |

### History Milestones (15-17)

| Milestone | Title | Required Documents | Additional Context |
|-----------|-------|-------------------|-------------------|
| **M15** | SQLite Setup | Core Planning, Implementation Guides, DATABASE-SCHEMA.md | DEPENDENCIES.md (modernc.org/sqlite), HISTORY-TESTING-GUIDE.md |
| **M16** | History Sync | Core Planning, Implementation Guides, DATABASE-SCHEMA.md | requirements.md (History section), HISTORY-TESTING-GUIDE.md, ACCEPTANCE-CRITERIA-M16-HISTORY-SYNC.md |
| **M17** | History Browser | Core Planning, Implementation Guides, project-structure.md (ui/history) | DATABASE-SCHEMA.md (Query Patterns), HISTORY-TESTING-GUIDE.md |

### Commands & Files Milestones (18-22)

| Milestone | Title | Required Documents | Additional Context |
|-----------|-------|-------------------|-------------------|
| **M18** | Command Registry | Core Planning, Implementation Guides, project-structure.md (commands domain) | COMMANDS-FILES-TESTING-GUIDE.md |
| **M19** | Command Palette UI | Core Planning, Implementation Guides, project-structure.md (ui/palette) | COMMANDS-FILES-TESTING-GUIDE.md |
| **M20** | File Finder | Core Planning, Implementation Guides, project-structure.md (platform/files) | DEPENDENCIES.md (go-gitignore), COMMANDS-FILES-TESTING-GUIDE.md |
| **M21** | Title Extraction | Core Planning, Implementation Guides, project-structure.md (platform/files) | COMMANDS-FILES-TESTING-GUIDE.md |
| **M22** | Batch Title Editor & Link Insertion | Core Planning, Implementation Guides, project-structure.md (ui/filereference) | COMMANDS-FILES-TESTING-GUIDE.md |

### Prompt Management Milestones (23-26)

| Milestone | Title | Required Documents | Additional Context |
|-----------|-------|-------------------|-------------------|
| **M23** | Prompt Validation | Core Planning, Implementation Guides, project-structure.md (library domain) | requirements.md (Validation section), PROMPT-MANAGEMENT-TESTING-GUIDE.md |
| **M24** | Validation Results Display | Core Planning, Implementation Guides, project-structure.md (ui/validation) | PROMPT-MANAGEMENT-TESTING-GUIDE.md |
| **M25** | Prompt Creator | Core Planning, Implementation Guides, project-structure.md (ui/promptcreator) | PROMPT-MANAGEMENT-TESTING-GUIDE.md |
| **M26** | Prompt Editor | Core Planning, Implementation Guides, project-structure.md (ui/prompteditor) | DEPENDENCIES.md (glamour), PROMPT-MANAGEMENT-TESTING-GUIDE.md |

### AI Integration Milestones (27-33)

| Milestone | Title | Required Documents | Additional Context |
|-----------|-------|-------------------|-------------------|
| **M27** | Claude API Client | Core Planning, Implementation Guides, project-structure.md (ai domain) | DEPENDENCIES.md (anthropic-sdk-go), AI-INTEGRATION-TESTING-GUIDE.md |
| **M28** | Context Selection Algorithm | Core Planning, Implementation Guides, project-structure.md (ai domain) | requirements.md (AI Context Window Management), AI-INTEGRATION-TESTING-GUIDE.md, ACCEPTANCE-CRITERIA-M28-CONTEXT-SELECTION.md |
| **M29** | Token Estimation & Budget | Core Planning, Implementation Guides, project-structure.md (ai domain) | requirements.md (AI Context Window Management), AI-INTEGRATION-TESTING-GUIDE.md |
| **M30** | Suggestion Parsing | Core Planning, Implementation Guides, project-structure.md (ai domain) | requirements.md (AI Suggestions section), AI-INTEGRATION-TESTING-GUIDE.md |
| **M31** | Suggestions Panel | Core Planning, Implementation Guides, project-structure.md (ui/suggestions) | AI-INTEGRATION-TESTING-GUIDE.md |
| **M32** | Diff Generation | Core Planning, Implementation Guides, project-structure.md (ai domain) | DEPENDENCIES.md (sergi/go-diff), AI-INTEGRATION-TESTING-GUIDE.md, ACCEPTANCE-CRITERIA-M32-DIFF-GENERATION.md |
| **M33** | Diff Application | Core Planning, Implementation Guides, project-structure.md (ai domain) | AI-INTEGRATION-TESTING-GUIDE.md, ACCEPTANCE-CRITERIA-M33-DIFF-APPLICATION.md |

### Vim Mode Milestones (34-35)

| Milestone | Title | Required Documents | Additional Context |
|-----------|-------|-------------------|-------------------|
| **M34** | Vim State Machine | Core Planning, Implementation Guides, project-structure.md (vim domain) | keybinding-system.md, VIM-MODE-TESTING-GUIDE.md |
| **M35** | Vim Keybindings | Core Planning, Implementation Guides, project-structure.md (vim domain) | keybinding-system.md, CONFIG-SCHEMA.md (vim_mode), VIM-MODE-TESTING-GUIDE.md, ACCEPTANCE-CRITERIA-M35-VIM-KEYBINDINGS.md |

### Polish Milestones (36-38)

| Milestone | Title | Required Documents | Additional Context |
|-----------|-------|-------------------|-------------------|
| **M36** | Settings Panel | Core Planning, Implementation Guides, project-structure.md (ui/settings), learnings/ui-domain.md | CONFIG-SCHEMA.md, POLISH-TESTING-GUIDE.md |
| **M37** | Responsive Layout | Core Planning, Implementation Guides, project-structure.md (ui/app), learnings/ui-domain.md | requirements.md (Split-Pane Layout), POLISH-TESTING-GUIDE.md, ACCEPTANCE-CRITERIA-M37-RESPONSIVE-LAYOUT.md |
| **M38** | Error Handling & Log Viewer | Core Planning, Implementation Guides, project-structure.md (platform/errors), learnings/error-handling.md | CONFIG-SCHEMA.md (logging section), POLISH-TESTING-GUIDE.md, ACCEPTANCE-CRITERIA-M38-ERROR-HANDLING.md |

---

## Document Checking Workflow

### Before Starting Any Milestone

1. **Read Core Planning Documents** (Always)
   - [`milestone-execution-prompt.md`](milestone-execution-prompt.md)
   - [`milestones.md`](milestones.md) - Read the specific milestone section
   - [`requirements.md`](requirements.md) - Read relevant sections
   - [`project-structure.md`](project-structure.md) - Read relevant domain sections

2. **Read Implementation Guides** (Always)
   - [`go-style-guide.md`](go-style-guide.md)
   - [`go-testing-guide.md`](go-testing-guide.md)

3. **Read Context-Specific Documents** (Per Matrix)
   - Check the matrix above for additional documents
   - Read only the relevant sections of those documents

4. **Create Implementation Plan**
   - Use the task list format from milestone-execution-prompt.md
   - Reference specific sections from all read documents
   - Include file paths from project-structure.md
   - Follow patterns from go-style-guide.md and go-testing-guide.md

### Document Reading Strategy

**Efficient Reading**:
- Use `read_file` tool to read multiple documents at once (up to 5 files)
- Prioritize documents by relevance to the milestone
- Focus on specific sections rather than entire documents

**Document Sections to Reference**:
- [`project-structure.md`](project-structure.md): Domain-specific package structures
- [`requirements.md`](requirements.md): Feature-specific requirements
- [`CONFIG-SCHEMA.md`](CONFIG-SCHEMA.md): Configuration fields for the feature
- [`DATABASE-SCHEMA.md`](DATABASE-SCHEMA.md): Relevant tables and queries
- [`DEPENDENCIES.md`](DEPENDENCIES.md): Specific packages to use

---

## Quick Reference by Domain

### Editor Domain
**Milestones**: M4, M5, M6, M11, M12, M13, M14
**Key Documents**: project-structure.md (editor section), go-style-guide.md, go-testing-guide.md, learnings/editor-domain.md

### Library Domain
**Milestones**: M7, M8, M9, M10, M23, M24, M25, M26
**Key Documents**: project-structure.md (library section), requirements.md (Library section), DEPENDENCIES.md (sahilm/fuzzy), learnings/architecture-patterns.md

### History Domain
**Milestones**: M15, M16, M17
**Key Documents**: project-structure.md (history section), DATABASE-SCHEMA.md, requirements.md (History section), learnings/history-domain.md

### AI Domain
**Milestones**: M27, M28, M29, M30, M31, M32, M33
**Key Documents**: project-structure.md (ai domain), DEPENDENCIES.md (anthropic-sdk-go), requirements.md (AI section), learnings/ai-domain.md

### UI Domain
**Milestones**: M2, M8, M19, M24, M25, M26, M31, M36, M37
**Key Documents**: project-structure.md (ui packages), go-testing-guide.md (Bubble Tea patterns), learnings/ui-domain.md

### Platform Domain
**Milestones**: M1, M3, M15, M20, M21, M38
**Key Documents**: project-structure.md (platform packages), CONFIG-SCHEMA.md, DATABASE-SCHEMA.md, learnings/go-fundamentals.md

### Vim Domain
**Milestones**: M34, M35
**Key Documents**: project-structure.md (vim domain), keybinding-system.md, CONFIG-SCHEMA.md (vim_mode), learnings/vim-domain.md

### Commands Domain
**Milestones**: M18, M19
**Key Documents**: project-structure.md (commands domain), go-style-guide.md (interface patterns), learnings/architecture-patterns.md

---

## Document Priority Levels

### Level 1: Always Required (Read First)
1. [`milestone-execution-prompt.md`](milestone-execution-prompt.md)
2. [`milestones.md`](milestones.md) - Specific milestone section
3. [`requirements.md`](requirements.md) - Relevant sections
4. [`project-structure.md`](project-structure.md) - Relevant domain sections
5. [`go-style-guide.md`](go-style-guide.md)
6. [`go-testing-guide.md`](go-testing-guide.md)

### Level 2: Context-Specific (Read if Applicable)
7. [`CONFIG-SCHEMA.md`](CONFIG-SCHEMA.md) - For config-related milestones
8. [`DATABASE-SCHEMA.md`](DATABASE-SCHEMA.md) - For database-related milestones
9. [`DEPENDENCIES.md`](DEPENDENCIES.md) - For external package usage
10. [`keybinding-system.md`](keybinding-system.md) - For vim mode milestones
11. [`learnings/go-fundamentals.md`](learnings/go-fundamentals.md) - For Go-specific patterns
12. [`learnings/editor-domain.md`](learnings/editor-domain.md) - For editor milestones
13. [`learnings/ui-domain.md`](learnings/ui-domain.md) - For UI/TUI milestones
14. [`learnings/error-handling.md`](learnings/error-handling.md) - For error handling
15. [`learnings/ai-domain.md`](learnings/ai-domain.md) - For AI milestones
16. [`learnings/vim-domain.md`](learnings/vim-domain.md) - For vim milestones
17. [`learnings/history-domain.md`](learnings/history-domain.md) - For history milestones
18. [`learnings/architecture-patterns.md`](learnings/architecture-patterns.md) - For architecture decisions

### Level 3: Supporting (Read as Needed)
11. [`BUILD.md`](BUILD.md) - For build/deployment questions
12. [`SETUP.md`](SETUP.md) - For environment setup issues
13. [`HOW-TO-USE.md`](HOW-TO-USE.md) - For workflow questions

---

## Implementation Plan Checklist

When creating an implementation plan for any milestone, ensure you have:

- [ ] Read the milestone definition from [`milestones.md`](milestones.md)
- [ ] Read relevant sections from [`requirements.md`](requirements.md)
- [ ] Read relevant domain structure from [`project-structure.md`](project-structure.md)
- [ ] Read [`go-style-guide.md`](go-style-guide.md) for coding standards
- [ ] Read [`go-testing-guide.md`](go-testing-guide.md) for testing patterns
- [ ] Read context-specific documents (CONFIG-SCHEMA, DATABASE-SCHEMA, etc.)
- [ ] Read relevant key learnings from [`learnings/`](learnings/) directory
- [ ] Referenced specific file paths from project structure
- [ ] Included relevant dependencies from DEPENDENCIES.md
- [ ] Applied style guide patterns to code examples
- [ ] Applied testing guide patterns to test examples
- [ ] Applied key learnings to implementation approach
- [ ] Noted any deviations from key learnings with justification
- [ ] Created task list with clear dependencies
- [ ] Created reference document with code examples

---

**Last Updated**: 2026-01-07  
**Status**: Active - Use this matrix for all milestone planning