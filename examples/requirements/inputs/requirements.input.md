# Project name: AI-Assisted Development Workflow Tool (planning example)

Short description (one sentence):
Create a CLI tool that generates validated Ralphy YAML plans from human-written requirements, with a code-generation fast path and optional AI-assisted flows.

Primary user persona(s):
- Senior developer (primary)

Primary goals (3–5):
- Produce quality-guaranteed Ralphy YAML plans from concise requirements.
- Provide a safe Build Mode that executes plans with pre-flight checks and commit-per-task safety.
- Make Plan Mode fast for routine tasks (code generation) and high-quality for complex tasks (hybrid/AI).

Success metrics (3):
- Plan quality score > 0.9 for 95% of generated plans.
- Second-run plan generation < 10s using cached knowledge.
- First-pass build success > 90% in MVP-simulated runs.

Non-goals (what we won't do):
- No automatic secret storage in the knowledge DB.
- No remote execution by default (local-only unless explicitly enabled).

Critical constraints (security, latency, cost, infra):
- Store API keys in OS secret stores or env vars; do not persist keys in DB.
- Default to local execution; require explicit opt-in for remote agent execution.

Must-have integrations (e.g., Git, CI, DB, third-party services):
- Git (local + remote awareness), SQLite for local knowledge cache, Ralphy (execution) integration as vendored shell for MVP.

Acceptable tech stack (optional):
- Node.js or Bun, TypeScript, SQLite (better-sqlite3), Commander.js

Privacy / secrets handling preferences:
- No secrets stored in knowledge DB; require runtime injection or OS secret store.

Initial scope (what will be delivered in milestone 1):
- A filled `examples/requirements/inputs/requirements.input.md` and template under `examples/requirements/templates/`.
- Plan Mode code-generation path should accept this input and produce a syntactically valid `tasks.yaml`.
- Documentation updated to reference these example files.

Out-of-scope (for milestone 1):
- AI generation, Ralphy runtime execution, and knowledge export/import.

Initial non-functional targets (latency, size, performance):
- Code generation <5s for small inputs; validation <1s.

Regulatory/compliance concerns (if any):
- None for the example project; follow project-level GDPR/security norms if handling user data.

Testing approach (unit, integration, TDD?):
- Manual-first for milestone 1; add unit tests after manual verification.

Definition of done for milestone 1 (acceptance criteria):
- `examples/requirements/inputs/requirements.input.md` exists and follows the template.
- `your-tool plan examples/requirements/inputs/requirements.input.md --method code` produces a `tasks.yaml` candidate.
- Documentation references updated and point to `examples/requirements`.

Notes / examples / reference links:
- Templates live in `examples/requirements/templates/` to keep docs self-contained but examples separate from code templates.
