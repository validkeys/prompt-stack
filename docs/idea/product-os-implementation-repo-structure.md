# Product-OS: Repository Structure Specification

**Date**: 2026-01-12  
**Status**: Implementation Specification  
**Purpose**: Formal specification for Product-OS repository structure and development setup, designed for AI-driven implementation.

**Related Documents**:
- [`product-os.md`](product-os.md) - Vision and concept
- [`product-os-data-structure.md`](product-os-data-structure.md) - Data structure concepts
- [`product-os-implementation-persistence.md`](product-os-implementation-persistence.md) - Database architecture
- [`product-os-implementation-extraction.md`](product-os-implementation-extraction.md) - Extraction pipeline

---

## 1. Repository Overview

### 1.1 Purpose
A TypeScript monorepo implementing Product-OS using modern development practices, service-command architecture, and polyglot database architecture.

### 1.2 Key Architectural Decisions
- **Monorepo**: TypeScript monorepo managed by TurboRepo
- **Service Architecture**: Service-command pattern using Contracted framework
- **Dependency Management**: No traditional DI container; dependencies provided directly to services
- **Persistence**: Modular monolith approach with shared persistence layer
- **Runtime**: Bun runtime with pnpm package manager
- **API**: Fastify-based API server as primary deployable application

### 1.3 Technology Stack
| Component | Technology | Purpose |
|-----------|------------|---------|
| **Language** | TypeScript | Primary development language |
| **Package Manager** | pnpm | Workspace and dependency management |
| **Runtime** | Bun | Execution and built-in test framework |
| **Build System** | TurboRepo | Monorepo task orchestration and caching |
| **Testing Framework** | Vite + Bun test | Testing and development builds |
| **API Framework** | Fastify | High-performance web server |
| **Service Framework** | Contracted | Contract-first service definitions |
| **Database ORM** | Kysely | Type-safe PostgreSQL queries |
| **Migration Tool** | Kysely CLI | Database schema management |
| **Type Generation** | kysely-autogen | Automatic TypeScript types from database |
| **Validation** | zod3 | Runtime type validation and schemas |
| **Infrastructure** | Docker-compose | Local development databases and services |
| **Namespace** | @know-s | All internal package namespace |

---

## 2. Repository Structure

### 2.1 Top-Level Directory Layout
```
/
├── apps/                    # Deployable applications
│   ├── api/                 # API server (Fastify) - PRIMARY
│   ├── ui/                  # React web application (FUTURE)
│   └── cli/                 # CLI tools (FUTURE)
├── services/                # Domain services
│   ├── knowledge/           # Knowledge management service
│   ├── extraction/          # Knowledge extraction pipeline
│   ├── search/              # Multi-level search service
│   └── [additional domains] # Other domain services
├── lib/                     # Internal libraries (modular monolith)
│   ├── contracts/           # Service contracts and shared interfaces
│   ├── persistence/         # Shared persistence layer
│   └── shared/              # Shared utilities and types
├── packages/                # Shared packages (optional)
├── docs/                    # Documentation
├── scripts/                 # Build and utility scripts
├── docker/                  # Docker configurations
├── .env.example            # Environment template
├── .env.local              # Local environment (gitignored)
├── .env.production         # Production environment (gitignored)
├── package.json            # Root package.json with workspaces
├── turbo.json              # TurboRepo configuration
├── tsconfig.base.json      # Base TypeScript configuration
├── docker-compose.yml      # Local development infrastructure
└── README.md               # Project overview
```

### 2.2 Apps Directory (`apps/`)

#### `apps/api/` - Primary API Server
```
apps/api/
├── src/
│   ├── app.ts              # Fastify application setup
│   ├── routes/             # API route definitions
│   ├── plugins/            # Fastify plugins
│   ├── hooks/              # Request/response hooks
│   └── index.ts            # Application entry point
├── test/                   # API-level E2E tests
├── package.json            # API-specific dependencies
├── tsconfig.json          # TypeScript configuration
└── README.md              # API documentation
```

**Responsibilities**:
- Dependency injection and service orchestration
- HTTP request/response handling
- Authentication and authorization middleware
- Service composition and lifecycle management
- E2E testing entry point

#### `apps/ui/` - Future Web UI (React)
```
apps/ui/
├── src/
│   ├── main.tsx           # Application entry point
│   ├── App.tsx            # Root component
│   ├── components/        # React components
│   ├── pages/             # Page components
│   ├── hooks/             # Custom React hooks
│   ├── store/             # State management
│   └── api/               # API client
├── public/                # Static assets
├── package.json           # UI-specific dependencies
├── tsconfig.json         # TypeScript configuration
├── vite.config.ts        # Vite configuration
└── README.md             # UI documentation
```

#### `apps/cli/` - Future CLI Tools
```
apps/cli/
├── src/
│   ├── index.ts          # CLI entry point
│   ├── commands/         # Command definitions
│   ├── utils/            # CLI utilities
│   └── types/            # CLI-specific types
├── package.json          # CLI-specific dependencies
├── tsconfig.json        # TypeScript configuration
└── README.md            # CLI documentation
```

### 2.3 Services Directory (`services/`)

#### Service Structure Template
```
services/[service-name]/
├── src/
│   ├── commands/         # Service command implementations
│   │   ├── [command-name].ts
│   │   └── index.ts
│   ├── types/            # Service-specific types
│   ├── utils/            # Service utilities
│   └── index.ts          # Service entry point
├── test/                 # Unit tests
│   ├── commands/         # Command tests
│   └── fixtures/         # Test fixtures
├── package.json          # Service dependencies
├── tsconfig.json        # TypeScript configuration
└── README.md            # Service documentation
```

**Service Requirements**:
- Each service must define its contract in `lib/contracts/`
- Services cannot directly import other services
- All inter-service communication through contracts
- Services receive dependencies via constructor/factory
- Unit tests must cover all commands

#### Example Services
- `services/knowledge/`: Knowledge graph management, CRUD operations, graph queries
- `services/extraction/`: AI-driven knowledge extraction from various sources
- `services/search/`: Multi-level search across all data dimensions
- `services/auth/`: Authentication and authorization
- `services/notification/`: Event notifications and subscriptions

### 2.4 Libraries Directory (`lib/`)

#### `lib/contracts/` - Service Contracts
```
lib/contracts/
├── knowledge/            # Knowledge service contracts
│   ├── commands/        # Command definitions
│   ├── types/           # Shared types
│   └── index.ts         # Contract exports
├── extraction/          # Extraction service contracts
├── search/              # Search service contracts
├── shared/              # Cross-service contracts
├── index.ts             # All contracts export
└── README.md            # Contracts documentation
```

**Contract Pattern**:
```typescript
// Example contract definition using Contracted framework
import { createCommand } from '@contracted/core';

export const createKnowledgeItem = createCommand({
  name: 'createKnowledgeItem',
  input: z.object({
    title: z.string(),
    content: z.string(),
    type: z.enum(['note', 'document', 'code']),
  }),
  output: z.object({
    id: z.string(),
    createdAt: z.date(),
  }),
  dependencies: {
    persistence: 'PersistenceClient',
    validation: 'ValidationService',
  },
});
```

#### `lib/persistence/` - Shared Persistence Layer
```
lib/persistence/
├── src/
│   ├── postgres/        # PostgreSQL client and queries
│   │   ├── client.ts    # Kysely client setup
│   │   ├── migrations/  # Database migrations
│   │   ├── queries/     # Query definitions
│   │   └── types/       # Generated database types
│   ├── neo4j/           # Neo4j graph database client
│   ├── redis/           # Redis cache client
│   ├── storage/         # Object storage (S3-compatible)
│   ├── client.ts        # Unified persistence client
│   └── index.ts         # Persistence exports
├── migrations/          # Database migration scripts
├── package.json         # Persistence dependencies
├── tsconfig.json       # TypeScript configuration
└── README.md           # Persistence documentation
```

**Persistence Requirements**:
- Single shared persistence client for all services
- Connection pooling and lifecycle management
- Transaction management across multiple databases
- Migration orchestration
- Type safety through Kysely and kysely-autogen

#### `lib/shared/` - Shared Utilities
```
lib/shared/
├── src/
│   ├── utils/          # General utilities
│   ├── errors/         # Error definitions and handling
│   ├── logging/        # Logging configuration
│   ├── config/         # Configuration utilities
│   ├── events/         # Event system
│   └── types/          # Shared TypeScript types
├── package.json        # Shared dependencies
├── tsconfig.json      # TypeScript configuration
└── README.md          # Shared library documentation
```

---

## 3. Build & Development Setup

### 3.1 Package Management

#### Root `package.json`
```json
{
  "name": "@know-s/product-os",
  "private": true,
  "workspaces": ["apps/*", "services/*", "lib/*"],
  "scripts": {
    "build": "turbo run build",
    "build:transpile": "turbo run build:transpile",
    "typecheck": "turbo run typecheck",
    "test": "turbo run test",
    "test:watch": "turbo run test:watch",
    "lint": "turbo run lint",
    "format": "turbo run format",
    "dev": "turbo run dev",
    "db:up": "docker-compose up -d",
    "db:down": "docker-compose down",
    "db:migrate": "turbo run db:migrate",
    "db:generate": "turbo run db:generate"
  },
  "devDependencies": {
    "@types/bun": "latest",
    "@types/node": "latest",
    "@typescript-eslint/eslint-plugin": "latest",
    "@typescript-eslint/parser": "latest",
    "eslint": "latest",
    "prettier": "latest",
    "turbo": "latest",
    "typescript": "latest"
  },
  "packageManager": "pnpm@latest",
  "engines": {
    "bun": ">=1.0.0"
  }
}
```

#### Package-Specific `package.json`
Each app, service, and library must have its own `package.json` with:
- Appropriate `name` following `@know-s/[type]-[name]` pattern
- Dependencies specific to the package
- Build and test scripts
- Reference to shared TypeScript configuration

### 3.2 TypeScript Configuration

#### Base Configuration (`tsconfig.base.json`)
```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "ESNext",
    "lib": ["ES2022"],
    "moduleResolution": "bundler",
    "strict": true,
    "skipLibCheck": true,
    "esModuleInterop": true,
    "resolveJsonModule": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,
    "outDir": "./dist",
    "rootDir": "./src",
    "baseUrl": ".",
    "paths": {
      "@know-s/contracts/*": ["../lib/contracts/src/*"],
      "@know-s/persistence/*": ["../lib/persistence/src/*"],
      "@know-s/shared/*": ["../lib/shared/src/*"]
    }
  },
  "exclude": ["node_modules", "dist", "**/*.test.ts", "**/*.spec.ts"]
}
```

#### Package-Specific Overrides
Each package extends the base configuration with package-specific paths and settings.

### 3.3 TurboRepo Configuration (`turbo.json`)
```json
{
  "$schema": "https://turbo.build/schema.json",
  "globalDependencies": ["**/.env.*", "tsconfig.base.json"],
  "pipeline": {
    "build": {
      "dependsOn": ["^build"],
      "outputs": ["dist/**"]
    },
    "build:transpile": {
      "dependsOn": ["^build:transpile"],
      "outputs": ["dist/**"]
    },
    "typecheck": {
      "dependsOn": ["^typecheck"]
    },
    "test": {
      "dependsOn": ["^build"],
      "outputs": []
    },
    "lint": {
      "outputs": []
    },
    "format": {
      "outputs": []
    },
    "dev": {
      "cache": false,
      "persistent": true
    },
    "db:migrate": {
      "cache": false
    },
    "db:generate": {
      "cache": false
    }
  }
}
```

### 3.4 Development Infrastructure (`docker-compose.yml`)
```yaml
version: '3.8'
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_DB: product_os
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./lib/persistence/migrations:/docker-entrypoint-initdb.d

  neo4j:
    image: neo4j:5
    environment:
      NEO4J_AUTH: neo4j/password
      NEO4J_PLUGINS: '["apoc"]'
    ports:
      - "7474:7474"
      - "7687:7687"
    volumes:
      - neo4j_data:/data

  redis:
    image: redis:7
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  minio:
    image: minio/minio
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio_data:/data

volumes:
  postgres_data:
  neo4j_data:
  redis_data:
  minio_data:
```

---

## 4. Service Architecture

### 4.1 Contracted Framework Pattern

#### Service Contract Definition
```typescript
// lib/contracts/knowledge/commands/create-knowledge-item.ts
import { createCommand } from '@contracted/core';
import { z } from 'zod';

export const createKnowledgeItem = createCommand({
  name: 'createKnowledgeItem',
  input: z.object({
    title: z.string().min(1).max(255),
    content: z.string().min(1),
    type: z.enum(['note', 'document', 'code', 'image']),
    metadata: z.record(z.any()).optional(),
  }),
  output: z.object({
    id: z.string().uuid(),
    createdAt: z.date(),
    updatedAt: z.date(),
  }),
  dependencies: {
    db: 'PersistenceClient',
    validate: 'ValidationService',
    events: 'EventService',
  },
  async execute(input, deps) {
    // Implementation in service
  },
});
```

#### Service Implementation
```typescript
// services/knowledge/src/commands/create-knowledge-item.ts
import { createKnowledgeItem } from '@know-s/contracts/knowledge/commands/create-knowledge-item';

export const createKnowledgeItemCommand = createKnowledgeItem.implement({
  async execute(input, deps) {
    const { db, validate, events } = deps;
    
    // Validate input using validation service
    const validated = await validate.knowledgeItem(input);
    
    // Execute database operation
    const result = await db.knowledgeItems.create({
      data: validated,
    });
    
    // Emit event
    await events.emit('knowledge:item:created', {
      itemId: result.id,
      type: input.type,
    });
    
    return result;
  },
});
```

#### Service Registration in API
```typescript
// apps/api/src/services/knowledge-service.ts
import { createKnowledgeItemCommand } from '@know-s/services/knowledge';

export function createKnowledgeService(dependencies: KnowledgeDependencies) {
  return {
    createKnowledgeItem: createKnowledgeItemCommand.bind(null, dependencies),
    // Other commands...
  };
}

// apps/api/src/app.ts
const knowledgeService = createKnowledgeService({
  db: persistenceClient,
  validate: validationService,
  events: eventService,
});

// Use in routes
app.post('/api/knowledge/items', async (request, reply) => {
  const result = await knowledgeService.createKnowledgeItem(request.body);
  return reply.send(result);
});
```

### 4.2 Dependency Management Pattern

**No Traditional DI Container**: Services receive dependencies directly via factory functions.

**Dependency Flow**:
1. API server initializes all shared dependencies (persistence, events, etc.)
2. API creates service instances by passing required dependencies
3. Services are stateless and receive dependencies at creation time
4. Commands are bound to dependencies for execution

**Benefits**:
- Explicit dependency declaration in contracts
- No runtime dependency resolution overhead
- Easy testing with mocked dependencies
- Clear dependency graph visualization

---

## 5. Development Workflow

### 5.1 Getting Started
```bash
# Clone repository
git clone <repository-url>
cd product-os

# Install dependencies
pnpm install

# Start development infrastructure
pnpm db:up

# Run database migrations
pnpm db:migrate

# Generate TypeScript types from database
pnpm db:generate

# Start development servers
pnpm dev

# Run tests
pnpm test

# Run type checking
pnpm typecheck

# Run linting
pnpm lint
```

### 5.2 Development Commands
| Command | Purpose |
|---------|---------|
| `pnpm dev` | Start all services in development mode |
| `pnpm build` | Build all packages for production |
| `pnpm build:transpile` | Transpile TypeScript to JavaScript |
| `pnpm typecheck` | Run TypeScript compiler checks |
| `pnpm test` | Run all tests |
| `pnpm test:watch` | Run tests in watch mode |
| `pnpm lint` | Run ESLint on all packages |
| `pnpm format` | Format code with Prettier |
| `pnpm db:up` | Start development databases |
| `pnpm db:down` | Stop development databases |
| `pnpm db:migrate` | Run database migrations |
| `pnpm db:generate` | Generate TypeScript types from database |

### 5.3 Environment Configuration

#### `.env.example`
```env
# Database
POSTGRES_URL=postgresql://postgres:postgres@localhost:5432/product_os
NEO4J_URL=bolt://localhost:7687
NEO4J_USER=neo4j
NEO4J_PASSWORD=password
REDIS_URL=redis://localhost:6379

# Storage
S3_ENDPOINT=http://localhost:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=product-os

# API
PORT=3000
NODE_ENV=development
LOG_LEVEL=info

# Security
JWT_SECRET=your-jwt-secret-key-here
```

#### Environment Validation
```typescript
// lib/shared/src/config/env.ts
import { z } from 'zod';

const envSchema = z.object({
  POSTGRES_URL: z.string().url(),
  NEO4J_URL: z.string().url(),
  REDIS_URL: z.string().url(),
  PORT: z.coerce.number().default(3000),
  NODE_ENV: z.enum(['development', 'test', 'production']).default('development'),
});

export const env = envSchema.parse(process.env);
```

---

## 6. Testing Strategy

### 6.1 Unit Testing (Service Level)
```typescript
// services/knowledge/test/commands/create-knowledge-item.test.ts
import { describe, test, expect, mock } from 'bun:test';
import { createKnowledgeItemCommand } from '../../src/commands/create-knowledge-item';

describe('createKnowledgeItem', () => {
  test('creates knowledge item with valid input', async () => {
    const mockDb = {
      knowledgeItems: {
        create: mock(async () => ({
          id: '123',
          createdAt: new Date(),
          updatedAt: new Date(),
        })),
      },
    };

    const mockValidate = {
      knowledgeItem: mock(async (input) => input),
    };

    const mockEvents = {
      emit: mock(async () => {}),
    };

    const command = createKnowledgeItemCommand({
      db: mockDb,
      validate: mockValidate,
      events: mockEvents,
    });

    const result = await command({
      title: 'Test Item',
      content: 'Test content',
      type: 'note',
    });

    expect(result.id).toBe('123');
    expect(mockDb.knowledgeItems.create).toHaveBeenCalled();
  });
});
```

### 6.2 Integration Testing (API Level)
```typescript
// apps/api/test/routes/knowledge.test.ts
import { describe, test, expect, beforeAll } from 'bun:test';
import { createApp } from '../src/app';

describe('Knowledge API', () => {
  let app: FastifyInstance;

  beforeAll(async () => {
    app = await createApp();
  });

  test('POST /api/knowledge/items creates item', async () => {
    const response = await app.inject({
      method: 'POST',
      url: '/api/knowledge/items',
      payload: {
        title: 'Test API Item',
        content: 'API test content',
        type: 'note',
      },
    });

    expect(response.statusCode).toBe(201);
    const body = JSON.parse(response.body);
    expect(body.id).toBeDefined();
  });
});
```

### 6.3 Contract Testing
```typescript
// lib/contracts/test/knowledge-contracts.test.ts
import { describe, test, expect } from 'bun:test';
import { createKnowledgeItem } from '../knowledge/commands/create-knowledge-item';

describe('Knowledge Contracts', () => {
  test('createKnowledgeItem contract validates input', () => {
    const validInput = {
      title: 'Valid Title',
      content: 'Valid content',
      type: 'note',
    };

    const invalidInput = {
      title: '', // Empty title
      content: 'Valid content',
      type: 'invalid', // Invalid type
    };

    expect(() => createKnowledgeItem.input.parse(validInput)).not.toThrow();
    expect(() => createKnowledgeItem.input.parse(invalidInput)).toThrow();
  });
});
```

---

## 7. Code Quality & Standards

### 7.1 ESLint Configuration
```javascript
// .eslintrc.js
module.exports = {
  root: true,
  extends: [
    'eslint:recommended',
    '@typescript-eslint/recommended',
    'prettier',
  ],
  parser: '@typescript-eslint/parser',
  plugins: ['@typescript-eslint'],
  env: {
    node: true,
    es2022: true,
  },
  rules: {
    '@typescript-eslint/explicit-function-return-type': 'warn',
    '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
    'no-console': ['warn', { allow: ['warn', 'error'] }],
  },
};
```

### 7.2 Prettier Configuration
```json
{
  "semi": true,
  "trailingComma": "es5",
  "singleQuote": true,
  "printWidth": 100,
  "tabWidth": 2,
  "useTabs": false
}
```

### 7.3 Git Hooks (Husky)
```json
{
  "husky": {
    "hooks": {
      "pre-commit": "pnpm lint-staged",
      "pre-push": "pnpm typecheck && pnpm test"
    }
  },
  "lint-staged": {
    "*.{ts,tsx}": ["eslint --fix", "prettier --write"],
    "*.{js,jsx}": ["eslint --fix", "prettier --write"],
    "*.{json,md,yml,yaml}": ["prettier --write"]
  }
}
```

### 7.4 Commit Convention
Follow [Conventional Commits](https://www.conventionalcommits.org/):
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting, etc.)
- `refactor:` Code refactoring
- `test:` Test changes
- `chore:` Build process or tooling changes

---

## 8. Deployment Considerations

### 8.1 Production Build
```bash
# Build all packages
pnpm build

# Run tests
pnpm test

# Type check
pnpm typecheck

# Create production bundle
cd apps/api && bun build --outdir=dist --target=node src/index.ts
```

### 8.2 Docker Production Image
```dockerfile
# Dockerfile for API
FROM oven/bun:1 as base
WORKDIR /app

# Install dependencies
COPY package.json pnpm-lock.yaml ./
RUN bun install --frozen-lockfile --production

# Copy built application
COPY dist/ ./dist/

# Run application
CMD ["bun", "run", "dist/index.js"]
```

### 8.3 Environment-Specific Configuration
- **Development**: `.env.local` (gitignored)
- **Testing**: `.env.test` (for CI/CD)
- **Production**: Environment variables only (no file)

---

## 9. Implementation Checklist

### 9.1 Repository Setup
- [ ] Initialize monorepo with pnpm workspaces
- [ ] Configure TurboRepo for task orchestration
- [ ] Set up TypeScript base configuration
- [ ] Create directory structure
- [ ] Add root package.json with shared scripts
- [ ] Configure ESLint and Prettier
- [ ] Set up Husky for git hooks
- [ ] Create docker-compose.yml for local development

### 9.2 Core Libraries
- [ ] Create lib/contracts with contract patterns
- [ ] Implement lib/persistence with database clients
- [ ] Set up lib/shared with utilities
- [ ] Configure environment validation
- [ ] Implement error handling patterns
- [ ] Set up logging configuration

### 9.3 Services
- [ ] Create knowledge service with contracts
- [ ] Implement extraction service with contracts
- [ ] Create search service with contracts
- [ ] Set up service testing patterns
- [ ] Implement command validation
- [ ] Create service documentation

### 9.4 API Application
- [ ] Set up Fastify application structure
- [ ] Implement dependency injection pattern
- [ ] Create API routes for all services
- [ ] Implement authentication middleware
- [ ] Set up request/response validation
- [ ] Create E2E test suite
- [ ] Configure production build

### 9.5 Database & Infrastructure
- [ ] Create PostgreSQL schema and migrations
- [ ] Set up Neo4j graph schema
- [ ] Configure Redis caching patterns
- [ ] Implement object storage integration
- [ ] Create database seed scripts
- [ ] Set up backup and recovery procedures

---

## 10. Future Considerations

### 10.1 Scaling Patterns
- **Horizontal Scaling**: Stateless services allow easy horizontal scaling
- **Database Sharding**: PostgreSQL sharding strategies for large datasets
- **Caching Layers**: Multi-level caching with Redis and CDN
- **Message Queues**: Event-driven architecture with message queues

### 10.2 Monitoring & Observability
- **Logging**: Structured logging with correlation IDs
- **Metrics**: Prometheus metrics for service monitoring
- **Tracing**: Distributed tracing with OpenTelemetry
- **Health Checks**: Comprehensive health check endpoints

### 10.3 Security Enhancements
- **API Security**: Rate limiting, request validation, input sanitization
- **Data Encryption**: Field-level encryption for sensitive data
- **Audit Logging**: Comprehensive audit trail for all operations
- **Compliance**: GDPR, CCPA, and other regulatory compliance features

### 10.4 Developer Experience
- **Local Development**: Improved local development tooling
- **Code Generation**: Scaffolding tools for new services and commands
- **Documentation**: Interactive API documentation with Swagger/OpenAPI
- **Debugging**: Enhanced debugging tools and error reporting

---

**Last Updated**: 2026-01-12  
**Specification Version**: 1.0  
**Maintained By**: Architecture Team