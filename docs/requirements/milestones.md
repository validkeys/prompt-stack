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
  1. Run the requirements prompt (copy the template from `docs/requirements/templates/requirements-prompt.md`) interactively and save output to `planning/inputs/requirements.input.md`.
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

5) M5 — Basic Validator and Review (non-AI) (2 days)
- Goal: Implement YAML schema validation, simple research-practice checks (style anchors present, task sizing ranges, affirmative constraints) and a `review` command that returns a score and warnings.
- Deliverables: `validator` + `reviewer` modules, `your-tool validate` and `your-tool review` working offline (no AI), JSON/text report output.
- Acceptance criteria:
  - `your-tool validate tasks.yaml` fails on broken schema; passes on valid YAML.
  - `your-tool review tasks.yaml` outputs a quality score and lists any missing TDD/workflow warnings.
- Manual test checklist:
  1. Mutate `tasks.yaml` to violate a rule (e.g., no style anchors) and run `your-tool review` to see warning.
  2. Add a TDD workflow to a task and confirm warning resolves.

6) M6 — Ralphy Executor (MVP shell) + vendor bundling (3 days)
- Goal: Implement `Executor` interface and `RalphyShellExecutor` that bundles `vendor/ralphy/ralphy.sh` using embed-like behavior (materialize script, make executable, run via `os/exec`), capture stdout/stderr to `report.txt` and `audit.log`.
- Deliverables: executor module, `your-tool build` invokes Ralphy script in a dry-run mode and captures outputs.
- Acceptance criteria:
  - `your-tool build tasks.yaml --dry-run` runs the vendored ralphy.sh and returns structured output (or stubbed output) captured to `./.your-tool/report.txt` and `./.your-tool/audit.log`.
  - Executor exposes `Capabilities()` including `supportsParallel` etc.
- Manual test checklist:
  1. Run `your-tool build tasks.yaml --dry-run` and inspect `report.txt` and `audit.log` for expected content.
  2. Verify the vendored `ralphy.sh` is executable at `./.your-tool/vendor/ralphy/ralphy.sh`.

7) M7 — Build Mode: Commit-per-task + pre-flight checks (2-3 days)
- Goal: Implement basic build orchestration: pre-flight checks (clean git tree), commit-per-task behavior (when enabled), verification hooks integration (run ESLint/tsc commands as configured), and enforcement of scope.
- Deliverables: build orchestration wiring commits and verification steps; drift detection (files modified outside scope cause abort).
- Acceptance criteria:
  - `your-tool build tasks.yaml` on a sample repo with a clean tree executes tasks (simulated) and if `--commit-per-task` is set, creates atomic commits limited to `files_in_scope`.
  - Build aborts if working tree is dirty.
- Manual test checklist:
  1. Create a sample git repo with initial commit; run `your-tool build tasks.yaml --dry-run` and inspect simulated commits.
  2. Make an unrelated change outside any `files_in_scope` and verify build flags drift and aborts.
  3. Toggle `--commit-per-task` and confirm commit behavior.

8) M8 — Context Optimization (minimal viable features) (4-6 days)
- Goal: Implement non-AI context optimizations that give high ROI: smart file filtering (ripgrep), line-range extraction, simple token estimation, and a context budget check that suggests splitting large tasks.
- Deliverables: `context` module providing `file-filter`, `budget-calc`, and integration into `plan` and `review` flows.
- Acceptance criteria:
  - `your-tool plan` reports context reduction percentage and warns when budget > limit.
  - `your-tool review` indicates budget usage per task.
- Manual test checklist:
  1. Run `your-tool plan` on a repo with many files; confirm top-N filtered files are used and reduction % is plausible.
  2. Create a very large task and verify recommendation to split.

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
