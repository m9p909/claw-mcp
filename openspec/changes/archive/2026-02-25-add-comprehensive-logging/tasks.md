## 1. Logging Infrastructure

- [x] 1.1 Create pkg/log/logger.go with Logger type and slog initialization
- [x] 1.2 Implement logger.Info(), logger.Warn(), logger.Error() functions with context support
- [x] 1.3 Add DEBUG_LEVEL environment variable support (INFO, DEBUG, TRACE)
- [x] 1.4 Create context helpers: WithRequestID(), WithSessionID(), requestIDFromContext(), sessionIDFromContext()


## 2. HTTP Middleware and Request Context

- [x] 2.1 Update main.go to initialize logger on startup
- [x] 2.2 Generate UUID request IDs in HTTP middleware
- [x] 2.3 Extract Mcp-Session-Id header and add to context
- [x] 2.4 Pass context with request ID and session ID to all tool handlers
- [x] 2.5 Add logging to HTTP middleware (method, path, request ID, session ID)

## 3. Tool Helper Functions

- [x] 3.1 Update pkg/tools/helpers.go errorResult() to log errors server-side before returning
- [x] 3.2 Add sanitization for file paths in error logs (mask actual paths, log operation type instead)

## 4. Filesystem Tool Logging

- [x] 4.1 Add entry/exit logging to HandleReadFile (path, result count, timing)
- [x] 4.2 Add operation logging to HandleReadFile (permission checks at DEBUG level)
- [x] 4.3 Add entry/exit logging to HandleWriteFile (path, bytes written, timing)
- [x] 4.4 Add operation logging to HandleWriteFile (path validation, hash validation at DEBUG level)
- [x] 4.5 Add entry/exit logging to HandleEditFile (start/end hash lookup, line count, timing)
- [x] 4.6 Add operation logging to HandleEditFile (hash validation at DEBUG level)

## 5. Command Execution Tool Logging

- [x] 5.1 Add entry/exit logging to HandleExecCommand for foreground execution
- [x] 5.2 Add operation logging to HandleExecCommand foreground (stdout/stderr sizes, exit code, timing)
- [x] 5.3 Add entry/exit logging to HandleExecCommand for background execution (session creation)
- [x] 5.4 Add goroutine logging to HandleExecCommand background (process monitoring at TRACE level)
- [x] 5.5 Add entry/exit logging to HandleManageProcess for each action (list, poll, send_keys, kill)
- [x] 5.6 Add operation logging to HandleManageProcess for session lookups (result counts, session status)

## 6. Memory Tool Logging

- [x] 6.1 Add entry/exit logging to HandleWriteMemory (category, content size, timing)
- [x] 6.2 Add entry/exit logging to HandleQueryMemory (result count, timing)
- [x] 6.3 Add operation logging to HandleQueryMemory (query string structure at DEBUG level)
- [x] 6.4 Add entry/exit logging to HandleMemorySearch (result count, limit applied, timing)

## 7. Testing and Validation

- [x] 7.1 Verify all tool handlers include entry/exit logging
- [x] 7.2 Verify all error paths log errors server-side
- [x] 7.3 Test DEBUG_LEVEL=INFO shows only important logs
- [x] 7.4 Test DEBUG_LEVEL=DEBUG shows operation details
- [x] 7.5 Test DEBUG_LEVEL=TRACE shows low-level details
- [x] 7.6 Verify request IDs and session IDs appear in all logs
- [x] 7.7 Verify no raw file paths or sensitive data in logs
