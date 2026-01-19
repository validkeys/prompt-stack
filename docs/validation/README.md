Validation report & CLI integration

This folder defines the standardized validation report schema used by `internal/validation` and consumed by Plan/Build flows, AI review prompts, and CI.

Running the validator

- Preferred: use the built-in validator via the main CLI:

  ```bash
  your-tool validate --input implementation-plan.yaml --schema docs/ralphy-inputs.schema.json --out ./.your-tool/reports/final_quality_report.json
  ```

- The validator produces a JSON report that conforms to `docs/validation/validation-report.schema.json`.
- Exit codes:
  - `0` = pass (overall_score >= configured threshold)
  - `2` = warnings (score < threshold but non-critical issues)
  - `3` = fail (critical issues present)

Integration notes

- Plan generation and review prompts should prefer running the built-in validator rather than generating validator source code.
- If the validator is not available in the repo, generation prompts should produce a validator *spec* (interface + test vectors) and mark implementation as `assumption: true` with remediation guidance.

Report schema

- `overall_score`: 0.0–1.0
- `overall_result`: `PASS`/`WARN`/`FAIL`/`APPROVED`/`NEEDS_REVISION`
- `component_scores`: map of per-check scores (anchors, sizing, schema, secrets, constraints)
- `issues[]`: list of issues with `severity`, `path`, `message`, and optional `fix_suggestion`

