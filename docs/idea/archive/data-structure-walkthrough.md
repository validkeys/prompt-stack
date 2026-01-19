# Data Structure Walkthrough: Hybrid PostgreSQL + Neo4j Architecture

**Date**: 2026-01-11  
**Participants**: Architect (User), Data Engineer (Assistant)  
**Goal**: Finalize data structure through iterative concept review  
**Status**: In Progress  

---

## Session Overview

This document captures our iterative design process, presenting one concept at a time for review, feedback, and finalization. We are designing a hybrid architecture for Knowledge OS that combines PostgreSQL (content storage) with Neo4j (relationship intelligence).

---

## Concept 1: The Two-Layer Architecture

### **Data Engineer Presentation**

**Problem Statement**: 
We need a system that can handle both structured content storage (like Notion for documents) and intelligent relationship queries (graph traversals for knowledge discovery). These are fundamentally different access patterns that no single database excels at simultaneously.

**Proposed Solution**: Split responsibilities between two specialized databases:

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

**Why This Split?**:
1. **PostgreSQL excels at**: 
   - Value-based queries (`WHERE number_value > 0.10`)
   - Document storage with rich text/JSON
   - ACID transactions for content updates
   - Multi-tenancy with row/schema/database isolation

2. **Neo4j excels at**:
   - Rule traversal ("Find all conditional rules affecting RRSP")
   - Impact analysis ("What breaks if Quebec tax rate changes?")
   - Relationship queries ("Show me concepts related to withholding tax")
   - Circular dependency detection

**Data Flow Example**:
```
User writes document → PostgreSQL (blocks)
AI extracts atoms → PostgreSQL (atoms table)
AI identifies molecules → PostgreSQL (metadata) + Neo4j (graph)
Queries routed → Appropriate database based on query type
```

**Architect Feedback**:
✓ **Approved**: "I think this makes sense."

**Final Decision**: Two-layer architecture approved. PostgreSQL for content storage and value queries, Neo4j for relationship intelligence and graph traversals.

---

## Concept 2: PostgreSQL Tables (Content Storage)

### **Data Engineer Presentation**

**Overview**: PostgreSQL will store all structured content, extracted values, and rule configurations. This includes Notion-inspired blocks, documents, atomic knowledge units, and molecule metadata.

**Proposed Tables**:

#### **2.1 `blocks` - Notion-style content units**
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

#### **2.2 `documents` - Collections of blocks**
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

#### **2.3 `atoms` - Extracted knowledge values**
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

#### **2.4 `molecule_metadata` - Business logic configuration**
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

#### **2.5 `domains` - Bounded contexts**
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

**Key Design Decisions**:
1. **Blocks use JSONB for content**: Allows flexible rich text storage
2. **Polymorphic atoms**: Single table with type-specific columns
3. **Molecule metadata separate from graph**: Business logic in PostgreSQL, relationships in Neo4j
4. **Tenant isolation**: All tables include `tenant_id`

**Architect Feedback**:
✓ **Approved**: PostgreSQL tables structure approved.
✓ **Raw Text Decision**: "I agree with your recommendation. Derive from blocks, we can extract later if needed post-POC. I don't think we need the raw_text."

**Clarifications from Architect Questions**:
1. **Tables are sufficient** for content storage in POC
2. **No `properties` table needed yet** - can add later if Notion property system becomes important
3. **Annotations/block-atom mappings**: Will be handled in Neo4j relationships (APPEARS_IN, PART_OF, etc.)
4. **Polymorphic atom design is acceptable** - matches existing schema

**Final Decision**: PostgreSQL tables approved as presented. Raw text will be derived from block JSONB content. No separate raw_text storage.

---

## Concept 3: Neo4j Graph Model (Relationship Intelligence)

### **Data Engineer Presentation**

**Overview**: Neo4j will store all relationships between knowledge entities - molecules, atoms, documents, and blocks. This is our single source of truth for relationship intelligence, enabling complex graph traversals, impact analysis, and dependency tracking.

**Proposed Graph Schema**:

#### **3.1 Node Types & Properties**

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

#### **3.2 Relationship Types & Properties**

```cypher
// COMPOSITION RELATIONSHIPS
(:Atom {id: 'atom-123'})-[:PART_OF {
  role: 'condition',                 // 'subject', 'condition', 'result', 'constraint'
  order: 1,                          // Position in molecule
  createdAt: timestamp()
}]->(:Molecule {id: 'molecule-456'})

(:Atom {id: 'atom-124'})-[:PART_OF {
  role: 'result',
  order: 2,
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

// BLOCK CONTAINMENT
(:Block {id: 'block-abc'})-[:PART_OF_DOCUMENT]->(:Document {id: 'doc-789'})
(:Block {id: 'block-def'})-[:CHILD_OF]->(:Block {id: 'block-abc'})  // Nested blocks

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

// DOMAIN ORGANIZATION
(:Atom {id: 'atom-123'})-[:IN_DOMAIN]->(:Domain {name: 'tax'})
(:Molecule {id: 'molecule-456'})-[:IN_DOMAIN]->(:Domain {name: 'tax'})
(:Document {id: 'doc-789'})-[:IN_DOMAIN]->(:Domain {name: 'tax'})
```

#### **3.3 Example: RRSP Tax Rule in Neo4j**

```cypher
// Create the graph structure for: "RRSP withholding tax is 10% if in Quebec"
CREATE (rrsp:Atom {
  id: 'atom-rrsp',
  type: 'TextAtom',
  value: 'RRSP',
  canonical: 'Registered Retirement Savings Plan',
  tenantId: 'tenant-xyz',
  domain: 'tax'
})

CREATE (withholding:Atom {
  id: 'atom-withholding',
  type: 'TextAtom', 
  value: 'withholding tax',
  canonical: 'Withholding Tax',
  tenantId: 'tenant-xyz',
  domain: 'tax'
})

CREATE (quebec:Atom {
  id: 'atom-quebec',
  type: 'TextAtom',
  value: 'Quebec',
  canonical: 'Quebec',
  tenantId: 'tenant-xyz',
  domain: 'tax'
})

CREATE (tenPercent:Atom {
  id: 'atom-10percent',
  type: 'NumberAtom',
  value: 0.10,
  canonical: '10%',
  tenantId: 'tenant-xyz',
  domain: 'tax'
})

CREATE (rule:Molecule {
  id: 'molecule-rrsp-rule',
  type: 'ConditionMolecule',
  ruleType: 'if_then',
  condition: "jurisdiction equals Quebec",
  result: "withholding tax rate is 10%",
  configuration: {
    conditionType: "jurisdiction_equals",
    jurisdictionValue: "Quebec",
    rateValue: 0.10,
    taxType: "withholding",
    subject: "RRSP"
  },
  enforcement: 'block',
  tenantId: 'tenant-xyz',
  domain: 'tax'
})

// Create relationships
CREATE (rrsp)-[:PART_OF {role: 'subject', order: 1}]->(rule)
CREATE (withholding)-[:PART_OF {role: 'tax_type', order: 2}]->(rule)
CREATE (quebec)-[:PART_OF {role: 'condition', order: 3}]->(rule)
CREATE (tenPercent)-[:PART_OF {role: 'result', order: 4}]->(rule)

// Link to document
CREATE (rrsp)-[:APPEARS_IN {
  spans: [{start: 29, end: 33}],
  blockId: 'block-main',
  confidence: 0.99
}]->(:Document {id: 'doc-tax-guide', title: 'RRSP Tax Guidelines'})

CREATE (quebec)-[:APPEARS_IN {
  spans: [{start: 47, end: 53}],
  blockId: 'block-main', 
  confidence: 0.99
}]->(:Document {id: 'doc-tax-guide', title: 'RRSP Tax Guidelines'})
```

#### **3.4 Key Query Patterns**

```cypher
// Find all molecules affecting a specific atom
MATCH (a:Atom {value: 'RRSP'})<-[:PART_OF]-(m:Molecule)
RETURN m.type, m.condition, m.result, m.enforcement

// Impact analysis: What depends on Quebec jurisdiction?
MATCH (quebec:Atom {value: 'Quebec'})<-[:PART_OF {role: 'condition'}]-(m:Molecule)
MATCH (m)-[:PART_OF {role: 'subject'}]->(subject:Atom)
RETURN subject.value AS affected_subject, m.result AS rule_result

// Find conflicting rules
MATCH (m1:Molecule)-[:CONFLICTS_WITH]->(m2:Molecule)
WHERE m1.domain = 'tax' AND m2.domain = 'tax'
RETURN m1.condition AS rule1, m2.condition AS rule2

// Document knowledge coverage
MATCH (d:Document {title: 'RRSP Tax Guidelines'})<-[:APPEARS_IN]-(a:Atom)
RETURN d.title, count(a) AS atom_count, collect(DISTINCT a.type) AS atom_types

// Multi-hop dependency traversal
MATCH path = (start:Atom {value: 'GDPR'})-[:REQUIRES*..3]->(dependent:Atom)
RETURN [node IN nodes(path) | node.value] AS dependency_chain
```

#### **3.5 Clarification: PostgreSQL vs Neo4j Relationship Split**

**Architect Question**: "When we say all relationships, we only mean the relationships we are graphing correct? We'll still be expressing certain relationships using traditional foreign keys or join tables in postgres correct?"

**Answer**: Correct. There's a clear distinction:

##### **PostgreSQL Relationships (Data Integrity & Simple Joins)**
```sql
-- FOREIGN KEY RELATIONSHIPS (PostgreSQL)
-- For data integrity and simple 1:1, 1:many joins

-- Hierarchical containment
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

##### **Neo4j Relationships (Knowledge Graph & Traversal)**
```cypher
-- GRAPH RELATIONSHIPS (Neo4j)  
-- For knowledge discovery, traversals, path finding

(:Atom)-[:PART_OF {role, order}]->(:Molecule)       -- Composition with semantics
(:Atom)-[:APPEARS_IN {spans, confidence}]->(:Document) -- Document appearance
(:Atom)-[:RELATED_TO {type, strength}]->(:Atom)     -- Semantic relationships
(:Molecule)-[:REQUIRES]->(:Atom)                    -- Dependency relationships
(:Molecule)-[:CONFLICTS_WITH]->(:Molecule)          -- Rule conflicts
(:Block)-[:PART_OF_DOCUMENT]->(:Document)           -- Document structure
(:Atom)-[:IN_DOMAIN]->(:Domain)                     -- Domain organization
```

**Purpose**: Graph traversals, impact analysis, multi-hop queries, relationship semantics.

##### **Decision Framework**
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

**Key Design Decisions**:
1. **Hybrid relationship model**: PostgreSQL for integrity, Neo4j for intelligence
2. **Rich relationship properties in Neo4j**: Spans, confidence, roles, order, timestamps
3. **Tenant isolation via properties**: `tenantId` on all nodes/relationships in both databases
4. **Bidirectional synchronization**: PostgreSQL ↔ Neo4j sync for entity consistency
5. **Domain organization**: Domain relationships in both systems (FK in PostgreSQL, graph in Neo4j)

**Architect Feedback**:
✓ **Approved**: "Yes this makes sense."

**Final Decision**: Neo4j graph model approved with clarified PostgreSQL vs Neo4j relationship split:
- **PostgreSQL**: Data integrity, foreign keys, simple joins, array membership
- **Neo4j**: Knowledge graph, multi-hop traversals, relationship semantics, impact analysis

**Open Question for Next Concept**: How to handle PostgreSQL ↔ Neo4j synchronization.

---

## Concept 4: Synchronization & Data Flow

### **Data Engineer Presentation**

**Overview**: With our hybrid PostgreSQL + Neo4j architecture, we need a robust synchronization system to keep both databases consistent. This concept covers synchronization patterns, data flow, consistency guarantees, and the sync service design.

#### **4.1 Synchronization Requirements**

**Entities to Sync**:
1. **Atoms**: PostgreSQL `atoms` ↔ Neo4j `Atom` nodes
2. **Molecule Metadata**: PostgreSQL `molecule_metadata` ↔ Neo4j `Molecule` nodes  
3. **Documents**: PostgreSQL `documents` ↔ Neo4j `Document` nodes
4. **Blocks**: PostgreSQL `blocks` ↔ Neo4j `Block` nodes
5. **Domains**: PostgreSQL `domains` ↔ Neo4j `Domain` nodes

**Sync Direction**: Bidirectional, but with clear source of truth per entity type:
- **PostgreSQL as source of truth**: Content, values, configurations
- **Neo4j as source of truth**: Relationships, graph structure

#### **4.2 Synchronization Patterns**

**Pattern 1: PostgreSQL → Neo4j (Primary Flow)**
```typescript
// When content changes in PostgreSQL, sync to Neo4j
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
// When relationships change in Neo4j, update PostgreSQL arrays
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

**Pattern 3: Dual-Write with Compensation**
```typescript
// For critical operations, write to both with compensation
async function dualWriteAtom(atom: Atom): Promise<void> {
  const transactionId = generateId();
  
  try {
    // 1. Write to PostgreSQL
    await postgresClient.query(
      'INSERT INTO atoms VALUES ($1, $2, ...)',
      [atom.id, atom.tenantId, ...]
    );
    
    // 2. Write to Neo4j
    await neo4jClient.writeTransaction(async tx => {
      await tx.run(`
        CREATE (a:Atom {
          id: $id,
          tenantId: $tenantId,
          type: $type,
          value: $value
        })
      `, atom);
    });
    
    // 3. Record successful sync
    await recordSyncCompletion(transactionId, 'success');
    
  } catch (error) {
    // 4. Compensation logic
    await compensateFailedWrite(transactionId, error);
    throw error;
  }
}
```

#### **4.3 Sync Service Architecture**

```typescript
// Core sync service components
class HybridSyncService {
  // PostgreSQL Change Data Capture (CDC)
  private postgresCDC: PostgresChangeDataCapture;
  
  // Neo4j change listeners
  private neo4jListeners: Neo4jChangeListeners;
  
  // Message queue for async processing
  private messageQueue: MessageQueue;
  
  // Sync state tracking
  private syncState: SyncStateRepository;
  
  async initialize(): Promise<void> {
    // 1. Set up PostgreSQL triggers/listeners
    await this.setupPostgresCDC();
    
    // 2. Set up Neo4j transaction event handlers
    await this.setupNeo4jListeners();
    
    // 3. Start sync workers
    await this.startSyncWorkers();
  }
  
  private async setupPostgresCDC(): Promise<void> {
    // Use PostgreSQL logical replication or triggers
    // Capture: INSERT/UPDATE/DELETE on atoms, molecules, documents, blocks
    // Publish to message queue for processing
  }
  
  private async setupNeo4jListeners(): Promise<void> {
    // Listen to Neo4j transaction events
    // Capture: Relationship CREATE/DELETE, node property updates
    // Transform to PostgreSQL array operations
  }
}
```

#### **4.4 Consistency Models**

**Strong Consistency Required**:
- **Atom values**: PostgreSQL is source of truth, Neo4j must match
- **Molecule configurations**: PostgreSQL is source of truth
- **Document/block content**: PostgreSQL is source of truth

**Eventual Consistency Acceptable**:
- **Relationship updates**: Neo4j changes can propagate to PostgreSQL arrays async
- **Graph structure**: Temporary inconsistencies ok for query routing

**Conflict Resolution**:
```typescript
// Conflict resolution strategy
enum ConflictResolution {
  POSTGRES_WINS = 'postgres_wins',      // For content/value conflicts
  NEO4J_WINS = 'neo4j_wins',            // For relationship conflicts
  MANUAL_RESOLUTION = 'manual_resolution', // For critical conflicts
  LAST_WRITE_WINS = 'last_write_wins'   // For non-critical metadata
}

// Conflict detection
async function detectConflicts(entityId: UUID): Promise<Conflict[]> {
  const postgresState = await getPostgresState(entityId);
  const neo4jState = await getNeo4jState(entityId);
  
  return compareStates(postgresState, neo4jState);
}
```

#### **4.5 Data Flow Examples**

**Example 1: Create Atom with Relationships**
```
1. User/AI extracts atom from content
2. INSERT INTO atoms (PostgreSQL) ✓
3. CDC captures change → message queue
4. Sync worker: CREATE (:Atom {...}) in Neo4j
5. Neo4j relationship events captured
6. Update molecule_metadata.atom_ids arrays
7. Consistency check completes
```

**Example 2: Update Molecule Configuration**
```
1. Expert edits molecule configuration
2. UPDATE molecule_metadata (PostgreSQL) ✓
3. CDC captures change → message queue  
4. Sync worker: UPDATE (:Molecule {...}) in Neo4j
5. Graph queries now see updated configuration
6. No relationship changes needed
```

**Example 3: Create Relationship in Neo4j**
```
1. AI identifies new relationship between atoms
2. CREATE (a)-[:RELATED_TO]->(b) in Neo4j ✓
3. Neo4j listener captures event → message queue
4. Sync worker: Update PostgreSQL arrays
5. PostgreSQL now reflects relationship
```

#### **4.6 Sync State Tracking**

```sql
-- PostgreSQL table for sync state tracking
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

-- Index for finding pending syncs
CREATE INDEX sync_state_pending_idx ON sync_state(tenant_id, sync_status) 
WHERE sync_status IN ('pending', 'failed');
```

#### **4.7 POC Implementation Approach**

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

**Key Design Decisions**:
1. **PostgreSQL as primary source**: Content and values originate here
2. **Async synchronization**: Eventual consistency acceptable for relationships
3. **Compensation transactions**: Handle partial failures
4. **Sync state tracking**: Auditable sync history
5. **POC progressive enhancement**: Start simple, add complexity as needed

**Architect Feedback**:
✓ **Approved**: "At this point we can iterate post POC. Let's go with your recommendation."

**Final Decision**: Synchronization approach approved. POC will implement PostgreSQL → Neo4j sync with triggers and simple sync service. Full bidirectional automation will be iterated post-POC.

---

## Concept 5: Repository Pattern for Persistence Abstraction

### **Data Engineer Presentation**

**Overview**: Instead of a complex query router, we'll use simple repository classes to abstract the persistence layer. Each entity type gets a repository that encapsulates database choice, making it easier for developers to reason about and use.

#### **5.1 Repository Pattern Design**

**Core Idea**: Each entity (Atom, Molecule, Document, Block) has a repository with clean methods. The repository implementation decides which database to use based on the operation type.

```typescript
// Base repository interface
interface Repository<T> {
  findById(id: UUID): Promise<T | null>;
  save(entity: T): Promise<T>;
  delete(id: UUID): Promise<void>;
}

// Atom repository - encapsulates database choice
class AtomRepository implements Repository<Atom> {
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
  
  // Private sync method
  private async syncToNeo4j(atom: Atom): Promise<void> {
    await this.neo4jClient.query(`
      MERGE (a:Atom {id: $id})
      SET a.value = $value,
          a.type = $type,
          a.tenantId = $tenantId,
          a.domain = $domain
    `, atom);
  }
}
```

#### **5.2 Repository Structure for All Entities**

```typescript
// Molecule repository
class MoleculeRepository implements Repository<Molecule> {
  private postgresClient: PostgresClient;
  private neo4jClient: Neo4jClient;
  
  async findById(id: UUID): Promise<Molecule | null> {
    // PostgreSQL for configuration
    return this.postgresClient.query(
      'SELECT * FROM molecule_metadata WHERE id = $1',
      [id]
    );
  }
  
  async getAtomsInMolecule(moleculeId: UUID): Promise<Atom[]> {
    // Get atom IDs from PostgreSQL array
    const molecule = await this.findById(moleculeId);
    if (!molecule) return [];
    
    // Get atom details
    return this.postgresClient.query(
      'SELECT * FROM atoms WHERE id = ANY($1)',
      [molecule.atom_ids]
    );
  }
  
  async impactAnalysis(moleculeId: UUID): Promise<ImpactResult> {
    // Neo4j for graph traversal
    return this.neo4jClient.query(`
      MATCH (m:Molecule {id: $moleculeId})
      MATCH (m)-[:REQUIRES|CONFLICTS_WITH*..5]-(impacted)
      RETURN impacted, count(*) as dependencyDepth
      ORDER BY dependencyDepth DESC
    `, { moleculeId });
  }
}

// Document repository
class DocumentRepository implements Repository<Document> {
  private postgresClient: PostgresClient;
  
  async findById(id: UUID): Promise<Document | null> {
    // PostgreSQL for document content
    return this.postgresClient.query(`
      SELECT d.*, json_agg(b.*) as blocks
      FROM documents d
      LEFT JOIN blocks b ON b.parent_id = d.id
      WHERE d.id = $1
      GROUP BY d.id
    `, [id]);
  }
  
  async getAtomsInDocument(documentId: UUID): Promise<Atom[]> {
    // Get atom IDs from blocks array
    const document = await this.findById(documentId);
    if (!document) return [];
    
    // Collect all atom IDs from blocks
    const atomIds = document.blocks.flatMap(block => block.atom_ids || []);
    
    return this.postgresClient.query(
      'SELECT * FROM atoms WHERE id = ANY($1)',
      [atomIds]
    );
  }
}

// Block repository
class BlockRepository implements Repository<Block> {
  private postgresClient: PostgresClient;
  
  async findById(id: UUID): Promise<Block | null> {
    return this.postgresClient.query(
      'SELECT * FROM blocks WHERE id = $1',
      [id]
    );
  }
  
  async getChildBlocks(parentId: UUID): Promise<Block[]> {
    return this.postgresClient.query(
      'SELECT * FROM blocks WHERE parent_id = $1 ORDER BY created_at',
      [parentId]
    );
  }
}
```

#### **5.3 Developer Usage (Simple & Clean)**

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

async function getDocumentWithContext(documentId: UUID) {
  // Get document content
  const document = await documentRepo.findById(documentId);
  
  // Get atoms in document
  const atoms = await documentRepo.getAtomsInDocument(documentId);
  
  // For each atom, get related concepts
  const atomContexts = await Promise.all(
    atoms.map(async atom => ({
      atom,
      related: await atomRepo.findRelatedConcepts(atom.id, 2)
    }))
  );
  
  return { document, atomContexts };
}

// Example: Complete workflow
async function handleUserQuery(userQuestion: string) {
  // Parse question (simplified)
  if (userQuestion.includes('related to')) {
    const concept = extractConcept(userQuestion);
    const atoms = await atomRepo.findByValue(concept);
    return await atomRepo.findRelatedConcepts(atoms[0].id);
  }
  
  if (userQuestion.includes('document')) {
    const docId = extractDocumentId(userQuestion);
    return await documentRepo.findById(docId);
  }
  
  if (userQuestion.includes('impact') || userQuestion.includes('affected if')) {
    const moleculeId = extractMoleculeId(userQuestion);
    return await moleculeRepo.impactAnalysis(moleculeId);
  }
}
```

#### **5.4 Benefits Over Complex Router**

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

**Progressive Enhancement**:
```typescript
// Start simple: Basic CRUD operations
class SimpleAtomRepository {
  async findById(id: UUID): Promise<Atom> {
    return postgres.query('SELECT * FROM atoms WHERE id = $1', [id]);
  }
}

// Add graph operations later
class EnhancedAtomRepository extends SimpleAtomRepository {
  async findRelatedConcepts(id: UUID): Promise<ConceptGraph> {
    return neo4j.query('MATCH (a:Atom {id: $id})-[:RELATED_TO*..2]-(related) RETURN related', { id });
  }
}
```

#### **5.5 Performance Considerations**

**Repository optimizations**:
```typescript
class CachedAtomRepository implements AtomRepository {
  private cache: Map<UUID, Atom> = new Map();
  
  async findById(id: UUID): Promise<Atom | null> {
    // Check cache first
    if (this.cache.has(id)) {
      return this.cache.get(id)!;
    }
    
    // Hit database
    const atom = await this.postgresClient.query(
      'SELECT * FROM atoms WHERE id = $1',
      [id]
    );
    
    // Cache result
    if (atom) this.cache.set(id, atom);
    
    return atom;
  }
}
```

**Database optimizations** (same as before, but encapsulated in repositories):
```sql
-- PostgreSQL indexes in migration files
CREATE INDEX atoms_value_idx ON atoms USING gin(to_tsvector('english', text_value));
CREATE INDEX atoms_tenant_idx ON atoms(tenant_id);

-- Neo4j indexes managed by repository initialization
// In repository constructor:
await neo4jClient.query('CREATE INDEX atom_id_idx IF NOT EXISTS FOR (a:Atom) ON (a.id)');
```

#### **5.6 Comparison: Repository vs Router**

| **Aspect** | **Complex Router** | **Simple Repositories** |
|------------|-------------------|-------------------------|
| **Developer experience** | Need to understand routing logic | Work with entity interfaces |
| **Testing** | Mock router + databases | Mock repositories |
| **Code organization** | Centralized routing logic | Distributed by entity |
| **Cognitive load** | High (which database?) | Low (just call repo methods) |
| **Flexibility** | Complex to extend | Easy to add new methods |
| **Performance** | Router overhead | Direct database calls |

**Key Insight**: Repositories match how developers already think - in terms of entities and operations on those entities, not database technologies.

**Architect Feedback**:
✓ **Approved**: "Yes, abstract the persistence query (postgres / neo4j) into abstractions but still allow for developers to create custom queries in a clear way."

**Final Decision**: Repository pattern approved. Each entity will have a repository that encapsulates database choice while allowing custom queries through clear abstractions.

---

## Concept 6: Multi-Tenancy Implementation

### **Data Engineer Presentation**

**Overview**: Knowledge OS needs to support multiple tenants (organizations/teams) with data isolation, tenant-specific configurations, and scalable resource management. This concept covers multi-tenancy approaches, tenant isolation strategies, and implementation details.

#### **6.1 Multi-Tenancy Approaches**

**Three Common Patterns**:
1. **Database per Tenant**: Separate database for each tenant
2. **Schema per Tenant**: Separate schema within shared database  
3. **Row-level Isolation**: Shared tables with `tenant_id` column

**Recommended Approach**: **Hybrid model** combining strengths:
- **PostgreSQL**: Schema-per-tenant for strong isolation
- **Neo4j**: Row-level isolation with `tenantId` property
- **Shared infrastructure**: For non-sensitive, scalable components

#### **6.2 PostgreSQL: Schema-per-Tenant Implementation**

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
    -- ... other columns from Concept 2
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

-- 4. Dynamic schema creation for new tenant
CREATE OR REPLACE FUNCTION create_tenant_schema(tenant_id UUID, tenant_name TEXT)
RETURNS VOID AS $$
DECLARE
    schema_name TEXT := 'tenant_' || replace(tenant_name, ' ', '_') || '_' || tenant_id;
BEGIN
    -- Create schema
    EXECUTE 'CREATE SCHEMA ' || quote_ident(schema_name);
    
    -- Copy template tables
    EXECUTE 'CREATE TABLE ' || quote_ident(schema_name) || '.blocks 
             (LIKE tenant_template.blocks INCLUDING ALL)';
    
    -- Update tenants table
    UPDATE tenants 
    SET schema_name = schema_name
    WHERE id = tenant_id;
END;
$$ LANGUAGE plpgsql;
```

**Benefits of Schema-per-Tenant**:
- **Strong isolation**: No accidental cross-tenant data leaks
- **Performance**: Tenant data co-located, better cache locality
- **Maintenance**: Easy tenant migration/archival (drop schema)
- **Security**: PostgreSQL role-based access control per schema

#### **6.3 Neo4j: Property-based Isolation**

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

#### **6.4 Repository Pattern with Multi-Tenancy**

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
    // Simple security wrapper - in production use more robust solution
    // Ensures all node patterns include {tenantId: $tenantId}
    return cypher.replace(/MATCH\s*\(([^:]+):/g, (match, varName) => {
      return `MATCH (${varName} {tenantId: $tenantId}:`;
    });
  }
}

// Atom repository with multi-tenancy
class TenantAtomRepository extends TenantAwareRepository<Atom> implements AtomRepository {
  async findById(id: UUID): Promise<Atom | null> {
    // PostgreSQL query in tenant schema
    return this.postgresQuery(
      'SELECT * FROM atoms WHERE id = $1',
      [id]
    );
  }
  
  async findRelatedConcepts(atomId: UUID, depth: number = 2): Promise<ConceptGraph> {
    // Neo4j query with tenant filter
    return this.neo4jQuery(`
      MATCH (start:Atom {id: $atomId})
      MATCH path = (start)-[:RELATED_TO*1..${depth}]-(related:Atom)
      RETURN nodes(path) as nodes, relationships(path) as relationships
    `, { atomId });
  }
}
```

#### **6.5 Tenant Context Management**

```typescript
// Tenant context middleware (Express/Next.js)
interface TenantContext {
  id: UUID;
  name: string;
  schemaName: string;
  planType: string;
  features: TenantFeatures;
}

class TenantContextManager {
  private context: Map<string, TenantContext> = new Map();
  
  // Extract tenant from request (subdomain, JWT, header, etc.)
  async resolveFromRequest(req: Request): Promise<TenantContext> {
    // 1. Try subdomain: acme.knowledgeos.app → 'acme'
    const hostname = req.headers.host || '';
    const subdomain = hostname.split('.')[0];
    
    // 2. Try X-Tenant-ID header
    const tenantHeader = req.headers['x-tenant-id'];
    
    // 3. Try JWT claim
    const authToken = req.headers.authorization?.replace('Bearer ', '');
    const tenantFromJwt = await this.extractTenantFromJwt(authToken);
    
    // 4. Lookup tenant in database
    const tenant = await this.lookupTenant(subdomain, tenantHeader, tenantFromJwt);
    
    // 5. Cache for request lifetime
    this.context.set(req.id, tenant);
    
    return tenant;
  }
  
  // Repository factory with tenant context
  createRepository<T>(req: Request, repoClass: new (tenantId: UUID, ...args: any[]) => T): T {
    const tenant = this.resolveFromRequest(req);
    const postgresClient = this.getPostgresClient();
    const neo4jClient = this.getNeo4jClient();
    
    return new repoClass(tenant.id, postgresClient, neo4jClient);
  }
}

// Usage in API endpoints
app.get('/api/documents/:id', async (req, res) => {
  const contextManager = new TenantContextManager();
  const documentRepo = contextManager.createRepository(req, TenantDocumentRepository);
  
  const document = await documentRepo.findById(req.params.id);
  res.json(document);
});
```

#### **6.6 Cross-Tenant Operations (Admin/Shared Knowledge)**

```typescript
// Admin repository (bypasses tenant isolation)
class AdminAtomRepository extends AtomRepository {
  private postgresClient: PostgresClient;
  private neo4jClient: Neo4jClient;
  
  async findAllTenants(): Promise<Tenant[]> {
    // Query public.tenants table
    return this.postgresClient.query('SELECT * FROM tenants');
  }
  
  async queryAcrossTenants(query: CrossTenantQuery): Promise<CrossTenantResult> {
    // Dynamic union across all tenant schemas
    const tenants = await this.findAllTenants();
    
    const results = await Promise.all(
      tenants.map(async tenant => {
        await this.postgresClient.query(`SET search_path TO ${tenant.schemaName}, public`);
        return this.postgresClient.query(query.sql, query.params);
      })
    );
    
    return { tenants, results };
  }
  
  async createSharedKnowledge(sharedAtom: SharedAtom): Promise<void> {
    // Create in all tenant schemas
    const tenants = await this.findAllTenants();
    
    for (const tenant of tenants) {
      await this.postgresClient.query(`SET search_path TO ${tenant.schemaName}, public`);
      await this.postgresClient.query(
        'INSERT INTO shared_atoms VALUES ($1, $2, $3)',
        [sharedAtom.id, sharedAtom.value, tenant.id]
      );
    }
    
    // Also create in Neo4j with multi-tenant flag
    await this.neo4jClient.query(`
      CREATE (a:SharedAtom {
        id: $id,
        value: $value,
        isShared: true,
        tenantIds: $tenantIds
      })
    `, { 
      id: sharedAtom.id, 
      value: sharedAtom.value,
      tenantIds: tenants.map(t => t.id)
    });
  }
}
```

#### **6.7 Migration & Data Isolation Strategy**

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

#### **6.8 Security & Compliance**

**Data Isolation Guarantees**:
1. **SQL Injection protection**: Tenant schema isolation prevents cross-tenant leaks
2. **Cypher injection protection**: Query rewriting ensures tenantId filter
3. **Audit logging**: All queries log tenant context
4. **Data residency**: Tenant schemas can be hosted in region-specific databases

**Compliance Features**:
```typescript
interface ComplianceConfig {
  gdprEnabled: boolean;
  dataRetentionDays: number;
  allowCrossTenantAnalytics: boolean;
  requireDataResidency: string[]; // ['EU', 'US', 'CA']
}

class ComplianceManager {
  async enforceDataRetention(tenantId: UUID): Promise<void> {
    const config = await this.getComplianceConfig(tenantId);
    
    if (config.gdprEnabled) {
      // Soft delete with deleted_at timestamp
      await this.postgresQuery(
        `UPDATE atoms SET deleted_at = NOW() 
         WHERE tenant_id = $1 AND created_at < NOW() - INTERVAL '${config.dataRetentionDays} days'`,
        [tenantId]
      );
    }
  }
  
  async handleDataSubjectRequest(tenantId: UUID, userId: UUID, requestType: 'access' | 'deletion'): Promise<void> {
    if (requestType === 'deletion') {
      // Anonymize or delete user data
      await this.postgresQuery(
        `UPDATE atoms SET text_value = '[REDACTED]', number_value = NULL
         WHERE tenant_id = $1 AND created_by = $2`,
        [tenantId, userId]
      );
    }
  }
}
```

#### **6.9 Performance Considerations**

**Connection Pooling**:
```typescript
class TenantConnectionPool {
  private pools: Map<UUID, Pool> = new Map();
  
  async getConnection(tenantId: UUID): Promise<PoolClient> {
    if (!this.pools.has(tenantId)) {
      // Create tenant-specific connection pool
      const pool = new Pool({
        database: `knowledgeos_${tenantId}`,
        // ... other config
      });
      this.pools.set(tenantId, pool);
    }
    
    return this.pools.get(tenantId)!.connect();
  }
}
```

**Caching with Tenant Isolation**:
```typescript
class TenantAwareCache {
  async get<T>(tenantId: UUID, key: string): Promise<T | null> {
    // Include tenantId in cache key
    const tenantKey = `tenant:${tenantId}:${key}`;
    return this.redis.get(tenantKey);
  }
  
  async set<T>(tenantId: UUID, key: string, value: T, ttl?: number): Promise<void> {
    const tenantKey = `tenant:${tenantId}:${key}`;
    await this.redis.set(tenantKey, JSON.stringify(value));
    if (ttl) await this.redis.expire(tenantKey, ttl);
  }
}
```

**Key Design Decisions**:
1. **Hybrid isolation**: PostgreSQL schema-per-tenant + Neo4j property-based
2. **Progressive migration**: Row-level → Schema-per-tenant → Dedicated database
3. **Repository pattern integration**: Tenant context baked into repositories
4. **Security first**: Automatic query rewriting for tenant isolation
5. **Compliance ready**: GDPR, data retention, subject requests

**Architect Feedback Requested**:
1. Does the hybrid isolation approach (PostgreSQL schema-per-tenant + Neo4j property-based) make sense?
2. Are the progressive migration phases appropriate?
3. Any concerns with the repository pattern integration for multi-tenancy?
4. Should we consider different approaches for enterprise-scale tenants?

---

## Next Concept (Pending Architect Review)