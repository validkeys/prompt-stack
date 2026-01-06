# React Code Review

A comprehensive prompt for conducting thorough React code reviews focused on idiomatic patterns, performance optimization, component architecture, and maintainability.

## Purpose

Use this prompt when you need a detailed React code review that emphasizes proper hook usage, component composition, performance optimization, elimination of code duplication, and adherence to React best practices.

## The Prompt

```
Please conduct a comprehensive React code review focusing on the following areas:

## 1. Hook Usage & Rules of Hooks
- **Conditional Hooks**: Flag ALL instances of hooks used conditionally or in loops (STRICTLY NOT ALLOWED)
- Verify hooks are only called at the top level of components/custom hooks
- Check that hooks are only called from React functions (not regular JS functions)
- Ensure custom hooks follow the `use` naming convention
- Verify proper dependency arrays in useEffect, useMemo, useCallback
- Look for missing or overly broad dependencies that could cause bugs or performance issues
- Check for unnecessary dependencies that cause excessive re-renders

## 2. Component Architecture & Composition
- **God Components**: Identify large components (>200-300 lines) that should be broken down
- Suggest decomposition strategies: extract sub-components, use composition, or implement context/Jotai
- Check for components with too many responsibilities (violating Single Responsibility Principle)
- Verify proper use of component composition over prop drilling
- Look for opportunities to extract reusable components
- Check if components are appropriately sized and focused
- Identify components that should be split by feature/concern

## 3. State Management
- Verify appropriate choice of state management (useState vs useReducer vs Context vs Jotai)
- Check for prop drilling that should use Context or Jotai atoms
- Look for derived state that should be computed values instead
- Identify state that's duplicated or could be lifted/lowered
- Verify state is colocated close to where it's used
- Check for unnecessary state that could be refs or local variables
- Look for server state being managed in client state (should use React Query/SWR)

## 4. Performance Optimization
- **useCallback**: Verify it's used ONLY when truly needed (passed to memoized child components, dependency for other hooks)
- **useMemo**: Check it's used for expensive computations, not simple operations
- Flag OVER-optimization: unnecessary useCallback/useMemo that adds complexity without benefit
- Identify missing React.memo on expensive child components that receive stable props
- Check for inline function definitions passed to memoized components
- Look for expensive operations in render that should be memoized
- Verify proper key usage in lists (not using index as key when list is dynamic)
- Check for unnecessary re-renders caused by object/array literals in JSX

## 5. Type & Constant Redefinitions
- Identify duplicate type definitions that should be shared
- Flag redefined constants across components
- Check for inline types that should be extracted and reused
- Verify shared types are in appropriate locations (types/ or co-located)
- Look for magic strings/numbers that should be constants (also relates to TypeScript review)
- Identify repeated prop type definitions that could use shared interfaces

## 6. Code Duplication & Logic Reuse
- Identify duplicate logic that should be extracted to custom hooks
- Flag similar components that could share a base component or use composition
- Look for repeated UI patterns that should be extracted as components
- Check for duplicate data fetching/transformation logic
- Identify repeated event handlers that could be abstracted
- Look for duplicate validation logic that should be shared

## 7. Clarity & Intent
- Verify descriptive component and hook names that express intent
- Check for clear prop names that indicate purpose
- Look for complex boolean logic that should be extracted to named variables
- Verify event handlers have descriptive names (handleSubmit, not onClick)
- Check for magic numbers/strings without explanation
- Ensure complex JSX is broken down or has explanatory comments
- Verify ternaries are readable (nested ternaries should be avoided)

## 8. Documentation
- **JSDoc Comments**: Check that complex components have JSDoc describing purpose, props, and usage
- Verify custom hooks have JSDoc explaining purpose, parameters, and return values
- Look for complex logic that needs explanatory comments
- Check that prop interfaces/types have descriptions for non-obvious props
- Verify examples are provided for complex component APIs
- Ensure public components have usage documentation

## 9. React Patterns & Idioms
- Check proper use of children prop for composition
- Verify appropriate use of render props vs hooks
- Look for opportunities to use compound components pattern
- Check proper use of forwardRef and useImperativeHandle when needed
- Verify proper error boundary implementation for critical components
- Check for appropriate use of Suspense and lazy loading
- Look for proper portal usage for modals/overlays

## 10. Effects & Side Effects
- Verify effects don't have missing cleanup functions
- Check that effects are focused and don't do too much
- Look for effects that should be event handlers instead
- Verify proper async handling in effects
- Check for race conditions in effects
- Identify effects that run too frequently

## 11. Common Anti-Patterns
- Flag direct DOM manipulation (should use refs)
- Identify state initialization from props without synchronization
- Check for inline object/array creation in dependency arrays
- Look for derived state that's not kept in sync
- Flag useState for values that don't trigger re-renders (should be useRef)
- Identify functions that close over stale state/props
- Check for missing exhaustive dependency warnings being ignored

## 12. Accessibility & Best Practices
- Verify semantic HTML usage
- Check for proper ARIA attributes where needed
- Look for missing keyboard navigation support
- Verify proper form labeling
- Check for focus management in modals/dialogs

## Output Format

For each issue found, provide:
1. **Location**: File, component name, and approximate line number
2. **Category**: Which review area it falls under
3. **Issue**: Clear description of the problem
4. **Impact**: Why this matters (performance, bugs, maintainability, UX)
5. **Recommendation**: Specific refactoring approach or code suggestion
6. **Priority**: Critical/High/Medium/Low based on potential impact

Also provide:
- Summary of overall code quality score
- List of positive patterns worth maintaining
- Prioritized refactoring roadmap
- Suggested component extraction opportunities
```

## When to Use

- Before merging pull requests with React component changes
- During refactoring efforts to improve component architecture
- When performance issues are reported
- When onboarding new team members to establish standards
- During code reviews of feature branches
- When planning component library architecture

## When NOT to Use

- For quick bug fixes that don't touch component logic
- When reviewing non-React files (utilities, types, etc.)
- During rapid prototyping where perfect architecture isn't priority
- For tiny components (<50 lines) with simple logic

## Usage Examples

### Full Component Review
```
[Paste the React code review prompt]

Please review all components in `src/features/dashboard/` - I'm concerned about performance and component size.
```

### Performance-Focused Review
```
[Paste the React code review prompt]

Please review `src/components/DataTable.tsx` with focus on sections 4 (Performance) and 10 (Effects). The component is re-rendering excessively.
```

### Architecture Review
```
[Paste the React code review prompt]

Please review `src/pages/UserProfile.tsx` focusing on sections 2 (Architecture) and 3 (State Management). This component has grown to 800 lines and needs refactoring.
```

### Hook Usage Audit
```
[Paste the React code review prompt]

Please review the custom hooks in `src/hooks/` focusing on sections 1 (Hook Rules) and 6 (Code Duplication). Ensure proper hook patterns and identify any duplicated logic.
```

## Tips

- **Enable ESLint React Hooks Plugin**: Run `eslint-plugin-react-hooks` before manual review to catch hook rule violations
- **Profile Before Optimizing**: Use React DevTools Profiler to identify actual performance issues before adding memoization
- **Context vs Jotai**: Use Context for theme/auth (infrequent changes), Jotai for complex state management
- **Measure Component Size**: Components over 200-300 lines are candidates for decomposition
- **Custom Hook Extraction**: If you copy-paste hook logic twice, extract a custom hook
- **Colocate State**: Keep state as close as possible to where it's used
- **Avoid Premature Optimization**: Don't add useCallback/useMemo without measuring first
- **Document Complex Logic**: Add comments explaining WHY, not WHAT
- **Use TypeScript**: Pair this review with TypeScript code review for comprehensive coverage

## Common Patterns to Watch For

### Conditional Hook Usage (NEVER ALLOWED)
```typescript
// ❌ CRITICAL: Hooks in conditions
function Component({ isEnabled }) {
  if (isEnabled) {
    const [value, setValue] = useState(0) // BREAKS RULES OF HOOKS
  }
}

// ✅ Good: Conditional logic inside hook
function Component({ isEnabled }) {
  const [value, setValue] = useState(0)
  const effectiveValue = isEnabled ? value : 0
}
```

### God Component Needing Decomposition
```typescript
// ❌ Bad: 500-line component doing everything
function UserDashboard() {
  // 100 lines of state
  // 200 lines of effects and handlers
  // 200 lines of JSX
}

// ✅ Good: Composed smaller components
function UserDashboard() {
  return (
    <>
      <UserProfile />
      <UserStats />
      <UserActivity />
      <UserSettings />
    </>
  )
}
```

### Unnecessary Optimization
```typescript
// ❌ Bad: Over-optimization
const Component = () => {
  const simpleValue = useMemo(() => props.x + props.y, [props.x, props.y]) // Overkill
  const handleClick = useCallback(() => {}, []) // Not passed to memoized child
}

// ✅ Good: Optimize only when needed
const Component = () => {
  const simpleValue = props.x + props.y // Just compute it
  const handleClick = () => {} // Fine if parent rarely re-renders
}
```

### Duplicate Logic to Extract
```typescript
// ❌ Bad: Duplicated fetch logic
function ComponentA() {
  const [data, setData] = useState(null)
  useEffect(() => { fetch('/api/data').then(setData) }, [])
}
function ComponentB() {
  const [data, setData] = useState(null)
  useEffect(() => { fetch('/api/data').then(setData) }, [])
}

// ✅ Good: Custom hook
function useApiData() {
  const [data, setData] = useState(null)
  useEffect(() => { fetch('/api/data').then(setData) }, [])
  return data
}
```

### Type Redefinition to Consolidate
```typescript
// ❌ Bad: Duplicate type definitions
// ComponentA.tsx
interface User { id: string; name: string }

// ComponentB.tsx
interface User { id: string; name: string }

// ✅ Good: Shared types
// types/user.ts
export interface User { id: string; name: string }
```

### State That Should Be Context/Jotai
```typescript
// ❌ Bad: Prop drilling through 5 levels
<Parent theme={theme}>
  <Child theme={theme}>
    <GrandChild theme={theme}>
      <GreatGrandChild theme={theme} />

// ✅ Good: Context for deeply shared values
const ThemeContext = createContext()
// Or Jotai for complex state
const themeAtom = atom('light')
```

## Related Workflows

- Combine with TypeScript code review for full coverage
- Use alongside component testing reviews
- Integrate with performance profiling workflows
- Reference when establishing component architecture patterns
- Pair with accessibility audits

## State Management Decision Tree

**Use useState when:**
- Simple, local component state
- Boolean flags, form inputs, UI state
- State doesn't need to be shared

**Use useReducer when:**
- Complex state logic with multiple sub-values
- Next state depends on previous state
- State transitions follow clear patterns

**Use Context when:**
- Data needed by many components at different nesting levels
- Infrequently changing values (theme, auth, locale)
- Want to avoid prop drilling

**Use Jotai when:**
- Complex global state management
- Need derived state and computed values
- Want atomic state updates
- Need state persistence or synchronization
- Building complex state machines
