Opencode Repo Rules

Purpose
- Repository-level rules and hard constraints for the Opencode agent. These are loaded at session start and should be authoritative.

General
- Agent: `opencode`
- Enforce strict TypeScript rules: NEVER use `any`. Prefer `unknown` if required.
- Runtime validation: Use Zod for all external inputs.
- Allowed stack: `node`, `typescript`, `zod`, `eslint`, `prettier`, `vitest`, `pnpm`.
- Disallowed: adding new dependencies without human approval.

Prompting and edits
- Always include the provided `style_anchors` from `docs/ralphy-inputs.md` when generating edits.
- Only edit files matching `outputs.allowed_file_edits` in Ralphy inputs unless explicitly permitted.
- Provide a brief justification for any non-trivial refactor in the commit body.

Type rules
- NEVER use `any` anywhere in code produced by the agent.
- NEVER inline-disable ESLint rules; if a rule must be disabled, add a comment explaining why and file a follow-up issue.
- ALWAYS type function parameters and return types explicitly.

Testing & TDD
- Tests required for feature changes when `tdd.required` is true in Ralphy inputs.
- Use test commands from `tdd.test_command` and include exact run instructions in the commit message.

Task sizing
- Honor `task_sizing.min_minutes` and `task_sizing.max_minutes`.
- Split tasks that exceed `task_sizing.max_minutes` or touch > `task_sizing.max_files`.

CI & linting
- Precommit hooks and CI must enforce `--max-warnings 0` for ESLint.
- Any `eslint-disable` occurrences should be counted and reported in CI.

Drift
- Follow `docs/drift-policy.md` for stop/revert actions.

Model behavior
- Prefer `minimal-diff` edits; avoid renames unless explicitly requested.
- If tests fail, return failing output and follow TDD instructions in `docs/tdd-checklist.md`.

Updating these rules
- Changes to this file require human review and must be recorded in the repository changelog.
