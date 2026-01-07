# TypeScript Code Review

A comprehensive prompt for conducting thorough TypeScript code reviews focused on type safety, maintainability, and best practices.

## Purpose

Use this prompt when you need a detailed TypeScript code review that emphasizes strong typing, elimination of magic values, appropriate use of type assertions, and adherence to TypeScript best practices.

## The Prompt

```
Please conduct a comprehensive TypeScript code review focusing on the following areas:

## 1. Type Safety
- Identify any usage of `any` types and suggest specific, strongly-typed alternatives
- Check for implicit `any` types that should be explicitly defined
- Verify return types are explicitly declared on functions (don't rely on inference for public APIs)
- Look for missing type annotations on function parameters
- Identify places where union types or discriminated unions would be more appropriate
- Check for proper use of `unknown` vs `any` when type is truly unknown

## 2. Magic Strings & Numbers
- Flag all magic strings and numbers that should be constants or enums
- Verify that existing constants/types/enums are being used instead of literals
- Check for string literals that should be string literal types or enums
- Identify repeated values that should be extracted to named constants
- Look for hardcoded configuration values that should be externalized

## 3. Type Assertions & Casting
- Identify all uses of type assertions (`as`, `<Type>`, non-null assertions `!`)
- For each type assertion, analyze and document WHY it's necessary
- Suggest safer alternatives (type guards, type narrowing, validation)
- Flag dangerous assertions that could cause runtime errors
- Check if type predicates or user-defined type guards would be better
- Verify that assertions are truly necessary and not masking type definition issues

## 4. TypeScript Best Practices
- **Interfaces vs Types**: Ensure appropriate use (prefer interfaces for object shapes, types for unions/intersections)
- **Readonly Properties**: Suggest `readonly` for properties that shouldn't change
- **Const Assertions**: Look for places where `as const` would improve type safety
- **Strict Null Checks**: Verify proper handling of null/undefined cases
- **Enum Usage**: Check if enums are appropriate or if const objects with `as const` would be better
- **Generic Constraints**: Ensure generics have appropriate constraints
- **Utility Types**: Identify opportunities to use built-in utility types (Partial, Pick, Omit, etc.)
- **Index Signatures**: Verify proper use and consider if Record<K, V> is more appropriate

## 5. Code Quality & Maintainability
- Check for overly complex type definitions that could be simplified
- Identify missing JSDoc comments on public APIs
- Look for type definitions that should be shared/reused
- Verify consistent naming conventions for types and interfaces
- Check for proper error handling with typed errors
- Identify opportunities for branded types for primitive values with semantic meaning

## 6. Potential Issues
- Look for type coercion that could cause bugs
- Check for improper optional chaining or nullish coalescing usage
- Identify places where stricter types would catch potential bugs
- Flag any `@ts-ignore` or `@ts-expect-error` comments and verify they're justified

## Output Format

For each issue found, provide:
1. **Location**: File and line number
2. **Category**: Which review area it falls under
3. **Issue**: Clear description of the problem
4. **Impact**: Why this matters (type safety, maintainability, potential bugs)
5. **Recommendation**: Specific code suggestion or approach to fix it
6. **Priority**: High/Medium/Low based on potential impact

Also provide:
- Summary of overall type safety score
- List of positive patterns worth maintaining
- Prioritized list of recommended improvements
```

## When to Use

- Before merging pull requests with TypeScript changes
- During refactoring efforts to improve type safety
- When onboarding new team members to establish standards
- When upgrading TypeScript versions to leverage new features
- During architectural reviews of TypeScript codebases

## When NOT to Use

- For quick bug fixes that don't touch type definitions
- When reviewing non-TypeScript files
- During rapid prototyping phases where strict typing would slow iteration

## Usage Examples

### Full Codebase Review
```
[Paste the TypeScript code review prompt]

Please review the entire `src/` directory, focusing especially on the API layer and data models.
```

### Specific File Review
```
[Paste the TypeScript code review prompt]

Please review `src/services/user-service.ts` - I'm particularly concerned about the type assertions in the data transformation functions.
```

### Pull Request Review
```
[Paste the TypeScript code review prompt]

Please review the changes in PR #123. The diff shows modifications to our authentication types and I want to ensure we haven't weakened our type safety.
```

## Tips

- **Run TypeScript in Strict Mode First**: Enable all strict flags before review to catch obvious issues
- **Provide Context**: Share relevant type definitions and interfaces that the reviewed code depends on
- **Prioritize**: Not all issues need immediate fixing - focus on high-impact type safety concerns first
- **Document Justified Assertions**: When type assertions are truly necessary, add comments explaining why
- **Use Type Guards**: Prefer writing type guard functions over assertions for reusable type narrowing
- **Consider Branded Types**: For primitives with semantic meaning (UserId, Email, etc.), consider branded types
- **Leverage DefinitelyTyped**: Before creating custom types for libraries, check if @types packages exist
- **Test Your Types**: Consider using tools like `tsd` to write tests for complex type definitions

## Common Patterns to Watch For

### Magic Strings that Should Be Types
```typescript
// ❌ Bad
function setStatus(status: string) { }
setStatus("active") // typo-prone

// ✅ Good
type Status = "active" | "inactive" | "pending"
function setStatus(status: Status) { }
```

### Unnecessary Type Assertions
```typescript
// ❌ Bad
const data = JSON.parse(response) as UserData

// ✅ Good - with validation
function isUserData(obj: unknown): obj is UserData {
  return typeof obj === "object" && obj !== null && "id" in obj
}
const data = JSON.parse(response)
if (isUserData(data)) { /* use data safely */ }
```

### Magic Numbers that Should Be Constants
```typescript
// ❌ Bad
if (user.age >= 18) { }
setTimeout(callback, 5000)

// ✅ Good
const LEGAL_AGE = 18
const DEBOUNCE_MS = 5000
if (user.age >= LEGAL_AGE) { }
setTimeout(callback, DEBOUNCE_MS)
```

## Related Workflows

- Combine with general code review workflows
- Use alongside linting and formatting checks
- Integrate findings into pull request templates
- Reference when updating TypeScript configuration
