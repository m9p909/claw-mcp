## MODIFIED Requirements

### Requirement: Bearer token authentication
The server SHALL require a valid Bearer token in the Authorization header for all requests to the MCP endpoint, except for initialization requests that are establishing a new session.

#### Scenario: Initialization request without authorization
- **WHEN** client sends an InitializeRequest without an Authorization header and without an Mcp-Session-Id
- **THEN** server allows the request to proceed to generate and return a session ID

#### Scenario: Subsequent request with valid bearer token
- **WHEN** client sends a request with both a valid Mcp-Session-Id (from initialization) and valid Bearer token
- **THEN** server authenticates the request and processes it normally

#### Scenario: Subsequent request without bearer token
- **WHEN** client sends a request with a valid Mcp-Session-Id but no Authorization header
- **THEN** server responds with HTTP 401 Unauthorized

#### Scenario: Invalid bearer token
- **WHEN** client sends a request with an invalid Bearer token
- **THEN** server responds with HTTP 401 Unauthorized and rejects the request
