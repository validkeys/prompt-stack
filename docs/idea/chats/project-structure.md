This is going to be a typescript monorepo.

## Rules

- Uses vite for testing
- Uses turbo repo for repo management
- Uses the same build, build:transpile and typecheck methods as ~/Sites/lwtree/develop
- Uses pnpm for package management
- Uses kysely for database queries and the kysely cli for database migrations. We will use kysely-autogen to auto generate the types from the database
- Setup a docker-compose file for databases, and whatever other 3rd party infra is needed
- all services should follow the service command architecture: /Users/kyledavis/Sites/lwtree/develop/specifications/system/core/service-command-architecture.md
- All services should be defined using the contracted framework (https://github.com/validkeys/contracted) defining their contracts in the contracts folder below.
- Services can not directly import one another, they must use DI
- The "app" in the apps folder is where all services are given their dependencies
- We will use zod3 for all schemas
- Will will fastify for our API server
- Will will use a namespace of @know-s for all packages
- We will use BUN for the runtime here

## Structure

- apps - deployable apps (api)
- services - domains mounted as services (knowledge etc..)
- lib - internal use only (postgres etc)
  - has a folder called "contracts" all service contracts, type definitions from the db and other shared interface go here (see ~/Sites/lwtree/develop/g-foundation/contracts as an example)

## Service Structure

- src
- handlers/
- lib/

Use /Users/kyledavis/Sites/lwtree/develop/d-modules/gapo as a good example of how to define a service and it's handlers.

We should also be using the technologies outlined in chats/typescript-llm-chat-export.md
