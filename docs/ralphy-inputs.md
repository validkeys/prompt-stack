Ralphy YAML: Required Inputs

Purpose
- Single, machine- and human-readable checklist of the repository-level inputs the Ralphy YAML generator must consume. Keep this in sync with `docs/ralphy-yaml-spec.md`.

Top-level fields (required by schema)
- `name`: short project name
- `version`: semantic version for the YAML inputs (pattern: `^\d+\.\d+\.\d+$`)
- `rules_file`: path to repo rule file (e.g., `docs/opencode-rules.md` or `.opencode/rules.yml`)
- `task_sizing`: object with `min_minutes`, `max_minutes`, `max_files` (optional)
- `tdd`: object with `required`, `test_command`, `failure_instruction` (optional)
- `model_preferences`: object with `primary`, `review_model`, `strategies` (optional)
- `outputs`: object with `allowed_file_edits`, `disallowed_file_edits`, `commit_policy` (optional)

Top-level fields (recommended)
- `description`: 1–2 sentence description of scope
- `style_anchors`: array of concrete example file paths to include in prompts (prefer repository examples with code + tests + README)
- `allowed_dependencies`: array of package names allowed for model edits; new dependencies must be human-approved
- `ci`:
  - `precommit`: list of commands run locally (e.g., `husky`, `lint-staged` entries)
  - `ci_checks`: list of commands run in CI (e.g., `pnpm test`, `pnpm lint`)
  - `count_eslint_disable`: boolean/threshold flag
- `drift_policy_ref`: relative path to drift policy doc (e.g., `docs/drift-policy.md`)
- `validation_schemas`: paths where Zod/validation schemas live
- `prompt_template`: object with `prefix`, `suffix`, and `placeholders` (used to assemble model prompts)

Metadata and validation
- Provide JSON Schema or simple key/value validation rules for Ralphy to validate inputs before generating YAML.
- Required fields should be marked; provide default fallbacks where safe.
- All schema-required top-level fields (`name`, `version`, `rules_file`, `task_sizing`, `tdd`, `model_preferences`, `outputs`) must be present to avoid schema validation failures.

Example minimal Ralphy inputs (YAML)

```yaml
name: prompt-stack
description: AI-assisted development policy and inputs
version: 0.1.0
rules_file: docs/opencode-rules.md
style_anchors:
  - examples/style/anchor-1.ts
allowed_dependencies:
  - zod
  - eslint
  - vitest
task_sizing:
  min_minutes: 30
  max_minutes: 150
  max_files: 5
tdd:
  required: true
  test_command: pnpm test
ci:
  precommit:
    - pnpm lint
  ci_checks:
    - pnpm test
    - pnpm lint
drift_policy_ref: docs/drift-policy.md
validation_schemas:
  - src/schemas/*.ts
prompt_template:
  prefix: "Project rules:\n{{rules_file}}\n\nTask:\n"
  suffix: "\nMake a minimal diff. Commit after tests pass."
model_preferences:
  primary: opencode
  review_model: gpt-4
outputs:
  allowed_file_edits:
    - src/**
    - tests/**
  disallowed_file_edits:
    - scripts/**
    - .github/**
commit_policy:
  prefix_rules:
    - feat:
    - fix:
    - test:
```

How to use
- Ralphy generator should load this file and validate against the schema in `docs/ralphy-yaml-spec.md`.
- When fields conflict with `docs/opencode-rules.md`, the repo-level rules file wins.

Key learnings from milestone 0
- Include all schema-required top-level fields to avoid validation failures
- Enforce 2–3 concrete style anchors per task (prefer repository examples with code + tests + README)
- Mark inferred additions with `assumption: true` and rationale
- Respect task-sizing constraints (30–150 minutes)
- Keep validation inline with `quality_score`, `issues`, and `approval` flags
- Use repository examples as high-leverage anchors

Maintenance
- Keep `style_anchors` and `allowed_dependencies` up to date; add a changelog entry to YAML or use `version`.
