## Why

AI agents need safe, read-only file system exploration without shell access. Current `exec_command` tool provides unrestricted shell access, creating security risks for untrusted agents (prompt injection, data exfiltration via curl/wget). This change adds dedicated file search and discovery tools based on a pure Go Boyer-Moore implementation, eliminating the need for bash when exploring codebases.

## What Changes

- Add 4 new MCP tools for file exploration:
  - `search_code`: Fast text/regex search with line hashes for edit compatibility
  - `find_files`: Find files by name/path patterns
  - `list_directory`: List directory contents
  - `tree_directory`: ASCII tree visualization of directory structure
- Implement Boyer-Moore string search algorithm in pure Go (adapted from healeycodes/tools)
- All search results include CRC32 line hashes matching `read_file` format for seamless `edit_file` workflow
- Worker pool-based concurrent search for performance

## Capabilities

### New Capabilities
- `code-search`: Fast substring and regex search across files with hash-compatible output
- `file-discovery`: Find files by name patterns, list directories, visualize trees
- `boyer-moore-search`: Efficient string matching using Boyer-Moore algorithm with worker pool

### Modified Capabilities
<!-- None - these are net-new capabilities -->

## Impact

**Added:**
- `pkg/tools/filesearch/` package with search engine and directory tools
- `pkg/tools/filesearch/boyer_moore.go` - string matching algorithm
- `pkg/tools/filesearch/search.go` - concurrent search orchestration
- `pkg/tools/filesearch/directory.go` - directory listing and tree generation
- Request/response models in `pkg/models/models.go`
- Tool handlers in `pkg/tools/file_search.go`
- Tool registration in `internal/server.go` (26 total tools)

**Modified:**
- Tool count increases from 22 to 26 in server initialization

**Dependencies:**
- No new external dependencies (pure Go implementation)
