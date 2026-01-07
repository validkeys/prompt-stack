# Diagnose Test Failure

A prompt decoration that enforces a disciplined approach to diagnosing failing tests, preventing premature solutions and ensuring thorough root cause analysis.

## Purpose

Use this prompt decoration when debugging conversations become chaotic or when the AI (or you) starts jumping to solutions without proper diagnosis. Inject this mid-conversation to reset the approach and enforce systematic investigation.

## The Prompt

```
** Be systematic about debugging **

First: State the current problem -- whether it's an error or otherwise.

Then:

1. Isolate a failing test with .only
2. Add comprehensive diagnostic console logging to the entire call stack
3. Run the test
4. Trace the diagnostic logging
5. Determine the root cause.
6. Present the root cause to the user along with a proposed solution. If there are several options, present the options for solving.
```

## Usage Example

Paste this prompt into an ongoing debugging conversation when:
- Multiple attempted fixes haven't worked
- The AI is proposing solutions without understanding the root cause
- You're getting lost in the call stack or error messages
- Debugging has become trial-and-error rather than methodical

The prompt will redirect the conversation toward systematic investigation before any fixes are attempted.

## Output Format

After injecting this prompt, expect the AI to:
1. Clearly articulate the current problem state
2. Isolate a single failing test case
3. Add detailed logging throughout the relevant code paths
4. Execute and analyze the logs methodically
5. Present findings with root cause analysis
6. Offer solution options only after diagnosis is complete

## Tips

- Use `.only` (Jest/Mocha) or equivalent to focus on one failing test at a time
- Comprehensive logging means logging: inputs, outputs, intermediate values, and control flow
- Don't skip the logging step - even when you think you know the issue
- Trace logs in execution order to understand the actual flow (not assumed flow)
- Present multiple solution options when trade-offs exist
- This works best with test-driven debugging - if no tests exist, create a minimal reproduction first
- Can be combined with the bug reporting workflow for more complex issues
