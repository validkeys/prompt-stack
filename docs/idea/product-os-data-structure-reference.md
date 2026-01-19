# Product-OS: Data Structure Reference

**Date**: 2026-01-11  
**Status**: Technical Reference - Approved Architecture  
**Related**: [`product-os-data-structure.md`](product-os-data-structure.md) (Concept document), [`product-os-implementation-persistence-operations.md`](product-os-implementation-persistence-operations.md) (Operational guidance)  
**Source**: Consolidated from [`data-structure-walkthrough.md`](archive/data-structure-walkthrough.md) (Iterative design process)

---

## PostgreSQL Schema Reference

### 1. `blocks` - Notion-style content units
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

### 2. `documents` - Collections of blocks
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

### 3. `atoms` - Extracted knowledge values
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

### 4. `molecule_metadata` - Business logic configuration
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

### 5. `domains` - Bounded contexts
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

### 6. `sync_state` - Synchronization tracking
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

**Purpose**: Tracks synchronization state between PostgreSQL and Neo4j for eventual consistency.

---

## Neo4j Graph Model Reference

### Node Types & Properties
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

### Relationship Types & Properties
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

### Example: RRSP Tax Rule in Neo4j
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
```

### Key Query Patterns
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

---

## Repository Pattern Implementation

### Atom Repository
```typescript
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

### Molecule Repository
```typescript
class MoleculeRepository {
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
```

### Document Repository
```typescript
class DocumentRepository {
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
```

### Developer Usage Example
```typescript
import { AtomRepository, MoleculeRepository, DocumentRepository } from '@knowledge-os/repositories';

// Initialize (dependency injection)
const atomRepo = new AtomRepository(postgresClient, neo4jClient);
const moleculeRepo = new MoleculeRepository(postgresClient, neo4jClient);
const documentRepo = new DocumentRepository(postgresClient);

// Use repositories naturally - no database decisions needed
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
```

---

## Multi-Tenancy Implementation

### PostgreSQL: Schema-per-Tenant
```sql
-- Tenant management table (in public schema)
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

-- Dynamic schema creation for new tenant
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

### Tenant-Aware Repository Base Class
```typescript
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

---

## Synchronization Patterns

### PostgreSQL → Neo4j Sync (Primary Flow)
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

### Neo4j → PostgreSQL Sync (Relationship Updates)
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

---

## Key Performance Indexes

### PostgreSQL Indexes
```sql
-- Blocks table
CREATE INDEX blocks_parent_idx ON blocks(parent_id, parent_type);
CREATE INDEX blocks_tenant_idx ON blocks(tenant_id);
CREATE INDEX blocks_atom_refs_idx ON blocks USING gin(atom_ids);
CREATE INDEX blocks_molecule_refs_idx ON blocks USING gin(molecule_ids);

-- Atoms table
CREATE INDEX atoms_tenant_idx ON atoms(tenant_id);
CREATE INDEX atoms_domain_idx ON atoms(domain_id);
CREATE INDEX atoms_value_search_idx ON atoms USING gin(to_tsvector('english', text_value));
CREATE INDEX atoms_type_idx ON atoms(atom_type);

-- Molecules table
CREATE INDEX molecules_tenant_idx ON molecule_metadata(tenant_id);
CREATE INDEX molecules_atom_refs_idx ON molecule_metadata USING gin(atom_ids);
CREATE INDEX molecules_type_idx ON molecule_metadata(molecule_type);

-- Documents table
CREATE INDEX documents_tenant_idx ON documents(tenant_id);
CREATE INDEX documents_domain_idx ON documents(domain_id);
CREATE INDEX documents_title_idx ON documents USING gin(to_tsvector('english', title));
```

### Neo4j Indexes
```cypher
// Tenant isolation indexes
CREATE INDEX tenant_atom_idx FOR (a:Atom) ON (a.tenantId, a.id);
CREATE INDEX tenant_molecule_idx FOR (m:Molecule) ON (m.tenantId, m.id);
CREATE INDEX tenant_document_idx FOR (d:Document) ON (d.tenantId, d.id);

// Query optimization indexes
CREATE INDEX atom_value_idx FOR (a:Atom) ON (a.value);
CREATE INDEX molecule_rule_type_idx FOR (m:Molecule) ON (m.ruleType);
CREATE INDEX document_title_idx FOR (d:Document) ON (d.title);

// Relationship traversal indexes
CREATE INDEX relationship_tenant_idx FOR ()-[r:PART_OF|APPEARS_IN|RELATED_TO|REQUIRES]-() ON (r.tenantId);
```

---

## Migration Strategy

### Phase 1: POC (Row-level isolation)
```typescript
class POCAtomRepository {
  async findById(tenantId: UUID, atomId: UUID): Promise<Atom> {
    return this.postgresClient.query(
      'SELECT * FROM atoms WHERE tenant_id = $1 AND id = $2',
      [tenantId, atomId]
    );
  }
}
```

### Phase 2: Schema-per-Tenant Migration
```typescript
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

---

**Status**: Technical reference - ready for implementation  
**Next**: Begin POC with PostgreSQL schema + basic Neo4j sync  
**Concept**: See [`product-os-data-structure.md`](product-os-data-structure.md) for architectural rationale