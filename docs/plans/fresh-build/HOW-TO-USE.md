# How to Use the Milestone Execution System

## Quick Start

### 1. Start a New Milestone

Copy the entire contents of [`milestone-execution-prompt.md`](milestone-execution-prompt.md:1) and paste it into a new chat with your AI assistant (in Code mode).

The AI will respond with:
```
I understand the milestone execution process. I will:
- Follow TDD strictly
- Complete one task at a time
- Stop after each task for verification
- Create comprehensive checkpoints
- Ensure build and tests always pass
- Follow all guidelines

Which milestone should I start with?
```

### 2. Specify the Milestone

Reply with:
```
Start with Milestone 1: Bootstrap & Config
```

The AI will:
1. Read [`DOCUMENT-REFERENCE-MATRIX.md`](DOCUMENT-REFERENCE-MATRIX.md:1) to identify required documents
2. Read [`milestones.md`](milestones.md:1) to get milestone details
3. Read testing guide for the milestone group (e.g., FOUNDATION-TESTING-GUIDE.md for M1-M6)
4. Read detailed acceptance criteria document if applicable (M16, M28, M32, M33, M35, M37, M38)
5. Generate `M1-task-list.md` (concise task breakdown)
6. Generate `M1-reference.md` (detailed implementation guide)
7. Wait for your approval

### 3. Review and Approve

Review the generated task list and reference document. The AI will have:
- Referenced the appropriate testing guide for the milestone group
- Referenced detailed acceptance criteria if applicable
- Applied the 5-category test criteria framework (Functional, Integration, Edge Cases, Performance, UX)

If satisfied, reply:
```
Approved. Begin Task 1.
```

### 4. Task Execution Cycle

For each task, the AI will:

1. **Write tests first** (TDD red phase)
2. **Implement code** (TDD green phase)
3. **Refactor** (TDD refactor phase)
4. **Run all checks** (build, tests, race detector, coverage)
5. **Verify guidelines** (style, testing, structure)
6. **Create checkpoint document**
7. **STOP and wait for your verification**

### 5. Verify Each Task

After each task, the AI will output:
```
✅ Task 1 Complete: Initialize Config Structure

📊 Metrics:
- Tests: 5 passing
- Coverage: 92%
- Build: ✅ Success

📄 Checkpoint: docs/plans/fresh-build/milestones/M1-checkpoints/task-1-checkpoint.md

⏸️  STOPPING for human verification.

Please review:
1. Checkpoint document
2. Code changes
3. Test results

Reply "continue" to proceed to Task 2, or provide feedback.
```

**Your options**:
- Reply `"continue"` to proceed to next task
- Provide feedback for corrections
- Request changes or clarifications

### 6. Complete the Milestone

After all tasks are complete, the AI will:
1. Run full test suite
2. Verify all enhanced test criteria (Functional, Integration, Edge Cases, Performance, UX)
3. Create milestone summary with test criteria verification
4. STOP for milestone review

### 7. Manual Testing

Follow the appropriate testing guide for the milestone group:

**Testing Guides by Milestone Group:**
- Foundation (M1-M6): [`FOUNDATION-TESTING-GUIDE.md`](milestones/FOUNDATION-TESTING-GUIDE.md)
- Library Integration (M7-M10): [`LIBRARY-INTEGRATION-TESTING-GUIDE.md`](milestones/LIBRARY-INTEGRATION-TESTING-GUIDE.md)
- Placeholders (M11-M14): [`PLACEHOLDER-TESTING-GUIDE.md`](milestones/PLACEHOLDER-TESTING-GUIDE.md)
- History (M15-M17): [`HISTORY-TESTING-GUIDE.md`](milestones/HISTORY-TESTING-GUIDE.md)
- Commands & Files (M18-M22): [`COMMANDS-FILES-TESTING-GUIDE.md`](milestones/COMMANDS-FILES-TESTING-GUIDE.md)
- Prompt Management (M23-M26): [`PROMPT-MANAGEMENT-TESTING-GUIDE.md`](milestones/PROMPT-MANAGEMENT-TESTING-GUIDE.md)
- AI Integration (M27-M33): [`AI-INTEGRATION-TESTING-GUIDE.md`](milestones/AI-INTEGRATION-TESTING-GUIDE.md)
- Vim Mode (M34-M35): [`VIM-MODE-TESTING-GUIDE.md`](milestones/VIM-MODE-TESTING-GUIDE.md)
- Polish (M36-M38): [`POLISH-TESTING-GUIDE.md`](milestones/POLISH-TESTING-GUIDE.md)

These guides provide:
- Integration tests
- End-to-end scenarios
- Performance benchmarks
- Go code examples

For complex milestones, also reference detailed acceptance criteria:
- M16: [`ACCEPTANCE-CRITERIA-M16-HISTORY-SYNC.md`](milestones/ACCEPTANCE-CRITERIA-M16-HISTORY-SYNC.md)
- M28: [`ACCEPTANCE-CRITERIA-M28-CONTEXT-SELECTION.md`](milestones/ACCEPTANCE-CRITERIA-M28-CONTEXT-SELECTION.md)
- M32: [`ACCEPTANCE-CRITERIA-M32-DIFF-GENERATION.md`](milestones/ACCEPTANCE-CRITERIA-M32-DIFF-GENERATION.md)
- M33: [`ACCEPTANCE-CRITERIA-M33-DIFF-APPLICATION.md`](milestones/ACCEPTANCE-CRITERIA-M33-DIFF-APPLICATION.md)
- M35: [`ACCEPTANCE-CRITERIA-M35-VIM-KEYBINDINGS.md`](milestones/ACCEPTANCE-CRITERIA-M35-VIM-KEYBINDINGS.md)
- M37: [`ACCEPTANCE-CRITERIA-M37-RESPONSIVE-LAYOUT.md`](milestones/ACCEPTANCE-CRITERIA-M37-RESPONSIVE-LAYOUT.md)
- M38: [`ACCEPTANCE-CRITERIA-M38-ERROR-HANDLING.md`](milestones/ACCEPTANCE-CRITERIA-M38-ERROR-HANDLING.md)

### 8. Move to Next Milestone

If satisfied with testing, reply:
```
next milestone
```

The AI will load Milestone 2 and repeat the process.

## File Structure Created

After completing Milestone 1, you'll have:

```
docs/plans/fresh-build/
├── milestones.md                          # All milestone definitions with enhanced test criteria
├── DOCUMENT-REFERENCE-MATRIX.md           # Milestone to document mapping
├── DOCUMENT-INDEX.md                     # Complete document index
├── milestone-execution-prompt.md          # Main execution workflow
└── milestones/
    ├── progress.md                        # Overall progress tracking
    ├── ENHANCED-TEST-CRITERIA-TEMPLATE.md # Acceptance criteria template
    ├── FOUNDATION-TESTING-GUIDE.md        # M1-M6 testing guide
    ├── LIBRARY-INTEGRATION-TESTING-GUIDE.md # M7-M10 testing guide
    ├── PLACEHOLDER-TESTING-GUIDE.md       # M11-M14 testing guide
    ├── HISTORY-TESTING-GUIDE.md           # M15-M17 testing guide
    ├── COMMANDS-FILES-TESTING-GUIDE.md    # M18-M22 testing guide
    ├── PROMPT-MANAGEMENT-TESTING-GUIDE.md # M23-M26 testing guide
    ├── AI-INTEGRATION-TESTING-GUIDE.md    # M27-M33 testing guide
    ├── VIM-MODE-TESTING-GUIDE.md          # M34-M35 testing guide
    ├── POLISH-TESTING-GUIDE.md            # M36-M38 testing guide
    ├── ACCEPTANCE-CRITERIA-M16-HISTORY-SYNC.md
    ├── ACCEPTANCE-CRITERIA-M28-CONTEXT-SELECTION.md
    ├── ACCEPTANCE-CRITERIA-M32-DIFF-GENERATION.md
    ├── ACCEPTANCE-CRITERIA-M33-DIFF-APPLICATION.md
    ├── ACCEPTANCE-CRITERIA-M35-VIM-KEYBINDINGS.md
    ├── ACCEPTANCE-CRITERIA-M37-RESPONSIVE-LAYOUT.md
    ├── ACCEPTANCE-CRITERIA-M38-ERROR-HANDLING.md
    ├── M1-task-list.md                    # Task breakdown
    ├── M1-reference.md                    # Implementation guide
    ├── M1-summary.md                      # Completion summary
    └── M1-checkpoints/
        ├── task-1-checkpoint.md           # Task 1 details
        ├── task-2-checkpoint.md           # Task 2 details
        └── ...
```

## Key Commands

| Command | When to Use | Effect |
|---------|-------------|--------|
| `"continue"` | After task completion | Proceed to next task |
| `"next milestone"` | After milestone completion | Load next milestone |
| `"rollback"` | If task fails critically | Revert changes and re-plan |
| `"show progress"` | Anytime | Display current progress |
| `"pause"` | Anytime | Save state and stop |

## Best Practices

### ✅ Do This

1. **Review every checkpoint** - Don't blindly approve
2. **Run tests yourself** - Verify the AI's test results
3. **Test manually** - Follow the testing guide at milestone end
4. **Provide feedback** - If something isn't right, say so
5. **Keep sessions focused** - One milestone per session is ideal
6. **Commit after each task** - Use git to track progress
7. **Read the reference docs** - They contain valuable context
8. **Consult DOCUMENT-REFERENCE-MATRIX.md** - Identifies all required documents for each milestone
9. **Use testing guides** - Provide integration tests, end-to-end scenarios, and performance benchmarks
10. **Reference detailed acceptance criteria** - For complex milestones (M16, M28, M32, M33, M35, M37, M38)

### ❌ Avoid This

1. **Don't skip verification** - Each stop point is important
2. **Don't rush** - Quality over speed
3. **Don't ignore warnings** - Address guideline violations
4. **Don't skip manual testing** - Automated tests aren't enough
5. **Don't continue with failing tests** - Fix issues immediately
6. **Don't modify the prompt mid-milestone** - Finish or restart

## Troubleshooting

### Issue: AI skips TDD
**Solution**: Remind it: "Please write tests first before implementation"

### Issue: AI doesn't stop after tasks
**Solution**: Remind it: "Stop after each task for verification"

### Issue: Tests are failing
**Solution**: 
1. Review the checkpoint document
2. Run tests yourself: `go test ./... -v`
3. If fixable, provide guidance
4. If not, request rollback

### Issue: Guidelines not followed
**Solution**: Point out specific violations and request fixes

### Issue: Lost context mid-milestone
**Solution**: 
1. Check the last checkpoint document
2. Review progress.md
3. Resume from last completed task

## Example Session Flow

```
You: [Paste milestone-execution-prompt.md]

AI: I understand the process. Which milestone should I start with?

You: Start with Milestone 1: Bootstrap & Config

AI: [Generates M1-task-list.md and M1-reference.md]
    Ready to begin. Please review and approve.

You: Approved. Begin Task 1.

AI: [Writes tests, implements code, runs checks]
    ✅ Task 1 Complete
    ⏸️  STOPPING for verification

You: [Reviews checkpoint, code, tests]
    continue

AI: [Proceeds to Task 2]
    ...

[After all tasks]

AI: 🎉 Milestone 1 Complete
    📄 Testing guide created
    ⏸️  STOPPING for milestone review

You: [Follows testing guide, manually tests]
    next milestone

AI: Loading Milestone 2: Basic TUI Shell
    [Generates M2-task-list.md and M2-reference.md]
```

## Tips for Success

1. **Set aside focused time** - Each milestone takes 2-4 hours
2. **Use git branches** - Create a branch per milestone
3. **Keep notes** - Document any deviations or learnings
4. **Test incrementally** - Don't wait until the end
5. **Trust the process** - The structure ensures quality
6. **Communicate clearly** - Be specific with feedback
7. **Celebrate progress** - Each completed task is an achievement

## Monitoring Progress

Check [`progress.md`](milestones/progress.md:1) anytime to see:
- Current milestone and task
- Overall completion percentage
- Test statistics
- Recent activity
- Metrics trends
- AC Complete status (✅ indicates acceptance criteria complete)
- Coverage percentage for each milestone

For complete documentation overview, see:
- [`DOCUMENT-INDEX.md`](DOCUMENT-INDEX.md:1) - Complete index of all documents
- [`DOCUMENT-REFERENCE-MATRIX.md`](DOCUMENT-REFERENCE-MATRIX.md:1) - Milestone to document mapping

## Getting Help

If you encounter issues:

1. **Check the checkpoint** - Contains detailed task info
2. **Review the reference doc** - Has implementation guidance
3. **Check style/testing guides** - Ensure guidelines are followed
4. **Consult DOCUMENT-REFERENCE-MATRIX.md** - Identifies all required documents
5. **Read the testing guide** - Provides integration tests and scenarios
6. **Check detailed acceptance criteria** - For complex milestones (M16, M28, M32, M33, M35, M37, M38)
7. **Ask specific questions** - The AI can clarify anything
8. **Start fresh if needed** - Sometimes a clean start helps

## Advanced Usage

### Resuming After a Break

1. Check `progress.md` for current state
2. Read the last checkpoint document
3. Tell the AI: "Resume from Task {N} of Milestone {M}"
4. The AI will pick up where you left off

### Customizing the Process

You can modify the prompt for your needs:
- Adjust checkpoint format
- Change verification requirements
- Add custom metrics
- Modify rollback strategy

Just update [`milestone-execution-prompt.md`](milestone-execution-prompt.md:1) before starting.

### Parallel Development

For multiple developers:
- Each works on different milestones
- Use separate branches
- Merge completed milestones sequentially
- Update progress.md after merges

---

## Summary

The milestone execution system provides:
- ✅ Structured, incremental development
- ✅ TDD enforcement
- ✅ Comprehensive documentation
- ✅ Quality checkpoints
- ✅ Clear progress tracking
- ✅ Manual testing guidance
- ✅ Enhanced 5-category test criteria (Functional, Integration, Edge Cases, Performance, UX)
- ✅ Testing guides for all milestone groups
- ✅ Detailed acceptance criteria for complex milestones
- ✅ Document reference matrix for easy navigation

Follow the process, trust the structure, and you'll build a high-quality, well-tested application.

**Ready to start?** Copy [`milestone-execution-prompt.md`](milestone-execution-prompt.md:1) and begin with Milestone 1!

**Key Documentation:**
- [`DOCUMENT-REFERENCE-MATRIX.md`](DOCUMENT-REFERENCE-MATRIX.md:1) - Start here to identify required documents
- [`DOCUMENT-INDEX.md`](DOCUMENT-INDEX.md:1) - Complete documentation overview
- [`ENHANCED-TEST-CRITERIA-TEMPLATE.md`](milestones/ENHANCED-TEST-CRITERIA-TEMPLATE.md) - Acceptance criteria framework