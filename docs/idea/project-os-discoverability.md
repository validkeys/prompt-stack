# Project OS: Discoverability in Composable Knowledge Systems

**Date**: 2026-01-09  
**Status**: Design Document - Discovery Patterns  
**Related**: [`project-os-atoms.md`](project-os-atoms.md) (Atomic Knowledge Design), [`product-os-v3-qa.md`](product-os-v3-qa.md) (Design Q&A)

---

## The Fundamental Problem: Composability ↔ Discovery Tension

**Core Insight**: The dream of composability dies at the altar of discovery. You can have infinite reusable pieces, but if no one can *find* them, they might as well not exist.

### Traditional Composability Failure Pattern
1. Team builds reusable components/libraries/knowledge
2. Initial few components are easy to remember and find
3. System grows to hundreds/thousands of pieces
4. Discovery becomes overwhelming ("I know we have this somewhere...")
5. Teams duplicate instead of reuse ("Easier to rebuild than find")
6. Composability benefits vanish despite the architecture

### Knowledge OS Discovery Philosophy
**Discovery isn't one tool—it's a stack** that:
- Translates human intent into structured queries
- Follows semantic relationships between knowledge units
- Presents results in contextually appropriate ways for each persona
- Gets *better* with scale (unlike traditional search)

---

## The Discovery Stack: Eight Complementary Layers

### Layer 1: Atomic Typing → Structured Querying
Every piece of knowledge has a known type, enabling precise queries beyond keywords.

```bash
# Find all tax rates above 10%
knowledge query --type=NumberAtom --domain=tax --value-gt=0.1

# Find all ExceptionMolecules in compliance domain  
knowledge query --molecule=ExceptionMolecule --domain=compliance

# Find all knowledge owned by Sarah
knowledge find --owner=sarah

# Find knowledge with specific metadata
knowledge query --metadata priority=P0 --status=active
```

**Key Benefit**: Query by semantic type (`TaxRuleOrganism`) not just keywords ("tax").

### Layer 2: Relationship Tracing → Impact Discovery
Every relationship (`REQUIRES`, `IMPLEMENTS`, `CONFLICTS_WITH`) is a discoverable path through the knowledge graph.

```bash
# What depends on Canadian GST rules? (forward tracing)
knowledge trace --from=/knowledge/domains/tax/canadian-gst --via=REQUIRES

# What code implements privacy-first value? (knowledge→code)
knowledge trace --from=/knowledge/values/privacy-first --via=IMPLEMENTS

# What knowledge constrains this feature? (code→knowledge)
knowledge trace --to=/capabilities/multi-currency-billing --via=CONSTRAINED_BY

# Show full dependency chain
knowledge impact --change=/knowledge/domains/tax/canadian-gst.yaml
```

**Key Benefit**: Discover knowledge through its connections, not just its content.

### Layer 3: Template-Based Browsing → Schema-Aware Discovery
Templates define what knowledge *should* exist. Discover by "kind of thing" not just "thing with these words."

```bash
# Show all tax rule documents (instances of TaxRuleTemplate)
knowledge browse --template=TaxRuleTemplate

# What capabilities do we have? (FeatureProposalTemplate instances)
knowledge browse --template=FeatureProposalTemplate

# Show all architectural decisions (DecisionTemplate instances)
knowledge browse --template=DecisionTemplate

# Find incomplete documents (missing required atoms/molecules)
knowledge validate --template=TaxRuleTemplate --show-incomplete
```

**Key Benefit**: Discover by category/schema, enabling "What kinds of things do we have?" exploration.

### Layer 4: AI-Assisted Semantic Discovery
Natural language queries that understand context, domain language, and relationships.

```sql
"Show me compliance rules about user data collection"
→ Finds: GDPR consent rules, CCPA requirements, internal privacy policies
→ Shows: Relationships between them (GDPR → CCPA → internal)
→ Suggests: Related capabilities (data deletion, consent management)

"What affects billing calculations for Canadian customers?"
→ Traces: GST rules, PST rules, tax exemption molecules, currency handling
→ Shows: Implementation files, test suites, dependent features
→ Warns: Stale knowledge (tax rates needing review)

"Who knows about tax compliance?"
→ Finds: Domain experts (Sarah), decision owners, recent contributors
→ Shows: Their knowledge contributions, review history
→ Suggests: Collaboration opportunities
```

**Key Benefit**: Bridge human language to structured knowledge without requiring query syntax expertise.

### Layer 5: Persona-Specific Discovery Interfaces
Different roles need different discovery tools optimized for their mental models.

#### Domain Expert (Tax Specialist, Compliance Officer)
- **Interface**: Document-focused search with domain filtering
- **Discovery pattern**: "Show me all tax knowledge" → gets documents, organisms, validation status
- **Tools**: Domain dashboard, review queue, validation reports
- **Example**: `knowledge domain tax --show-all --validate`

#### Developer (Software Engineer)
- **Interface**: Code↔knowledge tracing with impact analysis
- **Discovery pattern**: "What knowledge affects this file?" → gets related organisms, constraints
- **Tools**: IDE integration, CLI for code tracing, impact preview
- **Example**: `knowledge implements src/billing/tax/gst_calculator.go`

#### Product Owner / Business Analyst
- **Interface**: Capability mapping with business value focus
- **Discovery pattern**: "What compliance capabilities do we have?" → gets feature organisms, ROI estimates
- **Tools**: Capability dashboard, value visualization, dependency maps
- **Example**: `knowledge capabilities --domain=compliance --roi-gte=10000`

#### Compliance Officer / Security Specialist
- **Interface**: Constraint discovery with enforcement focus
- **Discovery pattern**: "What values enforce 'block' level?" → gets blocking constraints, violation history
- **Tools**: Constraint dashboard, violation alerts, approval workflows
- **Example**: `knowledge constraints --enforcement=block --domain=privacy`

#### AI Assistant
- **Interface**: Whole-graph semantic understanding
- **Discovery pattern**: "What's missing from our tax knowledge?" → finds gaps, suggests relationships
- **Tools**: Gap analysis, relationship suggestion, confidence scoring
- **Example**: (Internal system use, not user-facing)

**Key Benefit**: Right tool for the right job—no single interface tries to serve all personas.

### Layer 6: Graph Visualization → Spatial Discovery
Visual exploration of the knowledge graph for pattern recognition and relationship understanding.

```bash
# Interactive graph browser
knowledge graph --visual --domain=tax

# Zoom from high-level to detailed view
knowledge graph --organism=WithholdingTaxOrganism --show-atoms

# Filter by relationships
knowledge graph --filter="REQUIRES|IMPLEMENTS" --from=/knowledge/values/privacy-first

# Time-travel through knowledge evolution
knowledge graph --as-of="2025-06-01" --domain=billing
```

**Visualization Features**:
- **Zoom levels**: Domain → Organism → Molecule → Atom
- **Relationship highlighting**: Color-coded by type (REQUIRES=blue, CONFLICTS=red)
- **Cluster detection**: Automatically groups related knowledge
- **Temporal views**: See how knowledge evolved over time
- **Impact preview**: Visualize change propagation before committing

**Key Benefit**: Sometimes you need to *see* the relationships to understand what exists.

### Layer 7: Ownership & Responsibility Discovery
Discover knowledge through organizational structure, not just content.

```bash
# Who owns tax domain knowledge?
knowledge owners --domain=tax

# What knowledge is Sarah responsible for?
knowledge responsibilities --owner=sarah

# Show knowledge ownership by role (not individual)
knowledge owners --role=TaxSpecialist

# Find knowledge with unclear ownership
knowledge orphaned --domain=compliance

# Transfer discovery (when people leave/change roles)
knowledge transfer --from=sarah --to=mike --domain=tax
```

**Key Benefit**: Align discovery with organizational reality—people need to find what they're responsible for.

### Layer 8: Freshness & Confidence-Based Discovery
Discover knowledge that needs attention, not just knowledge that exists.

```bash
# Find knowledge that hasn't been reviewed in 6 months
knowledge stale --older-than=6months

# Find AI-extracted knowledge with low confidence
knowledge low-confidence --threshold=0.7

# Show knowledge with external dependencies needing verification
knowledge external --status=needs-verification

# Find values that haven't been applied to recent decisions
knowledge unused-values --since="2025-01-01"
```

**Key Benefit**: Proactive discovery—surfacing what needs human attention before it becomes a problem.

---

## Discovery Workflow Examples

### Example 1: New Feature Planning
**Scenario**: Product Owner wants to add EU VAT compliance.

```bash
# 1. Discover existing tax knowledge
knowledge browse --template=TaxRuleTemplate --region=*
→ Finds: Canadian GST, US sales tax, Australian GST

# 2. Find experts in tax domain
knowledge owners --domain=tax
→ Shows: Sarah (Tax Specialist), Mike (Compliance)

# 3. Discover related compliance capabilities
knowledge trace --from=/knowledge/compliance/gdpr --via=REQUIRES
→ Shows: Data protection, consent management, reporting

# 4. AI-assisted gap analysis
"what do we need for EU VAT compliance?"
→ Identifies: VAT rates per country, exemption rules, reporting requirements
→ Finds gaps: Missing EU country rates, no VAT reporting organisms
→ Suggests: Start with German VAT (largest EU market)

# 5. Impact discovery for implementation
knowledge impact --template=VATRuleTemplate --estimate
→ Shows: Billing system changes, reporting updates, test coverage needed
```

### Example 2: Code Change Impact Analysis
**Scenario**: Developer needs to update tax calculation logic.

```bash
# 1. Find knowledge affecting current file
knowledge implements src/billing/tax/calculator.py
→ Returns: CanadianGSTOrganism, WithholdingTaxOrganism, SalesTaxOrganism

# 2. Trace dependencies
knowledge trace --from=CanadianGSTOrganism --via=REQUIRES
→ Shows: Multi-currency billing, invoice generation, compliance reporting

# 3. Discover affected tests
knowledge test-coverage --organism=CanadianGSTOrganism
→ Shows: 12 test files, 47 test cases

# 4. Find decision history
knowledge history --organism=CanadianGSTOrganism --limit=5
→ Shows: Rate changes, exception additions, implementation updates

# 5. Validate change completeness
knowledge validate --change="rate: 0.05 → 0.06" --organism=CanadianGSTOrganism
→ Checks: All references updated, tests adjusted, documentation current
```

### Example 3: Compliance Audit Preparation
**Scenario**: Compliance officer needs audit evidence.

```bash
# 1. Find all blocking constraints
knowledge constraints --enforcement=block --domain=privacy
→ Shows: No-PII-sharing, consent-required, data-retention-limits

# 2. Discover implementation evidence
knowledge evidence --constraint=no-pii-sharing --since="2025-01-01"
→ Shows: Code reviews approving compliance, test results, deployment checks

# 3. Find decision rationale
knowledge rationale --decision="use-encryption-for-pii" --full
→ Shows: Alternatives considered, security team approval, risk assessment

# 4. Identify knowledge gaps
knowledge gaps --domain=compliance --standard=GDPR
→ Shows: Missing data-portability rules, incomplete consent documentation

# 5. Generate audit report
knowledge audit-report --domain=privacy --period="2025-Q4"
→ Produces: Compliance coverage, violation history, improvement recommendations
```

---

## The Composability ↔ Discovery Feedback Loop

### Positive Reinforcement Cycle
```
More composition → More structured knowledge → Better discovery
    ↓                                    ↑
More reuse ←─── Better discovery ←───────┘
```

1. **Team composes knowledge** using atomic framework (atoms → molecules → organisms)
2. **System captures structure** (types, relationships, metadata)
3. **Discovery tools leverage structure** for precise queries
4. **Team finds existing knowledge easily** (reduces duplication)
5. **Team reuses/composes more** (reinforcing the cycle)
6. **System gets richer structure** (improving discovery further)

### Unlike Traditional Systems
| Traditional Systems | Knowledge OS |
|-------------------|--------------|
| More pieces → harder to find | More pieces → better discovery (through structure) |
| Search degrades with scale | Discovery improves with scale (more relationships) |
| Duplication common ("can't find") | Reuse common ("easy to find") |
| Discovery = keyword search | Discovery = structured query + relationship tracing |

---

## Implementation Architecture

### Discovery Engine Components
1. **Query Parser**: Translates natural language + structured queries
2. **Graph Traversal**: Follows relationships across knowledge units
3. **Result Ranker**: Scores relevance based on context, confidence, freshness
4. **Persona Adapter**: Formats results for different user types
5. **Visualization Renderer**: Creates interactive graph views
6. **Impact Calculator**: Predicts change propagation

### Integration Points
- **Atomic Knowledge Graph**: Primary data source (atoms, molecules, organisms)
- **Code Repository**: For `IMPLEMENTS` relationship discovery
- **External Systems**: For freshness verification (tax law websites, etc.)
- **Collaboration Tools**: For ownership/responsibility discovery
- **AI/ML Systems**: For semantic understanding and gap detection

### Performance Considerations
- **Caching**: Frequently traversed relationships, popular queries
- **Indexing**: Atomic types, relationships, metadata, content
- **Incremental Updates**: As knowledge changes, update discovery indices
- **Distributed Processing**: For large-scale graph traversal

---

## Benefits Summary

### 1. **Reduced Duplication**
Teams find existing knowledge instead of recreating it. Example: 40% reduction in duplicate tax rule implementations.

### 2. **Faster Onboarding**
New team members discover organizational knowledge through structured exploration, not tribal knowledge transmission.

### 3. **Better Decision Making**
Discover all relevant constraints, precedents, and rationales before making decisions.

### 4. **Proactive Maintenance**
Discovery of stale, low-confidence, or orphaned knowledge before it causes issues.

### 5. **Audit Readiness**
Comprehensive discovery of compliance evidence, decision trails, and implementation verification.

### 6. **Knowledge Evolution**
Understanding how knowledge has changed over time through temporal discovery.

### 7. **Cross-Domain Insight**
Discovering unexpected connections between seemingly unrelated knowledge domains.

---

## Next Steps for Discovery System Design

1. **Prototype Query Language** - Implement `knowledge query` with basic atomic typing
2. **Build Graph Traversal Engine** - Start with `REQUIRES` and `IMPLEMENTS` relationships
3. **Design Persona Interfaces** - Create MVP for Developer and Domain Expert workflows
4. **Implement AI-Assisted Discovery** - Basic natural language to structured query translation
5. **Create Visualization Prototype** - Interactive graph browser for small knowledge sets
6. **Test with Real Examples** - Use withholding tax and GST examples from atomic framework
7. **Iterate Based on Feedback** - Refine discovery patterns based on user testing

---

## Related Documents

- [`project-os-atoms.md`](project-os-atoms.md) - Atomic Knowledge Design Framework
- [`product-os-v3-qa.md`](product-os-v3-qa.md) - Design Q&A (including visualization decisions)
- [`product-os-v3-concept.md`](product-os-v3-concept.md) - High-level concept and metaphor
- [`CONTINUATION.md`](CONTINUATION.md) - Current session status and next steps

---

**Status**: Design document - ready for prototyping discovery patterns  
**Key Insight**: Discovery is not a single feature but a **stack of complementary techniques** that together solve the composability discovery problem.