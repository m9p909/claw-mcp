## Why

Write tool failures occur with no visibility into root causes. Errors are returned to clients but not logged server-side, making debugging impossible. Without request/session tracking, failures from different operations cannot be correlated. The system needs structured logging that captures what happened, why it failed, and the context (request ID, session ID) to enable effective troubleshooting.

## What Changes

- Add a logging utility package with structured logging functions (Info, Warn, Error)
- Implement request ID generation (UUID) and propagation through context
- Extract and track session ID from HTTP headers
- Add conditional logging via DEBUG_LEVEL environment variable (INFO, DEBUG, TRACE)
- Log entry/exit of all tool handlers with inputs/outputs
- Log all errors server-side before returning to client
- Add operation-specific logging for filesystem operations (permissions, path validation, hash checks)
- Add operation-specific logging for command execution (startup, exit codes, timing)
- Add operation-specific logging for memory operations (writes, queries, search results)
- Ensure no raw sensitive data is logged (only error messages)

## Capabilities

### New Capabilities
- `structured-logging`: Logging infrastructure with context propagation, request/session tracking, and DEBUG_LEVEL control
- `filesystem-operation-logging`: Detailed logging for file read, write, edit operations including path validation and permission checks
- `command-execution-logging`: Detailed logging for command execution with timing, exit codes, and output sizes
- `memory-operation-logging`: Detailed logging for memory writes, queries, and searches

### Modified Capabilities

- `request-handling`: HTTP middleware now propagates request ID and session ID to tool handlers

## Impact

- HTTP middleware (main.go): Enhanced to generate and pass request IDs and session IDs
- All tool handlers (pkg/tools/*): Added entry/exit logging, error logging, operation-specific logging
- Tool helpers (pkg/tools/helpers.go): Enhanced errorResult to include more context
- New package: pkg/log/logger.go for structured logging utilities
