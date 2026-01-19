# Architecture

## Overview

prompt-stack is an AI-assisted development workflow tool that generates and validates Ralphy YAML files. The tool provides a CLI interface for planning, building, and validating implementation plans.

## Components

### CLI Layer (`cmd/prompt-stack/`)

The CLI layer provides user-facing commands using Cobra framework:

- **main.go**: Root command setup and entry point
- **plan.go**: Plan command for generating implementation plans
- **validate.go**: Validate command for schema and quality validation
- **build.go**: Build command for executing implementation plan tasks
- **review.go**: Review command for progress and quality metrics
- **init.go**: Init command for interactive requirements gathering

### Packages (`pkg/`)

#### Executor Package (`pkg/executor/`)

Manages execution of the vendored Ralphy shell script:

- **executor.go**: Core executor implementation
  - `Executor`: Main executor struct with working directory and dry-run configuration
  - `ExecutionConfig`: Configuration for task execution
  - `ExecutionResult`: Execution results with success status
  - `Execute()`: Runs ralphy.sh with provided configuration
  - `ValidateInputs()`: Validates execution parameters

- **dry_run.go**: Dry-run functionality
  - `executeDryRun()`: Validates and reports without actual execution
  - `DryRunValidator`: Comprehensive validation for dry-run mode
  - `generateDryRunReport()`: Generates execution reports
  - `writeReport()`: Writes reports to `.prompt-stack/report.txt`
  - `logAudit()`: Logs execution details to `.prompt-stack/audit.log`

#### Prompt Package (`pkg/prompt/`)

Interactive prompt system for requirements gathering:

- **prompt.go**: Core prompt implementation
  - `Prompt`: Manages question flow and response collection
  - `InterviewResult`: Captures responses and transcript
  - `Question`: Question definition with validation
  - `Run()`: Executes interactive interview with context cancellation
  - `NewPrompt()`: Creates new prompt instances

- **questions.go**: Default question set
  - `DefaultQuestions()`: Returns 14 milestone questions
  - Question validators for milestone ID, title, objectives, etc.

### Configuration (`.prompt-stack/`)

- `vendor/ralphy/ralphy.sh`: Vendored Ralphy shell script
- `report.txt`: Execution reports
- `audit.log`: Audit trail of executions

## Data Flow

### Init Workflow

1. User runs `prompt-stack init`
2. Prompt package asks questions sequentially
3. Responses are validated
4. Results saved as YAML and markdown transcript

### Plan Workflow

1. User runs `prompt-stack plan --input requirements.md`
2. Requirements are processed (placeholder)
3. Implementation plan generated (placeholder)

### Build Workflow

1. User runs `prompt-stack build --plan plan.yaml`
2. Executor reads tasks from plan
3. For each task:
   - Validate inputs
   - If dry-run: generate report
   - Otherwise: execute ralphy.sh
4. Log execution details

## Design Principles

- **Single Responsibility**: Each package and function has one clear purpose
- **Separation of Concerns**: CLI layer separate from business logic
- **Extensibility**: Easy to add new commands and validators
- **Testability**: All components are unit testable
- **Context Awareness**: Context.Context passed for cancellation support

## Future Enhancements

- Integrate reusable validators into `internal/validation`
- Add EventBus for before/after events
- Enhanced reporting with JSON schema validation
- Support for parallel task execution
- Integration with OpenCode for AI-assisted planning
