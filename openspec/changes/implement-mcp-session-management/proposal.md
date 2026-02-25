## Why

The claw MCP server currently uses the Go SDK's `NewStreamableHTTPHandler` without properly implementing the MCP Streamable HTTP Transport specification for session ID management. This causes Claude Code's MCP client to fail during connection with "Failed to reconnect to claw" errors. The spec requires the server to assign and return a unique session ID during initialization, which clients must then include in all subsequent requests. Without proper session ID handling, the MCP protocol cannot maintain state across multiple requests.

## What Changes

- Implement explicit session ID generation and tracking in the HTTP request handler
- Add session ID header handling in the authentication middleware to allow session ID headers to bypass auth for initialization requests
- Ensure session IDs are cryptographically secure (UUIDs) and properly returned in initialization responses
- Implement session validation logic to reject requests without valid session IDs (per spec)
- Add session cleanup mechanisms for terminated sessions

## Capabilities

### New Capabilities
- `mcp-session-management`: Manages MCP protocol session lifecycle (initialization, validation, termination) according to the Streamable HTTP Transport specification

### Modified Capabilities
- `bearer-token-auth`: Modify to allow unauthenticated initialization requests so clients can establish sessions before providing bearer tokens

## Impact

- **Code changes**: `main.go` (auth middleware), HTTP handler setup
- **Dependencies**: `crypto/rand` and `crypto/sha256` for session ID generation
- **Behavior**: Claude Code clients will now successfully connect to the claw MCP server
- **Breaking**: None - existing API behavior unchanged, only HTTP transport layer improved
