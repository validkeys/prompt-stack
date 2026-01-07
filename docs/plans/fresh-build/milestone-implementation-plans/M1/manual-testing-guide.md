# Milestone 1 Manual Testing Guide

**Milestone**: 1 - Bootstrap & Config
**Purpose**: Verify application functionality through manual testing
**Prerequisites**: Go 1.21+ installed, built application

---

## Test Environment Setup

### 1. Clean Test Environment

Before testing, ensure you have a clean environment:

```bash
# Remove existing config directory (if any)
rm -rf ~/.promptstack

# Verify removal
ls ~/.promptstack  # Should return "No such file or directory"
```

### 2. Build the Application

```bash
# Build the application
go build -o promptstack ./cmd/promptstack

# Verify build
./promptstack --help  # Should show help or run bootstrap
```

---

## Test Scenario 1: Fresh Install (No Config)

### Objective
Verify the setup wizard runs correctly when no config exists.

### Steps

1. **Run the application**
   ```bash
   ./promptstack
   ```

2. **Expected Behavior**
   - Application should detect missing config
   - Setup wizard should start automatically
   - You should see: "Welcome to PromptStack! Let's set up your configuration."

3. **Verify Setup Wizard Prompts**

   **Prompt 1: API Key**
   - Expected: "Enter your Claude API key (starts with 'sk-ant-'):"
   - Test valid input: `sk-ant-api03-test-key-12345`
   - Test invalid input: `invalid-key` (should reject with error)
   - Test empty input: Press Enter (should prompt again)

   **Prompt 2: Model Selection**
   - Expected: Numbered list of models (1, 2, 3, etc.)
   - Test valid selection: Enter `1` or `2`
   - Test invalid selection: Enter `99` (should reject)
   - Test empty input: Press Enter (should use default)

   **Prompt 3: Vim Mode**
   - Expected: "Enable vim mode? (y/N):"
   - Test valid input: `y` or `n`
   - Test empty input: Press Enter (should default to N)

4. **Verify Configuration Summary**
   - Should display all entered values
   - API key should be masked (show only first 8 chars)
   - Should ask for confirmation

5. **Confirm and Save**
   - Enter `y` to confirm
   - Should see: "Configuration saved to ~/.promptstack/config.yaml"

### Success Criteria
- [ ] Setup wizard starts automatically
- [ ] All prompts display correctly
- [ ] Invalid input is rejected with clear error messages
- [ ] Configuration summary shows masked API key
- [ ] Config file is created at `~/.promptstack/config.yaml`

---

## Test Scenario 2: Config File Verification

### Objective
Verify the config file is created correctly with proper structure.

### Steps

1. **Check config file exists**
   ```bash
   ls -la ~/.promptstack/config.yaml
   ```
   Expected: File exists with proper permissions

2. **View config file contents**
   ```bash
   cat ~/.promptstack/config.yaml
   ```

3. **Verify YAML Structure**
   Expected format:
   ```yaml
   version: "1.0.0"
   api_key: "sk-ant-api03-test-key-12345"
   model: "claude-3-5-sonnet-20241022"
   vim_mode: false
   log_level: "info"
   ```

4. **Verify directory permissions**
   ```bash
   ls -ld ~/.promptstack
   ```
   Expected: Directory with permissions 0755

5. **Verify file permissions**
   ```bash
   ls -l ~/.promptstack/config.yaml
   ```
   Expected: File with appropriate read/write permissions

### Success Criteria
- [ ] Config file exists at correct path
- [ ] YAML is valid and properly formatted
- [ ] All fields are present with correct values
- [ ] Version field is set to "1.0.0"
- [ ] Directory and file permissions are correct

---

## Test Scenario 3: Existing Config (No Wizard)

### Objective
Verify application loads existing config without running wizard.

### Steps

1. **Run the application again**
   ```bash
   ./promptstack
   ```

2. **Expected Behavior**
   - Setup wizard should NOT run
   - Application should load existing config
   - Should see: "Configuration loaded successfully"
   - Application should exit gracefully (or show TUI if implemented)

3. **Verify no new config created**
   ```bash
   ls -la ~/.promptstack/
   ```
   Expected: Only one config.yaml file (no duplicates)

### Success Criteria
- [ ] Setup wizard does not run
- [ ] Existing config is loaded
- [ ] No duplicate config files created
- [ ] Application starts without errors

---

## Test Scenario 4: Log File Inspection

### Objective
Verify logging infrastructure works correctly.

### Steps

1. **Check log file exists**
   ```bash
   ls -la ~/.promptstack/debug.log
   ```
   Expected: Log file exists

2. **View log file contents**
   ```bash
   cat ~/.promptstack/debug.log
   ```

3. **Verify Log Format**
   Expected JSON format with structured fields:
   ```json
   {"level":"info","timestamp":"2026-01-07T20:00:00.000Z","msg":"Logger initialized"}
   {"level":"info","timestamp":"2026-01-07T20:00:00.001Z","msg":"Configuration loaded successfully"}
   ```

4. **Verify Log Levels**
   - Check for INFO level messages (normal operations)
   - Check for ERROR level messages (if any errors occurred)
   - Verify no DEBUG messages unless LOG_LEVEL=debug is set

5. **Test Log Level Override**
   ```bash
   LOG_LEVEL=debug ./promptstack
   cat ~/.promptstack/debug.log | grep "level\":\"debug"
   ```
   Expected: Debug messages appear in log

### Success Criteria
- [ ] Log file is created at correct path
- [ ] Logs are in JSON format
- [ ] Structured fields are present (level, timestamp, msg)
- [ ] Log level can be controlled via LOG_LEVEL env var
- [ ] No string concatenation in log messages

---

## Test Scenario 5: Invalid Config Handling

### Objective
Verify application handles corrupted/invalid config gracefully.

### Steps

1. **Corrupt the config file**
   ```bash
   echo "invalid: yaml: content: [unclosed" > ~/.promptstack/config.yaml
   ```

2. **Run the application**
   ```bash
   ./promptstack
   ```

3. **Expected Behavior**
   - Application should detect invalid YAML
   - Should show clear error message
   - Should exit gracefully with exit code 1

4. **Restore valid config**
   ```bash
   rm ~/.promptstack/config.yaml
   ./promptstack  # Run setup wizard again
   ```

### Success Criteria
- [ ] Invalid config is detected
- [ ] Clear error message is displayed
- [ ] Application exits with error code 1
- [ ] No crash or panic occurs

---

## Test Scenario 6: Environment Variable Overrides

### Objective
Verify config can be overridden via environment variables.

### Steps

1. **Set environment variable**
   ```bash
   export LOG_LEVEL=debug
   ```

2. **Run application**
   ```bash
   ./promptstack
   ```

3. **Check log file**
   ```bash
   cat ~/.promptstack/debug.log | grep "level\":\"debug"
   ```
   Expected: Debug messages appear

4. **Verify config file unchanged**
   ```bash
   cat ~/.promptstack/config.yaml | grep log_level
   ```
   Expected: log_level in file is still "info" (not changed by env var)

### Success Criteria
- [ ] Environment variable overrides work
- [ ] Config file is not modified by env vars
- [ ] Override takes effect at runtime

---

## Test Scenario 7: API Key Validation

### Objective
Verify API key format validation works correctly.

### Steps

1. **Remove config to trigger wizard**
   ```bash
   rm ~/.promptstack/config.yaml
   ```

2. **Run application and test invalid keys**
   ```bash
   ./promptstack
   ```

3. **Test various invalid formats**
   - Empty string (press Enter)
   - `invalid-key`
   - `sk-` (too short)
   - `sk-ant` (missing rest)
   - `sk-ant-api03` (incomplete)

4. **Expected Behavior**
   - Each invalid key should be rejected
   - Error message: "Invalid API key. API key must start with 'sk-ant-'"
   - Wizard should prompt again for valid input

5. **Test valid key**
   - Enter: `sk-ant-api03-test-key-12345`
   - Should be accepted

### Success Criteria
- [ ] Invalid API keys are rejected
- [ ] Clear error messages are shown
- [ ] User can retry without restarting wizard
- [ ] Valid API key format is accepted

---

## Test Scenario 8: Keyboard Interrupt Handling

### Objective
Verify application handles Ctrl+C gracefully.

### Steps

1. **Remove config to trigger wizard**
   ```bash
   rm ~/.promptstack/config.yaml
   ```

2. **Run application**
   ```bash
   ./promptstack
   ```

3. **Press Ctrl+C during wizard**
   - Should see: "Setup cancelled by user"
   - Application should exit gracefully
   - No config file should be created

4. **Verify no partial config**
   ```bash
   ls ~/.promptstack/config.yaml
   ```
   Expected: "No such file or directory"

### Success Criteria
- [ ] Ctrl+C is handled gracefully
   - [ ] Clear cancellation message is shown
   - [ ] No partial config file is created
   - [ ] Application exits cleanly

---

## Test Scenario 9: Version Tracking

### Objective
Verify version tracking works correctly.

### Steps

1. **Check config file for version**
   ```bash
   cat ~/.promptstack/config.yaml | grep version
   ```
   Expected: `version: "1.0.0"`

2. **Test version mismatch (manual simulation)**
   ```bash
   # Edit config to have different version
   sed -i '' 's/version: "1.0.0"/version: "0.9.0"/' ~/.promptstack/config.yaml
   ```

3. **Run application**
   ```bash
   ./promptstack
   ```

4. **Expected Behavior**
   - Should log warning about version mismatch
   - Should still load config (migration placeholder)
   - Should not crash

5. **Restore correct version**
   ```bash
   sed -i '' 's/version: "0.9.0"/version: "1.0.0"/' ~/.promptstack/config.yaml
   ```

### Success Criteria
- [ ] Version is set to "1.0.0" on new config
- [ ] Version mismatch is logged as warning
- [ ] Application continues despite version mismatch
- [ ] No crash occurs on version mismatch

---

## Test Scenario 10: Log Rotation

### Objective
Verify log rotation works correctly (requires generating large logs).

### Steps

1. **Generate large log file**
   ```bash
   # Run application multiple times to generate logs
   for i in {1..100}; do
     LOG_LEVEL=debug ./promptstack
   done
   ```

2. **Check log file size**
   ```bash
   ls -lh ~/.promptstack/debug.log
   ```
   Expected: File size should be reasonable (rotation should kick in at 10MB)

3. **Check for rotated logs**
   ```bash
   ls -lh ~/.promptstack/debug.log*
   ```
   Expected: May see debug.log.1, debug.log.2 if rotation occurred

4. **Verify rotation settings**
   - Max size: 10MB
   - Max backups: 3
   - Max age: 30 days

### Success Criteria
- [ ] Log file is created
- [ ] Rotation occurs when file reaches 10MB
- [ ] Backup files are created with proper naming
- [ ] Old logs are cleaned up after 30 days

---

## Cleanup After Testing

```bash
# Remove test config and logs
rm -rf ~/.promptstack

# Remove built binary
rm -f promptstack
```

---

## Test Results Checklist

After completing all scenarios, verify:

- [ ] All 10 test scenarios completed
- [ ] All success criteria met
- [ ] No crashes or panics observed
- [ ] All error messages are clear and helpful
- [ ] Application behavior is consistent
- [ ] Documentation matches actual behavior

---

## Known Limitations

1. **Log Rotation Testing**: Requires generating 10MB+ of logs, which may take time
2. **Version Migration**: Only placeholder exists; actual migration logic not implemented yet
3. **TUI Functionality**: Not yet implemented (Milestone 2)

---

## Notes

- Record any unexpected behavior or bugs found
- Note any areas where UX could be improved
- Document any deviations from expected behavior
- Take screenshots of any issues for reference

---

**Last Updated**: 2026-01-07
**Tested By**: _______________
**Test Date**: _______________
**Overall Result**: [ ] Pass [ ] Fail