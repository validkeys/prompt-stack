# Product-OS: Persistence Operations

**Date**: 2026-01-11  
**Status**: Operational Guidance Document  
**Version**: 1.0  
**Related**: [`product-os-data-structure.md`](product-os-data-structure.md) (Concept document), [`product-os-data-structure-reference.md`](product-os-data-structure-reference.md) (Technical reference), [`product-os.md`](product-os.md) (Knowledge OS concept)

---

## Overview

**Operational guidance** for running and maintaining the polyglot persistence layer of Knowledge OS. This document covers migration strategies, backup procedures, performance optimization, security implementation, monitoring, and future evolution planning.

**Focus Areas**:
- **Migration**: Zero-downtime schema changes, tenant tier upgrades, database technology transitions
- **Backup & Recovery**: Multi-region deployment, backup strategies, point-in-time recovery
- **Performance**: Query optimization, connection pooling, read replicas, sharding
- **Security**: Encryption, row-level security, audit logging
- **Monitoring**: Metrics, alerting, performance tracing
- **Evolution**: Distributed SQL, vector databases, edge computing

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
#!/bin/bash
# Backup script for polyglot persistence

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
- **Data Structure**: See [`product-os-data-structure-reference.md`](product-os-data-structure-reference.md) for schema definitions and architecture

---

**Status**: Operational guidance - ready for implementation  
**Next**: Review with operations team for production deployment planning