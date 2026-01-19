# TypeScript LLM Libraries Research - Chat Export

## User Question
I want to build a typescript poc app that interacts with multiple llM's depending on user selection. Many of the responses I'll want from LLMs will be structured json (tokenization etc). I'll also want to interact with a dynamic number of mcp servers that I've created. What are the latest and greatest libraries for doing all of this with typescript. Please do some research

## Claude's Response

Based on my research, here are the latest and greatest TypeScript libraries for building your multi-LLM POC app with structured JSON outputs and MCP server integration:

### **For Structured JSON Outputs**

#### **1. Instructor-JS** (Highly Recommended)
- **GitHub**: `@instructor-ai/instructor`
- Built on top of Zod for schema validation
- Works seamlessly with OpenAI SDK and supports multiple providers via `llm-polyglot`
- Offers partial JSON streaming from LLMs
- Provides automatic parsing with retry logic
- Type-safe with full TypeScript support

#### **2. Vercel AI SDK** (Most Comprehensive)
- **Package**: `ai` with provider packages (`@ai-sdk/openai`, `@ai-sdk/anthropic`, `@ai-sdk/google`)
- **Latest Version**: AI SDK 6.0+ (released late 2025)
- Unified API across 20+ LLM providers
- Built-in structured output via `generateObject()` and `streamObject()`
- Native MCP integration support
- Includes `ToolLoopAgent` class for agentic workflows
- DevTools for debugging (`@ai-sdk/devtools`)
- Over 20M monthly downloads

#### **3. LangChain.js**
- **Package**: `@langchain/core`, `@langchain/openai`, etc.
- `.withStructuredOutput()` method with Zod schemas
- `StructuredOutputParser` for JSON validation
- Good for complex chains and workflows
- Note: Can be verbose for simple use cases

### **For MCP (Model Context Protocol) Integration**

#### **Official MCP TypeScript SDK**
- **Package**: `@modelcontextprotocol/sdk`
- **Latest Version**: 1.24.1 (December 2025, spec 2025-11-25)
- Official SDK from Anthropic
- Full MCP specification support
- Supports both client and server implementations
- Multiple transport options:
  - `StdioServerTransport` (for local servers)
  - `StreamableHTTPClientTransport` (for production/remote)
  - SSE for backwards compatibility
- Built-in Zod integration for schema validation
- Examples available in the repository

#### **Vercel AI SDK MCP Integration**
- Built-in MCP client creation via `createMCPClient()`
- Automatic tool conversion from MCP servers
- HTTP and SSE transport support
- Works seamlessly with AI SDK's tool system

### **Multi-LLM Provider Management**

#### **Best Options:**

1. **Vercel AI SDK** - Most mature
   - Simple provider switching: `model: 'anthropic/claude-opus-4.5'` or `model: 'openai/gpt-4o'`
   - AI Gateway for unified access
   - Consistent API regardless of provider

2. **LiteLLM** (via Node wrapper if needed)
   - Unified interface for 100+ LLM providers
   - Good for complex provider routing

### **Recommended Stack for Your POC**

```typescript
// Core dependencies
{
  "dependencies": {
    // For multi-LLM orchestration & structured outputs
    "ai": "^6.0.0",
    "@ai-sdk/openai": "latest",
    "@ai-sdk/anthropic": "latest",
    "@ai-sdk/google": "latest",
    
    // For enhanced structured output handling
    "@instructor-ai/instructor": "latest",
    
    // For MCP server integration
    "@modelcontextprotocol/sdk": "^1.24.1",
    
    // Schema validation (peer dependency)
    "zod": "^3.25.0",
    
    // Optional: DevTools for debugging
    "@ai-sdk/devtools": "latest"
  }
}
```

### **Architecture Suggestion**

```typescript
// Example structure combining everything:

// 1. Initialize multiple LLM providers
import { createOpenAI } from '@ai-sdk/openai';
import { createAnthropic } from '@ai-sdk/anthropic';
import { generateObject } from 'ai';
import { z } from 'zod';

// 2. MCP Server client
import { Client } from '@modelcontextprotocol/sdk/client/index.js';
import { StreamableHTTPClientTransport } from '@modelcontextprotocol/sdk/client/transport.js';

// 3. Structured output schema
const OutputSchema = z.object({
  result: z.string(),
  confidence: z.number(),
  metadata: z.record(z.unknown())
});

// 4. Dynamic provider selection
async function callLLM(provider: string, prompt: string) {
  const model = provider === 'openai' 
    ? 'openai/gpt-4o' 
    : 'anthropic/claude-opus-4.5';
    
  const result = await generateObject({
    model,
    schema: OutputSchema,
    prompt
  });
  
  return result.object;
}

// 5. Connect to MCP servers dynamically
async function connectToMCPServer(serverUrl: string) {
  const transport = new StreamableHTTPClientTransport(serverUrl);
  const client = new Client({
    name: "my-client",
    version: "1.0.0"
  }, {
    capabilities: {}
  });
  
  await client.connect(transport);
  return client;
}
```

### **Key Advantages of This Stack**

1. **Type Safety**: Full TypeScript support with Zod schemas
2. **Provider Flexibility**: Easy switching between OpenAI, Anthropic, Google, etc.
3. **Production Ready**: Battle-tested by major companies
4. **MCP Support**: Native integration with MCP servers
5. **Structured Outputs**: Multiple approaches (native API methods + Instructor for complex cases)
6. **Active Development**: All libraries actively maintained with recent updates

The Vercel AI SDK + MCP SDK combination gives you the most modern, well-documented, and production-ready foundation for your POC.

---

## Follow-up Question
how can i download this entire chat as a markdown file?

---

*Chat exported on January 11, 2026*
