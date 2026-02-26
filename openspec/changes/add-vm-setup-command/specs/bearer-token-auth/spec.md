## MODIFIED Requirements

### Requirement: Bearer token authentication on /mcp endpoint
The application SHALL authenticate requests to the `/mcp` endpoint using HTTP Bearer token authentication when `CLAW_TOKEN` environment variable is set. When `CLAW_TOKEN` is empty or unset, Bearer token validation SHALL be skipped.

#### Scenario: CLAW_TOKEN is set
- **WHEN** `CLAW_TOKEN` environment variable is set to a non-empty value
- **THEN** requests to `/mcp` without valid Bearer token in `Authorization: Bearer <token>` header are rejected with 401 Unauthorized

#### Scenario: CLAW_TOKEN is unset
- **WHEN** `CLAW_TOKEN` environment variable is unset or empty string
- **THEN** requests to `/mcp` endpoint do not require Bearer token; all requests are accepted

#### Scenario: Caddy basic auth is preferred over CLAW_TOKEN
- **WHEN** Caddy reverse proxy is configured with basic auth directive
- **THEN** `/mcp` endpoint authentication is delegated to Caddy; CLAW_TOKEN can be left empty

#### Scenario: /health endpoint remains unauthenticated
- **WHEN** request is made to `/health` endpoint
- **THEN** response is 200 OK without authentication requirement (both CLAW_TOKEN set or unset)
