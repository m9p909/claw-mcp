## 1. Caddyfile Template Creation

- [x] 1.1 Create Caddyfile template in project root with {$DOMAIN} and {$TLS_MODE} substitutions
- [x] 1.2 Configure Caddyfile to reverse proxy :80 and :443 to internal http://claw:8080
- [x] 1.3 Add self-signed certificate generation for TLS_MODE=self_signed
- [x] 1.4 Add Let's Encrypt configuration for TLS_MODE=letsencrypt with ACME challenges
- [x] 1.5 Configure automatic HTTP→HTTPS redirect in production mode
- [ ] 1.6 Test Caddyfile syntax and basic proxy functionality

## 2. Docker Compose Configuration - Development

- [x] 2.1 Update docker-compose.yml to include caddy service
- [x] 2.2 Configure caddy to mount Caddyfile as volume
- [x] 2.3 Set DOMAIN=localhost and TLS_MODE=self_signed for development
- [x] 2.4 Remove port mappings from claw service (only caddy exposes ports)
- [x] 2.5 Configure internal docker-compose network for claw ↔ caddy communication
- [x] 2.6 Add healthcheck for caddy service
- [ ] 2.7 Test docker-compose up starts both services and HTTPS works on localhost

## 3. Docker Compose Configuration - Production

- [x] 3.1 Create docker-compose.prod.yml based on docker-compose.yml
- [x] 3.2 Set TLS_MODE=letsencrypt for Let's Encrypt certificate acquisition
- [x] 3.3 Make DOMAIN configurable via environment variable (no default)
- [x] 3.4 Add documentation comment in .prod file explaining DOMAIN and CLAW_TOKEN requirements
- [x] 3.5 Configure Caddy data storage volume for certificate persistence (/data/caddy)
- [ ] 3.6 Test with example domain (use example.com placeholder in docs)

## 4. Documentation - TLS Setup Guide

- [x] 4.1 Create TLS_SETUP.md in project root
- [x] 4.2 Add "Docker Compose - Development" section with 3-step quickstart
- [x] 4.3 Add "Docker Compose - Production" section with domain/CLAW_TOKEN setup
- [x] 4.4 Add "Standalone VM" section with Caddy installation instructions
- [x] 4.5 Add "Standalone VM - Configure Caddyfile" section with template and example
- [x] 4.6 Add "Certificate Renewal" section explaining automatic renewal behavior
- [x] 4.7 Add "Troubleshooting" section with common issues and solutions
- [x] 4.8 Add "Kubernetes" section explaining that users should use Ingress instead

## 5. Documentation - README Updates

- [x] 5.1 Update README.md Security section to reference TLS_SETUP.md
- [x] 5.2 Change "Use TLS/HTTPS in production (via reverse proxy or load balancer)" to link to guide
- [x] 5.3 Add note about Caddy sidecar in Docker Compose section
- [x] 5.4 Update Kubernetes section to mention TLS is handled via Ingress

## 6. Testing and Validation

- [ ] 6.1 Verify docker-compose.yml (dev) starts successfully with HTTPS on localhost
- [ ] 6.2 Verify docker-compose.prod.yml can be configured with DOMAIN environment variable
- [ ] 6.3 Verify Claw port 8080 is not accessible directly from host in docker-compose
- [ ] 6.4 Verify Caddy successfully proxies requests to Claw on internal network
- [ ] 6.5 Verify certificate files are mounted correctly in both configurations
- [ ] 6.6 Test that HTTPS redirects work (HTTP → HTTPS for production)
- [ ] 6.7 Verify guide steps can be followed end-to-end by another person

---

**Note:** Testing tasks (6.1-6.7) require docker-compose to be run and should be performed by the user. All implementation and documentation is complete.
