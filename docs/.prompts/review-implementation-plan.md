# Implementation Plan Review Prompt

## Document References (Shorthand)
- **REVIEW**: review-implementation-plan.md (this file)
- **EXEC**: milestone-execution-prompt.md
- **MATRIX**: DOCUMENT-REFERENCE-MATRIX.md
- **MILESTONES**: milestones.md
- **TEMPLATE**: milestones/ENHANCED-TEST-CRITERIA-TEMPLATE.md
- **STYLE**: go-style-guide.md
- **TESTING**: go-testing-guide.md
- **STRUCT**: project-structure.md
- **DESIGN**: opencode-design-system.md
- **WORKFLOW**: DOCUMENT-CHECKING-WORKFLOW.md

---

## 1. Your Role

You are a senior Go developer and technical architect reviewing implementation plans for PromptStack. Your role is to thoroughly evaluate implementation plans for completeness, accuracy, and alignment with milestone requirements.

---

## 2. Process Flow

```
Load Milestone → Read Required Documents → Read Implementation Plan → 
Comprehensive Review → Generate Review Document → Output Summary
```

---

## 3. Step 1: Load Milestone Context

### 3.1 Read Core Documents

Read the following documents in parallel (up to 5 at a time):

```xml
<read_file>
<args>
  <file>
    <path>docs/plans/fresh-build/milestones.md</path>
  </file>
  <file>
    <path>docs/plans/fresh-build/DOCUMENT-REFERENCE-MATRIX.md</path>
  </file>
  <file>
    <path>docs/plans/fresh-build/go-style-guide.md</path>
  </file>
  <file>
    <path>docs/plans/fresh-build/go-testing-guide.md</path>
  </file>
  <file>
    <path>docs/plans/fresh-build/project-structure.md</path>
  </file>
</args>
</read_file>
```

### 3.2 Read Context-Specific Documents

Based on the milestone, read the required documents from the **MATRIX**. These may include:
- Design system documents (for UI milestones)
- Specific testing guides
- Detailed acceptance criteria documents
- Architecture specifications
- API specifications
- Database schemas
- Configuration schemas

---

## 4. Step 2: Read Implementation Plan

### 4.1 Read Plan Documents

Read the implementation plan for the given milestone:

```xml
<read_file>
<args>
  <file>
    <path>docs/plans/fresh-build/milestone-implementation-plans/M{N}/task-list.md</path>
  </file>
  <file>
    <path>docs/plans/fresh-build/milestone-implementation-plans/M{N}/reference.md</path>
  </file>
</args>
</read_file>
```

**Note**: If reference.md is split into multiple parts (reference-part-2.md, etc.), read all parts.

---

## 5. Step 3: Comprehensive Review

Perform a thorough review of the implementation plan against the following criteria:

### 5.1 Milestone Deliverables Coverage

**Review Checklist**:
- [ ] All deliverables from milestones.md are addressed
- [ ] Each deliverable has corresponding tasks in the task list
- [ ] No deliverables are missed or omitted
- [ ] Deliverable complexity is appropriately reflected in task estimates

**Findings to Document**:
- Missing deliverables (list them)
- Deliverables with insufficient coverage
- Deliverables that are over-engineered or under-planned

### 5.2 Task List Quality

**Review Checklist**:
- [ ] Tasks are ordered correctly by dependency
- [ ] Each task has a clear, descriptive title
- [ ] Task dependencies are correctly identified
- [ ] Task complexity estimates are reasonable
- [ ] Task descriptions are clear and actionable
- [ ] Each task has specific, measurable acceptance criteria
- [ ] Acceptance criteria use RFC 2119 keywords (MUST/SHOULD/MAY)
- [ ] Testing requirements are specified for each task
- [ ] Integration points are explicitly called out
- [ ] File paths match project-structure.md exactly

**Findings to Document**:
- Tasks with incorrect ordering
- Tasks with missing or incorrect dependencies
- Tasks with vague or incomplete descriptions
- Tasks with non-measurable acceptance criteria
- Tasks missing testing requirements
- Integration points not identified
- File path inconsistencies

### 5.3 Reference Document Quality

**Review Checklist**:
- [ ] Navigation guide is present and clear (Section 6.2.0)
- [ ] Architecture context is accurate
- [ ] Style guide references are correct with examples
- [ ] Testing guide references are accurate
- [ ] Key learnings are referenced and applied
- [ ] Design system references are included (if applicable)
- [ ] Code examples compile and follow style guide
- [ ] Code examples use correct imports from DEPENDENCIES.md
- [ ] Test examples match code signatures
- [ ] Test examples follow testing guide patterns
- [ ] Rollback scenarios are considered

**Findings to Document**:
- Missing navigation guide
- Incorrect architecture context
- Code examples that won't compile
- Code examples not following style guide
- Test examples not matching code
- Missing references to key learnings
- Missing design system references (for UI milestones)

### 5.4 Document Coverage Verification

**Review Checklist**:
- [ ] All required documents from MATRIX are referenced
- [ ] Core planning documents are referenced (STYLE, TESTING, STRUCT)
- [ ] Context-specific documents are referenced
- [ ] Testing guide for the milestone group is referenced
- [ ] Detailed acceptance criteria document is referenced (if applicable)
- [ ] Key learnings are referenced
- [ ] Key learnings are applied in the plan

**Findings to Document**:
- Required documents not referenced
- Documents referenced but not applied
- Missing references to testing guides
- Missing references to acceptance criteria documents

### 5.5 Pre-Implementation Checklist

**Review Checklist**:
- [ ] Pre-implementation checklist is present (Section 6.1.1)
- [ ] All checklist items are relevant to the milestone
- [ ] Checklist is comprehensive for the tasks planned
- [ ] Checklist covers package structure, dependency injection, documentation, testing, style, constants, design system

**Findings to Document**:
- Missing pre-implementation checklist
- Incomplete checklist sections
- Checklist items not relevant to tasks

### 5.6 Acceptance Criteria Quality

**Review Checklist**:
- [ ] Acceptance criteria use RFC 2119 keywords consistently
- [ ] Criteria are specific and measurable
- [ ] Criteria include both positive and negative scenarios
- [ ] Criteria define exact formats and constraints
- [ ] Criteria align with milestone goals
- [ ] Testing requirements have coverage targets
- [ ] Critical test scenarios are identified
- [ ] Edge cases are listed

**Findings to Document**:
- Acceptance criteria that are vague or subjective
- Missing negative test scenarios
- Acceptance criteria that don't align with milestone goals
- Missing coverage targets
- Insufficient critical test scenarios
- Edge cases not identified

### 5.7 Testing Alignment

**Review Checklist**:
- [ ] Testing guide for the milestone group is referenced
- [ ] Test examples follow testing guide patterns
- [ ] Table-driven tests are used where appropriate
- [ ] Test coverage targets are specified
- [ ] Integration testing is considered
- [ ] Manual testing scenarios are identified
- [ ] Test execution order is defined

**Findings to Document**:
- Testing guide not referenced
- Test patterns not following guide
- Missing coverage targets
- No integration testing plan
- Manual testing not considered

### 5.8 Design System Compliance (if applicable)

**Review Checklist**:
- [ ] Design system guidelines are referenced
- [ ] Color palette usage follows guidelines
- [ ] Spacing follows 1-character unit system
- [ ] Component structure matches patterns
- [ ] Keyboard shortcuts are supported
- [ ] Visual hierarchy matches examples
- [ ] Interactive elements have visual feedback

**Findings to Document**:
- Missing design system references
- Color usage violations
- Spacing inconsistencies
- Component structure deviations
- Missing keyboard shortcuts

### 5.9 Code Example Validation

**Review Checklist**:
- [ ] All code examples compile (if extracted)
- [ ] Imports match actual package usage
- [ ] Dependencies are correct per DEPENDENCIES.md
- [ ] No pseudo-code in examples
- [ ] Error handling follows error-handling.md patterns
- [ ] Comments on exported functions
- [ ] No stuttering in names
- [ ] Test examples match code signatures

**Findings to Document**:
- Code that won't compile
- Incorrect imports
- Wrong dependency names
- Pseudo-code in examples
- Missing error handling
- Name stuttering issues

### 5.10 Integration Points

**Review Checklist**:
- [ ] Integration points with existing code are identified
- [ ] Integration points with other milestones are identified
- [ ] Integration testing is planned
- [ ] Dependencies on previous milestones are noted
- [ ] Dependencies for future milestones are noted

**Findings to Document**:
- Missing integration points
- Integration points not tested
- Dependencies not documented

---

## 6. Step 4: Generate Review Document

Create a comprehensive review document at:
**Location**: `docs/plans/fresh-build/milestone-implementation-plans/M{N}/plan-reviews/{00N}-review.md`

### 6.1 Review Document Structure

```markdown
# Milestone {N} Implementation Plan Review

**Review Date**: {Date}
**Reviewer**: {AI Reviewer}
**Milestone**: {Milestone Number} - {Milestone Title}

---

## Executive Summary

### Overall Assessment
- Status: {✅ Approved | ⚠️ Approved with Recommendations | ❌ Needs Revision}
- Confidence: {High | Medium | Low}
- Critical Issues: {Number}

### Key Strengths
1. [Strength 1]
2. [Strength 2]
3. [Strength 3]

### Critical Issues
1. [Critical Issue 1 - must be fixed before implementation]
2. [Critical Issue 2 - must be fixed before implementation]

### Recommendations
1. [Recommendation 1]
2. [Recommendation 2]
3. [Recommendation 3]

---

## Detailed Findings

### 1. Milestone Deliverables Coverage

**Status**: {✅ Complete | ⚠️ Partial | ❌ Incomplete}

**Coverage Analysis**:
- Total deliverables: {N}
- Fully covered: {N}
- Partially covered: {N}
- Not covered: {N}

**Missing Deliverables**:
- [ ] [Deliverable] - [Why it's missing, impact]

**Under-Specified Deliverables**:
- [ ] [Deliverable] - [What's missing, what's needed]

**Over-Engineered Deliverables**:
- [ ] [Deliverable] - [Why it's over-planned,建议 simplification]

---

### 2. Task List Quality

**Status**: {✅ Excellent | ⚠️ Good | ⚠️ Fair | ❌ Poor}

**Task Ordering Issues**:
- [ ] [Task N] - [Issue with ordering, suggested fix]

**Missing Dependencies**:
- [ ] [Task N] - [Missing dependency, impact]

**Vague Descriptions**:
- [ ] [Task N] - [Why description is vague, suggested improvement]

**Non-Measurable Acceptance Criteria**:
- [ ] [Task N] - [Which criteria are non-measurable, suggested fix]

**Missing Testing Requirements**:
- [ ] [Task N] - [What testing is missing, what's needed]

**Missing Integration Points**:
- [ ] [Task N] - [Integration points not identified, impact]

**File Path Inconsistencies**:
- [ ] [File path in plan] vs [Correct path per STRUCT] - [Impact, fix]

**Complexity Estimate Issues**:
- [ ] [Task N] - [Why estimate is wrong, suggested adjustment]

---

### 3. Reference Document Quality

**Status**: {✅ Excellent | ⚠️ Good | ⚠️ Fair | ❌ Poor}

**Missing Navigation Guide**:
- [Issue description, impact]

**Incorrect Architecture Context**:
- [What's incorrect, what should be, impact]

**Code Example Issues**:
- [ ] [Example in reference.md line X-Y] - [Why it won't compile, fix needed]
- [ ] [Example in reference.md line X-Y] - [Style guide violation, fix needed]
- [ ] [Example in reference.md line X-Y] - [Incorrect import, fix needed]

**Test Example Issues**:
- [ ] [Test example in reference.md line X-Y] - [Doesn't match code signature, fix needed]
- [ ] [Test example in reference.md line X-Y] - [Doesn't follow testing guide, fix needed]

**Missing References**:
- [ ] [Document type] - [Which document, why it's needed]

**References Not Applied**:
- [ ] [Document/section] - [Referenced but not applied, how to apply]

**Missing Key Learnings**:
- [ ] [Key learning from learnings/] - [What learning, why it's relevant, how to apply]

**Missing Design System References**:
- [ ] [Design system pattern] - [Which pattern, why needed, impact]

---

### 4. Document Coverage Verification

**Status**: {✅ Complete | ⚠️ Partial | ❌ Incomplete}

**Required Documents from MATRIX**:
- [ ] [Document] - {✅ Referenced | ❌ Not referenced} - [Impact if missing]

**Core Planning Documents**:
- [ ] STYLE (go-style-guide.md) - {✅ Referenced | ❌ Not referenced}
- [ ] TESTING (go-testing-guide.md) - {✅ Referenced | ❌ Not referenced}
- [ ] STRUCT (project-structure.md) - {✅ Referenced | ❌ Not referenced}

**Context-Specific Documents**:
- [ ] [Document] - {✅ Referenced | ❌ Not referenced} - [Impact if missing]

**Testing Guide**:
- [ ] [Testing guide for milestone group] - {✅ Referenced | ❌ Not referenced} - [Impact if missing]

**Detailed Acceptance Criteria**:
- [ ] [ACCEPTANCE-CRITERIA-*.md] - {✅ Referenced | ❌ Not referenced | N/A} - [Impact if missing]

**Key Learnings**:
- [ ] [Key learning document] - {✅ Referenced | ❌ Not referenced} - [Impact if missing]
- [ ] [Key learning document] - {✅ Applied | ⚠️ Referenced but not applied | ❌ Not referenced}

---

### 5. Pre-Implementation Checklist

**Status**: {✅ Complete | ⚠️ Partial | ❌ Missing}

**Missing Sections**:
- [ ] [Section] - [Why needed, what to add]

**Incomplete Sections**:
- [ ] [Section] - [What's missing, what to add]

**Irrelevant Items**:
- [ ] [Item] - [Why not relevant, can be removed]

---

### 6. Acceptance Criteria Quality

**Status**: {✅ Excellent | ⚠️ Good | ⚠️ Fair | ❌ Poor}

**RFC 2119 Keyword Usage**:
- [ ] [Task N] - [Criteria not using keywords, suggested fix]

**Vague Criteria**:
- [ ] [Task N, Criteria] - [Why vague, suggested improvement]

**Missing Negative Scenarios**:
- [ ] [Task N] - [What negative tests are missing, impact]

**Alignment with Milestone Goals**:
- [ ] [Task N, Criteria] - [Doesn't align, suggested fix]

**Coverage Targets**:
- [ ] [Task N] - [Missing coverage target, what to add]

**Critical Test Scenarios**:
- [ ] [Task N] - [Missing critical scenarios, what to add]

**Edge Cases**:
- [ ] [Task N] - [Missing edge cases, what to add]

---

### 7. Testing Alignment

**Status**: {✅ Excellent | ⚠️ Good | ⚠️ Fair | ❌ Poor}

**Testing Guide Referenced**:
- [ ] [Testing guide] - {✅ Referenced | ❌ Not referenced} - [Impact if missing]

**Test Pattern Compliance**:
- [ ] [Test example] - [Doesn't follow guide, how to fix]

**Coverage Targets**:
- [ ] [Task N] - [Missing coverage target, what to add]

**Integration Testing**:
- [ ] [Task N] - [No integration testing plan, what's needed]

**Manual Testing**:
- [ ] [Manual testing not considered] - [What manual tests are needed]

**Test Execution Order**:
- [ ] [Test execution order not defined] - [How to define]

---

### 8. Design System Compliance (if applicable)

**Status**: {✅ Compliant | ⚠️ Partially Compliant | ❌ Non-Compliant | N/A}

**Missing References**:
- [ ] [Design system pattern] - [What's missing, impact]

**Color Usage**:
- [ ] [Component/element] - [Color violation, what's correct]

**Spacing**:
- [ ] [Component/element] - [Spacing violation, what's correct]

**Component Structure**:
- [ ] [Component] - [Structure deviation, correct pattern]

**Keyboard Shortcuts**:
- [ ] [Missing shortcut] - [What shortcut needed, impact]

**Visual Hierarchy**:
- [ ] [Element] - [Hierarchy issue, correct approach]

**Interactive Elements**:
- [ ] [Element] - [Missing feedback, what to add]

---

### 9. Code Example Validation

**Status**: {✅ All Valid | ⚠️ Some Issues | ❌ Many Issues}

**Compilation Issues**:
- [ ] [Example in reference.md line X-Y] - [Why won't compile, fix needed]

**Import Issues**:
- [ ] [Example in reference.md line X-Y] - [Wrong import, correct import]

**Dependency Issues**:
- [ ] [Example in reference.md line X-Y] - [Wrong dependency, correct one]

**Pseudo-Code**:
- [ ] [Example in reference.md line X-Y] - [Contains pseudo-code, real code needed]

**Error Handling**:
- [ ] [Example in reference.md line X-Y] - [Missing error handling, add pattern from error-handling.md]

**Missing Comments**:
- [ ] [Function name] - [Exported function missing comment]

**Name Stuttering**:
- [ ] [Name] - [Stuttering, better name]

---

### 10. Integration Points

**Status**: {✅ Well Identified | ⚠️ Partially Identified | ❌ Not Identified}

**Missing Integration Points**:
- [ ] [Integration] - [With what, impact, how to test]

**Integration Testing Missing**:
- [ ] [Integration point] - [No testing plan, what's needed]

**Dependencies Not Documented**:
- [ ] [Dependency] - [With what milestone, impact]

---

## Approval Criteria

### Must Pass (Critical)
- [ ] All milestone deliverables are covered
- [ ] All acceptance criteria are measurable using RFC 2119 keywords
- [ ] All code examples compile and follow style guide
- [ ] All required documents from MATRIX are referenced
- [ ] Testing guide for milestone group is referenced
- [ ] Key learnings are referenced and applied
- [ ] File paths match project-structure.md
- [ ] Integration points are identified

### Should Pass (Important)
- [ ] Task descriptions are clear and actionable
- [ ] Task ordering is correct by dependency
- [ ] Testing requirements have coverage targets
- [ ] Design system references are included (if applicable)
- [ ] Code examples use correct imports

### Nice to Have (Optional)
- [ ] Pre-implementation checklist is comprehensive
- [ ] Rollback scenarios are considered
- [ ] Test execution order is defined
- [ ] Manual testing scenarios are identified

---

## Recommendation

**Status**: {✅ Approved for Implementation | ⚠️ Approved with Minor Revisions | ❌ Requires Major Revisions}

### Summary
[Brief summary of the overall assessment and recommendation]

### Required Changes Before Implementation
1. [Critical change 1]
2. [Critical change 2]
3. [Critical change 3]

### Recommended Changes (Not Blocking)
1. [Recommendation 1]
2. [Recommendation 2]
3. [Recommendation 3]

### Next Steps
- If **Approved**: Proceed with implementation following the plan
- If **Approved with Minor Revisions**: Address minor revisions, then proceed
- If **Requires Major Revisions**: Address major revisions, resubmit for review

---

## Review Metrics

- Total deliverables: {N}
- Deliverables fully covered: {N} ({N}%)
- Tasks reviewed: {N}
- Tasks with issues: {N}
- Critical issues: {N}
- Recommendation issues: {N}
- Nice-to-have issues: {N}
- Overall confidence: {High | Medium | Low}

---

## Reviewer Notes

[Additional notes, context, or observations that don't fit in other sections]
```

---

## 7. Step 5: Output Summary

After generating the review document, output a concise summary:

```
📋 Milestone {N} Implementation Plan Review Complete

📄 Review Document: docs/plans/fresh-build/milestone-implementation-plans/M{N}/plan-reviews/{00N}-review.md

📊 Review Summary:
- Status: {✅ Approved | ⚠️ Approved with Recommendations | ❌ Needs Revision}
- Confidence: {High | Medium | Low}
- Critical Issues: {N}
- Recommendation Issues: {N}
- Overall Assessment: {1-2 sentence summary}

🔍 Key Findings:
✅ [Strength 1]
✅ [Strength 2]
⚠️ [Issue 1]
⚠️ [Issue 2]

🎯 Recommendation:
[Brief recommendation]

View the full review document for detailed findings and action items.
```

---

## 8. Getting Started

When you receive this prompt:

1. **Confirm you understand the review process**
2. **Ask which milestone to review** (e.g., "Which milestone implementation plan should I review?")
3. **Load the milestone context** (Step 3)
4. **Read the implementation plan** (Step 4)
5. **Perform comprehensive review** (Step 5)
6. **Generate review document** (Step 6)
7. **Output summary** (Step 7)

**Your first response should be**:
```
I understand the implementation plan review process. I will:
- Load milestone context from required documents
- Read the implementation plan (task-list.md and reference.md)
- Perform a comprehensive review across 10 categories
- Generate a detailed review document
- Provide a clear recommendation

Which milestone implementation plan should I review?
```

---

**Remember**: Thorough reviews ensure implementation plans are complete, accurate, and aligned with milestone requirements. This prevents implementation issues and ensures quality deliverables.
