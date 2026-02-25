## 1. Protocol Transport Fix

- [x] 1.1 Update `internal/server.go` to use `mcp.NewStreamableHTTPHandler` instead of SSE
- [x] 1.2 Verify MCP protocol initialization accepts streamable-http transport
- [x] 1.3 Test `/mcp` endpoint with curl to confirm streamable-http responses
- [x] 1.4 Document protocol change in CHANGELOG or release notes

## 2. Dockerfile Creation

- [x] 2.1 Create `Dockerfile` with Ubuntu 22.04 LTS base image
- [x] 2.2 Install Node.js and npm in Dockerfile
- [x] 2.3 Install Python 3 and pip in Dockerfile
- [x] 2.4 Install batteries-included Unix tools (git, curl, wget, grep, sed, awk, find)
- [x] 2.5 Create build stage that compiles Go binary from source
- [x] 2.6 Create runtime stage that copies compiled binary and dependencies
- [x] 2.7 Verify multi-stage build produces runtime image < 1GB

## 3. Entrypoint Script

- [x] 3.1 Create `scripts/entrypoint.sh` to initialize server
- [x] 3.2 Add logic to handle PORT environment variable override
- [x] 3.3 Add logic to create workspace and database directories
- [x] 3.4 Make entrypoint script executable and set as Docker ENTRYPOINT

## 4. Docker Build & Test

- [x] 4.1 Build Docker image locally with `docker build -t claw:latest .`
- [x] 4.2 Run container with port mapping: `docker run -p 8080:8080 claw:latest`
- [x] 4.3 Test health endpoint: `curl http://localhost:8080/health`
- [x] 4.4 Test MCP initialization endpoint with streamable-http protocol
- [x] 4.5 Verify Node.js available in container: `docker run claw:latest node --version`
- [x] 4.6 Verify Python available in container: `docker run claw:latest python3 --version`

## 5. Docker Compose Setup

- [x] 5.1 Create `docker-compose.yml` with Claw service
- [x] 5.2 Add persistent volume mount at `~/.mcpclaw` (contains workspace and database subdirectories)
- [x] 5.3 Add port mapping for 8080:8080
- [x] 5.4 Test Compose setup: `docker-compose up` and verify health endpoint
- [x] 5.5 Test volume persistence: write file to workspace, restart, verify file exists

## 6. Kubernetes Manifests

- [x] 6.1 Create `kubernetes/statefulset.yaml` defining StatefulSet with 1 replica
- [x] 6.2 Add single PersistentVolumeClaim for `~/.mcpclaw` (e.g., 10Gi for workspace + database)
- [x] 6.3 Mount PVC at `~/.mcpclaw` in container, ensuring `workspace` and `data` subdirectories are created
- [x] 6.4 Create `kubernetes/service.yaml` exposing port 8080
- [x] 6.5 Add environment variables or ConfigMap for port configuration
- [x] 6.6 Create `kubernetes/README.md` with deployment instructions

## 7. Documentation

- [x] 7.1 Update main README with Docker and Kubernetes sections
- [x] 7.2 Document how to build and run Docker image locally
- [x] 7.3 Document Docker Compose usage and volume setup
- [x] 7.4 Document Kubernetes deployment with persistent volumes
- [x] 7.5 Update API documentation to reflect streamable-http protocol
- [x] 7.6 Create migration guide for clients using old SSE transport

## 8. Testing & Validation

- [x] 8.1 Verify all 8 tools work in Dockerized environment
- [x] 8.2 Test multi-agent scenario: run two agents connecting simultaneously
- [x] 8.3 Verify workspace files shared between agents
- [x] 8.4 Verify memory data persists across container restarts
- [x] 8.5 Performance test: measure response time from containerized server
- [x] 8.6 Security review: document intentional insecurities and mitigation strategy
