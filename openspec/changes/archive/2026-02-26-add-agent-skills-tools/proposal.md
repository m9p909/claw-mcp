## Why

Agents need a way to discover and access **reusable skill packages** that contain instructions, scripts, and references. Currently, agents must be manually told about available capabilities. Agent Skills is an open standard format for packaging and discovering these resources. Exposing `list_skills` and `get_skill` tools enables agents to autonomously explore available skills and invoke their scripts via bash.

## What Changes

- Add two new MCP tools: `list_skills` and `get_skill`
- `list_skills`: Scans `~/.mcpclaw/skills/` directory, validates YAML frontmatter per Agent Skills spec, returns metadata for all valid skills
- `get_skill`: Retrieves a specific skill's full content (metadata + markdown body), returns skill directory path so agents can bash into scripts/
- Invalid/malformed skills are skipped with warning logs, not exposed to agents
- Agents can use returned `skill_directory` path with `read_file` and `exec_command` tools to explore and run skill scripts

## Capabilities

### New Capabilities
- `agent-skills-discovery`: List and retrieve Agent Skills format packages with strict YAML validation per spec
- `agent-skills-execution`: Support for agents to invoke skill scripts via bash using returned directory paths

### Modified Capabilities

(none)

## Impact

- **New MCP tools**: `list_skills`, `get_skill` registered in internal/server.go
- **New code**: pkg/tools/skills.go with handler functions and YAML parsing
- **New models**: Request/response types in pkg/models/models.go
- **Dependencies**: gopkg.in/yaml.v3 (YAML parsing)
- **User-facing**: Agents can now discover and use skills from `~/.mcpclaw/skills/` directory
