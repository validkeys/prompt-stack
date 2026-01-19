# Style Markers — exemplar Go projects

This document extracts concrete "style markers" we can reuse as anchors when asking models or writing code for Go projects. It uses three exemplar projects as templates: `kubernetes/kubernetes`, `gohugoio/hugo`, and `junegunn/fzf`.

Use these markers when creating prompts, PR checklists, or repository scaffolding. Follow the AI Codegen Best Practices (docs/best-practices.md) when using them: attach 2–3 exemplary files to prompts, reference exact paths, and scope tasks to small commits.

- Templates to include in prompts (touch ONLY these when asking for examples):
  - `kubernetes/kubernetes/README.md`
  - `gohugoio/hugo/README.md`
  - `junegunn/fzf/README.md`

What a "style marker" is

- Short, opinionated pattern you can point a model at to influence structure, naming, tests, and docs.
- Examples: single-binary CLI layout, Makefile-driven build targets, comprehensive README sections, shell-integration scripts, modular package layout, or strict CI checks.

How to use

1. Attach 2–3 exemplar files to the prompt (see the Templates above). The model should be told: `USE these files as style anchors, NOT to copy large blocks verbatim`.
2. Reference exact target paths and `touch ONLY` the example file(s) you want the model to mirror.
3. Break the work into 30m–2.5h tasks and commit after each task (see docs/best-practices.md Task sizing).
4. Ask for minimal diffs, explicit rationale for any deviation, and a short test plan.

Selected markers from each exemplar

- Kubernetes (`kubernetes/kubernetes`)
  - Project scale & governance: explicit contributor docs and community links in the README — keep README actionable (how to build, support links, developer start).
  - Makefile-first developer workflow: clear `make` targets for build/test/release used as developer affordances.
  - Modular layout: large codebase split into well-named packages and a staging area (treat package boundaries as API surfaces).
  - Explicit support notes: clearly state which packages are intended to be consumed as libraries and which are internal (document boundaries).

- Hugo (`gohugoio/hugo`)
  - Single-command UX with subcommands: `main` + `commands` style, clear CLI usage and examples in README.
  - Pluggable templates and content-first documentation: example content, themes, and a thorough examples section to show real usage.
  - Release & build instructions for multiple platforms (Go tooling + prebuilt binaries) — prefer reproducible build steps.
  - Friendly developer onboarding: simple `go install` / `make` steps and a short "to start developing" section.

- fzf (`junegunn/fzf`)
  - Portability & single-binary distribution: distributing a small portable binary and small install scripts (single-file install) is a deliberate design choice.
  - Shell integrations and examples: provide integration scripts (bin/), usage snippets, and recommended environment variables to customize behavior.
  - UX-first README: concise highlights, examples, and tips; short, executable examples that readers can copy/paste.
  - Config-driven defaults: environment variables and default options (`FZF_DEFAULT_OPTS`) documented with clear examples.

Concrete prompt snippets (copy into prompts)

- Minimal ask with anchors:

```text
Use these style anchors: kubernetes/kubernetes/README.md, gohugoio/hugo/README.md, junegunn/fzf/README.md
Create a single-binary CLI under cmd/mytool with:
- a lean README example showing install + one example
- a Makefile target: `make build` and `make test`
Touch ONLY: cmd/mytool/main.go, README.md, Makefile
Scope: complete in one 90m task. Provide tests.
```

- Enforce constraints with affirmative instructions:

```text
ONLY use: Go modules, standard library, and `spf13/cobra` for CLI (if necessary).
Do NOT add other dependencies. Add unit tests for exported packages. Follow the style markers above.
```

Checklist for reviewers / CI

- README includes: quickstart, build, contribute, support links (use the Kubernetes README pattern).
- CLI has: `--help`, clear examples, and at least one integration snippet (use the Hugo/fzf style).
- Build: `make build` and `make test` pass locally; binaries reproducible via Go modules.
- Packaging: if distributing a single binary, include an install script or instructions (`bin/` scripts like fzf).
- Docs: point out internal vs. external packages and include a short note about public API surfaces.

Next steps

1. Review `docs/style-markers.md` and confirm the three template files listed are the ones you want attached to prompts.
2. If yes, we can: (a) add short example files into `examples/style-anchor/` in this repo (one-liners or trimmed READMEs) to include in prompts; or (b) keep using remote repo files as references.
3. Optionally, I can create `examples/style-anchor/` with 2–3 minimal template files mirroring the markers (choose 1 or 2).

File created: `docs/style-markers.md`
