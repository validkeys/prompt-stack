# AI Codegen Best Practices (concise)

## Core principle
- Models optimize locally; enforce global constraints with layered verification (prompt → IDE → commit → CI → runtime).

## Style anchors
- Always include 2–3 exemplary files as templates in prompts.
- Reference exact paths and `touch ONLY` targets.
- Prefer concrete repository examples with code + tests + README (e.g., `examples/style-anchor/pkg/greeter/greeter.go`).
- Enforce anchors early to prevent architectural drift.

## Task sizing
- Split work into 30m–2.5h atomic tasks (30–150 minutes optimal).
- Limit scope to specific files; commit after each small task; revert immediately on drift.
- Respect task-sizing constraints; if a task is shorter than 30m, either increase estimate or split it and record rationale.

## Affirmative instructions
- State permitted actions explicitly (e.g., `ONLY use: express, zod, prisma`).
- Avoid negative framing.

## Tiered rules
- Global: user prefs (format, language, length).
- Project: persistent rules in `CLAUDE.md` or `.cursor/rules/` (loaded every session).
- Context-aware: auto-attached rules per directory or file pattern.

## Prompt-level Go rules (copy into CLAUDE.md)
- NEVER use `interface{}` (empty interface) unless absolutely necessary.
- NEVER ignore errors; always handle or explicitly return errors with context.
- NEVER use type assertions for external data — use proper validation.
- ALWAYS include doc comments for exported functions, types, and packages.

## IDE / linting
- Enforce strict Go linting rules:
  - Use `go vet` to catch common mistakes
  - Use `golangci-lint` with strict configuration
  - Ensure zero `gofmt` violations
  - Run `go test -race` for race condition detection
- Configure editor to surface these as errors.

## Commit & CI
- Pre-commit: Use `make lint` and `make test` with zero warnings.
- CI: quality gates that count `gofmt` violations and fail if thresholds exceeded.

## Runtime validation
- Validate all external inputs with proper error handling and validation instead of type assertions.

## TDD as anchor
- Require a TDD checklist before implementation (tests → minimal code → more tests → refactor).
- When tests fail, return failing output to the model and instruct: "Revise implementation to pass this test while keeping all previously passing tests. Do not modify the test. Do not add dependencies."

## Prompt positioning
- Put critical specs, style anchors, and hard rules at the beginning and reiterate them at the end of prompts (avoid burying requirements in the middle).

## Model strategies
- Claude: use for surgical, minimal-diff edits; request `research → plan → implement`, `minimal diff, no renames, explain each edit`; use thinking triggers (`think`, `think hard`, `ultrathink`).
- GPT: use for exploratory/greenfield work and code review; ask for tactical plans and side‑effect checks in modules X/Y.

## Self-consistency & AI-on-AI review
- Generate 3+ implementations (higher temperature), then ask the model to pick the most consistent and explain why.
- Use multi-model review (e.g., Claude writes, GPT/Gemini reviews) to catch subtle issues.

## Drift handling
- Stop and revert immediately on unexpected dependencies or unfamiliar patterns; do not fix mid-stream.
- Document learnings and update persistent rules after each session.

## Key learnings from milestone 0
- Make templates explicit vs concrete: include all schema-required top-level fields to avoid validation failures
- Enforce concrete style anchors early: include 2–3 concrete anchors (code + tests + README) on every planning task
- Mark inferred edits: use `assumption: true` with rationale for any inferred additions
- Respect task-sizing constraints: enforce 30–150m task estimates; split shorter tasks with rationale
- Keep validation inline: add `validation` summary with `quality_score`, `issues`, and `approval`
- Prefer concrete execution snippets: add explicit validator commands in `instructions`
- Scope implementation rules: keep implementation-only pattern checks scoped with `when: implementation_phase`
- Use repository examples as anchors: small, well-scoped examples are high-leverage anchors

## Quick practical checklist
1. Create `CLAUDE.md` or `.cursor/rules/` with prompt-level rules above
2. Add 2–3 concrete style anchors to prompts (prefer repository examples)
3. Rescope tasks to 30m–2.5h and commit per task
4. Convert negative constraints to affirmative instructions
5. Enforce Go linting zero-warnings, pre-commit hooks, and CI gates
6. Require TDD plans and tests before making changes
7. Use proper error handling for runtime validation and count `gofmt` violations in CI
8. Include all schema-required fields in generated YAML
9. Mark inferred additions with `assumption: true` and rationale

