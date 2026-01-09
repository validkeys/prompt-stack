# Continue: Viewport Bubbles Migration Implementation

## Context
Refactored Task 5 viewport migration plan based on review findings. All 10 critical issues and 6 recommendations addressed. Plan is ready for implementation.

## Key Documents

**Viewport Migration Plan:**
- Task List: `docs/implementation-plans/m4-viewport-bubbles-migration/task-list.md`
- Reference: `docs/implementation-plans/m4-viewport-bubbles-migration/reference.md`
- README: `docs/implementation-plans/m4-viewport-bubbles-migration/README.md`
- Final Summary: `docs/implementation-plans/m4-viewport-bubbles-migration/FINAL-SUMMARY.md`

**Original Documents:**
- Original Plan: `docs/implementation-plans/m4-viewport-bubbles-migration/task5-viewport-bubbles-migration-ORIGINAL-BACKUP.md`
- Review: `docs/implementation-plans/m4-viewport-bubbles-migration/task5-viewport-bubbles-migration-review.md`

**M4 Milestone Plan (Updated):**
- Main Task List: `docs/plans/fresh-build/milestone-implementation-plans/M4/task-list.md`
  - Task 5 now points to viewport migration plan
  - References detailed 7-task breakdown
- M4 Reference: `docs/plans/fresh-build/milestone-implementation-plans/M4/reference.md`

## What Was Done

✅ **All 10 critical review findings addressed:**
1. Split into two-file structure (task-list.md + reference.md)
2. Added ALL 17 required document references
3. Comprehensive pre-implementation checklist (50+ items)
4. RFC 2119 keywords (MUST/SHOULD/MAY) for 141 acceptance criteria
5. >90% test coverage targets specified
6. Integration points with M5, M6, M34-M35, M37 documented
7. Complete code examples with imports and error handling
8. Navigation guide added
9. Key learnings from ui-domain.md and editor-domain.md applied
10. M4 task-list.md updated to point to viewport plan

✅ **All 6 recommended improvements included**

## Next Steps

### When Ready to Implement

1. **Review the plan**:
   - Read: `docs/implementation-plans/m4-viewport-bubbles-migration/README.md`
   - Review: `docs/implementation-plans/m4-viewport-bubbles-migration/task-list.md`
   - Reference: `docs/implementation-plans/m4-viewport-bubbles-migration/reference.md`

2. **Start implementation**:
   ```bash
   git checkout -b m4-viewport-bubbles-migration
   ```

3. **Execute tasks sequentially** (from task-list.md):
   - Task 1: Dependency Setup (5 min)
   - Task 2: Workspace Model Migration (30 min)
   - Task 3: View Rendering Integration (20 min)
   - Task 4: Cursor Visibility Integration (15 min)
   - Task 5: Test Migration (45 min)
   - Task 6: Archive Custom Viewport (5 min)
   - Task 7: Documentation Updates (15 min)

4. **After each task**: Create checkpoint, run `go test ./...`, verify build passes

## Quick Reference

**Testing Requirements:**
- Coverage target: >90% for viewport integration code
- Testing guide: `docs/plans/fresh-build/milestones/FOUNDATION-TESTING-GUIDE.md`

**Key Learnings Applied:**
- Middle-third scrolling from `docs/plans/fresh-build/learnings/ui-domain.md` (Category 2)
- Editor patterns from `docs/plans/fresh-build/learnings/editor-domain.md`

**Integration Points:**
- M5 (Auto-save): Scroll position in auto-save state
- M6 (Undo/Redo): Scroll position in undo stack
- M34-M35 (Vim Mode): j/k alignment with vim state machine
- M37 (Responsive Layout): Viewport resize handling

## Current Status

- ✅ Implementation plan complete and refactored
- ✅ M4 task-list updated
- ⏳ **Awaiting approval** to begin implementation

## Command to Start

```bash
# Review plan first
cat docs/implementation-plans/m4-viewport-bubbles-migration/README.md

# When approved, create branch and start:
git checkout -b m4-viewport-bubbles-migration
```

---

**Say "continue" to begin Task 1 (Dependency Setup)**
