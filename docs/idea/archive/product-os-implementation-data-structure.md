# Product-OS: Data Structure Reference

**Date**: 2026-01-11  
**Status**: Technical Reference - Approved Architecture  
**Related**: [`product-os-implementation-persistence-operations.md`](product-os-implementation-persistence-operations.md) (Operational guidance), [`product-os-atoms.md`](product-os-atoms.md) (Atomic framework)  
**Source**: Derived from [`data-structure-walkthrough.md`](product-os-concept-work-in-progress/data-structure/data-structure-walkthrough.md) (Iterative design process)

---

## Overview

Hybrid PostgreSQL + Neo4j architecture combining content storage (PostgreSQL) with relationship intelligence (Neo4j). Approved through iterative concept review with 6 finalized concepts. For detailed design discussions and rationale, see the source walkthrough document.

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

**Decision**: Two-layer architecture approved. PostgreSQL for content storage and value queries, Neo4j for relationship intelligence and graph traversals. (See walkthrough Concept 1 for detailed rationale)

---

## 2. PostgreSQL Schema

### 2.1 `blocks` - Notion-style content units
```sql
CREATE TABLE blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    
    -- Hierarchy (parent can be block, document, or domain)
    parent_id UUID REFERENCES blocks(id) ON DELETE CASCADE,
    parent_type VARCHAR(20), -- 'block', 'document', 'domain'
    
    -- Content type
    block_type VARCHAR(50) NOT NULL,
    -- Types: 'paragraph', 'heading_1', 'heading_2', 'heading_3',
    --        'bulleted_list_item', 'numbered_list_item', 'code',
    --        'quote', 'callout', 'table', 'image', 'file'
    
    -- Content storage
    content JSONB NOT NULL, -- Rich text array or block-specific config
    
    -- Atomic knowledge references (for fast lookups)
    atom_ids UUID[],
    molecule_ids UUID[],
    
    -- Metadata
    has_children BOOLEAN DEFAULT FALSE,
    archived BOOLEAN DEFAULT FALSE,
    in_trash BOOLEAN DEFAULT FALSE,
    
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID,
    last_edited_by UUID
);
```

**Purpose**: Stores the actual content users create/edit, inspired by Notion's block model. Each block is a typed unit that can contain rich text, images, code, etc.

### 2.2 `documents` - Collections of blocks
```sql
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    domain_id UUID REFERENCES domains(id),
    
    -- Hierarchy
    parent_id UUID, -- Could reference domain, block, or another document
    parent_type VARCHAR(20), -- 'domain', 'block', 'document'
    
    -- Content root
    root_block_id UUID REFERENCES blocks(id),
    
    -- Schema-defined properties (like Notion database properties)
    properties JSONB, -- {property_name: value} conforming to template
    
    -- UI/UX elements
    icon JSONB, -- emoji or file reference
    cover JSONB, -- file reference for cover image
    title VARCHAR(500), -- denormalized for search
    
    -- Status
    archived BOOLEAN DEFAULT FALSE,
    in_trash BOOLEAN DEFAULT FALSE,
    
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID,
    last_edited_by UUID
);
```

**Purpose**: Acts as "pages" in Notion terminology. A document is a collection of blocks with additional metadata and properties.

### 2.3 `atoms` - Extracted knowledge values
```sql
CREATE TABLE atoms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    domain_id UUID REFERENCES domains(id) ON DELETE SET NULL,
    atom_type VARCHAR(50) NOT NULL,
    -- Types: 'TextAtom', 'NumberAtom', 'BooleanAtom', 'DateAtom', 'JSONAtom'
    
    -- Polymorphic value storage (one of these will be non-NULL)
    text_value TEXT,
    number_value NUMERIC,
    boolean_value BOOLEAN,
    date_value TIMESTAMPTZ,
    json_value JSONB,
    
    -- Metadata
    name VARCHAR(255),
    description TEXT,
    confidence FLOAT DEFAULT 1.0,
    source VARCHAR(50) DEFAULT 'ai_extracted',
    
    -- Versioning
    version INTEGER DEFAULT 1,
    previous_version UUID REFERENCES atoms(id),
    
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
```

**Purpose**: Stores atomic knowledge units extracted from content. These are the primitive values that molecules and rules operate on.

### 2.4 `molecule_metadata` - Business logic configuration
```sql
CREATE TABLE molecule_metadata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    domain_id UUID REFERENCES domains(id) ON DELETE SET NULL,
    molecule_type VARCHAR(50) NOT NULL,
    -- Types: 'ConditionMolecule', 'ConstraintMolecule', 'ExceptionMolecule', 'RequiresMolecule'
    
    -- Atom composition (for SQL queries)
    atom_ids UUID[] NOT NULL,
    
    -- Type-specific configuration
    configuration JSONB,
    -- Example: {rule_type: "if_then", condition: "jurisdiction=Quebec", result: "rate=0.10"}
    
    -- Validation and enforcement
    validation_rules JSONB,
    enforcement_level VARCHAR(20) DEFAULT 'warn', -- 'block', 'warn', 'inform'
    
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

**Purpose**: Stores the business logic and configuration for molecules. The actual graph relationships will be in Neo4j, but the rule definitions stay here for SQL queryability.

### 2.5 `domains` - Bounded contexts
```sql
CREATE TABLE domains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL, -- 'tax', 'compliance', 'billing'
    path VARCHAR(1000) NOT NULL, -- '/knowledge/domains/tax/'
    description TEXT,
    owner_id UUID, -- Domain expert/owner
    
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    UNIQUE(tenant_id, path)
);
```

**Purpose**: Organizes knowledge into DDD bounded contexts. Each domain represents a specific area of expertise.

### Key Design Decisions
1. **Blocks use JSONB for content**: Allows flexible rich text storage
2. **Polymorphic atoms**: Single table with type-specific columns
3. **Molecule metadata separate from graph**: Business logic in PostgreSQL, relationships in Neo4j
4. **Tenant isolation**: All tables include `tenant_id`

**Decision**: PostgreSQL tables approved as presented. Raw text will be derived from block JSONB content. No separate raw_text storage. (See walkthrough Concept 2 for table design discussions)

---

## 3. Neo4j Graph Model

### 3.1 Node Types & Properties
```cypher
// ATOM NODES (extracted knowledge values)
(:Atom {
  id: 'atom-123',                    // UUID matching PostgreSQL atoms.id
  type: 'NumberAtom',                // 'TextAtom', 'NumberAtom', 'BooleanAtom', etc.
  value: 0.10,                       // Type-specific value
  canonical: '10%',                  // Human-readable form
  name: 'Quebec Tax Rate',           // Optional descriptive name
  confidence: 0.98,                  // Extraction confidence
  source: 'ai_extracted',            // 'ai_extracted', 'manual', 'code_analysis'
  tenantId: 'tenant-xyz',            // Multi-tenancy isolation
  domain: 'tax',                     // Domain context
  createdAt: timestamp(),            // Creation timestamp
  updatedAt: timestamp()             // Last update timestamp
})

// MOLECULE NODES (business logic entities)
(:Molecule {
  id: 'molecule-456',                // UUID matching PostgreSQL molecule_metadata.id
  type: 'ConditionMolecule',         // 'ConditionMolecule', 'ConstraintMolecule', etc.
  moleculeType: 'ConditionMolecule', // Duplicate for query convenience
  ruleType: 'if_then',               // Type-specific rule classification
  condition: "jurisdiction equals Quebec", // Human-readable condition
  result: "rate is 10%",             // Human-readable result
  configuration: {                   // JSON configuration
    conditionType: "jurisdiction_equals",
    jurisdictionValue: "Quebec",
    rateValue: 0.10
  },
  enforcement: 'block',              // 'block', 'warn', 'inform'
  tenantId: 'tenant-xyz',
  domain: 'tax',
  createdAt: timestamp(),
  updatedAt: timestamp()
})

// DOCUMENT NODES (content collections)
(:Document {
  id: 'doc-789',                     // UUID matching PostgreSQL documents.id
  title: 'RRSP Withholding Tax Guidelines',
  type: 'document',                  // Could also be 'template', 'organism'
  tenantId: 'tenant-xyz',
  domain: 'tax',
  createdAt: timestamp(),
  updatedAt: timestamp()
})

// BLOCK NODES (content units)
(:Block {
  id: 'block-abc',                   // UUID matching PostgreSQL blocks.id
  blockType: 'paragraph',            // 'paragraph', 'heading_1', 'code', etc.
  tenantId: 'tenant-xyz',
  createdAt: timestamp(),
  updatedAt: timestamp()
})
```

### 3.2 Relationship Types & Properties
```cypher
// COMPOSITION RELATIONSHIPS
(:Atom {id: 'atom-123'})-[:PART_OF {
  role: 'condition',                 // 'subject', 'condition', 'result', 'constraint'
  order: 1,                          // Position in molecule
  createdAt: timestamp()
}]->(:Molecule {id: 'molecule-456'})

// DOCUMENT APPEARANCE RELATIONSHIPS
(:Atom {id: 'atom-123'})-[:APPEARS_IN {
  spans: [{start: 45, end: 56}],    // Character positions in document
  blockId: 'block-abc',              // Reference to specific block
  confidence: 0.98,
  source: 'ai',
  createdAt: timestamp()
}]->(:Document {id: 'doc-789'})

// KNOWLEDGE RELATIONSHIPS
(:Atom {id: 'atom-123'})-[:RELATED_TO {
  type: 'similar_to',                // 'similar_to', 'synonym_of', 'antonym_of'
  strength: 0.85,
  createdAt: timestamp()
}]->(:Atom {id: 'atom-125'})

(:Molecule {id: 'molecule-456'})-[:REQUIRES {
  type: 'prerequisite',              // 'prerequisite', 'dependency', 'conflict'
  enforcement: 'block',
  createdAt: timestamp()
}]->(:Atom {id: 'atom-126'})

(:Molecule {id: 'molecule-456'})-[:CONFLICTS_WITH {
  reason: 'contradictory rates',
  severity: 'high',
  createdAt: timestamp()
}]->(:Molecule {id: 'molecule-457'})
```

### 3.3 PostgreSQL vs Neo4j Relationship Split

**PostgreSQL Relationships (Data Integrity & Simple Joins)**
```sql
-- FOREIGN KEY RELATIONSHIPS
blocks.parent_id REFERENCES blocks(id)              -- Block nesting
documents.domain_id REFERENCES domains(id)          -- Document-domain association
documents.root_block_id REFERENCES blocks(id)       -- Document root block

-- Composition references  
molecule_metadata.atom_ids UUID[]                   -- Atom membership array
blocks.atom_ids UUID[]                              -- Block-atom references
blocks.molecule_ids UUID[]                          -- Block-molecule references

-- Versioning chains
atoms.previous_version REFERENCES atoms(id)         -- Atom version history
```

**Purpose**: Data integrity, ACID compliance, simple join queries, array operations.

**Neo4j Relationships (Knowledge Graph & Traversal)**
```cypher
-- GRAPH RELATIONSHIPS
(:Atom)-[:PART_OF {role, order}]->(:Molecule)       -- Composition with semantics
(:Atom)-[:APPEARS_IN {spans, confidence}]->(:Document) -- Document appearance
(:Atom)-[:RELATED_TO {type, strength}]->(:Atom)     -- Semantic relationships
(:Molecule)-[:REQUIRES]->(:Atom)                    -- Dependency relationships
(:Molecule)-[:CONFLICTS_WITH]->(:Molecule)          -- Rule conflicts
(:Block)-[:PART_OF_DOCUMENT]->(:Document)           -- Document structure
(:Atom)-[:IN_DOMAIN]->(:Domain)                     -- Domain organization
```

**Purpose**: Graph traversals, impact analysis, multi-hop queries, relationship semantics.

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

**Decision**: Neo4j graph model approved with clarified PostgreSQL vs Neo4j relationship split. (See walkthrough Concept 3 for detailed Cypher examples including RRSP tax rule implementation)

---

## 4. Synchronization & Data Flow

### 4.1 Synchronization Requirements

**Entities to Sync**:
1. **Atoms**: PostgreSQL `atoms` ↔ Neo4j `Atom` nodes
2. **Molecule Metadata**: PostgreSQL `molecule_metadata` ↔ Neo4j `Molecule` nodes  
3. **Documents**: PostgreSQL `documents` ↔ Neo4j `Document` nodes
4. **Blocks**: PostgreSQL `blocks` ↔ Neo4j `Block` nodes
5. **Domains**: PostgreSQL `domains` ↔ Neo4j `Domain` nodes

**Sync Direction**: Bidirectional, with clear source of truth:
- **PostgreSQL as source of truth**: Content, values, configurations
- **Neo4j as source of truth**: Relationships, graph structure

### 4.2 Synchronization Patterns

**Pattern 1: PostgreSQL → Neo4j (Primary Flow)**
```typescript
interface PostgresToNeo4jSync {
  trigger: 'postgres_change',  // INSERT/UPDATE/DELETE in PostgreSQL
  entity: 'atom' | 'molecule' | 'document' | 'block',
  operation: 'create' | 'update' | 'delete',
  postgresId: UUID,
  tenantId: UUID,
  timestamp: Date
}

// Sync actions:
// CREATE → Create node in Neo4j with all properties
// UPDATE → Update node properties in Neo4j
// DELETE → Remove node from Neo4j (cascade relationships)
```

**Pattern 2: Neo4j → PostgreSQL (Relationship Updates)**
```typescript
interface Neo4jToPostgresSync {
  trigger: 'neo4j_relationship_change',
  relationship: 'PART_OF' | 'APPEARS_IN' | 'RELATED_TO' | 'REQUIRES',
  sourceType: 'Atom' | 'Molecule' | 'Document',
  sourceId: UUID,
  targetType: 'Atom' | 'Molecule' | 'Document',
  targetId: UUID,
  operation: 'create' | 'delete',
  properties: Record<string, any>,  // roles, confidence, spans, etc.
  tenantId: UUID
}

// Sync actions:
// PART_OF created → Add atom_id to molecule_metadata.atom_ids
// APPEARS_IN created → Update block atom_ids array
// RELATIONSHIP deleted → Remove from PostgreSQL arrays
```

### 4.3 Sync State Tracking
```sql
CREATE TABLE sync_state (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    
    -- Entity being synced
    entity_type VARCHAR(50) NOT NULL,  -- 'atom', 'molecule', 'document', 'block'
    entity_id UUID NOT NULL,
    
    -- Sync state
    postgres_version INTEGER NOT NULL,
    neo4j_version INTEGER,
    
    -- Checksums for change detection
    postgres_checksum CHAR(64),  -- SHA-256 of PostgreSQL row
    neo4j_checksum CHAR(64),     -- SHA-256 of Neo4j node
    
    -- Sync status
    sync_status VARCHAR(20) DEFAULT 'pending',  -- 'pending', 'syncing', 'synced', 'failed'
    last_sync_attempt TIMESTAMPTZ,
    sync_error TEXT,
    
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    UNIQUE(tenant_id, entity_type, entity_id)
);

CREATE INDEX sync_state_pending_idx ON sync_state(tenant_id, sync_status) 
WHERE sync_status IN ('pending', 'failed');
```

### 4.4 Consistency Models

**Strong Consistency Required**:
- **Atom values**: PostgreSQL is source of truth, Neo4j must match
- **Molecule configurations**: PostgreSQL is source of truth
- **Document/block content**: PostgreSQL is source of truth

**Eventual Consistency Acceptable**:
- **Relationship updates**: Neo4j changes can propagate to PostgreSQL arrays async
- **Graph structure**: Temporary inconsistencies ok for query routing

### 4.5 POC Implementation Approach

**Phase 1: PostgreSQL → Neo4j Sync Only**
1. Implement PostgreSQL triggers for basic entities
2. Simple sync service that creates/updates Neo4j nodes
3. Manual relationship creation in Neo4j

**Phase 2: Bidirectional Sync**
1. Add Neo4j transaction listeners
2. Implement array updates in PostgreSQL
3. Add conflict detection/resolution

**Phase 3: Production Ready**
1. Add message queue for async processing
2. Implement retry logic with exponential backoff
3. Add monitoring and alerting
4. Performance optimization

**Decision**: Synchronization approach approved. POC will implement PostgreSQL → Neo4j sync with triggers and simple sync service. Full bidirectional automation will be iterated post-POC. (See walkthrough Concept 4 for sync patterns and compensation logic)

---

## 5. Repository Pattern

### 5.1 Repository Design

**Core Idea**: Each entity (Atom, Molecule, Document, Block) has a repository that encapsulates database choice.

```typescript
// Atom repository - encapsulates database choice
class AtomRepository {
  private postgresClient: PostgresClient;
  private neo4jClient: Neo4jClient;
  
  // Content operations → PostgreSQL
  async findById(id: UUID): Promise<Atom | null> {
    return this.postgresClient.query(
      'SELECT * FROM atoms WHERE id = $1',
      [id]
    );
  }
  
  async save(atom: Atom): Promise<Atom> {
    // Save to PostgreSQL (source of truth for values)
    const saved = await this.postgresClient.query(
      'INSERT INTO atoms VALUES ($1, $2, ...) RETURNING *',
      [atom.id, atom.tenantId, ...]
    );
    
    // Sync to Neo4j for graph relationships
    await this.syncToNeo4j(saved);
    
    return saved;
  }
  
  // Graph operations → Neo4j
  async findRelatedConcepts(atomId: UUID, depth: number = 2): Promise<ConceptGraph> {
    return this.neo4jClient.query(`
      MATCH (start:Atom {id: $atomId})
      MATCH path = (start)-[:RELATED_TO*1..${depth}]-(related:Atom)
      RETURN nodes(path) as nodes, relationships(path) as relationships
    `, { atomId });
  }
  
  async getMoleculesContainingAtom(atomId: UUID): Promise<Molecule[]> {
    // Get molecule IDs from Neo4j graph
    const moleculeIds = await this.neo4jClient.query(`
      MATCH (a:Atom {id: $atomId})<-[:PART_OF]-(m:Molecule)
      RETURN m.id as id
    `, { atomId });
    
    // Get molecule details from PostgreSQL
    return this.postgresClient.query(
      'SELECT * FROM molecule_metadata WHERE id = ANY($1)',
      [moleculeIds.map(m => m.id)]
    );
  }
}
```

### 5.2 Developer Usage
```typescript
// 1. Import repositories
import { AtomRepository, MoleculeRepository, DocumentRepository } from '@knowledge-os/repositories';

// 2. Initialize (dependency injection)
const atomRepo = new AtomRepository(postgresClient, neo4jClient);
const moleculeRepo = new MoleculeRepository(postgresClient, neo4jClient);
const documentRepo = new DocumentRepository(postgresClient);

// 3. Use repositories naturally - no database decisions needed
async function analyzeTaxRuleImpact() {
  // Find Quebec tax rate atom
  const quebecAtom = await atomRepo.findByValue('Quebec');
  
  // Get molecules containing this atom
  const molecules = await atomRepo.getMoleculesContainingAtom(quebecAtom.id);
  
  // Impact analysis for each molecule
  for (const molecule of molecules) {
    const impact = await moleculeRepo.impactAnalysis(molecule.id);
    console.log(`Molecule ${molecule.id} affects:`, impact);
  }
}
```

### 5.3 Benefits Over Complex Router

**Simpler Mental Model**:
- Developers work with **entities**, not databases
- Each repository method has clear responsibility
- No need to understand "which database for which query"

**Clean Separation**:
- **Repository layer**: Entity persistence and relationships
- **Service layer**: Business logic using repositories
- **API layer**: HTTP endpoints calling services

**Easier Testing**:
```typescript
// Mock repositories for testing
const mockAtomRepo: AtomRepository = {
  findById: jest.fn().mockResolvedValue(mockAtom),
  findRelatedConcepts: jest.fn().mockResolvedValue([])
};

// Test business logic without database dependencies
const taxService = new TaxService(mockAtomRepo, mockMoleculeRepo);
await taxService.analyzeImpact('quebec-tax-rate');
expect(mockAtomRepo.findById).toHaveBeenCalledWith('quebec-tax-rate');
```

**Decision**: Repository pattern approved. Each entity will have a repository that encapsulates database choice while allowing custom queries through clear abstractions. (See walkthrough Concept 5 for detailed repository implementations and comparison with complex router)

---

## 6. Multi-Tenancy Implementation

### 6.1 Hybrid Isolation Approach

**Recommended Approach**: Hybrid model combining strengths:
- **PostgreSQL**: Schema-per-tenant for strong isolation
- **Neo4j**: Row-level isolation with `tenantId` property
- **Shared infrastructure**: For non-sensitive, scalable components

### 6.2 PostgreSQL: Schema-per-Tenant
```sql
-- 1. Create tenant schema template
CREATE SCHEMA IF NOT EXISTS tenant_template;

-- 2. Template tables (same as main schema but with tenant isolation)
CREATE TABLE tenant_template.blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- No tenant_id column needed - implicit by schema
    parent_id UUID,
    block_type VARCHAR(50) NOT NULL,
    content JSONB NOT NULL,
    atom_ids UUID[],
    molecule_ids UUID[],
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
    -- ... other columns
);

-- 3. Tenant management table (in public schema)
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    subdomain VARCHAR(100) UNIQUE NOT NULL,  -- e.g., 'acme.knowledgeos.app'
    schema_name VARCHAR(100) NOT NULL,       -- e.g., 'tenant_acme_xyz'
    status VARCHAR(20) DEFAULT 'active',     -- 'active', 'suspended', 'deleted'
    plan_type VARCHAR(50) DEFAULT 'starter', -- 'starter', 'professional', 'enterprise'
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

**Benefits of Schema-per-Tenant**:
- **Strong isolation**: No accidental cross-tenant data leaks
- **Performance**: Tenant data co-located, better cache locality
- **Maintenance**: Easy tenant migration/archival (drop schema)
- **Security**: PostgreSQL role-based access control per schema

### 6.3 Neo4j: Property-based Isolation
```cypher
// All nodes include tenantId property
CREATE (a:Atom {
  id: 'atom-123',
  tenantId: 'tenant-acme-xyz',  // Tenant identifier
  value: 'Quebec',
  type: 'TextAtom'
})

// All queries filter by tenantId
MATCH (a:Atom {tenantId: $tenantId, value: 'Quebec'})
RETURN a

// Index for tenant isolation performance
CREATE INDEX tenant_atom_idx FOR (a:Atom) ON (a.tenantId, a.id);
CREATE INDEX tenant_molecule_idx FOR (m:Molecule) ON (m.tenantId, m.id);
```

**Why Property-based for Neo4j**:
- **Graph traversals**: Need to cross tenant boundaries for shared knowledge (optional)
- **Performance**: Single graph allows cross-tenant analytics (admin views)
- **Simplicity**: No need for multiple Neo4j instances/databases

### 6.4 Repository Pattern with Multi-Tenancy
```typescript
// Base repository with tenant context
abstract class TenantAwareRepository<T> {
  protected readonly tenantId: UUID;
  protected readonly postgresClient: PostgresClient;
  protected readonly neo4jClient: Neo4jClient;
  
  constructor(tenantId: UUID, postgresClient: PostgresClient, neo4jClient: Neo4jClient) {
    this.tenantId = tenantId;
    this.postgresClient = postgresClient;
    this.neo4jClient = neo4jClient;
  }
  
  // PostgreSQL query with schema selection
  protected async postgresQuery<T>(sql: string, params: any[] = []): Promise<T> {
    // Set search_path to tenant schema
    await this.postgresClient.query(`SET search_path TO tenant_${this.tenantId}, public`);
    
    // Execute query
    return this.postgresClient.query(sql, params);
  }
  
  // Neo4j query with tenant filter
  protected async neo4jQuery<T>(cypher: string, params: Record<string, any> = {}): Promise<T> {
    // Add tenantId to all queries
    const tenantParams = { ...params, tenantId: this.tenantId };
    
    // For security, all MATCH clauses should include tenantId
    const securedCypher = this.secureCypherQuery(cypher);
    
    return this.neo4jClient.query(securedCypher, tenantParams);
  }
  
  private secureCypherQuery(cypher: string): string {
    // Simple security wrapper - ensures all node patterns include {tenantId: $tenantId}
    return cypher.replace(/MATCH\s*\(([^:]+):/g, (match, varName) => {
      return `MATCH (${varName} {tenantId: $tenantId}:`;
    });
  }
}
```

### 6.5 Progressive Migration Strategy

**Phase 1: Single Tenant (POC)**
```typescript
// Simple tenant_id columns, no schema separation
class POCAtomRepository {
  async findById(tenantId: UUID, atomId: UUID): Promise<Atom> {
    return this.postgresClient.query(
      'SELECT * FROM atoms WHERE tenant_id = $1 AND id = $2',
      [tenantId, atomId]
    );
  }
}
```

**Phase 2: Schema-per-Tenant (Growth)**
```typescript
// Migrate from row-level to schema-per-tenant
async function migrateToSchemaPerTenant(tenantId: UUID): Promise<void> {
  // 1. Create tenant schema
  await createTenantSchema(tenantId);
  
  // 2. Copy data from public.atoms WHERE tenant_id = $1
  await this.postgresClient.query(`
    INSERT INTO tenant_${tenantId}.atoms
    SELECT * FROM public.atoms WHERE tenant_id = $1
  `, [tenantId]);
  
  // 3. Update Neo4j nodes with schema reference
  await this.neo4jClient.query(`
    MATCH (a:Atom {tenantId: $tenantId})
    SET a.schemaName = $schemaName
  `, { tenantId, schemaName: `tenant_${tenantId}` });
}
```

**Phase 3: Multi-Database (Enterprise)**
```typescript
// For large tenants, move to dedicated database
class EnterpriseTenantRepository extends TenantAwareRepository<Atom> {
  private tenantDbClient: PostgresClient; // Dedicated connection pool
  
  constructor(tenantId: UUID) {
    super(tenantId);
    this.tenantDbClient = this.connectToTenantDatabase(tenantId);
  }
  
  async findById(id: UUID): Promise<Atom> {
    // Query dedicated tenant database
    return this.tenantDbClient.query(
      'SELECT * FROM atoms WHERE id = $1',
      [id]
    );
  }
}
```

### Key Design Decisions
1. **Hybrid isolation**: PostgreSQL schema-per-tenant + Neo4j property-based
2. **Progressive migration**: Row-level → Schema-per-tenant → Dedicated database
3. **Repository pattern integration**: Tenant context baked into repositories
4. **Security first**: Automatic query rewriting for tenant isolation
5. **Compliance ready**: GDPR, data retention, subject requests

**Decision**: Multi-tenancy approach approved with hybrid isolation model. (See walkthrough Concept 6 for progressive migration strategies and compliance implementation)

---

## Open Questions for Refinement

1. **Organisms table**: Should we add explicit `organisms` table or treat them as special `documents`?
2. **Templates**: Need `templates` table for document schemas?
3. **Relationships**: PostgreSQL foreign keys vs Neo4j-only relationships?
4. **Vector embeddings**: Store in PostgreSQL (pgvector) or separate vector DB?
5. **Audit trail**: Separate `knowledge_events` table for event sourcing?
6. **Persona access**: Additional tables for persona-specific views/permissions?
7. **Capabilities**: Table linking knowledge to business capabilities?
8. **Code references**: Separate table for `implements` relationships to code?
9. **Domain relationships**: Should domains have parent-child hierarchy? (e.g., `tax/canadian`, `tax/us`)
10. **Domain templates**: Should domains enforce specific templates for their knowledge?

---

**Status**: Approved architecture - ready for POC implementation  
**Next**: Begin POC with PostgreSQL schema + basic Neo4j sync  
**Reference**: For complete design discussions, rationale, and detailed examples, consult the source [`data-structure-walkthrough.md`](product-os-concept-work-in-progress/data-structure/data-structure-walkthrough.md)