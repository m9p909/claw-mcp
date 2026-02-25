# Claw MCP Server

A Model Context Protocol (MCP) server written in Go that provides tools for file operations, command execution, and persistent memory management.

## Overview

Claw is a singleton MCP server designed for collaboration between multiple AI agents (Claude, GPT-4, etc.). It provides:

- **File Operations**: Read, write, and edit files with hash-based validation
- **Command Execution**: Execute commands synchronously or asynchronously
- **Process Management**: Manage background processes and streams
- **Memory Persistence**: Store and query persistent agent memory via SQLite

All agents connect to the same shared workspace at `~/.mcpclaw/workspace`, enabling real-time collaboration.

## Building from Source

### Prerequisites

- Go 1.25+
- SQLite development headers

### Build

```bash
go build -o mcpclaw .
```

### Run

```bash
./mcpclaw -port 8080
```

The server will start on `http://localhost:8080` with MCP endpoint at `/mcp`.

## Docker Deployment

### Building the Docker Image

```bash
docker build -t claw:latest .
```

This creates a multi-stage Docker image based on Ubuntu 22.04 with:
- Go 1.25 (builder stage)
- Node.js and npm (runtime)
- Python 3 and pip (runtime)
- Batteries-included Unix tools (git, curl, wget, gawk, etc.)

### Running with Docker

**Direct Docker:**

```bash
docker run -p 8080:8080 \
  -v ~/.mcpclaw:/root/.mcpclaw \
  claw:latest
```

**Docker Compose:**

```bash
docker-compose up
```

This starts the Claw server with persistent volume at `~/.mcpclaw`.

## Kubernetes Deployment

For production deployments in Kubernetes, use the provided manifests:

```bash
kubectl apply -f kubernetes/statefulset.yaml
kubectl apply -f kubernetes/service.yaml
```

This deploys Claw as a StatefulSet with:
- Single replica (singleton pattern)
- Persistent 10Gi storage for workspace and database
- Health checks and readiness probes
- Service exposing port 8080

See [kubernetes/README.md](kubernetes/README.md) for detailed instructions.

## API Documentation

### Health Check

```bash
curl http://localhost:8080/health
```

Response:
```json
{"status": "ok"}
```

### MCP Endpoint

The MCP endpoint is available at `/mcp` using the streamable-http protocol.

Connect your MCP client to `http://localhost:8080/mcp`.

## Tools

### Filesystem Tools

#### read_file
Read the contents of a file with line hashes.

#### write_file
Write or create a file with content.

#### edit_file
Edit a file by replacing a range of lines identified by hashes.

### Execution Tools

#### exec_command
Execute a command synchronously or asynchronously.

Parameters:
- `command` (string): Command to execute
- `args` (array): Command arguments
- `background` (boolean): Run in background
- `env` (object): Environment variables

#### manage_process
Manage background processes.

Actions:
- `list` - List all sessions
- `poll` - Get session status
- `send_keys` - Send input to process
- `kill` - Terminate process

### Memory Tools

#### write_memory
Store persistent memory.

Parameters:
- `category`: fact, todo, decision, or preference
- `content`: Memory content

#### query_memory
Query memory using SQL SELECT.

Parameters:
- `query`: SQL SELECT query

#### memory_search
Search memory with substring matching.

Parameters:
- `query`: Substring to search
- `limit`: Max results (0 = unlimited)

## Persistent Storage

Data is stored at `~/.mcpclaw/`:

```
~/.mcpclaw/
├── workspace/     # Shared workspace for all agents
└── data/          # SQLite database for memory
```

Both directories are created automatically on server startup.

## Environment Variables

- `PORT` - Server port (default: 8080)

## Protocol

Claw uses the **streamable-http** protocol from the MCP standard. The `/mcp` endpoint handles all MCP requests using HTTP with request/response streaming.

## Architecture

- **Single Process**: One Claw server handles all agent connections
- **Shared Workspace**: All agents access the same `~/.mcpclaw/workspace`
- **Persistent Memory**: SQLite database ensures memory survives process restarts
- **Stateless HTTP**: Uses streamable-http for protocol compatibility

## Security

**Note**: This server is intentionally insecure by design. It provides:
- No authentication or authorization
- No resource limits or sandboxing
- No encryption in transit

**Deployment recommendation**: Place behind an API Gateway or authentication proxy for production use.

## Development

### Project Structure

```
.
├── main.go                 # Server entrypoint
├── internal/
│   ├── server.go          # MCP server setup and tool registration
│   └── db.go              # SQLite database management
├── pkg/
│   ├── tools/             # Tool implementations
│   │   ├── filesystem.go
│   │   ├── execution.go
│   │   ├── memory.go
│   │   └── helpers.go
│   ├── models/            # Data structures and schemas
│   ├── storage/           # Persistence layer
│   └── hash/              # Hash utilities for file editing
├── Dockerfile             # Multi-stage Docker build
├── docker-compose.yml     # Local development compose
└── kubernetes/            # K8s manifests and docs
```

### Running Tests

Tests are defined in the spec documents at `openspec/changes/add-docker-setup/`.

## Migration Guide

### From SSE to Streamable-HTTP

Previous versions may have used SSE (Server-Sent Events) for MCP transport. This version uses the MCP-standard **streamable-http** protocol.

**For clients**: Update your MCP client to use the streamable-http transport when connecting to `/mcp`.

## Contributing

This is part of the Claw MCP Server project. For changes, see the OpenSpec change specification at `openspec/changes/add-docker-setup/`.

## License

TBD
