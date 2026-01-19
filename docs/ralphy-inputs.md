Ralphy YAML: Required Inputs

Purpose
- Single, machine- and human-readable checklist of the repository-level inputs the Ralphy YAML generator must consume. Keep this in sync with `docs/ralphy-yaml-spec.md`.

Top-level fields (recommended)
- `name`: short project name
- `description`: 1–2 sentence description of scope
- `version`: semantic version for the YAML inputs
- `rules_file`: path to repo rule file (e.g., `docs/opencode-rules.md` or `.opencode/rules.yml`)
- `style_anchors`: array of example file paths to include in prompts (e.g., `examples/auth/*.ts`)
- `allowed_dependencies`: array of package names allowed for model edits; new dependencies must be human-approved
- `task_sizing`:
  - `min_minutes`: 30
  - `max_minutes`: 150
  - `max_files`: recommended max touched files for a single task
- `tdd`:
  - `required`: true|false
  - `test_command`: command to run tests (string or array)
  - `failure_instruction`: text the model must receive when tests fail
- `ci`:
  - `precommit`: list of commands run locally (e.g., `husky`, `lint-staged` entries)
  - `ci_checks`: list of commands run in CI (e.g., `pnpm test`, `pnpm lint`)
  - `count_eslint_disable`: boolean/threshold flag
- `drift_policy_ref`: relative path to drift policy doc (e.g., `docs/drift-policy.md`)
- `validation_schemas`: paths where Zod/validation schemas live
- `prompt_template`: object with `prefix`, `suffix`, and `placeholders` (used to assemble model prompts)
- `model_preferences`:
  - `primary`: e.g., `opencode` (or a model name)
  - `review_model`: e.g., `gpt-4` or `gemini`
  - `strategies`: e.g., `minimal-diff`, `generate-3-variants`
- `outputs`:
  - `allowed_file_edits`: list of glob patterns the generator may edit
  - `disallowed_file_edits`: list of glob patterns
  - `commit_policy`: commit message prefixes and rules

Metadata and validation
- Provide JSON Schema or simple key/value validation rules for Ralphy to validate inputs before generating YAML.
- Required fields should be marked; provide default fallbacks where safe.

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

Maintenance
- Keep `style_anchors` and `allowed_dependencies` up to date; add a changelog entry to YAML or use `version`.
