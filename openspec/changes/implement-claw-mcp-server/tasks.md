## 1. Project Setup

- [ ] 1.1 Initialize Go modules and add dependencies (MCP SDK, SQLite)
- [ ] 1.2 Create directory structure (main.go, internal/, pkg/, openspec/)
- [ ] 1.3 Set up global database instance and initialization at startup

## 2. Hash Utilities & Models

- [ ] 2.1 Implement CRC32 line hashing (pkg/hash/hash.go)
- [ ] 2.2 Implement hash formatting/extraction utilities
- [ ] 2.3 Create request/response structs with jsonschema tags (pkg/models/)
- [ ] 2.4 Create error response types and constants

## 3. Database & Storage

- [ ] 3.1 Implement SQLite initialization and migrations (internal/db.go)
- [ ] 3.2 Create memories table schema (id, category, content, created_at)
- [ ] 3.3 Implement memory storage operations (pkg/storage/memory.go)
- [ ] 3.4 Implement in-memory process session tracking (pkg/storage/process.go)

## 4. MCP Server Setup

- [ ] 4.1 Create MCP server initialization (internal/server.go)
- [ ] 4.2 Implement HTTP handler setup and tool registration
- [ ] 4.3 Create main.go with server startup and CLI flags

## 5. Filesystem Tools

- [ ] 5.1 Implement read_file handler with hash generation
- [ ] 5.2 Implement write_file handler with hash validation
- [ ] 5.3 Implement edit_file handler with hash range validation
- [ ] 5.4 Register filesystem tools with MCP server

## 6. Execution Tools

- [ ] 6.1 Implement exec_command handler (foreground and background)
- [ ] 6.2 Implement background process spawning with session tracking
- [ ] 6.3 Implement manage_process handler (list, poll, send_keys, kill)
- [ ] 6.4 Implement process output capture (stdout/stderr buffering)
- [ ] 6.5 Register execution tools with MCP server

## 7. Memory Tools

- [ ] 7.1 Implement write_memory handler with validation
- [ ] 7.2 Implement query_memory handler (SQL execution, SELECT-only)
- [ ] 7.3 Implement memory_search handler (substring matching)
- [ ] 7.4 Register memory tools with MCP server

## 8. Testing & Validation

- [ ] 8.1 Test filesystem tools (read/write/edit with hashes)
- [ ] 8.2 Test execution tools (foreground, background, process management)
- [ ] 8.3 Test memory tools (write, query, search)
- [ ] 8.4 Test error cases and validation
- [ ] 8.5 Verify HTTP server startup and tool registration
