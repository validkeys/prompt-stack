assumption: true
reason: "`prompt-stack` CLI not found in local PATH; creating validator spec as a temporary artifact for CI or later implementation."

Validator CLI Spec (expected interface)

- Command: `prompt-stack validate --input <implementation-plan.yaml> --schema docs/ralphy-inputs.schema.json --out <output.json>`
- Exit codes: 0 = success (no schema errors), non-zero = validation issues found
- Output JSON shape (validation.json):
  {
    "valid": boolean,
    "errors": [
      {"path": "JSON-Pointer", "message": "human-readable message", "schemaKeyword": "required|type|pattern|..."}
    ],
    "summary": {
      "totalErrors": integer,
      "severityCounts": {"error": int, "warning": int}
    }
  }

Suggested adapter (todo: implement)
- Implement a small Go adapter under internal/validation/cli_adapter.go that wraps the repository's validator helpers and emits the JSON shape above.
- Place adapter or test fixtures under `./.prompt-stack/reports/m1/` and mark with `assumption: true` until internal/validation is available.

Usage example (CI):

prompt-stack validate --input docs/implementation-plan/m1/implementation-plan.yaml --schema docs/ralphy-inputs.schema.json --out ./.prompt-stack/reports/m1/validation.json

Notes:
- If `prompt-stack` is present in PATH in CI, prefer running the CLI instead of the adapter.
- This spec is informational and intended to be consumed by the planning pipeline; it does not implement validation logic.
