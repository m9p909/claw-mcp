## Context

The Claw MCP Server is currently a compiled binary that requires manual setup. We're adding containerization to support production Kubernetes deployments where multiple MCP agents (Claude, GPT-4, etc.) connect to a shared singleton server instance with a common workspace. The server uses streamable-http (MCP standard transport) and SQLite for memory persistence.

## Goals / Non-Goals

**Goals:**
- Create a production-ready Docker image with Ubuntu base, Node.js, Python, and batteries-included tools
- Enable Kubernetes deployment with persistent volumes for workspace and database
- Support Docker Compose for local development
- Fix transport protocol to use streamable-http from MCP standard
- Allow multiple MCP agents to connect simultaneously to a shared workspace singleton

**Non-Goals:**
- Security hardening or sandboxing (intentionally insecure by design)
- Image optimization for minimal size
- Custom authentication (delegated to API Gateway)
- Multi-node clustering or horizontal scaling

## Decisions

**1. Base Image: Ubuntu**
- **Choice**: Use Ubuntu 22.04 LTS as base
- **Rationale**: "Batteries included" approach - Ubuntu provides core Unix tools (git, curl, wget, grep, sed, awk, find) natively, reducing Dockerfile complexity and image customization
- **Alternatives**: Alpine (smaller, but missing many standard tools; would require manual installation), Debian (more minimal than Ubuntu but still rich toolset)

**2. Pre-installed Languages**
- **Choice**: Include both Node.js and Python in the base image
- **Rationale**: Agents need ability to execute Node and Python scripts in workspace; including languages avoids runtime installation delays
- **Alternatives**: Install-on-demand (slower for agents), separate base images per language (defeats singleton workspace goal)

**3. Workspace Architecture**
- **Choice**: Single persistent `~/.mcpclaw/workspace` directory mounted at runtime, shared by all connected agents
- **Rationale**: Simplifies multi-agent collaboration; agents see same files and directory state; keeps all Claw data under `~/.mcpclaw`
- **Alternatives**: Per-agent workspaces (prevents collaboration, complicates data sharing), separate `/workspace` directory (scatters Claw data)

**4. Database Persistence**
- **Choice**: SQLite database stored at `~/.mcpclaw/data` within the same persistent volume as workspace
- **Rationale**: Single persistent volume for all Claw state simplifies storage configuration; both workspace and database survive pod restarts
- **Alternatives**: Separate volumes (adds complexity), in-memory database (data lost on restart)

**5. Entrypoint Script**
- **Choice**: Create `scripts/entrypoint.sh` to handle initialization and port configuration
- **Rationale**: Cleanly separates container setup from binary execution; allows environment variable configuration
- **Alternatives**: Hardcode in Dockerfile (less flexible), rely on docker run flags (more fragile)

**6. Transport Protocol Fix**
- **Choice**: Update `internal/server.go` to use mcp.NewStreamableHTTPHandler instead of SSE
- **Rationale**: Streamable-http is the MCP standard transport; SSE was an incorrect assumption
- **Alternatives**: Keep SSE (violates MCP spec compatibility)

**7. Kubernetes Deployment Model**
- **Choice**: Single-instance StatefulSet with two PersistentVolumeClaims
- **Rationale**: Singleton ensures agents access shared workspace; StatefulSet provides stable pod identity and volume binding
- **Alternatives**: Deployment (simpler but less suitable for stateful workloads), DaemonSet (scales to every node, wastes resources)

**8. Compile Strategy**
- **Choice**: Multi-stage Docker build: compile Go binary in builder stage, copy to runtime stage
- **Rationale**: Runtime image contains only binary and dependencies, not Go SDK, reducing image size and attack surface
- **Alternatives**: Pre-compiled binary (requires separate build pipeline), compile at runtime (slower builds)

## Risks / Trade-offs

- **Large image size** → Ubuntu base is ~77MB; adding Node/Python makes image 800MB+. Mitigation: Use multi-stage build, document size expectations
- **Transport protocol change breaking compatibility** → Existing clients using SSE may fail. Mitigation: Not needed 
- **Shared workspace concurrency** → Multiple agents writing simultaneously could cause conflicts. Mitigation: Not Needed 
- **Single volume for workspace and database** → Requires careful directory structure. Mitigation: Ensure `~/.mcpclaw` contains both subdirectories with clear separation
- **Intentional insecurity** → No resource limits, sandboxing, or authentication. Mitigation: Document that this requires API Gateway protection
