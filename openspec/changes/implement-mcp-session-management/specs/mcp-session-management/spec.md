## ADDED Requirements

### Requirement: Server generates and assigns session IDs
The server SHALL generate a cryptographically secure, globally unique session ID when receiving an InitializeRequest without a session ID, and return this ID in the Mcp-Session-Id header of the InitializeResult response. The session ID MUST contain only visible ASCII characters (0x21 to 0x7E).

#### Scenario: First initialization request without session ID
- **WHEN** client sends an InitializeRequest without an Mcp-Session-Id header
- **THEN** server generates a new session ID and returns it in the Mcp-Session-Id response header

#### Scenario: Session ID is cryptographically secure
- **WHEN** server generates a session ID
- **THEN** the ID is a securely-generated UUID or equivalent (minimum 128-bit entropy)

### Requirement: Client includes session ID in subsequent requests
The client SHALL include the Mcp-Session-Id header returned during initialization in all subsequent HTTP requests to the MCP endpoint.

#### Scenario: Request with valid session ID
- **WHEN** client sends an HTTP request with a valid Mcp-Session-Id header from a previous initialization
- **THEN** server processes the request normally

### Requirement: Server validates session IDs on non-initialization requests
The server SHALL validate that non-initialization requests include a valid Mcp-Session-Id header. If a request lacks this header or contains an invalid session ID, the server SHALL respond with HTTP 400 Bad Request.

#### Scenario: Request without session ID header
- **WHEN** client sends a non-initialization request without an Mcp-Session-Id header
- **THEN** server responds with HTTP 400 Bad Request

#### Scenario: Request with terminated session ID
- **WHEN** client sends a request with an Mcp-Session-Id for a terminated session
- **THEN** server responds with HTTP 404 Not Found and client initiates a new session

### Requirement: Server maintains session state
The server SHALL track active sessions and their state, allowing recovery of session context across multiple requests.

#### Scenario: Session state persistence
- **WHEN** client makes multiple requests in the same session
- **THEN** server maintains the session context and correlates requests to the same session

### Requirement: Client can explicitly terminate sessions
The client MAY send an HTTP DELETE request to the MCP endpoint with an Mcp-Session-Id header to explicitly terminate a session.

#### Scenario: Explicit session termination
- **WHEN** client sends DELETE request with valid Mcp-Session-Id
- **THEN** server terminates the session (either HTTP 200 success or HTTP 405 if not supported)
