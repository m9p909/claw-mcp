## MODIFIED Requirements

### Requirement: MCP protocol transport uses streamable-http
The server SHALL implement the MCP protocol transport using streamable-http as defined in the MCP standard specification, not Server-Sent Events (SSE).

#### Scenario: Client connects via streamable-http
- **WHEN** a client initiates an MCP connection to the `/mcp` endpoint
- **THEN** the connection uses HTTP with streaming response bodies per MCP spec, not SSE

#### Scenario: Server responds with proper MCP messages
- **WHEN** the client sends an initialize request
- **THEN** the server responds with MCP-formatted JSON messages using the streamable-http protocol

### Requirement: Server implements full MCP protocol methods
The server SHALL support all standard MCP protocol methods: initialize, tools/list, tools/call, resources/list, resources/read, prompts/list, and prompts/get.

#### Scenario: Client lists available tools
- **WHEN** client sends `tools/list` method request
- **THEN** server responds with all 8 registered tools (read_file, write_file, edit_file, exec_command, manage_process, write_memory, query_memory, memory_search)

### Requirement: Server manages MCP sessions
The server SHALL maintain session state for connected clients, tracking protocol version, capabilities, and client identity throughout the connection lifecycle.

#### Scenario: Session persists across multiple requests
- **WHEN** a client connects and sends multiple tool calls
- **THEN** the session remains active and stateful until the client disconnects or times out
