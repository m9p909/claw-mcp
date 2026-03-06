## Context

Claw MCP server currently provides 26 tools across filesystem, execution, memory, browser, skills, and search categories. Agents connect via MCP but lack onboarding context about what Claw is, how to behave efficiently, or how to utilize features like Agent Skills stored at `~/.mcpclaw/skills/`.

The existing tool pattern (skills tools: `list_skills`, `get_skill`) provides a clear model: tool descriptions offer quick context, tool invocation returns full content.

## Goals / Non-Goals

**Goals:**
- Provide agents with essential context immediately visible in tool listings
- Guide agents toward professional, token-efficient behavior
- Explain Claw's identity as a personal Linux agent and the skills system
- Use existing patterns (follow skills tools model)
- Keep implementation minimal (single tool, single doc file)

**Non-Goals:**
- Auto-injection mechanisms beyond tool description visibility
- Comprehensive tool reference (MCP handles tool discovery natively)
- Dynamic documentation generation
- Multi-file documentation system

## Decisions

### Decision 1: Tool description contains essential context
**Choice:** Pack core guidance into tool description field (visible in tool listings)
**Rationale:** Agents see tool list on connect. Description provides immediate context without requiring tool invocation.
**Alternatives considered:**
- Server metadata extension: Not standardized in MCP, less discoverable
- Well-known skill: Requires agents to know to list/get skills first
- Separate initialization endpoint: Adds complexity, breaks MCP patterns

### Decision 2: Markdown file stores full content
**Choice:** Create `docs/AGENT_CONTEXT.md` with detailed imperative guidance
**Rationale:** Separates content from code, allows easy updates, follows docs convention
**Alternatives considered:**
- Embed in Go code: Harder to maintain, requires rebuilds for changes
- Multiple files: Overkill for current scope, adds complexity

### Decision 3: Follow skills tools pattern
**Choice:** Handler in `pkg/tools/`, simple read from `docs/`, return as text
**Rationale:** Consistent with existing `get_skill` implementation, minimal code
**Alternatives considered:**
- Streaming response: Not needed for static content
- Structured JSON: Markdown more readable for agents

### Decision 4: Tool description format
**Choice:** Imperative tone covering: identity, personality, token efficiency, skills basics
**Rationale:** Per requirements, direct instructions most effective
**Format:**
```
You are Claw, a personal agent executing in a real Linux environment.
Be professional and concise. Optimize for token efficiency.
Agent Skills at ~/.mcpclaw/skills/ provide reusable workflows (use list_skills/get_skill).
Call this tool for full context and usage guide.
```

## Risks / Trade-offs

**[Risk]** Tool description length limits → **Mitigation:** Keep under 500 chars, prioritize essentials
**[Risk]** Agents may not notice/call tool → **Mitigation:** Description itself provides core guidance
**[Trade-off]** Static content vs dynamic → Accept: Simple markdown easier to maintain than generation logic
**[Trade-off]** Single file vs modular → Accept: Current scope doesn't justify complexity
