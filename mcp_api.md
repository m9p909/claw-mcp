# Claw MCP Server - API Specification

Complete request/response schemas for all tools.

---

## Filesystem Tools

### read_file

**Request:**
```json
{
  "file_path": "/home/user/project/main.py",
  "offset": 1,
  "limit": 100
}
```

**Required:** `file_path`
**Optional:** `offset` (default: 1), `limit` (default: 2000 lines or 50KB)

**Success Response:**
```json
{
  "success": true,
  "content": "import os\nimport sys\n...",
  "file_path": "/home/user/project/main.py",
  "lines_read": 50,
  "total_lines": 500,
  "truncated": true
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "FileNotFoundError",
  "message": "No such file or directory: '/home/user/project/main.py'"
}
```

---

### write_file

**Request:**
```json
{
  "file_path": "/home/user/project/new_file.py",
  "content": "def main():\n    print('hello')\n"
}
```

**Required:** `file_path`, `content`

**Success Response:**
```json
{
  "success": true,
  "file_path": "/home/user/project/new_file.py",
  "bytes_written": 32,
  "created": true,
  "overwritten": false
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "PermissionError",
  "message": "Permission denied: '/root/protected.txt'"
}
```

---

### edit_file

**Request:**
```json
{
  "file_path": "/home/user/project/config.py",
  "old_text": "DEBUG = False",
  "new_text": "DEBUG = True"
}
```

**Required:** `file_path`, `old_text`, `new_text`

**Success Response:**
```json
{
  "success": true,
  "file_path": "/home/user/project/config.py",
  "replacements_made": 1,
  "line_number": 15
}
```

**Error Response (no match):**
```json
{
  "success": false,
  "error": "NoMatchError",
  "message": "Old text not found in file. Text must match exactly including whitespace."
}
```

**Error Response (multiple matches):**
```json
{
  "success": false,
  "error": "AmbiguousMatchError",
  "message": "Old text found 3 times in file. Use more specific text to match unique location."
}
```

---

## Execution Tools

### exec_command

**Request (foreground):**
```json
{
  "command": "ls -la",
  "cwd": "/home/user/project",
  "timeout": 30,
  "env": {
    "MY_VAR": "value"
  }
}
```

**Request (background):**
```json
{
  "command": "python3 train_model.py --epochs 100",
  "cwd": "/home/user/project",
  "background": true
}
```

**Request (PTY mode for TUI):**
```json
{
  "command": "vim main.py",
  "pty": true,
  "background": true
}
```

**Required:** `command`
**Optional:** `cwd` (default: workspace root), `background` (default: false), `pty` (default: false), `timeout` (default: 60s), `env` (default: {})

**Success Response (foreground):**
```json
{
  "success": true,
  "exit_code": 0,
  "stdout": "total 32\ndrwxr-xr-x 3 user user 4096 Feb 23 10:00 .\n...",
  "stderr": "",
  "duration_ms": 150
}
```

**Success Response (background):**
```json
{
  "success": true,
  "session_id": "exec_20240223_143052_a3f7",
  "command": "python3 train_model.py --epochs 100",
  "pid": 12345,
  "status": "running",
  "started_at": "2024-02-23T14:30:52Z"
}
```

**Error Response (timeout):**
```json
{
  "success": false,
  "error": "TimeoutError",
  "message": "Command timed out after 30 seconds",
  "stdout": "Partial output...",
  "stderr": "",
  "exit_code": null
}
```

---

### manage_process

**Request (list):**
```json
{
  "action": "list"
}
```

**Request (poll):**
```json
{
  "action": "poll",
  "session_id": "exec_20240223_143052_a3f7"
}
```

**Request (log):**
```json
{
  "action": "log",
  "session_id": "exec_20240223_143052_a3f7",
  "limit": 100
}
```

**Request (send_keys):**
```json
{
  "action": "send_keys",
  "session_id": "exec_20240223_143052_a3f7",
  "input": "Hello World"
}
```

**Request (kill):**
```json
{
  "action": "kill",
  "session_id": "exec_20240223_143052_a3f7"
}
```

**Required:** `action`
**Required for non-list:** `session_id`
**Optional for send_keys:** `input`
**Optional for log:** `limit` (default: 100)

**Success Response (list):**
```json
{
  "success": true,
  "sessions": [
    {
      "session_id": "exec_20240223_143052_a3f7",
      "command": "python3 train_model.py --epochs 100",
      "status": "running",
      "pid": 12345,
      "started_at": "2024-02-23T14:30:52Z",
      "pty": false
    }
  ],
  "count": 1
}
```

**Success Response (poll):**
```json
{
  "success": true,
  "session_id": "exec_20240223_143052_a3f7",
  "status": "completed",
  "exit_code": 0,
  "pid": 12345,
  "started_at": "2024-02-23T14:30:52Z",
  "ended_at": "2024-02-23T14:35:12Z",
  "runtime_seconds": 260
}
```

**Success Response (log):**
```json
{
  "success": true,
  "session_id": "exec_20240223_143052_a3f7",
  "stdout": [
    "Epoch 1/100: loss=0.542",
    "Epoch 2/100: loss=0.498",
    "..."
  ],
  "stderr": [
    "Warning: Using default learning rate"
  ],
  "lines_returned": 50,
  "total_lines": 250
}
```

**Success Response (send_keys):**
```json
{
  "success": true,
  "session_id": "exec_20240223_143052_a3f7",
  "bytes_sent": 11
}
```

**Success Response (kill):**
```json
{
  "success": true,
  "session_id": "exec_20240223_143052_a3f7",
  "status": "killed",
  "signal": "SIGTERM"
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "SessionNotFoundError",
  "message": "No session found with ID: exec_20240223_143052_a3f7"
}
```

---

## Memory Tools

### query_memory

**Request:**
```json
{
  "sql": "SELECT * FROM memories WHERE category = 'preference' ORDER BY created_at DESC LIMIT 10",
  "limit": 10
}
```

**Request (semantic search when implemented):**
```json
{
  "semantic_query": "What does the user like for breakfast?",
  "limit": 5
}
```

**Required:** One of `sql` or `semantic_query`
**Optional:** `limit` (default: 100)

**Success Response:**
```json
{
  "success": true,
  "query": "SELECT * FROM memories WHERE category = 'preference' ORDER BY created_at DESC LIMIT 10",
  "results": [
    {
      "id": 42,
      "session_key": "session_abc123",
      "created_at": "2024-02-23T10:30:00Z",
      "category": "preference",
      "content": "User prefers oatmeal for breakfast"
    },
    {
      "id": 38,
      "session_key": "session_abc123",
      "created_at": "2024-02-22T15:45:00Z",
      "category": "preference",
      "content": "User likes reading sci-fi books"
    }
  ],
  "count": 2
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "SQLError",
  "message": "no such column: category"
}
```

---

### write_memory

**Request:**
```json
{
  "category": "fact",
  "content": "User is building an MCP server in TypeScript",
  "session_key": "session_abc123"
}
```

**Required:** `category`, `content`
**Optional:** `session_key` (auto-generated if not provided)

**Categories:** `fact`, `todo`, `decision`, `preference`

**Success Response:**
```json
{
  "success": true,
  "id": 43,
  "category": "fact",
  "content": "User is building an MCP server in TypeScript",
  "session_key": "session_abc123",
  "created_at": "2024-02-23T14:35:00Z"
}
```

---

### memory_search

**Request:**
```json
{
  "query": "breakfast preferences",
  "limit": 5
}
```

**Required:** `query`
**Optional:** `limit` (default: 10)

**Success Response:**
```json
{
  "success": true,
  "query": "breakfast preferences",
  "results": [
    {
      "id": 42,
      "created_at": "2024-02-23T10:30:00Z",
      "category": "preference",
      "content": "User prefers oatmeal for breakfast",
      "similarity": 0.89
    }
  ],
  "count": 1
}
```

**Note:** Returns empty results with `success: true` if no matches. Requires embedding support.

---

## Context Tools

### load_soul

**Request:**
```json
{
  "soul_path": "./SOUL.md"
}
```

**Optional:** `soul_path` (default: "./SOUL.md")

**Success Response:**
```json
{
  "success": true,
  "soul_path": "/home/user/workspace/SOUL.md",
  "content": "# SOUL.md - Who You Are\n\nYou're a helpful AI assistant...",
  "loaded": true
}
```

**Fallback Response (file not found):**
```json
{
  "success": true,
  "soul_path": "./SOUL.md",
  "content": "# Default Soul\n\nYou are a helpful AI assistant with access to file system, execution, and memory tools.",
  "loaded": false,
  "note": "SOUL.md not found, using default persona"
}
```

---

### get_current_time

**Request:**
```json
{
  "timezone": "America/Toronto"
}
```

**Optional:** `timezone` (default: "America/Toronto")

**Success Response:**
```json
{
  "success": true,
  "iso": "2024-02-23T14:35:00-05:00",
  "formatted": "Friday, February 23, 2024 — 02:35 PM EST",
  "timezone": "America/Toronto",
  "unix_timestamp": 1708700100
}
```

**Error Response (invalid timezone):**
```json
{
  "success": false,
  "error": "InvalidTimezoneError",
  "message": "Unknown timezone: 'Mars/Phobos'"
}
```

---

## Session Management Tools

### list_sessions

**Request:**
```json
{
  "session_key": "session_abc123"
}
```

**Optional:** `session_key` (filters results)

**Success Response:**
```json
{
  "success": true,
  "sessions": [
    {
      "session_key": "session_abc123",
      "created_at": "2024-02-23T10:00:00Z",
      "last_active": "2024-02-23T14:35:00Z",
      "memory_count": 15
    }
  ],
  "count": 1
}
```

---

### send_to_session

**Request:**
```json
{
  "target_session": "session_abc123",
  "message": "Background task completed successfully"
}
```

**Required:** `target_session`, `message`

**Success Response:**
```json
{
  "success": true,
  "target_session": "session_abc123",
  "message": "Background task completed successfully",
  "delivered": true
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "SessionNotFoundError",
  "message": "Target session 'session_abc123' not found or inactive"
}
```

---

## Common Error Types

| Error Code | Description |
|------------|-------------|
| `FileNotFoundError` | File or directory doesn't exist |
| `PermissionError` | Insufficient permissions |
| `NoMatchError` | Text not found for edit_file |
| `AmbiguousMatchError` | Multiple matches for edit_file |
| `TimeoutError` | Command exceeded timeout |
| `SQLError` | Invalid SQL syntax or query |
| `SessionNotFoundError` | Invalid or expired session ID |
| `InvalidTimezoneError` | Unknown timezone name |

---

## Transport Notes

**MCP Protocol:** Tools are exposed via MCP's JSON-RPC interface. The actual request/response format follows MCP's tool call schema:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "exec_command",
    "arguments": {
      "command": "ls -la",
      "background": false
    }
  }
}
```

The schemas above represent the `arguments` object and the tool's internal response (which MCP wraps in its own envelope).