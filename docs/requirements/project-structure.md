# Project Structure Spec — AI-Orchestration CLI (Go)

## Purpose
- Provide a concise, enforceable layout for the Go implementation of the Plan/Build orchestration tool described in `docs/requirements/main.md`.
- Favor a flat, discoverable tree with clear domain boundaries so new contributors can map requirements to code quickly.
- Align with the style markers captured in `docs/style-markers.md`: single-binary CLI (`cmd/`), modular internal packages, explicit internal vs external boundaries, and auditable build tooling.

## High-Level Layout
```
.
├── cmd/
│   └── orchestrator/
│       └── main.go
├── internal/
│   ├── ai/
│   ├── audit/
│   ├── build/
│   ├── cli/
│   ├── config/
│   ├── executor/
│   ├── knowledge/
│   ├── plan/
│   ├── plugins/
│   ├── telemetry/
│   ├── validation/
│   └── shared/
├── configs/
├── docs/
├── examples/
├── scripts/
├── Makefile
├── go.mod
└── tools.go
```
- Keep the root shallow: command entrypoints, internal domains, docs, examples, build tooling.
- Prefer `internal/<domain>` packages over deep nesting; only introduce subpackages for cohesion (e.g., `internal/executor/ralphy`).
- Avoid `pkg/` unless we intentionally export APIs; today all code is internal-only per requirements.

## Directory Specs
- `cmd/orchestrator`: Cobra wiring, main binary. Only orchestrates; defer logic to internal packages.
- `internal/cli`: Command registration, flag parsing, shared Cobra helpers.
- `internal/plan`: Plan-mode orchestration, generation pipelines (code, AI, hybrid), task sizing utilities.
- `internal/build`: Build-mode execution flow, branch/commit policies, reconciliation of Ralphy runs.
- `internal/validation`: YAML schema checks, affirmative constraint enforcement, research-rule validators.
- `internal/ai`: Provider adapters (Anthropic, GPT), self-consistency logic, prompt construction aligned with style anchors.
- `internal/executor`: Interfaces and implementations for Ralphy/OpenCode execution; subpackages like `internal/executor/ralphy` embed scripts and manage process I/O.
- `internal/knowledge`: SQLite access layer, pattern discovery, import/export, confidence scoring.
- `internal/plugins`: Plugin host, discovery, manifests, event bus. Subdirs: `registry`, `loader`, `sandbox` as needed.
- `internal/audit`: Append-only audit log writer, report aggregation, compliance artifacts.
- `internal/telemetry`: Metrics, timing, cost tracking; optional sink adapters.
- `internal/config`: Config loading/merging (defaults, repo `.prompt-stack/config.yaml`, env overrides).
- `internal/shared`: Cross-cutting primitives (errors, logging interfaces, token budgeting helpers). Keep small; prefer domain-specific packages when possible.
- `configs/`: Sample `.prompt-stack/config.yaml`, plugin manifests, default policy bundles.
- `examples/`: Requirements samples, milestone manifests, demo YAML (mirrors `docs/requirements`).
- `scripts/`: Developer automation (`fmt`, lint, release); shell scripts must be invoked via Make targets.

## Package Boundaries & Dependencies
- One-way flow: `internal/cli` → domain packages; domains may depend on `shared`, but never on `cli` or binary code.
- `internal/plan` and `internal/build` depend on cross-cutting services (`config`, `knowledge`, `executor`, `ai`, `validation`, `audit`).
- `internal/executor` exposes interfaces consumed by `internal/build`; concrete implementations may depend on `shared` helpers.
- Disallow cycles by policy: add an automated check (e.g., `go list -deps`) to `Makefile lint`.
- Configuration and secrets flow inward: `internal/config` injects dependencies; packages request explicit structs/interfaces rather than reading env vars directly.

## Testing Conventions
- Co-locate unit tests with packages (`*_test.go`).
- Use `testdata/` under each domain for fixtures (YAML plans, review reports, SQLite seeds).
- Provide integration tests in `internal/build` and `internal/plan` using temporary directories to exercise cmd → internal flows.
- Shared mocks/utilities live in `internal/shared/testsupport` to avoid leaking helpers into production code.

## Generated & Runtime Artifacts
- Runtime state (knowledge DB, audit logs, task traces, vendored Ralphy script) lives under `./.prompt-stack/`; never inside `internal/`.
- Add `internal/executor/ralphy/embed.go` (Go `embed`) to ship the pinned `ralphy.sh`; extraction happens at runtime into `.prompt-stack/vendor/ralphy/`.
- CLI reports write to `reports/` (relative to repo root) consistent with requirements.

## Alignment with Requirements
- Plan vs Build separation maps to `internal/plan` and `internal/build`, honoring dual-mode behavior.
- Commit/branch policies from `docs/requirements/main.md` surface in `internal/build` and associated validator hooks.
- Research-backed checks (style anchors, affirmative constraints, TDD workflow) centralize in `internal/validation` and are consumed by both modes.
- Knowledge compounding (Phase 3) resides in `internal/knowledge`, feeding context minimization utilities shared by plan/build.

## Implementation Guardrails
- Add `Makefile` targets: `make build`, `make test`, `make lint`, `make fmt`, `make docs` (style marker alignment).
- Provide `tools.go` to track codegen/tool dependencies (`go:build tools`).
- Enforce `go test ./...` and `golangci-lint` (or `go vet` + staticcheck) in CI to preserve structure contracts.
- Document directory expectations in `docs/requirements/main.md` (link to this file) so future tasks stay consistent.
