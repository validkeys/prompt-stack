## CLI Commands for Knowledge Management

```bash
# Initialize knowledge base for repo
$ prompt-stack init
Creates .prompt-stack/knowledge.db

# Guided learning session
$ prompt-stack learn --interactive
Asks questions, scans codebase, builds cache

# Automatic background analysis
$ prompt-stack learn --auto
Silent analysis, no questions

# Show what's cached
$ prompt-stack knowledge list
Style Anchors:
  - AuthService.ts (service pattern)
  - user.schema.ts (Zod schema pattern)
  - auth.middleware.ts (middleware pattern)

Standards:
  - TypeScript: strict mode required
  - Testing: Jest with describe/it
  - Validation: Zod for external data

File Patterns:
  - Services: src/services/*.service.ts
  - Schemas: src/schemas/*.schema.ts
  - Tests: **/*.test.ts

# Add specific reference
$ prompt-stack knowledge add-anchor src/services/EmailService.ts \
  --category service \
  --description "Email service with queuing"

# Remove outdated reference
$ prompt-stack knowledge remove-anchor src/old/LegacyAuth.ts

# Export knowledge base (share with team)
$ prompt-stack knowledge export > team-knowledge.json

# Import knowledge base (onboard new dev)
$ prompt-stack knowledge import team-knowledge.json

# Validate cached knowledge (files still exist?)
$ prompt-stack knowledge validate
✓ 12/12 style anchors valid
✗ 2 standards reference deleted files
  Recommend: prompt-stack knowledge clean

# Search cached knowledge
$ prompt-stack knowledge search "authentication"
Found 3 matches:
  - AuthService.ts (style anchor)
  - "No any types" (standard)
  - auth.middleware.ts (style anchor)
```

## Smart Caching Strategies

### 1. **Confidence Scoring**
```javascript
// When tool auto-detects patterns
{
  "pattern": "Service uses dependency injection",
  "confidence": 0.85,  // High confidence
  "evidence": [
    "Constructor takes dependencies",
    "No direct instantiation",
    "Follows same pattern as 5 other services"
  ]
}

// Low confidence = ask user to verify
{
  "pattern": "Tests use snapshot testing",
  "confidence": 0.45,  // Low confidence
  "evidence": [
    "Found 2 .snap files",
    "Not consistent across tests"
  ],
  "action": "ASK_USER"  // Prompt for confirmation
}
```

### 2. **Usage-Based Prioritization**
```sql
-- Frequently used anchors suggested first
SELECT * FROM style_anchors 
WHERE category = 'service'
ORDER BY usage_count DESC, last_used DESC
LIMIT 3;

-- When generating YAML, use top 3 most relevant anchors
```

### 3. **Staleness Detection**
```javascript
// Check if cached files still exist
const staleAnchors = await db.query(`
  SELECT file_path FROM style_anchors 
  WHERE codebase_id = ? 
  AND file_path NOT IN (
    SELECT path FROM git_ls_files()
  )
`);

if (staleAnchors.length > 0) {
  console.warn("⚠ 3 cached references point to deleted files");
  console.log("Run: prompt-stack knowledge clean");
}
```

### 4. **Context Budget Prediction**
```javascript
// Learn from history what context size tasks need
const avgContext = await db.query(`
  SELECT AVG(context_tokens) 
  FROM task_history 
  WHERE task_description LIKE '%authentication%'
  AND success = true
`);

// Use this to set realistic context budgets
taskConfig.max_context_tokens = Math.ceil(avgContext * 1.2);
```

## The Power Move: Team Knowledge Sharing

```bash
# Developer 1 (senior) builds knowledge base
$ prompt-stack learn --interactive
[Answers all the questions, builds comprehensive cache]

$ prompt-stack knowledge export > .prompt-stack/team-knowledge.json
$ git add .prompt-stack/team-knowledge.json
$ git commit -m "Add team coding knowledge base"
$ git push

# Developer 2 (junior) clones repo
$ git clone repo
$ prompt-stack knowledge import .prompt-stack/team-knowledge.json

✓ Imported team knowledge:
  - 15 style anchors
  - 23 coding standards
  - 8 file patterns

# Developer 2 immediately benefits from senior's knowledge
$ prompt-stack plan "Add new feature"
✓ Using team patterns automatically
✓ No questions needed
✓ Generated plan matches team standards
```

## Advanced: Learn from Execution

```javascript
// After Ralphy completes a task
async function learnFromExecution(taskResult) {
  const { task, success, filesModified, contextUsed } = taskResult;
  
  // Store history
  await db.insert('task_history', {
    task_description: task.title,
    files_touched: JSON.stringify(filesModified),
    context_tokens: contextUsed,
    success: success,
    failure_reason: success ? null : task.error
  });
  
  // If successful, extract new patterns
  if (success) {
    for (const file of filesModified) {
      if (!await isCached(file)) {
        const shouldCache = await askUser(
          `${file} was successfully modified. Cache as style anchor?`
        );
        if (shouldCache) {
          await cacheStyleAnchor(file, await analyzeFile(file));
        }
      }
    }
  }
  
  // If failed, identify missing knowledge
  if (!success && task.error.includes("unknown pattern")) {
    console.log("💡 This failure suggests missing knowledge");
    console.log("Run: prompt-stack knowledge add-anchor <file>");
  }
}
```

## Storage Location Strategy

```bash
# Option 1: Per-repo (in .gitignore)
.prompt-stack/
  knowledge.db          # SQLite cache
  team-knowledge.json   # Exportable team knowledge (committed)

# Option 2: Global with repo mapping
~/.prompt-stack/
  knowledge.db          # All repos in one database
  
# Hybrid (best approach)
~/.prompt-stack/
  global-knowledge.db   # Cross-repo patterns
  
.prompt-stack/
  repo-knowledge.db     # This repo's specific patterns
  team-knowledge.json   # Shareable subset
```

## Why This Is Powerful

**Without JIT Discovery:**
```bash
$ prompt-stack plan "Add feature"
ERROR: Please specify:
  - Style anchor files (--anchor)
  - Coding standards (--standards)
  - Test patterns (--test-pattern)
  - Dependencies allowed (--deps)
```
→ User overwhelmed, gives up

**With JIT Discovery:**
```bash
$ prompt-stack plan "Add feature"

❓ Quick question: Do you have similar code? [Y/n]
> Y

📂 Found 3 candidates. Use AuthService.ts as reference? [Y/n]  
> Y

✓ Generated plan with learned patterns
  (Next time: zero questions, uses cache)
```
→ Friction minimized, knowledge compounds