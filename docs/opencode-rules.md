Opencode Repo Rules

Purpose
- Repository-level rules and hard constraints for the Opencode agent. These are loaded at session start and should be authoritative.

General
- Agent: `opencode`
- Enforce strict Go conventions: Use idiomatic Go patterns and standard library.
- Runtime validation: Use proper error handling and validation for all external inputs.
- Allowed stack: `go`, `cobra`, `testing`, `go vet`, `gofmt`, `golangci-lint`.
- Disallowed: adding new dependencies without human approval.

Prompting and edits
- Always include the provided `style_anchors` from `docs/ralphy-inputs.md` when generating edits.
- Only edit files matching `outputs.allowed_file_edits` in Ralphy inputs unless explicitly permitted.
- Provide a brief justification for any non-trivial refactor in the commit body.

Go conventions
- NEVER use `interface{}` (empty interface) unless absolutely necessary; prefer specific interfaces or generics.
- NEVER ignore errors; always handle or explicitly return errors with context.
- ALWAYS use proper Go naming conventions: PascalCase for exported identifiers, camelCase for unexported.
- ALWAYS include doc comments for exported functions, types, and packages.
- Use `go vet` and `golangci-lint` to catch common issues; fix all warnings.

Testing & TDD
- Tests required for feature changes when `tdd.required` is true in Ralphy inputs.
- Use test commands from `tdd.test_command` (typically `go test ./...`) and include exact run instructions in the commit message.
- Follow Go testing conventions: test files end with `_test.go`, use table-driven tests for multiple cases.
- Test coverage should be maintained; run `go test -cover ./...` to check coverage.

Task sizing
- Honor `task_sizing.min_minutes` and `task_sizing.max_minutes`.
- Split tasks that exceed `task_sizing.max_minutes` or touch > `task_sizing.max_files`.

CI & linting
- Precommit hooks and CI must enforce `gofmt -d` and `golangci-lint run` with zero warnings.
- Code must pass `go vet ./...` without issues.
- Use `make fmt`, `make lint`, and `make test` commands as defined in the project Makefile.

Drift
- Follow `docs/drift-policy.md` for stop/revert actions.

Model behavior
- Prefer `minimal-diff` edits; avoid renames unless explicitly requested.
- If tests fail, return failing output and follow TDD instructions in `docs/tdd-checklist.md`.

Updating these rules
- Changes to this file require human review and must be recorded in the repository changelog.
