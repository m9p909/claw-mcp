## Context

The project is an MCP (Model Context Protocol) server in Go that exposes tools to agents. Currently, tools are hardcoded and registered in `internal/server.go`. Agents need a way to discover skill packages (Agent Skills format) stored in `~/.mcpclaw/skills/`. Each skill is a directory containing `SKILL.md` (required, with YAML frontmatter) and optional `scripts/`, `references/`, `assets/` subdirectories.

The Agent Skills specification defines strict validation rules for the YAML frontmatter, including name format (lowercase alphanumeric + hyphens, no leading/trailing/consecutive hyphens), required description field, optional metadata/version/license/compatibility fields, and constraints on field lengths.

## Goals / Non-Goals

**Goals:**
- Implement `list_skills` tool to scan `~/.mcpclaw/skills/` and return validated skill metadata
- Implement `get_skill` tool to retrieve full skill content (frontmatter + markdown body)
- Validate all skills strictly per Agent Skills spec; skip invalid ones with warning logs
- Return `skill_directory` path in both responses so agents can bash into scripts/
- Support agents invoking skill scripts via `exec_command` using the returned directory path

**Non-Goals:**
- Remote skill registries or network fetching
- Caching of parsed skills
- Installing or managing skill dependencies
- Direct script execution (scripts invoked via bash, not direct tool calls)

## Decisions

### 1. YAML Parsing Library
**Decision**: Use `gopkg.in/yaml.v3` library for frontmatter parsing and validation.

**Rationale**:
- Avoids manual string parsing (error-prone, unmaintainable)
- Library is widely used and well-tested
- Handles complex nested structures (metadata maps) correctly
- Single new dependency, acceptable trade-off for clarity

**Alternatives Considered**:
- Manual frontmatter parsing: More control, but fragile; parsing YAML by hand is error-prone
- Standard library encoding/json: Only supports JSON, not YAML

### 2. Frontmatter Structure Validation
**Decision**: Require strict `---\n...\n---\n` structure. Malformed skills (missing closing `---`, invalid YAML) are skipped with warning logs; they do NOT appear in list_skills results.

**Rationale**:
- Spec defines exact structure; strict validation matches intent
- Malformed skills indicate authoring errors; skipping them prevents agents from encountering broken skills
- Warning logs allow operator to investigate and fix
- Clean API surface: agents only see valid, usable skills

**Alternatives Considered**:
- Lenient parsing: Accept malformed frontmatter if YAML is valid. Rejected because it masks authoring errors and could lead to incomplete skill metadata.

### 3. Metadata Preservation
**Decision**: Preserve metadata structure exactly as nested map[string]string from YAML. Do NOT flatten.

**Rationale**:
- Agent Skills spec allows arbitrary nesting in metadata
- Preserving structure lets skill authors organize metadata as needed (e.g., nested org.example.property)
- Agents can work with the full structure

### 4. Home Directory Expansion
**Decision**: Expand `~/.mcpclaw/skills/` to absolute path at tool invocation time using `os.UserHomeDir()` + `filepath.Join()`.

**Rationale**:
- Fresh expansion each call is simple, negligible performance cost
- No initialization state to manage
- Handles cases where `$HOME` changes (edge case, but cleaner to support)

### 5. Error Handling for Missing Skills
**Decision**: `get_skill` with non-existent skill name returns error result with descriptive message (e.g., "skill 'pdf-processing' not found in ~/.mcpclaw/skills/").

**Rationale**:
- Agents need clear feedback about what went wrong
- Helps debug typos or skills that were deleted
- Consistent with other tool error patterns in codebase

### 6. Fresh Reads (No Caching)
**Decision**: Read skill files fresh on every `list_skills` and `get_skill` call. No in-memory or persistent caching.

**Rationale**:
- Simplifies implementation (no cache invalidation)
- Skills directory is not expected to change frequently during operation
- If skills are added/removed, operators see changes immediately
- Trade-off: minimal performance impact vs. added complexity

## Risks / Trade-offs

**[Risk: Large skills directories impact latency]** → `list_skills` scales linearly with number of skills. Mitigation: typical use case (<50 skills); if needed, can add pagination in future.

**[Risk: Symlinks or special files in skills directory]** → Could cause unexpected behavior. Mitigation: use `os.IsDir()` to filter only directories; skip anything that doesn't match naming rules.

**[Risk: Skill authors don't follow spec]** → Malformed skills skipped silently (with logs). Mitigation: clear documentation, warning messages in logs help identify issues.

**[Trade-off: No caching vs. directory polling]** → Each call reads disk. Acceptable for discovery use case (infrequent); if agents call `list_skills` every few seconds, could revisit.

## Migration Plan

No breaking changes; new tools are additive.

1. **Add yaml.v3 dependency**: `go get gopkg.in/yaml.v3`
2. **Add models** (request/response types) to `pkg/models/models.go`
3. **Implement handlers** in `pkg/tools/skills.go`
4. **Register tools** in `internal/server.go`
5. **Update tool count** log message in `NewServer()` (from 20 to 22 tools)
6. **Test**: Verify parsing, validation, error handling
7. **Deploy**: No rollback needed; purely additive

## Open Questions

None. All decisions locked in per coworker's earlier input.
