# Enhanced Test Criteria Template

**Purpose**: Standardized template for comprehensive test criteria across all milestones.

**Usage**: Use this template when updating milestones.md with enhanced test criteria.

---

## Template Structure

```markdown
## Milestone X: [Milestone Name]
**Goal:** [Brief description of milestone objective]

**Deliverables:**
- [Deliverable 1]
- [Deliverable 2]
- [Deliverable 3]

**Test Criteria:**

### Functional Requirements
- [ ] [Specific functional test case]
- [ ] [Specific functional test case]
  - [ ] [Sub-test or edge case]
  - [ ] [Sub-test or edge case]
- [ ] [Specific functional test case]

### Integration Requirements
- [ ] [Integration test with previous milestone]
- [ ] [Integration test with specific component]
- [ ] [Integration test with specific component]

### Edge Cases & Error Handling
- [ ] [Edge case scenario]
- [ ] [Error condition handling]
- [ ] [Boundary condition test]
- [ ] [Invalid input handling]

### Performance Requirements
- [ ] [Performance metric with threshold]
- [ ] [Performance metric with threshold]
- [ ] [Resource usage constraint]

### User Experience Requirements
- [ ] [UX requirement]
- [ ] [UX requirement]
- [ ] [Accessibility or usability requirement]

**Files:** [`path/to/file1.go`], [`path/to/file2.go`]
```

---

## Section Guidelines

### Functional Requirements

**Purpose**: Verify that the feature works as specified.

**Guidelines**:
- **Specific**: Each test should be independently verifiable
- **Measurable**: Include concrete values where applicable (e.g., "<500ms", "100 levels")
- **Actionable**: Clear pass/fail criteria
- **Complete**: Cover all deliverables

**Examples**:
- ✅ Good: "App launches without errors on fresh install"
- ✅ Good: "Config file created at `~/.promptstack/config.yaml` with correct structure"
- ✅ Good: "Setup wizard prompts for API key and preferences"
  - [ ] Validates API key format before accepting
  - [ ] Shows error message for invalid keys
  - [ ] Allows re-entry of invalid fields
- ❌ Bad: "App works"
- ❌ Bad: "Config is good"

**Sub-tests**: Use nested checkboxes for related test cases:
```
- [ ] Setup wizard prompts for API key and preferences
  - [ ] Validates API key format before accepting
  - [ ] Shows error message for invalid keys
  - [ ] Allows re-entry of invalid fields
```

---

### Integration Requirements

**Purpose**: Verify that the feature integrates correctly with other components.

**Guidelines**:
- Test integration with previously completed milestones
- Test integration with other components in same milestone group
- Verify data flow between components
- Test state management across components

**Examples**:
- ✅ Good: "Config loading integrates with logging system"
- ✅ Good: "Setup wizard integrates with config persistence"
- ✅ Good: "Version tracking integrates with starter prompt extraction"
- ❌ Bad: "Works with other stuff"

**Integration Types**:
1. **Previous Milestone Integration**: Test with features from earlier milestones
2. **Same Group Integration**: Test with other features in the same milestone group
3. **Cross-Component Integration**: Test data flow between different components
4. **State Integration**: Test state management across components

---

### Edge Cases & Error Handling

**Purpose**: Verify robustness under unusual conditions and error scenarios.

**Guidelines**:
- **Boundary conditions**: Empty inputs, maximum values, minimum values
- **Invalid inputs**: Wrong types, malformed data, unexpected formats
- **Error conditions**: Network failures, file system errors, API errors
- **Concurrent operations**: Multiple simultaneous actions
- **Resource constraints**: Low memory, slow disk, network timeout

**Examples**:
- ✅ Good: "Handle missing config directory (create automatically)"
- ✅ Good: "Handle corrupted config file (show error, offer reset)"
- ✅ Good: "Handle invalid API key format (reject with clear message)"
- ✅ Good: "Handle read-only filesystem (show error, exit gracefully)"
- ✅ Good: "Handle interrupted setup wizard (resume on next launch)"
- ❌ Bad: "Handle errors"

**Edge Case Categories**:

1. **Boundary Conditions**:
   - Empty inputs
   - Maximum values (e.g., 100 undo levels)
   - Minimum values (e.g., 0 items)
   - Zero-length strings
   - Negative numbers

2. **Invalid Inputs**:
   - Wrong data types
   - Malformed YAML/JSON
   - Invalid file paths
   - Unexpected characters
   - Missing required fields

3. **Error Conditions**:
   - Network failures
   - File system errors (permission denied, disk full)
   - API errors (401, 429, 500)
   - Database connection failures
   - Timeout errors

4. **Concurrent Operations**:
   - Multiple simultaneous saves
   - Rapid keyboard input
   - Concurrent API requests
   - Simultaneous file operations

5. **Resource Constraints**:
   - Low memory conditions
   - Slow disk I/O
   - Network latency
   - Large datasets (1000+ items)

---

### Performance Requirements

**Purpose**: Verify that the feature meets performance expectations.

**Guidelines**:
- **Response time**: Maximum acceptable latency (e.g., "<10ms for fuzzy search")
- **Throughput**: Operations per second (e.g., "1000+ prompts/second")
- **Resource usage**: Memory limits, CPU usage, disk I/O
- **Scalability**: Behavior with large datasets (e.g., "100+ prompts", "1000+ history entries")

**Examples**:
- ✅ Good: "App startup time <500ms on fresh install"
- ✅ Good: "Config file read/write <50ms"
- ✅ Good: "Log file write <10ms per entry"
- ✅ Good: "Fuzzy search filters in real-time (<10ms)"
- ✅ Good: "Startup time <500ms for 100 prompts"
- ❌ Bad: "Fast enough"
- ❌ Bad: "Good performance"

**Performance Metrics**:

1. **Response Time**:
   - Startup time
   - Operation latency
   - UI responsiveness
   - API response time

2. **Throughput**:
   - Operations per second
   - Requests per second
   - Items processed per second

3. **Resource Usage**:
   - Memory footprint
   - CPU usage
   - Disk I/O
   - Network bandwidth

4. **Scalability**:
   - Behavior with 100+ items
   - Behavior with 1000+ items
   - Behavior with 10000+ items
   - Performance degradation rate

**Performance Testing Tips**:
- Use realistic data sizes
- Test on target hardware
- Measure multiple times and average
- Document test conditions (hardware, data size, etc.)

---

### User Experience Requirements

**Purpose**: Verify that the feature provides a good user experience.

**Guidelines**:
- **Visual feedback**: Status indicators, progress indicators, error messages
- **Keyboard shortcuts**: Consistent, discoverable, documented
- **Accessibility**: Clear error messages, readable text, keyboard navigation
- **Usability**: Intuitive workflows, minimal steps, clear affordances

**Examples**:
- ✅ Good: "Setup wizard provides clear instructions"
- ✅ Good: "Error messages are actionable and specific"
- ✅ Good: "Progress indicators shown during initialization"
- ✅ Good: "Keyboard navigation works in setup wizard"
- ✅ Good: "Status bar shows 'Saving...' then 'Saved'"
- ❌ Bad: "Good UX"
- ❌ Bad: "User friendly"

**UX Categories**:

1. **Visual Feedback**:
   - Status indicators (loading, saving, error)
   - Progress indicators (percentage, progress bar)
   - Error messages (clear, actionable, specific)
   - Success confirmations

2. **Keyboard Shortcuts**:
   - Consistent with common patterns
   - Discoverable (displayed in UI)
   - Documented in help
   - No conflicts

3. **Accessibility**:
   - Clear error messages
   - Readable text (contrast, size)
   - Keyboard navigation
   - Screen reader support (if applicable)

4. **Usability**:
   - Intuitive workflows
   - Minimal steps to complete tasks
   - Clear affordances (what can be clicked/pressed)
   - Consistent behavior across components

---

## Complete Example: Milestone 1

```markdown
## Milestone 1: Bootstrap & Config
**Goal:** Initialize application foundation

**Deliverables:**
- Config structure at `~/.promptstack/config.yaml`
- First-run interactive setup wizard
- Logging setup with zap
- Version tracking

**Test Criteria:**

### Functional Requirements
- [ ] App launches without errors on fresh install
- [ ] Config file created at `~/.promptstack/config.yaml` with correct structure
- [ ] Setup wizard prompts for API key and preferences
  - [ ] Validates API key format before accepting
  - [ ] Shows error message for invalid keys
  - [ ] Allows re-entry of invalid fields
- [ ] Logs written to `~/.promptstack/debug.log`
  - [ ] Log file created with correct permissions (0600)
  - [ ] Log entries include timestamp and level
  - [ ] Log rotation works at 10MB limit
- [ ] Version stored in config and compared on startup

### Integration Requirements
- [ ] Config loading integrates with logging system
- [ ] Setup wizard integrates with config persistence
- [ ] Version tracking integrates with starter prompt extraction

### Edge Cases & Error Handling
- [ ] Handle missing config directory (create automatically)
- [ ] Handle corrupted config file (show error, offer reset)
- [ ] Handle invalid API key format (reject with clear message)
- [ ] Handle read-only filesystem (show error, exit gracefully)
- [ ] Handle interrupted setup wizard (resume on next launch)

### Performance Requirements
- [ ] App startup time <500ms on fresh install
- [ ] Config file read/write <50ms
- [ ] Log file write <10ms per entry

### User Experience Requirements
- [ ] Setup wizard provides clear instructions
- [ ] Error messages are actionable and specific
- [ ] Progress indicators shown during initialization
- [ ] Keyboard navigation works in setup wizard

**Files:** [`internal/config/config.go`], [`internal/setup/wizard.go`], [`internal/logging/logger.go`], [`cmd/promptstack/main.go`]
```

---

## Best Practices

### Writing Good Test Criteria

1. **Be Specific**: Avoid vague terms like "works" or "good"
2. **Be Measurable**: Include concrete values and thresholds
3. **Be Actionable**: Each test should have clear pass/fail criteria
4. **Be Complete**: Cover all deliverables and edge cases
5. **Be Consistent**: Use similar language and structure across milestones

### Common Mistakes to Avoid

1. ❌ Too vague: "App works well"
2. ❌ Not measurable: "Fast performance"
3. ❌ Missing edge cases: Only testing happy path
4. ❌ No integration tests: Testing in isolation only
5. ❌ No performance criteria: No benchmarks or thresholds

### Review Checklist

Before finalizing test criteria for a milestone, verify:

- [ ] All deliverables have corresponding test criteria
- [ ] Functional requirements cover all features
- [ ] Integration requirements test with previous milestones
- [ ] Edge cases cover boundary conditions
- [ ] Error handling covers failure scenarios
- [ ] Performance requirements have measurable thresholds
- [ ] UX requirements cover user experience aspects
- [ ] All tests are specific and actionable
- [ ] All tests are independently verifiable
- [ ] No duplicate or redundant tests

---

## Template Usage Workflow

1. **Copy Template**: Start with the template structure
2. **Fill Deliverables**: List all milestone deliverables
3. **Write Functional Tests**: Create tests for each deliverable
4. **Add Integration Tests**: Identify integration points
5. **Identify Edge Cases**: List boundary conditions and error scenarios
6. **Define Performance Metrics**: Set measurable performance targets
7. **Specify UX Requirements**: Document user experience expectations
8. **Review**: Use the review checklist to verify completeness
9. **Update milestones.md**: Apply the enhanced test criteria

---

**Last Updated**: 2026-01-07
**Version**: 1.0