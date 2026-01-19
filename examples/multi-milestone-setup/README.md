# Multi-Milestone Setup (Concise Example)

This folder shows a minimal recommended layout and sample files for iterative milestone-based planning.

Structure
- planning/
  - inputs/
    - auth-v1.input.yaml        # per-milestone input (short)
  - milestones/
    - auth-v1.ralphy.yaml      # generated Ralphy YAML for milestone
  - reports/
    - auth-v1/final_quality_report.json
  - manifest.yaml              # lightweight project manifest referencing inputs

Example manifest (planning/manifest.yaml)
```
project_id: "example-project"
title: "Example Multi-Milestone Roadmap"
milestones:
  - id: "auth-v1"
    input: "planning/inputs/auth-v1.input.yaml"
    sequence: 1
    depends_on: []
```

Example input (planning/inputs/auth-v1.input.yaml)
```
id: "auth-v1"
title: "Authentication: Login & JWT"
short_description: "Add user login, sessions, and JWT issuance."
requirements_file: "docs/requirements/auth.md"
style_anchors:
  - "src/schemas/user.schema.ts"
  - "src/services/auth.service.ts"
stakeholders:
  product_owner:
    name: "Product Owner"
    email: "po@example.com"
timeline:
  start_date: "2026-02-01"
  target_completion: "2026-02-14"

# Keep inputs minimal — the planning template will expand these into a full Ralphy YAML.
```

Example commands

- Generate single milestone (recommended):
  `your-tool plan planning/inputs/auth-v1.input.yaml --method hybrid --review --output planning/milestones/auth-v1.ralphy.yaml`

- Manifest-driven integration check:
  `your-tool plan --manifest planning/manifest.yaml --integration-check`

Copy these files into your repo and adapt paths as needed.
