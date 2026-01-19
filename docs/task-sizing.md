Task Sizing Guidance

Purpose
- Keep model-directed work in predictable, reviewable chunks to reduce drift and unexpected scope growth.

Guidelines
- Target task duration: 30 minutes to 2.5 hours.
- Prefer making changes limited to 1–3 files where possible.
- If a requested change touches >5 files or will take >2.5 hours, split it into smaller tasks.

Process
1. Define the smallest useful deliverable and a clear acceptance test.
2. Provide example files and exact paths to edit.
3. After each completed task, run tests and commit with a clear message.

Breaking down large tasks
- Decompose into: (a) tests + scaffolding, (b) minimal implementation, (c) refactor & polish.
- Create a plan note in the issue or prompt listing the sub-tasks and acceptance criteria.

Examples
- Small: fix a bug in `src/utils/parse.ts` — 30–60 mins
- Medium: add new API endpoint `src/server/user.ts` with tests — 1–2.5 hours
- Large: migrate auth system — split into design + incremental PRs
