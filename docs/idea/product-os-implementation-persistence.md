# Product-OS: Polyglot Database Architecture

**Date**: 2026-01-11  
**Status**: Technical Reference Document  
**Version**: 1.0  
**Related**: [`product-os.md`](product-os.md) (Knowledge OS concept), [`product-os-data-structure.md`](product-os-data-structure.md) (Data structure concept), [`product-os-implementation-rendering.md`](product-os-implementation-rendering.md) (Standoff annotations), [`product-os-implementation-search.md`](product-os-implementation-search.md) (Multi-level search), [`product-os-implementation-extraction.md`](product-os-implementation-extraction.md) (Extraction pipeline), [`product-os-implementation-frontend.md`](product-os-implementation-frontend.md) (Frontend patterns)

---

## Overview

**Polyglot persistence** is the strategic use of multiple database technologies optimized for different data access patterns in Knowledge OS. This document covers the hybrid architecture combining PostgreSQL, Neo4j, Redis, and object storage to handle atomic knowledge, relationships, caching, and documents.

**Key Principles**:
- **Right tool for the job**: Match database technology to data characteristics
- **Tiered isolation**: Row-level → schema-per-tenant → database-per-tenant scaling
- **Eventual consistency**: Async synchronization between specialized stores
- **Zero-downtime migrations**: Live schema evolution and data redistribution

---

## Database Technology Selection

### 1. PostgreSQL (Primary Knowledge Store)
**Use case**: Atomic knowledge units, documents, transactions, vector search
```sql
-- Core tables for atomic knowledge
CREATE TABLE atoms (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  type VARCHAR(50) NOT NULL, -- 'TextAtom', 'NumberAtom', etc.
  value JSONB NOT NULL,
  metadata JSONB,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  -- Vector embedding for semantic search
  embedding VECTOR(1536)
);

CREATE TABLE molecules (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  type VARCHAR(50) NOT NULL, -- 'ConditionMolecule', 'ConstraintMolecule'
  atoms UUID[] NOT NULL, -- References to atom IDs
  structure JSONB NOT NULL,
  metadata JSONB
);

-- Enable vector similarity search
CREATE INDEX atoms_embedding_idx ON atoms 
USING ivfflat (embedding vector_cosine_ops);
```

**PostgreSQL Extensions**:
- **pgvector**: Vector similarity for semantic search
- **PostGIS**: Geographic atoms (location-based knowledge)
- **pg_trgm**: Fuzzy text matching for atom values
- **TimescaleDB**: Time-series atoms (metrics, monitoring)

### 2. Neo4j (Relationship Engine)
**Use case**: Knowledge graph relationships, dependency tracking, impact analysis
```cypher
// Graph model for atomic knowledge relationships
CREATE (a:Atom {
  id: 'atom-123',
  type: 'NumberAtom',
  value: 0.05,
  tenantId: 'tenant-xyz'
})

CREATE (m:Molecule {
  id: 'molecule-456',
  type: 'ConditionMolecule',
  tenantId: 'tenant-xyz'
})

CREATE (a)-[:PART_OF]->(m)
CREATE (m)-[:REQUIRES]->(other:Atom {id: 'atom-789'})
CREATE (m)-[:IMPLEMENTS]->(c:Code {file: 'tax_calculator.go', line: 45})
```

**Graph Patterns**:
- **Hierarchical traversal**: Atoms → Molecules → Organisms
- **Impact analysis**: Multi-hop relationship tracing
- **Circular dependency detection**: Graph algorithms for validation
- **Community detection**: Clustering related knowledge units

### 3. Redis (Caching & Real-time)
**Use case**: Annotation caching, session state, real-time collaboration
```redis
# Key patterns for atomic knowledge caching
SET "tenant:xyz:atom:atom-123" "{'value': 0.05, 'type': 'NumberAtom'}" EX 3600
HSET "tenant:xyz:document:doc-abc:annotations" "layer:atom" "[...]" 
ZADD "tenant:xyz:recent:atoms" 1641939200 "atom-123"
PUBLISH "tenant:xyz:updates" "{'type': 'atom_updated', 'id': 'atom-123'}"
```

**Redis Use Cases**:
- **Annotation position cache**: Fast document rendering
- **Real-time updates**: Pub/sub for collaborative editing
- **Rate limiting**: API request throttling per tenant
- **Session storage**: User preferences, persona configurations

### 4. Object Storage (Documents & Media)
**Use case**: Original documents, images, PDFs, versioned backups
```yaml
# S3-compatible storage structure
bucket: knowledge-os-tenant-xyz
├── documents/
│   ├── tax/
│   │   ├── canadian-gst.md
│   │   └── withholding-tax.pdf
│   └── compliance/
│       └── gdpr-requirements.docx
├── media/
│   ├── diagrams/
│   └── screenshots/
└── backups/
    ├── daily/
    └── weekly/
```

**Object Storage Features**:
- **Versioning**: Full history of document changes
- **Lifecycle policies**: Automatic archival to cold storage
- **Encryption**: At-rest and in-transit security
- **CDN integration**: Global delivery for frequently accessed documents

---

## Multi-Tenancy Architecture

### Tier 1: Row-level Isolation (Small Teams)
```sql
-- Single database with tenant_id filtering
SELECT * FROM atoms 
WHERE tenant_id = 'tenant-xyz' 
  AND type = 'NumberAtom'
  AND value->>'rate' > '0.1';

-- Row-level security policies
CREATE POLICY tenant_isolation_policy ON atoms
  USING (tenant_id = current_tenant_id());
```

**Pros**:
- Simple deployment
- Shared connection pooling
- Easy cross-tenant analytics (for platform admins)

**Cons**:
- Noisy neighbor risk
- Limited scaling (≈1k nodes per tenant)

### Tier 2: Schema-per-Tenant (Medium Organizations)
```sql
-- Dynamic schema switching
SET search_path TO tenant_xyz, public;

-- Per-tenant schema with identical structure
CREATE SCHEMA IF NOT EXISTS tenant_xyz;
CREATE TABLE tenant_xyz.atoms (...);
CREATE TABLE tenant_xyz.molecules (...);
```

**Pros**:
- Strong isolation
- Tenant-specific optimizations
- Easy tenant migration

**Cons**:
- Connection overhead
- Schema management complexity

### Tier 3: Database-per-Tenant (Large Enterprises)
```yaml
# Tenant database configuration
databases:
  tenant_xyz:
    host: db-tenant-xyz.cluster-abc.us-east-1.rds.amazonaws.com
    port: 5432
    database: knowledge_os_tenant_xyz
    pool_size: 20
    
  tenant_abc:
    host: db-tenant-abc.cluster-def.us-west-2.rds.amazonaws.com
    port: 5432
    database: knowledge_os_tenant_abc
    pool_size: 50
```

**Pros**:
- Maximum isolation
- Independent scaling
- Custom backup/retention policies

**Cons**:
- High operational overhead
- Cross-tenant queries impossible

### Automatic Tier Promotion
```typescript
// Monitor tenant size and promote tiers
class TenantTierManager {
  async evaluateTenant(tenantId: string) {
    const stats = await this.getTenantStats(tenantId);
    
    if (stats.atomCount > 10000 && stats.currentTier === 'row-level') {
      await this.migrateToSchemaPerTenant(tenantId);
    }
    
    if (stats.atomCount > 100000 && stats.currentTier === 'schema-per-tenant') {
      await this.migrateToDatabasePerTenant(tenantId);
    }
  }
}
```

---

## Data Models

### 1. Atom Storage Model
```sql
-- Polymorphic atom storage
CREATE TABLE atoms (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  atom_type VARCHAR(50) NOT NULL,
  
  -- Type-specific value storage
  text_value TEXT,
  number_value NUMERIC,
  boolean_value BOOLEAN,
  date_value TIMESTAMP WITH TIME ZONE,
  json_value JSONB,
  
  -- Common metadata
  name VARCHAR(255),
  description TEXT,
  confidence FLOAT DEFAULT 1.0,
  source VARCHAR(50), -- 'manual', 'ai_extracted', 'code_generated'
  
  -- Versioning
  version INTEGER DEFAULT 1,
  previous_version UUID REFERENCES atoms(id),
  
  -- Timestamps
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE
);

-- Efficient type-specific queries
CREATE INDEX atoms_type_number_idx ON atoms(tenant_id, atom_type) 
WHERE atom_type = 'NumberAtom' AND number_value IS NOT NULL;

CREATE INDEX atoms_text_search_idx ON atoms 
USING gin(to_tsvector('english', text_value));
```

### 2. Molecule Storage Model
```sql
-- Molecules as atom compositions
CREATE TABLE molecules (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  molecule_type VARCHAR(50) NOT NULL,
  
  -- Graph structure
  root_atom_id UUID REFERENCES atoms(id),
  atom_ids UUID[] NOT NULL,
  
  -- Type-specific configuration
  configuration JSONB,
  
  -- Validation rules
  validation_rules JSONB,
  enforcement_level VARCHAR(20), -- 'block', 'warn', 'inform'
  
  -- Relationships (denormalized for performance)
  requires_molecules UUID[],
  conflicts_with_molecules UUID[],
  implements_code_refs TEXT[],
  
  -- Timestamps
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Efficient atom membership queries
CREATE INDEX molecules_atom_ids_gin_idx ON molecules 
USING gin(atom_ids);
```

### 3. Document Storage Model
```sql
-- Documents with standoff annotations
CREATE TABLE documents (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  
  -- Content storage
  content TEXT NOT NULL,
  content_hash CHAR(64) NOT NULL, -- SHA-256 for deduplication
  format VARCHAR(20) NOT NULL, -- 'markdown', 'html', 'pdf', 'code'
  
  -- Metadata
  title VARCHAR(500),
  path VARCHAR(1000) NOT NULL, -- Virtual filesystem path
  mime_type VARCHAR(100),
  
  -- Versioning
  version INTEGER DEFAULT 1,
  previous_version_id UUID REFERENCES documents(id),
  
  -- Storage location
  storage_backend VARCHAR(20) DEFAULT 'database', -- 'database', 's3', 'git'
  storage_reference TEXT, -- External reference if stored elsewhere
  
  -- Timestamps
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  
  -- Full-text search
  tsvector_column TSVECTOR GENERATED ALWAYS AS (
    to_tsvector('english', COALESCE(title, '') || ' ' || content)
  ) STORED
);

CREATE INDEX documents_path_idx ON documents(tenant_id, path);
CREATE INDEX documents_tsvector_idx ON documents USING gin(tsvector_column);
```

### 4. Annotation Storage Model
```sql
-- Standoff annotations linking documents to knowledge
CREATE TABLE annotations (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  document_id UUID NOT NULL REFERENCES documents(id),
  
  -- Target position in document
  target_type VARCHAR(20) NOT NULL, -- 'character', 'line', 'element'
  start_offset INTEGER NOT NULL,
  end_offset INTEGER NOT NULL,
  css_selector TEXT,
  xpath TEXT,
  
  -- Knowledge reference
  knowledge_type VARCHAR(20) NOT NULL, -- 'atom', 'molecule', 'organism'
  knowledge_id UUID NOT NULL, -- Reference to atoms/molecules table
  
  -- Annotation metadata
  layer VARCHAR(20) DEFAULT 'atom',
  confidence FLOAT DEFAULT 1.0,
  source VARCHAR(50) DEFAULT 'manual', -- 'manual', 'ai', 'code_analysis'
  
  -- Visual styling
  style JSONB,
  
  -- Validation
  validated_by UUID, -- User who validated this annotation
  validated_at TIMESTAMP WITH TIME ZONE,
  
  -- Timestamps
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Efficient document annotation queries
CREATE INDEX annotations_document_idx ON annotations(tenant_id, document_id);
CREATE INDEX annotations_knowledge_idx ON annotations(tenant_id, knowledge_type, knowledge_id);
CREATE INDEX annotations_position_idx ON annotations(
  tenant_id, 
  document_id, 
  start_offset, 
  end_offset
);
```

---

## Synchronization Patterns

### 1. PostgreSQL ↔ Neo4j Sync
```typescript
// Change Data Capture (CDC) for graph synchronization
class GraphSyncService {
  async syncToNeo4j(change: DatabaseChange) {
    // Capture changes from PostgreSQL
    const changes = await this.captureChanges(change);
    
    // Transform to graph operations
    const graphOps = this.transformToGraphOperations(changes);
    
    // Apply to Neo4j with retry logic
    await this.applyGraphOperations(graphOps);
    
    // Update sync metadata
    await this.recordSyncCompletion(change);
  }
  
  async syncFromNeo4j(change: GraphChange) {
    // Reverse sync for graph-first changes
    const dbOps = this.transformToDatabaseOperations(change);
    await this.applyDatabaseOperations(dbOps);
  }
}
```

**Sync Strategies**:
- **Event-driven**: PostgreSQL triggers → message queue → Neo4j worker
- **Batch processing**: Periodic sync of changed records
- **Dual-write**: Write to both databases in transaction (with compensation)
- **Read-your-writes**: Cache layer for immediate consistency

### 2. Cache Invalidation Strategy
```typescript
// Cache coherence across polyglot stores
class CacheCoordinationService {
  async invalidateCaches(change: KnowledgeChange) {
    // Invalidate Redis caches
    await redis.del(`atom:${change.atomId}`);
    await redis.del(`document:${change.documentId}:annotations`);
    
    // Publish invalidation event
    await redis.publish('cache-invalidation', JSON.stringify({
      type: 'atom_updated',
      id: change.atomId,
      tenantId: change.tenantId
    }));
    
    // Update CDN cache headers
    await cdn.purge(`/api/atoms/${change.atomId}`);
  }
}
```

### 3. Event Sourcing for Audit Trail
```sql
-- Event store for all knowledge changes
CREATE TABLE knowledge_events (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  event_type VARCHAR(50) NOT NULL,
  aggregate_type VARCHAR(50) NOT NULL, -- 'atom', 'molecule', 'document'
  aggregate_id UUID NOT NULL,
  
  -- Event data
  event_data JSONB NOT NULL,
  metadata JSONB,
  
  -- Causality tracking
  causation_id UUID, -- What caused this event
  correlation_id UUID, -- Business process correlation
  
  -- User context
  user_id UUID,
  user_ip INET,
  
  -- Timestamps
  occurred_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Project current state from event stream
CREATE VIEW atoms_current AS
WITH ranked_events AS (
  SELECT *,
    ROW_NUMBER() OVER (
      PARTITION BY aggregate_id 
      ORDER BY occurred_at DESC
    ) as rn
  FROM knowledge_events
  WHERE aggregate_type = 'atom'
    AND deleted_at IS NULL
)
SELECT * FROM ranked_events WHERE rn = 1;
```

---

## Migration Strategies

### 1. Zero-Downtime Schema Migration
```sql
-- Example: Adding vector embeddings to atoms
-- Phase 1: Add nullable column
ALTER TABLE atoms ADD COLUMN embedding VECTOR(1536);

-- Phase 2: Backfill in batches
UPDATE atoms 
SET embedding = compute_embedding(value)
WHERE id IN (
  SELECT id FROM atoms 
  WHERE embedding IS NULL 
  LIMIT 1000
)
RETURNING id;

-- Phase 3: Make column NOT NULL after backfill complete
ALTER TABLE atoms ALTER COLUMN embedding SET NOT NULL;

-- Phase 4: Create index concurrently
CREATE INDEX CONCURRENTLY atoms_embedding_idx 
ON atoms USING ivfflat (embedding vector_cosine_ops);
```

### 2. Tenant Tier Migration
```typescript
// Live migration from row-level to schema-per-tenant
class TenantMigrationService {
  async migrateToSchemaPerTenant(tenantId: string) {
    // 1. Create new schema
    await this.createTenantSchema(tenantId);
    
    // 2. Copy data using logical replication
    await this.setupLogicalReplication(tenantId);
    
    // 3. Dual-write during migration window
    await this.enableDualWrites(tenantId);
    
    // 4. Switch reads to new schema
    await this.switchReadTraffic(tenantId);
    
    // 5. Disable old writes
    await this.disableOldWrites(tenantId);
    
    // 6. Clean up old data
    await this.cleanupOldData(tenantId);
  }
}
```

### 3. Database Technology Migration
```typescript
// Migrating from MongoDB to PostgreSQL
class DatabaseMigrationService {
  async migrateFromMongoDB() {
    // 1. Schema analysis
    const mongoSchema = await this.analyzeMongoSchema();
    const pgSchema = this.designPostgreSQLSchema(mongoSchema);
    
    // 2. Parallel read/write
    await this.enableDualWrites();
    
    // 3. Historical data migration
    await this.migrateHistoricalData();
    
    // 4. Cutover
    await this.switchReadsToPostgreSQL();
    await this.disableMongoWrites();
    
    // 5. Validation
    await this.validateMigration();
  }
}
```

---

## Backup & Disaster Recovery

### Multi-Region Deployment
```yaml
# Terraform configuration for multi-region persistence
resources:
  # Primary region (us-east-1)
  - postgresql_cluster:
      region: us-east-1
      replicas: 3
      backup_retention: 30
  
  # Secondary region (eu-west-1)  
  - postgresql_replica:
      region: eu-west-1
      async_replication: true
      read_only: true
  
  # Neo4j causal cluster
  - neo4j_cluster:
      core_servers:
        - us-east-1a
        - us-east-1b
        - eu-west-1a
      read_replicas: 2
  
  # Cross-region Redis
  - redis_global_datastore:
      primary_region: us-east-1
      secondary_region: eu-west-1
      automatic_failover: true
```

### Backup Strategy
```bash
# Backup script for polyglot persistence
#!/bin/bash

# PostgreSQL logical backup
pg_dump --format=directory \
  --jobs=4 \
  --file=/backups/postgres/$(date +%Y%m%d) \
  knowledge_os

# Neo4j offline backup
neo4j-admin backup --backup-dir=/backups/neo4j/$(date +%Y%m%d) \
  --name=knowledge_os \
  --check-consistency

# Redis RDB snapshot
redis-cli --rdb /backups/redis/dump.rdb

# S3 document sync
aws s3 sync s3://knowledge-os-documents \
  /backups/s3/documents/$(date +%Y%m%d)
```

### Point-in-Time Recovery
```sql
-- Recover knowledge state as of specific time
WITH recovery_point AS (
  SELECT '2026-01-10 14:30:00 UTC'::timestamptz AS recovery_time
)
SELECT 
  e.event_type,
  e.event_data,
  e.occurred_at
FROM knowledge_events e
CROSS JOIN recovery_point r
WHERE e.occurred_at <= r.recovery_time
  AND e.aggregate_id = 'atom-123'
ORDER BY e.occurred_at;

-- Reconstruct atom state
CREATE OR REPLACE FUNCTION reconstruct_atom_state(
  atom_id UUID, 
  as_of TIMESTAMPTZ
) RETURNS JSONB AS $$
DECLARE
  result JSONB;
BEGIN
  WITH events AS (
    SELECT event_data
    FROM knowledge_events
    WHERE aggregate_id = atom_id
      AND occurred_at <= as_of
    ORDER BY occurred_at
  )
  SELECT jsonb_object_agg(
    key, 
    CASE 
      WHEN event_type = 'atom_deleted' THEN NULL
      ELSE value
    END
  ) INTO result
  FROM events, jsonb_each(event_data);
  
  RETURN result;
END;
$$ LANGUAGE plpgsql;
```

---

## Performance Optimization

### 1. Query Optimization Patterns
```sql
-- Materialized views for expensive aggregations
CREATE MATERIALIZED VIEW atom_statistics AS
SELECT 
  tenant_id,
  atom_type,
  COUNT(*) as count,
  MIN(created_at) as first_created,
  MAX(updated_at) as last_updated
FROM atoms
WHERE deleted_at IS NULL
GROUP BY tenant_id, atom_type;

-- Refresh on schedule
REFRESH MATERIALIZED VIEW CONCURRENTLY atom_statistics;

-- Partial indexes for common filters
CREATE INDEX atoms_active_number_idx ON atoms(tenant_id, number_value)
WHERE atom_type = 'NumberAtom' 
  AND deleted_at IS NULL
  AND confidence > 0.8;
```

### 2. Connection Pooling
```yaml
# PgBouncer configuration for PostgreSQL connection pooling
[databases]
knowledge_os = host=localhost port=5432 dbname=knowledge_os

[pgbouncer]
pool_mode = transaction
max_client_conn = 1000
default_pool_size = 20
reserve_pool_size = 5
```

### 3. Read Replicas for Scaling
```typescript
// Read/write splitting in application layer
class DatabaseRouter {
  getWriteConnection(tenantId: string): Connection {
    // Always route writes to primary
    return this.primaryConnection;
  }
  
  getReadConnection(tenantId: string, query: string): Connection {
    // Route based on query type and tenant tier
    if (this.isAnalyticalQuery(query)) {
      return this.analyticalReplica;
    }
    
    if (this.tenantTier(tenantId) === 'enterprise') {
      return this.tenantSpecificReplica(tenantId);
    }
    
    return this.randomReplica();
  }
}
```

### 4. Sharding Strategy
```sql
-- Tenant-based sharding
CREATE TABLE tenant_shards (
  tenant_id UUID PRIMARY KEY,
  shard_id INTEGER NOT NULL,
  database_host VARCHAR(255) NOT NULL,
  database_name VARCHAR(255) NOT NULL
);

-- Application-level shard routing
CREATE OR REPLACE FUNCTION route_to_shard(tenant_id UUID)
RETURNS text AS $$
DECLARE
  shard_info record;
BEGIN
  SELECT database_host, database_name
  INTO shard_info
  FROM tenant_shards
  WHERE tenant_id = route_to_shard.tenant_id;
  
  RETURN format(
    'host=%s dbname=%s',
    shard_info.database_host,
    shard_info.database_name
  );
END;
$$ LANGUAGE plpgsql;
```

---

## Security & Compliance

### 1. Encryption at Rest
```sql
-- Transparent Data Encryption (TDE)
CREATE TABLE sensitive_atoms (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  
  -- Encrypted columns
  encrypted_value BYTEA NOT NULL,
  encryption_key_id VARCHAR(100) NOT NULL,
  
  -- Searchable encryption
  searchable_hash CHAR(64), -- Hash of value for equality checks
  range_prefix VARCHAR(10) -- For range queries on encrypted data
);

-- Key rotation
CREATE OR REPLACE FUNCTION rotate_encryption_keys()
RETURNS void AS $$
BEGIN
  -- Generate new key
  INSERT INTO encryption_keys (key_id, key_data, created_at)
  VALUES (gen_random_uuid(), gen_random_bytes(32), NOW());
  
  -- Re-encrypt data with new key in batches
  -- ... implementation ...
END;
$$ LANGUAGE plpgsql;
```

### 2. Row-Level Security
```sql
-- PostgreSQL RLS policies
ALTER TABLE atoms ENABLE ROW LEVEL SECURITY;

-- Tenant isolation policy
CREATE POLICY tenant_isolation ON atoms
  USING (tenant_id = current_tenant_id());

-- Persona-based access control
CREATE POLICY persona_access ON atoms
  USING (
    -- Domain experts see all atoms
    current_persona() = 'domain_expert'
    OR
    -- Developers only see implemented atoms
    (current_persona() = 'developer' AND EXISTS (
      SELECT 1 FROM implementations 
      WHERE atom_id = atoms.id
    ))
    OR
    -- Product owners see business-relevant atoms
    (current_persona() = 'product_owner' AND EXISTS (
      SELECT 1 FROM capability_mappings 
      WHERE atom_id = atoms.id
    ))
  );
```

### 3. Audit Logging
```sql
-- Comprehensive audit trail
CREATE TABLE audit_log (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL,
  user_id UUID,
  action VARCHAR(50) NOT NULL,
  resource_type VARCHAR(50) NOT NULL,
  resource_id UUID NOT NULL,
  
  -- Before/after state
  old_values JSONB,
  new_values JSONB,
  
  -- Context
  ip_address INET,
  user_agent TEXT,
  request_id UUID,
  
  -- Timestamps
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Trigger-based audit logging
CREATE OR REPLACE FUNCTION audit_atom_changes()
RETURNS trigger AS $$
BEGIN
  IF (TG_OP = 'UPDATE') THEN
    INSERT INTO audit_log (
      tenant_id, user_id, action, 
      resource_type, resource_id,
      old_values, new_values
    )
    VALUES (
      NEW.tenant_id,
      current_user_id(),
      'UPDATE',
      'atom',
      NEW.id,
      to_jsonb(OLD),
      to_jsonb(NEW)
    );
  END IF;
  
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER atoms_audit_trigger
  AFTER UPDATE ON atoms
  FOR EACH ROW
  EXECUTE FUNCTION audit_atom_changes();
```

---

## Monitoring & Observability

### 1. Database Metrics
```yaml
# Prometheus metrics for polyglot persistence
metrics:
  postgresql:
    - queries_per_second
    - replication_lag
    - cache_hit_ratio
    - dead_tuples
    - connection_count
  
  neo4j:
    - page_cache_hits
    - transaction_rate
    - query_duration
    - heap_usage
    
  redis:
    - hit_rate
    - memory_usage
    - connected_clients
    - evicted_keys
    
  s3:
    - request_latency
    - error_rate
    - bandwidth
    - storage_used
```

### 2. Alerting Rules
```yaml
# Alertmanager configuration
groups:
  - name: database-alerts
    rules:
      - alert: HighReplicationLag
        expr: pg_replication_lag_seconds > 30
        for: 5m
        annotations:
          summary: "PostgreSQL replication lag >30 seconds"
          
      - alert: RedisMemoryHigh
        expr: redis_memory_usage_percent > 85
        for: 10m
        annotations:
          summary: "Redis memory usage >85%"
          
      - alert: Neo4jQuerySlow
        expr: rate(neo4j_query_duration_seconds_sum[5m]) / rate(neo4j_query_duration_seconds_count[5m]) > 1
        for: 5m
        annotations:
          summary: "Neo4j query duration >1 second"
```

### 3. Performance Tracing
```typescript
// Distributed tracing for database operations
class TracedDatabaseClient {
  async query(sql: string, params: any[]) {
    const span = tracer.startSpan('database.query');
    span.setTag('db.statement', sql);
    span.setTag('db.params', JSON.stringify(params));
    
    try {
      const result = await this.client.query(sql, params);
      span.setTag('db.row_count', result.rowCount);
      return result;
    } catch (error) {
      span.setTag('error', true);
      span.setTag('error.message', error.message);
      throw error;
    } finally {
      span.finish();
    }
  }
}
```

---

## Future Evolution

### 1. Distributed SQL (CockroachDB, Yugabyte)
```sql
-- Future migration to distributed SQL
-- Current: PostgreSQL single node
-- Future: Distributed SQL with global consistency

CREATE TABLE atoms (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  -- ... existing columns ...
) LOCALITY REGIONAL BY ROW; -- Data pinned to tenant's region

-- Global secondary indexes for cross-region queries
CREATE INDEX atoms_global_type_idx ON atoms(atom_type)
  STORING (value, metadata)
  GLOBAL;
```

### 2. Vector Database Specialization
```python
# Offload vector operations to specialized database
# Current: PostgreSQL + pgvector
# Future: Pinecone/Weaviate for billion-scale vectors

class VectorDatabaseService:
    async def store_atom_embedding(self, atom_id: str, embedding: List[float]):
        # Store in vector database
        await pinecone.upsert(
            vectors=[(atom_id, embedding)],
            namespace=self.tenant_id
        )
    
    async def find_similar_atoms(self, query_embedding: List[float], k: int = 10):
        # Query vector database
        results = await pinecone.query(
            vector=query_embedding,
            top_k=k,
            namespace=self.tenant_id
        )
        
        # Retrieve full atom data from PostgreSQL
        atom_ids = [match.id for match in results.matches]
        return await self.fetch_atoms_by_ids(atom_ids)
```

### 3. Edge Computing Integration
```typescript
// Local persistence for offline capability
class EdgePersistence {
  async syncToCloud() {
    // Sync local changes to cloud
    const changes = await this.localDB.getUnsyncedChanges();
    
    for (const change of changes) {
      await this.cloudDB.applyChange(change);
      await this.localDB.markSynced(change.id);
    }
  }
  
  async syncFromCloud() {
    // Pull cloud changes to edge
    const lastSync = await this.localDB.getLastSyncTimestamp();
    const changes = await this.cloudDB.getChangesSince(lastSync);
    
    await this.localDB.applyChanges(changes);
  }
}
```

---

## Cross-References

- **Rendering**: See [`product-os-implementation-rendering.md`](product-os-implementation-rendering.md) for annotation positioning and display
- **Search**: See [`product-os-implementation-search.md`](product-os-implementation-search.md) for querying across polyglot stores
- **Extraction**: See [`product-os-implementation-extraction.md`](product-os-implementation-extraction.md) for AI-driven knowledge creation
- **Frontend**: See [`product-os-implementation-frontend.md`](product-os-implementation-frontend.md) for data fetching patterns

---

**Status**: Technical reference - architecture patterns defined  
**Next**: Review with database engineering team for implementation planning