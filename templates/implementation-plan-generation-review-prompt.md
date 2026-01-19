# Implementation Plan Generation Review Prompt (template)

Purpose: Evaluate a generated `make-implementation-plan.prd.yaml` for completeness, schema alignment, and adherence to repository rules and research-backed best practices. This prompt is run after the planning template produces a PRD to confirm the artifact is ready for downstream execution.

---

Initial context (required inputs)
- `prd_path` (path): absolute or repo-relative path to the `make-implementation-plan.prd.yaml` under review.
- `schema_path` (path): JSON Schema used to validate Ralphy inputs (default: `docs/ralphy-inputs.schema.json`).
- `rules_file` (path): canonical repo rules (default: `docs/opencode-rules.md`).
- `additional_refs` (list, optional): extra documents that should influence review (for example drift policy, TDD checklist).
- `commit_policy_doc` (path, optional): doc describing commit, branch, or drift rules (default: `docs/drift-policy.md`).
- `task_sizing_doc` (path, optional): task-sizing research (default: `docs/task-sizing.md`).
- `temp_artifacts_dir` (path): where temporary/sidecar artifacts produced during generation are stored (validation JSON, schema reports, generated tests, helper code, secrets-scan outputs). Default: `./.prompt-stack/reports/{{milestone_id}}/`.
- `validator_command` (string, optional): path or invocation to repository validator. If provided (or if `your-tool` is present in PATH), prefer running the validator and consuming its structured JSON output instead of regenerating validation code.

References (load before review)
- `prd_path`
- `docs/best-practices.md`
- `docs/ralphy-inputs.md`
- `schema_path`
- `rules_file`
- `commit_policy_doc`
- `task_sizing_doc`
- Any `additional_refs`
- `temp_artifacts_dir` (if present) — reviewer should consult full sidecar artifacts here when deeper inspection is required (schema_validation_report.json, secrets_scan_report.json, generated helper code).

Agent directives
1. Confirm the PRD is the planning-phase artifact for the specified milestone and contains the expected top-level sections: `metadata`, `global_constraints`, `tasks`, `instructions`, `research_compliance`, `performance_targets`, and inline `validation` summaries. Flag missing sections.
2. Validate metadata:
   - Ensure `name`, `description`, `version`, `rules_file`, `style_anchors`, `allowed_dependencies`, `prompt_template`, `model_preferences`, and `outputs` fields exist where required by `docs/ralphy-inputs.md`.
   - Check timestamps, generator identifiers, and `quality_target` fields for presence and plausibility.
   - Record any assumptions or defaults the template made; confirm they are flagged with `assumption: true` and include rationale.
3. Inspect `global_constraints` for alignment with `docs/best-practices.md`:
   - Style anchors required: 2–3 per task; confirm the rule exists and is restated at the end of the document.
   - Task sizing window: 30–150 minutes; verify affirmative framing for all constraints.
   - Forbidden/required pattern lists should apply to planning outputs; flag any constraint that targets implementation-only code (for example unconditional Zod import requirements).
4. Walk each task entry:
   - Confirm `id`, `title`, `description`, `outputs`, `files_in_scope`, `style_anchors` (2–3 unique anchors with reasons), `estimated_duration_minutes`, `verification`, and dependency data are present.
   - Check `files_in_scope` include `requirements_file` for relevant phases.
   - Ensure durations fall within the task sizing window; flag outliers and recommend splits/merges.
   - Confirm verification statements are affirmative ("Resolve ambiguity before completing" instead of "No ambiguity remaining").
5. Analyze `style_anchors` across the document:
   - Prefer concrete code or policy files over repeating research docs.
   - Ensure anchors exist on disk or are marked with `assumption: true` plus justification.
   - Highlight any task with fewer than two anchors or with redundant references.
6. Review inline `validation` blocks:
   - Confirm validation summaries, schema checks, YAML syntax confirmations, secrets scan results, and quality gating are embedded inside the PRD; sidecar file references should only exist if the file also stores the summarized results inline.
   - If `validator_command` or `your-tool` is available, run the validator and compare its `final_quality_report.json`/`validation.json` against the PRD's inline summary; if discrepancies exist, prefer the validator's structured output and update the PRD accordingly (apply minimal edits and record rationale in the review output).
   - Cross-check `quality_score`, `issues`, and approval flag; recompute if weights (anchors 30%, sizing 25%, schema 20%, secrets 15%, affirmative constraints 10%) were misapplied.
7. Verify multi-layer enforcement requirements:
   - Prompt-level rules, IDE/LSP checks, pre-commit hooks, CI checks, runtime validation, and commit policy references should all appear.
   - Validate presence of `files_in_scope` limits and commit-per-task expectations.
8. Confirm the PRD references execution guidance:
   - Usage snippet for regenerating the plan and running validators should be included and accurate.
   - Implementation guidance should avoid generating tests but should mention when future tasks require TDD or validation steps.
9. If any required input is missing (anchors, knowledge DB path, schema references), ensure the PRD marks the gap with `assumption: true` and a remediation note.
10. Collect open questions for the stakeholder when information is insufficient; record them in the review output.
11. Make all obvious refinements directly to the PRD document: apply minimal, precise edits to fix formatting, missing fields that can be safely inferred, typos, and alignment with repository rules. For each applied edit, include a one-line rationale in the review output and preserve original content with `assumption: true` where applicable.
12. If any information is uncertain or ambiguous, interview the stakeholder one question at a time; pause after each question and proceed only after receiving the user's response. Do not batch multiple clarification questions into a single message.

Output format (return to user)
- `summary`: 2–3 sentence overview of PRD health.
- `compliance_table`: Markdown table listing review categories (`Metadata`, `Global Constraints`, `Task Structure`, `Style Anchors`, `Validation`, `Enforcement`, `Affirmative Tone`, `Schema Alignment`, `Secrets`, `Execution Guidance`) with `status` (`OK`, `WARN`, `FAIL`) and one-line notes.
- `issues`: numbered list of findings. For each, include `severity` (`HIGH`, `MEDIUM`, `LOW`), `section/task`, description, and explicit fix recommendation.
- `missing_data`: bullet list of required inputs or assumptions that must be resolved.
- `follow_up_questions`: bullet list of stakeholder questions (if any).
- `quality_score`: recomputed numeric score and approval verdict (`APPROVED` or `NEEDS_REVISION`).

Reviewer behavior
- Keep language affirmative and actionable.
- When recommending fixes, cite exact paths or YAML keys (for example `tasks.planning-006.style_anchors`).
- Prefer precise, minimal-diff corrections over broad rewrites.
- Encourage aligning the generation template with this review template when mismatches are detected.
