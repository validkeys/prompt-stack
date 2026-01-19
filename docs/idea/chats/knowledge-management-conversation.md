# Knowledge Management System with Atomic Design - Complete Conversation

## Initial Question: Token Types and Relationships

**Human:** Imagine an application where the various domain experts keep internal knowledge up to date in a database. The application is then setup to tokenize and find relationships between other knowledge. How would define the types of tokens? What vocabulary would you use? How would you define and find relationships between other documents

**Assistant:**

I'd approach this knowledge management system with a multi-layered tokenization and relationship strategy:

### Token Types

**Semantic Tokens:**
- **Entities**: People, organizations, products, locations, technical components
- **Concepts**: Domain-specific ideas, methodologies, principles
- **Events**: Milestones, incidents, decisions, changes
- **Attributes**: Properties, metrics, classifications

**Structural Tokens:**
- **Document metadata**: Author, department, date, version, security level
- **Section types**: Problem statement, solution, rationale, constraints
- **References**: Citations to other internal documents or external sources

**Relational Tokens:**
- **Dependency markers**: "requires", "supersedes", "implements"
- **Temporal markers**: "before", "during", "after"
- **Causal indicators**: "causes", "results in", "prevents"

### Vocabulary Strategy

I'd use a **hybrid vocabulary** approach:

**Domain-Specific Ontology**: Work with each domain expert group to establish canonical terms. For example:
- Engineering: "deployment pipeline", "tech debt", "API contract"
- Legal: "compliance requirement", "risk assessment", "audit trail"
- Product: "user story", "feature flag", "A/B test"

**Cross-Domain Bridge Terms**: Identify concepts that appear across domains but may use different terminology (e.g., "customer requirement" vs "user need" vs "specification")

**Controlled Synonyms**: Map variations to canonical forms while preserving context-specific usage

### Relationship Discovery

**Explicit Relationships** (defined by experts):
- Document templates with relationship fields
- Tagged dependencies when creating/updating documents
- Formal linking between versions and related work

**Implicit Relationships** (discovered algorithmically):

*Co-occurrence Analysis*: Documents frequently referencing the same entities or concepts likely share context

*Semantic Similarity*: Use embeddings to find documents discussing similar topics even with different terminology

*Temporal Patterns*: Documents created in sequence by the same team often build on each other

*Citation Networks*: Analyze which documents reference which others to build knowledge graphs

*Contradiction Detection*: Flag when newer documents conflict with established knowledge

**Relationship Types I'd Track:**
- **Hierarchical**: parent/child, general/specific
- **Sequential**: prerequisite, next step, evolution
- **Complementary**: alternative approach, supporting evidence
- **Conflicting**: contradiction, deprecation, correction

### Implementation Considerations

The system should surface relationships contextually - when someone views a document about API design, automatically suggest related docs about security requirements, deployment procedures, and past architectural decisions. The key is making the relationships actionable for users, not just mathematically interesting.

---

## Question 2: Atomic Design Data Structure

**Human:** Let's supposed we were using the atomic design for how we label these data components. Suppose we had a workflow for extracting these atoms/molecules/organisms, how would we store this text with the individual atoms highlighted in the browser? How is the data persisted instead of just raw text?

**Assistant:**

Great question - atomic design principles work really well for knowledge decomposition. Here's how I'd structure the storage and rendering:

### Data Model

**Atom-Level Storage:**
```json
{
  "id": "atom_uuid_123",
  "type": "entity|concept|event|attribute",
  "canonical_form": "API Gateway",
  "variants": ["api gateway", "API-GW", "gateway service"],
  "domain": "engineering",
  "metadata": {
    "created": "2024-01-15",
    "confidence": 0.95,
    "source_documents": ["doc_456", "doc_789"]
  }
}
```

**Molecule-Level Storage** (atoms in context):
```json
{
  "id": "mol_uuid_456",
  "type": "relationship|pattern|phrase",
  "atoms": [
    {"atom_id": "atom_123", "role": "subject"},
    {"atom_id": "atom_124", "role": "action"},
    {"atom_id": "atom_125", "role": "object"}
  ],
  "context": "The API Gateway authenticates all requests",
  "relationship_type": "performs_function"
}
```

**Document Storage** (organisms):
```json
{
  "id": "doc_uuid_789",
  "title": "Authentication Architecture",
  "raw_text": "...",
  "annotated_spans": [
    {
      "start": 45,
      "end": 56,
      "atom_id": "atom_123",
      "atom_type": "entity",
      "canonical": "API Gateway"
    },
    {
      "start": 120,
      "end": 145,
      "molecule_id": "mol_456",
      "relationship": "authentication_flow"
    }
  ],
  "sections": [
    {
      "type": "problem_statement",
      "spans": [0, 250],
      "key_atoms": ["atom_123", "atom_130"]
    }
  ]
}
```

### Browser Rendering Strategy

**HTML Generation with Standoff Markup:**

Instead of embedding markup in text (which creates versioning nightmares), use standoff annotations:

```html
<div class="knowledge-doc" data-doc-id="doc_789">
  <p data-span="0-250">
    <span class="text-chunk" data-offset="0">The </span>
    <span class="atom entity" 
          data-atom-id="atom_123"
          data-canonical="API Gateway"
          data-offset="45-56">API Gateway</span>
    <span class="text-chunk" data-offset="57"> authenticates all requests...</span>
  </p>
</div>
```

**Interactive Layer via JavaScript:**

```javascript
// On hover/click, enrich with:
{
  atomId: "atom_123",
  canonical: "API Gateway",
  related_atoms: [...],
  appears_in: ["doc_456", "doc_789", "doc_890"],
  relationships: [
    {type: "performs", target: "authentication"},
    {type: "depends_on", target: "OAuth Service"}
  ],
  expert_notes: [...]
}
```

### Persistence Architecture

**Multi-Store Approach:**

**1. Document Store** (PostgreSQL/MongoDB):
- Original raw text with versions
- Metadata (author, date, domain)
- Section boundaries

**2. Graph Database** (Neo4j/DGraph):
- Atoms as nodes
- Relationships as edges
- Enables traversal queries like "find all documents that discuss concepts related to this one"

**3. Vector Store** (Pinecone/Weaviate):
- Embeddings for atoms, molecules, sections
- Semantic search and similarity
- Cross-domain concept matching

**4. Annotation Index** (Elasticsearch):
- Fast lookup: "which documents contain atom_123?"
- Span-based search
- Faceted filtering by atom type, domain, etc.

### Extraction Workflow Persistence

**Extraction Pipeline Results:**

```json
{
  "extraction_id": "ext_uuid_555",
  "document_id": "doc_789",
  "version": 3,
  "timestamp": "2024-01-15T10:30:00Z",
  "extractor": "expert_manual|ai_assisted|automated",
  "atoms_extracted": [
    {
      "text": "API Gateway",
      "span": [45, 56],
      "proposed_atom_id": "atom_123",
      "confidence": 0.98,
      "status": "approved|pending|rejected"
    }
  ],
  "molecules_extracted": [...],
  "review_status": "pending_review",
  "reviewer": "expert_jane"
}
```

### Browser Update Strategy

**Real-Time Highlighting Without Page Reload:**

When experts approve new atoms/molecules:

1. **WebSocket notification** to active viewers
2. **Incremental DOM update** using data-offset attributes
3. **Highlight animation** for newly recognized atoms
4. **Tooltip update** with enriched relationship data

**Handling Text Edits:**

This is tricky - when the document text changes, spans shift. Solutions:

- **Content-addressable anchors**: Hash text chunks before/after atoms
- **Fuzzy matching**: Re-anchor spans using surrounding context
- **Version locking**: Lock atom positions to specific document versions
- **Conflict resolution**: Expert review when text changes invalidate spans

### Example User Experience

Expert views document → atoms highlighted in subtle colors → hover shows canonical form + quick stats → click opens sidebar with:
- All relationships
- Other documents using this atom
- Alternate phrasings
- Edit/refine options
- "Why was this extracted?" provenance

The key advantage: you can update the canonical definition of an atom or add new relationships, and **all documents automatically reflect this** without re-processing the text. The annotations are living metadata, not baked-in markup.

---

## Question 3: Optimal Database Selection

**Human:** What type of database is optimal for storing data on this way?

**Assistant:**

For this atomic knowledge system, you'll want a **polyglot persistence** approach - no single database handles all requirements optimally. Here's what I'd recommend:

### Primary Recommendation: Multi-Database Architecture

#### 1. **Graph Database (Neo4j or Amazon Neptune)** - *The Core*

**Why it's optimal:**
- Atoms, molecules, and organisms are fundamentally graph structures
- Relationship traversal queries are first-class operations
- Natural fit for "find all concepts related to X within 2 hops"
- ACID compliance for critical knowledge integrity

**Schema example:**
```cypher
// Atoms as nodes
(:Atom {id, canonical, type, domain})
(:Document {id, title, version})

// Relationships
(:Atom)-[:APPEARS_IN {spans: [{start, end}]}]->(:Document)
(:Atom)-[:RELATED_TO {type: "performs|requires|conflicts"}]->(:Atom)
(:Atom)-[:PART_OF]->(:Molecule)
(:Document)-[:CITES]->(:Document)
(:Document)-[:SUPERSEDES]->(:Document)
```

**Queries this enables:**
```cypher
// Find all documents discussing related concepts
MATCH (a:Atom {canonical: "API Gateway"})-[:RELATED_TO*1..2]-(related:Atom)
      -[:APPEARS_IN]->(doc:Document)
RETURN DISTINCT doc

// Find knowledge conflicts
MATCH (a1:Atom)-[r:CONFLICTS_WITH]->(a2:Atom)
WHERE a1.domain = "engineering" AND a2.domain = "security"
```

#### 2. **Document Store (PostgreSQL with JSONB)** - *The Source of Truth*

**Why:**
- Excellent for versioned document storage
- JSONB handles flexible annotation structures
- Strong consistency for document updates
- Full-text search capabilities built-in

**Schema:**
```sql
CREATE TABLE documents (
  id UUID PRIMARY KEY,
  title TEXT,
  raw_text TEXT,
  annotations JSONB,  -- standoff markup
  metadata JSONB,
  version INTEGER,
  created_at TIMESTAMP,
  created_by UUID
);

CREATE INDEX idx_annotations ON documents 
  USING GIN (annotations jsonb_path_ops);

-- Query atoms by span
SELECT * FROM documents 
WHERE annotations @> '[{"atom_id": "atom_123"}]';
```

**Alternative: MongoDB** if you need more schema flexibility, but PostgreSQL's JSONB + relational capabilities usually win for this use case.

#### 3. **Vector Database (Pinecone, Weaviate, or pgvector)** - *Semantic Layer*

**Why:**
- Semantic similarity across different terminology
- "Find documents similar to this one" queries
- Cross-domain concept matching
- Embedding-based search

**Storage:**
```python
# Each atom gets embedded
{
  "id": "atom_123",
  "vector": [0.123, -0.456, ...],  # 1536 dims from OpenAI
  "metadata": {
    "canonical": "API Gateway",
    "domain": "engineering",
    "doc_count": 47
  }
}

# Query: "security authentication patterns"
# Returns semantically similar atoms even if exact terms don't match
```

**pgvector advantage**: Keep everything in PostgreSQL if you want to minimize infrastructure complexity.

#### 4. **Search Engine (Elasticsearch)** - *Discovery Layer*

**Why:**
- Fast full-text search across all documents
- Faceted search (filter by domain, date, expert, atom type)
- Aggregations for analytics
- Highlighting of search terms in context

**Index structure:**
```json
{
  "document_id": "doc_789",
  "title": "Authentication Architecture",
  "full_text": "...",
  "atoms": [
    {"id": "atom_123", "canonical": "API Gateway", "type": "entity"}
  ],
  "domain": "engineering",
  "expert": "jane_doe",
  "created": "2024-01-15"
}
```

### Alternative: All-in-One Solutions

If you want to minimize operational complexity:

#### **Option A: PostgreSQL + Extensions** (Good for MVP)
- **pgvector**: Vector similarity search
- **JSONB**: Flexible document storage
- **pg_trgm**: Fuzzy text matching
- **Apache AGE**: Graph queries on PostgreSQL

**Pros:** Single database, simpler ops, ACID everywhere
**Cons:** Graph queries not as elegant, scaling limitations

#### **Option B: ArangoDB** (Multi-Model)
- Native graph + document + key-value support
- AQL query language handles all patterns
- Single database to manage

**Pros:** True multi-model, easier to reason about
**Cons:** Smaller ecosystem, less mature than specialists

#### **Option C: DGraph** (GraphQL-Native Graph DB)
- Built-in full-text search
- GraphQL schema = your data model
- Distributed by default

**Pros:** Modern API, good for app integration
**Cons:** Less mature than Neo4j, smaller community

### My Recommended Stack for Production

**Tier 1 (Essential):**
- **Neo4j**: Atom/molecule relationships, knowledge graph
- **PostgreSQL**: Document versions, annotations, metadata

**Tier 2 (High Value):**
- **Elasticsearch**: Search and discovery
- **pgvector or Pinecone**: Semantic search

**Tier 3 (Nice to Have):**
- **Redis**: Caching frequently accessed atoms, real-time notifications
- **S3/Object Store**: Large document attachments, document history

### Data Flow Example

```
1. Expert creates document
   ↓
2. Save to PostgreSQL (source of truth)
   ↓
3. Extract atoms/molecules
   ↓
4. Create/link nodes in Neo4j
   ↓
5. Generate embeddings → Vector DB
   ↓
6. Index in Elasticsearch
   ↓
7. Cache hot paths in Redis

Query flow:
- Text search → Elasticsearch
- "Related concepts" → Neo4j
- "Similar documents" → Vector DB
- Document retrieval → PostgreSQL
```

### Storage Estimates

For a knowledge base with:
- 10,000 documents
- 50,000 unique atoms
- 200,000 relationships

**Neo4j**: ~500MB-1GB (nodes + relationships)
**PostgreSQL**: ~5-10GB (full documents + annotations)
**Vector DB**: ~300MB (50k atoms × 1536 dims × 4 bytes)
**Elasticsearch**: ~8-15GB (inverted indices)

### Decision Framework

**Start with Neo4j + PostgreSQL if:**
- Relationship queries are central to your use case
- You need complex graph traversals
- Team has graph database experience

**Start with PostgreSQL + pgvector if:**
- Simplicity and operational overhead matter most
- Budget constrained
- Graph queries are relatively simple (2-3 hops max)

**Add Elasticsearch when:**
- You have >1000 documents
- Users need faceted search and filtering
- Full-text search performance becomes critical

The polyglot approach gives you the best tool for each job, but requires more operational maturity. For an MVP, I'd start with **PostgreSQL + pgvector + a lightweight graph layer**, then migrate to Neo4j when graph complexity demands it.

---

## Question 4: Rendering, React Libraries, and Search Implementation

**Human:** In working with this data structure, how I am printing a well formatted document of text for a user? Are there good react libraries that make working with this type of data structure easier? What about search, how do I leverage this data structure to make searching fast?

**Assistant:**

Great questions - let's tackle document rendering, React tooling, and search optimization.

### Rendering Documents with Annotations

#### The Core Challenge

You have text stored separately from its annotations (standoff markup). You need to:
1. Reconstruct highlighted text from spans
2. Handle overlapping annotations (an atom inside a molecule)
3. Make it interactive (hover, click, navigate)
4. Keep it performant

#### Approach 1: Pre-Processing Server-Side

**Build an annotated structure before sending to client:**

```javascript
// Server-side transformation
function buildAnnotatedDocument(doc) {
  const { raw_text, annotated_spans } = doc;
  
  // Sort spans by position
  const sorted = annotated_spans.sort((a, b) => a.start - b.start);
  
  // Build segment tree for overlapping spans
  const segments = [];
  let lastEnd = 0;
  
  sorted.forEach(span => {
    // Add plain text before this span
    if (span.start > lastEnd) {
      segments.push({
        type: 'text',
        content: raw_text.slice(lastEnd, span.start)
      });
    }
    
    // Add annotated span
    segments.push({
      type: 'atom',
      content: raw_text.slice(span.start, span.end),
      atomId: span.atom_id,
      canonical: span.canonical,
      metadata: span.metadata
    });
    
    lastEnd = Math.max(lastEnd, span.end);
  });
  
  // Add remaining text
  if (lastEnd < raw_text.length) {
    segments.push({
      type: 'text',
      content: raw_text.slice(lastEnd)
    });
  }
  
  return segments;
}
```

**Handle overlapping/nested annotations:**

```javascript
// Build a tree structure for nested spans
function buildAnnotationTree(text, spans) {
  const sorted = [...spans].sort((a, b) => {
    if (a.start !== b.start) return a.start - b.start;
    return b.end - a.end; // Longer spans first
  });
  
  const root = { start: 0, end: text.length, children: [] };
  
  function insertSpan(parent, span) {
    // Find children that this span contains
    const contained = parent.children.filter(
      child => child.start >= span.start && child.end <= span.end
    );
    
    // Remove contained children from parent
    parent.children = parent.children.filter(
      child => !contained.includes(child)
    );
    
    // Add contained children to this span
    span.children = contained;
    
    // Add this span to parent
    parent.children.push(span);
  }
  
  sorted.forEach(span => insertSpan(root, { ...span, children: [] }));
  
  return root;
}
```

### React Libraries & Patterns

#### **Recommended: Slate.js or Lexical**

Both are excellent for rich text with custom annotations:

**Slate.js** (more mature, larger ecosystem):

```jsx
import { Slate, Editable, withReact } from 'slate-react';
import { createEditor } from 'slate';

// Define custom leaf rendering
const Leaf = ({ attributes, children, leaf }) => {
  if (leaf.atom) {
    return (
      <span
        {...attributes}
        className="atom-highlight"
        data-atom-id={leaf.atomId}
        onMouseEnter={() => showAtomTooltip(leaf.atomId)}
        onClick={() => navigateToAtom(leaf.atomId)}
      >
        {children}
        <AtomPopover atomId={leaf.atomId} />
      </span>
    );
  }
  
  return <span {...attributes}>{children}</span>;
};

function KnowledgeDocument({ document }) {
  const editor = useMemo(() => withReact(createEditor()), []);
  
  // Convert your annotation structure to Slate format
  const initialValue = useMemo(() => 
    convertToSlateFormat(document), 
    [document]
  );
  
  return (
    <Slate editor={editor} value={initialValue}>
      <Editable
        readOnly
        renderLeaf={props => <Leaf {...props} />}
      />
    </Slate>
  );
}

// Convert your data structure
function convertToSlateFormat(doc) {
  const segments = buildAnnotatedDocument(doc);
  
  return [{
    type: 'paragraph',
    children: segments.map(seg => ({
      text: seg.content,
      ...(seg.type === 'atom' && {
        atom: true,
        atomId: seg.atomId,
        canonical: seg.canonical
      })
    }))
  }];
}
```

**Lexical** (Meta's newer framework, more performant):

```jsx
import { LexicalComposer } from '@lexical/react/LexicalComposer';
import { RichTextPlugin } from '@lexical/react/LexicalRichTextPlugin';
import { ContentEditable } from '@lexical/react/LexicalContentEditable';

// Custom node for atoms
class AtomNode extends TextNode {
  __atomId;
  __canonical;
  
  static getType() {
    return 'atom';
  }
  
  createDOM(config) {
    const span = document.createElement('span');
    span.className = 'atom-highlight';
    span.dataset.atomId = this.__atomId;
    return span;
  }
  
  updateDOM(prevNode, dom) {
    return false;
  }
}

function KnowledgeDocumentLexical({ document }) {
  const initialConfig = {
    namespace: 'KnowledgeDoc',
    nodes: [AtomNode],
    editable: false,
    theme: {
      // Custom styling
    }
  };
  
  return (
    <LexicalComposer initialConfig={initialConfig}>
      <RichTextPlugin
        contentEditable={<ContentEditable />}
        placeholder={null}
      />
    </LexicalComposer>
  );
}
```

#### **Alternative: ProseMirror**

More low-level but extremely powerful for custom annotation UIs:

```javascript
import { Schema } from 'prosemirror-model';
import { EditorState } from 'prosemirror-state';
import { EditorView } from 'prosemirror-view';

const schema = new Schema({
  nodes: {
    doc: { content: 'block+' },
    paragraph: { content: 'inline*', group: 'block' },
    text: { group: 'inline' }
  },
  marks: {
    atom: {
      attrs: { atomId: {}, canonical: {} },
      toDOM: mark => ['span', {
        class: 'atom-highlight',
        'data-atom-id': mark.attrs.atomId
      }]
    }
  }
});
```

#### **Simple Solution: Custom React Component**

If you don't need editing capabilities, roll your own:

```jsx
function AnnotatedText({ segments, onAtomClick, onAtomHover }) {
  const [hoveredAtom, setHoveredAtom] = useState(null);
  
  return (
    <div className="annotated-document">
      {segments.map((segment, idx) => {
        if (segment.type === 'text') {
          return <span key={idx}>{segment.content}</span>;
        }
        
        return (
          <AtomSpan
            key={idx}
            segment={segment}
            isHovered={hoveredAtom === segment.atomId}
            onMouseEnter={() => {
              setHoveredAtom(segment.atomId);
              onAtomHover?.(segment.atomId);
            }}
            onMouseLeave={() => setHoveredAtom(null)}
            onClick={() => onAtomClick?.(segment.atomId)}
          />
        );
      })}
    </div>
  );
}

function AtomSpan({ segment, isHovered, onMouseEnter, onMouseLeave, onClick }) {
  return (
    <span
      className={`atom-highlight atom-${segment.metadata.type} ${
        isHovered ? 'hovered' : ''
      }`}
      data-atom-id={segment.atomId}
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
      onClick={onClick}
    >
      {segment.content}
      {isHovered && (
        <AtomTooltip
          atomId={segment.atomId}
          canonical={segment.canonical}
        />
      )}
    </span>
  );
}
```

#### **For Tooltips/Popovers: Floating UI**

```jsx
import { useFloating, offset, flip, shift } from '@floating-ui/react';

function AtomTooltip({ atomId, canonical, reference }) {
  const { data: atomData } = useAtomDetails(atomId); // React Query
  
  const { refs, floatingStyles } = useFloating({
    placement: 'top',
    middleware: [offset(10), flip(), shift()]
  });
  
  return (
    <div ref={refs.setFloating} style={floatingStyles} className="atom-tooltip">
      <h4>{canonical}</h4>
      <p>Type: {atomData?.type}</p>
      <p>Appears in {atomData?.documentCount} documents</p>
      <div className="related-atoms">
        {atomData?.relatedAtoms.map(rel => (
          <AtomChip key={rel.id} atom={rel} />
        ))}
      </div>
    </div>
  );
}
```

### Leveraging Data Structure for Fast Search

#### 1. **Multi-Level Search Strategy**

```javascript
// Search API endpoint
async function search(query, filters = {}) {
  // Level 1: Elasticsearch full-text search
  const textResults = await elasticsearchClient.search({
    index: 'documents',
    body: {
      query: {
        bool: {
          must: [
            { multi_match: {
              query: query,
              fields: ['title^3', 'full_text', 'atoms.canonical^2'],
              type: 'best_fields'
            }}
          ],
          filter: buildFilters(filters) // domain, date, expert
        }
      },
      highlight: {
        fields: {
          full_text: { fragment_size: 150, number_of_fragments: 3 }
        }
      },
      size: 20
    }
  });
  
  // Level 2: Vector similarity for semantic search
  const embedding = await generateEmbedding(query);
  const semanticResults = await vectorDB.query({
    vector: embedding,
    topK: 10,
    filter: { domain: filters.domain }
  });
  
  // Level 3: Graph expansion for related concepts
  const relatedAtoms = await neo4j.run(`
    MATCH (a:Atom)-[:APPEARS_IN]->(d:Document)
    WHERE d.id IN $docIds
    MATCH (a)-[:RELATED_TO*1..2]-(related:Atom)-[:APPEARS_IN]->(relatedDoc:Document)
    WHERE relatedDoc.id NOT IN $docIds
    RETURN DISTINCT relatedDoc
    LIMIT 10
  `, { docIds: textResults.hits.map(h => h._id) });
  
  // Merge and rank results
  return mergeAndRank(textResults, semanticResults, relatedAtoms);
}
```

#### 2. **Atom-Aware Search UI**

```jsx
function SearchInterface() {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState([]);
  const [atomFilters, setAtomFilters] = useState([]);
  
  // Debounced search with atom suggestion
  const { data: suggestions } = useQuery(
    ['atom-suggestions', query],
    () => suggestAtoms(query),
    { enabled: query.length > 2 }
  );
  
  return (
    <div className="search-container">
      <SearchBar
        value={query}
        onChange={setQuery}
        suggestions={suggestions}
        onAtomSelect={atom => setAtomFilters([...atomFilters, atom])}
      />
      
      <ActiveAtomFilters
        atoms={atomFilters}
        onRemove={atom => setAtomFilters(atomFilters.filter(a => a !== atom))}
      />
      
      <SearchResults results={results}>
        {results.map(result => (
          <SearchResultCard
            key={result.id}
            document={result}
            highlights={result.highlights}
            matchedAtoms={result.matchedAtoms}
          />
        ))}
      </SearchResults>
    </div>
  );
}
```

#### 3. **Pre-computed Search Indexes**

Store materialized views in Elasticsearch for instant results:

```javascript
// When document is created/updated
async function indexDocument(doc) {
  const atoms = await getDocumentAtoms(doc.id);
  const relatedDocs = await getRelatedDocuments(doc.id);
  
  await elasticsearchClient.index({
    index: 'documents',
    id: doc.id,
    body: {
      title: doc.title,
      full_text: doc.raw_text,
      atoms: atoms.map(a => ({
        id: a.id,
        canonical: a.canonical,
        type: a.type,
        domain: a.domain
      })),
      atom_canonical_list: atoms.map(a => a.canonical), // For fast filtering
      related_doc_ids: relatedDocs.map(d => d.id),
      expert: doc.created_by,
      domain: doc.metadata.domain,
      created_at: doc.created_at,
      // Pre-compute popular facets
      popularity_score: calculatePopularity(doc),
      recency_score: calculateRecency(doc)
    }
  });
}
```

#### 4. **Cached Graph Queries**

```javascript
// Redis cache for common relationship queries
async function getRelatedAtoms(atomId) {
  const cacheKey = `related:${atomId}`;
  
  // Check cache first
  const cached = await redis.get(cacheKey);
  if (cached) return JSON.parse(cached);
  
  // Query Neo4j
  const result = await neo4j.run(`
    MATCH (a:Atom {id: $atomId})-[:RELATED_TO*1..2]-(related:Atom)
    RETURN related, count(*) as strength
    ORDER BY strength DESC
    LIMIT 20
  `, { atomId });
  
  // Cache for 1 hour
  await redis.setex(cacheKey, 3600, JSON.stringify(result));
  
  return result;
}
```

#### 5. **Incremental Search with Typeahead**

```jsx
function AtomTypeahead({ onSelect }) {
  const [input, setInput] = useState('');
  
  // Use Elasticsearch completion suggester
  const { data: suggestions } = useQuery(
    ['typeahead', input],
    async () => {
      const response = await fetch('/api/suggest', {
        method: 'POST',
        body: JSON.stringify({
          prefix: input,
          field: 'atoms.canonical.completion'
        })
      });
      return response.json();
    },
    { enabled: input.length > 1, staleTime: 5000 }
  );
  
  return (
    <Autocomplete
      options={suggestions || []}
      getOptionLabel={opt => opt.canonical}
      renderOption={opt => (
        <AtomOption atom={opt} matchedText={input} />
      )}
      onSelect={onSelect}
    />
  );
}
```

### Performance Optimizations

#### Virtual Scrolling for Long Documents

```jsx
import { FixedSizeList } from 'react-window';

function LargeDocument({ segments }) {
  const Row = ({ index, style }) => (
    <div style={style}>
      <AnnotatedText segments={segments[index]} />
    </div>
  );
  
  return (
    <FixedSizeList
      height={800}
      itemCount={segments.length}
      itemSize={50}
    >
      {Row}
    </FixedSizeList>
  );
}
```

#### Search Result Streaming

```javascript
// Server-side streaming
async function* streamSearchResults(query) {
  // Quick results first
  const quickResults = await elasticsearchClient.search({...});
  yield { type: 'quick', results: quickResults };
  
  // Semantic results
  const semanticResults = await vectorSearch(query);
  yield { type: 'semantic', results: semanticResults };
  
  // Graph expansion
  const graphResults = await neo4jSearch(query);
  yield { type: 'graph', results: graphResults };
}

// Client-side
function SearchResults() {
  const [results, setResults] = useState([]);
  
  useEffect(() => {
    const eventSource = new EventSource(`/api/search/stream?q=${query}`);
    
    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setResults(prev => [...prev, ...data.results]);
    };
    
    return () => eventSource.close();
  }, [query]);
}
```

### Recommended Stack

**For MVP:**
- Custom React component for rendering
- Elasticsearch for search
- React Query for data fetching
- Floating UI for tooltips

**For Production:**
- Lexical or Slate for rich text
- Elasticsearch + pgvector for hybrid search
- Redis for caching
- Server-side pre-processing of annotations
- Streaming search results

The key is **separating concerns**: annotations stay as data, rendering logic stays in React, and search leverages pre-computed indexes. This keeps your system fast and maintainable.

---

## Summary

This conversation covered the design and implementation of an atomic knowledge management system:

1. **Token Types & Relationships**: Semantic, structural, and relational tokens with explicit and implicit relationship discovery
2. **Data Structure**: Atomic design pattern with atoms, molecules, and organisms stored using standoff annotations
3. **Database Architecture**: Polyglot persistence using Neo4j (graph), PostgreSQL (documents), Vector DB (semantic search), and Elasticsearch (full-text search)
4. **Frontend Implementation**: React libraries (Slate.js, Lexical, ProseMirror) for rendering annotated documents with interactive tooltips and highlights
5. **Search Optimization**: Multi-level search strategy combining full-text, semantic, and graph-based approaches with caching and streaming results

The system enables domain experts to collaboratively build and maintain a knowledge base where concepts are automatically linked, searchable, and contextually enriched across all documents.
