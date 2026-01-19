Architecture — concise spec

Purpose
- A compact, actionable architecture spec for the host CLI that produces and executes Ralphy plans and integrates with OpenCode. This file distills the working draft at `docs/requirements/architecture.interview.md` into a concise reference for implementers and reviewers.

Core intent
- Pluggable, testable, auditable meta-orchestration CLI with two primary modes: Plan (YAML generation + review) and Build (execution orchestration).
- Optimize non-AI concerns first (parsing, filtering, caching) and expose clear adapter boundaries for AI, executors, git, and storage.
- Plan and Build rely on the Ralphy protocol; a working Ralphy (bundled or installed) is required.

Primary components
- cli: Command surface and orchestration (plan, review, validate, build, knowledge).
- core: Pure business logic (generators, validators, reviewers, orchestration primitives).
- context/knowledge: Discovery, ingestion, and SQLite-backed knowledge DB.
- adapters/plugins: External integrations (AI providers, Ralphy/OpenCode, Git, CI, secrets, storage).
- plugin-host: Discovery, manifest, lifecycle, and eventBus for hook points.
- yaml/telemetry: Generator templates, audit.log, reports.

Adapter contracts (examples)
- AIProvider: generate, validate, review
- Executor: Validate, Plan, Build, Capabilities
  - Capabilities detail: supportsParallel, supportsBranchPerTask, supportsCommitPerTask, supportsValidationOnly
- CodeImplementer (OpenCode), GitProvider, StorageProvider, PromptStore

Plugin model & discovery
- Filesystem-discoverable plugins (e.g. `./plugins`, `./.prompt-stack/plugins`) and executable manifests for discoverability.
- Minimal manifest fields: `{ id, name, description, version, hostApiVersion, provides, configSchema }`.
- Protocols: plugins are separate executables communicating over stdin/stdout RPC (JSON-RPC or protobuf encouraged) — host provides logger/db/eventBus handles.
- Discovery/lifecycle: install → register → enable → health-check → run → uninstall.
- Versioning: plugin semantic version + declared `hostApiMinimumVersion`.
- Security: run plugins with least privilege, require explicit consent for network or repo-mutating operations, and prefer signed/trusted plugins for remote installs.

Extension points
- Generation strategy (code | ai | hybrid)
- AI validation & review hooks (multi-model pipelines)
- Discovery/ingestion connectors (MCP, docs, interactive)
- Knowledge import/export converters (org-specific formats)
- Build-time hooks (pre-task / post-task / verifier plugins)

Host eventBus (hook points)
- Example events the host emits: `beforePlan`, `afterPlan`, `beforeTask`, `taskSucceeded`, `taskFailed`, `beforeCommit`, `afterCommit`.

CLI & config
- Project config: repo-root `./.prompt-stack/config.yaml` (opt-in from `prompt-stack init`); allows enabling/disabling plugins and setting provider preferences.
- Global defaults: environment variables and OS secret stores for API keys (never stored in DB).
- Interactive UX: one-question-at-a-time with `--auto` (skip questions) and `--interactive` (force questions) flags.

Persistence & defaults
- Default per-repo layout under `./.prompt-stack/`:
  - `config.yaml`, `knowledge.db`, `audit.log`, `reports/`, `task-trace/`, `vendor/ralphy/`.
- SQLite primary store; `StorageProvider` interface allows remote replacement.
- Key tables (example): `codebase_knowledge`, `style_anchors`, `coding_standards`, `task_history`, `plugin_registry`, `audit_log`.
- Knowledge export/import: JSON for team sharing; never include secrets in exported artifacts.

Security & observability
- Do not store secrets in DB; use OS secret stores or environment variables.
- Require explicit user consent for repo/network mutations; audit.log is append-only by default and should support signed entries where possible.
- Optional remote mode: TLS + token auth.
- Concurrency defaults: 3 parallel agents (configurable).
- Performance targets (from working draft): file scanning <5s for 1k files, code gen <5s, AI gen <60s, validation <1s.
- Instrumentation hooks for timing, cost estimation, and telemetry (opt-in).

Executor strategy (MVP)
- Define `Executor` interface with methods like:
  - `Validate(ctx, planPath) -> ValidationResult`
  - `Plan(ctx, requirementsPath|stdin, options) -> PlanResult`
  - `Build(ctx, planPath, options) -> BuildResult`
  - `Capabilities() -> { ... }`
- Capabilities should include at minimum: `supportsParallel`, `supportsBranchPerTask`, `supportsCommitPerTask`, `supportsValidationOnly`.
- Implementations:
  - `RalphyShellExecutor` (MVP): embed `vendor/ralphy/ralphy.sh` via Go `embed`, extract to `./.prompt-stack/vendor/ralphy/ralphy.sh`, ensure executable, invoke via `os/exec`, capture stdout/stderr for audit and parsing structured signals when available.
  - `RalphyGoExecutor` (future): native Go port honoring same interface for Windows and improved portability.

Bundling approach (MVP)
- Use Go `embed` to ship `vendor/ralphy/ralphy.sh` in the binary.
- On run, materialize to `./.prompt-stack/vendor/ralphy/ralphy.sh`, make executable, and invoke.
- Pin the bundled Ralphy to a commit hash and provide an upgrade command (e.g. `prompt-stack ralphy upgrade`) to update the vendored script.
- Capture and persist stdout/stderr to `report.txt` and `audit.log`; parse structured outputs when the script emits machine-readable signals.
- Document POSIX dependency expectations (e.g., `bash`, `jq`, optional `yq`, `gh`).

Implementation notes (Go-first)
- Language/runtime: Go (Go 1.20+ recommended).
- Database: SQLite via a modern Go driver (e.g. `mattn/go-sqlite3` or `modernc/sqlite`).
- CLI: Cobra for commands/subcommands.
- Interactive prompts: `survey` or `promptui` for one-question UX.
- Validation: `go-playground/validator` and YAML validation via `gopkg.in/yaml.v3` plus schema tooling as needed.
- Plugin model: external executables (stdin/stdout RPC) for cross-platform, version-robust extensibility.
- AI integrations: adapters call HTTP SDKs/CLIs (Anthropic, OpenAI, OpenCode, custom MCP clients); keep keys in OS secret stores.

Context & discovery
- AST-based code analysis is deferred/removed from MVP; rely on external discovery: MCP servers, document ingestion, and interactive questioning to populate the knowledge graph.
- Discovery flows: MCP connector, document ingestion, and interactive interrogation when confidence is low.

Key implementation choices (from interview)
1) Language: Go (single static binary, concurrency). 2) Plugins: external executables. 3) Knowledge DB: per-repo `./.prompt-stack/knowledge.db`. 4) Per-repo runtime artifacts in `./.prompt-stack/`. 5) MCP: optional connectors. 6) Branch-per-task: default OFF. 7) Commit-per-task: default ON. 8) Plan Mode: AI-first with template fallback. 9) Ralphy: hard dependency, bundled for MVP. 10) Bundle `ralphy.sh` and expose Executor interface. 11) MVP platforms: macOS + Linux. 12) Default engine: OpenCode (`--opencode`). 13) AI integrations: direct SDKs; MCP optional.

Where this came from
- See `docs/requirements/architecture.interview.md` for the full working draft and the recorded interview rationale.

Next steps
1) Review and approve this concise spec or request changes (reply with specific areas to expand).
2) If approved: create an initial Go module, add `Executor` and `AIProvider` interface sketches, and scaffold `./.prompt-stack/` layout.
3) Optional: commit `docs/requirements/architecture.md` and link it from `docs/requirements/main.md`.
