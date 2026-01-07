# Phase 3: AI Features Integration - Testing Guide

## Overview

Phase 3 implements AI-powered suggestions for prompt compositions using the Claude API. This document provides a comprehensive testing guide for all AI features.

## Implementation Summary

### What Was Implemented

1. **AI Client Initialization** ([`internal/bootstrap/bootstrap.go`](internal/bootstrap/bootstrap.go:236))
   - AI client initialized from config (`ClaudeAPIKey` and `Model` fields)
   - Graceful degradation if API key not configured
   - Context selector for intelligent prompt selection

2. **App Model Integration** ([`ui/app/model.go`](ui/app/model.go:41))
   - Added `aiClient` and `contextSelector` fields
   - New message types for AI workflow:
     - `TriggerAISuggestionsMsg` - Triggers suggestion generation
     - `AISuggestionsGeneratedMsg` - Delivers generated suggestions
     - `AISuggestionsErrorMsg` - Handles AI errors

3. **Suggestion Generation Workflow** ([`ui/app/model.go:generateAISuggestions`](ui/app/model.go:497))
   - Extracts keywords from composition
   - Scores library prompts based on relevance
   - Selects top prompts within token budget
   - Sends request to Claude API
   - Parses structured suggestions response

4. **Token Budget Enforcement** ([`ui/app/model.go:259`](ui/app/model.go:259))
   - Checks composition against 25% context limit
   - Blocks suggestions if budget exceeded
   - Shows warning at 15% threshold

5. **Command Integration** ([`internal/commands/core.go:50`](internal/commands/core.go:50))
   - "Get AI Suggestions" command in command palette
   - Triggers `TriggerAISuggestionsMsg` when executed

6. **UI Integration** ([`ui/app/model.go:273`](ui/app/model.go:273))
   - Displays generated suggestions in split view
   - Shows status messages for AI operations
   - Handles errors gracefully

## Testing Prerequisites

### 1. API Key Configuration

The application requires a valid Anthropic API key to test AI features.

**Setup:**
```bash
# Edit config file
~/.promptstack/config.yaml
```

**Required fields:**
```yaml
claude_api_key: "your-api-key-here"
model: "claude-3-sonnet-20240229"
```

**Get API Key:**
- Visit https://console.anthropic.com/
- Create an account or sign in
- Navigate to API Keys section
- Generate a new API key
- Copy and paste into config file

### 2. Library Prompts

AI suggestions work best with a populated library of prompts.

**Check library:**
```bash
# View library contents
ls ~/.promptstack/data/
```

**Expected structure:**
```
~/.promptstack/data/
├── commands/
│   ├── code-review.md
│   ├── create-commit.md
│   └── ...
├── prompt-decorations/
│   └── ...
└── workflows/
    └── ...
```

**If library is empty:**
- Starter prompts are extracted on first run
- Check logs for extraction errors: `~/.promptstack/debug.log`

## Test Scenarios

### Test 1: Basic AI Suggestion Generation

**Objective:** Verify AI suggestions can be generated for a simple composition.

**Steps:**
1. Launch application: `cd cmd/promptstack && go run .`
2. Type a simple composition in workspace:
   ```
   Write a function to sort an array of integers in ascending order.
   ```
3. Press `Ctrl+P` to open command palette
4. Type "ai" to filter commands
5. Select "Get AI Suggestions" and press Enter
6. Wait for AI response (may take 5-10 seconds)

**Expected Results:**
- Status bar shows "Generated X suggestion(s)"
- Suggestions panel appears on right side (30% of screen)
- Each suggestion shows:
  - Icon (💡, 🔍, 📝, ⚠️, 🔍, 🔄)
  - Title
  - Description (truncated if long)
  - Type label (e.g., "Formatting • 2 changes")
- Workspace remains editable

**Success Criteria:**
- ✅ No errors in status bar
- ✅ Suggestions panel visible
- ✅ At least 1 suggestion generated
- ✅ Suggestions are relevant to composition

**Failure Indicators:**
- ❌ "AI client not initialized" - Check API key in config
- ❌ "Failed to generate suggestions" - Check debug.log for details
- ❌ "Composition exceeds token budget" - Reduce composition length

### Test 2: Token Budget Enforcement

**Objective:** Verify token budget limits are enforced correctly.

**Steps:**
1. Create a very long composition (copy-paste multiple times):
   ```
   Write a function to sort an array of integers in ascending order.
   [Repeat 100+ times]
   ```
2. Press `Ctrl+P` → "Get AI Suggestions"

**Expected Results:**
- Status bar shows warning or error message
- No suggestions generated if budget exceeded

**Test Cases:**

| Composition Size | Expected Behavior |
|----------------|------------------|
| < 50K tokens | Suggestions generated normally |
| ~ 50K tokens | Warning: "approaching token budget" |
| > 50K tokens | Error: "exceeds token budget" |

**Success Criteria:**
- ✅ Warning shown at 15% threshold
- ✅ Block at 25% threshold
- ✅ Clear error messages

### Test 3: Suggestion Application

**Objective:** Verify suggestions can be applied via diff viewer.

**Steps:**
1. Generate suggestions (Test 1)
2. Navigate suggestions with arrow keys or j/k
3. Press `a` to apply selected suggestion
4. Review diff in modal
5. Press Enter to accept or Esc to reject

**Expected Results:**
- Diff viewer modal appears
- Shows unified diff format:
  ```
  --- original (X lines)
  +++ new (Y lines)
  @@ -line,count +line,count @@
  -old line
  +new line
  ```
- Workspace locked (read-only) during review
- Status shows "✨ AI is applying suggestion..."

**Accept Flow:**
1. Press Enter to accept
2. Changes applied to workspace
3. Workspace unlocked
4. Suggestion marked as "[Applied ✓]"
5. Changes can be undone with Ctrl+Z

**Reject Flow:**
1. Press Esc to reject
2. Changes discarded
3. Workspace unlocked
4. Suggestion marked as "[Dismissed]"

**Success Criteria:**
- ✅ Diff viewer shows correct changes
- ✅ Accept applies changes correctly
- ✅ Reject discards changes
- ✅ Undo works after acceptance

### Test 4: Suggestion Dismissal

**Objective:** Verify suggestions can be dismissed without applying.

**Steps:**
1. Generate suggestions
2. Select a suggestion
3. Press `d` to dismiss

**Expected Results:**
- Suggestion marked as "[Dismissed]"
- No changes to workspace
- Suggestion remains in list but not applicable

**Success Criteria:**
- ✅ Dismissed suggestions show status
- ✅ Cannot apply dismissed suggestions
- ✅ Workspace unchanged

### Test 5: Error Handling

**Objective:** Verify graceful error handling for various failure scenarios.

**Test Cases:**

**A. Invalid API Key:**
1. Set invalid API key in config
2. Try to generate suggestions
3. Expected: "AI client not initialized" or authentication error

**B. Network Error:**
1. Disconnect network
2. Try to generate suggestions
3. Expected: "Failed to generate suggestions: network error"
4. Reconnect and retry

**C. API Rate Limit:**
1. Generate suggestions rapidly (10+ times)
2. Expected: Rate limit error with retry option
3. Wait and retry succeeds

**D. Malformed Response:**
1. (Simulated) AI returns invalid JSON
2. Expected: "Failed to parse suggestions" error
3. No suggestions displayed

**Success Criteria:**
- ✅ All errors show clear messages
- ✅ Application doesn't crash
- ✅ Workspace remains functional
- ✅ Errors logged to debug.log

### Test 6: Context Selection

**Objective:** Verify relevant library prompts are selected as context.

**Steps:**
1. Create composition about "code review"
2. Ensure library has code-review prompts
3. Generate suggestions
4. Check debug.log for context selection

**Expected Results:**
- Code review prompts scored higher
- Top 3-5 prompts selected
- Context fits within token budget

**Verification:**
```bash
# Check logs for scoring
grep "ScorePrompts" ~/.promptstack/debug.log
grep "SelectTopPrompts" ~/.promptstack/debug.log
```

**Success Criteria:**
- ✅ Relevant prompts selected
- ✅ Scoring algorithm works
- ✅ Token budget respected

### Test 7: Multiple Suggestion Types

**Objective:** Verify all suggestion types can be generated.

**Suggestion Types:**
1. 💡 **Recommendation** - Suggest relevant library prompts
2. 🔍 **Gap** - Identify missing context
3. 📝 **Formatting** - Suggest better structure
4. ⚠️ **Contradiction** - Identify conflicts
5. 🔍 **Clarity** - Point out ambiguity
6. 🔄 **Reformatting** - Alternative structures

**Steps:**
1. Create various compositions:
   - Simple: "Hello world"
   - Complex: Multi-step task
   - Ambiguous: "Fix the bug"
   - Conflicting: "Make it fast and secure"
2. Generate suggestions for each
3. Observe suggestion types

**Expected Results:**
- Different compositions yield different suggestion types
- Each type has appropriate icon and label
- Suggestions are actionable

**Success Criteria:**
- ✅ All 6 types can appear
- ✅ Types match composition needs
- ✅ Icons and labels correct

### Test 8: UI Responsiveness

**Objective:** Verify UI remains responsive during AI operations.

**Steps:**
1. Generate suggestions (long operation)
2. While waiting:
   - Try to edit workspace
   - Try to navigate suggestions
   - Try to open command palette
3. After suggestions arrive:
   - Navigate suggestions
   - Apply a suggestion

**Expected Results:**
- Workspace editable during generation
- UI doesn't freeze
- Status updates show progress

**Success Criteria:**
- ✅ No UI blocking
- ✅ Smooth animations
- ✅ Responsive keyboard input

## Performance Benchmarks

### Expected Performance

| Operation | Target | Acceptable |
|------------|---------|-------------|
| AI client init | < 100ms | < 500ms |
| Keyword extraction | < 50ms | < 200ms |
| Prompt scoring | < 100ms | < 500ms |
| API request | 5-10s | < 30s |
| Response parsing | < 50ms | < 200ms |
| Total workflow | 5-15s | < 45s |

### Measuring Performance

```bash
# Check timing in logs
grep "AI client initialized" ~/.promptstack/debug.log
grep "Generated.*suggestion" ~/.promptstack/debug.log
```

## Debugging

### Enable Debug Logging

```bash
# Set environment variable
export PROMPTSTACK_DEBUG=1

# Run application
cd cmd/promptstack && go run .
```

### Check Logs

```bash
# View recent logs
tail -f ~/.promptstack/debug.log

# Search for AI-related logs
grep -i "ai\|suggestion\|claude" ~/.promptstack/debug.log
```

### Common Issues

**Issue:** "AI client not initialized"
- **Cause:** Missing or invalid API key
- **Fix:** Check `~/.promptstack/config.yaml` for `claude_api_key`

**Issue:** "Failed to generate suggestions"
- **Cause:** Network error, API error, or invalid response
- **Fix:** Check debug.log for details, verify network, check API key

**Issue:** "Composition exceeds token budget"
- **Cause:** Composition too long (> 50K tokens)
- **Fix:** Reduce composition length or split into smaller parts

**Issue:** No suggestions generated
- **Cause:** AI returned empty list or parsing failed
- **Fix:** Check debug.log, try different composition

## Integration Testing

### End-to-End Workflow

**Complete User Journey:**
1. Launch application
2. Write composition
3. Insert library prompt (Ctrl+B)
4. Fill placeholders
5. Generate AI suggestions (Ctrl+P → "Get AI Suggestions")
6. Review suggestions
7. Apply relevant suggestion
8. Edit result
9. Save composition
10. View in history (Ctrl+H)

**Success Criteria:**
- ✅ All features work together
- ✅ No crashes or errors
- ✅ Smooth user experience
- ✅ Data persists correctly

## Known Limitations

1. **Token Estimation:** Uses rough approximation (~4 chars = 1 token)
2. **Context Selection:** Basic keyword matching, no semantic analysis
3. **Suggestion Quality:** Depends on AI model and prompt quality
4. **Offline Mode:** AI features require internet connection
5. **Rate Limiting:** API has rate limits (check Anthropic docs)

## Next Steps

After Phase 3 testing is complete:

1. **Phase 4:** Prompt Management (create, edit, validate prompts)
2. **Phase 5:** Vim Mode (universal vim keybindings)
3. **Phase 6:** Polish & Settings (settings panel, status bar)
4. **Phase 7:** Testing & Documentation (comprehensive tests, docs)

## Test Results Template

```markdown
## Test Results - [Date]

### Environment
- OS: [macOS/Linux/Windows]
- Go Version: [go version]
- API Model: [claude-3-sonnet-20240229]
- Library Size: [X prompts]

### Test Results

| Test | Status | Notes |
|------|--------|-------|
| Basic AI Suggestion Generation | ✅/❌ | |
| Token Budget Enforcement | ✅/❌ | |
| Suggestion Application | ✅/❌ | |
| Suggestion Dismissal | ✅/❌ | |
| Error Handling | ✅/❌ | |
| Context Selection | ✅/❌ | |
| Multiple Suggestion Types | ✅/❌ | |
| UI Responsiveness | ✅/❌ | |

### Issues Found
1. [Description]
   - Severity: [Low/Medium/High]
   - Steps to reproduce: [...]
   - Expected: [...]
   - Actual: [...]

### Performance Metrics
- AI client init: [X]ms
- Keyword extraction: [X]ms
- Prompt scoring: [X]ms
- API request: [X]s
- Total workflow: [X]s

### Overall Assessment
- Phase 3 Status: [Complete/Partial/Failed]
- Ready for Phase 4: [Yes/No]
```

## Conclusion

Phase 3 AI Features Integration provides a solid foundation for AI-powered prompt composition. The implementation includes:

- ✅ AI client initialization and configuration
- ✅ Intelligent context selection from library
- ✅ Token budget enforcement
- ✅ Structured suggestion generation
- ✅ Diff-based suggestion application
- ✅ Comprehensive error handling
- ✅ UI integration with split view

The application is now ready for Phase 4: Prompt Management.