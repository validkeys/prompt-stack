# Product-OS Data Structure: Notion Influence Analysis

**Date**: 2026-01-11  
**Status**: Research & Analysis - Work in Progress  
**Related**: [`product-os-implementation-data-structure.md`](../product-os-implementation-data-structure.md) (Current schema), [`product-os-implementation-persistence.md`](../product-os-implementation-persistence.md) (Polyglot architecture)

---

## Overview

Analysis of Notion's data structure to inform Product-OS schema design. Notion provides a proven model for document-based knowledge management with strong user experience patterns.

## Notion's Core Data Model

### Hierarchy
```
Database → Data Source → Page → Block
```

### Key Components
1. **Database**: Container for data sources (permissions, organization)
2. **Data Source**: Schema definition (properties) + page container
3. **Page**: Instance with property values + content blocks  
4. **Block**: Typed content units (paragraph, heading, list, etc.)

### Block Structure (Content Storage)
```json
{
  "object": "block",
  "id": "uuid",
  "parent": {"type": "page_id", "page_id": "..."},
  "type": "paragraph",
  "paragraph": {
    "rich_text": [
      {
        "type": "text",
        "text": {"content": "text with formatting"},
        "annotations": {"bold": true, "italic": false}
      }
    ]
  },
  "has_children": false
}
```

### Property System (Schema)
- **Data Sources** define property schemas (title, rich_text, number, select, multi_select, date, relation)
- **Pages** contain property values conforming to schema
- **Rich Text** with inline annotations (no standoff markup)

---

## Key Differences from Current Product-OS Schema

### 1. Content Storage Pattern
| Aspect | Notion | Current Product-OS |
|--------|--------|-------------------|
| **Content storage** | Blocks (typed units) | Plain text + annotations |
| **Format handling** | Inline rich text | Standoff markup |
| **Hierarchy** | Explicit parent references | Implicit via `domain_id` + `path` |
| **Schema definition** | Properties at data source level | Templates (proposed) + domains |

### 2. Atomic Knowledge Integration
- **Notion**: No built-in atomic knowledge extraction
- **Product-OS**: Atoms + Molecules as first-class citizens
- **Relationship storage**: Notion uses relations property type; we use Neo4j graph

### 3. Span Calculation Problem
- **Notion**: No span calculations (blocks are units)
- **Current approach**: Text spans + annotations (complex when text changes)
- **Solution**: Annotate blocks instead of character spans

---

## Proposed Revisions Inspired by Notion

### 1. Add `blocks` Table (Notion-style content storage)
```sql
CREATE TABLE blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    
    -- Hierarchy
    parent_id UUID REFERENCES blocks(id) ON DELETE CASCADE,
    parent_type VARCHAR(20), -- 'block', 'document', 'domain'
    
    -- Content type
    block_type VARCHAR(50) NOT NULL, 
    -- Types: 'paragraph', 'heading_1', 'heading_2', 'heading_3', 
    --        'bulleted_list_item', 'numbered_list_item', 'code',
    --        'quote', 'callout', 'table', 'image', 'file'
    
    -- Content storage
    content JSONB NOT NULL, -- Rich text array or block-specific config
    
    -- Atomic knowledge references
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

### 2. Revise `documents` Table (Notion Pages)
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
    
    -- Schema-defined properties
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

### 3. Add `properties` Table (Schema Definitions)
```sql
CREATE TABLE properties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    domain_id UUID REFERENCES domains(id),
    
    -- Property definition
    name VARCHAR(255) NOT NULL,
    property_type VARCHAR(50) NOT NULL, 
    -- Types: 'title', 'rich_text', 'number', 'boolean', 'select', 
    --        'multi_select', 'date', 'relation', 'formula'
    
    -- Configuration
    config JSONB, -- Type-specific (options for select, formula, etc.)
    
    -- Metadata
    sort_order INTEGER DEFAULT 0,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### 4. Keep But Revise Atomic Tables
```sql
-- Atoms become typed property values
CREATE TABLE property_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    property_id UUID REFERENCES properties(id),
    document_id UUID REFERENCES documents(id),
    
    -- Polymorphic value (atomic)
    atom_type VARCHAR(50),
    text_value TEXT,
    number_value NUMERIC,
    boolean_value BOOLEAN,
    date_value TIMESTAMPTZ,
    json_value JSONB,
    
    -- Versioning
    version INTEGER DEFAULT 1,
    previous_version UUID REFERENCES property_values(id),
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Molecules become relationships between property values
CREATE TABLE molecule_relations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    molecule_type VARCHAR(50), -- 'ConditionMolecule', 'RequiresMolecule'
    
    -- Source and target property values
    source_value_id UUID REFERENCES property_values(id),
    target_value_id UUID REFERENCES property_values(id),
    
    -- Configuration
    config JSONB,
    enforcement_level VARCHAR(20), -- 'block', 'warn', 'inform'
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### 5. Update `annotations` Table
```sql
-- Replace span-based annotations with block-based
CREATE TABLE annotations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    
    -- Reference to block instead of document + spans
    block_id UUID REFERENCES blocks(id) ON DELETE CASCADE,
    
    -- Knowledge reference
    knowledge_type VARCHAR(20) NOT NULL, -- 'property_value', 'molecule'
    knowledge_id UUID NOT NULL,
    
    -- Context within block (optional)
    block_context JSONB, -- Could store text segment if needed
    
    -- Metadata
    confidence FLOAT DEFAULT 1.0,
    source VARCHAR(50) DEFAULT 'ai', -- 'ai', 'manual', 'code_analysis'
    validated_by UUID,
    validated_at TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## Hybrid Architecture: Notion Structure + Atomic Semantics

### Two-Layer Design
```
USER LAYER (Notion-style)
├── Documents (Pages)
├── Blocks (Content units)
└── Properties (Schema + values)

ATOMIC LAYER (Knowledge extraction)
├── Property Values (Atoms)
├── Molecule Relations (Relationships)
└── Annotations (Block ↔ Knowledge mapping)
```

### Data Flow
1. **User creates content** → Blocks stored in PostgreSQL
2. **AI extracts knowledge** → Property values + molecules created
3. **Annotations created** → Link blocks to atomic knowledge
4. **Graph relationships** → Sync to Neo4j for traversal
5. **Search indexes** → Update Elasticsearch/vector DB

### Advantages
- **User experience**: Familiar block-based editing
- **Content management**: No span recalculation issues  
- **Knowledge extraction**: Atomic layer remains for semantics
- **Flexibility**: Blocks can contain multiple atoms/molecules
- **Performance**: Block-based queries simpler than span-based

---

## Recommendation

### Adopt Notion's Block Model for Content Storage
**Why**: 
1. **Proven UX**: Billions of users understand block-based editors
2. **Technical simplicity**: No complex span calculations
3. **Future compatibility**: Easier to support rich media, tables, etc.
4. **Performance**: Block queries more efficient than text span queries

### But Keep Atomic Layer for Knowledge
**Why**:
1. **Semantic extraction**: AI identifies atoms in block content
2. **Knowledge graph**: Neo4j relationships enable complex queries
3. **Cross-document intelligence**: Find related concepts across documents
4. **Change impact analysis**: Trace dependencies through atomic relationships

### Implementation Strategy
1. **Phase 1**: Implement blocks + documents tables
2. **Phase 2**: Migrate existing documents to block structure
3. **Phase 3**: Update atomic layer to reference blocks
4. **Phase 4**: Update UI to block-based editor

---

## Open Questions

1. **Block granularity**: How fine-grained should blocks be? Paragraph-level? Sentence-level?
2. **Atomic extraction**: Should atoms be extracted per-block or across document?
3. **Versioning**: Block-level versioning vs document-level versioning?
4. **Template integration**: How do templates map to block structures?
5. **Performance**: Will block-based queries scale better than span-based?

---

## Next Steps

1. **Prototype block structure** with sample documents
2. **Test atomic extraction** from block content
3. **Compare query performance** vs current span-based approach
4. **Design migration path** from current schema
5. **Evaluate block-based editor options** (Slate.js, Lexical, ProseMirror)

---

**Status**: Research complete - ready for prototyping  
**Next**: Create prototype schema with Notion-inspired structure