## Context

Currently, tool failures only appear in MCP responses returned to clients—not in server logs. When a write operation fails, the error message reaches the client but the server has no record of what went wrong. Without request/session correlation, it's impossible to trace a failure back through the system. The HTTP layer logs only method/path; tool handlers return errors but don't log them.

The system needs structured logging that captures:
- What operation was attempted (tool name, inputs)
- What went wrong and why (error code, message, context)
- How to correlate it across requests (request ID, session ID)
- When and how long it took

## Goals / Non-Goals

**Goals:**
- Add request ID (UUID) generation and propagation through the stack
- Extract and propagate session ID from HTTP headers (Mcp-Session-Id)
- Log all tool handler entry/exit with inputs/outputs and execution time
- Log all errors server-side before returning to client
- Add operation-specific logging for filesystem, execution, and memory tools
- Support conditional logging via DEBUG_LEVEL env var (INFO, DEBUG, TRACE)
- Ensure no raw sensitive data is logged (sanitize paths in logs, only error messages)
- Standard library logging only (log/slog for structured logging)
- Stream to stdout (no file rotation)

**Non-Goals:**
- Log aggregation or external systems
- Custom log formatting (use slog defaults)
- Per-tool log levels (single DEBUG_LEVEL applies globally)
- Async logging or buffering

## Decisions

### 1. Use log/slog for structured logging
**Decision**: Use Go 1.21+ standard library `log/slog` for structured, leveled logging.

**Rationale**:
- Standard library (no external dependencies)
- Built-in support for log levels (DEBUG, INFO, WARN, ERROR)
- Structured key-value fields enable easy correlation by request/session ID
- JSON output available for parsing

**Alternative considered**: logrus/zap would add external dependencies; standard log package lacks structure.

### 2. Request IDs via UUID, propagated in context
**Decision**: Generate UUID v4 for each MCP request in middleware. Pass via `context.Context` to tool handlers.

**Rationale**:
- Unique per-request identification for tracing
- Context is Go idiomatic for request-scoped data
- Easy to add to every log statement

**Propagation flow**:
```
HTTP Middleware (generate UUID)
  → context.WithValue(ctx, "requestID", uuid)
  → tool handler receives ctx
  → all log calls include requestID
```

### 3. Session ID from HTTP headers
**Decision**: Extract `Mcp-Session-Id` header in middleware, pass alongside request ID via context.

**Rationale**:
- Session ID already in HTTP headers (MCP spec)
- Enables correlation across multiple requests in same session
- Already validated by MCP SDK

### 4. DEBUG_LEVEL env var for conditional logging
**Decision**: Support three levels: INFO (default), DEBUG, TRACE. Single global setting via `DEBUG_LEVEL` env var.

**Rationale**:
- Reduces log noise in production (default INFO only logs warnings, errors, key events)
- DEBUG adds entry/exit logging, operation details
- TRACE adds all low-level details (hash validation, permission checks, etc.)

### 5. Log all errors server-side before returning
**Decision**: In errorResult() and all handlers, log errors server-side (with full context) before returning error response to client.

**Rationale**:
- Ensures errors are always visible in logs, not just in client responses
- Provides full context (request ID, operation, inputs) that client doesn't have

### 6. No raw sensitive data in logs
**Decision**:
- Log file paths as `path="<sanitized>"` (hash or mask, not the actual path)
- Log error messages in full (they are user-facing and contain intent)
- Log operation names, command names (not arguments)

**Rationale**:
- Prevents accidental leakage of file paths that might contain secrets
- Error messages are already shown to users, safe to log
- Balances debuggability with security

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| Log bloat at TRACE level | Document DEBUG_LEVEL; only use TRACE for active debugging |
| UUID generation overhead per request | Minimal (microseconds); trade-off worth it for correlation |
| Context passing through all handlers | Go idiom; minimal performance impact |
| Sanitizing paths reduces debuggability | Still log error messages and operation context; paths only masked |

## Migration Plan

1. Create pkg/log/logger.go with slog initialization and helper functions
2. Update main.go: initialize logger, generate request IDs, pass via context
3. Update pkg/tools/helpers.go: enhance errorResult() to log server-side
4. Update all tool handlers: add entry/exit logging, operation logging
5. No breaking changes; feature-additive only
6. Rollback: remove logger calls (handlers still function normally)

## Open Questions

- Should request ID also be included in HTTP response headers for client correlation? (Can implement later if needed)
