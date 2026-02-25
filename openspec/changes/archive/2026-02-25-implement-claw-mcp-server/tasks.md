## 1. Project Setup

- [x] 1.1 Initialize Go modules and add dependencies (MCP SDK, SQLite)
- [x] 1.2 Create directory structure (main.go, internal/, pkg/, openspec/)
- [x] 1.3 Set up global database instance and initialization at startup

## 2. Hash Utilities & Models

- [x] 2.1 Implement CRC32 line hashing (pkg/hash/hash.go)
- [x] 2.2 Implement hash formatting/extraction utilities
- [x] 2.3 Create request/response structs with jsonschema tags (pkg/models/)
- [x] 2.4 Create error response types and constants

## 3. Database & Storage

- [x] 3.1 Implement SQLite initialization and migrations (internal/db.go)
- [x] 3.2 Create memories table schema (id, category, content, created_at)
- [x] 3.3 Implement memory storage operations (pkg/storage/memory.go)
- [x] 3.4 Implement in-memory process session tracking (pkg/storage/process.go)

## 4. MCP Server Setup

- [x] 4.1 Create MCP server initialization (internal/server.go)
- [x] 4.2 Implement HTTP handler setup and tool registration
- [x] 4.3 Create main.go with server startup and CLI flags

## 5. Filesystem Tools

- [x] 5.1 Implement read_file handler with hash generation
- [x] 5.2 Implement write_file handler with hash validation
- [x] 5.3 Implement edit_file handler with hash range validation
- [x] 5.4 Register filesystem tools with MCP server

## 6. Execution Tools

- [x] 6.1 Implement exec_command handler (foreground and background)
- [x] 6.2 Implement background process spawning with session tracking
- [x] 6.3 Implement manage_process handler (list, poll, send_keys, kill)
- [x] 6.4 Implement process output capture (stdout/stderr buffering)
- [x] 6.5 Register execution tools with MCP server

## 7. Memory Tools

- [x] 7.1 Implement write_memory handler with validation
- [x] 7.2 Implement query_memory handler (SQL execution, SELECT-only)
- [x] 7.3 Implement memory_search handler (substring matching)
- [x] 7.4 Register memory tools with MCP server

## 8. Testing & Validation

- [x] 8.1 Test filesystem tools (read/write/edit with hashes)
- [x] 8.2 Test execution tools (foreground, background, process management)
- [x] 8.3 Test memory tools (write, query, search)
- [x] 8.4 Test error cases and validation
- [x] 8.5 Verify HTTP server startup and tool registration
