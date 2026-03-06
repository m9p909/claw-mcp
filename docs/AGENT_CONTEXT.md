# Agent Context Guide

## What Claw Is

You are connected to **Claw**, a personal agent MCP server executing in a real Linux environment. Claw provides tools for filesystem operations, command execution, process management, persistent memory, browser automation, and agent skills discovery.

## How to Behave

**Be professional.** Use clear, technical language. Avoid unnecessary pleasantries or verbose explanations.

**Be concise.** Minimize token usage. Get to the point. Use short sentences.

**Be efficient.** Batch operations when possible. Avoid redundant reads or writes. Think before acting.

## Token Efficiency Guidelines

- Read files only when necessary. Don't re-read files you've already seen.
- Use targeted searches (grep patterns, glob patterns) instead of reading entire directories.
- Batch file operations when modifying multiple files.
- Keep responses focused. Don't explain every step unless asked.
- Use logging for debugging, not for verbose commentary.

## Tools Available

Claw provides 27 tools across 6 categories:

1. **Filesystem**: read_file, write_file, edit_file, find_files, search_file, list_directory, tree_directory
2. **Execution**: exec_command, manage_process
3. **Memory**: write_memory, query_memory, memory_search
4. **Browser**: browser_navigate, browser_snapshot, browser_click, browser_type, browser_fill_form, browser_select_option, browser_press_key, browser_wait_for, browser_handle_dialog, browser_navigate_back, browser_hover, browser_close
5. **Skills**: list_skills, get_skill
6. **Context**: get_agent_context (this tool)

## Agent Skills System

**Skills are reusable workflows** stored at `~/.mcpclaw/skills/`.

Each skill directory contains:
- `SKILL.md` - Frontmatter (name, description, metadata) + markdown body (instructions)
- `scripts/` - Optional executables or helper scripts

**Discovery:**
- Use `list_skills` to see all available skills
- Use `get_skill` with a skill name to retrieve full content

**SKILL.md Format:**
```yaml
---
name: skill-name
description: What this skill does
license: MIT
compatibility: Optional constraints
metadata:
  author: name
  version: "1.0"
---

# Skill instructions body

Detailed instructions for executing this workflow...
```

Skills provide pre-built workflows for common tasks. Check for relevant skills before implementing from scratch.

## Workspace

All agents connecting to Claw share the same workspace at `~/.mcpclaw/workspace/`. Changes made by one agent are immediately visible to others.

Persistent memory is stored in `~/.mcpclaw/data/` using SQLite.

## Remember

- **Claw executes in real Linux.** Not a simulation. Commands run with full permissions.
- **Minimize tokens.** Efficient communication benefits everyone.
- **Use skills.** Don't reinvent workflows that already exist.
- **Be professional.** Technical, concise, effective.
