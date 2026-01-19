# Requirements Gathering — Interview Workflow Prompt (reusable canonical)

Purpose: run a short, one-question-at-a-time interview with a stakeholder, record each Q&A, then produce a filled planning input (based on `templates/planning-phase.input.yaml`) and save it as the milestone requirements document.

Initial step (project context)
- Ask the stakeholder for brief context and links to any relevant files, docs, or repos that will help you understand what they're trying to build. Example prompts you can use:
  - "Give me a short project summary and links to any existing requirements, design docs, or repos."
  - "List 2–3 files or URLs that best show the code/design style we should follow."
- After the stakeholder replies, repeat back a concise summary of what you understood from their entry (one or two sentences). Confirm you have the correct context before proceeding.

Before you begin (confirm references)
- Ask which reference documents you should consult for this interview (e.g. project requirements, milestone list, or a planning template). If the stakeholder has local paths or repo URLs, record them.
- Confirm the milestone identifier/slug to use (e.g. `m1`). All interview transcripts and the final requirements file will be placed under `docs/implementation-plan/<milestone>/` unless the stakeholder specifies a different location.

Interview rules
1. Ask one question at a time. Wait for the stakeholder's answer before asking the next question.
2. Record every question and answer verbatim in `docs/implementation-plan/<milestone>/requirements-interview.md` using this Q/A markdown format:

   Q: <question text>

   A: <answer text>

3. After each answer, confirm the captured answer back in one sentence ("Captured: ...") and ask the next logical question.
4. Continue until you have at least the required fields from the planning input template (see list below) or the stakeholder says they are done.
5. When finished, generate a filled planning input document (YAML preferred) using the collected answers and save it to `docs/implementation-plan/<milestone>/requirements.md`. Include a short `usage` snippet that demonstrates how to run the plan generator against the saved file.

Required fields to gather (map these to `templates/planning-phase.input.yaml`)
- id (short slug)
- title (one-line)
- short_description
- requirements_file (path — can point to this saved requirements.md)
- style_anchors (1-3 files or references)
- stakeholders.product_owner (name/email or handle)
- objectives (1-5)
- success_metrics (metric + target)

Helpful additional fields (ask when available)
- background
- constraints and assumptions
- scope: in_scope / out_of_scope
- deliverables
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
7. "Which files or code areas should be used as style anchors? (2-3 recommended)"
8. "What should be included in scope for this milestone?"
9. "What is explicitly out of scope?"
10. "Any critical constraints or assumptions? (security, infra, timelines)"
11. "What deliverables do you expect?" (e.g., task_breakdown.yaml, quality_report.json)
12. "Do you have desired timelines or dates? (start / target completion)"
13. "Testing requirements: unit tests, integration, TDD preference?"
14. "Any privacy, compliance, or secrets handling notes?"

Finalization
- Produce a filled planning input YAML using the collected answers and save it to `docs/implementation-plan/<milestone>/requirements.md`.
- Save the interview transcript to `docs/implementation-plan/<milestone>/requirements-interview.md` (Q/A format).
- Append a short `usage` section in the saved requirements.md showing the `plan` command that consumes it; for example:

```bash
# Generate a candidate plan (code path)
your-tool plan docs/implementation-plan/<milestone>/requirements.md --method code --output planning/milestones/<id>.ralphy.yaml
```

Notes and best practices
- Keep questions short and actionable; prefer affirmative phrasing.
- If the stakeholder is unsure about a field, capture the assumption as a note in the interview transcript and mark the field as `assumption: <text>` in the final YAML.
- When possible, prefer file paths or URLs for style anchors so the plan generator can inspect concrete examples.
- All final artifacts should be human-readable and committed under `docs/implementation-plan/` for traceability.

---

This is the canonical, project-agnostic interview prompt. Use it from `templates/requirements-prompt.md` in code, CI, or manual workflows.