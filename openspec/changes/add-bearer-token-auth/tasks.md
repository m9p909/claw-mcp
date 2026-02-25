## 1. Core Implementation

- [x] 1.1 Create auth middleware function in `main.go`
- [x] 1.2 Read `CLAW_TOKEN` environment variable at startup and fail fast if missing
- [x] 1.3 Extract and validate Bearer token from Authorization header
- [x] 1.4 Return 401 with JSON error body on auth failure

## 2. Server Integration

- [x] 2.1 Wrap the mux with auth middleware
- [x] 2.2 Verify `/mcp` endpoint requires valid token
- [x] 2.3 Verify `/health` endpoint bypasses auth

## 3. Docker & Kubernetes

- [x] 3.1 Update Dockerfile to include CLAW_TOKEN in build documentation
- [x] 3.2 Update docker-compose.yml to pass CLAW_TOKEN environment variable
- [x] 3.3 Update kubernetes/statefulset.yaml to read CLAW_TOKEN from Secret
- [x] 3.4 Create example kubernetes secret manifest for documentation

## 4. Testing & Verification

- [x] 4.1 Verify server fails to start without CLAW_TOKEN set
- [x] 4.2 Test valid token provides access to `/mcp`
- [x] 4.3 Test invalid token returns 401 on `/mcp`
- [x] 4.4 Test missing header returns 401 on `/mcp`
- [x] 4.5 Test `/health` works without authentication
- [x] 4.6 Test Docker image with CLAW_TOKEN env var
- [x] 4.7 Test Kubernetes deployment with secret-based token

## 5. Documentation

- [x] 5.1 Update README.md with authentication section
- [x] 5.2 Document CLAW_TOKEN environment variable requirement
- [x] 5.3 Update deployment instructions for Docker and Kubernetes
