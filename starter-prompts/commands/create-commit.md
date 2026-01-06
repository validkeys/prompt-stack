# Create Commit

A concise prompt for creating short, descriptive commit messages that follow project conventions without unnecessary attributions.

## Purpose

Use this prompt when you need to create a git commit with a brief, focused message that describes what was done while adhering to any commitlint rules defined in the project.

## The Prompt

```
Create a git commit with the following requirements:

1. Write a short, concise commit message that briefly describes what was done
2. If commitlint configuration exists in the project, follow those rules (conventional commits, type prefixes, etc.)
3. Focus on the "what" and "why" in a brief format
4. Do NOT include any co-author attributions or mentions of AI assistance
5. Keep the message under 72 characters for the subject line when possible
6. Use imperative mood (e.g., "Add feature" not "Added feature")

The commit should capture the essence of the changes without unnecessary verbosity or attributions.
```

## When to Use

- After completing a focused set of changes that should be committed together
- When you want a clean commit message that follows project standards
- During feature development with logical checkpoints
- When making bug fixes or small improvements
- After refactoring work that should be documented

## When NOT to Use

- When changes are incomplete or would leave the codebase in a broken state
- For work-in-progress commits that are just for backup purposes
- When you haven't reviewed the changes being committed
- If there are multiple unrelated changes that should be separate commits

## Usage Examples

### Basic Commit
```
[Paste the create commit prompt]

Please commit the authentication changes we just made.
```

### With Specific Files
```
[Paste the create commit prompt]

Commit only the changes to src/auth/ directory.
```

### Bug Fix Commit
```
[Paste the create commit prompt]

Create a commit for the pagination bug fix.
```

## Tips

- **Review Changes First**: Always run `git status` and `git diff` before committing to understand what's being committed
- **Atomic Commits**: Each commit should represent a single logical change
- **Commitlint Rules**: If your project uses commitlint, the prompt will automatically follow those conventions (e.g., `feat:`, `fix:`, `chore:`)
- **Subject Line Length**: Aim for commit messages under 72 characters for the first line
- **No Noise**: Avoid generic messages like "updates" or "changes" - be specific about what changed
- **Imperative Mood**: Use imperative verbs like "Add", "Fix", "Update", "Remove", not past tense
- **Staging**: Use `git add` to stage specific files before committing if you don't want to commit everything

## Common Commit Types

If your project uses conventional commits (commitlint), common types include:

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that don't affect code meaning (formatting, whitespace)
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **perf**: Code change that improves performance
- **test**: Adding or correcting tests
- **chore**: Changes to build process or auxiliary tools

## Examples of Good Commit Messages

```
feat: add user authentication middleware

fix: resolve pagination bug on search results

refactor: extract validation logic to separate module

docs: update API endpoint documentation

test: add unit tests for date formatting utility
```

## Related Prompts

- Use `atomic-commit-planning.md` when planning multiple commits for a larger feature
- Combine with code review prompts before committing to ensure quality
- Reference project's commitlint configuration for specific formatting rules
