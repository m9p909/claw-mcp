## ADDED Requirements

### Requirement: Bearer token validation on MCP endpoint

The MCP server SHALL require a valid Bearer token in the `Authorization` header for all requests to the `/mcp` endpoint. The token SHALL be validated against the `CLAW_TOKEN` environment variable.

#### Scenario: Valid token provided

- **WHEN** a request to `/mcp` includes `Authorization: Bearer <valid-token>`
- **THEN** the request is processed normally by the MCP handler

#### Scenario: Missing authorization header

- **WHEN** a request to `/mcp` is made without an `Authorization` header
- **THEN** the server returns 401 Unauthorized with body `{"error": "Unauthorized"}`

#### Scenario: Invalid token

- **WHEN** a request to `/mcp` includes `Authorization: Bearer <invalid-token>`
- **THEN** the server returns 401 Unauthorized with body `{"error": "Unauthorized"}`

#### Scenario: Malformed authorization header

- **WHEN** a request to `/mcp` includes an `Authorization` header without "Bearer " prefix
- **THEN** the server returns 401 Unauthorized with body `{"error": "Unauthorized"}`

### Requirement: Environment variable configuration

The server SHALL read the `CLAW_TOKEN` environment variable at startup and fail immediately if it is not set.

#### Scenario: CLAW_TOKEN is set

- **WHEN** the server starts with `CLAW_TOKEN` environment variable defined
- **THEN** the server initializes successfully and listens for connections

#### Scenario: CLAW_TOKEN is not set

- **WHEN** the server starts without `CLAW_TOKEN` environment variable
- **THEN** the server logs an error and exits with a non-zero status code

### Requirement: Health endpoint remains unauthenticated

The `/health` endpoint SHALL remain accessible without authentication.

#### Scenario: Health check without token

- **WHEN** a request to `/health` is made without an `Authorization` header
- **THEN** the server returns 200 OK with body `{"status":"ok"}`

#### Scenario: Health check with token

- **WHEN** a request to `/health` is made with a valid `Authorization: Bearer` header
- **THEN** the server returns 200 OK with body `{"status":"ok"}`
