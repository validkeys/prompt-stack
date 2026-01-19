# AI-Assisted Development Workflow Tool: Requirements Document

**Version:** 2.0 Final  
**Date:** January 2026  
**Purpose:** Meta-orchestration tool with dual modes (Plan/Build) that generates and validates perfect Ralphy YAML files

---

## Interview Findings

- Primary users: senior developer (tool is intended as a standalone developer tool).
- Primary integrations at launch: Ralphy and OpenCode. Minimal recommended set for MVP: Git (local + remote awareness), Ralphy (execution), Anthropic/Claude (AI), SQLite (local knowledge cache). Optional: CI providers, LSP/editor integrations.
- Access & security defaults:
  - Store API keys in OS-provided secret stores or require env vars; do not persist keys in the knowledge DB.
  - Default to local execution; provide an explicit opt-in remote mode with TLS + token auth.
  - Require explicit user confirmation before performing commits or creating branches; use minimal CI token scopes.
  - Maintain a local append-only `audit.log`; remote telemetry is optional and opt-in.
  - Never include secrets in `tasks.yaml`; support secure vault integration for runtime secrets during Build Mode.
- Config and UX expectations: sensible built-in defaults with optional repo-root configuration or CLI override; POC should be interactive CLI only and rely on sensible defaults.
- Artifacts produced by runs (MVP): `tasks.yaml`, `review-report.json`, `audit.log`, optional `task-trace/` for per-task execution traces, and a short human-readable `report.txt` printed to stdout and saved.
- Build Mode defaults per interview: `--commit-per-task` enabled by default; `--branch-per-task` disabled by default.
- Non-functional targets (MVP suggestions): file scanning up to 1,000 files in <5s on an SSD dev machine; code generation <5s; hybrid plan generation <30s; YAML validation <1s; memory footprint <512MB; SQLite DB <50MB initially; default concurrency: 3 parallel agents.

## Executive Summary

### What This Tool Does

Creates quality-guaranteed implementation plans through **dual-mode operation**:

**Plan Mode**: Generates + validates Ralphy YAML files  
**Build Mode**: Executes validated plans through Ralphy/OpenCode

(See `/docs/initial-claude-conversation/split-files/03-opencode-integrations-discussion.md` for detailed OpenCode integration patterns, agent examples, MCP server usage, and recommended `.opencode` project configs. Also see `/docs/initial-claude-conversation/split-files/04-cli-tool-building-discussion.md` for a scoped assessment and recommended approaches for building or extending a CLI tool that integrates Ralphy and OpenCode, and `/docs/initial-claude-conversation/split-files/06-knowledge-management.md` for CLI knowledge-management commands, storage strategies, and team-sharing examples.)

### The Core Innovation

**Self-validating, self-improving YAML generation:**

```
Requirements → Plan Mode → YAML
                  ↓
            [Code Generation: Fast, deterministic]
                  ↓
            [AI Validation: Quality check]
                  ↓
            [AI Review: Research compliance]
                  ↓
            Perfect YAML → Build Mode → Quality Code
```

### Key Value

- **90% context reduction**: Non-AI filtering before sending to models
- **95% first-pass success**: Research-backed constraints prevent drift
- **10-20x faster**: Compared to manual YAML authoring + debugging
- **Team knowledge sharing**: Cached patterns benefit entire organization

### Origin: The Initial Problem Statement

This project originated from research into solving AI code generation drift and quality issues experienced by senior developers. The initial problem statement documented in `/docs/initial-claude-conversation/split-files/01-initial-problem-statement.md` identified key challenges:

- **Architectural drift** across incremental task execution
- **Type safety violations** in TypeScript strict mode
- **Context retention issues** when breaking work into sequential tasks
- **Regression from initial plans** during AI implementation

This tool was designed specifically to address these pain points through research-backed constraints, multi-layer validation, and self-validating YAML generation. The research findings documented in `/docs/initial-claude-conversation/split-files/02-research-report-preventing-drift.md` (see `/Users/kyledavis/Sites/prompt-stack/docs/best-practices.md`) form the foundation for the best practices implemented in this tool. For knowledge management practices, CLI commands, and team-sharing examples see `/docs/initial-claude-conversation/split-files/06-knowledge-management.md`.

---

## Core Modes

### Plan Mode

**Purpose**: Generate validated Ralphy YAML from requirements

**Three Generation Methods**:

1. **Code Generation** (Fast Path - Default)
   - Deterministic template-based generation
    - Uses cached knowledge base (see `/docs/initial-claude-conversation/split-files/05-jit-caching.md` for JIT caching details and `/docs/initial-claude-conversation/split-files/06-knowledge-management.md` for CLI knowledge-management commands, caching strategies, and team-sharing examples)
   - 2-5 seconds execution
   - Best for: Similar tasks, simple requirements

2. **AI Generation** (Quality Path)
   - Uses Ralphy to generate YAML via meta-PRD
   - AI-powered reasoning for complex analysis
   - 30-60 seconds execution
   - Best for: Complex requirements, novel patterns

3. **Hybrid** (Smart Path - Recommended)
   - Code generation → AI validation → AI review
   - Regenerates with AI if quality < 0.9
   - 5-30 seconds execution
   - Best for: Production use, unknown complexity

**Quality Gates in Plan Mode**:

```
Step 1: Generate YAML (code or AI)
Step 2: Validate YAML
  ✓ Schema valid
  ✓ File references exist
  ✓ Dependencies resolvable
  ✓ Context budgets under limits
  ✓ Task sizes 30min-2.5hr

Step 3: AI Review Against Research Best Practices
  ✓ Has 2-3 style anchors per task
  ✓ Uses affirmative constraints
  ✓ Includes TDD workflow
  ✓ Multi-layer verification
  ✓ Context optimization applied
  ✓ Critical specs positioned correctly
  
Step 4: Generate Review Report
  Quality Score: 0.96
  Issues: 0 critical, 1 warning
  Recommendation: APPROVED
```

### Build Mode

**Purpose**: Execute validated YAML with Ralphy/OpenCode

**Execution Flow**:

```bash
$ your-tool build tasks.yaml

Pre-flight checks:
  ✓ YAML validation passed
  ✓ Git working tree clean
  ✓ All dependencies available
  
Executing with Ralphy:
  ✓ Spawned 3 parallel agents
  ✓ Agent 1: Task auth-001 (schemas)
  ✓ Agent 2: Task auth-003 (routes)
  ✓ Agent 3: Task auth-005 (middleware)
  
Post-execution:
  ✓ 8/8 tasks complete
  ✓ All verifications passed
  ✓ Learning from execution...
```

**Build Mode Features**:
- Pre-flight validation
- Parallel agent orchestration via Ralphy
- Real-time progress monitoring
- Post-execution learning
    - Automatic knowledge base updates (see `/docs/initial-claude-conversation/split-files/06-knowledge-management.md` for knowledge export/import, validation, and background learning strategies)
- Cost/time tracking

---

## Research-Backed Best Practices

*Based on findings from `/docs/initial-claude-conversation/split-files/02-research-report-preventing-drift.md (see /Users/kyledavis/Sites/prompt-stack/docs/best-practices.md)`*

**These must be embedded in every generated YAML:**

### 1. Style Anchors (Critical)
- **Requirement**: 2-3 reference files per task
- **Detection**: Auto-find via AST analysis + cached patterns
- **Format**: File path + line ranges + reason
- **Quality impact**: 40%+ improvement
- **Research basis**: Based on findings in `/docs/initial-claude-conversation/split-files/02-research-report-preventing-drift.md (see /Users/kyledavis/Sites/prompt-stack/docs/best-practices.md)` showing style anchors as the "single highest-leverage technique" for preventing AI drift

### 2. Task Sizing
- **Requirement**: 30 minutes to 2.5 hours maximum
- **Detection**: Estimate from file count, complexity, history
- **Action**: Auto-split tasks exceeding limit
- **Reason**: Prevents context overflow
- **Research basis**: Research in `/docs/initial-claude-conversation/split-files/02-research-report-preventing-drift.md (see /Users/kyledavis/Sites/prompt-stack/docs/best-practices.md)` shows tasks larger than 2.5 hours force AI models to "invent structure from generic patterns" rather than follow established conventions

### 3. Affirmative Constraints
- **Requirement**: State what TO do, not what NOT to do
- **Example**: ✅ "Use unknown with type guards" vs ❌ "Don't use any"
- **Quality impact**: 40%+ better compliance
- **Research basis**: Research in `/docs/initial-claude-conversation/split-files/02-research-report-preventing-drift.md (see /Users/kyledavis/Sites/prompt-stack/docs/best-practices.md)` demonstrates AI models "follow affirmative framing more reliably than prohibitions"

### 4. Multi-Layer Enforcement
```
Layer 1: Prompt (in YAML constraints)
Layer 2: IDE (LSP + ESLint config)
Layer 3: Commit (pre-commit hooks)
Layer 4: CI (quality gates)
Layer 5: Runtime (Zod schemas)
```
- **Research basis**: Research in `/docs/initial-claude-conversation/split-files/02-research-report-preventing-drift.md (see /Users/kyledavis/Sites/prompt-stack/docs/best-practices.md)` identifies multi-layer enforcement as critical because "prompt-level rules alone won't prevent drift"

### 5. TypeScript Strict Mode
```yaml
forbidden_patterns:
  - pattern: "\\bany\\b"
    message: "Use unknown with type guards"
  - pattern: "@ts-ignore"
    message: "Fix type errors properly"
  
required_patterns:
  - "import.*zod"  # When validating external data
```

### 6. TDD Workflow
```yaml
workflow: "test-first"
tdd_steps:
  1. "Write failing test"
  2. "Implement minimum to pass"
  3. "Verify all tests pass"
  4. "Refactor if needed"
```
- **Research basis**: Research in `/docs/initial-claude-conversation/split-files/02-research-report-preventing-drift.md (see /Users/kyledavis/Sites/prompt-stack/docs/best-practices.md)` shows TDD serves as "natural language specifications" that constrain AI generation more effectively than prose requirements

### 7. Context Positioning
- Critical specs: Beginning + end of prompts
- Never bury requirements in middle
- Reason: 20%+ performance drop for middle content
- **Research basis**: Research in `/docs/initial-claude-conversation/split-files/02-research-report-preventing-drift.md (see /Users/kyledavis/Sites/prompt-stack/docs/best-practices.md)` reveals the "lost in the middle" problem where AI model performance drops by 20%+ for information positioned in the middle of context

### 8. Context Budget Management
- Calculate tokens per task
- Trigger compaction at 85% full
- Default limit: 5,000 tokens per task
- Target reduction: 80-95% vs naive approach

### 9. Model-Specific Strategies
- **Claude**: Precision edits, surgical fixes (generation)
- **GPT**: Code review, bug catching (validation)
- **Strategy**: Use both for critical tasks
- **Research basis**: Research in `/docs/initial-claude-conversation/split-files/02-research-report-preventing-drift.md (see /Users/kyledavis/Sites/prompt-stack/docs/best-practices.md)` shows Claude models have "65% less likelihood to engage in shortcuts/loopholes" while GPT is better for exploratory work and code review

### 10. Self-Consistency
- Generate 3 solutions for critical tasks
- AI votes on best approach
- Reduces errors 30%+
- **Research basis**: Research in `/docs/initial-claude-conversation/split-files/02-research-report-preventing-drift.md (see /Users/kyledavis/Sites/prompt-stack/docs/best-practices.md)` identifies self-consistency checking as a technique that "dramatically reduces drift" by having AI models review their own outputs

---

## Architecture

### Component Separation

```
┌────────────────────────────────────────────────┐
│  YOUR TOOL (Meta-Orchestration)                │
│  • Plan Mode: YAML generation + review         │
│  • Build Mode: Execution orchestration         │
│  • Knowledge management                        │
│  • Context optimization                        │
└─────────────┬──────────────────────────────────┘
              │
              ├─→ [Plan Mode: Code Path]
              │     Fast template generation
              │
              ├─→ [Plan Mode: AI Path]
              │     Ralphy meta-loop for YAML
              │
              └─→ [Build Mode]
                    ↓
┌────────────────────────────────────────────────┐
│  RALPHY (Execution Layer)                      │
│  • Parallel agent orchestration                │
│  • Git workflow automation                     │
│  • Task execution                              │
└─────────────┬──────────────────────────────────┘
              ↓
┌────────────────────────────────────────────────┐
│  OPENCODE (Implementation Layer)               │
│  • Code generation                             │
│  • LSP integration                             │
│  • File editing                                │
└────────────────────────────────────────────────┘
```

### Data Flow: Plan Mode

```
Requirements
    ↓
[Method Selection: code/ai/hybrid]
    ↓
Code Path:                    AI Path:
  Template + Cache   OR       Meta-PRD Generation
  2-5 seconds                 ↓
    ↓                        Ralphy Execution
    ↓                        30-60 seconds
    ↓                         ↓
    └────────┬────────────────┘
             ↓
    [AI Validation]
    Quality score: 0.0-1.0
             ↓
    if score < 0.9 → regenerate with AI
             ↓
    [AI Review vs Research]
    Compliance check
             ↓
    Perfect YAML
```

### Data Flow: Build Mode

```
Validated YAML
    ↓
Pre-flight Checks
    ↓
Ralphy Execution
    ↓
OpenCode Implementation
    ↓
Verification Gates
    ↓
Learning & Feedback
    ↓
Knowledge Base Update
```

Note: For concrete OpenCode integration patterns (LSP setup, MCP servers, custom commands, project-level `OPENCODE.md`, agent configurations, and Husky/lint-staged examples) see `/docs/initial-claude-conversation/split-files/03-opencode-integrations-discussion.md`.

---

## JIT Discovery System

(See `/docs/initial-claude-conversation/split-files/05-jit-caching.md` for the full JIT discovery flow, progressive interrogation examples, and SQLite caching schema.)

### Progressive Knowledge Building

**First Run**:
```
❓ Questions asked: 3-5 (blocking only)
✓ Patterns cached: 5-10
⏱ Time: 2-3 minutes
```

**Second Run**:
```
❓ Questions asked: 0-1
✓ Patterns used: 5-10 (from cache)
⏱ Time: 5-10 seconds
```

**Team Knowledge Sharing**:
```bash
# Senior dev exports knowledge
$ your-tool knowledge export > .your-tool/team-knowledge.json
$ git commit -m "Add team patterns"

# Junior dev imports
$ your-tool knowledge import .your-tool/team-knowledge.json
✓ Immediate benefit from senior's knowledge
```

### Discovery Strategies

**1. Smart File Suggestion**
```javascript
// AI finds similar files based on:
- Task description keywords
- File naming patterns
- Import dependencies
- Recent modifications
- Usage frequency in successful tasks
```

**2. Confidence-Based Interrogation**
```
Confidence > 0.9: Use without asking
Confidence 0.7-0.9: Use, confirm if interactive
Confidence 0.5-0.7: Always ask
Confidence < 0.5: Don't suggest
```

**3. Learning from Execution**
```javascript
// After task completes:
- Update usage statistics
- Identify new patterns
- Refine confidence scores
- Cache successful approaches
```

---

## SQLite Schema (Minimal)

(See `/docs/initial-claude-conversation/split-files/05-jit-caching.md` for a fuller SQLite schema and JIT caching design.)

```sql
-- Core tables only (full schema in implementation)

CREATE TABLE codebase_knowledge (
    id INTEGER PRIMARY KEY,
    repo_path TEXT UNIQUE,
    primary_language TEXT
);

CREATE TABLE style_anchors (
    id INTEGER PRIMARY KEY,
    codebase_id INTEGER REFERENCES codebase_knowledge(id),
    file_path TEXT,
    category TEXT,  -- 'service', 'schema', 'test', 'api'
    pattern_summary TEXT,
    usage_count INTEGER DEFAULT 0,
    confidence REAL DEFAULT 1.0
);

CREATE TABLE coding_standards (
    id INTEGER PRIMARY KEY,
    codebase_id INTEGER REFERENCES codebase_knowledge(id),
    rule_type TEXT,  -- 'required', 'forbidden'
    rule TEXT,
    priority INTEGER
);

CREATE TABLE task_history (
    id INTEGER PRIMARY KEY,
    codebase_id INTEGER REFERENCES codebase_knowledge(id),
    task_description TEXT,
    success BOOLEAN,
    context_tokens INTEGER,
    execution_time_minutes INTEGER
);
```

---

## Enhanced YAML Structure

### Minimal Example (Code Generated)

```yaml
# Fast generation: 3 seconds

tasks:
  - id: "auth-001"
    title: "Create auth schemas"
    
    files_in_scope:
      - "src/schemas/auth.schema.ts"
    
    style_anchors:
      - file: "src/schemas/user.schema.ts"
        reason: "Zod pattern"
    
    verification:
      pre_commit:
        - "eslint --max-warnings=0 {files}"
```

### Complete Example (AI Generated)

```yaml
# AI generation: 45 seconds
# Includes full research compliance

metadata:
  generated_by: "your-tool v2.0"
  method: "ai-generation"
  quality_score: 0.96
  review_status: "APPROVED"

global_constraints:
  forbidden_patterns:
    - pattern: "\\bany\\b"
      message: "Use unknown with type guards"
  
  required_patterns:
    - pattern: "import.*zod"
      when: "external_data"

tasks:
  - id: "auth-001"
    title: "Create authentication Zod schemas"
    
    # SCOPE (90% context reduction)
    files_in_scope:
      - "src/schemas/auth.schema.ts"
    estimated_context_tokens: 1200
    
    # STYLE ANCHORS (research requirement)
    style_anchors:
      - file: "src/schemas/user.schema.ts"
        lines: [1, 45]
        reason: "Follow Zod pattern"
      - file: "src/types/user.types.ts"
        lines: [1, 20]
        reason: "Type extraction pattern"
    
    # CONSTRAINTS (affirmative)
    required_patterns:
      - "export const.*Schema = z\\.object"
      - "export type.* = z\\.infer"
    
    # TDD WORKFLOW
    workflow: "test-first"
    test_file: "src/schemas/auth.schema.test.ts"
    
    # MULTI-LAYER VERIFICATION
    verification:
      pre_commit:
        - "eslint src/schemas/auth.schema.ts --max-warnings=0"
        - "tsc --noEmit"
      post_commit:
        - "npm test -- auth.schema.test.ts"
    
    # ACCEPTANCE CRITERIA
    acceptance_criteria:
      - "LoginSchema validates email"
      - "100% test coverage"
      - "Zero ESLint warnings"
    
    estimated_duration_minutes: 45
    
    completion_signal: "<promise>COMPLETE</promise>"
```

---

## Non-AI Context Optimization

(See `/docs/initial-claude-conversation/split-files/04-cli-tool-building-discussion.md` for CLI-focused strategies and subagent architectures to optimize AI context windows.)

**Goal**: 90% token reduction before sending to AI

### Techniques

**1. AST-Based Dependency Analysis**
```javascript
// Build dependency graph
// Include only required imports
// Result: 47 files → 4 files
```

**2. Git Diff Context**
```javascript
// For modifications:
// Send diff + 5 lines context
// Not full file
// Result: 500 lines → 50 lines
```

**3. Symbol Table Compression**
```javascript
// Send signatures, not implementations
// Result: 200 lines → 20 lines of function signatures
```

**4. Smart File Filtering**
```javascript
// Use ripgrep + scoring
// Relevance = name_match + content_match + recency
// Top 10 files only
```

**5. Line Range Extraction**
```javascript
// Extract relevant sections only
// Not entire files
// Result: 80% reduction per file
```

### Context Budget Calculation

```javascript
function calculateBudget(task) {
  const base = 500;  // Task description
  const files = estimateFileTokens(task.files);
  const anchors = estimateAnchorTokens(task.anchors);
  
  const total = base + files + anchors;
  const budget = total * 1.2;  // 20% buffer
  
  if (budget > 5000) {
    warn("Task too large, recommend split");
  }
  
  return budget;
}
```

---

## CLI Interface

### Core Commands

```bash
# Initialize repo
$ your-tool init

# Plan Mode: Generate YAML
$ your-tool plan requirements.md [--method code|ai|hybrid]

# Plan Mode: Validate YAML
$ your-tool validate tasks.yaml

# Plan Mode: Review against research
$ your-tool review tasks.yaml

# Build Mode: Execute YAML
$ your-tool build tasks.yaml

# Knowledge Management
$ your-tool knowledge <list|export|import|search>

# Refinement
$ your-tool refine tasks.yaml --task <id>
```

### Detailed: Plan Mode

```bash
$ your-tool plan [requirements] [options]

Options:
  --method <code|ai|hybrid>    Generation method (default: hybrid)
  --file, -f <file>            Read requirements from file
  --output, -o <file>          Output YAML file (default: tasks.yaml)
  --interactive, -i            Ask questions during generation
  --auto                       No questions (use cache only)
  --review                     Include AI review step (default: true)
  --max-tasks <n>              Maximum tasks to generate
  --template <name>            Use template

Examples:
  # Hybrid mode (recommended)
  $ your-tool plan requirements.md
  
  # Fast code generation
  $ your-tool plan "Add auth" --method code
  
  # High-quality AI generation
  $ your-tool plan requirements.md --method ai --review
  
  # Automated batch processing
  $ your-tool plan batch.txt --auto --method code
```

### Detailed: Build Mode

```bash
$ your-tool build <yaml> [options]

Options:
  --parallel <n>        Number of parallel agents (default: 3)
  --dry-run            Simulate execution
  --continue           Continue from last checkpoint
  --watch              Monitor progress in real-time

Examples:
  # Standard execution
  $ your-tool build tasks.yaml
  
  # Dry run first
  $ your-tool build tasks.yaml --dry-run
  
  # Continue failed execution
  $ your-tool build tasks.yaml --continue
```

### AI Review Command

```bash
$ your-tool review <yaml> [options]

Reviews YAML against research best practices.

Options:
  --strict              Fail on warnings
  --fix                 Auto-fix issues where possible
  --report <format>     text|json|html

Output:
  Quality Score: 0.96/1.00
  
  ✓ Style anchors: 2-3 per task
  ✓ Task sizing: All tasks 30-150 min
  ✓ Affirmative constraints: 100%
  ✓ Multi-layer verification: Present
  ✓ Context optimization: 91% reduction
  ⚠ TDD workflow: Missing in 1 task
  
  Recommendation: APPROVED (1 warning)
```

---

## Plan Mode: AI Review Process

### Meta-PRD for Review

When `--review` is enabled, tool generates review PRD:

```yaml
# review-prd.yaml (executed by Ralphy)

tasks:
  - id: "review-style-anchors"
    title: "Verify 2-3 style anchors per task"
    description: |
      Check each task has 2-3 style anchor file references.
      Anchors must exist and be relevant to task type.
      Report any tasks missing anchors.
    
    files_in_scope:
      - "tasks.yaml"
    
    verification:
      post_commit:
        - "All tasks have 2-3 anchors"
    
    output_format: "JSON report"
  
  - id: "review-task-sizing"
    title: "Validate task duration estimates"
    description: |
      Check all tasks are 30-150 minutes.
      Flag any oversized tasks.
      Recommend splits if needed.
    
    verification:
      post_commit:
        - "All tasks within size limits"
  
  - id: "review-constraints"
    title: "Check affirmative constraint format"
    description: |
      Verify constraints use affirmative language.
      Flag any "don't" or "never" patterns.
      Suggest rewrites.
  
  - id: "review-context-optimization"
    title: "Validate context reduction"
    description: |
      Calculate expected context per task.
      Verify 80%+ reduction vs naive approach.
      Flag tasks exceeding 5000 token budget.
  
  - id: "generate-review-report"
    title: "Create final review report"
    depends_on: [
      "review-style-anchors",
      "review-task-sizing",
      "review-constraints",
      "review-context-optimization"
    ]
    description: |
      Aggregate all review results.
      Calculate quality score (0.0-1.0).
      Generate recommendations.
    
    output: "review-report.json"
    
    completion_signal: "<promise>REVIEW_COMPLETE</promise>"
```

### Review Execution

```bash
# In plan mode with --review:

$ your-tool plan requirements.md --method hybrid --review

⚡ Generating YAML (hybrid mode)...
✓ Generated candidate YAML (5 seconds)

🤖 AI validating quality...
✓ Validation passed (score: 0.92)

🔍 Reviewing against research best practices...
  Generating review PRD...
  Executing review with Ralphy...
  
  Review Task 1: Style anchors... ✓
  Review Task 2: Task sizing... ✓
  Review Task 3: Constraints... ⚠ 1 issue
  Review Task 4: Context optimization... ✓
  Review Task 5: Generate report... ✓

📊 Review Results:
  Quality Score: 0.96/1.00
  Status: APPROVED
  Issues: 0 critical, 1 warning
  
  Warnings:
    - Task auth-003: Missing TDD workflow steps
  
  Recommendations:
    - Add test-first workflow to auth-003
    - Otherwise ready for execution

Auto-fix warnings? [Y/n] Y

✓ Applied fixes
✓ Final score: 0.98

Saved to: tasks.yaml
Execute with build mode? [Y/n]
```

---

## Plan Mode: Code vs AI Generation

### Code Generation (Fast)

**When to use**:
- Similar to previous tasks
- Simple requirements
- Time-sensitive
- Batch processing

**How it works**:
```javascript
1. Parse requirements
2. Query knowledge DB for cached patterns
3. Load appropriate template
4. Fill template with:
   - Tasks from requirements
   - Style anchors from cache
   - Standards from cache
   - Verification commands from config
5. Generate YAML (2-5 seconds)
```

**Quality**: 85% (good for routine tasks)

### AI Generation (High Quality)

**When to use**:
- Complex requirements
- Novel patterns
- Production-critical
- Learning mode

**How it works**:
```javascript
1. Generate meta-PRD (YAML generation instructions)
2. Execute meta-PRD with Ralphy
   - Task 1: Analyze requirements
   - Task 2: Find style anchors (semantic search)
   - Task 3: Optimize context (AI reasoning)
   - Task 4: Apply standards
   - Task 5: Generate final YAML
3. Ralphy outputs tasks.yaml (30-60 seconds)
```

**Quality**: 95% (excellent for complex tasks)

### Hybrid (Recommended)

**How it works**:
```javascript
1. Generate with code (fast)
2. AI validates quality
3. If score < 0.9:
   - Regenerate with AI
   - Use AI feedback
4. AI reviews against research
5. Return best result
```

**Quality**: 90-95% with speed optimization

---

## Implementation Phases

### Phase 0: Requirements Gathering (1-3 days)

- Goal: Capture, formalize, and validate project requirements as the first milestone. Produce a concise requirements input file and templates that Plan Mode can consume; this milestone intentionally avoids Ralphy/OpenCode usage and focuses on human-driven requirements collection and static validation.
- Deliverables: `examples/requirements/inputs/requirements.input.md` (or `requirements.md`), `examples/requirements/templates/requirements-prompt.md` (prompt + templates used to gather requirements), an updated `docs/requirements/main.md` entry documenting the gathered requirements, and a short `planning/manifest.yaml` entry referencing the milestone.
- Acceptance criteria:
  - A requirements input file exists at `examples/requirements/inputs/requirements.input.md` and follows the project's input template.
  - The example prompt and templates used to gather requirements are committed under `docs/requirements/templates/`.
  - `your-tool plan examples/requirements/inputs/requirements.input.md --method code` produces a syntactically valid `tasks.yaml` candidate (code-generation path only; no AI required).
  - The requirements document (`docs/requirements/main.md`) includes the example prompt and links to the templates.
- Manual test checklist:
  1. Run the requirements prompt (copy the template from `examples/requirements/templates/requirements-prompt.md`) interactively and save output to `examples/requirements/inputs/requirements.input.md`.
  2. Run `your-tool plan examples/requirements/inputs/requirements.input.md --method code` and verify `tasks.yaml` is produced.
  3. Confirm `docs/requirements/main.md` references this milestone and contains the example prompt or link to `docs/requirements/templates/`.

### Phase 1: MVP (1-2 weeks)

**Plan Mode - Code Generation**:
- [ ] CLI scaffold
- [ ] Requirements parsing
- [ ] Template-based YAML generation
- [ ] Basic validation
- [ ] SQLite schema + basic caching

**Deliverable**: Working plan mode with code generation

### Phase 2: Context Optimization (2 weeks)

- [ ] AST parsing
- [ ] Dependency graph
- [ ] File filtering (ripgrep)
- [ ] Token estimation
- [ ] Context budget validation

**Deliverable**: 80%+ context reduction

### Phase 3: JIT Discovery (2 weeks)

- [ ] Interactive interrogation
- [ ] Pattern detection
- [ ] Confidence scoring
- [ ] Knowledge export/import
- [ ] Team sharing

**Deliverable**: Zero-question second runs

### Phase 4: AI Generation (2 weeks)

- [ ] Meta-PRD generation
- [ ] Ralphy integration for YAML generation
- [ ] Hybrid mode implementation
- [ ] Quality validation

**Deliverable**: High-quality AI generation option

### Phase 5: AI Review (1 week)

- [ ] Review PRD generation
- [ ] Research compliance checking
- [ ] Quality scoring
- [ ] Auto-fix suggestions

**Deliverable**: Self-validating plans

### Phase 6: Build Mode (1 week)

- [ ] Ralphy execution wrapper
- [ ] Progress monitoring
- [ ] Post-execution learning
- [ ] Knowledge updates

**Deliverable**: Complete plan → build workflow

---

## Integration Points

### With Ralphy

**Plan Mode**:
```bash
# Generate YAML (AI method)
your-tool plan → meta-PRD → ralphy.sh → tasks.yaml

# Review YAML
your-tool review → review-PRD → ralphy.sh → report.json
```

**Build Mode**:
```bash
# Execute YAML
your-tool build → ralphy.sh → code implementation
```

### With OpenCode

**Indirect** (via Ralphy):
```
Ralphy calls OpenCode agents
OpenCode follows YAML constraints
```

**Direct** (optional):
```bash
# Export custom commands
your-tool export-commands --format opencode
→ Creates .opencode/commands/*.md
```

See detailed OpenCode integration patterns, recommended configurations, and example command templates in `/docs/initial-claude-conversation/split-files/03-opencode-integrations-discussion.md` (covers LSP integration, MCP servers, custom commands, `OPENCODE.md` project config, agent roles, and Husky/lint-staged recommendations).  

### With Git

**Hooks Installation**:
```bash
your-tool init --install-hooks
→ Creates .husky/pre-commit with zero-warning checks
```

**Branch Management**:
```bash
# Via Ralphy in build mode
--branch-per-task flag creates branches automatically
```

### With CI/CD

```yaml
# .github/workflows/validate-plan.yml
- name: Validate YAML
  run: your-tool validate tasks.yaml --strict

- name: Review against research
  run: your-tool review tasks.yaml --strict
```

---

## Technical Specifications

### Stack

- **Runtime**: Node.js / Bun (for speed)
- **Language**: TypeScript
- **Database**: SQLite (better-sqlite3)
- **CLI**: Commander.js
- **Prompts**: Inquirer
- **Validation**: Zod
- **AST**: @typescript-eslint/parser
- **Search**: ripgrep-js
- **Git**: simple-git
- **AI**: Anthropic SDK (Claude)

### Project Structure
See `docs/requirements/project-structure.md` for the authoritative Go layout (flat tree, `cmd/` entrypoint, `internal/<domain>` packages, testing guardrails, runtime artifacts).

your-tool/
├── src/
│   ├── cli/
│   │   └── commands/
│   │       ├── plan.ts       # Plan mode
│   │       ├── build.ts      # Build mode
│   │       ├── review.ts     # AI review
│   │       └── knowledge.ts
│   ├── core/
│   │   ├── code-generator.ts     # Template-based
│   │   ├── ai-generator.ts       # Meta-PRD
│   │   ├── validator.ts          # YAML validation
│   │   └── reviewer.ts           # Research compliance
│   ├── context/
│   │   ├── ast-analyzer.ts
│   │   ├── file-filter.ts
│   │   └── budget-calc.ts
│   ├── knowledge/
│   │   ├── database.ts
│   │   ├── discovery.ts
│   │   └── patterns.ts
│   └── yaml/
│       ├── generator.ts
│       └── templates/
├── templates/
│   ├── meta-prd-generation.yaml   # For AI YAML generation
│   └── meta-prd-review.yaml       # For AI review
└── tests/
```

### Performance Targets

- **Code generation**: <5 seconds
- **AI generation**: <60 seconds
- **Validation**: <1 second
- **Review**: <45 seconds
- **File scanning**: <5 seconds (1000 files)
- **Knowledge query**: <100ms

---

## Multi-Milestone Projects

This project supports iterative planning by treating each milestone as a first-class planning unit. Recommended approach:

- Per-milestone planning files: store a single concise input per milestone (follow `templates/planning-phase.input.yaml`) and generate one Ralphy YAML per milestone. Keep each milestone focused so AI context budgets remain small.
- Use a small project manifest to declare milestone order and cross-milestone dependencies (the manifest is lightweight and only references per-milestone input files).
- Run Plan Mode per milestone during development (`your-tool plan <milestone.input.yaml>`) and run a separate integration pass when you need a full-project view or cross-milestone validation.
- Enforce per-milestone quality gates (YAML syntax, schema, secrets scan, style-anchor and sizing validations). Only promote a milestone to `APPROVED` once it meets the quality target.
- Store artifacts under `planning/` or `examples/multi-milestone-setup/` for clarity (see the concise example in `examples/multi-milestone-setup`).

Example layout (recommended):

- planning/
  - inputs/
    - <milestone>.input.yaml
  - milestones/
    - <milestone>.ralphy.yaml
  - reports/
    - <milestone>/final_quality_report.json
  - manifest.yaml  # lightweight list of milestones + deps

CLI examples:

- Generate a single milestone (hybrid + review):
  `your-tool plan planning/inputs/auth-v1.input.yaml --method hybrid --review --output planning/milestones/auth-v1.ralphy.yaml`

- Run a manifest-driven integration check (validate ordering and cross-milestone constraints):
  `your-tool plan --manifest planning/manifest.yaml --integration-check`

See `examples/multi-milestone-setup` for a concrete minimal example you can copy and adapt.

## Complete Workflow Examples

### Example 1: First Time (Full Learning)

```bash
$ cd my-project
$ your-tool init

✓ Created knowledge base
✓ Cached 15 patterns (3 minutes)

$ cat requirements.md
# Add user authentication with JWT

$ your-tool plan requirements.md

❓ Do you have auth code? > yes
📂 Use AuthService.ts as anchor? > yes
✓ Cached patterns

⚡ Generating YAML (hybrid)...
✓ Code generation (3s)
🤖 AI validation (2s) - Score: 0.94
🔍 AI review (45s) - Score: 0.97

✓ Saved to tasks.yaml

$ your-tool build tasks.yaml

Pre-flight: ✓
Executing with Ralphy...
  ✓ 8/8 tasks complete
  ✓ 3.2 hours, $8.30
  
✓ Learning from execution...
```

### Example 2: Second Time (Zero Questions)

```bash
$ your-tool plan "Add notification system"

✓ Using cached patterns (0 questions)
⚡ Generated (5s)
✓ Review score: 0.96

$ your-tool build tasks.yaml
  ✓ 6/6 tasks complete
```

### Example 3: AI Generation for Complex Task

```bash
$ your-tool plan complex-refactor.md --method ai

🤖 AI generation mode...
  Generating meta-PRD...
  Executing with Ralphy...
    ✓ Task 1: Analyze requirements
    ✓ Task 2: Find patterns
    ✓ Task 3: Optimize context
    ✓ Task 4: Generate YAML
    
✓ High-quality YAML (45s)
🔍 Review score: 0.98

$ your-tool build tasks.yaml
```

### Example 4: Failed Task Refinement

```bash
$ your-tool build tasks.yaml

Task 3 failed: Missing JWT config

$ your-tool refine tasks.yaml --task 3 \
  --add-context "Need JWT_SECRET in env"

✓ Updated task 3
✓ Added configuration setup subtask

$ your-tool build tasks.yaml --continue
  ✓ Task 3 (retry): Complete
```

---

## Key Principles

### 1. Quality Over Speed (But Optimize Both)

- Default to hybrid mode: fast when possible, AI when needed
- Always review against research best practices
- Never sacrifice quality for marginal speed gains

### 2. Knowledge Compounds

- Every execution teaches the system
- Team shares knowledge via export/import
- Second+ runs should be near-instant

### 3. AI as Validator AND Generator

- Use AI to check AI-generated content
- Meta-loop: AI improves AI inputs
- Self-validating, self-improving system

### 4. Non-AI First

- 90% of optimization is non-AI (parsing, filtering, caching)
- Reserve AI for reasoning, not mechanical tasks
- Cheaper, faster, more reliable

### 5. Fail Fast, Learn Fast

- Validation before execution
- Review before build
- Learn from failures immediately

---

## Success Metrics

**Plan Mode**:
- Quality score >0.9 for 95% of plans
- Context reduction >80%
- Second-run generation <10 seconds

**Build Mode**:
- First-pass success >90%
- Zero drift from plan
- Cost reduction >10x vs manual

**Overall**:
- Developer velocity: 3-5x
- Bug reduction: 50%+
- Team onboarding: <1 hour (vs days)

---

## Conclusion

This tool creates a **closed-loop quality system**:

1. **Plan Mode** generates perfect specifications
2. **AI Review** validates against research
3. **Build Mode** executes with guarantees
4. **Learning** improves future plans

**The innovation**: Plans that guarantee their own quality through multi-tier validation and research-backed constraints embedded in executable YAML.

**The result**: Predictable, efficient, high-quality AI-assisted development at scale.

---

## Commit & Revert Policy (Required)

- **Default branch-per-task**: Build Mode runs with `--branch-per-task` enabled by default; each task executes on its own branch to contain changes and make reviews simple.
- **Atomic commits per task**: After a task's agents complete and all verification gates pass for that task, the system creates a single atomic commit containing ONLY the files listed in the task's `files_in_scope`.
- **Verification before merge**: A task's branch must pass all post-commit verification (tests, ESLint with `--max-warnings=0`, type checks, acceptance criteria) before it is merged into the main integration branch.
- **Fail-fast and revert**: If a task fails any verification step during Build Mode, execution stops for dependent tasks. The tool will, by default, abort further execution and either:
  - leave the task branch for investigation, or
  - when `--fail-fast-and-revert` (or `--abort-on-drift`) is supplied, automatically revert the working tree to the pre-build state for that task and reset the branch to the last successful commit.
- **Scope enforcement**: Any file modifications outside a task's `files_in_scope` are flagged as a drift violation. The build aborts and the violating changes are not committed.
- **Command-line options**: Additions to Build Mode options:
  - `--commit-per-task` (boolean) — enforce atomic per-task commits (default: true)
  - `--fail-fast-and-revert` (boolean) — abort execution and revert working tree on verification failure (default: false)

Add this behavior to the Build Mode pre-flight checks and the example flows so builds are safe, auditable, and reversible.

---

## Nice-to-haves (Optional)

These recommendations are strongly encouraged but not required for MVP. They are helpful for long-term project hygiene and stricter enforcement.

1) Persistent project rules file (CLAUDE.md / .cursor/rules/) — nice to have

- Description: A repo-level rules file that your-tool loads at the start of Plan/Build sessions to provide persistent global constraints (style anchors, forbidden patterns, affirmative rules).
- Suggested example (to be created by `your-tool init` as an opt-in flag):

```yaml
# CLAUDE.md - Project rules loaded for every Plan/Build session
GLOBAL:
  - NEVER use "any"; prefer "unknown" + type guards
  - NEVER add inline eslint-disable or @ts-ignore without explicit justification
  - ALWAYS validate external input with Zod
  - ALWAYS explicitly type function params and returns
  - TASK_SIZING: 30-150m
  - STYLE_ANCHORS_PER_TASK: 2-3
```

2) CI check for inline ESLint/TS bypasses — nice to have

- Description: A concrete CI step (or reviewer rule) that counts `eslint-disable` comments and `@ts-ignore` occurrences and fails the run (or raises a high-severity warning) if any are present beyond an agreed threshold.
- Minimal example (convert to a lint rule or script for production):

```bash
# Fail if any eslint-disable or @ts-ignore occurrences in src/
if grep -R --line-number -E "//\s*eslint-disable|@ts-ignore" src | grep -q .; then
  echo "Found eslint-disable/@ts-ignore in src/ — fail build."
  exit 1
fi
```

These two items can be surfaced in `your-tool init --install-hooks` as optional opts and wired into `your-tool review` and CI templates.

---

*End of Requirements Document*
