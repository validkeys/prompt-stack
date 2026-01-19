# Requirements Interview Log

Purpose: Track Q&A while clarifying business requirements for the AI-Assisted Development Workflow Tool.

Started: 2026-01-18

---

Instructions: I will ask one question at a time and record your answers below. When you're ready, answer each question and I'll append it to this file.


Q1: Who are the primary users and stakeholders for this tool?
A1: senior dev

Notes: User provided a single role. We'll ask about secondary users next.

Q2: Are there secondary users (e.g., junior devs, QA, DevOps, PMs) who will also use or be impacted by the tool? If so, list them and their level of interaction (direct user, reviewer, CI consumer, etc.).
A2: this will be a standalone tool for developers to use.

Q3: What are the primary goals and success metrics for the tool beyond what's in the doc? For example: reduce PR cycle time by X, improve first-pass build success to Y, or lower manual YAML edits by Z. Provide numeric targets if you have them.
A3: i well i don't know that it matters right? This is simply an orchestration tool for developer as close to perfect implementation plans, and then using ralphy to implement those plans

Notes: User indicates the primary goal is orchestration and high-quality plan generation rather than tailoring to specific codebases. Default support should be broad; ask about specifics if needed.

Q4: Which codebases (languages, repo sizes, monorepo vs single repo) must the tool support initially? Any constraints like strict TypeScript, legacy JS, or polyglot repos?
A4: That's the one i just answered. The existing type of repository to support shouldn't matter i don't believe. This is a meta tool on top of any existing codebase that helpers developers create finely tuned best practice yaml implementation plans and then implement, should they choose, with ralphy

Notes: Tool should be language-agnostic; focus on providing anchors, context reduction, and verification workflows that adapt to repo contents. If language-specific rules are needed, implement as optional plugins/configs.

Q5: What are the minimum required integrations at launch? (e.g., Git, CI providers, Ralphy, OpenCode, Anthropic/Claude). Which are optional initially?
A5: suggest

Notes: User asked for a suggested minimal integration set for the MVP.

Suggested Minimal Integration Set for MVP:
- Required: Git (local + remote awareness), Ralphy (execution), Anthropic/Claude (AI review/generation), SQLite (local knowledge cache)
- Optional: OpenCode (direct implementation), CI providers (GitHub Actions/GitLab CI), Husky/pre-commit hooks, LSP integration

Rationale: Keep MVP focused on plan generation and review with the ability to execute via Ralphy; add CI and editor integrations later.

Q6: What access & security considerations are required? (e.g., store API keys, local vs remote execution, repo permissions, audit logs)
A6: Opencode, Ralphy

Notes: User clarified Q5: primary integrations at launch are OpenCode and Ralphy. Update suggested set if needed.

Q7: What access & security considerations are required? (e.g., store API keys, local vs remote execution, repo permissions, audit logs)
A7: recommend

Notes: User asked for recommendations on access and security defaults.

Recommended Access & Security Defaults:
- API keys: store in OS-provided secret stores or require env vars; avoid persisting in the knowledge DB. Provide `your-tool auth` to configure keys.
- Execution: default to local execution; optional remote mode with explicit opt-in and TLS + token auth.
- Repo permissions: read/write only when running Build Mode; require user confirmation before commits/branches. Use minimal scopes for CI tokens.
- Audit logs: local append-only log; optional remote telemetry with explicit opt-in. Record actions like plan generation, AI review results, build runs, and commits (timestamp, user, git-sha).
- Secrets in tasks: never include secrets in tasks.yaml; provide secure vault integration (optional) for runtime secrets during Build Mode.

Q8: What level of configurability do you expect for defaults like task sizing, anchors per task, or token budgets? (global config, per-repo file, per-plan flags)
A8: I think there should be sensible defaults built in with an ability to add a repo-root configuration file or pass a custom configuration file path via cli argument. For POC, sensible defaults.

Notes: Defaults + optional per-repo config and CLI override. POC: rely on sensible defaults only.

Q9: Any non-functional requirements (performance, latency, file scanning limits) beyond the document? For example: max files to scan, response time targets, or memory limits.
A9: propose

Notes: User asked for proposed non-functional requirements for MVP.

Proposed Non-Functional Requirements (MVP):
- File scanning: handle codebases up to 1,000 files with a full scan in <5s on developer machine (SSD).
- Code generation latency: code-generation path <5s; hybrid plan generation <30s.
- Validation latency: YAML validation <1s.
- Memory: keep working memory use modest (<512MB) on typical dev machines.
- Disk: local SQLite use only; database size should be small (<50MB initially).
- Concurrency: Build Mode supports 3 parallel agents by default.
- Reliability: graceful failure with clear error messages; do not auto-modify files without explicit consent.

Q10: Any required reporting or artifacts to produce from runs? (e.g., review-report.json, audit logs, task traces)
A10: propose

Proposed artifacts for MVP:
- `tasks.yaml` — final generated plan
- `review-report.json` — AI review results, quality score, issues/warnings
- `audit.log` — append-only log of actions (plan, review, build) with timestamps and git-sha
- `task-trace/` — per-task execution traces when Build Mode runs (optional)
- Human-readable summary printed to stdout and saved as `report.txt` (short)

Q11: Any constraints on user experience? (interactive CLI only, non-interactive batch, GUI later?)
A11: Interactive CLI only

Notes: Q9 & Q10 confirmed as propose. User wants interactive CLI only for POC.

Q12: Any compliance or legal constraints? (data residency, telemetry opt-in, IP retention rules)
A12: No

Notes: No specific compliance or legal constraints.

Q13: For Build Mode, do you want branch-per-task and commit-per-task enabled by default? (document suggests yes; confirm)
A13: Commit per task by default, not branch per task

Notes: Updated to reflect user's preference: `--commit-per-task` enabled by default; `--branch-per-task` not enabled by default.

Q14: Anything else I should ask to clarify business requirements before we move to architecture? (priority items, constraints, or open decisions)
A14: I don't think so

Notes: No further clarifications requested; ready to move to architecture.
