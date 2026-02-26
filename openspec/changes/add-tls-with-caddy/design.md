## Context

Claw currently runs plain HTTP and relies on users to add TLS at the deployment layer. The project already documents this (see README.md): "Use TLS/HTTPS in production (via reverse proxy or load balancer)". We're implementing this pattern with Caddy as the reverse proxy.

Caddy was chosen because:
- Auto-renews Let's Encrypt certificates (zero maintenance)
- Simple, readable Caddyfile configuration
- Works identically in Docker Compose and on standalone VMs
- Decouples certificate lifecycle from application lifecycle

## Goals / Non-Goals

**Goals:**
- Provide TLS termination for Docker Compose deployments
- Support both development (self-signed, localhost) and production (Let's Encrypt) domains
- Enable standalone VM users to set up TLS with minimal steps
- Maintain zero code changes to Claw application
- Create clear, repeatable deployment patterns

**Non-Goals:**
- Kubernetes TLS (users manage via Ingress)
- Built-in application-level HTTPS (reverse proxy pattern is cleaner)
- ACME certificate management inside the app
- Rate limiting, WAF, or other advanced proxy features

## Decisions

### Decision 1: Reverse Proxy Pattern (Caddy)
**Choice**: Caddy as sidecar reverse proxy
**Rationale**:
- Clean separation: app handles business logic, proxy handles TLS/certs
- Automatic renewal via ACME (Let's Encrypt)
- Single Caddyfile template works for all deployment patterns
- Well-tested, battle-hardened for production use

**Alternatives Considered**:
- Built-in TLS in Go app: Adds cert management complexity, renewal logic, restart requirements
- Nginx: More powerful but requires manual cert renewal scripting
- Traefik: Overkill for singleton application pattern

### Decision 2: Two Docker Compose Files
**Choice**: Separate `docker-compose.yml` (development) and `docker-compose.prod.yml` (production)
**Rationale**:
- Development: Uses self-signed certs via `tls.dns.provider = "" ` (Caddy manual HTTPS)
- Production: Uses Let's Encrypt with domain binding
- Keeps configuration explicit and separate
- Users clearly choose dev vs prod

**Alternatives Considered**:
- Single file with env var toggle: Harder to read, mixing concerns
- Runtime selection: Extra complexity in entrypoint script

### Decision 3: Caddyfile Template with Environment Variables
**Choice**: Single Caddyfile template mounted in both configurations, using `{$DOMAIN}` substitution
**Rationale**:
- DRY principle: One Caddyfile template instead of duplicating config
- Docker Compose and standalone VMs can use the same template
- Environment variables `DOMAIN`, `TLS_MODE` control behavior
- Easier to maintain one canonical config

**Alternatives Considered**:
- Separate Caddyfiles: Duplication, harder to maintain
- Hardcoded domains: Inflexible, requires code changes per deployment

### Decision 4: Claw Container Port Binding
**Choice**: Claw internal port (8080) not exposed to host; only Caddy exposes 80/443
**Rationale**:
- Enforces TLS at network boundary
- Prevents accidental unencrypted access to Claw
- Cleaner architecture: users only interact with Caddy on public ports

**Alternatives Considered**:
- Expose both Claw and Caddy ports: Allows bypassing TLS, defeats security goal

### Decision 5: Standalone VM Setup
**Choice**: Simple guide: install Caddy, use provided Caddyfile template, mount cert storage
**Rationale**:
- Minimal steps (install → configure → run)
- Same Caddyfile template as Docker Compose (consistency)
- Users handle systemd/init integration per their preference

## Risks / Trade-offs

**[Risk] Caddy version incompatibility**
- Mitigation: Pin Caddy version in Dockerfile; test against latest stable release quarterly

**[Risk] Let's Encrypt rate limits in production**
- Mitigation: Document rate limits in guide; recommend staging environment testing before production

**[Risk] Docker Compose users can't quickly toggle between dev/prod**
- Mitigation: Acceptable trade-off for clarity; users explicitly choose which compose file to use

**[Risk] Standalone VM users need to manually restart Caddy on cert renewal**
- Mitigation: Caddy's auto-reload handles most cases; document manual restart procedure if needed

**[Trade-off] Network latency: Two containers instead of one**
- Acceptable: Network overhead negligible compared to Claw operations; orchestration benefit worth it

## Migration Plan

1. Create Caddyfile template in repo root
2. Update docker-compose.yml (dev) to include caddy sidecar
3. Create docker-compose.prod.yml for production
4. Create TLS_SETUP.md with:
   - Docker Compose section (dev & prod)
   - Standalone VM section
   - Troubleshooting
5. Update README.md to link to TLS_SETUP.md
6. No breaking changes; existing HTTP deployments continue to work

## Open Questions

None at this time. Design is stable for implementation.
