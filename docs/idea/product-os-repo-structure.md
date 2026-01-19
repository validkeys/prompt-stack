# Product-OS: Repository Structure Requirements

**Date**: 2026-01-12  
**Status**: Requirements Document  
**Related**: [`product-os.md`](product-os.md) (Vision), [`product-os-data-structure.md`](product-os-data-structure.md) (Data Structure), [`chats/project-structure.md`](chats/project-structure.md) (Initial Requirements)

---

## Executive Summary

**Purpose**: Define the repository structure, build system, and development tooling for Product-OS implementation. This document captures requirements for a TypeScript monorepo using modern development practices and established patterns from existing projects.

**Key Requirements**:
- TypeScript monorepo with TurboRepo for management
- Service-command architecture with dependency injection
- Contract-first service definitions using Contracted framework
- Polyglot database architecture (PostgreSQL + Neo4j + Redis)
- Modern development tooling (Vite, pnpm, Bun runtime)

---

## Current Requirements (from project-structure.md)

### Technology Stack
- **Language**: TypeScript
- **Package Manager**: pnpm
- **Runtime**: Bun
- **Build System**: Vite for testing, TurboRepo for monorepo management
- **Database**: Kysely for PostgreSQL queries, Kysely CLI for migrations, kysely-autogen for type generation
- **API Server**: Fastify
- **Validation**: zod3
- **Framework**: Contracted for service contracts
- **Namespace**: @know-s for all packages
- **Infrastructure**: Docker-compose for databases and 3rd party services

### Repository Structure
```
- apps/           # Deployable applications (API server, etc.)
- services/       # Domain services (knowledge, extraction, etc.)
- lib/            # Internal libraries (postgres, contracts, etc.)
  - contracts/    # Service contracts, type definitions, shared interfaces
```

### Service Structure
```
- src/
- handlers/
- lib/
```

### Architectural Patterns
- **Service-command architecture**: Follows `/Users/kyledavis/Sites/lwtree/develop/specifications/system/core/service-command-architecture.md`
- **Dependency Injection**: Services cannot directly import each other, must use DI
- **Contract-first**: All services defined using Contracted framework with contracts in lib/contracts
- **Dependency Resolution**: The "app" in apps folder is where all services receive their dependencies
- **Reference Implementation**: Use `/Users/kyledavis/Sites/lwtree/develop/d-modules/gapo` as example

### Build & Development
- Uses same build, build:transpile and typecheck methods as `~/Sites/lwtree/develop`
- Technologies outlined in `chats/typescript-llm-chat-export.md`

---

## Interview Notes

### Q1: Apps & Services Scope
**Question**: What specific apps and services do you envision for Product-OS?
**Answer**: For the POC there would just be an API, future would have UI (React app) and CLI.
**Implications**: 
- Initial focus on API server as primary deployable app
- Future expansion to include React-based web UI and CLI tools
 - Repository structure should accommodate these future additions

### Q2: Database Access Patterns
**Question**: How should services interact with polyglot databases (PostgreSQL + Neo4j + Redis)?
**Answer**: Prefer modular monolith approach with lib/persistence vs folders for individual types of persistence.
**Implications**:
- Centralized persistence layer in lib/persistence/ for shared database access
- Avoid microservice-style per-service database clients
- Shared connection pooling and client management
 - Unified migration and schema management approach

### Q3: Dependency Injection & Contracts
**Question**: How should DI container and Contracted framework be organized?
**Answer**: Contracted library doesn't use DI container; dependencies are provided to each service in the app.
**Implications**:
- No traditional DI container; services receive dependencies directly via constructor/factory
- Contracted framework defines service contracts
- App (api) responsible for instantiating services with their dependencies
- Services remain loosely coupled through contracts

### Q4: Monorepo Build Setup
**Question**: How should package.json files and workspaces be structured?
**Answer**: Per-package configuration.
**Implications**:
- Each app, service, and library has its own package.json
- Root package.json defines workspaces and shared dependencies
- TurboRepo manages workspace dependencies and builds
- Shared build scripts can be defined at root with per-package overrides

### Q5: Environment & Deployment Configuration
**Question**: How should environment configuration and deployment be managed?
**Answer**: Environment variables + .env files.
**Implications**:
- Use environment variables for runtime configuration
- .env files for local development with environment-specific variants
- No centralized config service for initial implementation
- Simple deployment configuration approach

### Q6: Testing Strategy
**Question**: What about testing strategy and CI/CD?
**Answer**: All services should have their own unit tests using bun's test framework. We can then do e2e testing in the apps.
**Implications**:
- Unit testing at service level using Bun's built-in test framework
- E2E testing at application level (api, ui, cli)
- Contract-based testing for service interfaces
- No separate CI/CD configuration question answered yet

### Q7: Code Quality & Documentation Tooling
**Question**: What about code quality tools and documentation?
**Answer**: Question about standards/nest tools to use in a Bun environment.
**Implications**:
- Need to research Bun-compatible tooling standards
- Should include linting, formatting, commit hooks
- Documentation structure needs definition
- Tooling choices should align with Bun ecosystem

---

## Detailed Requirements

### Repository Structure
Based on the initial requirements and Q1:

#### Top-Level Structure
```
- apps/           # Deployable applications
  - api/          # API server (Fastify)
  - ui/           # React web application (future)
  - cli/          # CLI tools (future)
- services/       # Domain services
  - knowledge/    # Knowledge management service
  - extraction/   # Knowledge extraction pipeline  
  - search/       # Multi-level search service
  - [other domains]
 - lib/            # Internal libraries (modular monolith approach)
   - contracts/    # Service contracts and shared interfaces
   - persistence/  # Shared persistence layer (PostgreSQL, Neo4j, Redis)
   - shared/       # Shared utilities and types
```

#### Apps Definition
- **api**: Primary API server built with Fastify, responsible for dependency injection and service orchestration
- **ui**: Future React application for web-based interfaces
- **cli**: Future CLI tools for developer workflows

#### Services Definition
- Domain services follow service-command architecture using Contracted framework
- Each service is independently testable and deployable
- Services communicate through contracts defined in lib/contracts
- **Contracted Framework Pattern**:
  - No traditional DI container; services receive dependencies directly
  - Service contracts define command interfaces with typed dependencies
  - App (api) instantiates services by providing their required dependencies
  - Services remain loosely coupled through contract interfaces
   - Follows contract-first development with clear separation of interfaces and implementations

### Build & Development Setup
Based on the initial requirements and Q4:

#### Package Structure
- **Root-level**: `package.json` with workspaces definition, shared dev dependencies, TurboRepo configuration
- **Per-package**: Each app, service, and library has its own `package.json` with specific dependencies
- **Workspace Dependencies**: Managed via pnpm workspaces with `@know-s/*` namespace

#### Build System
- **TurboRepo**: For monorepo task orchestration and caching
- **Vite**: For testing and development builds
- **TypeScript**: Shared `tsconfig.json` base with package-specific overrides
- **Build Scripts**: Follow patterns from `~/Sites/lwtree/develop` (build, build:transpile, typecheck)

#### Development Tooling
- **Runtime**: Bun for execution and package management
- **Database**: Kysely CLI for migrations, kysely-autogen for type generation
- **Validation**: zod3 for all schemas
- **Infrastructure**: Docker-compose for databases and 3rd party services

#### Testing
Based on Q6:
- **Unit Testing**: Bun's built-in test framework for service-level unit tests
- **E2E Testing**: Application-level testing for api, ui, and cli apps
- **Contract Testing**: Leverage Contracted framework for interface testing
- **Coverage**: Integrated test coverage reporting
- **Mocking**: Contract-based mocking using Contracted framework patterns

#### Environment & Deployment
Based on Q5:
- **Configuration**: Environment variables + .env files for local development
- **Secrets**: Managed through environment variables (not committed)
- **Deployment**: Simple approach with app-specific deployment configurations
- **Docker**: Docker-compose for local development infrastructure
- **Environment Files**: `.env`, `.env.local`, `.env.production` patterns
- **Validation**: Environment variable validation using zod3 schemas

#### Code Quality & Documentation
Based on Q7:
- **Linting**: ESLint with TypeScript support (Bun-compatible configuration)
- **Formatting**: Prettier with standard configuration
- **Commit Hooks**: Husky for pre-commit and pre-push hooks
- **Type Checking**: TypeScript compiler checks integrated into build process
- **Documentation**:
  - Root README.md with project overview and getting started
  - Package-specific README.md files for apps, services, and libraries
  - API documentation generated from TypeScript types and JSDoc
  - Architecture documentation in `/docs` folder
 - **Research Needed**: Specific Bun-compatible tooling standards to be determined

---

## Implementation Plan

### Phase 1: Foundation Setup
1. Initialize monorepo with pnpm workspaces and TurboRepo
2. Configure root package.json with shared dependencies and scripts
3. Set up TypeScript configuration with base tsconfig.json
4. Establish Bun runtime configuration and tooling

### Phase 2: Build & Development Tooling
1. Configure Vite for testing and development builds
2. Set up Bun test framework configuration
3. Implement code quality tools (ESLint, Prettier, Husky)
4. Configure Docker-compose for local development infrastructure

### Phase 3: Core Architecture
1. Create lib/contracts structure for service contracts
2. Set up lib/persistence with PostgreSQL, Neo4j, Redis clients
3. Implement shared utilities and types in lib/shared
4. Establish Contracted framework patterns and conventions

### Phase 4: Service Implementation
1. Create initial services (knowledge, extraction, search)
2. Implement service contracts and command definitions
3. Set up service testing with Bun test framework
4. Configure service dependency injection patterns

### Phase 5: Application Development
1. Build API server (apps/api) with Fastify
2. Configure dependency injection and service orchestration
3. Set up environment configuration (.env files, validation)
4. Implement E2E testing for API endpoints

### Phase 6: Deployment & Documentation
1. Configure deployment scripts and environment setups
2. Create comprehensive documentation structure
3. Set up CI/CD pipeline (to be determined)
4. Establish monitoring and observability patterns

### Research Items
- Bun-compatible tooling standards for linting, formatting, etc.
- CI/CD pipeline options for Bun/TurboRepo monorepo
- Deployment strategies for modular monolith architecture
- Monitoring and observability tooling integration

---

**Last Updated**: 2026-01-12  
**Maintained By**: Architecture Team