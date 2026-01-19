# Implementation Plan Generation Prompt (template)

Purpose: Generate a validated Ralphy YAML (planning-phase) inputs file from a project requirements document. This prompt uses `templates/planning-phase.prd-template.yaml` as the structural template and references canonical docs for style, schema, and verification.

Use case: run this prompt (interactive or programmatically) with a filled requirements file (e.g. `docs/implementation-plan/m0/requirements.md`) to produce a single `make-implementation-plan.prd.yaml` artifact (validation-ready) suitable for downstream planning and Ralphy execution.

---

Initial context (required inputs)
- `requirements_file` (path): path to the saved milestone requirements file (YAML or markdown) you should read and consume.
- `milestone_id` (slug): short identifier (example: `m0`, `m1`). Used to name outputs and paths.
- `output_dir` (path): where to save generated artifacts. Default: `docs/implementation-plan/<milestone>/`.
- `knowledge_db_path` (path): optional path to knowledge DB. Default: `./.prompt-stack/knowledge.db`.
- `style_anchors` (list): 2–3 canonical file paths or URLs used as style anchors (each with a one-line reason). If omitted, search repo for candidate anchors but record them as `assumption` in the output.
- `reference_docs` (list, optional): prioritized documents (paths or URLs) that codify project rules, quality gates, or research. If not provided, inspect the repository index (for example `docs/index.md`) or ask for recommended references before generating outputs.
- `temp_artifacts_dir` (path): where to write temporary/sidecar artifacts created during generation (validation JSON, schema reports, generated tests, helper code, secrets-scan outputs). Default: `./.prompt-stack/reports/{{milestone_id}}/` (create subfolders per milestone).
- `validator_command` (string, optional): path or invocation to repository validator. Default: `prompt-stack validate`. If provided (or if `prompt-stack` is present in PATH), the generator must prefer and call the validator instead of producing new validator code.
- (No `prd_output_path` input): the agent will place the final PRD document alongside the provided `requirements_file` as `make-implementation-plan.prd.yaml`.

References (consult these files in the repo)
- `reference_docs` list (if provided) — read these first; they codify project-specific rules and research.
- `templates/planning-phase.prd-template.yaml` — primary template to follow for structure and phases (replace with the project-specific template path if different).
- `docs/ralphy-inputs.md` and `docs/ralphy-yaml-spec.md` — canonical Ralphy input guidance (or equivalent references supplied by the stakeholder).
- `docs/ralphy-inputs.schema.json` — JSON Schema to validate final YAML (or the schema path provided via inputs).
- `docs/best-practices.md`, `docs/task-sizing.md`, or their equivalents — research-backed guidance for anchors, affirmative constraints, and task sizing.
- Additional policy/governance docs (for example repository rules, drift policies, TDD checklists) discovered via `reference_docs` or the repository index (`docs/index.md` if available).
- `tools/validate_yaml.go` (or the project-provided validator) — preferred YAML syntax validator.

Prompt rules and behavior (agent must follow exactly)

- Reuse shared helpers: Do not regenerate repository-level validation helpers. Prefer the general-purpose helpers in `internal/validation/helpers.go` for common tasks (file existence, JSON/YAML reading, string utilities, task-section extraction, repo root detection). If you need additional glue code, write minimal adapters into `temp_artifacts_dir` and mark them with `assumption: true` so future milestones reuse the canonical helpers.
- If the validator CLI (`prompt-stack validate`) is available prefer invoking it rather than emitting new validator code. If the CLI is not present, produce a validator spec that references the helpers in `internal/validation`.

Prompt rules and behavior (agent must follow exactly)
1. Ask the user: "Where is the requirements document (path) that I should read?" and record the answer. The agent will create the final PRD file named `make-implementation-plan.prd.yaml` in the same folder as the provided `requirements_file`.
2. Read `requirements_file` fully and produce a concise 1-paragraph summary (one sentence summary + 2–3 bullet highlights: objectives, success metrics, constraints).
3. Load `templates/planning-phase.prd-template.yaml` and use it as the authoritative structure. Replace placeholders (`{{...}}`) with concrete values derived from `requirements_file` and provided inputs.
4. Generate a single `make-implementation-plan.prd.yaml` artifact (top-level sections: `metadata`, `global_constraints`, `tasks`, `instructions`, `validation`) that:
   - Includes project `name`, `version`, `description`, and milestone data from the requirements file.
   - Embeds `requirements_file` under `files_in_scope` for all relevant phases (at minimum planning-001 and planning-004).
   - Populates `style_anchors` with the provided anchors (2–3 per task) and records `file` plus `reason` for each entry.
   - Ensures each task has: `id`, `title`, `description`, `outputs`, `files_in_scope`, `style_anchors`, `estimated_duration_minutes`, `verification`, and any dependency references.
   - Applies `global_constraints` consistent with `templates/planning-phase.prd-template.yaml` (affirmative constraints, task sizing ranges, forbidden/required patterns).
   - Records inline validation and quality findings inside the same document (e.g. under `validation.reports` and `metadata.final_quality_report`).
5. For any missing but required inputs (style anchors, knowledge_db_path, schema references, or `reference_docs`), add the field with an `assumption:` note and set a sensible default. Mark those placeholders in the YAML with `assumption: true` and include a short justification inside the same document.
6. Inline validation steps (performed by the agent) must summarize their findings within the `validation` section of the YAML rather than emitting external files.
6.a. All temporary or sidecar artifacts produced by the generator (intermediate JSON reports, schema validation artifacts, generated test files or helper code, secrets-scan outputs, etc.) MUST be written to `temp_artifacts_dir` (organize by milestone) and MUST NOT be written into `output_dir` or `docs/implementation-plan/<milestone>/`.
6.b. The only files created in `output_dir` are:
   - `make-implementation-plan.prd.yaml` (the PRD)
   - `implementation-plan.yaml` (the final execution plan for downstream execution)
6.c. Summaries of validation and scans must be embedded inside the PRD under `validation.reports` and `metadata.assumptions`; full sidecar artifacts may be stored in `temp_artifacts_dir` for debugging or CI inspection.
7. Quality gates (the generated YAML must meet these or the agent should flag warnings and produce a remedial plan) should be documented inline under `validation.gates`.
8. Produce a `final_quality_report` summary block within the YAML metadata including: `quality_score` (0.0–1.0), `issues` (list), and `approval` flag (`APPROVED` if >= 0.95, else `NEEDS_REVISION`). Provide a short explanation of how the score was computed (weights: anchors 30%, sizing 25%, schema 20%, secrets 15%, affirmative constraints 10%).
9. Output artifact to write to disk: a single `make-implementation-plan.prd.yaml` file saved to `output_dir` (no secondary copies or sidecar reports in that folder). Additionally, produce `implementation-plan.yaml` (the final execution plan) saved to `output_dir`. All validation, schema, and secrets-scan notes must be embedded inside the YAML under `validation.reports` or `metadata.assumptions`.

Validator integration (new)
- If a repository validator exists (recommended) prefer it over generating validator code in the PRD flow. Default command: `prompt-stack validate`.
- Behavior when `validator_command` or `prompt-stack` is available:
  - Run the validator after generating `implementation-plan.yaml`:
    `{{validator_command}} --input {{output_dir}}/implementation-plan.yaml --schema docs/ralphy-inputs.schema.json --out {{temp_artifacts_dir}}/validation.json`
  - Read `validation.json` and embed a concise summary into `validation.reports` and `metadata.final_quality_report` in the PRD.
  - Do NOT invent or output full validator source code when the validator CLI is available.
- Behavior when validator is not present:
  - Perform inline syntactic/schema checks and produce a validator *spec* (interface, expected report shape, and test fixtures) and write the spec to `temp_artifacts_dir` marked with `assumption: true` and `todo: implement validator as internal/validation`.

Suggested output path (default)
- `docs/implementation-plan/<milestone>/make-implementation-plan.prd.yaml`
- `docs/implementation-plan/<milestone>/implementation-plan.yaml` (final execution plan)

Usage snippet (also include in generated `make-implementation-plan.prd.yaml` comments or metadata):

```
# Generate the plan (code path)
prompt-stack plan docs/implementation-plan/<milestone>/requirements.md --method code --output docs/implementation-plan/<milestone>/implementation-plan.yaml

# Preferred: run built-in validator to produce structured JSON summary and embed it back into the PRD
prompt-stack validate --input docs/implementation-plan/<milestone>/implementation-plan.yaml --schema docs/ralphy-inputs.schema.json --out ./.prompt-stack/reports/{{milestone_id}}/validation.json
```

Deliverable format & placeholders
- The prompt consumer should replace or accept runtime mappings for these placeholders:
  - `{{requirements_file}}`, `{{milestone_id}}`, `{{output_dir}}`, `{{knowledge_db_path}}`, `{{style_anchor_*}}`, `{{temp_artifacts_dir}}`, `{{validator_command}}`.
- The produced YAML must be valid UTF-8, adhere to the schema at `docs/ralphy-inputs.schema.json`, and store validation notes, assumptions, and quality scoring inline (no external sidecar files required for summary).


Edge cases & remediation guidance (agent responsibilities)
- If the requirements document is ambiguous or missing an objective/success metric: ask one clarifying question and record the Q/A in `generation_log.md`. If the user is unreachable, set a conservative assumption and mark it with `assumption: true` and provide the rationale.
- If style anchors are not found in the repo, recommend 2–3 best-effort anchors from `docs/best-practices.md` sections and mark them as `assumption` in the YAML.
- If large tasks (>150 minutes) are detected, automatically propose a split into smaller tasks and include both the proposed new tasks and a `split_reason` in `generation_log.md`.

Notes for reviewers (you)
- This is a reusable template prompt. Review that the referenced files exist in the repo (`docs/ralphy-inputs.md`, `docs/ralphy-inputs.schema.json`, `docs/best-practices.md`, `tools/validate_yaml.go`) and adjust paths if your repo layout differs.
- The prompt leans on `templates/planning-phase.prd-template.yaml` for task/phases structure; update the template reference if that file is moved or renamed.

---

End of prompt template. Use this as `docs/requirements/templates/implementation-plan-generation-prompt.md` for generating Ralphy planning inputs from milestone requirements.