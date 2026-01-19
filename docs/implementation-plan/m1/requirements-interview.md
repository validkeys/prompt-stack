# Requirements Interview - Milestone m1

**Interview conducted**: 2026-01-19

## Initial Context Gathering

Q: Give me a short project summary and links to any existing requirements, design docs, or repos. List 2–3 files or URLs that best show the code/design style we should follow.

A: Read @docs/project-summary.md and milestone 1 from @docs/requirements/milestones.md. Check all docs in @docs/index.md

Captured: Building "prompt-stack" - a Go CLI tool for AI-assisted development workflow with Plan Mode (YAML generation) and Build Mode (execution).

## Before You Begin (Confirm References)

Q: Which reference documents should I consult for this interview? Please provide any local paths or repo URLs. What is the milestone identifier/slug to use (e.g., `m1`)? Where should outputs be stored? Default: `docs/implementation-plan/<milestone>/`.

A: Use @docs/index.md templates/ examples/. 2. m1 3. yes

Captured: Will consult docs from `docs/index.md`, templates from `templates/`, and examples from `examples/`. Milestone identifier is `m1`. Outputs will be stored in `docs/implementation-plan/m1/`.

## Interview Questions

Q: What is the milestone id/slug you want to use for this work?

A: m1

Q: Give me a one-line title for the milestone.

A: M2 — Repo Init & CLI scaffold (1-2 days)

Q: Provide a short (1-2 sentence) description of the milestone goal.

A: asic repo + CLI surface implemented; commands wired but minimal behavior

Q: Who is the primary stakeholder or product owner for this milestone? (capture contact)

A: @kyledavis Much of this work may already be done from milestone 0

Q: What are the top objectives for this milestone? (3 max)

A: the help and init commands as well as a verification of our @docs/requirements/project-structure.md

Q: How will we measure success? Give metrics and targets.

A: - Goal: Basic repo + CLI surface implemented; commands wired but minimal behavior.
- Deliverables: `prompt-stack` CLI scaffold, `init`, `plan`, `validate`, `review`, `build` commands (stubs), project layout and `./.prompt-stack/` default structure.
- Acceptance criteria:
  - `prompt-stack --help` lists core commands.
  - `prompt-stack init` creates `./.prompt-stack/config.yaml` and `./.prompt-stack/knowledge.db` (empty) and prints instructions.
- Manual test checklist:
  1. Run `prompt-stack --help` and inspect listed commands.
  2. Run `prompt-stack init` in a sample repo; verify files created under `./.prompt-stack/`.
  3. Confirm no secrets are written to the DB.

Q: What is the canonical project name and current version we should record?

A: prompt-stack v1

Q: Which documents or schema files must planners reference (paths/URLs)? (validation assets)

A: Read @docs/index.md and recommend

Captured: Must reference: `docs/requirements/project-structure.md`, `docs/ralphy-inputs.schema.json`, `docs/requirements/main.md`, `docs/requirements/milestones.md`.

Q: Which runtime resources/tools does the plan assume (knowledge DB path, validator scripts)? (execution resources)

A: it currently assumes nothing

Q: Which files or code areas should be used as style anchors? Provide 2–3 file paths and a short reason for each.

A: /Users/kyledavis/Sites/prompt-stack/examples/style-anchor and @docs/style-markers.md

Captured: Style anchors: 1) `examples/style-anchor/` (Go CLI project template), 2) `docs/style-markers.md` (style guidance), 3) `docs/requirements/project-structure.md` (project layout spec).

Q: What should be included in scope for this milestone?

A: Goal: Basic repo + CLI surface implemented; commands wired but minimal behavior.
- Deliverables: `prompt-stack` CLI scaffold, `init`, `plan`, `validate`, `review`, `build` commands (stubs), project layout and `./.prompt-stack/` default structure.
- Acceptance criteria:
  - `prompt-stack --help` lists core commands.
  - `prompt-stack init` creates `./.prompt-stack/config.yaml` and `./.prompt-stack/knowledge.db` (empty) and prints instructions.
- Manual test checklist:
  1. Run `prompt-stack --help` and inspect listed commands.
  2. Run `prompt-stack init` in a sample repo; verify files created under `./.prompt-stack/`.
  3. Confirm no secrets are written to the DB.

Q: What is explicitly out of scope?

A: anything not explicitly stated here

Q: Any critical constraints or assumptions? (security, infra, timelines)

A: no

Q: What deliverables do you expect? For each deliverable capture: name, one-line description, owner (name/email), format (yaml/json/docx), and desired due date. Example: {name: task_breakdown.yaml, description: task list, owner: alice@example.com, format: yaml, due: 2026-02-01}

A: - Acceptance criteria:
  - `prompt-stack --help` lists core commands.
  - `prompt-stack init` creates `./.prompt-stack/config.yaml` and `./.prompt-stack/knowledge.db`

Captured: Deliverables: 1) `prompt-stack` CLI binary (working help command), 2) `init` command implementation, 3) `.prompt-stack/config.yaml`, 4) `.prompt-stack/knowledge.db`.

Q: Do you have desired timelines or dates? (start / target completion)

A: asap

Q: Testing requirements: unit tests, integration, TDD preference?

A: all should be tdd, prefer integration tests, unit where aboslutely necessary

Q: Any privacy, compliance, or secrets handling notes?

A: no

Q: Are there quality/acceptance thresholds we must meet (e.g., quality score >= 0.95)? For each threshold, also describe how it will be validated (automated test, manual review, smoke test).

A: Review /Users/kyledavis/Sites/prompt-stack/docs/implementation-plan/m0/make-implementation-plan.prd.yaml and use the quality/acceptance thresholds from that document

Captured: Quality threshold: quality score ≥0.95. Validation: YAML syntax validation, JSON Schema validation, secrets scan, style anchors compliance, task sizing compliance, affirmative constraints validation, multi-layer enforcement validation.