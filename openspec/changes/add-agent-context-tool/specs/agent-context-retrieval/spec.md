## ADDED Requirements

### Requirement: Provide essential context in tool description
The system SHALL include essential agent guidance in the `get_agent_context` tool description field, visible in MCP tool listings.

#### Scenario: Agent lists available tools on connect
- **WHEN** an agent connects and lists MCP tools
- **THEN** the `get_agent_context` tool description contains imperative guidance on Claw identity (personal agent in Linux environment), professional behavior, token efficiency, and skills system basics

#### Scenario: Tool description length is concise
- **WHEN** the tool description is rendered
- **THEN** it MUST be under 500 characters and prioritize essential information

### Requirement: Return full agent context documentation
The system SHALL return complete agent-oriented documentation when `get_agent_context` is invoked.

#### Scenario: Agent calls get_agent_context
- **WHEN** an agent invokes the `get_agent_context` tool with no parameters
- **THEN** the system returns the full content of `docs/AGENT_CONTEXT.md` as text

#### Scenario: Documentation file does not exist
- **WHEN** an agent invokes `get_agent_context` and `docs/AGENT_CONTEXT.md` is missing
- **THEN** the system returns an error indicating the documentation file could not be found

### Requirement: Documentation content follows imperative tone
The documentation SHALL use imperative voice to provide direct instructions to agents.

#### Scenario: Content structure covers key areas
- **WHEN** the documentation is read
- **THEN** it MUST cover: Claw identity and environment, recommended personality (professional, concise), token efficiency guidelines, and Agent Skills system explanation with usage

#### Scenario: Skills system explanation includes discovery
- **WHEN** the skills section is read
- **THEN** it MUST explain that skills are at `~/.mcpclaw/skills/`, describe the SKILL.md format, and instruct agents to use `list_skills` and `get_skill` tools for discovery

### Requirement: Tool registration follows existing patterns
The system SHALL register `get_agent_context` tool following the same pattern as existing tools like `list_skills` and `get_skill`.

#### Scenario: Tool registered in server initialization
- **WHEN** the MCP server initializes in `internal/server.go`
- **THEN** the `get_agent_context` tool is registered with appropriate name, description, and handler function

#### Scenario: Handler implementation in pkg/tools
- **WHEN** the handler is implemented
- **THEN** it MUST be located in `pkg/tools/` directory, follow naming conventions (e.g., `HandleGetAgentContext`), and return `*mcp.CallToolResult`
