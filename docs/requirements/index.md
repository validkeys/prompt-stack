# Requirements — Index

A concise entry point for AI and humans to quickly find and inspect requirement-related documents. Each entry links to the source document and includes a one-line description for fast lookup.

- `docs/requirements/main.md` — Core business and product requirements. Primary reference for the CLI behavior, user flows, style anchors, example usage (`your-tool plan ...`), and non-functional constraints.
- `docs/requirements/architecture.md` — Compact, actionable architecture spec for the host CLI (plug-in points, lifecycle, core APIs). Distills decisions captured in the interview draft into implementable guidance.
- `docs/requirements/architecture.interview.md` — Working architecture draft captured during the one-question-at-a-time interview. Contains rationale, incremental decisions, and implementation sketches.
- `docs/requirements/milestones.md` — Test-oriented milestone plan that breaks requirements into an MVP path with deliverables and manual verification steps.
- `docs/requirements/interview.md` — Q&A tracker used to clarify business requirements; useful for rationale, open questions, and decision history.

Notes:
- Prefer `docs/requirements/main.md` for product-level questions and `docs/requirements/architecture.md` when designing or changing system components.
- For implementation rationale and incremental decisions consult `architecture.interview.md` and `interview.md`.

If you'd like, I can:
1. Add short anchors or a table-of-contents for quicker navigation.
2. Include cross-references to related research files (e.g. `docs/initial-claude-conversation/*`).
3. Commit this file on your behalf with a concise commit message.
