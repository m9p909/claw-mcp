## Context

The MCP (Model Context Protocol) server exposes filesystem, execution, and memory tools via HTTP. This enables AI models (like Claude) to interact with a local development environment in a controlled, standardized way. The server is the gateway between LLM requests and system-level operations.

## Goals / Non-Goals

**Goals:**
- Implement complete MCP HTTP server with 9 tools across 3 categories
- Use hashline-based editing to eliminate whitespace brittleness in file edits
- Provide persistent memory storage with lightweight in-memory search
- Handle background process execution with polling/management
- Maintain clean, minimal architecture (max 5 functions per module)

**Non-Goals:**
- Multi-session/user support (single session only)
- Semantic embeddings or vector search (substring matching sufficient)
- PTY/TUI support for interactive processes
- Persistent process state (restarts lose background jobs)
- Fine-grained access control or sandboxing

## Decisions

### 1. HTTP Transport via MCP Go SDK
**Decision:** Use official MCP Go SDK with HTTP handler instead of custom JSON-RPC.

**Rationale:**
- Spec-compliant, maintained by MCP authors
- Automatic schema generation from Go structs + jsonschema tags
- Reduces boilerplate, avoids reimplementing protocol logic

**Alternatives Considered:**
- Custom JSON-RPC: More control, but duplicates protocol logic
- Stdio transport: Simpler protocol but incompatible with remote use cases

### 2. Hashline-Based File Editing
**Decision:** Every line gets a 2-3 character CRC32 hash. Edit requests use hash ranges, not text matching.

**Rationale:**
- Eliminates whitespace brittleness (exact string matching fails on indentation changes)
- Provides stable anchor even if file content shifts slightly
- Matches blog.can.ac performance improvements across 16 LLM models
- Hash mismatch = safe rejection, preventing accidental corruption

**Alternatives Considered:**
- String replace (status quo): High failure rate on whitespace
- Patch format (unified diff): Complex to parse, model-specific
- Line numbers only: No validation, unsafe if file changes between read/edit

### 3. SQLite for Memory Storage
**Decision:** Use embedded SQLite at `~/.mcpclaw/data`, load all data for in-memory search.

**Rationale:**
- No external service dependency (dev-friendly)
- File-based persistence survives restarts
- In-memory filtering is fast for typical memory sizes (<100k entries)
- SQL queries support flexible filtering

**Alternatives Considered:**
- PostgreSQL: Overkill, adds deployment complexity
- JSON file: No query capability
- Redis: Requires running service, memory-only (not persistent by default)

### 4. In-Memory Process Management
**Decision:** Store running processes in a global map, not persisted to disk.

**Rationale:**
- Processes die on server restart anyway
- In-memory tracking avoids DB round-trips for polling
- Simple (single global map), minimal overhead
- Matches typical use case (short-lived background tasks)

**Alternatives Considered:**
- Disk-persisted sessions: Adds complexity, stale state on crashes
- Pub/sub for notifications: Overkill for single-session scenario

### 5. Substring Search for Memory
**Decision:** Load all memories, filter client-side with case-insensitive substring match.

**Rationale:**
- No indexing complexity
- Fast for typical corpus sizes
- Predictable performance
- Can add ranking later (exact matches first)

**Alternatives Considered:**
- Full-text search: More complex, overkill for simple queries
- Semantic search: Requires embeddings, out of scope

### 6. Global DB Instance
**Decision:** Single global SQLite connection initialized at startup.

**Rationale:**
- Simplifies tool implementations (no dependency injection)
- SQLite handles concurrent access via locking
- Connection pooling unnecessary for single server

**Alternatives Considered:**
- Per-request DB connections: More overhead
- Dependency injection: More verbose, no real benefit here

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| **Large file reads/writes** (>1GB) consume significant memory | No limits imposed; document as expected behavior. Users responsible for filesize. |
| **Unlimited process output buffering** causes OOM on long-running tasks | Document limitation. Kill process if output exceeds memory. Future: implement circular buffer. |
| **Hash collisions** on CRC32 (unlikely but possible) cause edit rejection | Collision rate negligible for typical file sizes. Can switch to longer hash if needed. |
| **In-memory process state lost on restart** breaks long-running tasks | Expected behavior. Sessions are transient. Document clearly. |
| **No access control** exposes full filesystem/execution to HTTP clients | For local-only use. Document security risk for remote deployment. |

## Open Questions

1. Should hashes use CRC32 or full SHA256 prefix? (CRC32 sufficient for now)
2. Should we validate hashes match file before write_file, or only for edit_file? (Both, to catch user errors)
3. Process output size limits? (Unlimited for now, can add later)
