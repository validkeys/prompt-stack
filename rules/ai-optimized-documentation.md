# AI-Optimized Documentation

Guidelines for creating documentation that AI systems can effectively retrieve, understand, and utilize.

## Why This Matters

Poorly structured documentation leads to:
- **Inaccurate retrievals**: AI retrieves wrong documents
- **Incomplete context**: AI gets partial information, leading to incorrect answers
- **Confusion between states**: AI can't distinguish current vs. planned vs. deprecated
- **Poor chunking**: Documents split at illogical boundaries during processing

## Length Guidelines

AI systems process documents in chunks. Document length directly affects retrieval accuracy.

### Target Length
- **Target: 300-800 lines per document** for optimal chunking and retrieval
- **Maximum: 1000 lines** before considering a split
- **Each major section: 50-200 lines** to ensure digestible chunks

### When to Split Documents
Split documents exceeding 1000 lines when:
1. **Multiple distinct topics**: Document covers several unrelated concepts
2. **Clear logical boundaries**: Natural separation points exist
3. **Different audiences**: Some sections target different skill levels or roles
4. **Maintenance burden**: Document becomes difficult to keep current

### How to Split Documents
1. **Identify logical sections**: Find natural topic boundaries
2. **Create focused documents**: Each new document should have a single clear purpose
3. **Maintain self-containment**: Each split document must be independently comprehensible
4. **Add navigation**: Link related documents clearly
5. **Update index**: Ensure the main index reflects the new structure

## Metadata Requirements

Every document should include YAML frontmatter with structured metadata.

### YAML Frontmatter Format

```yaml
---
title: Full Document Title
domain: category/subcategory
keywords: [term1, term2, term3, pattern-name, synonym]
related: [related-doc-1, related-doc-2]
last_updated: 2025-01-05
status: active
---
```

### Required Fields

**title**: Full document title, should match the H1 heading

**domain**: Location in documentation hierarchy using slash notation
- Examples: `system/core`, `features/auth`, `guides/testing`

**keywords**: Array of 3-8 terms for semantic search
- Include both technical terms and domain concepts
- Use terms developers would naturally search for
- Include synonyms and alternative terminology
- Tag architectural patterns explicitly

**related**: Array of related document slugs (without file extension)
- Enables graph-based navigation
- Helps AI understand document relationships

**last_updated**: Date of last significant update (yyyy-mm-dd)
- Used for freshness ranking

**status**: Current state of the document
- Values: `active`, `planned`, `in_progress`, `draft`, `deprecated`, `historical`
- Critical for filtering current vs. future vs. past

## Heading Structure

Clear heading hierarchy is critical for AI to understand document structure.

### Heading Guidelines

1. **Maximum depth of 3 levels**: H1, H2, H3 only
   - Deeper nesting degrades retrieval accuracy

2. **Semantic meaning**: Each heading should be self-descriptive
   - Good: `## Authorization Validation`
   - Bad: `## Step 3`

3. **Logical hierarchy**: Follow document outline structure
   ```markdown
   # Document Title (H1)
   ## Major Section (H2)
   ### Subsection (H3)
   ### Another Subsection (H3)
   ## Another Major Section (H2)
   ```

4. **Balance section sizes**: Avoid heading siblings with vastly different lengths

## Formatting for RAG

### Summary Sections

Include summary sections for documents longer than 400 lines or after complex sections.

**Format**:
```markdown
## Summary

Key takeaways:
- Point one explaining main concept
- Point two covering important detail
- Point three noting critical relationship
```

### Terminology Precision

Use consistent, unambiguous terminology.

**Avoid pronoun ambiguity**:
```markdown
❌ Bad: The service calls the command. It validates the request.
✅ Good: The service calls the command. The command validates the request.
```

**Define acronyms on first use**:
```markdown
The Inversion of Control (IoC) container manages dependencies.
```

**Use consistent terms**: Pick one term and stick with it throughout
- ✅ Good: Always use "Request Context"
- ❌ Bad: Switching between "Request Context", "context object", "auth context"

### Context Anchoring

Begin each major section with context to help AI understand purpose and scope.

**Example**:
```markdown
## Authorization System

The authorization system validates user permissions before allowing access.
This section covers validation patterns, common scenarios, and error handling.

### Request Context Structure
...
```

### Code Block Annotations

Always include:
1. **Language identifier**: For syntax highlighting
2. **File path comment**: Where this code lives
3. **Explanatory comments**: What the code does and why
4. **Import statements**: Show dependencies

**Example**:
```typescript
// src/services/account-service.ts
import { injectable } from 'inversify';

/**
 * Account service demonstrating authorization validation.
 */
@injectable()
export class AccountService {
  async getAccount(accountId: string, context: RequestContext) {
    // Validate user has access to this account
    context.requireAccountAccess(accountId);

    return this.accountRepo.findById(accountId);
  }
}
```

## Explicit vs. Implicit Information

AI systems struggle with implicit information. Make relationships and dependencies explicit.

### Bad: Implicit Dependencies
```markdown
Services should use the standard pattern. Authorization is required.
```

**Problems**: "Standard pattern" is undefined, "required" doesn't explain how.

### Good: Explicit Dependencies
```markdown
Services follow the service-command pattern. Key requirements:

1. **Authorization**: Every service method must validate context before
   executing business logic using `context.requireAccountAccess()`.

2. **Error Handling**: Services catch command errors and transform them
   into appropriate HTTP responses.
```

## Optimizing for Search and Retrieval

### Question-Answer Format

Structure content to answer common questions.

**Example**:
```markdown
## How Do I Implement Authorization?

To implement authorization in a service:
1. Accept context as a parameter
2. Call authorization check before business logic
3. Handle authorization errors appropriately

## What Happens If Authorization Fails?

When authorization fails, the system throws an error with:
- User ID
- Resource ID
- Required permission
```

### Semantic Relationships

Explicitly state relationships between concepts.

**Example**:
```markdown
Request Context **is the foundation for** authorization. Authorization
validation **depends on** Request Context to determine permissions.

The relationship:
- Request Context **provides** hierarchical permission data
- Authorization validators **consume** Request Context
- Services **must create** Request Context before authorization
```

## Common Pitfalls

### Pitfall 1: Documents Too Long
**Problem**: 2000-line documents covering multiple topics
**Solution**: Split into focused 300-800 line documents

### Pitfall 2: Missing Context
**Problem**: Sections assume knowledge from previous sections
**Solution**: Each major section should be self-contained with brief context

### Pitfall 3: Inconsistent Terminology
**Problem**: Using "service", "provider", "manager" interchangeably
**Solution**: Define and use consistent terms throughout

### Pitfall 4: Implicit Dependencies
**Problem**: "Use the standard approach" without defining it
**Solution**: Name and link to specific patterns and guidelines

### Pitfall 5: No Lifecycle Separation
**Problem**: Active and planned features in same document
**Solution**: Use status field and separate directories for planning/archive

## Summary

Key principles for AI-optimized documentation:

- **Keep documents 300-800 lines** for optimal chunking
- **Include YAML frontmatter** with title, domain, keywords, related, last_updated, status
- **Use 3-level heading hierarchy** maximum (H1, H2, H3)
- **Add summary sections** for documents >400 lines
- **Avoid pronoun ambiguity** - repeat nouns explicitly
- **Define acronyms** on first use in each document
- **Annotate code blocks** with file paths and explanatory comments
- **Make relationships explicit** - state dependencies clearly
- **Separate lifecycle states** using status field
- **Structure for questions** - organize content around common queries
