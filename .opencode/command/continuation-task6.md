# M4 Task 6: Apply Theme Styles to Workspace

## Current Status
- **Milestone**: 4 (Basic Text Editor)
- **Progress**: 5/7 tasks completed (Tasks 1-5 done)
- **Current Task**: Task 6 - Apply Theme Styles to Workspace
- **Last Commit**: 796daef "Update M4 task list: Mark Task 5 as completed"

## Task Overview
Apply OpenCode design system styles to workspace using centralized theme package. Implements consistent color palette, spacing, and visual hierarchy.

## Key Documents

**Implementation Plan:**
- `docs/plans/fresh-build/milestone-implementation-plans/M4/task-list.md` (Task 6 section, lines 493-543)

**Design System:**
- `ui/theme/theme.go` - Theme package with color constants and style functions
- `docs/plans/fresh-build/opencode-design-system.md` - Color System, Typography & Layout

**Reference Guides:**
- `docs/plans/fresh-build/milestones/FOUNDATION-TESTING-GUIDE.md` - Foundation testing guide

## What to Do

1. **Review Task 6 requirements** in M4 task-list.md (lines 493-543)
2. **Read design system** in opencode-design-system.md (color palette, spacing)
3. **Check existing theme** in ui/theme/theme.go (available styles)
4. **Apply styles** to ui/workspace/view.go:
   - Use theme.BackgroundPrimary for editor background
   - Use theme.ForegroundPrimary for editor text
   - Use theme.CursorStyle() for cursor highlighting
   - Use theme.StatusStyle() for status bar
5. **Add tests** for theme integration
6. **Verify** test coverage >75%
7. **Run tests** with `go test ./ui/workspace -v`

## Acceptance Criteria (RFC 2119)
- FR1-FR7: Must use theme package (no hardcoded colors)
- EC1-EC3: Must handle 256-color and truecolor terminals
- PR1-PR2: Must complete in <16ms, <1ms style application
- UX1-UX4: High contrast cursor, readable text, distinct status bar

## Integration Points
- Uses `ui/theme/theme.go` for color constants
- Integrates with workspace View() rendering

## Estimated Complexity
Low (~1-2 hours)

## Next After Task 6
Task 7: Integration Testing and Manual Verification (final M4 task)
