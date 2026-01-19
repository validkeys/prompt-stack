# Requirements Gathering Prompt Template

Use this template interactively with a stakeholder or paste into an editor and fill in answers. Save the filled output as `examples/requirements/inputs/requirements.input.md`.

## Prompt (Template)

Project name:
Short description (one sentence):
Primary user persona(s):
Primary goals (3–5):
Success metrics (3):
Non-goals (what we won't do):
Critical constraints (security, latency, cost, infra):
Must-have integrations (e.g., Git, CI, DB, third-party services):
Acceptable tech stack (optional):
Privacy / secrets handling preferences:
Initial scope (what will be delivered in milestone 1):
Out-of-scope (for milestone 1):
Initial non-functional targets (latency, size, performance):
Regulatory/compliance concerns (if any):
Testing approach (unit, integration, TDD?):
Definition of done for milestone 1 (acceptance criteria):
Notes / examples / reference links:

---

## Example usage

1. Copy this file to `examples/requirements/inputs/requirements.input.md` and fill it in.
2. Run `your-tool plan examples/requirements/inputs/requirements.input.md --method code` to generate a `tasks.yaml` candidate.

Save the filled input alongside your project planning artifacts so Plan Mode can consume it directly.
