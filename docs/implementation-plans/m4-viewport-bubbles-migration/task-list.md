# Task 5: Viewport with Scrolling - Bubbles Migration Task List

**Date:** 2026-01-09  
**Status:** PENDING APPROVAL  
**Author:** AI Agent  
**Reviewers:** Kyle Davis

---

## Overview

### Goal
Replace custom viewport implementation with production-ready `charmbracelet/bubbles/viewport` component to leverage battle-tested scrolling logic, automatic edge case handling, and improved performance for large documents.

### Deliverables
1. Bubbles viewport integration in workspace model
2. Middle-third scrolling strategy for cursor visibility
3. Automated bounds checking and edge case handling
4. Performance optimization for large documents (>10,000 lines)
5. Comprehensive test coverage (>90%)

### Dependencies
- **Prerequisite**: Task 4 completed (Basic Text Editor with character and line counting)
- **All tests passing** before migration
- **Go 1.21+** installed
- **Integration with M6**: Scroll position SHOULD be part of undo state (future consideration)
- **Integration with M34-M35**: j/k keys SHOULD align with vim state machine (future consideration)
- **Integration with M37**: Viewport MUST handle terminal resize events

### Document References
- [`milestones.md`](../../milestones.md) - Task 5 section
- [`requirements.md`](../../requirements.md) - Editor scrolling requirements
- [`project-structure.md`](../../project-structure.md) - UI and editor domain structure
- [`go-style-guide.md`](../../go-style-guide.md) - Go coding standards
- [`go-testing-guide.md`](../../go-testing-guide.md) - Testing patterns
- [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md) - Foundation milestone testing
- [`bubble-tea-testing-best-practices.md`](../../bubble-tea-testing-best-practices.md) - Bubble Tea testing
- [`bubble-tea-best-practices.md`](../../../../patterns/bubble-tea-best-practices.md) - Bubble Tea architecture
- [`DEPENDENCIES.md`](../../DEPENDENCIES.md) - Bubbles package version
- [`opencode-design-system.md`](../../opencode-design-system.md) - Viewport styling guidelines
- [`learnings/ui-domain.md`](../../learnings/ui-domain.md) - Cursor and viewport management patterns
- [`learnings/editor-domain.md`](../../learnings/editor-domain.md) - Editor integration patterns
- [`learnings/error-handling.md`](../../learnings/error-handling.md) - Error handling patterns
- [`learnings/architecture-patterns.md`](../../learnings/architecture-patterns.md) - Component extraction patterns

---

## Pre-Implementation Checklist

Before writing any code, verify:

### Package Structure
- [ ] All file paths match [`project-structure.md`](../../project-structure.md)
- [ ] Viewport integration in `ui/workspace/` package (not in `internal/editor/`)
- [ ] No packages in wrong locations

### Dependency Injection
- [ ] No global state variables planned
- [ ] Viewport initialized in workspace model constructor
- [ ] All dependencies passed through constructors
- [ ] Logger passed explicitly if needed

### Documentation
- [ ] Package comments planned for new viewport integration
- [ ] Exported function comments planned (with export comments)
- [ ] Error messages follow lowercase, no punctuation style
- [ ] Archive comment planned for old viewport implementation

### Testing
- [ ] Test files planned alongside implementation files
- [ ] Table-driven test structure planned per [`go-testing-guide.md`](../../go-testing-guide.md)
- [ ] Mock interfaces identified (viewport interface for testing)
- [ ] Test coverage target: >90% for viewport integration

### Style Compliance
- [ ] Constructor naming follows New() or NewType() pattern
- [ ] Error wrapping uses %w consistently per [`error-handling.md`](../../learnings/error-handling.md)
- [ ] Method receivers are consistent (all pointer or all value)
- [ ] No stuttering in names (e.g., viewport.ViewportScroll → viewport.Scroll)
- [ ] Import statements use correct bubbles package: `github.com/charmbracelet/bubbles/viewport`

### Constants
- [ ] Middle-third scroll strategy constants defined (e.g., SCROLL_MARGIN_RATIO = 1/3)
- [ ] No magic numbers in scroll calculations
- [ ] Viewport size constants if needed

### Design System Compliance
- [ ] Viewport styling follows [`opencode-design-system.md`](../../opencode-design-system.md)
- [ ] Color usage follows OpenCode palette guidelines
- [ ] Spacing follows 1-character unit system
- [ ] Visual feedback for scroll operations
- [ ] Keyboard shortcuts (j/k) follow OpenCode defaults

### Key Learnings Application
- [ ] Middle-third scrolling strategy from [`ui-domain.md`](../../learnings/ui-domain.md) applied
- [ ] Cursor visibility patterns from [`ui-domain.md`](../../learnings/ui-domain.md) applied
- [ ] Editor integration patterns from [`editor-domain.md`](../../learnings/editor-domain.md) applied
- [ ] Deviations from key learnings documented with justification

**If any item is unchecked, review and adjust plan before proceeding.**

---

## Tasks

### Task 1: Dependency Setup

**Dependencies:** None

**Files:**
- `go.mod` (per [`DEPENDENCIES.md`](../../DEPENDENCIES.md))
- `go.sum` (auto-generated)

**Integration Points:**
- None - foundational task

**Estimated Complexity:** Low (5 minutes)

**Description:**
Add `github.com/charmbracelet/bubbles` dependency to project. Verify version compatibility with existing bubbletea installation per [`DEPENDENCIES.md`](../../DEPENDENCIES.md).

**Acceptance Criteria (RFC 2119):**
1. **FR-1.1**: Bubbles package MUST be added to go.mod at version v0.18.0 or later
2. **FR-1.2**: `go mod tidy` MUST complete without errors
3. **FR-1.3**: `go build ./...` MUST succeed after dependency addition
4. **EC-1.1**: MUST handle case where bubbles version conflicts with bubbletea version
5. **UX-1.1**: Terminal output SHOULD clearly indicate dependency installation progress

**Testing Requirements:**
- **Coverage Target:** N/A (dependency only)
- **Critical Scenarios:**
  - Verify `go.mod` contains correct bubbles version
  - Verify `go build ./...` succeeds
- **Test Execution Order:** Build verification only
- **Edge Cases:**
  - Version conflicts with bubbletea
  - Network unavailable during `go get`

**Testing Guide Reference:** [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md) - Section on dependency management

---

### Task 2: Workspace Model Migration

**Dependencies:** Task 1

**Files:**
- `ui/workspace/model.go` (per [`project-structure.md`](../../project-structure.md))
- `ui/workspace/model_test.go`

**Integration Points:**
- Integrates with editor buffer in `internal/editor/buffer.go`
- Updates workspace initialization
- Affects cursor visibility logic

**Estimated Complexity:** Medium (30 minutes)

**Description:**
Migrate workspace model from custom viewport (`internal/editor/viewport.go`) to bubbles viewport (`github.com/charmbracelet/bubbles/viewport`). Update model structure, initialization, and viewport management methods. Apply patterns from [`ui-domain.md`](../../learnings/ui-domain.md) for cursor and viewport coordination.

**Acceptance Criteria (RFC 2119):**
1. **FR-2.1**: Model MUST replace `editor.Viewport` with `viewport.Model` from bubbles
2. **FR-2.2**: `New()` constructor MUST initialize bubbles viewport with `viewport.New(width, height)`
3. **FR-2.3**: Viewport MUST be updated on `tea.WindowSizeMsg` to handle terminal resize
4. **FR-2.4**: Buffer content MUST sync to viewport using `viewport.SetContent()`
5. **IR-2.1**: Integration with editor buffer MUST maintain existing API contracts
6. **IR-2.2**: Viewport MUST correctly integrate with status bar for line count display
7. **EC-2.1**: MUST handle zero-width or zero-height viewport without panic
8. **EC-2.2**: MUST handle nil buffer gracefully
9. **PR-2.1**: Viewport initialization MUST complete in <10ms
10. **UX-2.1**: Viewport SHOULD maintain current scroll position during resize when possible

**Testing Requirements:**
- **Coverage Target:** >85% for model changes
- **Critical Scenarios:**
  - Initialize viewport with various dimensions
  - Resize viewport from large to small
  - Sync buffer content with 0, 1, and 1000+ lines
- **Test Execution Order:** Unit → Integration
- **Edge Cases:**
  - Zero-dimension viewport
  - Nil buffer
  - Empty buffer
  - Very large buffer (>100,000 lines)

**Testing Guide Reference:** [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md) - Section 3: Model Testing

---

### Task 3: View Rendering Integration

**Dependencies:** Task 2

**Files:**
- `ui/workspace/view.go` (per [`project-structure.md`](../../project-structure.md))
- `ui/workspace/view_test.go`

**Integration Points:**
- Integrates with status bar rendering
- Affects cursor position display
- Uses bubbles viewport View() method

**Estimated Complexity:** Low (20 minutes)

**Description:**
Update view rendering to use bubbles viewport's built-in `View()` method. Remove manual visible line calculation and slicing logic. Maintain cursor position display and status bar integration. Follow UI patterns from [`ui-domain.md`](../../learnings/ui-domain.md).

**Acceptance Criteria (RFC 2119):**
1. **FR-3.1**: View MUST use `viewport.View()` for content rendering
2. **FR-3.2**: View MUST call `viewport.SetContent()` with buffer content before rendering
3. **FR-3.3**: Cursor position MUST remain visible in rendered output
4. **FR-3.4**: Status bar MUST display correct line counts after migration
5. **IR-3.1**: View rendering MUST integrate with existing status bar component
6. **EC-3.1**: MUST handle empty viewport content without panic
7. **EC-3.2**: MUST render correctly with unicode characters (emoji, CJK)
8. **PR-3.1**: Rendering 1000 lines MUST complete in <50ms
9. **PR-3.2**: Frame rate MUST maintain >=30fps during scrolling
10. **UX-3.1**: Content MUST render without flickering or visual artifacts

**Testing Requirements:**
- **Coverage Target:** >90% for view rendering logic
- **Critical Scenarios:**
  - Render empty content
  - Render single line
  - Render 1000+ lines with scrolling
  - Render with unicode characters
- **Test Execution Order:** Unit → Visual validation
- **Edge Cases:**
  - Empty content
  - Very long single line (>1000 chars)
  - Content with only newlines
  - Unicode edge cases (zero-width joiners, RTL text)

**Testing Guide Reference:** [`bubble-tea-testing-best-practices.md`](../../bubble-tea-testing-best-practices.md) - View testing patterns

---

### Task 4: Cursor Visibility Integration

**Dependencies:** Task 3

**Files:**
- `ui/workspace/model.go` (update)
- `ui/workspace/viewport_cursor.go` (new - cursor visibility logic)
- `ui/workspace/viewport_cursor_test.go` (new)

**Integration Points:**
- Integrates with buffer cursor position
- Affects viewport scroll position (YOffset)
- Critical for M6 (Undo/Redo) - scroll position SHOULD be part of undo state
- Critical for M34-M35 (Vim Mode) - j/k keys MUST work consistently

**Estimated Complexity:** Medium (15 minutes)

**Description:**
Implement custom cursor visibility logic using middle-third scrolling strategy. Since bubbles viewport doesn't have built-in `EnsureVisible()`, implement custom logic per patterns from [`ui-domain.md`](../../learnings/ui-domain.md) Category 2. This ensures cursor remains visible and scrolls smoothly.

**Acceptance Criteria (RFC 2119):**
1. **FR-4.1**: Cursor MUST remain visible during all cursor movement operations
2. **FR-4.2**: Scrolling MUST use middle-third strategy (cursor in middle 1/3 of viewport)
3. **FR-4.3**: Viewport MUST scroll when cursor exits middle third region
4. **FR-4.4**: YOffset MUST NEVER be negative (bounds checking)
5. **FR-4.5**: YOffset MUST NEVER exceed (total_lines - viewport_height)
6. **IR-4.1**: Cursor visibility MUST work with buffer navigation (arrow keys, home, end)
7. **IR-4.2**: Scroll position SHOULD integrate with future undo/redo system (M6)
8. **IR-4.3**: Scrolling MUST work consistently with future vim mode (M34-M35) j/k keys
9. **EC-4.1**: MUST handle cursor at document start (line 0)
10. **EC-4.2**: MUST handle cursor at document end (last line)
11. **EC-4.3**: MUST handle viewport taller than document
12. **EC-4.4**: MUST handle viewport height of 1 line
13. **PR-4.1**: Scroll adjustment MUST complete in <1ms per operation
14. **PR-4.2**: Rapid cursor movement (10 operations/sec) MUST maintain smooth scrolling
15. **UX-4.1**: Scrolling SHOULD feel natural with no jarring jumps
16. **UX-4.2**: Scroll position SHOULD maintain context (2-3 lines visible above/below cursor)

**Testing Requirements:**
- **Coverage Target:** >95% for cursor visibility logic (critical functionality)
- **Critical Scenarios:**
  - Cursor moves down from top of viewport → scrolls after exiting middle third
  - Cursor moves up from bottom of viewport → scrolls after exiting middle third
  - Cursor at line 0 → viewport at offset 0
  - Cursor at last line → viewport at max offset
  - Viewport taller than document → no scrolling
  - Rapid cursor movement (100 operations) → no negative offsets
- **Test Execution Order:** Unit → Integration → Performance
- **Edge Cases:**
  - Zero-height viewport
  - One-line viewport
  - Cursor position beyond document length (error case)
  - Negative cursor position (error case)
  - Document with 0 lines
  - Document with 1 line

**Testing Guide Reference:** [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md) - Section 4: Viewport Testing

**Key Learnings Applied:**
- Middle-third scrolling from [`ui-domain.md`](../../learnings/ui-domain.md) Category 2
- Cursor/viewport coordination patterns from [`ui-domain.md`](../../learnings/ui-domain.md)

---

### Task 5: Test Migration and Enhancement

**Dependencies:** Task 4

**Files to Delete:**
- `internal/editor/viewport_test.go`
- `internal/editor/viewport_debug_test.go`
- `internal/editor/viewport_debug2_test.go`
- `internal/editor/viewport_debug3_test.go`

**Files to Create:**
- `ui/workspace/viewport_integration_test.go` (per [`project-structure.md`](../../project-structure.md))

**Integration Points:**
- Tests integration with buffer
- Tests integration with cursor movement
- Tests integration with status bar

**Estimated Complexity:** High (45 minutes)

**Description:**
Remove legacy custom viewport tests and create comprehensive integration tests for bubbles viewport integration. Follow table-driven test patterns from [`go-testing-guide.md`](../../go-testing-guide.md) and [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md). Focus on cursor visibility, middle-third scrolling, and edge cases.

**Acceptance Criteria (RFC 2119):**
1. **FR-5.1**: All legacy viewport test files MUST be deleted
2. **FR-5.2**: New integration tests MUST use table-driven test pattern per [`go-testing-guide.md`](../../go-testing-guide.md)
3. **FR-5.3**: Tests MUST cover all acceptance criteria from Tasks 1-4
4. **FR-5.4**: Tests MUST verify middle-third scrolling behavior
5. **FR-5.5**: Tests MUST verify bounds checking (no negative YOffset)
6. **IR-5.1**: Integration tests MUST verify buffer-viewport synchronization
7. **IR-5.2**: Integration tests MUST verify cursor movement triggers correct scrolling
8. **IR-5.3**: Integration tests SHOULD verify scroll position with undo/redo (if M6 available)
9. **EC-5.1**: Tests MUST cover empty document case
10. **EC-5.2**: Tests MUST cover single line document case
11. **EC-5.3**: Tests MUST cover large document (>10,000 lines) case
12. **EC-5.4**: Tests MUST cover viewport smaller than document
13. **EC-5.5**: Tests MUST cover viewport larger than document
14. **PR-5.1**: Code coverage MUST exceed 90% for viewport integration code
15. **PR-5.2**: All tests MUST pass with race detector (`go test -race`)
16. **UX-5.1**: Test names SHOULD clearly describe scenario being tested

**Testing Requirements:**
- **Coverage Target:** >90% for all viewport integration code
- **Critical Scenarios:**
  - Cursor visibility maintained during navigation (10+ test cases)
  - Middle-third scrolling triggered correctly (5+ test cases)
  - Bounds checking prevents negative offsets (3+ test cases)
  - Large document performance (1 benchmark test)
  - Buffer synchronization (3+ test cases)
- **Test Execution Order:** Unit → Integration → Performance/Benchmark
- **Edge Cases:**
  - All edge cases from Tasks 2-4 MUST have dedicated test cases
  - Unicode handling test cases
  - Terminal resize during scrolling

**Testing Guide Reference:** [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md) - Complete guide for Foundation milestone testing

**Test Example Structure (per [`go-testing-guide.md`](../../go-testing-guide.md)):**
```go
func TestViewportCursorVisibility(t *testing.T) {
    tests := []struct {
        name           string
        viewportHeight int
        cursorY        int
        currentOffset  int
        expectedOffset int
    }{
        // Test cases here
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

---

### Task 6: Archive Custom Viewport

**Dependencies:** Task 5 (all tests passing)

**Files:**
- `internal/editor/viewport.go` → `docs/archived/viewport-custom-v1.go`
- Create archive directory if needed

**Integration Points:**
- None - cleanup task

**Estimated Complexity:** Low (5 minutes)

**Description:**
Archive the custom viewport implementation for future reference. Add context comment explaining why it was replaced and when. Follow archival best practices.

**Acceptance Criteria (RFC 2119):**
1. **FR-6.1**: `docs/archived/` directory MUST exist
2. **FR-6.2**: Custom viewport code MUST be copied to `docs/archived/viewport-custom-v1.go`
3. **FR-6.3**: Archived file MUST include header comment with:
   - Date archived (2026-01-09)
   - Reason for replacement (migrated to bubbles viewport)
   - Related task (Task 5)
   - Original file path (internal/editor/viewport.go)
4. **FR-6.4**: Original `internal/editor/viewport.go` MUST be deleted
5. **FR-6.5**: `go build ./...` MUST succeed after deletion
6. **EC-6.1**: MUST verify no remaining imports of custom viewport package
7. **UX-6.1**: Archive comment SHOULD explain what bubbles viewport provides that custom didn't

**Testing Requirements:**
- **Coverage Target:** N/A (cleanup only)
- **Critical Scenarios:**
  - Verify no build errors after deletion
  - Verify no remaining imports of custom viewport
  - Verify archived file is accessible
- **Test Execution Order:** Build verification → Import verification
- **Edge Cases:** None

**Testing Guide Reference:** N/A

---

### Task 7: Documentation Updates

**Dependencies:** Task 6

**Files:**
- `docs/patterns/bubble-tea-best-practices.md` (already updated)
- `docs/architecture/viewport-design.md` (new)

**Integration Points:**
- References bubbles viewport documentation
- Documents middle-third scrolling strategy
- Provides examples for future developers

**Estimated Complexity:** Low (15 minutes)

**Description:**
Create architectural documentation for viewport design. Document the middle-third scrolling strategy, bubbles viewport integration patterns, and best practices for future development. Reference [`opencode-design-system.md`](../../opencode-design-system.md) for styling guidelines.

**Acceptance Criteria (RFC 2119):**
1. **FR-7.1**: `docs/architecture/viewport-design.md` MUST be created
2. **FR-7.2**: Documentation MUST explain middle-third scrolling strategy
3. **FR-7.3**: Documentation MUST provide API reference for viewport methods:
   - `viewport.SetContent(content string)`
   - `viewport.YOffset` (read/write)
   - `viewport.View()` → string
   - `viewport.SetYOffset(offset int)`
   - `viewport.LineUp(n int)` / `viewport.LineDown(n int)`
4. **FR-7.4**: Documentation MUST include code examples with correct imports
5. **FR-7.5**: Documentation MUST reference [`ui-domain.md`](../../learnings/ui-domain.md) for patterns
6. **IR-7.1**: Documentation SHOULD explain integration points with M6 (Undo/Redo)
7. **IR-7.2**: Documentation SHOULD explain integration points with M34-M35 (Vim Mode)
8. **IR-7.3**: Documentation SHOULD explain integration points with M37 (Responsive Layout)
9. **UX-7.1**: Documentation SHOULD include diagrams or ASCII art showing middle-third regions
10. **UX-7.2**: Documentation SHOULD provide troubleshooting section

**Testing Requirements:**
- **Coverage Target:** N/A (documentation only)
- **Critical Scenarios:**
  - Documentation is readable and well-formatted
  - Code examples are valid Go code
  - Links to other documents are correct
- **Test Execution Order:** Manual review
- **Edge Cases:** None

**Testing Guide Reference:** N/A

---

## Integration Points

### Integration with M6: Undo/Redo System
**Status**: Future consideration (M6 not yet implemented)

**Integration Point**: Scroll position should be part of undo state

**Details:**
- When implementing M6, consider adding `viewport.YOffset` to undo state
- Undoing a text edit SHOULD restore scroll position to where user was
- Redoing SHOULD restore scroll position forward

**Testing Strategy**:
- Add integration test in M6: Undo edit → verify viewport position restored
- Test: Redo edit → verify viewport position restored forward

**How to Test:**
1. Make edit at bottom of document (scroll down)
2. Undo edit
3. Verify viewport scrolls back to previous position

### Integration with M34-M35: Vim Mode
**Status**: Future consideration (M34-M35 not yet implemented)

**Integration Point**: j/k keys MUST work consistently in normal mode

**Details:**
- j/k keys in vim normal mode should trigger same cursor movement as arrow keys
- Cursor visibility logic MUST work with vim navigation commands
- Middle-third scrolling MUST apply to vim navigation

**Testing Strategy**:
- Add integration test in M34-M35: Press j/k 50 times → verify smooth scrolling
- Test: Vim motions (10j, 5k) → verify correct cursor/viewport coordination

**How to Test:**
1. Enter vim normal mode
2. Press j 100 times rapidly
3. Verify cursor stays visible
4. Verify middle-third scrolling applies

### Integration with M37: Responsive Layout
**Status**: Future consideration (M37 not yet implemented)

**Integration Point**: Viewport MUST resize correctly on terminal resize

**Details:**
- Terminal resize SHOULD preserve scroll position when possible
- Viewport MUST handle resize from large to small gracefully
- Status bar MUST update after resize

**Testing Strategy**:
- Add integration test in M37: Resize terminal → verify viewport adjusts
- Test: Resize to minimum size (80x24) → verify no crashes
- Test: Resize to very small (20x10) → verify graceful degradation

**How to Test:**
1. Open document with 100 lines
2. Scroll to line 50
3. Resize terminal from 120x40 to 80x24
4. Verify viewport adjusts and cursor stays visible
5. Verify scroll position preserved proportionally

---

## File Impact Summary

### Files to Modify
1. `go.mod` - Add bubbles dependency
2. `ui/workspace/model.go` - Migrate to bubbles viewport
3. `ui/workspace/view.go` - Update rendering logic
4. `ui/workspace/viewport_cursor.go` - Add cursor visibility logic (new file)

### Files to Delete
1. `internal/editor/viewport.go` - Custom viewport implementation
2. `internal/editor/viewport_test.go` - Legacy tests
3. `internal/editor/viewport_debug_test.go` - Debug tests
4. `internal/editor/viewport_debug2_test.go` - Debug tests
5. `internal/editor/viewport_debug3_test.go` - Debug tests

### Files to Create
1. `docs/archived/viewport-custom-v1.go` - Archived custom implementation
2. `ui/workspace/viewport_cursor.go` - Cursor visibility logic
3. `ui/workspace/viewport_cursor_test.go` - Cursor visibility tests
4. `ui/workspace/viewport_integration_test.go` - Integration tests
5. `docs/architecture/viewport-design.md` - Architecture documentation

---

## Timeline Estimate

| Task | Description | Time | Cumulative |
|------|-------------|------|------------|
| 1 | Dependency Setup | 5 min | 5 min |
| 2 | Workspace Migration | 30 min | 35 min |
| 3 | View Integration | 20 min | 55 min |
| 4 | Cursor Visibility | 15 min | 70 min |
| 5 | Test Migration | 45 min | 115 min |
| 6 | Archive Custom Code | 5 min | 120 min |
| 7 | Documentation | 15 min | 135 min |
| **Buffer Time** | | 15 min | **150 min** |

**Total Estimated Time:** 2.5 hours

---

## Risk Assessment

### Low Risk ✅
- Bubbles viewport is production-proven (used in Glow, Charm's official editor)
- Well-documented API with active maintenance
- No breaking changes to external API (internal refactoring only)
- Can rollback by restoring archived viewport (<10 minutes)

### Medium Risk ⚠️
- Need to implement custom middle-third scrolling logic (not built into bubbles)
- May need adjustments for cursor visibility edge cases
- **Mitigation:** Thorough testing with various document sizes per Task 5

### Identified Risks
- **Risk 1:** Bubbles viewport behavior differs from custom implementation
  - **Impact:** Medium - May require adjustment period
  - **Mitigation:** Comprehensive test suite, manual testing
- **Risk 2:** Performance regression with very large documents (>100,000 lines)
  - **Impact:** Low - Bubbles is optimized for large content
  - **Mitigation:** Performance benchmarks in test suite
- **Risk 3:** Integration breaks with future milestones (M6, M34-M35, M37)
  - **Impact:** Low - Integration points documented
  - **Mitigation:** Integration points section above, documented considerations

---

## Rollback Plan

If issues arise after implementation:

1. Restore `internal/editor/viewport.go` from `docs/archived/viewport-custom-v1.go`
2. Revert changes to `ui/workspace/model.go` and `ui/workspace/view.go`
3. Revert `go.mod` changes (remove bubbles dependency)
4. Run `go mod tidy`
5. Verify all tests pass
6. Document reason for rollback in checkpoint

**Rollback Time Estimate:** <10 minutes

---

## Post-Implementation Verification

After completing all tasks:

1. [ ] Run full test suite: `go test ./... -v -race -cover`
2. [ ] Verify code coverage >90% for viewport code
3. [ ] Manual testing with large files (>10,000 lines)
4. [ ] Performance benchmarking (typing latency <10ms, scroll <1ms)
5. [ ] Terminal resize testing (various sizes)
6. [ ] Visual inspection for flickering or artifacts
7. [ ] Code review focusing on:
   - Middle-third scrolling implementation
   - Bounds checking correctness
   - Error handling patterns
   - Test coverage
8. [ ] Update task tracking document
9. [ ] Mark Task 5 as complete in milestones/progress.md

---

## Success Metrics

### Functional Metrics
- ✅ All 7 tasks completed
- ✅ All tests passing (unit + integration)
- ✅ No build errors or warnings
- ✅ Manual testing scenarios pass

### Quality Metrics
- ✅ Code coverage >90% for viewport integration
- ✅ No race conditions detected
- ✅ All RFC 2119 acceptance criteria met (MUST/SHOULD/MAY)

### Performance Metrics
- ✅ Scroll operations <1ms (Task 4 PR-4.1)
- ✅ Rendering 1000 lines <50ms (Task 3 PR-3.1)
- ✅ Frame rate >=30fps during scrolling (Task 3 PR-3.2)

### Documentation Metrics
- ✅ Architecture docs created and reviewed
- ✅ Old implementation archived with context
- ✅ Integration points documented for future milestones

---

## Approval

**Status:** ⏳ AWAITING APPROVAL

**Approver:** Kyle Davis  
**Date Approved:** _________________  
**Signature:** _________________

---

## Next Steps After Approval

1. Execute Task 1 (Dependency Setup)
2. Create implementation branch: `task5-bubbles-viewport`
3. Follow tasks sequentially (1→2→3→4→5→6→7)
4. Stop after each task for human verification (create checkpoint)
5. Request review after Task 5 (Test Migration)
6. Complete documentation and archive in Tasks 6-7
