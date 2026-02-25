## Why

The Claw MCP Server currently requires manual setup and compilation. To enable production Kubernetes deployments and multi-agent collaboration, we need a containerized image that packages the server with all dependencies, allowing teams to deploy Claw as a shared sandbox where Claude, GPT-4, and other MCP agents can collaborate in a common workspace at `~/.mcpclaw/workspace`.

## What Changes

- Create a Dockerfile based on Ubuntu with pre-installed Node.js, Python, and batteries-included Unix tools (git, curl, wget, grep, sed, awk, find, etc.)
- Build and publish a Docker image containing the compiled Claw MCP Server binary
- Create Kubernetes deployment manifests for running Claw as a singleton with persistent volume for `~/.mcpclaw` (containing both workspace and database)
- Add Docker Compose configuration for local multi-container development
- Create entrypoint script handling server initialization and volume mounts

## Capabilities

### New Capabilities

- `docker-image-build`: Build and containerize the Claw MCP Server with batteries-included tools
- `kubernetes-deployment`: Deploy Claw as a Kubernetes singleton with persistent storage
- `docker-compose-local`: Local development environment with Docker Compose

### Modified Capabilities

- `mcp-server-core`: Update transport to use streamable-http from MCP standard (fixing SSE assumption)

## Impact

- Affected code: `main.go`, `internal/server.go` (protocol transport update)
- New files: `Dockerfile`, `kubernetes/statefulset.yaml`, `kubernetes/service.yaml`, `docker-compose.yml`, `scripts/entrypoint.sh`
- Dependencies: Docker, Kubernetes (optional), Docker Compose (optional)
- Systems: Container registry (for image storage), Kubernetes cluster (for production)
- Storage: Single PersistentVolume for `~/.mcpclaw` containing workspace and database subdirectories
