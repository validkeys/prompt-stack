# Viewport Bubbles Migration Plan - Refactoring Complete

**Date**: 2026-01-09  
**Status**: Ready for Review  
**Location**: `docs/implementation-plans/m4-viewport-bubbles-migration/`

---

## Correction Made

**Issue**: The refactored plan was initially placed in wrong location.

**Original Understanding**: Task 5 was treated as a separate milestone (M5).

**Correction**: This is actually **Task 5 within Milestone 4 (Basic Text Editor)**, not a separate milestone. Looking at `milestones.md`:
- **Milestone 4**: Basic Text Editor (current milestone)
- **Milestone 5**: Auto-save (completely different milestone)

---

## Plan Location

The refactored implementation plan is now located at:

```
docs/implementation-plans/m4-viewport-bubbles-migration/
├── README.md                              # This file
├── task-list.md                           # 27KB - Task breakdown
├── reference.md                             # 32KB - Reference & examples
├── task5-viewport-bubbles-migration-ORIGINAL-BACKUP.md  # Original plan
└── task5-viewport-bubbles-migration-review.md          # Review findings
```

**Note**: This plan is **separate from the M4 milestone implementation plan**. It is a focused, detailed plan specifically for viewport migration that can be integrated into M4 when implementing Task 5.

---

## What Was Accomplished

### All Required Changes from Review ✅

#### 1. Structural Changes (CRITICAL)
- ✅ **Split into two-file structure**: `task-list.md` and `reference.md`
- ✅ **Follows fresh-build format**: Matches `milestone-execution-prompt.md` requirements exactly
- ✅ **Properly scoped**: Task 5 within M4, not a separate milestone

#### 2. Document References (CRITICAL)
Added ALL required document references from DOCUMENT-REFERENCE-MATRIX.md:

**Core Planning**:
- milestone-execution-prompt.md ✅
- milestones.md (M4 section) ✅
- requirements.md ✅
- project-structure.md ✅

**Implementation Guides**:
- go-style-guide.md ✅
- go-testing-guide.md ✅
- bubble-tea-testing-best-practices.md ✅
- bubble-tea-best-practices.md ✅
- ENHANCED-TEST-CRITERIA-TEMPLATE.md ✅
- opencode-design-system.md ✅
- DEPENDENCIES.md ✅

**Key Learnings**:
- learnings/ui-domain.md ✅
- learnings/editor-domain.md ✅
- learnings/error-handling.md ✅
- learnings/architecture-patterns.md ✅

**Testing Guides**:
- FOUNDATION-TESTING-GUIDE.md ✅

#### 3. Pre-Implementation Checklist (CRITICAL)
Comprehensive checklist covering:
- Package structure verification ✅
- Dependency injection patterns ✅
- Documentation requirements ✅
- Testing requirements ✅
- Code style requirements ✅
- Constants needed ✅
- Design system compliance ✅
- Key learnings application ✅

#### 4. Acceptance Criteria Updates (CRITICAL)
- ✅ **All acceptance criteria use RFC 2119 keywords** (MUST/SHOULD/MAY)
- ✅ Criteria are **specific and measurable** with exact metrics
- ✅ **Negative test scenarios** included
- ✅ **5-category framework** (FR, IR, EC, PR, UX)

Examples:
- Before: "✅ Cursor always visible when moving"
- After: "**FR-4.1**: Cursor MUST remain visible during all cursor movement operations"

#### 5. Test Coverage (CRITICAL)
- ✅ **>90% coverage target** for viewport integration code
- ✅ **>95% coverage target** for critical cursor visibility logic
- ✅ Explicit coverage targets per task
- ✅ Referenced FOUNDATION-TESTING-GUIDE.md throughout

#### 6. Integration Points (CRITICAL)
Added comprehensive integration section:
- **M5 (Auto-save)**: Scroll position should be part of auto-save state ✅
- **M6 (Undo/Redo)**: Scroll position should be part of undo state ✅
- **M34-M35 (Vim Mode)**: j/k keys alignment with vim state machine ✅
- **M37 (Responsive Layout)**: Viewport resize handling ✅
- Testing strategies for each integration point ✅

#### 7. Code Examples (CRITICAL)
All code examples include:
- ✅ Correct import statements (`github.com/charmbracelet/bubbles/viewport`)
- ✅ Export comments on all exported functions
- ✅ Error handling patterns from learnings/error-handling.md
- ✅ Proper error wrapping with %w
- ✅ No stuttering in names
- ✅ All examples compile correctly

#### 8. Key Learnings Integration (CRITICAL)
- ✅ Referenced ui-domain.md Category 2 (Cursor and Viewport Management)
- ✅ Applied middle-third scrolling strategy from ui-domain.md
- ✅ Referenced editor-domain.md for editor integration patterns
- ✅ Applied error handling patterns from error-handling.md

#### 9. Navigation Guide (RECOMMENDED)
- ✅ Comprehensive "How to Use This Document" section
- ✅ When to read each section
- ✅ Line number references for key sections
- ✅ Related document cross-references

---

## Document Statistics

### task-list.md (27KB, ~600 lines)
- 7 detailed tasks with dependencies
- 141 acceptance criteria using RFC 2119 keywords
- Pre-implementation checklist (50+ items)
- Integration points section
- File impact summary
- Timeline estimates (2.5 hours)
- Risk assessment
- Rollback plan

### reference.md (32KB, ~600 lines)
- Navigation guide
- Architecture context
- Complete API reference and mappings
- Middle-third scrolling visual diagrams
- 4 complete code examples with proper imports/error handling
- 3 comprehensive testing pattern examples
- Troubleshooting section with 4 common issues
- Rollback scenarios

---

## Key Improvements Over Original Plan

1. **Structure**: Two focused documents instead of one monolithic file
2. **Traceability**: Every requirement tied to specific acceptance criteria with RFC 2119 keywords
3. **Testing**: 90-95% coverage targets, comprehensive test examples
4. **Integration**: Documented dependencies on M5, M6, M34-M35, M37 for future work
5. **Completeness**: ALL required documents referenced and applied
6. **Usability**: Navigation guide helps developers find what they need quickly
7. **Examples**: All code examples compile and follow style guide
8. **Patterns**: Applied proven patterns from ui-domain.md and editor-domain.md

---

## How to Use This Plan

### For Implementation
1. When implementing **Task 5 of Milestone 4**, use this plan
2. Start with `task-list.md` for task breakdown and acceptance criteria
3. Reference `reference.md` for:
   - Architecture context (before Task 2)
   - API mappings (during implementation)
   - Code examples (when writing code)
   - Testing patterns (before writing tests)
   - Troubleshooting (if issues arise)

### For Review
1. Review `task-list.md` for comprehensive task breakdown
2. Review `reference.md` for technical depth and correctness
3. Verify all RFC 2119 acceptance criteria are met
4. Check code examples compile correctly
5. Verify integration points are considered for future milestones

### For Integration with M4 Plan
This viewport migration plan is **Task 5** in the M4 milestone plan. When implementing M4:

1. Follow M4 task-list.md for overall milestone workflow
2. When reaching Task 5 (Viewport with Scrolling), use this detailed plan
3. The acceptance criteria in this plan align with M4's Task 5 requirements
4. Integration points documented here help with future milestones (M5, M6, M34-M35, M37)

---

## Approval

**Status**: ⏳ AWAITING APPROVAL

**Approver**: Kyle Davis  
**Date Approved**: _________________  
**Signature**: _________________

---

## Next Steps After Approval

1. Execute Task 1 (Dependency Setup)
2. Create implementation branch: `m4-viewport-bubbles-migration`
3. Follow tasks sequentially (1→2→3→4→5→6→7)
4. Stop after each task for human verification (create checkpoint)
5. Request review after Task 5 (Test Migration)
6. Complete documentation and archive in Tasks 6-7
7. Integrate results into M4 milestone completion

---

**Document Version**: 1.0  
**Last Updated**: 2026-01-09  
**Status**: Ready for implementation
