**Preventing AI Code Generation Drift in Multi-Step Development**

````markdown
# Preventing AI code generation drift in multi-step development

Your workflow—generating implementation plans, creating tracking documents with cross-references, committing incrementally, and asking models to quote back requirements—is already sophisticated. But research reveals several critical techniques you're likely missing that explain persistent drift.

## The core problem: AI models optimize locally, not globally

AI code generation models are autoregressive—they optimize for **locally plausible tokens**, not global architectural consistency. Without extremely tight constraints and rich context, each step forces the model to "invent" structure from generic patterns rather than your project's established conventions. This fundamental limitation explains why drift occurs even with good documentation: the model isn't truly "remembering" your architecture—it's pattern-matching against whatever context fits in its attention window.

The solution isn't more documentation. It's **constrained execution with multiple verification layers**.

## Critical gaps in your current approach

### You're missing style anchors

The single highest-leverage technique identified in the research: reference **2-3 exemplary files** in every prompt as pattern templates. ThoughtWorks Radar recommends "anchoring coding agents to a reference application"—a complete, realistic template the AI can pattern-match against. Instead of describing patterns abstractly in your tracking document, include explicit file references:

```
When implementing UserService, follow the exact patterns in:
@src/services/AuthService.ts  # Service structure
@src/api/healthcheck.ts       # API endpoint pattern
Touch ONLY these files: src/services/UserService.ts, src/services/index.ts
```

Henrik Jernevad's key insight: "Provide idiomatic examples. Refer to 2–3 high-quality snippets from the codebase as style anchors." Your tracking documents describe what to build—but style anchors show *how* it should look.

### Your tasks are probably too large

Research converges on a critical constraint: tasks should be executable in **30 minutes to 2.5 hours maximum**. Larger steps force the model to invent structure from generic patterns. Break your implementation plan into smaller atomic units, each touching minimal files with explicit boundaries:

```markdown
### Task 4.2: Add validation to CreateUserHandler
SCOPE: src/handlers/CreateUserHandler.ts ONLY
PRESERVE: All public APIs unchanged
DO NOT: Add dependencies, modify other files
REFERENCE: Follow validation pattern in @src/handlers/UpdateUserHandler.ts:L45-L67
VERIFY: Must pass existing test suite before proceeding
```

If drift appears (unexpected dependencies, unfamiliar patterns), **stop and revert immediately**. Version control becomes your safety net—commit after every small task.

### You're using negative constraints instead of affirmative instructions

Research shows stating what TO do performs significantly better than what NOT to do. Instead of "don't add unnecessary dependencies," use "only use existing project dependencies: express, zod, prisma." The model follows affirmative framing more reliably than prohibitions.

Transform your tracking documents from "do not modify files outside /src/services/" to "touch ONLY: src/services/UserService.ts, src/services/index.ts."

### You lack a tiered rules architecture

Experienced Cursor and Claude Code users implement **three-tier rule systems** that persist across sessions:

- **Global rules** (user settings): Universal preferences—output format, language, response length
- **Project rules** (.cursor/rules/ or CLAUDE.md): Architecture patterns, coding standards, strict TypeScript requirements
- **Context-aware rules**: Auto-attached based on file patterns (triggered only when AI touches specific directories)

Your tracking documents exist per-task, but you may lack persistent project-level constraints that apply to every interaction. Create a CLAUDE.md or .cursor/rules/ directory that's always loaded.

## Advanced constraint enforcement for TypeScript

AI models notoriously reach for shortcuts—`any` types, disabled ESLint rules, `@ts-ignore` comments. Steve Kinney observes: "Almost all models will quickly reach for disabling ESLint rules—either for a particular line or wholesale, skipping tests rather than fixing them, and liberally using `any`."

### Multi-layer enforcement stack

Prompt-level rules alone won't prevent drift. Implement enforcement at every layer:

**Layer 1 - Prompt**: Explicit rules in your CLAUDE.md or cursor rules:
```markdown
NEVER use `any`. Use `unknown` if absolutely necessary.
NEVER disable ESLint rules inline without explanation.
NEVER use type assertions (`as`) for external data—use Zod schemas.
ALWAYS explicitly type function parameters and return types.
```

**Layer 2 - IDE**: ESLint with `@typescript-eslint/strictTypeChecked` config plus these critical rules:
```javascript
'@typescript-eslint/no-explicit-any': 'error',
'@typescript-eslint/no-unsafe-argument': 'error',
'@typescript-eslint/no-unsafe-assignment': 'error',
'@typescript-eslint/no-unsafe-return': 'error',
'@typescript-eslint/explicit-function-return-type': 'error',
```

**Layer 3 - Commit**: Pre-commit hooks with Husky/lint-staged using `--max-warnings 0`. This is critical—most configs only *warn* about `any` types; zero-warning tolerance forces fixes.

**Layer 4 - CI**: Quality gates that count `eslint-disable` comments and fail if they exceed a threshold. AI-generated code frequently introduces these; automated tracking catches drift.

**Layer 5 - Runtime**: Zod schemas for all external data instead of type assertions. This is a key gap—`as User` provides compile-time safety without runtime validation. AI models default to type assertions unless explicitly constrained.

## Test-driven development as the missing anchor

TDD is underutilized as a constraint mechanism for AI code generation. Tests serve as **natural language specifications** that constrain generation more effectively than prose requirements.

### The TDD plan generation pattern

Before implementation, ask the model to generate a step-by-step TDD plan:

```
I'm implementing UserService based on these requirements.
Break this into a Test-Driven Development flow. Output as markdown checklist:
[ ] Write test for basic user creation
[ ] Implement minimal code to pass
[ ] Write test for validation (invalid email)
[ ] Update implementation
[ ] Add test for edge case (duplicate user)
[ ] Refactor for clarity

Check boxes as you progress. DO NOT FORGET THIS.
```

Research from ACM confirms: providing tests WITH problem statements consistently improves code generation success. Test names become implicit prompts—`it('should return only valid emails from a mixed list')` guides implementation better than prose descriptions.

### Using test failures as feedback loops

When tests fail, feed failures back as explicit constraints:
```
Unit test TC-003 is failing. [Include error output]
Revise implementation to pass this test while maintaining ALL previously passing tests.
Do not modify the test. Do not add new dependencies.
```

This creates a tight verification loop where human defines expected behavior; AI generates implementation; tests verify alignment.

## Model-specific strategies matter significantly

The models you're using have meaningfully different behaviors for multi-step coding:

### Claude models excel at surgical, consistent edits

Claude 4 models show **65% less likelihood to engage in shortcuts/loopholes** than earlier versions. Sourcegraph describes Claude as "staying on track longer, understanding problems more deeply." Claude is preferred for surgical fixes—one file, one bug—with minimal drift.

**Best practices for Claude**:
- Use the explicit research → plan → implement workflow—without instruction, Claude "tends to jump straight to coding"
- Leverage CLAUDE.md files for persistent project context—repository-specific optimization yields **+11% better code** in benchmarks
- For complex problems, use trigger words for extended thinking: "think" < "think hard" < "think harder" < "ultrathink"
- Prompt for "minimal diff, no renames, explain each edit"

### GPT models are more aggressive but context-degraded

GPT-4.1's accuracy **drops from 84% with 8K tokens to 50% with 1M tokens**—significant degradation with context size. GPT models tend toward larger refactors with more context. Developer observations: "GPT-4 sometimes ignores instructions or gives code that doesn't work."

**Best practices for GPT**:
- Better for exploratory/greenfield work than precision edits
- Use for code **review** rather than generation
- Prompt for "tactical plan, check side-effects in modules X/Y"

### The "lost in the middle" problem

Stanford/Berkeley research reveals a distinctive U-shaped curve: models are much better at using information at the **beginning and end** of context. Performance drops by more than 20% for information in the middle.

**Positioning strategy**: Place your most critical specifications (architectural constraints, style anchors, strict rules) at the **beginning** and reiterate key constraints at the **end** of prompts. Don't bury requirements in the middle of your tracking document.

## Self-consistency and AI-on-AI review

Two techniques that dramatically reduce drift:

### Universal self-consistency

Generate multiple solutions, then have the model select the most consistent:
1. Generate 3+ implementations with slightly higher temperature
2. Bundle responses into a new prompt
3. Ask: "Which implementation is most consistent with the requirements and existing codebase patterns? Select one and explain why."

This catches drift before it reaches your codebase.

### Multi-model review

Use different models to review each other's output. Addy Osmani's 2026 workflow: "I might have Claude write the code and then ask Gemini, 'Can you review this function for any errors or improvements?' This can catch subtle issues."

The pattern: Claude for code generation (precision), GPT for code review (catching bugs). This dynamic has remained consistent for approximately a year in practitioner reports.

## Practical implementation checklist

These are the specific techniques to add to your current workflow:

**Immediate (this week)**:
- Create CLAUDE.md or .cursor/rules/ with persistent project constraints
- Add 2-3 style anchor files to every implementation prompt
- Restructure tasks to 30min-2.5hr scope with explicit file boundaries
- Convert negative constraints to affirmative instructions
- Configure ESLint with `--max-warnings 0` and no-unsafe-* rules

**Short-term (this month)**:
- Implement TDD plan generation before each feature
- Position critical specs at beginning AND end of prompts
- Add pre-commit hooks with zero-warning tolerance
- Create Zod schemas for all external data sources
- Set up CI quality gates counting eslint-disable comments

**Ongoing practice**:
- Use AI-on-AI review (Claude generates, GPT reviews)
- Generate multiple solutions and vote on consistency
- Revert immediately when drift appears—don't try to fix mid-stream
- At end of sessions, document learnings in persistent rules files

## Conclusion

Your current workflow addresses documentation but not constraint enforcement. The key insight: AI models don't truly maintain context—they pattern-match against whatever's in their attention window. The solution is multi-layer enforcement (prompt → IDE → commit → CI → runtime), anchoring to reference implementations rather than prose descriptions, TDD as specification, and leveraging different models for different tasks.

The gap isn't in your documentation—it's in your verification and enforcement mechanisms. Style anchors, tiered rules architecture, self-consistency checking, and zero-tolerance linting gates transform AI code generation from "hopeful" to "constrained." The most sophisticated teams treat AI like a junior developer with exceptional speed but poor judgment: tight guardrails, small tasks, immediate verification, and multiple review layers.
````