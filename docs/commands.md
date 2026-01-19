# Commands Reference

## Overview

prompt-stack provides 5 core commands for AI-assisted development workflow management.

## Commands

### init

Initialize a new milestone with interactive requirements gathering.

```sh
prompt-stack init [flags]
```

#### Flags

- `--output-dir, -o`: Directory to save output files (default: `docs/implementation-plan/m0`)

#### Description

Runs an interactive interview to gather milestone requirements. Asks 14 questions about:

- Milestone ID, title, description
- Stakeholder information
- Objectives and success metrics
- Style anchors to reference
- Scope and out-of-scope items
- Constraints and assumptions
- Deliverables and timeline
- Testing requirements
- Privacy and security considerations

#### Output

- `requirements-interview.md`: Complete transcript of the interview
- `requirements.md`: Formatted requirements in YAML structure

#### Example

```sh
./dist/prompt-stack init
./dist/prompt-stack init --output-dir docs/implementation-plan/m1
```

#### Features

- Sequential question flow with validation
- Required and optional questions
- Helpful error messages for invalid inputs
- Graceful cancellation with Ctrl+C

---

### plan

Generate implementation plans from requirements.

```sh
prompt-stack plan [flags]
```

#### Flags

- `--input`: Input requirements file (required)
- `--output`: Output directory for plan (optional)

#### Description

Generate implementation plans from requirements or templates. (Currently placeholder - full implementation in progress)

#### Example

```sh
./dist/prompt-stack plan --input docs/requirements.md --output docs/implementation-plan/m0/
```

---

### validate

Validate implementation plans.

```sh
prompt-stack validate [flags]
```

#### Flags

- `--input`: Input implementation plan file (required)
- `--output`: Output file for validation report (optional)

#### Description

Validate implementation plans against schema and quality standards. (Currently placeholder - full implementation in progress)

#### Example

```sh
./dist/prompt-stack validate --input docs/implementation-plan/m0/final_implementation-plan.yaml
```

---

### build

Build project from implementation plan.

```sh
prompt-stack build [flags]
```

#### Flags

- `--plan`: Implementation plan file (required)
- `--dry-run`: Validate without executing (optional)

#### Description

Build project components based on implementation plan tasks. Uses the executor package to orchestrate Ralphy script execution.

#### Example

```sh
./dist/prompt-stack build --plan docs/implementation-plan/m0/final_implementation-plan.yaml
./dist/prompt-stack build --plan docs/implementation-plan/m0/final_implementation-plan.yaml --dry-run
```

#### Features

- Reads tasks from implementation plan
- Validates inputs before execution
- Dry-run mode generates reports without executing
- Execution results logged to `.prompt-stack/audit.log`
- Reports written to `.prompt-stack/report.txt`

---

### review

Review implementation progress.

```sh
prompt-stack review [flags]
``

#### Description

Review implementation progress and quality metrics. (Currently placeholder - full implementation in progress)

#### Example

```sh
./dist/prompt-stack review
```

---

## Global Flags

All commands support Cobra's global flags:

- `--help, -h`: Show help for the command
- `--version, -v`: Show version information

## Exit Codes

- `0`: Success
- `1`: Error occurred
- `2`: Invalid usage or arguments

## Error Handling

All commands wrap errors with context for debugging:

```
Error: interview failed: failed to create output directory: permission denied
```

Use the error message to identify the issue and consult the documentation for resolution steps.
