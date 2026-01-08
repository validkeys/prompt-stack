# Bubble Tea Compliance - Implementation Plan

**Audit:** 001-bubble-tea-compliance  
**Created:** 2026-01-08  
**Status:** Ready for Implementation

## Overview

Implementation plan for addressing Bubble Tea compliance issues identified in [`audit-tracking.md`](./audit-tracking.md:1). Focus: fix critical architectural violations, add effect-based testing, refactor workspace model.

## Priority Matrix

| Phase | Focus | Timeline | Risk |
|-------|-------|----------|------|
| P1 | State mutations + critical tests | Week 1-2 | High |
| P2 | Workspace refactor + coverage | Week 3-4 | Medium |
| P3 | Performance + infrastructure | Month 2 | Low |

## Phase 1: Critical Fixes (Week 1-2)

### 1.1 Fix State Mutations
**Files:** [`app/model.go`](../../../../archive/code/ui/app/model.go:1), [`workspace/model.go`](../../../../archive/code/ui/workspace/model.go:1), [`palette/model.go`](../../../../archive/code/ui/palette/model.go:1)

**Pattern:**
```go
// BEFORE: m.field = value; return m, nil
// AFTER: newM := m; newM.field = value; return newM, nil
```

**Code Examples:** See [`references.md#code-examples`](./references.md:8) for before/after examples for each file

**Tasks:**
- [ ] Convert pointer receivers to value receivers in Update methods
- [ ] Ensure all Update paths return new model instances
- [ ] Add copy helpers for complex models ([`references.md#copy-helper-pattern`](./references.md:68))
- [ ] Verify no direct mutations remain

### 1.2 Create Test Utilities
**Location:** `internal/testutil/bubbletea.go`

**Required Helpers:**
```go
TypeText(model, text) tea.Model
PressKey(model, keyType) tea.Model
AssertState(t, got, want)
SimulateWorkflow(model, steps) tea.Model
```

**Tasks:**
- [ ] Create `internal/testutil/` package
- [ ] Implement typing simulation helper
- [ ] Implement key press helper
- [ ] Implement state assertion helper
- [ ] Add workflow simulation helper

**Verification:** See [`references.md#verification-commands`](./references.md:438) for test commands

### 1.3 Add Critical Tests
**Target:** [`workspace/model.go`](../../../../archive/code/ui/workspace/model.go:1), [`app/model.go`](../../../../archive/code/ui/app/model.go:1)

**Test Types:**
1. **Effect Tests** - Verify state changes
2. **Command Tests** - Verify commands returned
3. **Message Tests** - Cover all message types

**Test Templates:** See [`references.md#test-templates`](./references.md:76) for complete test function examples

**Tasks:**
- [ ] Create `workspace/model_test.go` with effect tests ([`references.md#effect-testing-template`](./references.md:76))
- [ ] Create `app/model_test.go` with command tests ([`references.md#command-verification-template`](./references.md:95))
- [ ] Add table-driven message tests ([`references.md#message-coverage-template`](./references.md:119))
- [ ] Test edge cases ([`references.md#edge-case-template`](./references.md:145), [`references.md#specific-test-cases`](./references.md:453))

## Phase 2: Refactoring (Week 3-4)

### 2.1 Decompose Workspace Model
**Target:** Reduce from 1292 → <400 lines

**Extraction Plan:**
| Component | New Location | Lines | Responsibility |
|-----------|-------------|-------|----------------|
| Cursor | `internal/editor/cursor.go` | ~150 | Position, navigation |
| Viewport | `internal/editor/viewport.go` | ~200 | Scroll, visible lines |
| Placeholder | `internal/editor/placeholder.go` | ~250 | Detection, editing |
| FileIO | `internal/editor/fileio.go` | ~100 | Save/load ops |

**Component Specs:** See [`references.md#component-interface-specifications`](./references.md:197) for complete interface definitions

**Integration:** See [`references.md#integration-points`](./references.md:357) for workspace refactoring example

**Tasks:**
- [ ] Extract cursor component with tests ([`references.md#cursor-component`](./references.md:197))
- [ ] Extract viewport component with tests ([`references.md#viewport-component`](./references.md:247))
- [ ] Extract placeholder system with tests ([`references.md#placeholder-component`](./references.md:297))
- [ ] Extract file I/O operations with tests ([`references.md#file-io-component`](./references.md:357))
- [ ] Update workspace to use components ([`references.md#workspace-model-after-refactoring`](./references.md:357))
- [ ] Verify integration tests pass

**Dependencies:** See [`references.md#task-dependencies`](./references.md:423) for task ordering

### 2.2 Improve Test Coverage
**Target:** 80% coverage across all models

**Priority Models:**
1. [`workspace/model.go`](../../../../archive/code/ui/workspace/model.go:1) - 0% → 80%
2. [`suggestions/model.go`](../../../../archive/code/ui/suggestions/model.go:1) - 0% → 80%
3. [`palette/model.go`](../../../../archive/code/ui/palette/model.go:1) - 0% → 80%

**Tasks:**
- [ ] Add unit tests for all Update methods
- [ ] Add integration tests for workflows
- [ ] Add edge case tests
- [ ] Add performance regression tests

## Phase 3: Advanced (Month 2)

### 3.1 Performance Optimization
**Tasks:**
- [ ] Profile view rendering for large docs
- [ ] Optimize string operations in Update
- [ ] Implement efficient diff algorithms
- [ ] Add rendering caching

### 3.2 Testing Infrastructure
**Tasks:**
- [ ] Set up CI with coverage reporting
- [ ] Add benchmark tests
- [ ] Create golden file tests for views
- [ ] Implement property-based testing

### 3.3 Developer Experience
**Tasks:**
- [ ] Document testing patterns
- [ ] Create model template generator
- [ ] Add linting rules for Bubble Tea
- [ ] Create comprehensive test helpers

## Success Metrics

| Metric | Current | Target | Phase |
|--------|---------|--------|-------|
| Test Coverage | 30% | 80% | P2 |
| State Mutations | 15+ | 0 | P1 |
| Workspace Lines | 1292 | <400 | P2 |
| Effect Tests | 0% | 100% | P1 |
| Arch Compliance | 70% | 95% | P2 |

## Implementation Checklist

### Week 1
- [ ] Create testutil package
- [ ] Fix state mutations in app model
- [ ] Fix state mutations in palette model
- [ ] Add effect tests for workspace basics

### Week 2
- [ ] Fix state mutations in workspace model
- [ ] Add command verification tests
- [ ] Add message coverage tests
- [ ] Begin cursor extraction

### Week 3
- [ ] Complete cursor extraction
- [ ] Extract viewport component
- [ ] Extract placeholder system
- [ ] Update workspace integration

### Week 4
- [ ] Extract file I/O operations
- [ ] Achieve 80% test coverage
- [ ] Verify all integration tests pass
- [ ] Document new architecture

## Risk Mitigation

| Risk | Mitigation | Owner |
|------|-----------|-------|
| Breaking changes | Comprehensive integration tests | Dev |
| Performance regression | Benchmark tests before/after | Dev |
| Test complexity | Use testutil helpers | Dev |
| Incomplete refactor | Phase-based approach | PM |

## Validation Criteria

**Phase 1 Complete When:**
- ✅ Zero state mutation violations
- ✅ All critical tests passing
- ✅ Testutil package functional

**Phase 2 Complete When:**
- ✅ Workspace <400 lines
- ✅ 80% test coverage achieved
- ✅ All integration tests passing

**Phase 3 Complete When:**
- ✅ CI pipeline operational
- ✅ Performance benchmarks established
- ✅ Documentation complete

## References

- [`references.md`](./references.md:1) - **Detailed implementation examples, test templates, component specs**
- [`audit-tracking.md`](./audit-tracking.md:1) - Full audit findings
- [`bubble-tea-testing-best-practices.md`](../../bubble-tea-testing-best-practices.md:1) - Testing guidelines
- [`bubble-tea-best-practices.md`](../../../../bubble-tea-best-practices.md:1) - Architecture guidelines
- [`go-testing-guide.md`](../../go-testing-guide.md:1) - Go testing patterns

**Quick Links to References:**
- [`Code Examples`](./references.md:8) - Before/after mutation fixes
- [`Test Templates`](./references.md:76) - Complete test function examples
- [`Component Specs`](./references.md:197) - Interface definitions for extracted components
- [`Integration Points`](./references.md:357) - Workspace refactoring example
- [`Task Dependencies`](./references.md:423) - Task ordering and dependencies
- [`Verification Commands`](./references.md:438) - Commands to verify each phase
- [`Specific Test Cases`](./references.md:453) - 10+ edge cases to test
- [`Error Handling Strategy`](./references.md:483) - Debugging and rollback procedures

## Notes

- All file paths relative to `/Users/kyledavis/Sites/prompt-stack`
- Use testutil helpers for all new tests
- Follow immutable Update pattern strictly
- Document architectural decisions
- Update this plan as work progresses

---

**Status:** Ready for Implementation  
**Next Review:** After Phase 1 completion  
**Owner:** Development Team