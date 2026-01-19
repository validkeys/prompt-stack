# Requirements Gathering — Interview Workflow Prompt (reusable canonical)

Purpose: run a short, one-question-at-a-time interview with a stakeholder, record each Q&A, then produce a filled planning input (based on `templates/planning-phase.input.yaml`) and save it as the milestone requirements document.

Initial step (project context)
- Ask the stakeholder for brief context and links to any relevant files, docs, or repos that will help you understand what they're trying to build. Example prompts you can use:
  - "Give me a short project summary and links to any existing requirements, design docs, or repos."
  - "List 2–3 files or URLs that best show the code/design style we should follow."
- After the stakeholder replies, repeat back a concise summary of what you understood (one sentence). Confirm before proceeding.

Before you begin (confirm references)
- Ask which reference documents you should consult for this interview and record any local paths or repo URLs.
- Confirm the milestone identifier/slug to use (e.g. `m1`).
- Ask where outputs should be stored. Default: `docs/implementation-plan/<milestone>/`.

Interview rules
1. Ask one question at a time. Wait for the stakeholder's answer before asking the next question.
2. Record every question and answer verbatim in `docs/implementation-plan/<milestone>/requirements-interview.md` using this Q/A markdown format:

   Q: <question text>

   A: <answer text>

3. After each answer, confirm the captured answer back in one sentence ("Captured: ...") and ask the next logical question.
4. Continue until you have all required fields or the stakeholder says they are done.
5. When finished, generate a filled planning input document (YAML preferred) using the collected answers and save it to `docs/implementation-plan/<milestone>/requirements.md`. Include a short `usage` snippet showing the `plan` command that consumes it.

Required fields to gather (map these to `templates/planning-phase.input.yaml`)
- id (short slug)
- title (one-line)
- short_description
- requirements_file (path — can point to this saved requirements.md)
- style_anchors (2–3 files or references; include a short reason per anchor)
- stakeholders.product_owner (name/email or handle)
- objectives (1-5)
- success_metrics (metric + target)
- project_metadata: `name` and `version` (semantic version)
- execution_resources: paths for knowledge DB, validators, or key scripts (e.g. `.prompt-stack/knowledge.db`, `tools/validate_yaml.go`)
- validation_assets: schema/spec files the plan should reference (e.g. `docs/ralphy-inputs.schema.json`)
- quality_targets: quality/acceptance thresholds (e.g. quality_score >= 0.95)

Helpful additional fields (ask when available)
- background
- constraints and assumptions
- scope: in_scope / out_of_scope
- deliverables (capture structured fields: name, description, owner, format, due)

- timeline.start_date and timeline.target_completion
- testing expectations
- data_classification / secrets_included

Suggested question sequence (ask in order, but follow up when answers need clarification)
1. "What is the milestone id/slug you want to use for this work?"
2. "Give me a one-line title for the milestone."
3. "Provide a short (1-2 sentence) description of the milestone goal."
4. "Who is the primary stakeholder or product owner for this milestone?" (capture contact)
5. "What are the top objectives for this milestone? (3 max)"
6. "How will we measure success? Give metrics and targets."
7. "What is the canonical project name and current version we should record?"
8. "Which documents or schema files must planners reference (paths/URLs)?" (validation assets)
9. "Which runtime resources/tools does the plan assume (knowledge DB path, validator scripts)?" (execution resources)
10. "Which files or code areas should be used as style anchors? Provide 2–3 file paths and a short reason for each."
11. "What should be included in scope for this milestone?"
12. "What is explicitly out of scope?"
13. "Any critical constraints or assumptions? (security, infra, timelines)"
14. "What deliverables do you expect? For each deliverable capture: name, one-line description, owner (name/email), format (yaml/json/docx), and desired due date. Example: {name: task_breakdown.yaml, description: task list, owner: alice@example.com, format: yaml, due: 2026-02-01}"
15. "Do you have desired timelines or dates? (start / target completion)"
16. "Testing requirements: unit tests, integration, TDD preference?"
17. "Any privacy, compliance, or secrets handling notes?"
18. "Are there quality/acceptance thresholds we must meet (e.g., quality score >= 0.95)? For each threshold, also describe how it will be validated (automated test, manual review, smoke test)."

Style anchor guidance
- Request concrete file paths or URLs for anchors and a short reason for relevance.
- Prefer 2–3 anchors per task; record them verbatim so templates can reference them directly.

Finalization
- Produce a filled planning input YAML using the collected answers and save it to `docs/implementation-plan/<milestone>/requirements.md`.
- Save the interview transcript to `docs/implementation-plan/<milestone>/requirements-interview.md` (Q/A format).
- Map captured placeholders to concrete values before generating the planning input (e.g. `{{requirements_file}}`, `{{knowledge_db_path}}`, `{{style_anchor_*}}`). If the stakeholder defers, record a sensible default in the YAML and mark the field as an `assumption`.
- Append a short `usage` section in the saved requirements.md showing the `plan` command that consumes it; for example:

```bash
# Generate a candidate plan (code path)
prompt-stack plan docs/implementation-plan/<milestone>/requirements.md --method code --output planning/milestones/<id>.ralphy.yaml
```

Notes and best practices
- Keep questions short and actionable; prefer affirmative phrasing.
- If the stakeholder is unsure about a field, capture the assumption as a note and mark it in the final YAML (`assumption: <text>`).
- When possible, prefer file paths or URLs for style anchors so the plan generator can inspect concrete examples.
- Ensure validation assets and execution resources are recorded so downstream tasks can run schema/DB checks.
- All final artifacts should be human-readable and committed under `docs/implementation-plan/` for traceability.

---

This is the canonical, project-agnostic interview prompt. Use it from `templates/requirements-gathering-prompt.md` in code, CI, or manual workflows.
