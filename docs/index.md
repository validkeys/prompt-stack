# Docs index

- Requirements
  - `docs/requirements/main.md` — Full requirements for the AI-Assisted Development Workflow Tool (Plan Mode, Build Mode, QA gates).
  - `docs/requirements/index.md` — Requirements entry: quick lookup referencing all requirements documents.
  - `docs/requirements/architecture.md` — Architecture, component separation, and data flows.
  - `docs/requirements/project-structure.md` — Flat Go project layout spec: cmd entrypoint, internal domain packages, testing guardrails.
  - `docs/requirements/milestones.md` — Project milestones and delivery phases.
  - `docs/requirements/interview.md` — Stakeholder interview notes and action items.
  - `docs/requirements/architecture.interview.md` — Architecture interview notes and follow-ups.

- Policies & Best Practices
  - `docs/best-practices.md` — Research-backed best practices (style anchors, task sizing, affirmative constraints).
  - `docs/style-markers.md` — Style markers, tokens and anchors for consistent prompts and docs.
  - `docs/tdd-checklist.md` — TDD checklist for AI-assisted edits: failing-tests-first loop, run/fix instructions, CI expectations.
  - `docs/task-sizing.md` — Task sizing guidance: keep model-driven tasks between 30 minutes and 2.5 hours; splitting guidance.
  - `docs/drift-policy.md` — Drift policy: stop/revert criteria, incident note guidance, and remediation workflow.
  - `docs/opencode-rules.md` — Repository-level Opencode rules and hard constraints loaded at session start (TypeScript rules, allowed stack, testing and CI requirements).

- Ralphy & Agent Config
  - `docs/ralphy-yaml-spec.md` — Ralphy YAML spec, examples, and validation rules.
  - `docs/ralphy-inputs.md` — Ralphy inputs checklist: required top-level fields the YAML generator must consume and an example.
  - `docs/ralphy-inputs.schema.json` — JSON Schema to validate Ralphy inputs (used to validate generated YAML before applying).
  - `docs/prompt-stack-config.md` — Repo-level role-based model selection design, examples, and integration notes.

- Origins & Research Conversation
  - `docs/initial-claude-conversation/split-files/01-initial-problem-statement.md` — Origin and initial problem statement.
  - `docs/initial-claude-conversation/split-files/02-research-report-preventing-drift.md` — Research on preventing AI drift and quality evidence.
  - `docs/initial-claude-conversation/split-files/03-opencode-integrations-discussion.md` — OpenCode integration patterns and agent examples.
  - `docs/initial-claude-conversation/split-files/04-cli-tool-building-discussion.md` — CLI design guidance and implementation notes.
  - `docs/initial-claude-conversation/split-files/05-jit-caching.md` — JIT discovery, caching, and SQLite schema.
  - `docs/initial-claude-conversation/split-files/06-knowledge-management.md` — Knowledge management commands, exports, and team-sharing examples.
  - `docs/initial-claude-conversation/chat.md` — Full initial chat transcript and context.

- Reports & Misc
  - `docs/report.md` — Research report and summaries.
  - `docs/tmp/from-idea-to-plan.md` — Drafts and early notes: idea→plan flow (temporary/proposal content).

- Templates
  - `templates/meta-prd-generation.yaml` — Meta-PRD template for AI YAML generation.
  - `templates/meta-prd-review.yaml` — Meta-PRD template for AI review generation.
  - `templates/requirements-gathering-prompt.md` — Prompt template for gathering requirements.
  - `examples/requirements/templates/requirements-prompt.md` — Example requirements prompt template (points to canonical).

- Index
  - `docs/index.md` — Document index (this file): concise, AI-friendly lookup of core docs.

