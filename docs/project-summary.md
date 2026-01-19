# Project Summary

## What We're Building

A dual-mode AI-assisted development workflow tool that generates and validates perfect Ralphy YAML files, then executes them with guaranteed quality.

**Plan Mode**: Creates validated Ralphy YAML from requirements through three generation methods:
- **Code Generation** (Fast Path - 2-5s): Template-based using cached knowledge base
- **AI Generation** (Quality Path - 30-60s): Ralphy-powered for complex requirements  
- **Hybrid** (Smart Path - 5-30s): Fast generation → AI validation → AI review (recommended)

**Build Mode**: Executes validated YAML through Ralphy/OpenCode with pre-flight checks, parallel agent orchestration, and automatic learning.

## Core Innovation

Self-validating, self-improving YAML generation with multi-layer validation:

```
Requirements → Plan Mode → YAML
              ↓ (Code: 2-5s / AI: 30-60s / Hybrid: 5-30s)
        [Code Generation + AI Validation + AI Review]
              ↓
        Perfect YAML → Build Mode → Quality Code
```

## Key Benefits

- **90% context reduction** through AST-based filtering before sending to AI
- **95% first-pass success** via research-backed constraints preventing drift
- **10-20x faster** than manual YAML authoring + debugging
- **Team knowledge sharing** with cached patterns benefiting entire organization

## Research-Backed Best Practices

Every generated YAML embeds 10 proven techniques from drift-prevention research:
1. Style anchors (2-3 per task) - 40%+ quality improvement
2. Task sizing (30min-2.5hr max)
3. Affirmative constraints (what TO do, not what NOT to do)
4. Multi-layer enforcement (prompt → IDE → commit → CI → runtime)
5. TypeScript strict mode with forbidden/required patterns
6. TDD workflow
7. Critical context positioning (beginning/end, not middle)
8. Context budget management (5,000 tokens per task, 80-95% reduction)
9. Model-specific strategies (Claude for precision, GPT for review)
10. Self-consistency checking (3 solutions, AI votes)

## Architecture

```
YOUR TOOL (Meta-Orchestration)
  • Plan Mode: YAML generation + review
  • Build Mode: Execution orchestration
  • Knowledge management (SQLite + JIT caching)
  • Context optimization
        ↓
RALPHY (Execution Layer)
  • Parallel agent orchestration
  • Git workflow automation
        ↓
OPENCODE (Implementation Layer)
  • Code generation
```

## Tech Stack

- Runtime: Go (single static binary)
- Language: Go 1.20+
- Database: SQLite (mattn/go-sqlite3 or modernc/sqlite)
- CLI: Cobra
- AI: Direct SDKs/HTTP clients (Anthropic, OpenAI, OpenCode)
- Plugin model: External executables (stdin/stdout RPC)
- Target: Interactive CLI for senior developers

## Success Metrics

- Plan Mode: Quality score >0.9 for 95% of plans
- Build Mode: First-pass success >90%, zero drift from plan
- Overall: 3-5x developer velocity, 50%+ bug reduction, <1 hour team onboarding
