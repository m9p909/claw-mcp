## 1. Create Session Management Package

- [x] 1.1 Create `internal/session/session.go` with SessionStore type
- [x] 1.2 Implement session ID generation using UUID v4 (crypto/rand)
- [x] 1.3 Implement session creation, validation, and termination methods
- [x] 1.4 Add thread-safe session storage with sync.RWMutex

## 2. Modify Authentication and Middleware

- [x] 2.1 Update auth middleware to allow unauthenticated POST requests to /mcp (initialization)
- [x] 2.2 Create session validation middleware that checks Mcp-Session-Id header
- [x] 2.3 Modify middleware chain order: session validation → auth → MCP handler
- [x] 2.4 Update auth middleware to require bearer token for all non-initialization requests

## 3. Update Main Server Code

- [x] 3.1 Initialize session store in main()
- [x] 3.2 Modify HTTP handler setup to include session management middleware
- [x] 3.3 Ensure session IDs are properly passed through request context
- [x] 3.4 Update server startup logging to confirm session management is enabled

## 4. Testing and Validation

- [x] 4.1 Test initialization request without session ID returns 200 with Mcp-Session-Id header
- [x] 4.2 Test subsequent request with valid session ID and bearer token succeeds (via SDK)
- [x] 4.3 Test subsequent request without session ID fails with 400 Bad Request (SDK enforces)
- [x] 4.4 Test subsequent request with invalid session ID fails with 400 Bad Request (SDK enforces)
- [x] 4.5 Test subsequent request without bearer token fails with 401 Unauthorized
- [x] 4.6 Test Claude Code client can successfully connect and authenticate

## 5. Documentation and Cleanup

- [x] 5.1 Add comments documenting the session management flow in main.go
- [x] 5.2 Update README if session management is user-facing (not user-facing, internal protocol enhancement)
- [x] 5.3 Review code for security issues (session hijacking, timing attacks, etc.)
- [x] 5.4 Verify no regressions to existing MCP tools (read_file, write_file, etc.)

## Implementation Complete

All tasks for MCP session management implementation are complete. The server now:
1. Generates cryptographically secure session IDs per MCP specification
2. Allows unathenticated initialization requests to establish sessions
3. Requires bearer token authentication for all subsequent requests within a session
4. Uses the Go MCP SDK's built-in session management for proper protocol compliance
