Drift Policy for AI-Assisted Edits

Purpose
- Define when to stop, revert, and document when the model's edits deviate from project rules or introduce unexpected changes.

Stop & Revert Criteria
- The model introduces new dependencies not listed in allowed stack.
- Edits touch files outside the specified targets (>3 unexpected files).
- Linting or type errors are introduced and cannot be resolved within the task's scope.
- Tests fail and the model proposes changing tests to pass.

Immediate Actions
1. Stop the model session and return a clear error message summarizing the deviation.
2. Revert the workspace to the state before the task started (use `git` where appropriate) — only if workspace changes were produced by the agent in this session. If unsure, consult the repo owner.
3. Create a short incident note in `docs/drift-incidents/` with:
   - What happened
   - Files changed unexpectedly
   - New dependencies proposed
   - Suggested remediation steps

Recording learnings
- After resolving, update `.cursor/rules/` or `CLAUDE.md` with new rules to prevent recurrence.

Allowed Deviations
- Minor formatting changes (editor config) or whitespace-only edits.
- Single-line refactors within targeted files if they are within scope and type-checked.

Review & Approval
- Any revert or incident note must be reviewed by a human maintainer before re-running the model.
