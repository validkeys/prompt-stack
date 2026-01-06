# Build Implementation Plan Workflow

A structured workflow for gathering requirements, clarifying ambiguities through iterative Q&A, and creating a comprehensive implementation plan with code samples.

## Purpose

Use this workflow when starting a new feature or project that requires detailed requirements gathering and planning. This workflow enforces a disciplined approach: document requirements first, clarify through questions, then plan the implementation with AI-optimized documentation and organized code samples.

## The Workflow

```
I need to build an implementation plan using a structured requirements-gathering process:

**Phase 1: Initial Requirements**
Create a file called `requirements.md` in the project root (or appropriate folder).
I will provide my initial requirements in plain text.

**Phase 2: Requirements Clarification**
Read the requirements document, then interview me with clarifying questions:
- Ask ONE question at a time
- After each answer, update `requirements.md` with:
  - The question you asked (as a heading or bold text)
  - My answer (as the content under that question)
- Continue until you have no more questions
- When I say "I'm done" or "no more questions", proceed to Phase 3

**Phase 3: Requirements Review**
Review the updated requirements document for:
- Gaps in understanding
- Missing edge cases
- Unclear acceptance criteria
- Ambiguous terminology
Ask any final clarifying questions if needed.

**Phase 4: Implementation Plan**
Create `implementation-plan.md` following these rules:
- Follow the AI-optimized documentation guidelines (proper headings, metadata, summaries)
- Include YAML frontmatter with: title, domain, keywords, related, last_updated, status
- Use maximum 3-level heading hierarchy (H1, H2, H3)
- NO code samples in this document - reference them instead
- Include sections for:
  - Overview and objectives
  - Architecture and approach
  - Component breakdown
  - Data models and relationships
  - API endpoints or interfaces
  - State management approach
  - Error handling strategy
  - Testing strategy
  - Security considerations
  - Performance considerations
- Reference code samples using format: "See code-samples/001-example.sample.tsx"

**Phase 5: Code Samples**
Create a `code-samples/` subfolder containing:
- One file per code sample
- File naming: `NNN-descriptive-name.sample.ext`
  - NNN = auto-incrementing 3-digit prefix (001, 002, 003, etc.)
  - descriptive-name = kebab-case description
  - sample = literal text to indicate this is a sample
  - ext = appropriate file extension (.tsx, .ts, .py, etc.)
- Each sample should:
  - Include file path comment showing where it would live in the real project
  - Include explanatory comments
  - Be complete and runnable (or clearly marked as partial if showing a snippet)
  - Include imports and type definitions

**Phase 6: Tracking Document**
Create `tracking.md` with:
- Milestone-based structure using H2 headings
- Task lists under each milestone using markdown checkboxes
- Format:
  ```markdown
  ## Milestone 1: Foundation
  - [ ] Task description
  - [ ] Another task

  ## Milestone 2: Core Features
  - [ ] Task description
  - [ ] Another task
  ```
- Tasks should be specific and actionable
- Order milestones logically (foundation → core → polish)

**Final Deliverables:**
- `requirements.md` - Fully clarified requirements with Q&A
- `implementation-plan.md` - AI-optimized implementation plan
- `code-samples/` - Folder with numbered code sample files
- `tracking.md` - Milestone-based task tracking

Wait for my initial requirements to begin Phase 1.
```

## Usage Example

Start by asking your AI assistant to follow this workflow:

```
Let's use the build-implementation-plan workflow. I'm ready to provide my initial requirements.
```

Then provide your requirements in plain text. The AI will create the `requirements.md` file and begin asking clarifying questions. Answer one at a time until you're satisfied, then proceed through the phases.

## Output Format

The workflow creates a complete planning package:

```
project-or-feature-folder/
├── requirements.md           # Clarified requirements with Q&A
├── implementation-plan.md    # AI-optimized plan document
├── tracking.md               # Milestone-based task list
└── code-samples/
    ├── 001-component-example.sample.tsx
    ├── 002-api-endpoint.sample.ts
    ├── 003-data-model.sample.ts
    └── ...
```

## Tips

- **Start broad, then narrow**: Provide high-level requirements first, let the AI ask about details
- **Answer honestly**: If you don't know something, say so - it helps identify what needs research
- **Review before proceeding**: After Q&A, read the updated requirements.md before moving to planning
- **Reference samples liberally**: In the implementation plan, reference specific code samples by number and name
- **Keep samples focused**: Each code sample should demonstrate one concept or pattern
- **Use realistic names**: Code samples should use realistic variable/function names from your domain
- **Iterate if needed**: After seeing the plan, you can ask the AI to update requirements and regenerate
- **Copy tracking.md**: The tracking document can be copied into your project management tool or GitHub issues
- **Version control**: Commit all planning documents - they're valuable context for future changes
- **Update as you learn**: As implementation reveals new requirements, update the planning documents
