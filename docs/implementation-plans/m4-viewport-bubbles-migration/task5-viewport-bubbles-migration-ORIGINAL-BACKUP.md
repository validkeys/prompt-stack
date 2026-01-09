# Task 5: Viewport with Scrolling - Bubbles Migration Plan

**Date:** 2026-01-09  
**Status:** PENDING APPROVAL  
**Author:** AI Agent  
**Reviewers:** Kyle Davis

---

## Executive Summary

Replace custom viewport implementation with production-ready `charmbracelet/bubbles/viewport` component to leverage battle-tested scrolling logic, automatic edge case handling, and improved performance for large documents.

---

## Problem Statement

### Current Issues
1. Custom viewport has bounds checking bug (negative topLine values)
2. Missing edge case handling for small viewports
3. Reinventing already-solved problems
4. Increased maintenance burden

### Research Findings
- Bubbles viewport is the production standard for Bubble Tea applications
- Used by thousands of applications including Glow, Charm's official editor
- Includes automatic bounds checking and edge case handling
- Supports middle-third scrolling strategy
- Better performance with large documents (>10,000 lines)

---

## Implementation Plan

### Phase 1: Dependency Setup
**Files to modify:**
- `go.mod`

**Actions:**
1. Add `github.com/charmbracelet/bubbles` dependency
2. Run `go mod tidy`
3. Verify dependency installation

**Estimated time:** 5 minutes

---

### Phase 2: Workspace Model Migration
**Files to modify:**
- `ui/workspace/model.go`

**Current structure:**
```go
type Model struct {
    buffer   editor.Buffer
    viewport editor.Viewport  // Custom implementation
    // ...
}
```

**New structure:**
```go
import "github.com/charmbracelet/bubbles/viewport"

type Model struct {
    buffer   editor.Buffer
    viewport viewport.Model  // Bubbles implementation
    // ...
}
```

**Method updates required:**
- `New()` - Initialize bubbles viewport with `viewport.New(width, height)`
- `adjustViewport()` - Use `viewport.SetContent()` and `viewport.YOffset`
- `syncViewportLines()` - Use `viewport.SetContent()` to sync with buffer
- `View()` - Use `viewport.View()` for rendering
- `Update()` - Forward viewport messages to bubbles viewport

**Key API mappings:**
| Custom Viewport | Bubbles Viewport |
|----------------|------------------|
| `viewport.TopLine()` | `viewport.YOffset` |
| `viewport.ScrollTo(line)` | `viewport.SetYOffset(line)` |
| `viewport.EnsureVisible(line)` | Custom logic with `viewport.LineUp()/LineDown()` |
| `viewport.SetTotalLines(n)` | Handled by `viewport.SetContent()` |
| `viewport.VisibleLines()` | `viewport.YOffset` to `viewport.YOffset + viewport.Height` |

**Estimated time:** 30 minutes

---

### Phase 3: View Rendering Integration
**Files to modify:**
- `ui/workspace/view.go`

**Changes required:**
1. Remove manual visible line calculation
2. Use `viewport.SetContent(bufferContent)` 
3. Use `viewport.View()` for rendering
4. Maintain cursor position display
5. Keep status bar integration

**Current approach:**
```go
func (m Model) getVisibleLines() []string {
    start, end := m.viewport.VisibleLines()
    // Manual slicing...
}
```

**New approach:**
```go
func (m Model) View() string {
    // Set viewport content from buffer
    m.viewport.SetContent(m.buffer.String())
    
    // Render viewport (handles scrolling automatically)
    return m.viewport.View()
}
```

**Estimated time:** 20 minutes

---

### Phase 4: Cursor Visibility Integration
**Files to modify:**
- `ui/workspace/model.go`

**Challenge:** Bubbles viewport doesn't have built-in `EnsureVisible()` method

**Solution:** Implement custom cursor tracking:
```go
func (m Model) ensureCursorVisible() Model {
    _, cursorY := m.buffer.CursorPosition()
    
    // Get current viewport bounds
    viewportTop := m.viewport.YOffset
    viewportBottom := viewportTop + m.viewport.Height
    
    // Calculate middle third
    third := m.viewport.Height / 3
    middleTop := viewportTop + third
    middleBottom := viewportBottom - third
    
    // Scroll if cursor outside middle third
    if cursorY < middleTop && cursorY >= 0 {
        newOffset := max(0, cursorY - third)
        m.viewport.SetYOffset(newOffset)
    } else if cursorY >= middleBottom {
        newOffset := cursorY - m.viewport.Height + third
        m.viewport.SetYOffset(newOffset)
    }
    
    return m
}
```

**Estimated time:** 15 minutes

---

### Phase 5: Test Migration
**Files to modify:**
- Remove `internal/editor/viewport_test.go`
- Remove `internal/editor/viewport_debug_test.go`
- Remove `internal/editor/viewport_debug2_test.go`
- Remove `internal/editor/viewport_debug3_test.go`
- Update `ui/workspace/model_test.go` (if exists)

**New test strategy:**
1. Test cursor visibility logic (middle-third scrolling)
2. Test buffer-viewport synchronization
3. Test edge cases (empty buffer, single line, etc.)
4. Integration tests for cursor movement with scrolling

**Test files to create:**
- `ui/workspace/viewport_integration_test.go`

**Estimated time:** 45 minutes

---

### Phase 6: Archive Custom Viewport
**Files to archive:**
- `internal/editor/viewport.go` → `docs/archived/viewport-custom-v1.go`

**Rationale:** Keep for reference in case we need to understand original logic

**Actions:**
1. Create `docs/archived/` directory
2. Copy viewport.go with timestamp and context comment
3. Delete `internal/editor/viewport.go`

**Estimated time:** 5 minutes

---

### Phase 7: Documentation Updates
**Files to update:**
- `docs/patterns/bubble-tea-best-practices.md` (already updated)
- Create `docs/architecture/viewport-design.md`

**New documentation:**
```markdown
# Viewport Design

## Overview
Uses `charmbracelet/bubbles/viewport` for production-ready scrolling.

## Middle-Third Scrolling Strategy
- Cursor stays in middle third of viewport
- Scrolls only when cursor exits middle third
- Provides optimal context without jarring jumps

## API Reference
- `viewport.SetContent()` - Update viewport content
- `viewport.YOffset` - Current scroll position
- `viewport.View()` - Render visible content
```

**Estimated time:** 15 minutes

---

## File Impact Summary

### Files to Modify
1. `go.mod` - Add bubbles dependency
2. `ui/workspace/model.go` - Migrate to bubbles viewport
3. `ui/workspace/view.go` - Update rendering
4. `ui/workspace/model_test.go` - Update tests

### Files to Delete
1. `internal/editor/viewport.go`
2. `internal/editor/viewport_test.go`
3. `internal/editor/viewport_debug_test.go`
4. `internal/editor/viewport_debug2_test.go`
5. `internal/editor/viewport_debug3_test.go`

### Files to Create
1. `docs/archived/viewport-custom-v1.go` - Archive old implementation
2. `ui/workspace/viewport_integration_test.go` - New tests
3. `docs/architecture/viewport-design.md` - Documentation

---

## Risk Assessment

### Low Risk
- ✅ Bubbles viewport is production-proven
- ✅ Well-documented API
- ✅ Active maintenance and support
- ✅ No breaking changes to external API

### Medium Risk
- ⚠️ Need to implement custom middle-third scrolling logic
- ⚠️ May need adjustments for cursor visibility
- **Mitigation:** Thorough testing with various document sizes

### Minimal Risk
- No user-facing API changes
- Internal refactoring only
- Can rollback by restoring archived viewport

---

## Testing Strategy

### Unit Tests
- [ ] Cursor visibility in middle third
- [ ] Cursor scrolling at top boundary
- [ ] Cursor scrolling at bottom boundary
- [ ] Empty document edge case
- [ ] Single line document
- [ ] Large document (>1000 lines)

### Integration Tests
- [ ] Buffer insert triggers viewport update
- [ ] Buffer delete triggers viewport update
- [ ] Cursor movement maintains visibility
- [ ] Status bar shows correct line counts

### Manual Testing
- [ ] Navigate large file (>10,000 lines)
- [ ] Rapid cursor movement (j/k keys)
- [ ] Jump to beginning/end of document
- [ ] Resize terminal window
- [ ] Insert/delete text while scrolling

---

## Success Criteria

### Functional Requirements
1. ✅ Cursor always visible when moving
2. ✅ Middle-third scrolling behavior maintained
3. ✅ No negative scroll positions
4. ✅ Smooth scrolling without jumps
5. ✅ Works with documents >10,000 lines

### Technical Requirements
1. ✅ All tests passing
2. ✅ No custom viewport code remaining
3. ✅ Code coverage >90% for new viewport logic
4. ✅ Performance: <1ms for scroll operations

### Documentation Requirements
1. ✅ Architecture docs updated
2. ✅ Old implementation archived
3. ✅ Best practices documented

---

## Timeline Estimate

| Phase | Time | Cumulative |
|-------|------|------------|
| 1. Dependency Setup | 5 min | 5 min |
| 2. Workspace Migration | 30 min | 35 min |
| 3. View Integration | 20 min | 55 min |
| 4. Cursor Visibility | 15 min | 70 min |
| 5. Test Migration | 45 min | 115 min |
| 6. Archive Custom Code | 5 min | 120 min |
| 7. Documentation | 15 min | 135 min |
| **Buffer Time** | 15 min | **150 min** |

**Total Estimated Time:** 2.5 hours

---

## Rollback Plan

If issues arise:
1. Restore `internal/editor/viewport.go` from `docs/archived/`
2. Revert changes to `ui/workspace/model.go`
3. Revert `go.mod` changes
4. Run `go mod tidy`
5. Verify all tests pass

**Rollback time:** <10 minutes

---

## Post-Implementation Tasks

1. [ ] Run full test suite
2. [ ] Manual testing with large files
3. [ ] Performance benchmarking
4. [ ] Code review
5. [ ] Update task tracking document
6. [ ] Mark Task 5 as complete

---

## Dependencies

### Prerequisites
- Task 4 completed (Character and Line Counting)
- All tests passing before migration
- Go 1.21+ installed

### External Dependencies
- `github.com/charmbracelet/bubbles` v0.18.0+
- `github.com/charmbracelet/bubbletea` (already in go.mod)

---

## Open Questions

1. **Q:** Should we keep cursor centering option for future?
   **A:** No, middle-third is optimal based on research

2. **Q:** Should we add mouse wheel support immediately?
   **A:** Yes, bubbles viewport includes it automatically

3. **Q:** What about horizontal scrolling?
   **A:** Defer to future task (not in current scope)

---

## References

- [Bubbles Viewport Documentation](https://github.com/charmbracelet/bubbles/tree/master/viewport)
- [Bubble Tea Best Practices](../patterns/bubble-tea-best-practices.md)
- [Task 5 Requirements](../tasks/task5-viewport-scrolling.md)

---

## Approval

**Status:** ⏳ AWAITING APPROVAL

**Approver:** Kyle Davis  
**Date Approved:** _________________  
**Signature:** _________________

---

## Notes

- Consider this a **refactoring task**, not a feature addition
- No user-facing changes expected
- Focus on code quality and maintainability
- Leverage battle-tested components over custom code

---

**Next Steps After Approval:**
1. Execute Phase 1 (Dependency Setup)
2. Create implementation branch: `task5-bubbles-viewport`
3. Follow phases sequentially
4. Request review after Phase 5 (Test Migration)
5. Complete documentation and archive in Phase 6-7
