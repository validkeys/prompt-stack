# Debug Implementation Code

A prompt decoration that enforces a systematic approach to debugging implementation issues, emphasizing diagnostic logging and collaborative problem-solving.

## Purpose

Use this prompt decoration when debugging implementation code issues where the AI is making assumptions, proposing untested solutions, or going in circles. Inject this mid-conversation to enforce methodical investigation and encourage asking for guidance when uncertain.

## The Prompt

```
Be systematic about debugging this issue. If it's a test failure, ensure that you isolate the test with test.only

Add as much diagnostic console logging to the entire call stack to narrow your focus and truly understand what's going on. Identify the root cause of the problem and present to me for review. I don't want to go in circles, so if you are unsure of something or are looking for guidance, please ask me.
```

## Usage Example

Paste this prompt into an ongoing debugging conversation when:
- The AI is proposing multiple fixes without understanding the root cause
- Debugging has become trial-and-error or going in circles
- You need the AI to slow down and investigate systematically
- The AI should ask you questions but isn't doing so
- Implementation code is behaving unexpectedly

The prompt redirects the conversation toward methodical investigation and collaborative problem-solving.

## Output Format

After injecting this prompt, expect the AI to:
1. Add comprehensive diagnostic logging throughout the relevant code paths
2. Run the code and analyze the logs systematically
3. Trace execution flow to understand actual behavior (not assumed behavior)
4. Identify the root cause of the issue
5. Present findings and root cause for your review before proposing solutions
6. Ask for guidance when uncertain rather than making assumptions

## Tips

- Diagnostic logging should cover: inputs, outputs, intermediate values, state changes, and control flow decisions
- For complex issues, log at multiple levels of the call stack
- Use descriptive log labels to make tracing easier (e.g., `console.log('[Component.method] input:', input)`)
- Encourage the AI to ask questions early - it saves time over incorrect assumptions
- If tests exist, isolate with `.only` to focus debugging efforts
- If no tests exist, create a minimal reproduction case first
- Compare expected vs actual behavior at each step in the logs
- This pairs well with the bug reporting workflow for tracking complex issues
- Remember: understanding the problem thoroughly is faster than multiple failed fix attempts
