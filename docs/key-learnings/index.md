Key Learnings — PRD → Ralphy Inputs

- Make templates explicit vs concrete: when producing a concrete Ralphy inputs file, include all schema-required top-level fields (`name`, `version`, `rules_file`, `task_sizing`, `tdd`, `model_preferences`, `outputs`) to avoid schema validation failures.

- Enforce anchors early: include 2–3 concrete style anchors (code + tests + README) on every planning task to prevent architectural drift and give Ralphy exact patterns to follow.

- Mark inferred edits: when filling templates automatically, mark inferred additions with `assumption: true` and a short rationale so reviewers can quickly find and accept/reject them.

- Respect task-sizing constraints: enforce 30–150m task estimates; if a task is shorter, either increase estimate or split it and record rationale.

- Keep validation inline: add a `validation` summary block in PRDs with `quality_score`, `issues`, and `approval` so automation and humans read a single canonical result.

- Prefer concrete execution snippets: add explicit commands to run YAML/schema/secrets validators in `instructions` to improve reproducibility.

- Scope implementation rules: keep implementation-only pattern checks (e.g., Zod imports) scoped to implementation-phase PRDs or mark them `when: implementation_phase` to avoid blocking planning artifacts.

- Use repository examples as anchors: small, well-scoped examples (like `examples/style-anchor/pkg/greeter/greeter.go`) are high-leverage anchors; include tests and README for TDD/context.

(End)