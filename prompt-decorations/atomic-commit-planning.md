# Atomic Commit Planning

A prompt decoration that structures implementation plans into atomic, reviewable commits with clear tracking documentation.

## Purpose

Use this prompt decoration when planning a feature or refactor to ensure the implementation follows atomic commit principles. This creates a clear implementation path where each commit is independently reviewable, testable, and provides a logical step in the overall change.

## The Prompt

```
Structure this implementation plan following atomic commit principles.

For each commit:
1. Define a single, focused change that makes sense on its own
2. Ensure each commit leaves the codebase in a working state
3. Write a clear commit message describing the "why" not just the "what"
4. List the specific files that will change

Present the plan as a tracking document with:
- [ ] Commit checkboxes for progress tracking
- Clear dependencies between commits (if any)
- Estimated scope (small/medium/large) per commit
- Notes on testing approach for each commit

Each commit should be independently reviewable and should not break existing functionality.
```

## Usage Example

Use this prompt when:
- Planning a multi-step feature implementation
- Breaking down a large refactor into manageable pieces
- Creating an implementation plan that needs to be reviewed before coding
- You want clear progress tracking with meaningful milestones
- Working on changes that will go through code review
- Need to ensure rollback points if something goes wrong

## Output Format

After using this prompt, expect a structured implementation plan like:

```
## Implementation Plan: [Feature Name]

### Commit 1: [Short description]
**Scope**: Small
**Why**: [Explanation of purpose]
**Files**:
- path/to/file1.ts
- path/to/file2.ts
**Testing**: [How this commit will be tested]

- [ ] Complete implementation
- [ ] Tests passing
- [ ] Ready for review

### Commit 2: [Short description]
**Depends on**: Commit 1
**Scope**: Medium
[...]
```

## Tips

- Each commit should pass all existing tests - no "fix tests later" commits
- Commits should tell a story - someone reviewing should understand the progression
- Aim for commits that take 30 minutes to 2 hours to implement
- If a commit seems too large (>5 files), consider breaking it down further
- Dependencies between commits should be explicit - parallelizable commits save time
- Include refactoring as separate commits from feature additions when possible
- Documentation updates can be their own commit or bundled with the related code change
- The "why" in commit messages should explain the business or technical reasoning
- Mark commits as small/medium/large to help with time estimation and prioritization
- This pairs well with feature branches and pull requests for team collaboration
