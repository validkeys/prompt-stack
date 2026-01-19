TDD Checklist — repo-level

Purpose
- Anchor tests-first development for AI-assisted edits. Ensure the model always receives a clear test plan and run/fix loop.

Before implementation
1. Describe the behavior to be tested in 1–3 sentences.
2. List test types required (unit, integration, e2e, property) and minimal scope for now.
3. Provide failing tests first (code + steps to run). Include exact commands (e.g., `pnpm test --filter <pattern>`).
4. Note any required test doubles/mocks and where they live.

Implementation loop (TDD)
1. Write failing test(s).
2. Implement the minimal code to make the test pass.
3. Run test suite locally; fix until passing.
4. Add more tests to cover edge cases and refactor.
5. Run full suite and linters, then commit.

Commit policy
- One logical change per commit; prefer small commits that keep tests green.
- Include `Test:` prefix in commit messages when the change is primarily test-related.

Model instructions for failing tests
- When tests fail during model-driven work, return the failing test output to the model and instruct:
  "Revise implementation to pass this test while keeping all previously passing tests. Do not modify the test. Do not add dependencies."

Local verification
- Provide exact commands to run relevant tests and linters (examples):
  - `pnpm test` (run all)
  - `pnpm test -- -t "UserService"` (run specific tests)
  - `pnpm lint` (ESLint)

CI expectations
- CI must run the same commands; failing CI blocks merges.
- Coverage targets optional but document if enforced.

Templates
- Minimal test template (JS/TS):

```ts
// tests/<feature>.test.ts
import { describe, it, expect } from 'vitest'
import { myFn } from '../src/myFn'

describe('myFn', () => {
  it('does the expected thing', () => {
    expect(myFn(input)).toEqual(expected)
  })
})
```

Notes
- Don't change tests to make model output look good; change implementation only.
- If a failing test requires external infra (DB, API), provide a reproducible local mock or docker-compose snippet.
