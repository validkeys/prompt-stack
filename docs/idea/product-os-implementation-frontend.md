# Product-OS: Frontend Component Patterns & Libraries

**Date**: 2026-01-11  
**Status**: Technical Reference Document  
**Version**: 1.0  
**Related**: [`product-os.md`](product-os.md) (Knowledge OS concept), [`product-os-data-structure.md`](product-os-data-structure.md) (Data structure concept), [`product-os-implementation-rendering.md`](product-os-implementation-rendering.md) (Standoff annotations), [`product-os-implementation-persistence.md`](product-os-implementation-persistence.md) (Polyglot database), [`product-os-implementation-search.md`](product-os-implementation-search.md) (Multi-level search), [`product-os-implementation-extraction.md`](product-os-implementation-extraction.md) (Extraction pipeline)

---

## Overview

**Frontend architecture** for Knowledge OS implements persona-specific dashboards, atomic knowledge editors, and collaborative interfaces using modern web technologies. This document covers component libraries, state management, and UI patterns that enable human-AI collaboration across different roles.

**Key Principles**:
- **Persona-first design**: Optimized interfaces for domain experts, developers, product owners
- **Progressive disclosure**: Simple defaults with advanced options available
- **Real-time collaboration**: Multi-user editing with conflict resolution
- **Accessibility first**: WCAG 2.1 AA compliance for all user interfaces

---

## Technology Stack

### Core Framework Decisions
```typescript
// Technology choices based on requirements
const techStack = {
  // UI Framework: React 18+ with TypeScript
  framework: {
    name: 'React',
    version: '18+',
    features: [
      'Concurrent Rendering',
      'Server Components',
      'Suspense for Data Fetching'
    ]
  },
  
  // Styling: Tailwind CSS + CSS Modules
  styling: {
    utility: 'Tailwind CSS',
    components: 'Headless UI + Radix UI',
    theming: 'CSS Variables + Dark Mode'
  },
  
  // State Management: Zustand + TanStack Query
  state: {
    client: 'Zustand (lightweight stores)',
    server: 'TanStack Query (data fetching/caching)',
    forms: 'React Hook Form + Zod validation'
  },
  
  // Real-time: Socket.io + CRDTs
  realtime: {
    websockets: 'Socket.io',
    collaboration: 'Yjs (CRDTs)',
    presence: 'Liveblocks'
  },
  
  // Visualization: D3.js + Three.js
  visualization: {
    charts: 'Recharts',
    graphs: 'Cytoscape.js',
    '3d': 'Three.js + React Three Fiber'
  }
};
```

### Project Structure
```
src/
├── personas/                 # Persona-specific code
│   ├── domain-expert/
│   ├── developer/
│   ├── product-owner/
│   ├── project-manager/
│   └── qa-engineer/
│
├── components/              # Shared components
│   ├── atomic/             # Atomic knowledge UI
│   │   ├── AtomDisplay/
│   │   ├── MoleculeEditor/
│   │   └── OrganismViewer/
│   │
│   ├── collaboration/      # Real-time collaboration
│   │   ├── PresenceIndicator/
│   │   ├── CollaborativeEditor/
│   │   └── CommentThread/
│   │
│   ├── visualization/      # Data visualization
│   │   ├── KnowledgeGraph/
│   │   ├── ValueFlowChart/
│   │   └── DependencyGraph/
│   │
│   └── layout/            # Layout components
│       ├── PersonaLayout/
│       ├── SidebarNavigation/
│       └── ResponsiveGrid/
│
├── hooks/                  # Custom React hooks
│   ├── useKnowledgeQuery/
│   ├── useRealTimeUpdates/
│   └── usePersonaPreferences/
│
├── stores/                 # Zustand stores
│   ├── knowledge.store.ts
│   ├── collaboration.store.ts
│   └── persona.store.ts
│
├── lib/                    # Utility libraries
│   ├── knowledge-parser/
│   ├── annotation-renderer/
│   └── search-client/
│
└── api/                    # API clients
    ├── knowledge.client.ts
    ├── search.client.ts
    └── realtime.client.ts
```

---

## Persona-Specific Interfaces

### Domain Expert Dashboard
```typescript
// Domain expert focus: Knowledge refinement and validation
const DomainExpertDashboard = () => {
  const { persona } = usePersona();
  const { pendingReviews, recentExtractions } = useKnowledge();
  
  return (
    <PersonaLayout persona="domain-expert">
      {/* Main content area */}
      <div className="grid grid-cols-3 gap-6">
        {/* Validation queue */}
        <ValidationQueue 
          items={pendingReviews}
          onReview={handleReview}
          autoSort="confidence"
        />
        
        {/* Knowledge refinement */}
        <KnowledgeRefinement 
          extractions={recentExtractions}
          onApprove={handleApprove}
          onEdit={handleEdit}
          onReject={handleReject}
        />
        
        {/* Domain-specific tools */}
        <DomainTools 
          domain={persona.domain}
          templates={domainTemplates}
          onInjectRule={handleInjectRule}
        />
      </div>
      
      {/* Sidebar with context */}
      <Sidebar>
        <DomainContext 
          domain={persona.domain}
          recentChanges={recentChanges}
          expertNetwork={expertNetwork}
        />
        
        <ConfidenceMetrics 
          extractions={recentExtractions}
          validationHistory={validationHistory}
        />
        
        <QuickActions 
          actions={[
            { label: 'Review AI extractions', icon: 'check', action: reviewAI },
            { label: 'Inject new rule', icon: 'rule', action: injectRule },
            { label: 'Validate compliance', icon: 'shield', action: validate }
          ]}
        />
      </Sidebar>
    </PersonaLayout>
  );
};
```

**Domain Expert UI Patterns**:
- **Validation workflows**: Streamlined approve/edit/reject interfaces
- **Rule injection**: Visual editors for constraint molecules
- **Confidence indicators**: Clear visual cues for AI-extracted knowledge
- **Domain context**: Relevant references, regulations, standards

### Developer Workspace
```typescript
// Developer focus: Implementation mapping and impact analysis
const DeveloperWorkspace = () => {
  const { currentTask, knowledgeContext } = useDeveloper();
  
  return (
    <PersonaLayout persona="developer">
      {/* Split view: Code + Knowledge */}
      <SplitPane defaultSize="60%">
        {/* Code editor */}
        <CodeEditor 
          file={currentTask.file}
          annotations={knowledgeContext.annotations}
          onEdit={handleCodeEdit}
        />
        
        {/* Knowledge panel */}
        <KnowledgePanel>
          <AtomMapping 
            codeSelection={codeSelection}
            knowledgeAtoms={knowledgeContext.atoms}
            onMap={handleMapAtom}
          />
          
          <ImpactAnalysis 
            changes={pendingChanges}
            affectedKnowledge={affectedKnowledge}
            onAnalyze={handleAnalyzeImpact}
          />
          
          <TestGeneration 
            knowledge={knowledgeContext}
            onGenerateTests={handleGenerateTests}
          />
        </KnowledgePanel>
      </SplitPane>
      
      {/* Bottom panel: CLI emulator */}
      <CLIEmulator
        commands={developerCommands}
        onExecute={handleCLICommand}
        output={commandOutput}
      />
    </PersonaLayout>
  );
};
```

**Developer UI Patterns**:
- **Dual-pane interfaces**: Code alongside knowledge context
- **Annotation overlays**: Inline knowledge highlighting in code
- **Impact visualization**: Graph views of dependencies
- **CLI integration**: Terminal access for power users

### Product Owner Dashboard
```typescript
// Product owner focus: Value flow and capability planning
const ProductOwnerDashboard = () => {
  const { capabilities, valueMetrics } = useProduct();
  
  return (
    <PersonaLayout persona="product-owner">
      {/* Value flow visualization */}
      <ValueFlowChart 
        capabilities={capabilities}
        metrics={valueMetrics}
        onSelectCapability={handleSelectCapability}
      />
      
      {/* Capability planning */}
      <div className="grid grid-cols-2 gap-6 mt-6">
        <CapabilityPlanner 
          capabilities={capabilities}
          onPrioritize={handlePrioritize}
          onDecompose={handleDecompose}
        />
        
        <ROITracker 
          investments={investments}
          returns={returns}
          forecasts={forecasts}
        />
      </div>
      
      {/* Stakeholder alignment */}
      <StakeholderAlignment 
        stakeholders={stakeholders}
        requirements={requirements}
        feedback={feedback}
      />
    </PersonaLayout>
  );
};
```

**Product Owner UI Patterns**:
- **Value flow visualization**: ROI tracking across capabilities
- **Capability decomposition**: Natural language to task breakdown
- **Stakeholder alignment**: Requirement mapping and prioritization
- **Business metrics**: KPI dashboards with forecasting

### Project Manager Dashboard
```typescript
// Project manager focus: Dependency tracking and resource allocation
const ProjectManagerDashboard = () => {
  const { dependencies, resources } = useProject();
  
  return (
    <PersonaLayout persona="project-manager">
      {/* Dependency graph with time travel */}
      <DependencyGraph 
        dependencies={dependencies}
        timeRange={timeRange}
        onTimeTravel={handleTimeTravel}
      />
      
      {/* Resource allocation */}
      <div className="grid grid-cols-3 gap-6">
        <ResourceAllocation 
          resources={resources}
          onAllocate={handleAllocate}
          onOptimize={handleOptimize}
        />
        
        <BlockerPrediction 
          dependencies={dependencies}
          predictions={predictions}
          onMitigate={handleMitigate}
        />
        
        <TeamCoordination 
          teams={teams}
          assignments={assignments}
          onCoordinate={handleCoordinate}
        />
      </div>
      
      {/* Attention scheduling */}
      <AttentionScheduler 
        notifications={notifications}
        onSchedule={handleSchedule}
        onRoute={handleRoute}
      />
    </PersonaLayout>
  );
};
```

**Project Manager UI Patterns**:
- **Time-travel graphs**: Past/present/future dependency views
- **Resource visualization**: AI token, human attention, compute budgets
- **Blocker prediction**: ML-based risk identification
- **Attention scheduling**: Intelligent notification routing

### QA Engineer Dashboard
```typescript
// QA engineer focus: Quality gates and validation
const QAEngineerDashboard = () => {
  const { qualityGates, testSuites } = useQA();
  
  return (
    <PersonaLayout persona="qa-engineer">
      {/* Quality gates dashboard */}
      <QualityGatesDashboard 
        gates={qualityGates}
        status={gateStatus}
        onConfigure={handleConfigureGate}
      />
      
      {/* Test generation and management */}
      <div className="grid grid-cols-2 gap-6">
        <TestGenerator 
          knowledge={knowledgeForTesting}
          onGenerate={handleGenerateTests}
          coverage={testCoverage}
        />
        
        <RegressionPredictor 
          changes={recentChanges}
          predictions={regressionPredictions}
          onValidate={handleValidatePrediction}
        />
      </div>
      
      {/* Compliance validation */}
      <ComplianceValidator 
        rules={complianceRules}
        violations={violations}
        onApprove={handleApproveViolation}
        onBlock={handleBlockViolation}
      />
    </PersonaLayout>
  );
};
```

**QA Engineer UI Patterns**:
- **Quality gate configuration**: Visual editor for acceptance criteria
- **Test generation**: AI-assisted test creation from knowledge
- **Regression prediction**: Risk visualization and validation
- **Compliance tracking**: Real-time violation detection

---

## Atomic Knowledge Components

### Atom Display Component
```typescript
// Generic atom display with type-specific rendering
const AtomDisplay = ({ atom, editable = false, onEdit }) => {
  const Component = atomComponents[atom.type];
  
  if (!Component) {
    return <FallbackAtomDisplay atom={atom} />;
  }
  
  return (
    <div className={`atom-display atom-type-${atom.type}`}>
      <Component 
        atom={atom}
        editable={editable}
        onEdit={onEdit}
      />
      
      {/* Metadata */}
      <AtomMetadata 
        confidence={atom.confidence}
        source={atom.source}
        owner={atom.owner}
        updatedAt={atom.updatedAt}
      />
      
      {/* Actions */}
      {editable && (
        <AtomActions
          onEdit={onEdit}
          onDelete={() => handleDelete(atom.id)}
          onCopy={() => handleCopy(atom)}
        />
      )}
    </div>
  );
};

// Type-specific atom components
const atomComponents = {
  TextAtom: ({ atom, editable, onEdit }) => (
    <div className="text-atom">
      {editable ? (
        <TextEditor 
          value={atom.value}
          onChange={onEdit}
          placeholder="Enter text..."
        />
      ) : (
        <div className="text-content">{atom.value}</div>
      )}
    </div>
  ),
  
  NumberAtom: ({ atom, editable, onEdit }) => (
    <div className="number-atom">
      {editable ? (
        <NumberInput 
          value={atom.value}
          onChange={onEdit}
          min={atom.min}
          max={atom.max}
          step={atom.step}
        />
      ) : (
        <div className="number-value">
          {atom.value}
          {atom.unit && <span className="unit">{atom.unit}</span>}
        </div>
      )}
    </div>
  ),
  
  // ... more atom types
};
```

### Molecule Editor Component
```typescript
// Visual editor for molecules (groups of atoms)
const MoleculeEditor = ({ molecule, onChange }) => {
  const [atoms, setAtoms] = useState(molecule.atoms);
  const [structure, setStructure] = useState(molecule.structure);
  
  const handleAddAtom = (atomType) => {
    const newAtom = createAtom(atomType);
    const newAtoms = [...atoms, newAtom];
    setAtoms(newAtoms);
    onChange({ ...molecule, atoms: newAtoms });
  };
  
  const handleConnectAtoms = (sourceId, targetId, relationship) => {
    const newStructure = addConnection(structure, sourceId, targetId, relationship);
    setStructure(newStructure);
    onChange({ ...molecule, structure: newStructure });
  };
  
  return (
    <div className="molecule-editor">
      {/* Atom palette */}
      <AtomPalette 
        availableTypes={molecule.allowedAtomTypes}
        onSelect={handleAddAtom}
      />
      
      {/* Visual structure editor */}
      <div className="structure-canvas">
        {atoms.map(atom => (
          <DraggableAtom 
            key={atom.id}
            atom={atom}
            position={getPosition(structure, atom.id)}
            onDrag={(position) => updatePosition(structure, atom.id, position)}
          />
        ))}
        
        {/* Connection lines */}
        <svg className="connection-layer">
          {structure.connections.map(conn => (
            <ConnectionLine 
              key={conn.id}
              source={conn.source}
              target={conn.target}
              type={conn.type}
            />
          ))}
        </svg>
      </div>
      
      {/* Properties panel */}
      <MoleculeProperties 
        molecule={molecule}
        onChange={onChange}
      />
    </div>
  );
};
```

### Organism Viewer Component
```typescript
// Complete knowledge unit viewer
const OrganismViewer = ({ organism, mode = 'view' }) => {
  const [view, setView] = useState('composite'); // 'composite' | 'graph' | 'document'
  
  return (
    <div className="organism-viewer">
      {/* View mode selector */}
      <ViewModeSelector 
        modes={['composite', 'graph', 'document', 'code']}
        current={view}
        onChange={setView}
      />
      
      {/* Main content based on view mode */}
      {view === 'composite' && (
        <CompositeView 
          organism={organism}
          onSelectAtom={handleSelectAtom}
          onEditMolecule={handleEditMolecule}
        />
      )}
      
      {view === 'graph' && (
        <GraphView 
          organism={organism}
          onNodeClick={handleNodeClick}
          onRelationshipClick={handleRelationshipClick}
        />
      )}
      
      {view === 'document' && (
        <DocumentView 
          organism={organism}
          editable={mode === 'edit'}
          onEdit={handleDocumentEdit}
        />
      )}
      
      {view === 'code' && (
        <CodeView 
          implementations={organism.implementations}
          onNavigateToCode={handleNavigateToCode}
        />
      )}
      
      {/* Side panel with details */}
      <OrganismDetails 
        organism={organism}
        relationships={organism.relationships}
        validation={organism.validation}
      />
    </div>
  );
};
```

---

## Collaboration Components

### Real-time Presence Indicator
```typescript
// Show who's viewing/editing a knowledge unit
const PresenceIndicator = ({ knowledgeId }) => {
  const { others, self } = usePresence(knowledgeId);
  
  return (
    <div className="presence-indicator">
      {/* Self indicator */}
      <UserAvatar 
        user={self}
        isSelf={true}
        status="online"
      />
      
      {/* Others viewing */}
      {others.map(user => (
        <UserAvatar 
          key={user.id}
          user={user}
          status={user.status}
          activity={user.activity}
          tooltip={`${user.name} (${user.activity})`}
        />
      ))}
      
      {/* Collaboration tools */}
      <CollaborationTools
        onInvite={handleInvite}
        onStartCall={handleStartCall}
        onShareScreen={handleShareScreen}
      />
    </div>
  );
};
```

### Collaborative Editor
```typescript
// Real-time collaborative editing using CRDTs
const CollaborativeEditor = ({ documentId, initialContent }) => {
  const { yDoc, provider, awareness } = useYjs(documentId);
  const [content, setContent] = useState(initialContent);
  
  // Sync with Yjs
  useEffect(() => {
    const yText = yDoc.getText('content');
    
    // Update local state when remote changes
    const updateHandler = () => {
      setContent(yText.toString());
    };
    
    yText.observe(updateHandler);
    
    // Cleanup
    return () => {
      yText.unobserve(updateHandler);
    };
  }, [yDoc]);
  
  const handleChange = (newContent) => {
    const yText = yDoc.getText('content');
    
    // Calculate diff and apply to Yjs
    const diff = calculateDiff(content, newContent);
    applyDiffToYText(yText, diff);
    
    setContent(newContent);
  };
  
  return (
    <div className="collaborative-editor">
      {/* Presence sidebar */}
      <PresenceIndicator 
        awareness={awareness}
        provider={provider}
      />
      
      {/* Editor */}
      <RichTextEditor 
        value={content}
        onChange={handleChange}
        extensions={[
          Collaboration.configure({
            document: yDoc
          }),
          CollaborationCursor.configure({
            provider,
            user: awareness.getLocalState()
          })
        ]}
      />
      
      {/* Version history */}
      <VersionHistory 
        documentId={documentId}
        onRevert={handleRevert}
      />
    </div>
  );
};
```

### Comment Thread System
```typescript
// Contextual commenting on knowledge units
const CommentThread = ({ targetId, targetType }) => {
  const { comments, addComment, resolveComment } = useComments(targetId);
  const [newComment, setNewComment] = useState('');
  
  const handleSubmit = async () => {
    if (!newComment.trim()) return;
    
    await addComment({
      targetId,
      targetType,
      content: newComment,
      mentions: extractMentions(newComment),
      references: extractReferences(newComment)
    });
    
    setNewComment('');
  };
  
  return (
    <div className="comment-thread">
      {/* Comment list */}
      <div className="comments-list">
        {comments.map(comment => (
          <Comment 
            key={comment.id}
            comment={comment}
            onReply={handleReply}
            onResolve={() => resolveComment(comment.id)}
          />
        ))}
      </div>
      
      {/* New comment input */}
      <CommentInput 
        value={newComment}
        onChange={setNewComment}
        onSubmit={handleSubmit}
        onMention={handleMention}
        onReference={handleReference}
      />
    </div>
  );
};
```

---

## Visualization Components

### Knowledge Graph Visualization
```typescript
// Interactive graph visualization of knowledge relationships
const KnowledgeGraph = ({ 
  nodes, 
  edges, 
  onNodeClick, 
  onEdgeClick 
}) => {
  const [layout, setLayout] = useState('force');
  const [filter, setFilter] = useState({});
  const [selected, setSelected] = useState(null);
  
  const filteredNodes = applyFilter(nodes, filter);
  const filteredEdges = applyFilter(edges, filter);
  
  return (
    <div className="knowledge-graph">
      {/* Controls */}
      <GraphControls 
        layout={layout}
        onLayoutChange={setLayout}
        filter={filter}
        onFilterChange={setFilter}
      />
      
      {/* Visualization */}
      <CytoscapeComponent
        elements={CytoscapeComponent.normalizeElements({
          nodes: filteredNodes.map(node => ({
            data: {
              id: node.id,
              label: node.label,
              type: node.type,
              ...node.properties
            }
          })),
          edges: filteredEdges.map(edge => ({
            data: {
              id: edge.id,
              source: edge.source,
              target: edge.target,
              label: edge.type,
              ...edge.properties
            }
          }))
        })}
        layout={{
          name: layout,
          ...layoutConfigs[layout]
        }}
        stylesheet={graphStyles}
        cy={(cy) => {
          cy.on('tap', 'node', (evt) => {
            const node = evt.target;
            setSelected(node.data());
            onNodeClick?.(node.data());
          });
          
          cy.on('tap', 'edge', (evt) => {
            const edge = evt.target;
            onEdgeClick?.(edge.data());
          });
        }}
      />
      
      {/* Details panel */}
      {selected && (
        <GraphDetailsPanel 
          element={selected}
          onClose={() => setSelected(null)}
        />
      )}
    </div>
  );
};
```

### Value Flow Chart
```typescript
// Business value visualization
const ValueFlowChart = ({ capabilities, metrics }) => {
  const [timeRange, setTimeRange] = useState('month');
  const [groupBy, setGroupBy] = useState('capability');
  
  const data = transformData(capabilities, metrics, {
    timeRange,
    groupBy
  });
  
  return (
    <div className="value-flow-chart">
      {/* Header with controls */}
      <div className="chart-header">
        <h3>Value Flow Analysis</h3>
        <div className="controls">
          <TimeRangeSelector 
            value={timeRange}
            onChange={setTimeRange}
            options={['week', 'month', 'quarter', 'year']}
          />
          
          <GroupBySelector 
            value={groupBy}
            onChange={setGroupBy}
            options={['capability', 'team', 'persona', 'domain']}
          />
        </div>
      </div>
      
      {/* Chart */}
      <ResponsiveContainer width="100%" height={400}>
        <ComposedChart data={data}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="name" />
          <YAxis />
          <Tooltip />
          <Legend />
          
          {/* Investment bars */}
          <Bar 
            dataKey="investment" 
            fill="#8884d8" 
            name="Investment"
          />
          
          {/* Return line */}
          <Line 
            type="monotone" 
            dataKey="return" 
            stroke="#82ca9d" 
            name="Return"
          />
          
          {/* ROI scatter */}
          <Scatter 
            dataKey="roi" 
            fill="#ff7300" 
            name="ROI"
          />
        </ComposedChart>
      </ResponsiveContainer>
      
      {/* Insights */}
      <ValueInsights 
        data={data}
        metrics={calculateMetrics(data)}
      />
    </div>
  );
};
```

### Dependency Graph with Time Travel
```typescript
// Time-aware dependency visualization
const TimeTravelDependencyGraph = ({ dependencies }) => {
  const [currentTime, setCurrentTime] = useState('present');
  const [viewMode, setViewMode] = useState('graph'); // 'graph' | 'timeline' | 'gantt'
  
  const filteredDeps = filterByTime(dependencies, currentTime);
  
  return (
    <div className="time-travel-graph">
      {/* Time controls */}
      <TimeControls 
        currentTime={currentTime}
        onChange={setCurrentTime}
        modes={['past', 'present', 'future']}
      />
      
      {/* Main visualization */}
      {viewMode === 'graph' && (
        <DependencyGraph 
          dependencies={filteredDeps}
          timeContext={currentTime}
        />
      )}
      
      {viewMode === 'timeline' && (
        <TimelineView 
          dependencies={dependencies}
          currentTime={currentTime}
          onSelectTime={setCurrentTime}
        />
      )}
      
      {viewMode === 'gantt' && (
        <GanttChart 
          dependencies={dependencies}
          currentTime={currentTime}
        />
      )}
      
      {/* Prediction panel for future view */}
      {currentTime === 'future' && (
        <PredictionPanel 
          dependencies={dependencies}
          predictions={generatePredictions(dependencies)}
        />
      )}
    </div>
  );
};
```

---

## State Management Patterns

### Knowledge Store (Zustand)
```typescript
// Centralized knowledge state management
interface KnowledgeState {
  // State
  atoms: Record<string, Atom>;
  molecules: Record<string, Molecule>;
  organisms: Record<string, Organism>;
  selectedId: string | null;
  loading: boolean;
  error: string | null;
  
  // Actions
  fetchKnowledge: (id: string) => Promise<void>;
  updateAtom: (id: string, updates: Partial<Atom>) => void;
  createMolecule: (molecule: Molecule) => Promise<string>;
  deleteOrganism: (id: string) => Promise<void>;
  selectKnowledge: (id: string | null) => void;
  
  // Derived state
  selectedKnowledge: () => KnowledgeUnit | null;
  relatedKnowledge: () => KnowledgeUnit[];
  validationStatus: () => ValidationStatus;
}

const useKnowledgeStore = create<KnowledgeState>()((set, get) => ({
  // Initial state
  atoms: {},
  molecules: {},
  organisms: {},
  selectedId: null,
  loading: false,
  error: null,
  
  // Actions
  fetchKnowledge: async (id: string) => {
    set({ loading: true, error: null });
    
    try {
      const knowledge = await knowledgeAPI.fetch(id);
      
      set(state => ({
        loading: false,
        atoms: { ...state.atoms, ...knowledge.atoms },
        molecules: { ...state.molecules, ...knowledge.molecules },
        organisms: { ...state.organisms, [knowledge.id]: knowledge }
      }));
    } catch (error) {
      set({ loading: false, error: error.message });
    }
  },
  
  updateAtom: (id: string, updates: Partial<Atom>) => {
    set(state => ({
      atoms: {
        ...state.atoms,
        [id]: { ...state.atoms[id], ...updates }
      }
    }));
    
    // Real-time update
    knowledgeAPI.updateAtom(id, updates);
  },
  
  // ... more actions
  
  // Derived state (computed)
  selectedKnowledge: () => {
    const state = get();
    if (!state.selectedId) return null;
    
    // Find in appropriate collection
    if (state.atoms[state.selectedId]) return state.atoms[state.selectedId];
    if (state.molecules[state.selectedId]) return state.molecules[state.selectedId];
    if (state.organisms[state.selectedId]) return state.organisms[state.selectedId];
    
    return null;
  },
  
  relatedKnowledge: () => {
    const selected = get().selectedKnowledge();
    if (!selected) return [];
    
    // Find relationships
    const relatedIds = selected.relationships?.map(r => r.targetId) || [];
    
    return relatedIds.map(id => 
      get().atoms[id] || get().molecules[id] || get().organisms[id]
    ).filter(Boolean);
  }
}));
```

### Persona Preference Store
```typescript
// Persona-specific UI preferences
interface PersonaPreferences {
  // UI preferences
  layout: 'grid' | 'list' | 'detailed';
  density: 'compact' | 'normal' | 'comfortable';
  theme: 'light' | 'dark' | 'auto';
  
  // Feature preferences
  enabledFeatures: string[];
  hiddenPanels: string[];
  shortcutOverrides: Record<string, string>;
  
  // Domain preferences
  defaultDomain: string;
  favoriteTemplates: string[];
  expertNetworks: string[];
  
  // Actions
  setLayout: (layout: PersonaPreferences['layout']) => void;
  toggleFeature: (feature: string) => void;
  addFavoriteTemplate: (template: string) => void;
  resetToDefaults: () => void;
}

const usePersonaStore = create<PersonaPreferences>()((set) => ({
  // Default preferences
  layout: 'grid',
  density: 'normal',
  theme: 'auto',
  enabledFeatures: ['validationQueue', 'ruleInjection', 'impactAnalysis'],
  hiddenPanels: [],
  shortcutOverrides: {},
  defaultDomain: 'tax',
  favoriteTemplates: ['TaxRuleTemplate', 'DecisionTemplate'],
  expertNetworks: ['tax-experts', 'compliance-officers'],
  
  // Actions
  setLayout: (layout) => set({ layout }),
  toggleFeature: (feature) => set(state => ({
    enabledFeatures: state.enabledFeatures.includes(feature)
      ? state.enabledFeatures.filter(f => f !== feature)
      : [...state.enabledFeatures, feature]
  })),
  addFavoriteTemplate: (template) => set(state => ({
    favoriteTemplates: [...state.favoriteTemplates, template]
  })),
  resetToDefaults: () => set({
    layout: 'grid',
    density: 'normal',
    theme: 'auto',
    enabledFeatures: ['validationQueue', 'ruleInjection', 'impactAnalysis'],
    hiddenPanels: [],
    shortcutOverrides: {},
    defaultDomain: 'tax',
    favoriteTemplates: ['TaxRuleTemplate', 'DecisionTemplate'],
    expertNetworks: ['tax-experts', 'compliance-officers']
  })
}));
```

---

## Data Fetching Patterns

### TanStack Query for Knowledge Data
```typescript
// React Query hooks for knowledge data
const useKnowledgeQueries = () => {
  // Fetch atom with automatic caching
  const atomQuery = useQuery({
    queryKey: ['atoms', atomId],
    queryFn: () => knowledgeAPI.fetchAtom(atomId),
    staleTime: 5 * 60 * 1000, // 5 minutes
    cacheTime: 30 * 60 * 1000, // 30 minutes
  });
  
  // Fetch related molecules
  const moleculesQuery = useQueries({
    queries: atomQuery.data?.moleculeIds?.map(moleculeId => ({
      queryKey: ['molecules', moleculeId],
      queryFn: () => knowledgeAPI.fetchMolecule(moleculeId),
      enabled: !!atomQuery.data
    })) || []
  });
  
  // Mutations for updates
  const updateMutation = useMutation({
    mutationFn: (updates: AtomUpdate) => 
      knowledgeAPI.updateAtom(atomId, updates),
    onSuccess: () => {
      // Invalidate queries
      queryClient.invalidateQueries(['atoms', atomId]);
    },
    onError: (error) => {
      toast.error(`Update failed: ${error.message}`);
    }
  });
  
  // Real-time subscription
  useSubscription({
    channel: `atoms:${atomId}`,
    onMessage: (message) => {
      queryClient.setQueryData(['atoms', atomId], 
        (old) => mergeUpdates(old, message)
      );
    }
  });
  
  return {
    atom: atomQuery,
    molecules: moleculesQuery,
    updateAtom: updateMutation.mutate,
    isLoading: atomQuery.isLoading || moleculesQuery.some(q => q.isLoading),
    error: atomQuery.error || moleculesQuery.find(q => q.error)?.error
  };
};
```

### Optimistic Updates for Collaboration
```typescript
// Optimistic updates for real-time collaboration
const useOptimisticKnowledgeUpdate = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async (update: KnowledgeUpdate) => {
      // Send update to server
      return await knowledgeAPI.update(update);
    },
    
    onMutate: async (update) => {
      // Cancel outgoing refetches
      await queryClient.cancelQueries(['knowledge', update.id]);
      
      // Snapshot previous value
      const previous = queryClient.getQueryData(['knowledge', update.id]);
      
      // Optimistically update
      queryClient.setQueryData(['knowledge', update.id], 
        (old: any) => mergeUpdate(old, update)
      );
      
      // Return context for rollback
      return { previous };
    },
    
    onError: (error, update, context) => {
      // Rollback on error
      queryClient.setQueryData(['knowledge', update.id], context?.previous);
      toast.error(`Update failed: ${error.message}`);
    },
    
    onSettled: () => {
      // Refetch to ensure consistency
      queryClient.invalidateQueries(['knowledge']);
    }
  });
};
```

---

## Accessibility Implementation

### Accessible Component Library
```typescript
// Base accessible components
const AccessibleButton = React.forwardRef(({ children, ...props }, ref) => {
  return (
    <button
      ref={ref}
      {...props}
      className={clsx('accessible-button', props.className)}
      // ARIA attributes
      aria-pressed={props['aria-pressed']}
      aria-label={props['aria-label'] || getAccessibleLabel(children)}
      // Keyboard navigation
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault();
          props.onClick?.(e);
        }
        props.onKeyDown?.(e);
      }}
    >
      {children}
    </button>
  );
});

const AccessibleDialog = ({ isOpen, onClose, title, children }) => {
  const dialogRef = useRef<HTMLDialogElement>(null);
  
  useEffect(() => {
    if (isOpen) {
      dialogRef.current?.showModal();
      // Trap focus
      trapFocus(dialogRef.current);
    } else {
      dialogRef.current?.close();
    }
  }, [isOpen]);
  
  return (
    <dialog
      ref={dialogRef}
      className="accessible-dialog"
      aria-labelledby="dialog-title"
      aria-modal="true"
      onClose={onClose}
    >
      <div className="dialog-header">
        <h2 id="dialog-title">{title}</h2>
        <button
          className="close-button"
          onClick={onClose}
          aria-label="Close dialog"
        >
          ×
        </button>
      </div>
      
      <div className="dialog-content">
        {children}
      </div>
    </dialog>
  );
};
```

### Screen Reader Announcements
```typescript
// Dynamic announcements for screen readers
const LiveRegion = () => {
  const [announcements, setAnnouncements] = useState<string[]>([]);
  
  // Subscribe to announcement events
  useEffect(() => {
    const handleAnnouncement = (event: CustomEvent) => {
      setAnnouncements(prev => [...prev, event.detail.message]);
      
      // Remove after announcement
      setTimeout(() => {
        setAnnouncements(prev => prev.filter(m => m !== event.detail.message));
      }, 5000);
    };
    
    window.addEventListener('announcement', handleAnnouncement);
    return () => window.removeEventListener('announcement', handleAnnouncement);
  }, []);
  
  return (
    <div 
      className="live-region"
      aria-live="polite"
      aria-atomic="false"
      // Visually hidden but available to screen readers
      style={{
        position: 'absolute',
        width: '1px',
        height: '1px',
        padding: 0,
        margin: '-1px',
        overflow: 'hidden',
        clip: 'rect(0, 0, 0, 0)',
        whiteSpace: 'nowrap',
        border: 0
      }}
    >
      {announcements.map((message, index) => (
        <div key={index}>{message}</div>
      ))}
    </div>
  );
};

// Utility to trigger announcements
const announce = (message: string, priority: 'polite' | 'assertive' = 'polite') => {
  const event = new CustomEvent('announcement', {
    detail: { message, priority }
  });
  window.dispatchEvent(event);
};
```

---

## Performance Optimization

### Virtualized Lists for Large Datasets
```typescript
// Virtualized list for thousands of knowledge items
const VirtualizedKnowledgeList = ({ items, renderItem }) => {
  const listRef = useRef();
  const { width, height } = useResizeObserver(listRef);
  
  const rowVirtualizer = useVirtualizer({
    count: items.length,
    getScrollElement: () => listRef.current,
    estimateSize: () => 60, // Estimated row height
    overscan: 5 // Render extra items for smooth scrolling
  });
  
  return (
    <div 
      ref={listRef}
      className="virtualized-list"
      style={{ height: '600px', overflow: 'auto' }}
    >
      <div
        style={{
          height: `${rowVirtualizer.getTotalSize()}px`,
          width: '100%',
          position: 'relative'
        }}
      >
        {rowVirtualizer.getVirtualItems().map(virtualRow => (
          <div
            key={virtualRow.key}
            style={{
              position: 'absolute',
              top: 0,
              left: 0,
              width: '100%',
              height: `${virtualRow.size}px`,
              transform: `translateY(${virtualRow.start}px)`
            }}
          >
            {renderItem(items[virtualRow.index], virtualRow.index)}
          </div>
        ))}
      </div>
    </div>
  );
};
```

### Code Splitting by Persona
```typescript
// Dynamic imports for persona-specific code
const PersonaRouter = () => {
  const { persona } = usePersona();
  
  // Dynamic imports for persona dashboards
  const PersonaDashboard = lazy(() => {
    switch (persona) {
      case 'domain-expert':
        return import('./personas/domain-expert/Dashboard');
      case 'developer':
        return import('./personas/developer/Workspace');
      case 'product-owner':
        return import('./personas/product-owner/Dashboard');
      case 'project-manager':
        return import('./personas/project-manager/Dashboard');
      case 'qa-engineer':
        return import('./personas/qa-engineer/Dashboard');
      default:
        return import('./personas/default/Dashboard');
    }
  });
  
  return (
    <Suspense fallback={<PersonaLoadingSkeleton />}>
      <PersonaDashboard />
    </Suspense>
  );
};

// Route-based code splitting
const routes = [
  {
    path: '/domain-expert/*',
    element: lazy(() => import('./personas/domain-expert/Routes'))
  },
  {
    path: '/developer/*',
    element: lazy(() => import('./personas/developer/Routes'))
  },
  // ... more persona routes
];
```

---

## Testing Strategy

### Component Testing
```typescript
// Component tests with Testing Library
describe('AtomDisplay', () => {
  test('renders TextAtom correctly', () => {
    const atom: Atom = {
      id: 'atom-123',
      type: 'TextAtom',
      value: 'Canadian GST is 5%',
      confidence: 0.9
    };
    
    render(<AtomDisplay atom={atom} />);
    
    expect(screen.getByText('Canadian GST is 5%')).toBeInTheDocument();
    expect(screen.getByLabelText('Confidence: 90%')).toBeInTheDocument();
  });
  
  test('allows editing when editable prop is true', () => {
    const atom: Atom = {
      id: 'atom-123',
      type: 'TextAtom',
      value: 'Initial value'
    };
    
    const onEdit = jest.fn();
    
    render(<AtomDisplay atom={atom} editable={true} onEdit={onEdit} />);
    
    const editButton = screen.getByRole('button', { name: /edit/i });
    fireEvent.click(editButton);
    
    const input = screen.getByRole('textbox');
    fireEvent.change(input, { target: { value: 'Updated value' } });
    fireEvent.blur(input);
    
    expect(onEdit).toHaveBeenCalledWith('Updated value');
  });
});
```

### Integration Testing
```typescript
// End-to-end persona workflow tests
describe('Domain Expert Workflow', () => {
  test('completes knowledge validation workflow', async () => {
    // Setup test data
    const extraction: AtomicExtraction = {
      id: 'extraction-123',
      atoms: [/* test atoms */],
      confidence: 0.75
    };
    
    // Render validation interface
    render(<ValidationInterface task={extraction} />);
    
    // Expert reviews and approves
    const approveButton = screen.getByRole('button', { name: /approve/i });
    fireEvent.click(approveButton);
    
    // Add comment
    const commentInput = screen.getByPlaceholderText(/add validation comments/i);
    fireEvent.change(commentInput, { 
      target: { value: 'Looks accurate to me' } 
    });
    
    const submitButton = screen.getByRole('button', { name: /submit/i });
    fireEvent.click(submitButton);
    
    // Verify knowledge was updated
    await waitFor(() => {
      expect(mockKnowledgeAPI.update).toHaveBeenCalledWith(
        expect.objectContaining({
          status: 'approved',
          validatedBy: 'expert-123'
        })
      );
    });
  });
});
```

### Visual Regression Testing
```typescript
// Screenshot-based visual tests
describe('Visual Regression', () => {
  test('DomainExpertDashboard matches snapshot', async () => {
    const { container } = render(<DomainExpertDashboard />);
    
    // Wait for async content
    await screen.findByText(/validation queue/i);
    
    // Compare screenshot
    expect(container).toMatchImageSnapshot({
      customSnapshotIdentifier: 'domain-expert-dashboard',
      failureThreshold: 0.01,
      failureThresholdType: 'percent'
    });
  });
});
```

---

## Deployment & Monitoring

### Build Configuration
```javascript
// Vite configuration for optimized builds
export default defineConfig({
  plugins: [
    react(),
    // Code splitting by persona
    splitVendorChunkPlugin(),
    // Bundle visualization
    visualizer({
      filename: 'dist/stats.html'
    })
  ],
  
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          // Persona-specific chunks
          'domain-expert': ['./src/personas/domain-expert'],
          'developer': ['./src/personas/developer'],
          'product-owner': ['./src/personas/product-owner'],
          // Vendor chunks
          'react-vendor': ['react', 'react-dom'],
          'ui-vendor': ['@radix-ui/react-dialog', '@radix-ui/react-dropdown-menu'],
          'visualization-vendor': ['d3', 'cytoscape', 'three']
        }
      }
    },
    
    // Performance optimizations
    target: 'es2020',
    minify: 'terser',
    sourcemap: true,
    chunkSizeWarningLimit: 1000
  },
  
  // Environment variables
  define: {
    __APP_VERSION__: JSON.stringify(packageJson.version)
  }
});
```

### Performance Monitoring
```typescript
// Frontend performance monitoring
class FrontendMonitor {
  private vitals = {
    CLS: new PerformanceMetric('CLS'),
    FID: new PerformanceMetric('FID'),
    LCP: new PerformanceMetric('LCP'),
    FCP: new PerformanceMetric('FCP'),
    TTI: new PerformanceMetric('TTI')
  };
  
  private errors: ErrorMetric[] = [];
  private userActions: UserAction[] = [];
  
  constructor() {
    this.setupPerformanceObservers();
    this.setupErrorTracking();
    this.setupUserActionTracking();
  }
  
  private setupPerformanceObservers() {
    // Core Web Vitals
    const observer = new PerformanceObserver((entryList) => {
      for (const entry of entryList.getEntries()) {
        const metric = this.vitals[entry.name];
        if (metric) {
          metric.record(entry.value);
          
          // Report to analytics
          if (metric.isPoor()) {
            this.reportPoorMetric(entry.name, entry.value);
          }
        }
      }
    });
    
    observer.observe({ entryTypes: ['paint', 'largest-contentful-paint', 'layout-shift', 'first-input'] });
  }
  
  private setupErrorTracking() {
    window.addEventListener('error', (event) => {
      this.errors.push({
        message: event.message,
        stack: event.error?.stack,
        timestamp: new Date(),
        url: window.location.href
      });
      
      // Report to error tracking service
      errorTrackingService.report(event.error);
    });
    
    // Unhandled promise rejections
    window.addEventListener('unhandledrejection', (event) => {
      this.errors.push({
        message: event.reason?.message || 'Unhandled promise rejection',
        stack: event.reason?.stack,
        timestamp: new Date(),
        url: window.location.href
      });
    });
  }
}
```

---

## Cross-References

- **Rendering**: See [`product-os-implementation-rendering.md`](product-os-implementation-rendering.md) for standoff annotation rendering components
- **Persistence**: See [`product-os-implementation-persistence.md`](product-os-implementation-persistence.md) for data models and API clients
- **Search**: See [`product-os-implementation-search.md`](product-os-implementation-search.md) for search UI components
- **Extraction**: See [`product-os-implementation-extraction.md`](product-os-implementation-extraction.md) for validation interfaces

---

**Status**: Technical reference - frontend architecture defined  
**Next**: Review with frontend engineering team for implementation planning