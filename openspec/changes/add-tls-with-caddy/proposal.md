## Why

The Claw MCP Server currently runs over plain HTTP. For production deployments, users need TLS/HTTPS encryption in transit. Rather than adding complexity to the application itself, we'll use Caddy as a reverse proxy sidecar. This provides automated certificate management via Let's Encrypt, clean separation of concerns, and easy deployment across Docker Compose and standalone VMs.

## What Changes

- Add Caddy reverse proxy as sidecar in Docker Compose setup
- Create templated Caddyfile that supports both development and production domains
- Add two docker-compose configurations: one for development (self-signed) and one for production (Let's Encrypt)
- Create comprehensive TLS setup guide covering Docker Compose and standalone VM deployments
- Update main README to reference TLS guide for production users
- Kubernetes users continue using their existing Ingress configuration

## Capabilities

### New Capabilities
- `tls-with-caddy-docker-compose`: TLS termination via Caddy in Docker Compose environments with automatic Let's Encrypt certificate management
- `tls-setup-guide`: Documentation for setting up TLS on standalone VMs using Caddy and production Docker Compose

### Modified Capabilities

## Impact

- Docker Compose: Adds caddy service, changes internal port bindings (claw no longer exposes 8080 directly)
- Documentation: New TLS_SETUP.md, updated README.md
- No code changes to Claw application itself
- Kubernetes deployments unaffected (users manage TLS via Ingress)
