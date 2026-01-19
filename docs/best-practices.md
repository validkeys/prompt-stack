# AI Codegen Best Practices (concise)

## Core principle
- Models optimize locally; enforce global constraints with layered verification (prompt → IDE → commit → CI → runtime).

## Style anchors
- Always include 2–3 exemplary files as templates in prompts.
- Reference exact paths and `touch ONLY` targets.

## Task sizing
- Split work into 30m–2.5h atomic tasks.
- Limit scope to specific files; commit after each small task; revert immediately on drift.

## Affirmative instructions
- State permitted actions explicitly (e.g., `ONLY use: express, zod, prisma`).
- Avoid negative framing.

## Tiered rules
- Global: user prefs (format, language, length).
- Project: persistent rules in `CLAUDE.md` or `.cursor/rules/` (loaded every session).
- Context-aware: auto-attached rules per directory or file pattern.

## Prompt-level TypeScript rules (copy into CLAUDE.md)
- NEVER use `any`. Use `unknown` if absolutely necessary.
- NEVER disable ESLint rules inline without explanation.
- NEVER use type assertions (`as`) for external data — use Zod schemas.
- ALWAYS explicitly type function parameters and return types.

## IDE / linting
- Enforce strict ESLint/@typescript-eslint rules such as:
  - `@typescript-eslint/no-explicit-any: 'error'`
  - `@typescript-eslint/no-unsafe-argument: 'error'`
  - `@typescript-eslint/no-unsafe-assignment: 'error'`
  - `@typescript-eslint/no-unsafe-return: 'error'`
  - `@typescript-eslint/explicit-function-return-type: 'error'`
- Configure editor to surface these as errors.

## Commit & CI
- Pre-commit: Husky + lint-staged with `--max-warnings 0`.
- CI: quality gates that count `eslint-disable`/`eslint-disable-next-line` and fail if thresholds exceeded.

## Runtime validation
- Validate all external inputs with Zod schemas instead of `as` assertions.

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

## Quick practical checklist
1. Create `CLAUDE.md` or `.cursor/rules/` with prompt-level rules above
2. Add 2–3 style anchors to prompts
3. Rescope tasks to 30m–2.5h and commit per task
4. Convert negative constraints to affirmative instructions
5. Enforce ESLint zero-warnings, pre-commit hooks, and CI gates
6. Require TDD plans and tests before making changes
7. Use Zod for runtime validation and count `eslint-disable` in CI

