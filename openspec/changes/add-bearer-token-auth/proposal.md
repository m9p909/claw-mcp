## Why

Claw is currently deployed without authentication, making it suitable only for trusted, isolated environments. This creates a security risk if the server is accidentally exposed or deployed to less-controlled environments. Adding Bearer token authentication makes Claw safe for production use without requiring external authentication proxies.

## What Changes

- Add `CLAW_TOKEN` environment variable for configuring the authentication token
- Require Bearer token in `Authorization` header on `/mcp` endpoint
- Fail fast at startup if `CLAW_TOKEN` is not set (required for all deployments)
- Return `401 Unauthorized` with error body on invalid/missing token for `/mcp` requests
- Leave `/health` endpoint unauthenticated for liveness/readiness probes
- Update Docker and Kubernetes manifests to include `CLAW_TOKEN` configuration

## Capabilities

### New Capabilities
- `mcp-bearer-token-auth`: Bearer token validation on MCP endpoint requests

### Modified Capabilities
<!-- No existing capability requirements are changing -->

## Impact

- **Code**: `main.go` (HTTP middleware), Docker/K8s configs
- **Breaking Change**: Clients must now provide valid `Authorization: Bearer <token>` header
- **Deployment**: Requires setting `CLAW_TOKEN` environment variable before startup
- **Operations**: No external auth dependency needed; token lifecycle is admin-managed
