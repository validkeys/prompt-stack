# Bug Reporting Workflow

A structured workflow for systematically identifying, diagnosing, and fixing bugs using a multi-document approach.

## Purpose

Use this workflow when you discover a bug and want to follow a disciplined process from initial report through diagnosis, solution proposal, and implementation. This prevents jumping to solutions prematurely and ensures thorough documentation.

## The Workflow

```
I need to create a structured bug report with the following organization:

Create a folder named [bug-name or issue-X] with these files:

1. **readme.md** - Workflow checklist containing:
   - Bug report structure overview
   - Step-by-step checklist (format report → create tests → diagnose → propose solution → get approval → implement → review & enhance comments → update docs)
   - Definition of done (all tests pass, no regressions, enterprise-grade comments)
   - Patterns to follow (use surrounding tests for consistency, reference bug report in code comments)
   - Final context section (to be filled when complete)

2. **report.md** - Initial bug report with:
   - Summary
   - Expected vs actual behavior
   - Steps to reproduce
   - Impact assessment
   - Additional context

3. **diagnosis.md** - Investigation tracking with:
   - Investigation steps (what was checked and findings)
   - Root cause analysis
   - Supporting evidence (code snippets, logs)

4. **solution.md** - Solution proposal with:
   - Proposed solution description
   - Implementation approach
   - Changes required (checklist)
   - Testing strategy
   - Risks and considerations
   - Alternatives considered

5. **plan.md** - Implementation tracking with:
   - Implementation tasks (checklist)
   - Progress log with dates
   - Testing results
   - Documentation updates (checklist)

Follow this process strictly:
- Report the bug first, don't diagnose yet
- Create tests that expose the bug before fixing
- Diagnose thoroughly, adding diagnostic logging if needed, but DO NOT fix
- Present findings and ask for approval before proposing solutions
- Wait for approval on the solution before implementing
- Update plan.md as you implement
- Review implementation files and enhance comments to enterprise grade (explain decisions, edge cases, why certain approaches were chosen)
- Ensure all tests pass before marking as complete
```

## Usage Example

When you encounter a bug, ask your AI assistant to generate the bug report structure using the workflow above. Then work through each phase systematically, ensuring each document is complete before moving to the next phase.

## Output Format

The workflow creates a self-contained bug report folder with all necessary documentation to track the issue from discovery through resolution. Each document has a specific purpose and should be updated at the appropriate phase.

## Tips

- Don't skip phases - the separation prevents premature solutions
- Create tests first - they validate your fix and prevent regressions
- Add diagnostic logging during diagnosis - don't attempt fixes yet
- Wait for approval before implementing - solutions may need discussion
- Reference the bug report in code comments - helps future maintainers
- Use surrounding code patterns for consistency
- Review implementation and enhance comments after fixing - explain why decisions were made, document edge cases, and ensure comments are enterprise grade
- Update the Final Context section in readme.md when complete - creates a summary for future reference
- Keep plan.md updated in real-time during implementation
