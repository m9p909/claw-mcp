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
  -e CLAW_TOKEN="your-secret-token" \
  -v ~/.mcpclaw:/root/.mcpclaw \
  claw:latest
```

**Docker Compose (Development with TLS):**

```bash
export CLAW_TOKEN="your-secret-token"
docker-compose up
```

This starts Claw with Caddy reverse proxy providing TLS encryption on localhost. For production with Let's Encrypt, see [TLS Setup Guide](TLS_SETUP.md).

**Docker Compose (Production with Let's Encrypt):**

```bash
export DOMAIN=claw.example.com
export CLAW_TOKEN=$(openssl rand -base64 32)
docker-compose -f docker-compose.prod.yml up
```

For detailed TLS setup instructions, see [TLS Setup Guide](TLS_SETUP.md).

## Kubernetes Deployment

For production deployments in Kubernetes, create the authentication secret first, then deploy:

```bash
# Create secret with your token
kubectl create secret generic claw-token --from-literal=token="your-secret-token"

# Or use the example manifest as a template
kubectl apply -f kubernetes/secret-example.yaml

# Deploy Claw
kubectl apply -f kubernetes/statefulset.yaml
kubectl apply -f kubernetes/service.yaml
```

This deploys Claw as a StatefulSet with:
- Single replica (singleton pattern)
- Persistent 10Gi storage for workspace and database
- Bearer token authentication via Kubernetes Secret
- Health checks and readiness probes
- Service exposing port 8080
- TLS/HTTPS via Ingress controller

See [kubernetes/README.md](kubernetes/README.md) for detailed instructions, including TLS configuration with cert-manager and Let's Encrypt.

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
- `CLAW_TOKEN` - **Required**. Bearer token for authenticating MCP endpoint requests

## Protocol

Claw uses the **streamable-http** protocol from the MCP standard. The `/mcp` endpoint handles all MCP requests using HTTP with request/response streaming.

## Architecture

- **Single Process**: One Claw server handles all agent connections
- **Shared Workspace**: All agents access the same `~/.mcpclaw/workspace`
- **Persistent Memory**: SQLite database ensures memory survives process restarts
- **Stateless HTTP**: Uses streamable-http for protocol compatibility

## Authentication

Claw requires Bearer token authentication on the `/mcp` endpoint. The `/health` endpoint remains unauthenticated for container orchestration probes.

### Configuration

Set the `CLAW_TOKEN` environment variable before starting the server. The server will fail immediately if this variable is not set.

```bash
export CLAW_TOKEN="your-secret-token"
./mcpclaw
```

### Using the API

Include the token in the `Authorization` header:

```bash
curl -H "Authorization: Bearer your-secret-token" http://localhost:8080/mcp
```

### Health Check (No Auth Required)

```bash
curl http://localhost:8080/health
```

## Security

Claw implements Bearer token authentication to protect the `/mcp` endpoint from unauthorized access. Additional considerations:

- **Token Storage**: Use secret management systems (Kubernetes Secrets, Docker secrets, environment managers) rather than hardcoding tokens
- **Token Rotation**: Restart the server with a new `CLAW_TOKEN` to rotate credentials
- **No Resource Limits**: Claw does not implement rate limiting or resource quotas
- **Encryption in Transit**: Docker Compose deployments include Caddy for TLS/HTTPS. See [TLS Setup Guide](TLS_SETUP.md) for production configuration with Let's Encrypt. For Kubernetes, use your Ingress controller.
- **No Sandboxing**: Commands executed via `exec_command` have full permissions of the Claw process

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
