## Context

Claw is an MCP server that currently accepts all requests without authentication. The server uses HTTP with a multiplexed handler for `/mcp` and `/health` endpoints. We need to add Bearer token authentication to the `/mcp` endpoint only, using a token provided via the `CLAW_TOKEN` environment variable.

## Goals / Non-Goals

**Goals:**
- Require valid Bearer token on all `/mcp` requests
- Fail at startup if `CLAW_TOKEN` is not set (no insecure mode)
- Return 401 with JSON error body on authentication failure
- Keep `/health` endpoint unauthenticated (for container orchestration probes)
- Minimal code changes using Go's standard HTTP middleware pattern

**Non-Goals:**
- Multi-user or role-based access control
- Token rotation or expiration mechanisms
- Integration with external identity providers
- Per-tool or per-agent authorization

## Decisions

**Decision 1: HTTP Middleware vs Handler Wrapper**
- **Chosen**: HTTP Middleware (wrap the mux with a middleware function)
- **Rationale**: Cleaner separation of concerns; allows us to check auth before routing. Alternative (wrapping just the MCP handler) would require deeper changes to the streamable-http integration.

**Decision 2: Token Storage Location**
- **Chosen**: Environment variable `CLAW_TOKEN` read at startup
- **Rationale**: Simple, no persistence layer needed. Matches existing pattern (PORT env var). Alternative (database) adds complexity for no benefit with a single global token.

**Decision 3: Startup Behavior**
- **Chosen**: Fail fast if `CLAW_TOKEN` is missing
- **Rationale**: Prevents accidental insecure deployments. Forces explicit token provisioning before running.

**Decision 4: Error Response Format**
- **Chosen**: JSON `{"error": "Unauthorized"}` with 401 status
- **Rationale**: Consistent with `/health` endpoint style. Provides minimal info (don't leak token existence).

**Decision 5: Route Exclusion**
- **Chosen**: Only `/mcp` requires auth; `/health` is open
- **Rationale**: Health checks from orchestration systems (Kubernetes, Docker) need unauthenticated access.

## Risks / Trade-offs

**Risk: Token in Environment Variable**
- **Mitigation**: Document that `CLAW_TOKEN` should be set via secret management tools (K8s Secrets, Docker secrets, etc.), not hardcoded in container images.

**Risk: Single Global Token**
- **Trade-off**: Cannot isolate agents or audit per-client access. Acceptable for now (per-agent identity is future work).

**Risk: Token Comparison Timing**
- **Mitigation**: Use constant-time comparison to prevent timing attacks (Go's `bytes.Equal` is constant-time, but we're comparing strings; use `subtle.ConstantTimeCompare` for best practice).

**Risk: No Token Rotation**
- **Trade-off**: Acceptable. Rotation can be handled by restarting server with new token.

## Deployment Impact

- Docker: Add `CLAW_TOKEN` to environment in Dockerfile or compose file
- Kubernetes: Use Secret and environment variable in StatefulSet
- Development: Must set `CLAW_TOKEN` before running (no dev bypass)
