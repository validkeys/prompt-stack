# Milestone m1 Requirements - M2 — Repo Init & CLI scaffold (1-2 days)

## Planning Input YAML

```yaml
id: "m1"
title: "M2 — Repo Init & CLI scaffold (1-2 days)"
short_description: "Basic repo + CLI surface implemented; commands wired but minimal behavior"
background: "This is milestone 2 from the project milestones document, focusing on CLI scaffold implementation"

# Stakeholders
stakeholders:
  product_owner: "@kyledavis"

# Objectives — measurable, max 5
objectives:
  - "Implement help command that lists core commands"
  - "Implement init command that creates .prompt-stack/config.yaml and .prompt-stack/knowledge.db"
  - "Verify project structure matches docs/requirements/project-structure.md spec"

# Success metrics (metric: target)
success_metrics:
  - metric: "prompt-stack --help lists core commands"
    target: "Lists init, plan, validate, review, build commands"
  - metric: "prompt-stack init creates required files"
    target: "Creates .prompt-stack/config.yaml and .prompt-stack/knowledge.db"
  - metric: "No secrets written to database"
    target: "Zero secrets in knowledge.db"

# Primary files/inputs the generator will read (required)
requirements_file: "docs/implementation-plan/m1/requirements.md"
style_anchors:            # 2-3 files that define design/style/architecture
  - "examples/style-anchor/"  # Go CLI project template
  - "docs/style-markers.md"   # Style guidance
  - "docs/requirements/project-structure.md"  # Project layout spec

# Timeline — provide target dates where known
timeline:
  start_date: "2026-01-19"
  target_completion: "2026-01-21"

# Scope (concise lists)
scope:
  in_scope:
    - "CLI scaffold with help command"
    - "init command implementation"
    - "plan, validate, review, build commands (stubs)"
    - "Project layout according to docs/requirements/project-structure.md"
    - ".prompt-stack/ default structure creation"
    - "config.yaml and knowledge.db creation"
  out_of_scope:
    - "Anything not explicitly stated in scope"

# Constraints & assumptions
constraints:
  - "TDD preference with integration tests preferred"
  - "Unit tests only where absolutely necessary"
assumptions:
  - "No runtime resources/tools currently assumed"
  - "No critical constraints or assumptions beyond stated scope"

# Deliverables (artifacts to produce) — structured
deliverables:
  - name: "prompt-stack CLI binary"
    description: "Working CLI with help command"
    owner: "@kyledavis"
    format: "binary"
    due: "2026-01-21"
  - name: "init command implementation"
    description: "Command that creates .prompt-stack/config.yaml and .prompt-stack/knowledge.db"
    owner: "@kyledavis"
    format: "go code"
    due: "2026-01-21"
  - name: ".prompt-stack/config.yaml"
    description: "Default configuration file"
    owner: "@kyledavis"
    format: "yaml"
    due: "2026-01-21"
  - name: ".prompt-stack/knowledge.db"
    description: "Empty SQLite knowledge database"
    owner: "@kyledavis"
    format: "sqlite"
    due: "2026-01-21"

# Acceptance criteria — user-facing verification rules
acceptance_criteria:
  - id: "AC-1"
    title: "Help command lists core commands"
    scenario: "User runs prompt-stack --help to see available commands"
    expected_outcome: "Command lists init, plan, validate, review, build commands"
    validation_method: "manual-review"
    stakeholder_signoff: "required"
    related_deliverables:
      - "prompt-stack CLI binary"
  - id: "AC-2"
    title: "Init command creates required files"
    scenario: "User runs prompt-stack init in a sample repository"
    expected_outcome: "Creates .prompt-stack/config.yaml and .prompt-stack/knowledge.db files"
    validation_method: "manual-review"
    stakeholder_signoff: "required"
    related_deliverables:
      - "init command implementation"
      - ".prompt-stack/config.yaml"
      - ".prompt-stack/knowledge.db"
  - id: "AC-3"
    title: "No secrets in database"
    scenario: "After init command runs, inspect knowledge.db contents"
    expected_outcome: "No API keys, tokens, or secrets present in database"
    validation_method: "manual-review"
    stakeholder_signoff: "required"
    related_deliverables:
      - ".prompt-stack/knowledge.db"

# Tech and integrations (short)
tech_stack:
  languages: ["Go"]
  frameworks: ["Cobra"]
  infra: ["SQLite"]

integrations:
  - system: "Git"
    notes: "Repository initialization"

# Access & attachments
attachments:
  - "docs/implementation-plan/m1/requirements-interview.md"
repo_access:
  repo: "/Users/kyledavis/Sites/prompt-stack"
  read_only: true

# Testing & QA expectations
testing:
  require_unit_tests: false
  require_integration_tests: true
  require_e2e: false
  qa_gates:
    - "TDD workflow"
    - "Manual test checklist completion"

# Privacy & secrets handling
data_classification: "internal"
secrets_included: false

# Project metadata
project_metadata:
  name: "prompt-stack"
  version: "v1"

# Execution resources
execution_resources:
  knowledge_db_path: ".prompt-stack/knowledge.db"
  validator_scripts: []

# Validation assets
validation_assets:
  - "docs/requirements/project-structure.md"
  - "docs/ralphy-inputs.schema.json"
  - "docs/requirements/main.md"
  - "docs/requirements/milestones.md"

# Quality targets
quality_targets:
  - metric: "quality_score"
    threshold: ">=0.95"
    validation_method: "automated-test"
    description: "YAML syntax validation, JSON Schema validation, secrets scan, style anchors compliance, task sizing compliance, affirmative constraints validation, multi-layer enforcement validation"

# Minimal checklist the user must confirm
confirmations:
  - "I provided a valid requirements_file path"
  - "I included 1-3 style_anchor files or code anchors"
  - "I granted read access to the repo if code discovery is required"

# Optional generator hints (keeps this short)
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
```

## Usage

```bash
# Generate a candidate plan (code path)
prompt-stack plan docs/implementation-plan/m1/requirements.md --method code --output planning/milestones/m1.ralphy.yaml
```

## Manual Test Checklist
1. Run `prompt-stack --help` and inspect listed commands.
2. Run `prompt-stack init` in a sample repo; verify files created under `./.prompt-stack/`.
3. Confirm no secrets are written to the DB.

## References
- Interview transcript: `docs/implementation-plan/m1/requirements-interview.md`
- Project structure: `docs/requirements/project-structure.md`
- Style guidance: `docs/style-markers.md`
- Example template: `examples/style-anchor/`