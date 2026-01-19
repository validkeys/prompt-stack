# Requirements Gathering — Interview Workflow Prompt (canonical)

Follow the interview workflow in `docs/.prompts/1-make-milestone-requirements.md`.

Purpose: run a short, one-question-at-a-time interview with a stakeholder, record each Q&A, then produce a filled planning input (based on `templates/planning-phase.input.yaml`) and save it as the milestone requirements document.

Before you begin
- Repeat back the end goal, the steps you will take, and the documents you will consult (confirm you will use `/docs/requirements/main.md`, `/docs/requirements/milestones.md`, and `templates/planning-phase.input.yaml`).
- Ask the user to confirm the milestone number (e.g. `m1`). All interview transcripts and the final requirements file will be placed under `docs/implementation-plan/m{n}/` (replace `{n}` with the milestone slug).

Interview rules
1. Ask one question at a time. Wait for the user's answer before asking the next question.
2. Record every question and answer verbatim in `docs/implementation-plan/m{n}/requirements-interview.md` using a simple Q/A markdown format:

   Q: <question text>

   A: <answer text>

3. After each answer, confirm the captured answer back in one sentence ("Captured: ...") and ask the next logical question.
4. Continue until you have at least the required fields from `templates/planning-phase.input.yaml` (see list below) or the stakeholder says they are done.
5. When finished, generate a filled `planning-phase.input.yaml` document using the collected answers and save it to `docs/implementation-plan/m{n}/requirements.md` (YAML or JSON is acceptable — YAML preferred). Also copy a usage example command showing how to run `your-tool plan` on the saved file.

Required fields to gather (map these to `planning-phase.input.yaml`):
- id (short slug)
- title (one-line)
- short_description
- requirements_file (path — can point to this saved requirements.md)
- style_anchors (1-3 files or references)
- stakeholders.product_owner (name/email or github handle)
- objectives (1-5)
- success_metrics (metric + target)

Helpful additional fields (ask when available):
- background
- constraints and assumptions
- scope: in_scope / out_of_scope
- deliverables
- timeline.start_date and timeline.target_completion
- testing expectations
- data_classification / secrets_included

Suggested question sequence (ask in order, but follow up when answers need clarification):
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
- Once you have captured the required fields, produce a filled `planning-phase.input.yaml` (YAML) using the collected answers.
- Save the file to `docs/implementation-plan/m{n}/requirements.md` and also save the interview transcript to `docs/implementation-plan/m{n}/requirements-interview.md`.
- Append a short `usage` section in the saved requirements.md showing the `your-tool plan` command that consumes it, e.g.:

```bash
# Generate a candidate plan (code path)
your-tool plan docs/implementation-plan/m{n}/requirements.md --method code --output planning/milestones/{id}.ralphy.yaml
```

Example: if the milestone is `m1`, save to:
- `docs/implementation-plan/m1/requirements-interview.md`
- `docs/implementation-plan/m1/requirements.md`

Notes
- Keep questions short and actionable; prefer affirmative phrasing.
- If the stakeholder is unsure about a field, capture the assumption as a note in the interview transcript and mark the field as "assumption: <text>" in the final YAML.
- All final artifacts should be human-readable and committed under `docs/implementation-plan/` for traceability.

---

This is the canonical template. Use it from `templates/requirements-prompt.md` in code or CI.
