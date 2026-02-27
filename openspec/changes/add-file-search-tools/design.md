## Context

Current MCP server has 22 tools including `exec_command` which provides unrestricted shell access. For code RAG systems and untrusted agents, this creates security risks (prompt injection → curl/wget data exfiltration, arbitrary command execution). Existing filesystem tools (`read_file`, `write_file`, `edit_file`) use CRC32 line hashing for precise editing but lack discovery capabilities.

We explored embedding ripgrep binary but found a pure Go implementation (healeycodes/tools) that adapts Go's internal Boyer-Moore algorithm. This gives us fast search without external dependencies.

## Goals / Non-Goals

**Goals:**
- Enable safe file exploration without shell access
- Maintain hash compatibility between `search_code` and `read_file`/`edit_file` workflow
- Achieve acceptable performance for typical codebases (pure Go implementation)
- Provide 4 distinct tools: search, find, list, tree
- Support both literal string and regex matching

**Non-Goals:**
- Feature parity with ripgrep (no .gitignore awareness, file type filters - can add later)
- Multi-platform support (Linux AMD64 only for now)
- Replacing existing `exec_command` tool (keep for trusted agents)
- Performance matching compiled ripgrep (pure Go trade-off accepted)

## Decisions

### 1. Pure Go vs Embedded Binary
**Decision:** Use pure Go Boyer-Moore implementation adapted from healeycodes/tools

**Alternatives:**
- Embed ripgrep binary (initially proposed)
- Use existing Go search libraries (couldn't find fast enough options)

**Rationale:**
- No external binaries to manage/extract
- ~300 LOC, easy to maintain and customize
- Boyer-Moore algorithm is proven fast for substring search
- Can add incremental features (filters, ignore files) as needed

### 2. Hash Calculation Strategy
**Decision:** Read actual file lines to calculate CRC32 hashes for search results

**Alternatives:**
- Hash the ripgrep/search output directly (faster but inconsistent)

**Rationale:**
- Guarantees hash compatibility with `read_file` output
- Critical for `edit_file` workflow (edit_file requires exact hash matches)
- Slight performance trade-off acceptable for correctness

### 3. Package Structure
**Decision:** Create `pkg/tools/filesearch/` package with 3 files:
- `boyer_moore.go` - string matching algorithm (adapted from healeycodes)
- `search.go` - concurrent search orchestration with worker pool
- `directory.go` - list/tree directory operations

**Rationale:**
- Separates search algorithm from directory operations
- Keeps boyer_moore.go isolated for future replacement/optimization
- Follows existing package structure in `pkg/tools/`

### 4. Worker Pool Size
**Decision:** Default 128 workers (matching healeycodes implementation), configurable

**Rationale:**
- Proven default from source implementation
- File I/O bound workload benefits from high concurrency
- Can tune based on performance testing

### 5. Binary File Detection
**Decision:** Check first buffer for NUL byte, report "Binary file X matches" and skip

**Rationale:**
- Matches ripgrep behavior
- Prevents printing unprintable content
- Fast detection (only first read needed)

## Risks / Trade-offs

**[Risk] Pure Go slower than ripgrep** → Acceptable trade-off for zero dependencies. Can revisit with embedded binary if performance unacceptable in practice.

**[Risk] No .gitignore awareness** → Future enhancement. Initial version searches all files. Can add ignore logic incrementally.

**[Risk] Memory usage on large directories** → Worker pool limits concurrency. Can add max file size limits if needed.

**[Risk] Hash calculation doubles I/O** → Read file twice (once for search, once for hash). Optimization: cache line content during search. Not critical for initial version.

**[Trade-off] Separate tools vs unified** → More tools (26 vs 22) but clearer API. Each tool has single responsibility.

## Migration Plan

No migration needed - this is additive. New tools register alongside existing ones. Agents can continue using `exec_command` or switch to new search tools.

**Deployment:**
1. Add new tool handlers and models
2. Register in `internal/server.go`
3. Restart MCP server
4. Tool count increases from 22 → 26

**Rollback:** Remove tool registrations, restart server.

## Open Questions

None - design is complete for initial implementation.
