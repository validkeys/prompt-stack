# Ralphy YAML Specification

This document describes the YAML file structures used by Ralphy (https://github.com/michaelshimeles/ralphy).
Place the project config at `.ralphy/config.yaml` and task lists (PRD) in a YAML file (default `tasks.yaml`, configurable with `--yaml`).

Notes
- The primary script is `ralphy.sh` and it uses `yq` to read/write YAML.
- This spec documents the fields Ralphy reads/writes and how the script expects them to be shaped.

1) .ralphy/config.yaml (project configuration)

Purpose: store project metadata, commonly-run commands, AI rules and file boundaries that the agent must respect.

Top-level keys
- project (map) — project metadata
  - name: string — project name (auto-detected by `ralphy --init`)
  - language: string — e.g. "JavaScript", "TypeScript", "Python", "Go", "Rust"; defaults to "Unknown" when not detected
  - framework: string — optional, comma-separated frameworks detected (e.g. "Next.js, React")
  - description: string — free text description (optional)

- commands (map) — recommended commands for CI / local developer workflow
  - test: string — command to run tests (e.g. "npm test" or "pytest")
  - lint: string — command to run linter (e.g. "npm run lint" or "ruff check .")
  - build: string — build command (e.g. "npm run build" or "cargo build")

- rules (sequence of strings) — instructions the AI MUST follow. Injected into prompts.
  - type: sequence
  - default: []
  - examples:
    - "Always use TypeScript strict mode"
    - "All API endpoints must validate input with Zod"

- boundaries (map) — paths the agent should not modify or should always test
  - never_touch: sequence of strings — globs or paths which must not be modified by the agent (default: [])
    - examples: "src/legacy/**", "migrations/**", "*.lock"
  - always_test: sequence of strings — (supported by the script’s boundary loader) files/globs that should always be covered by tests (optional)
    - NOTE: `always_test` is referenced by the script as an accepted boundary type; the initializer only creates `never_touch` by default but `always_test` may be used by users.

Example `.ralphy/config.yaml`

```yaml
# Ralphy Configuration
project:
  name: "my-app"
  language: "TypeScript"
  framework: "Next.js"
  description: "Backend for foo service"

commands:
  test: "npm test"
  lint: "npm run lint"
  build: "npm run build"

rules:
  - "Always follow repository coding conventions"
  - "Write unit tests for new behavior"

boundaries:
  never_touch:
    - "src/legacy/**"
    - "migrations/**"
  always_test:
    - "src/core/**"
```

Implementation notes (how the script uses these fields)
- `ralphy.sh` reads project fields with `yq -r '.project.name'` etc.
- `rules` are loaded with `yq -r '.rules // [] | .[]' $CONFIG_FILE` and injected into AI prompts.
- `boundaries` values are loaded with `yq -r ".boundaries.$boundary_type // [] | .[]" $CONFIG_FILE`.
- To add a rule programmatically the script runs `yq -i '.rules += [env(RULE)]' $CONFIG_FILE`.

2) tasks YAML (PRD format — default file `tasks.yaml` or other file passed via `--yaml`)

Purpose: define the list of tasks (PRD) Ralphy will iterate through in PRD mode.

Top-level structure
- tasks: sequence of task objects

Task object fields
- title: string — required; short task title used to pick and display the task
- completed: boolean — optional; when `true` the task is considered done; default behavior treats missing or `false` as incomplete
- parallel_group: integer — optional; tasks sharing the same positive integer are considered part of the same parallel group and can run in parallel; default is 0 (no group)
- (open) additional fields — any other fields are ignored by Ralphy core but may be useful for users (e.g. `description`, `estimate`, `owner`).

Behavioral notes
- The script selects incomplete tasks using `yq -r '.tasks[] | select(.completed != true) | .title'`.
- Marking a task complete is done by `yq -i "(.tasks[] | select(.title == \"$task\")).completed = true" $PRD_FILE`.
- Parallel groups are read with `.parallel_group // 0` and tasks in a group are selected with equality checks.

Example `tasks.yaml`

```yaml
tasks:
  - title: "Add user registration endpoint"
    completed: false
    parallel_group: 0
    description: "Implement POST /users, input validation and unit tests"

  - title: "Add login rate limiting"
    completed: false
    parallel_group: 1

  - title: "Add forgot-password email flow"
    completed: false
    parallel_group: 1
```

Minimal example (only titles)

```yaml
tasks:
  - title: "Fix login bug"
  - title: "Add dark mode toggle"
```

3) Markdown PRD (alternate PRD source)

- The script also supports a Markdown PRD (`PRD.md` by default). Tasks are lines that start with `- [ ] ` (unchecked) and completed tasks use `- [x] `.
- The script marks tasks complete by replacing `- [ ] <task>` with `- [x] <task>` using `sed`.

4) GitHub issues PRD (alternate PRD source)

- Ralphy can fetch tasks from open GitHub issues in a repo supplied via `--github owner/repo` and optional `--github-label` filter. It uses `gh issue list --repo <repo> --state open --json number,title` and treats each issue as a task in the form `number:title`.

5) YQ compatibility notes

- Ralphy expects `yq` (Mike Farah's yq) to be installed and available on PATH. Many lookups use `yq -r` to read scalar values and array elements.
- When editing `config.yaml` or `tasks.yaml` manually, ensure proper quoting of titles that contain double quotes because the script uses `yq -i` and shell interpolation when marking tasks complete.

6) Recommended practices

1. Keep `title` values unique enough to be safely matched by `yq` when marking complete.
2. Prefer explicit `completed: false` or `true` for clarity, although missing `completed` is treated as incomplete.
3. Use `parallel_group` only when tasks can truly be executed independently; group 0 (or omitted) means single-task execution.
4. Use `boundaries.never_touch` to protect generated, vendor, or sensitive files; use glob-style paths.
5. Add repository-specific `rules` to guide AI behavior and reduce regressions.

7) Quick reference — YAML path snippets (yq)

- Ralphy inputs file: `docs/ralphy-inputs.md` (human checklist) and `docs/ralphy-inputs.schema.json` (JSON Schema for validation). The generator should prefer `.ralphy/config.yaml` when present but fall back to `docs/ralphy-inputs.md`/example YAML when initializing.


- Project name: `.project.name`
- Test command: `.commands.test`
- Rules (iterate): `.rules[]`
- Boundaries (never_touch): `.boundaries.never_touch[]`
- All incomplete task titles: `.tasks[] | select(.completed != true) | .title`
- Mark a task complete (yq in-place): `yq -i "(.tasks[] | select(.title == \"<TITLE>\")).completed = true" tasks.yaml`

If you want, I can:
1) open a PR adding `docs/ralphy-yaml-spec.md` to this repository (requires push access), or
2) run a quick validation script that checks a repository's `tasks.yaml` and `.ralphy/config.yaml` for common problems (duplicates, missing fields).

