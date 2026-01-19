project_metadata:
  name: "prompt-stack"
  version: "0.1.0"

id: "m0"
title: "Initial setup with requirements gathering phase"
short_description: "Scaffold a Go CLI (Cobra) and interactive requirements prompt; bundle Ralphy shell executor and produce a completed planning input for downstream plan generation."
background: "Project-level requirements and architecture exist under docs/; this milestone creates an initial Go-based CLI scaffold and a prompt-driven requirements capture that outputs a planning input file for Plan Mode."

objectives:
  - "Scaffold Go CLI with plan, validate, review, build, and init commands"
  - "Bundle Ralphy shell executor (vendor/ralphy/ralphy.sh) with dry-run reporting"
  - "Interactive requirements prompt that saves a filled planning input YAML"

success_metrics:
  - metric: "CLI help surface"
    target: "--help lists core commands and exits 0"
  - metric: "Requirements output"
    target: "Produces docs/implementation-plan/m0/requirements.md matching template schema"
  - metric: "Executor dry-run"
    target: "Writes ./.prompt-stack/report.txt and ./.prompt-stack/audit.log"

requirements_file: "docs/implementation-plan/m0/requirements.md"

style_anchors:
  - file: "docs/requirements/architecture.md"
    reason: "Project architecture and problem framing"
  - file: "docs/best-practices.md"
    reason: "Coding and review guidelines to prevent drift"
  - file: "templates/planning-phase.input.yaml"
    reason: "Input template reference for mapping fields"

timeline:
  start_date: "2026-01-19"
  target_completion: "2026-01-21"

scope:
  in_scope:
    - "Scaffold Go CLI (Cobra)"
    - "Implement RalphyShellExecutor via Go embed and materialization"
    - "Interactive one-question-at-a-time requirements prompt and transcript saving"
    - "Minimal README and Makefile with build/test targets"
  out_of_scope:
    - "Full AI provider integration (Anthropic/OpenAI) — adapters may be stubbed"
    - "AST-based context optimization and token budget logic"
    - "Full CI/CD pipelines or cross-platform packaging beyond macOS/Linux"

constraints:
  - "POSIX bash required for bundled ralphy.sh"
  - "Secrets must not be stored in repo or SQLite DB; use OS secret store or env vars"
assumptions:
  - "Author/committer (@kyledavis) is the primary stakeholder"
  - "Examples/ style anchors exist under examples/ and docs/style-markers.md"
  - "Timeline is 1-3 days for M0"

deliverables:
  - "docs/implementation-plan/m0/requirements.md"
  - "docs/implementation-plan/m0/requirements-interview.md"
  - "cmd/prompt-stack/main.go (Cobra scaffold)"
  - ".prompt-stack/vendor/ralphy/ralphy.sh (materialized at runtime)"
  - "Makefile, README.md, and basic unit tests"
  - ".prompt-stack/report.txt and .prompt-stack/audit.log (after dry-run)"

tech_stack:
  languages: ["Go"]
  frameworks: ["Cobra"]
  infra: ["SQLite (optional)"]

integrations:
  - system: "OpenCode (executor target)"
    notes: "Executor will call vendored ralphy.sh which orchestrates OpenCode agents; real OpenCode integration is out of scope for M0 beyond dry-run."

attachments:
  - "docs/requirements/main.md"
  - "docs/requirements/architecture.md"
repo_access:
  repo: "/Users/kyledavis/Sites/prompt-stack"
  read_only: false

execution_resources:
  knowledge_db_path: ".prompt-stack/knowledge.db"
  yaml_validator: "tools/validate_yaml.go"
  ralphy_script: ".prompt-stack/vendor/ralphy/ralphy.sh"

validation_assets:
  - "docs/ralphy-inputs.schema.json"
  - "docs/ralphy-yaml-spec.md"

testing:
  require_unit_tests: true
  require_integration_tests: false
  require_e2e: false
  qa_gates:
    - "manual checklist"

quality_targets:
  quality_score: 0.95

data_classification: "internal"
secrets_included: false

confirmations:
  - "I provided a valid requirements_file path"
  - "I included 2-3 style_anchor files or code anchors with reasons"
  - "I provided knowledge_db_path or confirmed default"
  - "I listed validation_assets (schemas/specs) the plan must reference"

hints:
  prefer_affirmative_constraints: true
  preferred_task_size_minutes: 30-120
  required_top_level_fields:
    - id
    - title
    - short_description
    - requirements_file
    - style_anchors
    - stakeholders.product_owner
    - project_metadata
    - execution_resources
    - validation_assets
    - quality_targets

stakeholders:
  product_owner: "@kyledavis"

usage: |
  # Generate a candidate plan (code path)
  prompt-stack plan docs/implementation-plan/m0/requirements.md --method code --output planning/milestones/m0.ralphy.yaml
