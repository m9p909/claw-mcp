## Why


THe purpose of this project is to easily turn normal AI agent software into claw-like AI software.

The MCP server enables Claude and other AI models to interact with your filesystem, execute commands, and maintain persistent memory. This is the foundation for creating AI-powered development tools that can read code, run tests, and remember context across conversations.


## What Changes

- Implement HTTP-based MCP server in Go exposing 9 tools across 3 categories
- Add filesystem operations with hashline-based editing (eliminates whitespace brittleness)
- Add execution engine for spawning and managing background processes
- Add SQLite-backed persistent memory system with in-memory search
- Server listens on localhost:3000, SQLite database at ~/.mcpclaw/data

## Capabilities

### New Capabilities
- `filesystem-tools`: Read, write, and edit files with content hashing for precise edits
- `process-execution`: Execute commands synchronously or in background, poll and manage running processes
- `memory-persistence`: Store and query memories in SQLite, search with substring matching

### Modified Capabilities
<!-- No existing capabilities modified -->

## Impact

- New module dependencies: MCP Go SDK, SQLite driver
- New persistent storage: ~/.mcpclaw/data directory
- Exposes local filesystem and command execution via HTTP (security implications for deployment)