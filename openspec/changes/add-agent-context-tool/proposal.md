## Why

Agents connecting to Claw need immediate context about what Claw is, how to behave efficiently, and how to use the skills system. Without this, agents waste tokens on verbose responses and may not discover or utilize critical features like Agent Skills.

## What Changes

- Add `get_agent_context` MCP tool that returns agent-oriented documentation
- Create `docs/AGENT_CONTEXT.md` with imperative guidance on Claw identity, personality, token efficiency, and skills system
- Tool description contains essential context (visible in tool listings), full content available on fetch
- Documentation positions Claw as a personal agent executing in a real Linux environment
- Guides agents to be professional, concise, and token-efficient

## Capabilities

### New Capabilities

- `agent-context-retrieval`: Provide agents with essential context about Claw's identity, recommended behavior, and available features through a dedicated MCP tool

### Modified Capabilities

<!-- No existing capabilities modified -->

## Impact

- New tool registered in `internal/server.go`
- New handler in `pkg/tools/` to serve documentation
- New markdown file at `docs/AGENT_CONTEXT.md`
- Agents gain immediate access to guidance on first connect (via tool description) and detailed reference (via tool invocation)
