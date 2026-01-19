# Repo-level Prompt-Stack Config — Role-based Model Selection

Purpose
- Describe a repository-level `prompt-stack.yaml` design where authors declare model profiles and roles (e.g. `planner`, `coder`, `reviewer`, `qa`). The agent assigns roles to tasks (heuristic or classifier) and selects models by capability/cost matching.

Design goals
- Keep config role-centric rather than task-specific so the AI can assign roles to tasks.
- Provide capability metadata per model (token limits, determinism, cost) for deterministic selection.
- Support advisory and strict policies per-role; start advisory for telemetry.
- Make it easy for Ralphy to emit task metadata (`intent`, `est_tokens`, `complexity`, optional `role_hint`).

Config model (concept)
- Top-level sections:
  - `models`: model profiles with provider, token_limit, cost_rank, deterministic, latency_hint
  - `roles`: each role defines `required` capabilities, candidate model aliases and `policy` (strict/advisory)
  - `defaults`: fallback models and global policy flags

Example (compact)
```yaml
version: 1
models:
  openai-gpt-4:
    provider: openai
    token_limit: 8192
    cost_rank: 3
    deterministic: false
  local-small:
    provider: local
    token_limit: 2048
    cost_rank: 1
roles:
  planner:
    required:
      token_limit: 4000
    candidates: [openai-gpt-4, local-small]
    policy: advisory
  coder:
    required:
      token_limit: 4096
    candidates: [openai-gpt-4]
    policy: advisory
  reviewer:
    required:
      deterministic: true
    candidates: [openai-gpt-4, local-small]
    policy: advisory
  qa:
    required:
      deterministic: true
      token_limit: 8000
    candidates: [openai-gpt-4]
    policy: advisory
defaults:
  fallback_models: [local-small]
```

Runtime selection (high level)
- Input: task object with `prompt`, `intent`, `est_tokens`, `complexity`, optional `role_hint`.
- Steps:
  1. Role = `role_hint` or `heuristic_map(task)` or `small-classifier(task)`.
  2. Candidates = config.roles[role].candidates
  3. Score each candidate by capability match − cost_penalty.
  4. Choose best available; if none, use `defaults.fallback_models`.
- Prefer deterministic scoring and log reason for audit.

Ralphy integration
- Have Ralphy emit per-task metadata in generated plan (`intent`, `est_tokens`, `complexity`, optional `role_hint`).
- Optionally, Ralphy can propose a `role` or `model_choice` for human review.

Validation & enforcement
- Provide a JSON Schema for `prompt-stack.yaml` and a small CLI validator for CI.
- `policy: strict` per-role means CI should fail if a task uses a non-allowed model.
- Record telemetry: chosen model, selection_reason, and fallback events for audits.

Rollout recommendation
- Start advisory (log-only) rollout to gather telemetry; flip to strict once stable.
- Maintain an authoritative source of model capability metadata; consider automation later.

Notes
- Heuristics (keywords, file paths, estimates) are a cheap starting point.
- A classifier can improve role inference but adds a small extra cost.
- Keep `prompt-stack.yaml` example under version control and update as providers evolve.
