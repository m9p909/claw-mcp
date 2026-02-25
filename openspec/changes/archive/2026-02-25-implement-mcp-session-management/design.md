## Context

The claw MCP server is built on the Go MCP SDK's `NewStreamableHTTPHandler`, which implements the MCP Streamable HTTP Transport specification. However, the current implementation doesn't properly handle session ID management as required by the spec. The Go SDK's handler expects session IDs to be managed by middleware or a wrapper layer, but no such layer currently exists.

Claude Code's HTTP MCP client follows the specification: it expects the server to generate and return a session ID during initialization, then includes that ID in all subsequent requests. Without this, the client cannot maintain protocol state and fails with "Failed to reconnect to claw" errors.

## Goals / Non-Goals

**Goals:**
- Implement proper session ID generation and tracking per MCP Streamable HTTP specification
- Allow Claude Code's HTTP MCP client to successfully connect and communicate with the server
- Maintain security by requiring bearer token authentication for all requests (even after session establishment)
- Implement session lifecycle management (initialization, validation, termination)

**Non-Goals:**
- Change the MCP tools or capabilities (read_file, write_file, etc.)
- Modify the database or storage layer
- Implement session persistence across server restarts
- Add new authentication methods beyond bearer tokens

## Decisions

### Decision 1: Session ID Generation Strategy
**Choice**: Use UUID v4 (via `crypto/rand`) for session IDs, stored in memory with a map protected by RWMutex.

**Rationale**: UUIDs are simple, cryptographically secure, and provide the required global uniqueness. In-memory storage is sufficient for single-instance deployment; for distributed setups, this can be extended to use a session store.

**Alternatives Considered**:
- JWT tokens: More complex, but no additional benefit for this use case
- Database-backed sessions: Not needed for current deployment model
- Timestamp-based: Insufficient entropy, vulnerable to collision attacks

### Decision 2: Middleware Architecture
**Choice**: Create a session management middleware that wraps the authentication middleware. Order: Session validation → Bearer token auth → MCP handler.

**Rationale**: Session management must happen before auth so we can extract the session ID, and auth must happen before the MCP handler to ensure security. This layering keeps concerns separated.

**Implementation**:
- Session middleware checks if request contains valid Mcp-Session-Id header
- If no session ID and is POST to /mcp (initialization), allow to proceed to auth
- If has session ID, validate it exists in session store before proceeding to auth
- Auth middleware checks Bearer token as before

### Decision 3: Initialization vs. Subsequent Requests
**Choice**: Distinguish initialization (POST with no session ID) from regular requests by checking request headers in middleware.

**Rationale**: MCP spec allows unauthenticated initialization so clients can establish sessions. After initialization, all requests require both session ID (from init response) and bearer token.

**Alternatives Considered**:
- Require bearer token even for initialization: Would violate MCP spec which expects clients to auth after session establishment
- Allow session-less requests: Would violate spec requirements

### Decision 4: Session Storage
**Choice**: In-memory map with goroutine-safe access (sync.RWMutex) tracking: session ID → session metadata (creation time, last activity).

**Rationale**: Simple, fast, and sufficient for current needs. For scale, can be replaced with Redis or similar.

**Metadata tracked**:
- Session ID
- Created timestamp
- Last activity timestamp
- Optional: request counter for debugging

## Risks / Trade-offs

**[Risk] In-memory session loss on restart** → Session store doesn't persist across server restarts. Clients must re-initialize. This is acceptable for development/testing.

**[Risk] Session ID collision** → Extremely unlikely with UUIDs (1 in 2^122 probability), but mitigated by checking session store before accepting requests.

**[Risk] Session table memory growth** → Long-running server with many clients could accumulate session entries. **Mitigation**: Implement session cleanup (e.g., expire sessions after 24 hours of inactivity).

**[Risk] Middleware ordering complexity** → Getting middleware order wrong could bypass authentication. **Mitigation**: Document middleware chain clearly and add tests.

**[Trade-off] Simplicity vs. Feature-richness** → We're not implementing session timeout, token refresh, or other advanced features. Acceptable because current use case is human-interactive clients.

## Migration Plan

1. Add session management package with session store and ID generation
2. Modify auth middleware to allow initialization requests without bearer token
3. Add session validation middleware
4. Deploy updated binary with CLAW_TOKEN already configured (no env changes needed)
5. Claude Code automatically re-connects and establishes session

**Rollback**: If issues occur, revert to previous binary - clients will fail to connect but no data loss.

## Open Questions

- Should we implement session expiration? (Recommendation: Add in follow-up if needed)
- Should sessions persist in database for debugging/audit? (Recommendation: Later enhancement)
