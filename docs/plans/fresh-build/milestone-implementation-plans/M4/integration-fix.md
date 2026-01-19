# M4 Integration Fix

**Issue**: Application launched with basic TUI shell (M2) instead of text editor workspace (M4)

**Root Cause**: `cmd/promptstack/main.go` was using `app.New()` from M2 instead of `workspace.New()` from M4

**Fix Applied**:
1. Changed import from `ui/app` to `ui/workspace`
2. Changed `appModel := app.New()` to `workspaceModel := workspace.New()`
3. Updated Bubble Tea program to use `workspaceModel` instead of `appModel`

**Verification**:
- App builds successfully: `go build -o promptstack ./cmd/promptstack`
- M4 integration tests pass: `go test ./test/integration/m4_editor_test.go`
- All 21 test scenarios pass

**Next Step**: Manual testing should now work - text editor workspace will launch with full text editing functionality

**Files Modified**:
- `cmd/promptstack/main.go` - Lines 11, 54-55, 61

**Date**: 2026-01-09
