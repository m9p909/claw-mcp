# Claw MCP Server - Implementation Plan

## Overview

Building a Go-based MCP server exposing three tool categories: Filesystem (with hashline-based editing), Execution (background processes), and Memory (SQLite-backed).

**Tech Stack:**
- Language: Go
- Transport: HTTP via MCP Go SDK
- Database: SQLite at `~/.mcpclaw/data`
- Hashing: CRC32 (line-based)

---

## Phase 1: Project Setup & Foundation

### 1.1 Directory Structure
```
awesomeProject/
├─ main.go                           # HTTP server entry point
├─ go.mod
├─ go.sum
├─ internal/
│  ├─ server.go                      # MCP server initialization
│  ├─ config.go                      # Config + DB path setup
│  └─ db.go                          # SQLite initialization + migrations
├─ pkg/
│  ├─ tools/
│  │  ├─ filesystem.go               # read_file, write_file, edit_file
│  │  ├─ execution.go                # exec_command, manage_process
│  │  └─ memory.go                   # write_memory, query_memory, memory_search
│  ├─ storage/
│  │  ├─ memory.go                   # Memory table operations
│  │  └─ process.go                  # Process session tracking (in-memory)
│  ├─ hash/
│  │  └─ hash.go                     # CRC32 line hashing utilities
│  └─ models/
│     ├─ requests.go                 # Request structs with jsonschema tags
│     ├─ responses.go                # Response structs
│     └─ errors.go                   # Error code + message pattern
└─ README.md
```

### 1.2 Dependencies
```
go get github.com/modelcontextprotocol/go-sdk
go get github.com/mattn/go-sqlite3
```

### 1.3 Main.go Structure
```go
package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"

    "github.com/modelcontextprotocol/go-sdk/mcp"
    "awesomeProject/internal"
    "awesomeProject/pkg/tools"
)

func main() {
    port := flag.String("port", "3000", "HTTP port")
    flag.Parse()

    // Initialize DB
    db, err := internal.InitDB()
    if err != nil {
        log.Fatalf("Failed to init DB: %v", err)
    }
    defer db.Close()

    // Create MCP server
    server := mcp.NewServer(&mcp.Implementation{
        Name:    "claw",
        Version: "1.0.0",
    }, nil)

    // Register tools
    tools.RegisterFilesystemTools(server)
    tools.RegisterExecutionTools(server)
    tools.RegisterMemoryTools(server)

    // Start HTTP server
    handler := mcp.NewHTTPHandler(server)
    addr := fmt.Sprintf("localhost:%s", *port)
    log.Printf("Claw MCP server listening on http://%s", addr)

    if err := http.ListenAndServe(addr, handler); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

### 1.4 Config & DB Setup
- **Config struct**: DB path resolution (`~/.mcpclaw/data`)
- **Global DB instance**: Initialized in main, used by all tools
- **Migrations**: Create `memories` table on startup
  - Columns: `id (PK)`, `category`, `content`, `created_at`

---

## Phase 2: Hashing & Models

### 2.1 Hash Utilities (pkg/hash/hash.go)
```go
// ComputeLineHash returns 2-3 char hex hash of line content
func ComputeLineHash(line string) string

// FormatWithHashes adds inline hashes to content
func FormatWithHashes(content string) string
// Input: "line1\nline2\nline3"
// Output: "1:a3|line1\n2:f1|line2\n3:0e|line3"

// ExtractHash pulls hash from formatted line
func ExtractHash(formattedLine string) (string, string, error)
// Input: "1:a3|content here"
// Output: ("a3", "content here", nil)

// ValidateHashRange checks if hashes exist in file at positions
func ValidateHashRange(filepath string, startHash, endHash string) (int, int, error)
// Returns (startLine, endLine, error) or error if hashes don't match
```

### 2.2 Request/Response Models (pkg/models/)

**Filesystem requests:**
```go
type ReadFileRequest struct {
    FilePath string `json:"file_path" jsonschema:"Absolute path to file"`
    Offset   int    `json:"offset,omitempty" jsonschema:"Line offset (1-based)"`
    Limit    int    `json:"limit,omitempty" jsonschema:"Max lines to read"`
}

type WriteFileRequest struct {
    FilePath string `json:"file_path" jsonschema:"Absolute path to file"`
    Content  string `json:"content" jsonschema:"File content (may include hashes)"`
}

type EditFileRequest struct {
    FilePath   string `json:"file_path" jsonschema:"Absolute path to file"`
    Range      struct {
        StartHash string `json:"start_hash" jsonschema:"Hash of first line to replace"`
        EndHash   string `json:"end_hash" jsonschema:"Hash of last line to replace"`
    } `json:"range"`
    NewContent string `json:"new_content" jsonschema:"Replacement content"`
}
```

**Execution requests:**
```go
type ExecCommandRequest struct {
    Command   string            `json:"command" jsonschema:"Shell command to execute"`
    Cwd       string            `json:"cwd,omitempty" jsonschema:"Working directory"`
    Timeout   int               `json:"timeout,omitempty" jsonschema:"Timeout in seconds (default 60)"`
    Env       map[string]string `json:"env,omitempty" jsonschema:"Environment variables"`
    Background bool             `json:"background,omitempty" jsonschema:"Run in background"`
}

type ManageProcessRequest struct {
    Action    string `json:"action" jsonschema:"list|poll|send_keys|kill"`
    SessionID string `json:"session_id,omitempty" jsonschema:"Process session ID"`
    Input     string `json:"input,omitempty" jsonschema:"Input for send_keys"`
    Limit     int    `json:"limit,omitempty" jsonschema:"Output limit for log action"`
}
```

**Memory requests:**
```go
type WriteMemoryRequest struct {
    Category string `json:"category" jsonschema:"fact|todo|decision|preference"`
    Content  string `json:"content" jsonschema:"Memory content"`
}

type QueryMemoryRequest struct {
    SQL   string `json:"sql,omitempty" jsonschema:"SQL query"`
    Limit int    `json:"limit,omitempty" jsonschema:"Result limit (default 100)"`
}

type MemorySearchRequest struct {
    Query string `json:"query" jsonschema:"Search string (substring match)"`
    Limit int    `json:"limit,omitempty" jsonschema:"Result limit (default 10)"`
}
```

### 2.3 Error Pattern (pkg/models/errors.go)
```go
type ErrorResponse struct {
    Success bool   `json:"success"`
    Error   string `json:"error"`
    Message string `json:"message"`
}

// Error codes as constants
const (
    ErrFileNotFound      = "FileNotFoundError"
    ErrPermissionDenied  = "PermissionError"
    ErrHashMismatch      = "HashMismatchError"
    ErrNoMatchFound      = "NoMatchError"
    ErrSessionNotFound   = "SessionNotFoundError"
    ErrTimeout           = "TimeoutError"
    ErrSQLError          = "SQLError"
    ErrInvalidInput      = "InvalidInputError"
)

func NewErrorResponse(code, message string) *ErrorResponse
```

---

## Phase 3: Filesystem Tools

### 3.1 read_file
**Behavior:**
- Read file at `file_path`
- If `offset` + `limit` provided, slice lines
- Format output with inline hashes
- Return `lines_read`, `total_lines`, `truncated`

**Implementation:**
```go
func HandleReadFile(ctx context.Context, req *mcp.CallToolRequest, params *ReadFileRequest) (*mcp.CallToolResult, any, error)
```

**Error cases:**
- File not found → `FileNotFoundError`
- Permission denied → `PermissionError`

**Edge cases:**
- Offset > total lines → return empty, `truncated: false`
- Large files → still return all (no limits imposed)

### 3.2 write_file
**Behavior:**
- Create or overwrite file at `file_path`
- If content has hashes (format `1:a3|line...`), validate hashes match file before write
- If no existing file, skip hash validation
- Return `bytes_written`, `created`, `overwritten`

**Implementation:**
```go
func HandleWriteFile(ctx context.Context, req *mcp.CallToolRequest, params *WriteFileRequest) (*mcp.CallToolResult, any, error)
```

**Edge cases:**
- Content with hashes but file doesn't exist → strip hashes and write
- Content with hashes but file changed → return error
- Parent directory missing → return `PermissionError`

### 3.3 edit_file
**Behavior:**
- Find lines by hash range (strict match required)
- Replace those lines with `new_content`
- Return `replacements_made`, `line_number` (first affected line)

**Implementation:**
```go
func HandleEditFile(ctx context.Context, req *mcp.CallToolRequest, params *EditFileRequest) (*mcp.CallToolResult, any, error)
```

**Algorithm:**
1. Read file into memory
2. Find line with `start_hash` → record line number
3. Find line with `end_hash` → record line number
4. If either hash not found OR hashes don't match content → return `HashMismatchError`
5. Replace lines [start..end] with `new_content`
6. Write back to file
7. Return success with line number

**Error cases:**
- `HashMismatchError` → hash not found or content mismatch
- `FileNotFoundError` → file doesn't exist
- `PermissionError` → can't write

---

## Phase 4: Execution Tools

### 4.1 Process Session Model (pkg/storage/process.go)
```go
type ProcessSession struct {
    ID        string       // "exec_20240223_143052_a3f7"
    Command   string
    PID       int
    Status    string       // "running", "completed"
    ExitCode  *int
    Stdout    []string     // Unlimited buffer
    Stderr    []string     // Unlimited buffer
    StartedAt time.Time
    EndedAt   *time.Time
}

// In-memory map, NOT persisted
var globalSessions map[string]*ProcessSession
```

### 4.2 exec_command
**Behavior (foreground):**
- Run command synchronously
- Capture all stdout/stderr
- Return exit code + output

**Behavior (background):**
- Spawn process
- Return session ID immediately
- Store session in memory
- Process runs detached

**Implementation:**
```go
func HandleExecCommand(ctx context.Context, req *mcp.CallToolRequest, params *ExecCommandRequest) (*mcp.CallToolResult, any, error)
```

**Details:**
- Use `os/exec.CommandContext(ctx, ...)`
- Timeout via context cancellation
- Default timeout: 60s
- Custom env: merge with `os.Environ()`
- CWD: use `cmd.Dir = cwd` or default to workspace root

**Error cases:**
- Command not found → `error launching process`
- Timeout → return partial output + `TimeoutError`

### 4.3 manage_process
**Actions:**

**list:**
- Return all sessions (active + completed)

**poll:**
- Check session by ID
- If still running, return `status: "running"`
- If completed, return `status: "completed"`, `exit_code`, `runtime_seconds`

**send_keys:**
- Write input to process stdin
- Only works for running processes
- Return bytes sent

**kill:**
- Send SIGTERM to process
- Update status to "killed"
- Return signal sent

**log:**
- Return last N lines of stdout + stderr
- `limit` defaults to 100
- Return `lines_returned`, `total_lines`

**Implementation:**
```go
func HandleManageProcess(ctx context.Context, req *mcp.CallToolRequest, params *ManageProcessRequest) (*mcp.CallToolResult, any, error)
```

**Error cases:**
- Session not found → `SessionNotFoundError`
- Action invalid → `InvalidInputError`

---

## Phase 5: Memory Tools

### 5.1 Storage Layer (pkg/storage/memory.go)
```go
// DB operations
func WriteMemory(category, content string) (int, error)
func QueryMemory(sql string, limit int) ([]Memory, error)
func GetAllMemories() ([]Memory, error)

type Memory struct {
    ID        int       `json:"id"`
    Category  string    `json:"category"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}
```

**DB Schema:**
```sql
CREATE TABLE IF NOT EXISTS memories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 5.2 write_memory
**Behavior:**
- Insert into `memories` table
- Return `id`, `category`, `content`, `created_at`

**Validation:**
- Category must be one of: `fact`, `todo`, `decision`, `preference`
- Content required

### 5.3 query_memory
**Behavior:**
- Execute SQL query against memories table
- Apply limit (default 100)
- Return results as JSON array

**Error cases:**
- SQL syntax error → `SQLError` with message
- Invalid query → `SQLError`

**Allowed queries:**
- SELECT only (no INSERT/UPDATE/DELETE)
- Can filter by category, date, etc.

### 5.4 memory_search
**Behavior:**
- Load ALL memories into memory
- Filter by substring match (case-insensitive)
- Sort by relevance (exact match first, then matches at start of content)
- Return top N results

**Implementation:**
```go
func MemorySearch(query string, limit int) []Memory {
    all, _ := GetAllMemories()
    var matches []Memory

    // Substring matching (case-insensitive)
    for _, m := range all {
        if strings.Contains(strings.ToLower(m.Content), strings.ToLower(query)) {
            matches = append(matches, m)
        }
    }

    // Simple sorting: exact match first
    sort.SliceStable(matches, func(i, j int) bool {
        iExact := strings.ToLower(m.Content) == strings.ToLower(query)
        jExact := strings.ToLower(m.Content) == strings.ToLower(query)
        if iExact != jExact {
            return iExact
        }
        return false
    })

    return matches[:min(len(matches), limit)]
}
```

---

## Phase 6: Tooling & Integration

### 6.1 Tool Registration Pattern
Each tool package exports a register function:

```go
// pkg/tools/filesystem.go
func RegisterFilesystemTools(server *mcp.Server) {
    mcp.AddTool(server, &mcp.Tool{
        Name:        "read_file",
        Description: "Read file contents with line hashes",
    }, HandleReadFile)

    mcp.AddTool(server, &mcp.Tool{
        Name:        "write_file",
        Description: "Write file (validates hashes if present)",
    }, HandleWriteFile)

    mcp.AddTool(server, &mcp.Tool{
        Name:        "edit_file",
        Description: "Edit file using hashline ranges",
    }, HandleEditFile)
}
```

### 6.2 Logging
- All tools log to stdout via `log.Printf()`
- Log format: `[TOOL_NAME] action: details`
- Example: `[read_file] Reading /home/user/file.go (500 lines)`

---

## Implementation Order

1. **main.go + internal/config.go + internal/db.go** - Server setup
2. **pkg/hash/hash.go** - Hash utilities
3. **pkg/models/** - Request/response structs
4. **pkg/storage/memory.go** - Memory table operations
5. **pkg/tools/filesystem.go** - Filesystem tools
6. **pkg/tools/execution.go** - Execution tools + process tracking
7. **pkg/tools/memory.go** - Memory tools
8. **Testing** - Integration test with sample requests

---

## Testing Strategy

- Manual HTTP requests via `curl` or Postman
- Sample requests matching mcp_api.md spec
- Verify hashes work correctly for edit_file
- Verify process backgrounding works
- Verify memory search filtering works

---

## Open Questions Resolved

1. ✅ Hash format: CRC32, 2-3 char hex, line-based
2. ✅ Edit format: Range-based with start/end hash
3. ✅ Write validation: Validate hashes match before overwriting
4. ✅ DB access: Global instance
5. ✅ Logging: stdout
6. ✅ Error pattern: Code + message
7. ✅ MCP registration: Struct-based with jsonschema tags
8. ✅ Process persistence: In-memory only, no disk
9. ✅ Memory search: Substring match, in-memory filtering