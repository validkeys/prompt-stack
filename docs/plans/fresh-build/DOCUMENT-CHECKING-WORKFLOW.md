# Document Checking Workflow

**Purpose**: Systematic workflow to ensure AI checks all relevant documents before creating implementation plans for each milestone.

---

## Overview

This workflow ensures that when AI creates an implementation plan for any milestone, it systematically reads and references all relevant documentation. This prevents missing critical information and ensures consistency across the codebase.

---

## Workflow Steps

### Step 1: Identify Milestone Context

Before reading any documents, identify:
1. **Milestone Number** (e.g., M1, M15, M27)
2. **Milestone Title** (e.g., "Bootstrap & Config", "SQLite Setup")
3. **Primary Domain** (e.g., config, history, ai, ui)
4. **Key Features** (from milestones.md deliverables)

### Step 2: Consult Document Reference Matrix

Use [`DOCUMENT-REFERENCE-MATRIX.md`](DOCUMENT-REFERENCE-MATRIX.md) to identify required documents:

1. **Find milestone row** in the matrix
2. **List all documents** in "Required Documents" column
3. **Note any additional context** from "Additional Context" column

### Step 3: Read Core Planning Documents (Always Required)

Read these documents first (up to 5 at a time using `read_file`):

1. [`milestone-execution-prompt.md`](milestone-execution-prompt.md)
   - Understand the execution process
   - Review task list format requirements
   - Note checkpoint documentation requirements

2. [`milestones.md`](milestones.md)
   - Read the specific milestone section
   - Extract goal, deliverables, test criteria
   - Note file paths mentioned

3. [`requirements.md`](requirements.md)
   - Read relevant sections for the milestone's features
   - Extract specific requirements
   - Note any constraints or edge cases

4. [`project-structure.md`](project-structure.md)
   - Read relevant domain sections
   - Identify package structure for the milestone
   - Note file paths and package organization

5. [`go-style-guide.md`](go-style-guide.md)
   - Review relevant coding patterns
   - Note style requirements for the domain
   - Extract applicable examples

6. [`go-testing-guide.md`](go-testing-guide.md)
   - Review testing patterns for the domain
   - Note TDD requirements
   - Extract applicable test examples

### Step 4: Read Context-Specific Documents

Based on the milestone, read additional documents (up to 5 at a time):

**For Config Milestones (M1, M36)**:
- [`CONFIG-SCHEMA.md`](CONFIG-SCHEMA.md) - Read relevant sections

**For Database Milestones (M15, M16, M17)**:
- [`DATABASE-SCHEMA.md`](DATABASE-SCHEMA.md) - Read relevant tables, queries, triggers

**For AI Milestones (M27-M33)**:
- [`DEPENDENCIES.md`](DEPENDENCIES.md) - Read anthropic-sdk-go section
- [`requirements.md`](requirements.md) - Read AI Context Window Management section

**For Vim Milestones (M34, M35)**:
- [`keybinding-system.md`](keybinding-system.md) - Read relevant sections
- [`CONFIG-SCHEMA.md`](CONFIG-SCHEMA.md) - Read vim_mode section

**For File Operations (M3, M20, M21)**:
- [`DEPENDENCIES.md`](DEPENDENCIES.md) - Read go-gitignore section
- [`project-structure.md`](project-structure.md) - Read platform/files section

**For UI Components**:
- [`project-structure.md`](project-structure.md) - Read relevant UI package sections
- [`go-testing-guide.md`](go-testing-guide.md) - Read Bubble Tea testing patterns

### Step 5: Extract and Organize Information

Create a structured summary of all information gathered:

```markdown
## Milestone Context
- **Number**: M{N}
- **Title**: {Title}
- **Domain**: {Primary Domain}
- **Goal**: {From milestones.md}

## Requirements Summary
- {Extracted from requirements.md}
- {Specific to this milestone}

## Architecture Context
- {From project-structure.md}
- {Package structure}
- {File paths}

## Style Guide References
- {From go-style-guide.md}
- {Relevant patterns}
- {Code examples}

## Testing Guide References
- {From go-testing-guide.md}
- {Test patterns}
- {Test examples}

## Technical Specifications
- {From CONFIG-SCHEMA.md or DATABASE-SCHEMA.md}
- {Relevant sections}
- {Data structures}

## Dependencies
- {From DEPENDENCIES.md}
- {Packages to use}
- {Version requirements}
```

### Step 6: Create Implementation Plan

Using all gathered information, create the implementation plan following the format in [`milestone-execution-prompt.md`](milestone-execution-prompt.md):

1. **Task List Document** (`M{N}-task-list.md`)
   - Include all tasks with dependencies
   - Reference file paths from project-structure.md
   - Apply style guide patterns
   - Include testing guide patterns

2. **Reference Document** (`M{N}-reference.md`)
   - Include architecture context
   - Include code examples from style guide
   - Include test examples from testing guide
   - Include technical specifications

### Step 7: Verify Completeness

Before finalizing the plan, verify:

- [ ] All required documents have been read
- [ ] All deliverables from milestones.md are addressed
- [ ] All file paths match project-structure.md
- [ ] Code examples follow go-style-guide.md
- [ ] Test examples follow go-testing-guide.md
- [ ] Technical specifications are correctly applied
- [ ] Dependencies are correctly identified
- [ ] Integration points are noted

---

## Document Reading Strategy

### Efficient Batch Reading

Use the `read_file` tool to read multiple documents at once (up to 5 files):

```xml
<read_file>
<args>
  <file>
    <path>doc1.md</path>
  </file>
  <file>
    <path>doc2.md</path>
  </file>
  <file>
    <path>doc3.md</path>
  </file>
  <file>
    <path>doc4.md</path>
  </file>
  <file>
    <path>doc5.md</path>
  </file>
</args>
</read_file>
```

### Reading Order Priority

1. **Core Planning Documents** (Always first)
   - milestone-execution-prompt.md
   - milestones.md (specific milestone)
   - requirements.md (relevant sections)
   - project-structure.md (relevant domains)
   - go-style-guide.md
   - go-testing-guide.md

2. **Context-Specific Documents** (Second)
   - CONFIG-SCHEMA.md (if config-related)
   - DATABASE-SCHEMA.md (if database-related)
   - DEPENDENCIES.md (if external packages needed)
   - keybinding-system.md (if vim-related)

3. **Supporting Documents** (As needed)
   - BUILD.md (if build questions)
   - SETUP.md (if environment issues)
   - HOW-TO-USE.md (if workflow questions)

### Focused Reading

Don't read entire documents if not needed. Focus on:

- **Specific sections** relevant to the milestone
- **Domain-specific** information
- **Code examples** that apply to the milestone
- **Test patterns** for the milestone's domain

Example: For M15 (SQLite Setup), read:
- DATABASE-SCHEMA.md: Focus on "Tables" and "Query Patterns" sections
- project-structure.md: Focus on "internal/history/" section
- DEPENDENCIES.md: Focus on "modernc.org/sqlite" section

---

## Common Pitfalls to Avoid

### ❌ Don't Skip Documents

**Wrong**: Only read milestones.md and start coding
**Right**: Read all required documents from the matrix first

### ❌ Don't Read Everything

**Wrong**: Read all 15+ documents for every milestone
**Right**: Read only the 6-10 documents relevant to the milestone

### ❌ Don't Ignore Context

**Wrong**: Read documents but don't reference them in the plan
**Right**: Explicitly reference document sections in task descriptions

### ❌ Don't Mix Domains

**Wrong**: Apply editor domain patterns to AI domain code
**Right**: Use the correct domain's patterns from project-structure.md

### ❌ Don't Forget Testing

**Wrong**: Create implementation plan without test examples
**Right**: Include test examples from go-testing-guide.md for each task

---

## Verification Checklist

Before presenting the implementation plan, ensure:

### Document Coverage
- [ ] Core planning documents read (6 documents)
- [ ] Context-specific documents read (2-4 documents)
- [ ] All documents from reference matrix consulted
- [ ] Document sections properly referenced

### Plan Quality
- [ ] All deliverables from milestones.md addressed
- [ ] File paths match project-structure.md
- [ ] Code follows go-style-guide.md patterns
- [ ] Tests follow go-testing-guide.md patterns
- [ ] Technical specifications correctly applied
- [ ] Dependencies correctly identified
- [ ] Integration points noted

### Plan Completeness
- [ ] Task list created with clear dependencies
- [ ] Reference document created with examples
- [ ] Acceptance criteria are testable
- [ ] File paths are explicit
- [ ] Integration points are documented

---

## Example Workflow: Milestone 15 (SQLite Setup)

### Step 1: Identify Context
- Milestone: M15
- Title: SQLite Setup
- Domain: history
- Key Features: Database schema, CRUD operations

### Step 2: Consult Matrix
From DOCUMENT-REFERENCE-MATRIX.md:
- Required: Core Planning + DATABASE-SCHEMA.md + DEPENDENCIES.md

### Step 3: Read Core Documents (Batch 1)
```xml
<read_file>
<args>
  <file>
    <path>milestone-execution-prompt.md</path>
  </file>
  <file>
    <path>milestones.md</path>
  </file>
  <file>
    <path>requirements.md</path>
  </file>
  <file>
    <path>project-structure.md</path>
  </file>
  <file>
    <path>go-style-guide.md</path>
  </file>
</args>
</read_file>
```

### Step 4: Read Context Documents (Batch 2)
```xml
<read_file>
<args>
  <file>
    <path>go-testing-guide.md</path>
  </file>
  <file>
    <path>DATABASE-SCHEMA.md</path>
  </file>
  <file>
    <path>DEPENDENCIES.md</path>
  </file>
</args>
</read_file>
```

### Step 5: Extract Information
- From milestones.md: Goal, deliverables, test criteria
- From requirements.md: History section requirements
- From project-structure.md: internal/history/ package structure
- From DATABASE-SCHEMA.md: Table schemas, query patterns
- From DEPENDENCIES.md: modernc.org/sqlite package details
- From go-style-guide.md: Error handling patterns
- From go-testing-guide.md: Database testing patterns

### Step 6: Create Plan
- Task list with database setup tasks
- Reference document with SQL examples
- Include test patterns for database operations

### Step 7: Verify
- All required documents read ✓
- All deliverables addressed ✓
- File paths correct ✓
- Style guide followed ✓
- Testing guide followed ✓

---

## Quick Reference

### Document Categories

| Category | Documents | When to Read |
|-----------|-----------|---------------|
| **Core Planning** | milestone-execution-prompt.md, milestones.md, requirements.md, project-structure.md, go-style-guide.md, go-testing-guide.md | Always (every milestone) |
| **Config** | CONFIG-SCHEMA.md | M1, M36 |
| **Database** | DATABASE-SCHEMA.md | M15, M16, M17 |
| **AI** | DEPENDENCIES.md (anthropic-sdk-go), requirements.md (AI section) | M27-M33 |
| **Vim** | keybinding-system.md, CONFIG-SCHEMA.md (vim_mode) | M34, M35 |
| **Files** | DEPENDENCIES.md (go-gitignore) | M3, M20, M21 |
| **UI** | project-structure.md (ui packages), go-testing-guide.md (Bubble Tea) | M2, M8, M19, M24, M25, M26, M31, M36, M37 |

### Milestone Groups

| Group | Milestones | Key Documents |
|-------|-----------|---------------|
| **Foundation** | M1-M6 | Core Planning + CONFIG-SCHEMA.md (M1) |
| **Library** | M7-M10 | Core Planning + DEPENDENCIES.md (sahilm/fuzzy) |
| **Placeholders** | M11-M14 | Core Planning |
| **History** | M15-M17 | Core Planning + DATABASE-SCHEMA.md |
| **Commands** | M18-M22 | Core Planning + DEPENDENCIES.md (go-gitignore) |
| **Prompt Mgmt** | M23-M26 | Core Planning + DEPENDENCIES.md (glamour) |
| **AI** | M27-M33 | Core Planning + DEPENDENCIES.md (anthropic-sdk-go) |
| **Vim** | M34-M35 | Core Planning + keybinding-system.md |
| **Polish** | M36-M38 | Core Planning + CONFIG-SCHEMA.md (M36) |

---

## Integration with Milestone Execution Prompt

This workflow integrates directly with [`milestone-execution-prompt.md`](milestone-execution-prompt.md):

### Before Step 1: Load Current Milestone
- Execute this document checking workflow
- Read all required documents
- Extract and organize information

### Before Step 2: Generate Task List
- Use extracted information to create task list
- Reference all read documents
- Include file paths and patterns

### Before Step 3: Execute Each Task
- Reference style guide for code patterns
- Reference testing guide for test patterns
- Reference technical specifications

---

**Last Updated**: 2026-01-07  
**Status**: Active - Use this workflow for all milestone planning