# Project Milestones

This document breaks the requirements in `docs/requirements/main.md` and `docs/requirements/architecture.md` into a concise, test-oriented milestone plan that prioritizes an MVP path and supports manual verification at each step.

Principles
- Ship the smallest end-to-end vertical slice first: Plan Mode (template-based code generation) → basic validation → persist artifacts.
- Keep milestones testable: include short manual test checklist and clear exit criteria for each milestone.
- Optimize for fast iterations and developer-led manual verification before automating tests.
- Prefer non-AI implementations first; add AI flows after the code path is reliable.

Milestones (ordered)

1) M1 — Requirements Gathering (1-3 days)
- Goal: Capture, formalize, and validate project requirements as the first milestone. Produce a concise requirements input file and templates that Plan Mode can consume; this milestone intentionally avoids Ralphy/OpenCode usage and focuses on human-driven requirements collection and static validation.
- Deliverables: `examples/requirements/inputs/requirements.input.md` (or `requirements.md`), `examples/requirements/templates/requirements-prompt.md` (prompt + templates used to gather requirements), an updated `docs/requirements/main.md` entry documenting the gathered requirements, and a short `planning/manifest.yaml` entry referencing the milestone.
- Acceptance criteria:
  - A requirements input file exists at `planning/inputs/requirements.input.md` and follows the project's input template.
  - The example prompt and templates used to gather requirements are committed under `docs/requirements/templates/`.
  - `your-tool plan planning/inputs/requirements.input.md --method code` produces a syntactically valid `tasks.yaml` candidate (code-generation path only; no AI required).
- Manual test checklist:
   1. Run the requirements prompt (copy the template from `templates/requirements-gathering-prompt.md`) interactively and save output to `planning/inputs/requirements.input.md`.
  2. Run `your-tool plan planning/inputs/requirements.input.md --method code` and verify `tasks.yaml` is produced.
  3. Confirm `docs/requirements/main.md` references this milestone and contains the example prompt or link to `docs/requirements/templates/`.

2) M2 — Repo Init & CLI scaffold (1-2 days)
- Goal: Basic repo + CLI surface implemented; commands wired but minimal behavior.
- Deliverables: `your-tool` CLI scaffold, `init`, `plan`, `validate`, `review`, `build` commands (stubs), project layout and `./.your-tool/` default structure.
- Acceptance criteria:
  - `your-tool --help` lists core commands.
  - `your-tool init` creates `./.your-tool/config.yaml` and `./.your-tool/knowledge.db` (empty) and prints instructions.
- Manual test checklist:
  1. Run `your-tool --help` and inspect listed commands.
  2. Run `your-tool init` in a sample repo; verify files created under `./.your-tool/`.
  3. Confirm no secrets are written to the DB.

3) M3 — Requirements Parser + Template-based Plan Mode (2-4 days)
- Goal: Implement requirements parser and template-based YAML generation (code generation fast path).
- Deliverables: `plan` command that takes `requirements.md` or prompt string and writes a `tasks.yaml` using templates and cached defaults.
- Acceptance criteria:
  - `your-tool plan requirements.md` produces syntactically valid `tasks.yaml` matching the Enhanced YAML Structure examples.
  - Basic validation step (schema) runs and succeeds for generated YAML.
- Manual test checklist:
  1. Prepare a simple `requirements.md` and run `your-tool plan requirements.md`.
  2. Open `tasks.yaml`; verify top-level fields: `metadata`, `global_constraints`, `tasks`.
  3. Run `your-tool validate tasks.yaml` — expect pass.
  4. Confirm that `./.your-tool/audit.log` includes an entry for the generation run.

4) M4 — SQLite knowledge DB + caching (2-3 days)
- Goal: Add a small SQLite schema and store/lookup simple patterns (style anchors, coding rules).
- Deliverables: `knowledge` module with `database` and `patterns` APIs; plan uses cached anchors when present.
- Acceptance criteria:
  - `your-tool init` creates/initializes `knowledge.db` with core tables.
  - `your-tool plan` uses cached anchors when available and logs cache hits in `audit.log`.
- Manual test checklist:
  1. Run `your-tool init` and inspect `./.your-tool/knowledge.db` (sqlite3 CLI or `sqlitebrowser`).
  2. Insert a sample style anchor via CLI or small script; run `your-tool plan` and confirm anchor usage.
  3. Verify `task_history` rows are created after plan runs.


9a) M9a — Role-based Model Selection & Repo Config (2-4 days)
- Goal: Add a repo-level `prompt-stack.yaml` that defines model profiles and roles; implement a docs page and basic validator. Enable the agent to map Ralphy tasks to roles and select an appropriate model via capability+cost scoring (advisory rollout).
- Deliverables:
  - `docs/prompt-stack-config.md` describing role-based design and examples.
  - Example `prompt-stack.yaml` in repo root (example only).
  - JSON Schema `docs/prompt-stack.schema.json` and a small `your-tool validate-config` validator CLI (advisory by default).
  - Integration notes: how Ralphy should emit `intent`/`est_tokens`/`role_hint`.
- Acceptance criteria:
  - `docs/prompt-stack-config.md` exists and is linked from `docs/index.md`.
  - `prompt-stack.yaml` example validates against `docs/prompt-stack.schema.json`.
  - `your-tool validate-config prompt-stack.yaml` returns pass/fail and logs the intended role->model mapping for at least one sample `tasks.yaml`.
- Manual test checklist:
  1. Add example `prompt-stack.yaml` to repo root; run `your-tool validate-config prompt-stack.yaml` and confirm pass.
  2. Generate a sample `tasks.yaml` with Ralphy (or hand-edit) that includes `intent` and `est_tokens`; run the selection routine (locally simulated) and confirm it outputs selected model + reason.
  3. Toggle role policy to `strict` in a role and run the validator to confirm it flags disallowed models (validator only; no CI enforcement yet).

9) M9 — Hybrid & AI flow hooks (optional for MVP or staged next sprint) (2-3 weeks)
- Goal: Add meta-PRD templates, Ralphy-based AI generation paths, AI validation and review stages (integrations with OpenCode/Anthropic/OpenAI), and auto-regenerate loop based on score thresholds.
- Deliverables: `ai-generator`, `ai-validator`, templates in `templates/`, `--method ai|hybrid` flag behavior.
- Acceptance criteria:
  - `your-tool plan --method hybrid` runs code path then AI-review; when score < threshold, triggers AI regeneration.
- Manual test checklist:
  1. Run hybrid plan with networked AI providers (or mocked) and confirm the flow and score-driven regeneration.
  2. Confirm `review-report.json` is produced and saved.

10) M10 — Knowledge export/import, team sharing, and simple CI templates (2 days)
- Goal: Provide `knowledge export`/`import`, sample CI workflow for validate/review, and `your-tool init --install-hooks` wiring.
- Deliverables: CLI knowledge commands, `.github/workflows/validate-plan.yml` template, optional Husky hook installer.
- Acceptance criteria:
  - `your-tool knowledge export` writes JSON; `import` restores patterns.
  - Provided CI templates run `your-tool validate` successfully in a sample environment.
- Manual test checklist:
 1. Export knowledge, remove patterns, re-import, and verify patterns return.
 2. Install sample CI config in a test repo and run validate locally.


Testing / QA approach
- Manual-first: every milestone includes a short manual test checklist; require the author (you) to run through these before moving on.
- Create small, focused unit tests after the manual checklist passes for automation in CI.
- Keep each milestone small enough so manual verification takes <30 minutes.

Estimations & priorities
- High-priority MVP path: M2 → M3 → M4 → M5 → M6 → M7. These give an end-to-end Plan Mode (code-gen) and a safe Build Mode with vendored Ralphy dry-run and commit-per-task simulation.
- Mid-term: M8 (context optimization) yields the 80%+ token reduction target.
- Longer: M9 (AI flows) and M10 (team sharing + CI) come afterwards.

Next actions (pick one)
1. Approve this milestone plan and I will create `docs/requirements/milestones.md` in the repo (done).
2. Ask to adjust scope/durations or reorder milestones (tell me which milestone to change).
3. Start implementation: I can scaffold the CLI (M2) next — tell me to proceed.
