# Product-OS: Data Structure Concept

**Date**: 2026-01-11  
**Status**: Approved Architecture - Concept Document  
**Related**: [`product-os-data-structure-reference.md`](product-os-data-structure-reference.md) (Technical reference), [`product-os-implementation-persistence-operations.md`](product-os-implementation-persistence-operations.md) (Operational guidance)  
**Source**: Consolidated from [`data-structure-walkthrough.md`](archive/data-structure-walkthrough.md) (Iterative design process)

---

## Overview

Hybrid PostgreSQL + Neo4j architecture combining content storage (PostgreSQL) with relationship intelligence (Neo4j). Designed for knowledge management with atomic extraction, semantic relationships, and multi-tenant isolation.

**Core Design Principles**:
- **PostgreSQL**: Content storage, value queries, ACID transactions, multi-tenancy
- **Neo4j**: Graph traversals, impact analysis, relationship intelligence, dependency tracking
- **Repository Pattern**: Simple abstractions over persistence layer
- **Schema-per-Tenant**: PostgreSQL isolation with property-based Neo4j isolation
- **POC-First**: Progressive implementation with PostgreSQL → Neo4j sync

---

## 1. Hybrid Architecture

### Two-Layer Design

```
┌─────────────────────────────────────────────────────────────┐
│         CONTENT LAYER (PostgreSQL - Notion-inspired)       │
│  • Blocks: Typed content units (paragraphs, headings, etc.) │
│  • Documents: Collections of blocks with properties          │
│  • Atoms: Extracted values from content                      │
│  • Molecule Metadata: Rule configurations, enforcement       │
│                                                             │
│  STRENGTHS: Value queries, ACID transactions,               │
│            multi-tenancy, SQL analytics                      │
└──────────────────────────────┬──────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────┐
│        RELATIONSHIP LAYER (Neo4j - Graph Intelligence)     │
│  • Molecule Nodes: Business logic entities                  │
│  • All Relationships: PART_OF, RELATED_TO, DEPENDS_ON, etc. │
│  • Relationship Properties: Confidence, spans, timestamps   │
│                                                             │
│  STRENGTHS: Rule traversal, impact analysis,               │
│            relationship queries, dependency graphs          │
└─────────────────────────────────────────────────────────────┘
```

### Database Responsibilities

**PostgreSQL excels at**:
- Value-based queries (`WHERE number_value > 0.10`)
- Document storage with rich text/JSON
- ACID transactions for content updates
- Multi-tenancy with row/schema/database isolation

**Neo4j excels at**:
- Rule traversal ("Find all conditional rules affecting RRSP")
- Impact analysis ("What breaks if Quebec tax rate changes?")
- Relationship queries ("Show me concepts related to withholding tax")
- Circular dependency detection

### Data Flow
```
User writes document → PostgreSQL (blocks)
AI extracts atoms → PostgreSQL (atoms table)
AI identifies molecules → PostgreSQL (metadata) + Neo4j (graph)
Queries routed → Appropriate database based on query type
```

**Decision**: Two-layer architecture approved. PostgreSQL for content storage and value queries, Neo4j for relationship intelligence and graph traversals.

---

## 2. Core Entities

### 2.1 Atomic Knowledge Hierarchy
```
Atoms → Molecules → Organisms → Templates → Documents
```

**Atoms**: Primitive data types extracted from content
- `TextAtom`: Text values ("Quebec", "RRSP")
- `NumberAtom`: Numeric values (0.10, 5000)
- `BooleanAtom`: True/false values
- `DateAtom`: Date/time values
- `JSONAtom`: Structured data

**Molecules**: Business logic entities composed of atoms
- `ConditionMolecule`: If-then rules ("If jurisdiction=Quebec, then rate=10%")
- `ConstraintMolecule`: Business constraints ("Must be > 18 years old")
- `ExceptionMolecule`: Special cases ("Except for registered charities")
- `RequiresMolecule`: Dependency relationships

**Documents**: Notion-style pages with blocks and properties
- Collections of typed content blocks
- Schema-defined properties
- Domain organization

**Blocks**: Notion-inspired content units
- Paragraphs, headings, lists, code, tables, images
- Rich text storage in JSONB
- Hierarchical nesting

---

## 3. Key Design Decisions

### 3.1 PostgreSQL vs Neo4j Relationship Split

**PostgreSQL Relationships (Data Integrity & Simple Joins)**:
- Foreign key constraints for data integrity
- Simple 1-hop joins for document-block associations
- Array membership for atom-molecule composition
- Versioning chains for audit trails

**Neo4j Relationships (Knowledge Graph & Traversal)**:
- Multi-hop traversals for dependency analysis
- Path finding for impact analysis
- Relationship semantics (roles, confidence, spans)
- Circular dependency detection

### Decision Framework
| **Use Case** | **Database** | **Example** |
|--------------|--------------|-------------|
| **Data integrity** | PostgreSQL | Foreign key constraints |
| **Simple 1-hop joins** | PostgreSQL | `JOIN domains ON documents.domain_id = domains.id` |
| **Array membership** | PostgreSQL | `WHERE 'atom-123' = ANY(molecule_metadata.atom_ids)` |
| **Multi-hop traversals** | Neo4j | `MATCH (a)-[:REQUIRES*..3]->(b)` |
| **Path finding** | Neo4j | `MATCH path = shortestPath((a)-[*]-(b))` |
| **Relationship semantics** | Neo4j | Roles, confidence, spans on relationships |
| **Circular dependency detection** | Neo4j | Graph algorithms |
| **Impact analysis** | Neo4j | "What breaks if X changes?" |

**Key Insight**: PostgreSQL manages the **data**, Neo4j manages the **knowledge graph**.

### 3.2 Notion-Inspired Content Storage

**Why adopt Notion's block model**:
1. **Proven UX**: Billions of users understand block-based editors
2. **Technical simplicity**: No complex span calculations
3. **Future compatibility**: Easier to support rich media, tables, etc.
4. **Performance**: Block queries more efficient than text span queries

**But keep atomic layer for knowledge**:
1. **Semantic extraction**: AI identifies atoms in block content
2. **Knowledge graph**: Neo4j relationships enable complex queries
3. **Cross-document intelligence**: Find related concepts across documents
4. **Change impact analysis**: Trace dependencies through atomic relationships

---

## 4. Multi-Tenancy Approach

### Hybrid Isolation Model

**PostgreSQL**: Schema-per-tenant for strong isolation
- Separate schema for each tenant
- No accidental cross-tenant data leaks
- Easy tenant migration/archival
- Role-based access control per schema

**Neo4j**: Property-based isolation with `tenantId`
- Single graph with tenant filtering
- Allows cross-tenant analytics (admin views)
- Simpler deployment and maintenance
- Indexes on `tenantId` for performance

### Progressive Migration Strategy
1. **Phase 1 (POC)**: Row-level isolation with `tenant_id` columns
2. **Phase 2 (Growth)**: Schema-per-tenant migration
3. **Phase 3 (Enterprise)**: Dedicated databases for large tenants

---

## 5. Synchronization Strategy

### Bidirectional Sync with Clear Source of Truth

**PostgreSQL as source of truth**:
- Content, values, configurations
- Strong consistency required
- Primary flow: PostgreSQL → Neo4j

**Neo4j as source of truth**:
- Relationships, graph structure
- Eventual consistency acceptable
- Secondary flow: Neo4j → PostgreSQL (array updates)

### POC Implementation Approach
**Phase 1**: PostgreSQL → Neo4j Sync Only
- PostgreSQL triggers for basic entities
- Simple sync service creates/updates Neo4j nodes
- Manual relationship creation in Neo4j

**Phase 2**: Bidirectional Sync
- Neo4j transaction listeners
- PostgreSQL array updates from relationships
- Conflict detection/resolution

**Phase 3**: Production Ready
- Message queue for async processing
- Retry logic with exponential backoff
- Monitoring and alerting

---

## 6. Repository Pattern

### Simple Abstractions Over Persistence

**Core Idea**: Each entity (Atom, Molecule, Document, Block) has a repository that encapsulates database choice.

**Benefits**:
- **Simpler mental model**: Developers work with entities, not databases
- **Clean separation**: Repository layer → Service layer → API layer
- **Easier testing**: Mock repositories without database dependencies
- **Progressive enhancement**: Start simple, add graph operations later

**Example Usage**:
```typescript
// Developers use repositories naturally
const quebecAtom = await atomRepo.findByValue('Quebec');
const molecules = await atomRepo.getMoleculesContainingAtom(quebecAtom.id);
const impact = await moleculeRepo.impactAnalysis(molecule.id);
```

---

## 7. Open Questions for Refinement

1. **Organisms table**: Should we add explicit `organisms` table or treat them as special `documents`?
2. **Templates**: Need `templates` table for document schemas?
3. **Vector embeddings**: Store in PostgreSQL (pgvector) or separate vector DB?
4. **Audit trail**: Separate `knowledge_events` table for event sourcing?
5. **Persona access**: Additional tables for persona-specific views/permissions?
6. **Capabilities**: Table linking knowledge to business capabilities?
7. **Code references**: Separate table for `implements` relationships to code?
8. **Domain relationships**: Should domains have parent-child hierarchy? (e.g., `tax/canadian`, `tax/us`)
9. **Domain templates**: Should domains enforce specific templates for their knowledge?

---

## 8. Implementation Roadmap

### POC Phase (Weeks 1-4)
1. PostgreSQL schema implementation (blocks, documents, atoms, molecules, domains)
2. Basic Neo4j node creation sync
3. Simple repository implementations
4. Row-level multi-tenancy

### Growth Phase (Weeks 5-8)
1. Schema-per-tenant migration
2. Bidirectional sync implementation
3. Graph query patterns
4. Performance optimization

### Production Phase (Weeks 9-12)
1. Message queue integration
2. Monitoring and alerting
3. Backup and recovery procedures
4. Security hardening

---

**Status**: Approved architecture - ready for POC implementation  
**Next**: Begin POC with PostgreSQL schema + basic Neo4j sync  
**Reference**: For complete technical details, see [`product-os-data-structure-reference.md`](product-os-data-structure-reference.md)