# Knowledge OS v3 - Design Q&A Session

**Date**: 2026-01-09  
**Status**: In Progress - Design Discovery  
**Purpose**: Capture design decisions and rationale for Knowledge OS v3 document

---

## Core Vision Established

**Fundamental Insight**: Modern software development is about **creating value FROM knowledge**, not just code from requirements.

- **Knowledge is the primary asset** (durable, versioned, maintained)
- **Code is an expression of knowledge** (somewhat ephemeral, regenerable)
- **AI supplements humans** in productizing knowledge into software
- **Outcome**: Institutional memory preservation, value lock-in, preventing corporate amnesia

---

## Design Decisions

### Q1: Knowledge Entry Points
**Question**: How should knowledge enter the system?

**Answer**: **All three methods equally important:**
- A) Explicit capture (users write structured docs)
- B) Conversation capture (AI listens, extracts, asks to add)
- C) Code observation (system notices logic, prompts documentation)

**Rationale**: Different knowledge enters through different channels. All must be supported from day one.

---

### Q2: Values Enforcement - How Prescriptive?
**Question**: When code violates an organizational value, what should happen?

**Answer**: **D) Configurable per value**

**Rationale**: Balance enforcement with trust in humans:
- Critical compliance values can **block** (require override authority)
- Operational guidelines can **warn** (flag but allow)
- Informational principles can **inform** (log for review)

**Example**:
```yaml
value:
  id: privacy-first
  enforcement: block  # Prevents merge without Privacy Officer approval
  
value:
  id: code-quality-standards
  enforcement: warn   # Flags but allows with justification
```

---

### Q3: Knowledge Graph Structure
**Question**: Should node types and relationships be prescribed or flexible?

**Answer**: **C) Hybrid - Prescribe core types, allow extension**

**Core prescribed types**:
- Values (organizational constitution)
- Domain rules (business logic, compliance)
- Capabilities (product features)
- Decisions (architectural choices, rationale)
- Technical patterns (implementation approaches)

**Core relationships**:
- REQUIRES (dependency)
- CONFLICTS_WITH (constraints)
- IMPLEMENTS (code-to-knowledge)
- OWNED_BY (responsibility)
- AFFECTS (impact)

**Extension**: Organizations can add custom node types and relationships as needs evolve.

**Rationale**: Prescribing essentials enables AI reasoning and faster bootstrap. Extension allows unique organizational needs.

---

### Q4: Knowledge Ownership & Maintenance
**Question**: Who owns and maintains knowledge over time?

**Answer**: **D) Hybrid ownership model**

**Options supported**:
- **Individual owners**: "Sarah owns Canadian tax rules"
- **Role-based**: "Tax Specialist role owns tax domain"
- **Team-based**: "Finance team owns billing knowledge"

**Rationale**: Different knowledge suits different ownership models:
- Specialized expertise → Individual
- Rotating responsibilities → Role
- Collaborative domains → Team

When Sarah (individual) leaves, her knowledge transfers to role/team seamlessly.

---

### Q5: Knowledge Freshness - External Monitoring
**Question**: How to keep external knowledge (tax laws, regulations) current?

**Answer**: **D) Hybrid approach**

**Two mechanisms**:
1. **Time-based review prompts**: "Tax rules last verified 6 months ago - time to review?"
2. **Automated monitoring** (optional): System checks external URLs, flags significant changes

**Rationale**: 
- Time-based prompts work for all knowledge (simple, reliable)
- Automated monitoring adds value for critical knowledge with stable sources
- Experts ultimately verify accuracy (AI assists, humans validate)

---

### Q6: First Knowledge to Capture
**Question**: When new org adopts Knowledge OS, what should they capture first?

**Answer**: **E) Let organizations choose based on pain points**

**Provide starting templates for**:
- A) Organizational values (cultural foundation)
- B) Domain knowledge (business rules, compliance)
- C) Architectural decisions (existing system rationale)
- D) Current work in progress (immediate value)

**Recommended path**: A → B → C (values first establishes foundation, but flexible)

**Rationale**: Different orgs have different urgent needs:
- Scaling startup → Capture values before culture dilutes
- Regulated industry → Domain/compliance rules critical
- Inherited codebase → Decision lineage prevents "why?" questions
- Active projects → Prove value immediately

---

### Q7: Code-Knowledge Relationship (CRITICAL INSIGHT)
**Question**: When knowledge changes, what happens to code?

**Answer**: **Knowledge changes trigger "Knowledge Pull Request" workflow**

**Key Insight**: "Multi-layered git for knowledge and productization of that knowledge"

**Workflow**:
```
Knowledge Change Proposed
    ↓
Requirements Analysis (What needs to change?)
    ↓
Impact Assessment (Which code/systems affected? Configuration? Code? Architecture?)
    ↓
Design Proposals (How should we implement?)
    ↓
Expert Reviews (Domain experts weigh in)
    ↓
Implementation Plan (Code changes mapped, or just DB update?)
    ↓
Human Approval & Execution
    ↓
Knowledge + Implementation merged together
    ↓
Historical Record (full traceability)
```

**Impact varies greatly**:
- **Simple**: Update database value (e.g., tax rate in config table)
- **Medium**: Code logic change (e.g., calculation algorithm)
- **Complex**: Architecture change (e.g., new service needed)

**Example - Tax Rate Change**:
```yaml
# PROPOSED CHANGE to /knowledge/domains/tax/canadian-gst.yaml
- rate: 0.05  # current
+ rate: 0.06  # proposed
rationale: "CRA announced GST increase effective 2026-04-01"
source: https://canada.ca/revenue-agency/gst-update

# System analyzes impact
IMPACT ANALYSIS:
Implementation type: database_config
Affected: config.tax_rates.CA_GST_RATE
Code references: src/billing/tax/gst_calculator.go:45 (reads from DB)
Tests affected: 12 tests need expectation updates
Required approvals: Tax Expert, Compliance Officer
Estimated effort: 30 minutes
Risk: LOW
```

**Critical**: AI SUPPLEMENTS this process (provides analysis), humans DECIDE (determine actual impact and approach).

**Unsolved Design Challenge**: How does system know knowledge↔implementation mapping?
- Explicit linking in knowledge files?
- Code annotations referencing knowledge?
- AI inference via semantic analysis?
- Manual architect review?
- All of the above?

**Likely answer**: Combination of approaches depending on situation. To be refined.

---

### Q8: Knowledge-to-Implementation Mapping
**Question**: How does system know WHERE and HOW knowledge is implemented?

**Answer**: **To be determined - likely D) Multiple approaches**

**Options being considered**:
- Explicit linking in knowledge files
- Code annotations
- AI inference
- Manual architect review

**Key Principle Established**: AI suggests, humans decide. The impact assessment requires human judgment - it's contextual and varies by situation.

**Rationale**: This is inherently complex. A tax rate might be:
- A database value (simple update)
- Hardcoded constant (code change)
- Calculated dynamically (algorithm change)
- Part of external API (integration change)

No one-size-fits-all solution. System should help humans figure this out, not try to be fully autonomous.

---

### Q9: Value Amendment Process
**Question**: How should organizational values change over time?

**Answer**: **B) Versioned with history, like a Pull Request**

**Process**:
1. Propose value change (with rationale)
2. Impact analysis (what code/decisions affected?)
3. Review process (who must approve?)
4. Merge with full history preserved
5. Transition period (how to handle conflicts during changeover?)

**Example**:
```
VALUE CHANGE PR #12
Date: 2026-03-15
Proposer: New CEO

Current:
  id: privacy-first
  principle: "Customer privacy over growth metrics"

Proposed:
  id: privacy-first
  principle: "Balanced privacy with sustainable monetization"
  
Rationale: "Company needs ad revenue to survive. Privacy remains important but can't be absolute."

Impact:
  - 23 decisions previously made under old value
  - 8 code constraints may need review
  - 2 vendor contracts rejected under old policy

Approvals Required: Board, Privacy Officer, CTO
Discussion: [link to conversation]
Vote: 4 approve, 1 dissent (with recorded reasoning)

Status: MERGED 2026-03-20
History: This is VERSION 2 of privacy-first value
  - v1: 2023-01-10 to 2026-03-20 (3.2 years)
  - v2: 2026-03-20 to present
```

**Outcome**: 
- Prevents silent value changes
- Full accountability and transparency
- Future employees understand the journey
- Can reference "under which value version was this decided?"

---

### Q10: AI's Role in Knowledge Entry
**Question**: When should AI surface potential knowledge from conversations/code?

**Answer**: **B and D together - Prompt after conversation + Batch suggestions**

**Two-tier approach**:

**Tier 1 - Real-time prompts (B)**: For high-confidence, critical knowledge
```
[After conversation pause]
💡 Knowledge Detected

I noticed a compliance rule mentioned:
"GDPR requires consent before storing location data"

Should I add this to /knowledge/compliance/gdpr.yaml?
[Yes] [Edit First] [No, ignore]
```

**Tier 2 - Batch review (D)**: For lower-confidence suggestions
```
📊 Weekly Knowledge Review

I extracted 8 potential knowledge items from this week's activities:

1. ✓ Business rule: Refund window is 30 days (HIGH confidence)
   Source: Product team discussion, mentioned 3 times
   
2. ? Technical decision: Using Redis for session storage (MEDIUM)
   Source: Code commit by @jane, no doc explanation
   
3. ? Value: "Ship fast, iterate later" (LOW confidence)
   Source: CEO comment in standup
   
[Review All] [Auto-approve HIGH] [Dismiss]
```

**Rationale**: 
- Don't interrupt flow with low-value suggestions
- Do capture critical knowledge before it's forgotten
- Batch review reduces notification fatigue
- Confidence scoring helps prioritize human attention

---

### Q11: Knowledge Graph Visualization
**Question**: How should users explore the knowledge graph?

**Answer**: **E) Multiple interfaces - meet users where they are**

**Different personas, different needs**:

**Developers**:
- Filesystem interface: `ls /knowledge/domains/`, `cat /knowledge/values/privacy.yaml`
- CLI commands: `knowledge find --owner=sarah --type=tax`
- Git-like operations: `knowledge diff v1..v2`, `knowledge blame privacy-first`

**Product Owners**:
- Visual graph interface (interactive, zoom, filter)
- Capability maps (business value view)
- Value dashboard (ROI, impact)

**Domain Experts**:
- Filtered views of their domains
- Review queue (pending approvals)
- Refinement dashboard (AI suggestions in their area)

**Everyone**:
- AI-assisted queries: "What compliance rules affect user data collection?"
- Search: Full-text + semantic search
- History: Time-travel through knowledge evolution

**Rationale**: 
- One interface won't fit all personas
- Developers comfortable with CLI/files
- Executives need visual dashboards
- Everyone benefits from AI assistance

---

---

### Q12: Knowledge Versioning Granularity
**Question**: What granularity should knowledge versioning happen at?

**Answer**: **Atomic-level versioning with document-level organization**

**Key Decision**: Adopted **Atomic Design Framework** (see [`project-os-atoms.md`](project-os-atoms.md))

**Versioning Strategy**:
- **Atoms**: Individual atomic units are versioned (TextAtom, NumberAtom, etc.)
- **Molecules**: Groups of atoms versioned together
- **Organisms**: Complete knowledge units versioned as cohesive unit
- **Documents**: Collection of organisms versioned with git semantics
- **Graph Snapshots**: Graph state captured at document commit boundaries

**Atomic-level versioning example**:
```bash
# Show version history for specific atom
knowledge blame --atom=NumberAtom --id=rate /knowledge/tax/gst.yaml

> v3.2: Changed from 0.05 to 0.06 (by sarah, 2026-01-09)
> v2.1: Changed from 0.04 to 0.05 (by mike, 2025-06-15)
> v1.0: Initial value 0.04 (by sarah, 2024-01-10)
```

**Graph versioning**:
- Each commit creates a graph snapshot
- Temporal queries: "Show knowledge state as of 2026-01-01"
- Lineage tracking: Full provenance of atom/molecule/organism evolution
- Rollback: Can revert specific atoms, molecules, or entire organisms

**Rationale**: 
- Atomic granularity enables precise diffing (what exactly changed?)
- Document-level organization makes sense to humans (files are natural mental model)
- Graph snapshots enable complex queries across entire knowledge base at a point in time
- Git semantics are familiar and battle-tested

**See detailed atomic framework**: [`project-os-atoms.md`](project-os-atoms.md)

---

### Q21: Knowledge Graph Technology Choice

**Question**: What technology should power the knowledge graph for the atomic framework (Atoms → Molecules → Organisms → Documents), considering multi-tenancy and scalability?

**Answer**: **Hybrid approach with abstraction layer**

**Technology Stack**:
- **PostgreSQL + pgvector + AGE** for atoms, documents, vector search, and transactional operations
- **Neo4j** for relationship traversals, dependency graphs, and complex impact analysis
- **Abstraction layer** that supports both, with ability to start with PostgreSQL-only and add Neo4j later

**Rationale for Hybrid**:
1. **Multi-tenant isolation**: PostgreSQL row-level security + separate databases per large client
2. **Query optimization**:
   - PostgreSQL excels at: atomic CRUD, full-text search, vector embeddings, complex transactions
   - Neo4j excels at: deep relationship traversals (`knowledge trace --via=REQUIRES`), graph algorithms, impact analysis
3. **Scalability path**:
   - Phase 1: PostgreSQL only (simpler ops)
   - Phase 2: Add Neo4j for organizations with complex dependency graphs
   - Phase 3: Optional Neo4j per high-usage tenant
4. **Cost efficiency**: Small orgs use PostgreSQL only; large orgs can enable Neo4j

**Schema Mapping**:
```yaml
# PostgreSQL tables
atoms(id, tenant_id, type, value, created_at, version)
molecules(id, tenant_id, type, created_at)
atom_molecules(atom_id, molecule_id)
documents(id, tenant_id, path, template_type)

# Neo4j graph
(:Atom {id, type, value})-[:PART_OF]->(:Molecule {id, type})
(:Molecule)-[:REQUIRES]->(:Molecule)
(:Document)-[:CONTAINS]->(:Molecule)
```

**Multi-tenancy Architecture**:
1. **Database-per-tenant**: Large enterprises (>10k nodes)
2. **Schema-per-tenant**: Medium organizations (1k-10k nodes)  
3. **Row-level isolation**: Small teams (<1k nodes)
4. **Hybrid**: Mix based on tier and requirements

**Migration Path**:
1. Start with PostgreSQL abstraction (current SQLite → PostgreSQL)
2. Add Neo4j driver behind same interface
3. Use feature flags to enable Neo4j per tenant
4. Sync data via change data capture or dual-writes

**Performance Targets**:
- Relationship traversals (5+ hops): <100ms (Neo4j)
- Vector similarity search: <50ms (pgvector)
- Atomic updates: <10ms (PostgreSQL)
- Graph impact analysis: <200ms (Neo4j)

**Implementation Priority**:
1. **Abstract storage interface** supporting both SQL and graph operations
2. **PostgreSQL implementation** first (covers 80% of use cases)
3. **Neo4j implementation** for high-tier tenants
4. **Query router** that directs queries to appropriate backend

**Decision**: Hybrid approach provides maximum flexibility for multi-tenant scaling while allowing cost-optimized deployment per organization size.

---

## Open Questions (Still to Discuss)

13. **TBD** - Conflict resolution when multiple experts disagree (deferred for future iteration)
14. Knowledge privacy/security (who can see what?) - Hybrid: Directory-based + ABAC
15. **TBD** - Multi-tenancy for consultancies serving multiple clients (deferred for future iteration, leaning toward B or C)
16. Integration with existing tools (Notion, Confluence, Jira)
17. Bootstrapping experience for organizations with zero documentation
18. Measuring knowledge graph health (coverage, freshness, etc.)
19. Knowledge deprecation vs. archival
20. Cost model / pricing approach

---

## Next Steps

Once Q&A complete:
1. Draft Knowledge OS v3 document incorporating all decisions
2. Create visual diagrams (knowledge graph structure, workflows)
3. Develop concrete examples for each knowledge type
4. Write "Day in" Life" scenarios from knowledge perspective
5. Outline implementation roadmap

## Related Documents Created

- **[`project-os-atoms.md`](project-os-atoms.md)** - Atomic Design Framework for Knowledge
  - Atoms → Molecules → Organisms → Templates → Documents
  - Detailed atomic types, molecules, organisms, templates
  - Knowledge graph mapping from atomic primitives
  - AI assistance patterns using atomic framework

---

**Status**: 14/21 questions answered (2 deferred)  
**Next**: Continue design Q&A to completion
