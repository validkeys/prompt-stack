# Architecture & Interview (Working Draft)

Source: docs/requirements/main.md

This file is a working architecture plan derived from the requirements document at `docs/requirements/main.md`. It aims to define a pluggable architecture (extension points, plugin lifecycle, adapters) and will be refined via a one-question-at-a-time interview. Each answer will be recorded below and used to evolve the architecture.

--

Core intent
- Provide a pluggable, testable, and auditable meta-orchestration CLI for Plan Mode (YAML generation + review) and Build Mode (execution orchestration).
- Keep non-AI optimizations first (parsing, filtering, caching) and provide clear adapter boundaries for AI providers, execution engines (Ralphy/OpenCode), git integrations, and storage backends.

High-level components (pluggable)
- cli: Command definitions and top-level orchestration (plan, review, validate, build, knowledge)
- core: Pure business logic (generators, validators, reviewers, orchestration primitives)
- context: Discovery & knowledge ingestion (MCP servers, document prompts, user interrogation)
- knowledge: SQLite-backed knowledge DB, import/export, pattern store
- adapters/plugins: Concrete implementations for external systems (AI providers, Ralphy, OpenCode, Git, CI, secret stores)
- plugin-host: Plugin discovery, registration, lifecycle (enable/disable, versioning)
- yaml: Generator + templates
- telemetry/audit: Append-only audit.log and optional opt-in telemetry

Pluggability goals and patterns
- Clearly typed adapter interfaces (Go: interfaces) for each external integration:
  - AIProvider (generate, validate, review)
  - Executor (Ralphy wrapper / local executor)
  - CodeImplementer (OpenCode adapter)
  - GitProvider (commit, branch, hooks)
  - StorageProvider (SQLite, optional remote cache)
  - PromptStore (templates, meta-PRDs)
- Plugin discovery: filesystem-based (`./plugins` or `.prompt-stack/plugins`), plus executable manifest files for discoverability
- Plugin lifecycle: install -> register -> enable -> health-check -> run -> uninstall
- Versioning: plugin semver + declared host-api-minimum-version
- Security: run plugins with least privilege; require explicit consent for network or repo-mutating operations; prefer signed/trusted plugins for remote installs

Extension points (examples)
- Generation strategy (code | ai | hybrid) — plugins can provide strategies
- AI validation & review hooks — allow multi-model pipelines to be composed
- Discovery/ingestion hooks — plugins provide connectors to MCP servers, document stores, or interactive question flows
- Knowledge import/export formats — plugin-provided converters (e.g., org-specific patterns)
- Build-time hooks — pre-task / post-task hooks, verifier plugins

Plugin API (sketch)
- Every plugin ships a small manifest: { id, name, description, version, hostApiVersion, provides: ["ai-provider"], configSchema }
- Plugins are separate executables that communicate with the host via a simple RPC protocol (stdin/stdout JSON-RPC or protobuf) — host provides services (logger, db access via API surface, eventBus)
- Host exposes eventBus events: beforePlan, afterPlan, beforeTask, taskSucceeded, taskFailed, beforeCommit, afterCommit

CLI & config
- Project config: repo-root `.prompt-stack/config.yaml` (opt-in from `prompt-stack init`); allows enabling/disabling plugins and setting provider preferences
- Global defaults: environment variables, OS secret stores for API keys (never stored in DB)
- Interactive mode: one-question-at-a-time UX in CLI; support `--auto` to skip questions and `--interactive` to force them

Persistence & data model
- SQLite as primary local store; provide StorageProvider interface so a plugin can replace it with a remote DB
- Key tables: codebase_knowledge, style_anchors, coding_standards, task_history, plugin_registry, audit_log
- Knowledge export/import JSON for team sharing

Security
- Enforce not storing secrets in DB; use OS secret store or env vars
- TLS + token auth for optional remote mode
- Audit.log append-only by default; signed entries where possible
- Require user confirmation for branch/commit creation unless configured

Observability & performance
- Concurrency defaults: 3 parallel agents, configurable
- Performance targets from requirements: file scanning <5s for 1k files, code gen <5s, AI gen <60s, validation <1s
- Instrumentation hooks for timing and cost estimation plugins

Implementation notes (Go-first)
- Runtime: Go
- Language: Go 1.20+ (or latest stable supported)
- Database: SQLite via modern Go driver (e.g. mattn/go-sqlite3 or modernc/sqlite)
- CLI: Cobra (commands/subcommands)
- Prompts: survey or promptui for interactive one-question UX
- Validation: lightweight schema validation in Go (e.g. go-playground/validator) and YAML schema verification via go-yaml + z-schema equivalents
- Plugin model: external executables (RPC) — preferred for cross-platform, version-robust extensibility
- AI integrations: adapters call HTTP SDKs/CLIs (Anthropic, OpenAI, custom MCP clients) — keep keys in OS secret stores

Context & Discovery (changed)
- AST-based analysis is deferred/removed. Instead, the tool will rely on external discovery: MCP servers, document ingestion, and direct interactive questions to populate the knowledge graph.
- Discovery flows:
  - MCP connector: ask MCP for relevant style anchors / file suggestions
  - Document ingestion: accept user-provided documents, patterns, or exported knowledge JSON
  - Interactive interrogation: one-question-at-a-time follow-ups when confidence is low

Next steps: interview-driven decisions
- We'll use an interactive one-question-at-a-time approach to choose defaults and refine the architecture. Each answer will be appended to this file under an "Interview" section, and used to update component decisions.

Interview (recorded answers)
- Q1 (implementation language): Chosen: Go
  - Rationale: single static binary distribution makes it easy for anyone to download and run from different directories; good performance and concurrency; fits a design where most heavy work is delegated to external services (Ralphy/OpenCode/MCPs).
  - Long-term tradeoffs vs TypeScript: Go trades off some ecosystem conveniences (TypeScript's native AST tooling, ESLint/tsc integration, dynamic npm plugin model) for distribution simplicity, performance, and operational robustness. For TypeScript-specific analysis (TS AST), the tool will delegate to external Node-based services or MCP servers rather than implementing TypeScript parsing natively in Go.

- Q2 (plugin approach): Chosen: 1) External plugins
  - Rationale: cross-platform, version-robust extensibility; keeps the host binary stable while allowing org-specific integrations.
  - MVP note: no plugin ecosystem work required for MVP beyond reserving the interfaces and a minimal discovery mechanism (can be feature-flagged).

- Q3 (knowledge DB scope): Chosen: 1) Per-repo by default
  - Default path: `./.prompt-stack/knowledge.db`
  - Rationale: portable with repo; easy mental model; enables team sharing by committing/exporting selected knowledge artifacts (not secrets).
  - Note: still support `--global-cache` (or config) later if you want cross-repo warm caches.

- Q4 (config + runtime artifacts): Chosen: 1) Per-repo `.prompt-stack/`
  - Default paths:
    - Config: `./.prompt-stack/config.yaml`
    - Knowledge DB: `./.prompt-stack/knowledge.db`
    - Audit log: `./.prompt-stack/audit.log`
    - Reports: `./.prompt-stack/reports/` (e.g. `review-report.json`, `report.txt`)
    - Task traces: `./.prompt-stack/task-trace/` (optional)
  - Rationale: everything needed to understand/reproduce a run lives with the repo; supports working from different directories without relying on global state.

- Q5 (MCP integration posture): Chosen: 1) Optional connectors
  - Rationale: keeps the tool usable in any repo with zero setup; MCP becomes an accelerant for discovery/knowledge ingestion when present.

- Q6 (Build Mode git workflow default): Chosen: 2) `--branch-per-task` default OFF
  - Rationale: keep the default workflow lightweight; allow teams to opt into branch-per-task when desired.

- Q7 (Build Mode commit behavior): Chosen: 1) `--commit-per-task` default ON
  - Rationale: aligns with auditability and reproducibility; commits become the primary unit of progress and rollback.

- Q8 (Plan Mode default generation method): Chosen: 2) AI-first
  - Rationale: the tool's value is context injection + best-practices enforcement; templates exist but are secondary.
  - Behavior:
    - If an AI provider is configured (direct LLM or via MCP), Plan Mode uses it by default.
    - If not configured, fallback to a minimal template-based generator with reduced guarantees.
  - Note: both Plan and Build ultimately rely on the Ralphy protocol for execution/orchestration.

- Q9 (Ralphy coupling): Chosen: 1) Hard dependency
  - Rationale: Ralphy is the execution substrate; this tool focuses on producing high-quality Ralphy plans (YAML) and invoking them consistently.
  - Implication: a working `ralphy` installation (or bundled Ralphy) is required for both Plan and Build flows.

- Q10 (Ralphy acquisition): Chosen: 3) Bundle Ralphy with this tool
  - Feasibility: yes; Ralphy is currently a shell script (`ralphy.sh`) so bundling means vendoring the script (and any required assets/config templates) into this repo and making the Go binary able to materialize + execute it.
  - Pinning: because upstream has no releases, bundle a pinned commit hash and provide an explicit upgrade command later (e.g. `prompt-stack ralphy upgrade`).
  - Platform note: the bundled Ralphy remains a bash toolchain, so it requires a POSIX shell and dependencies like `jq` (and optionally `yq`, `gh`, etc.).

- Q11 (bundling mechanism): Chosen: bundle `ralphy.sh`, but keep an abstraction boundary
  - Decision: bundle upstream `ralphy.sh` initially.
  - Constraint: all interaction with Ralphy goes through an internal `Executor` interface so we can swap in a future native Go implementation with minimal changes.

Executor abstraction (required)
- Define a stable internal contract for execution:
  - `Executor.Validate(ctx, planPath) -> ValidationResult`
  - `Executor.Plan(ctx, requirementsPath | stdin, options) -> PlanResult`
  - `Executor.Build(ctx, planPath, options) -> BuildResult`
  - `Executor.Capabilities() -> { supportsParallel, supportsBranchPerTask, supportsCommitPerTask, ... }`
- Implementations:
  - `RalphyShellExecutor` (MVP): embeds/extracts bundled `ralphy.sh` and calls it via `os/exec`.
  - `RalphyGoExecutor` (future): native Go port that honors the same interface.

Bundling approach (MVP)
- Use Go `embed` to ship `vendor/ralphy/ralphy.sh` in the binary.
- On run, extract to `./.prompt-stack/vendor/ralphy/ralphy.sh`, ensure executable, then invoke it.
- Capture stdout/stderr for audit + `report.txt` and parse structured signals when available.

- Q12 (shell dependency posture): Chosen: 1) macOS + Linux only for MVP
  - Rationale: bundled `ralphy.sh` + POSIX dependencies are first-class on macOS/Linux; Windows support can be added later via WSL guidance or a native Go executor.

- Q13 (engine selection default): Chosen: 2) OpenCode (`--opencode`)
  - Rationale: prefer OpenCode as the default execution engine for higher-quality code generation and closer integration with our implementation layer.

- Q14 (AI provider integration): Chosen: 1) Direct SDKs + MCP optional
  - Rationale: Architect for direct OpenCode delegation today while allowing MCP connectors as optional adapters later for advanced context injection.
  - MVP behavior: the tool delegates implementations to OpenCode (via Ralphy/OpenCode flags). We will expose adapter interfaces so an MCP connector (or other LLM provider SDKs) can be added later without changing core planning/build logic.

- Note: Current reality — Plan/Build delegates primarily to OpenCode via Ralphy. MCP connectors are supported conceptually and can be added as plugin executables or as remote adapters that implement the AIProvider interface.