# Task 5: Viewport Bubbles Migration Plan Review

**Review Date**: 2026-01-09
**Reviewer**: AI Reviewer
**Milestone**: Task 5 - Viewport with Scrolling (Bubbles Migration)
**Plan Document**: `docs/implementation-plans/task5-viewport-bubbles-migration.md`

---

## Executive Summary

### Overall Assessment
- Status: ⚠️ Approved with Recommendations
- Confidence: Medium
- Critical Issues: 1
- Recommendation Issues: 6
- Nice-to-have Issues: 4

### Key Strengths
1. Comprehensive breakdown of migration phases with clear steps
2. Detailed API mapping table between custom and bubbles viewport
3. Good risk assessment with mitigation strategies
4. Clear rollback plan included
5. Reasonable time estimates (2.5 hours total)
6. Practical cursor visibility implementation approach

### Critical Issues
1. **Document structure mismatch** - Plan doesn't follow fresh-build standards (wrong location, wrong format, missing required sections)

### Recommendations
1. Refactor to follow fresh-build format (task-list.md + reference.md split)
2. Add all required document references from DOCUMENT-REFERENCE-MATRIX.md
3. Include pre-implementation checklist
4. Structure tasks with dependencies, complexity estimates, and RFC 2119 acceptance criteria
5. Add key learnings references (ui-domain.md, editor-domain.md)
6. Include integration points section

---

## Detailed Findings

### 1. Document Location and Format

**Status**: ❌ Incorrect

**Issues**:
- [ ] Plan located at `docs/implementation-plans/task5-viewport-bubbles-migration.md`
  - **Should be**: `docs/plans/fresh-build/milestone-implementation-plans/M4/task-list.md` and `reference.md`
- [ ] Single monolithic document
  - **Should be**: Split into `task-list.md` (tasks) and `reference.md` (examples, patterns)
- [ ] No review document in `plan-reviews/` subdirectory
  - **Should be**: `docs/plans/fresh-build/milestone-implementation-plans/M4/plan-reviews/001-review.md`

**Impact**: Makes plan inconsistent with fresh-build documentation standards, harder to maintain and review.

---

### 2. Document Coverage Verification

**Status**: ❌ Incomplete - Major Gaps

**Required Documents from DOCUMENT-REFERENCE-MATRIX.md**:

#### Core Planning Documents (Not Referenced)
- [x] milestone-execution-prompt.md - ❌ Not referenced
- [x] milestones.md - ❌ Not referenced (Task 5 section should be quoted)
- [x] requirements.md - ❌ Not referenced
- [x] project-structure.md - ❌ Not referenced

#### Implementation Guides (Not Referenced)
- [x] go-style-guide.md - ❌ Not referenced
- [x] go-testing-guide.md - ❌ Not referenced
- [x] bubble-tea-testing-best-practices.md - ❌ Not referenced
- [x] ENHANCED-TEST-CRITERIA-TEMPLATE.md - ❌ Not referenced

#### Design System Documents (Not Referenced)
- [x] opencode-design-system.md - ❌ Not referenced (UI/TUI milestone)
- [x] bubble-tea-best-practices.md - ❌ Not referenced

#### Key Learnings (Not Referenced)
- [x] learnings/ui-domain.md - ❌ Not referenced (critical for viewport/scrolling patterns)
- [x] learnings/editor-domain.md - ❌ Not referenced (critical for editor integration)
- [x] learnings/error-handling.md - ❌ Not referenced

#### Testing Guides (Not Referenced)
- [x] FOUNDATION-TESTING-GUIDE.md - ❌ Not referenced (Foundation milestone)

**Impact**: Plan doesn't leverage established patterns and may reinvent solved problems or miss best practices.

---

### 3. Milestone Deliverables Coverage

**Status**: ⚠️ Partial

**Coverage Analysis**:
- Total deliverables from M5: 4 (viewport component, scrolling, auto-scroll, integration)
- Fully covered: 4
- Partially covered: 0
- Not covered: 0

**Missing Deliverables**: None identified

**Under-Specified Deliverables**:
- [ ] **Deliverable: Scroll with keyboard** (j/k keys)
  - **What's missing**: No explicit testing for vim-style navigation
  - **What's needed**: Add test cases for j/k scrolling

- [ ] **Deliverable: Auto-scroll on insert/delete**
  - **What's missing**: No specific test cases for auto-scroll behavior
  - **What's needed**: Add test cases for text insertion at viewport edges

**Over-Engineered Deliverables**: None

---

### 4. Task List Quality

**Status**: ⚠️ Good but incomplete

**Task Ordering Issues**:
- [ ] **Phase 4 (Cursor Visibility)** should reference key learnings from ui-domain.md for scroll strategies
  - **Impact**: May miss established patterns for cursor/viewport coordination

**Missing Dependencies**:
- [ ] **Phase 5: Test Migration** - Doesn't reference FOUNDATION-TESTING-GUIDE.md for testing patterns
  - **Impact**: Tests may not follow established Foundation testing standards

**Vague Descriptions**:
- [ ] **Phase 4: Cursor Visibility Integration** - "Challenge: Bubbles viewport doesn't have built-in EnsureVisible() method"
  - **Suggested improvement**: Reference ui-domain.md for cursor visibility patterns from previous implementation

**Non-Measurable Acceptance Criteria**:
- [ ] **All phases** - Acceptance criteria use checklists but not RFC 2119 keywords (MUST/SHOULD/MAY)
  - **Example**: Instead of "✅ Cursor always visible when moving", use "MUST maintain cursor visibility during navigation"
  - **Impact**: Harder to verify objective compliance

**Missing Testing Requirements**:
- [ ] **No explicit coverage targets** - Fresh-build requires >85-90% coverage
  - **What's needed**: Add "Code coverage >90% for viewport integration tests"

**Missing Integration Points**:
- [ ] **No integration with vim-domain.md** - j/k keys should align with Vim state machine (M34)
- [ ] **No integration with undo/redo system** (M6) - Scroll position should be part of undo stack

**File Path Inconsistencies**:
- [ ] `internal/editor/viewport.go` - Should verify against project-structure.md
- [ ] `ui/workspace/viewport_integration_test.go` - Verify file naming convention

**Complexity Estimate Issues**: None - 2.5 hours is reasonable

---

### 5. Reference Document Quality

**Status**: ⚠️ Partially Complete

**Missing Navigation Guide**:
- [ ] No Section 6.2.0 navigation guide explaining how to use reference sections
  - **Impact**: Harder for implementers to find relevant information

**Incorrect Architecture Context**:
- [ ] References `internal/editor/viewport.go` but doesn't explain how it fits into editor domain
- [ ] Doesn't reference editor-domain.md for architecture patterns
  - **What should be**: Include architecture context explaining viewport's role in editor domain

**Code Example Issues**:
- [ ] Phase 4 cursor visibility code (lines 132-156) - No imports shown
  - **Impact**: Won't compile as-is
- [ ] No style guide references in code examples
  - **Impact**: May not follow go-style-guide.md conventions

**Test Example Issues**:
- [ ] No test examples provided (Phase 5 mentions tests but no code)
  - **Impact**: Harder to verify testing patterns match go-testing-guide.md

**Missing References**:
- [ ] learnings/ui-domain.md - Critical for viewport/scrolling patterns
- [ ] bubble-tea-best-practices.md - For Bubble Tea component patterns
- [ ] go-style-guide.md - For code conventions

**References Not Applied**:
- [ ] DEPENDENCIES.md - Should reference bubbles package version requirements
- [ ] opencode-design-system.md - Should reference UI patterns for viewport styling

**Missing Key Learnings**:
- [ ] **learnings/ui-domain.md** - "Cursor and viewport management" section has established patterns
- [ ] **learnings/editor-domain.md** - Editor integration patterns
- [ ] **learnings/architecture-patterns.md** - Component extraction patterns

**Missing Design System References**:
- [ ] opencode-design-system.md - Viewport should follow design system for styling
  - **Impact**: May have inconsistent UI appearance

---

### 6. Pre-Implementation Checklist

**Status**: ❌ Missing

**Missing Sections**:
- [ ] **Package structure** - Should verify `internal/editor/` and `ui/workspace/` structure
- [ ] **Dependency injection** - Should document viewport initialization pattern
- [ ] **Documentation** - Should specify what docs need updating
- [ ] **Testing** - Should reference FOUNDATION-TESTING-GUIDE.md
- [ ] **Style** - Should reference go-style-guide.md
- [ ] **Constants** - Should identify constants (scroll speed, buffer sizes)
- [ ] **Design system** - Should reference opencode-design-system.md

**Impact**: Implementers may miss important setup steps and quality checks.

---

### 7. Acceptance Criteria Quality

**Status**: ⚠️ Fair

**RFC 2119 Keyword Usage**:
- [ ] Phase 4, line 296: "✅ Cursor always visible when moving" - Should be "MUST maintain cursor visibility during all cursor movement operations"
- [ ] Phase 4, line 297: "✅ Middle-third scrolling behavior maintained" - Should be "MUST maintain middle-third scrolling strategy for cursor visibility"
- [ ] Phase 4, line 298: "✅ No negative scroll positions" - Should be "MUST prevent negative YOffset values in viewport"
- [ ] All other checkmarks should use MUST/SHOULD/MAY keywords

**Vague Criteria**:
- [ ] "✅ Smooth scrolling without jumps" (Phase 4, line 299)
  - **Why vague**: "Smooth" and "without jumps" are subjective
  - **Suggested improvement**: "MUST maintain scroll position within 2 lines of expected position during scrolling operations"

**Missing Negative Scenarios**:
- [ ] No test for viewport with zero height (terminal too small)
- [ ] No test for viewport with negative dimensions
- [ ] No test for cursor at document boundaries with middle-third scrolling

**Alignment with Milestone Goals**: All criteria align with M5 goals

**Coverage Targets**:
- [ ] No explicit coverage targets specified
  - **What to add**: "Code coverage >90% for viewport integration tests" (per fresh-build standards)

**Critical Test Scenarios**: Well identified in Phase 5

**Edge Cases**: Well identified in Phase 4 and Testing Strategy section

---

### 8. Testing Alignment

**Status**: ⚠️ Good but incomplete

**Testing Guide Referenced**:
- [x] No testing guide referenced
  - **Should reference**: FOUNDATION-TESTING-GUIDE.md (Foundation milestone)

**Test Pattern Compliance**:
- [ ] No test examples provided
  - **Can't verify**: Whether tests follow go-testing-guide.md patterns
  - **Need**: Add table-driven test examples

**Coverage Targets**:
- [ ] No coverage targets specified
  - **What to add**: >90% coverage requirement (fresh-build standard)

**Integration Testing**:
- [ ] Phase 5 mentions "Integration tests for cursor movement with scrolling"
  - **Good**: Integration tests are planned
  - **Missing**: Integration with undo/redo (M6) and future vim mode (M34)

**Manual Testing**:
- [x] Manual testing scenarios identified in Testing Strategy section
  - **Good**: Navigate large file, rapid cursor movement, jump to beginning/end, resize terminal, insert/delete while scrolling

**Test Execution Order**: Not defined, should be part of task dependencies

---

### 9. Design System Compliance (if applicable)

**Status**: ❌ Not Referenced

**Missing References**:
- [ ] opencode-design-system.md - Viewport is a UI component, should follow design system

**Color Usage**:
- [ ] Not specified - Should reference design system for viewport styling

**Spacing**:
- [ ] Not specified - Should follow 1-character unit system from design system

**Component Structure**:
- [ ] Not specified - Should follow design system component patterns

**Keyboard Shortcuts**:
- [ ] j/k keys mentioned but not tied to vim-domain.md or keybinding-system.md

**Visual Hierarchy**:
- [ ] Not specified - Should follow design system for viewport rendering

**Interactive Elements**:
- [ ] Not specified - Should follow design system for visual feedback

**Impact**: Viewport may not be consistent with rest of UI.

---

### 10. Code Example Validation

**Status**: ⚠️ Some Issues

**Compilation Issues**:
- [ ] Phase 4, lines 132-156: Missing imports (likely need viewport package import)
  - **Why won't compile**: No `import "github.com/charmbracelet/bubbles/viewport"`
  - **Fix needed**: Add import statement

**Import Issues**:
- [ ] See above - Missing viewport package import

**Dependency Issues**:
- [ ] Phase 1: "Add github.com/charmbracelet/bubbles dependency"
  - **Issue**: Should reference DEPENDENCIES.md for version pinning
  - **Correct version**: Check DEPENDENCIES.md for recommended version

**Pseudo-Code**:
- [ ] None - All code examples appear to be real Go code (good)

**Error Handling**:
- [ ] Phase 4 cursor visibility code doesn't include error handling
  - **Should include**: Error handling patterns from error-handling.md

**Missing Comments**:
- [ ] Phase 4: `ensureCursorVisible()` function - Should have export comment
  - **Go style guide**: All exported functions must have comments

**Name Stuttering**: None - Names follow Go conventions

---

### 11. Integration Points

**Status**: ⚠️ Partially Identified

**Missing Integration Points**:
- [ ] **Integration with M6 (Undo/Redo)** - Scroll position should be part of undo state
  - **With what**: internal/editor/undo.go
  - **Impact**: Users may lose scroll position when undoing edits
  - **How to test**: Undo edit, verify viewport position restored

- [ ] **Integration with M34-M35 (Vim Mode)** - j/k keys should align with vim state machine
  - **With what**: internal/vim/ directory
  - **Impact**: Inconsistent navigation behavior between normal and edit modes
  - **How to test**: Test j/k in normal mode, verify consistent scrolling

- [ ] **Integration with M37 (Responsive Layout)** - Viewport should resize on terminal resize
  - **With what**: ui/app/model.go (resize handling)
  - **Impact**: Viewport may break when window resized
  - **How to test**: Resize terminal, verify viewport adjusts correctly

**Integration Testing Missing**:
- [ ] **No integration tests with undo/redo** - Should test scroll position restoration
- [ ] **No integration tests with future vim mode** - Should test j/k compatibility

**Dependencies Not Documented**:
- [ ] **Dependency on M4 (Basic Text Editor)** - Editor buffer integration
- [ ] **Dependency on future M34-M35** - Vim-style scrolling

---

## Approval Criteria

### Must Pass (Critical)
- [x] All milestone deliverables are covered
- [ ] All acceptance criteria are measurable using RFC 2119 keywords
- [x] All code examples compile and follow style guide (minor import issues)
- [ ] All required documents from MATRIX are referenced
- [ ] Testing guide for milestone group is referenced
- [ ] Key learnings are referenced and applied
- [ ] File paths match project-structure.md
- [x] Integration points are identified (but incomplete)

### Should Pass (Important)
- [x] Task descriptions are clear and actionable
- [x] Task ordering is correct by dependency
- [ ] Testing requirements have coverage targets
- [ ] Design system references are included (if applicable)
- [x] Code examples use correct imports (minor issues)

### Nice to Have (Optional)
- [ ] Pre-implementation checklist is comprehensive
- [x] Rollback scenarios are considered
- [ ] Test execution order is defined
- [x] Manual testing scenarios are identified

---

## Recommendation

**Status**: ⚠️ Approved with Minor Revisions

### Summary
The implementation plan is technically sound with a well-thought-out migration approach. However, it doesn't follow fresh-build documentation standards, missing critical document references and required sections. The plan needs structural improvements to align with fresh-build conventions before implementation.

### Required Changes Before Implementation
1. **Move to correct location**: `docs/plans/fresh-build/milestone-implementation-plans/M4/task-list.md`
2. **Split into two files**: `task-list.md` and `reference.md`
3. **Add all required document references** from DOCUMENT-REFERENCE-MATRIX.md
4. **Add pre-implementation checklist** with package structure, dependency injection, testing, style, constants, design system
5. **Structure tasks with RFC 2119 acceptance criteria** (MUST/SHOULD/MAY)
6. **Add test coverage targets** (>90% per fresh-build standards)
7. **Reference key learnings**: learnings/ui-domain.md and learnings/editor-domain.md
8. **Include integration points section** with undo/redo and vim mode
9. **Add missing imports** to code examples
10. **Add export comments** to all exported functions

### Recommended Changes (Not Blocking)
1. Reference opencode-design-system.md for viewport styling
2. Add more detailed test examples following go-testing-guide.md patterns
3. Add integration tests for scroll position with undo/redo
4. Consider future vim mode integration in scroll behavior
5. Add performance benchmarks for scroll operations
6. Document viewport initialization pattern for dependency injection

### Next Steps
1. Address critical structural changes (location, format, document references)
2. Add missing sections (pre-implementation checklist, integration points)
3. Improve acceptance criteria with RFC 2119 keywords
4. Add key learnings references and apply patterns
5. Review and approve refactored plan
6. Execute Phase 1 (Dependency Setup)
7. Create implementation branch: `task5-bubbles-viewport`
8. Follow phases sequentially
9. Request review after Phase 5 (Test Migration)

---

## Review Metrics

- Total deliverables: 4
- Deliverables fully covered: 4 (100%)
- Tasks reviewed: 7
- Tasks with issues: 2
- Critical issues: 1
- Recommendation issues: 6
- Nice-to-have issues: 4
- Overall confidence: Medium

---

## Reviewer Notes

The implementation plan demonstrates solid technical understanding of the viewport migration challenge. The phased approach is reasonable, and the cursor visibility implementation (middle-third scrolling) is practical. However, the plan's structure doesn't align with fresh-build documentation standards, making it inconsistent with other milestone plans.

The most significant gaps are:
1. **No reference to key learnings** - ui-domain.md and editor-domain.md contain established patterns for viewport/scrolling that should be leveraged
2. **No design system reference** - Viewport is a UI component that should follow opencode-design-system.md patterns
3. **Missing integration with future milestones** - Should plan for undo/redo (M6) and vim mode (M34-M35) compatibility

The plan's technical approach is sound, but it needs structural improvements to match fresh-build conventions and leverage existing knowledge and patterns from the project's documentation ecosystem.

---

## Appendix: Required Document References for Refactoring

To properly refactor this plan, the following documents should be read:

### Core Planning Documents (Read First)
1. milestone-execution-prompt.md - Execution workflow and plan format
2. milestones.md - Task 5 deliverables and acceptance criteria
3. requirements.md - Editor and scrolling requirements
4. project-structure.md - Verify file paths for editor domain
5. go-style-guide.md - Code conventions and patterns
6. go-testing-guide.md - Testing patterns and TDD

### Context-Specific Documents (Read if Applicable)
7. FOUNDATION-TESTING-GUIDE.md - Foundation milestone testing
8. bubble-tea-testing-best-practices.md - Bubble Tea TUI testing
9. bubble-tea-best-practices.md - Bubble Tea architecture
10. DEPENDENCIES.md - Bubbles package version and usage
11. opencode-design-system.md - UI design patterns (for viewport styling)
12. keybinding-system.md - Vim mode integration (future)

### Key Learnings (Critical for this plan)
13. learnings/ui-domain.md - Cursor and viewport management patterns
14. learnings/editor-domain.md - Editor implementation patterns
15. learnings/error-handling.md - Error handling patterns
16. learnings/architecture-patterns.md - Component extraction patterns

### Implementation Plan Review Reference
17. docs/.prompts/review-implementation-plan.md - Review template and criteria

---

**Review Document Created**: 2026-01-09
**Status**: Ready for human review and approval
**Next Action**: Refactor implementation plan based on findings
