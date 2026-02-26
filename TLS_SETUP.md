# TLS Setup Guide for Claw MCP Server

This guide covers setting up TLS/HTTPS for Claw in Docker Compose (development and production) and on standalone virtual machines.

For Kubernetes deployments, use your existing Ingress controller to handle TLS termination—see the [Kubernetes section](#kubernetes).

---

## Docker Compose - Development (Self-Signed TLS)

For local development with self-signed HTTPS certificates:

### Prerequisites
- Docker and Docker Compose installed
- Claw Docker image built: `docker build -t claw:latest .`

### Quick Start

```bash
# 1. Set a development token
export CLAW_TOKEN="dev-token-12345"

# 2. Start Claw with self-signed TLS on localhost
docker-compose up

# 3. Access via HTTPS
# Browser: https://localhost
# Accept the self-signed certificate warning (it's expected)
# Health check: curl -k https://localhost/health
```

### What Happens

- Caddy automatically generates a self-signed certificate for `localhost`
- Traffic flows: Browser (HTTPS) → Caddy (443) → Claw (8080, internal only)
- Claw's port 8080 is not exposed to the host—all external traffic goes through Caddy's TLS

### Accepting Self-Signed Certificates

**Browser:**
- Navigate to `https://localhost`
- Click "Advanced" or equivalent button
- Click "Proceed to localhost (unsafe)" or similar

**curl:**
```bash
curl -k https://localhost/health
# or
curl --insecure https://localhost/health
```

**MCP clients:**
- Configure client to accept self-signed certificates (usually a flag like `--insecure`)

---

## Docker Compose - Production (Let's Encrypt TLS)

For production deployment with automatic Let's Encrypt certificate renewal:

### Prerequisites

- Docker and Docker Compose installed
- A public domain (e.g., `claw.example.com`)
- DNS A record pointing to your server
- Claw Docker image built: `docker build -t claw:latest .`
- Secure token generated: `openssl rand -base64 32`

### Setup

```bash
# 1. Set your production domain
export DOMAIN=claw.example.com

# 2. Set a secure authentication token
export CLAW_TOKEN=$(openssl rand -base64 32)
# Save this token securely (use a secrets manager in production)

# 3. Verify DNS A record is set
nslookup claw.example.com
# Should resolve to your server's IP address

# 4. Start Claw with production TLS
docker-compose -f docker-compose.prod.yml up -d

# 5. Check certificate status
docker logs claw-caddy-prod | grep "certificate"
```

### What Happens

- Caddy requests a certificate from Let's Encrypt using ACME
- ACME validation uses HTTP-01 challenge on port 80
- Certificate is issued and stored in the `caddy_data` volume
- Caddy automatically renews the certificate 30 days before expiration
- Traffic flows: Browser (HTTPS) → Caddy (443, Let's Encrypt) → Claw (8080, internal only)

### Verification

```bash
# Check certificate details
curl -vI https://claw.example.com

# Test MCP endpoint
curl -H "Authorization: Bearer $CLAW_TOKEN" https://claw.example.com/health
```

### Rate Limits

Let's Encrypt has rate limits:
- **50 certificates per domain per week**
- **5 certificates per week per IP address** (for new domains)

If you hit these limits:
- Wait before retrying (usually a week)
- Use staging Let's Encrypt server for testing:
  - Add `ca https://acme-staging-v02.api.letsencrypt.org/directory` to Caddyfile
  - Replace with production URL when ready

### Renewal Behavior

- Caddy automatically checks for renewal 30 days before expiration
- Renewal happens without downtime (Caddy hot-reloads certificates)
- No manual intervention required

To verify renewal is working:
```bash
docker logs claw-caddy-prod | tail -20
# Look for messages about certificate issuance or renewal
```

---

## Standalone VM

For production deployments on a virtual machine without Docker:

### Prerequisites

- Linux or macOS system
- `curl` installed
- A public domain with DNS A record pointing to your server
- Claw binary compiled: `go build -o claw .`

### Install Caddy

**Linux (apt):**
```bash
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.tailscale.com/stable/ubuntu/focal.noarmor.gpg' | sudo tee /usr/share/keyrings/tailscale-archive-keyring.gpg >/dev/null
curl -1sLf 'https://dl.tailscale.com/stable/ubuntu/focal.tailscale.list' | sudo tee /etc/apt/sources.list.d/tailscale.list

curl https://apt.fury.io/caddy/gpg.key | sudo tee /usr/share/keyrings/caddy-archive-keyring.gpg >/dev/null
echo "deb [signed-by=/usr/share/keyrings/caddy-archive-keyring.gpg] https://apt.fury.io/caddy/ any main" | sudo tee /etc/apt/sources.list.d/caddy-list.list

sudo apt update
sudo apt install caddy
```

**macOS (Homebrew):**
```bash
brew install caddy
```

**Other systems:** See [caddy.com/docs/install](https://caddyserver.com/docs/install)

### Configure Caddyfile

Copy the `Caddyfile` from the repository to your system:

```bash
# Option 1: Copy from repo
cp /path/to/claw/Caddyfile /etc/caddy/Caddyfile

# Option 2: Create manually
sudo tee /etc/caddy/Caddyfile > /dev/null <<'EOF'
yourdomain.com {
    reverse_proxy http://localhost:8080 {
        header_up X-Forwarded-For {remote_host}
        header_up X-Forwarded-Proto {scheme}
        header_up X-Forwarded-Host {host}
    }

    @http {
        protocol http
    }
    redir @http https://{host}{uri} permanent

    log {
        output stdout
        format json
    }
}
EOF
```

Replace `yourdomain.com` with your actual domain.

### Start Claw

```bash
# In one terminal, start Claw on localhost:8080
export CLAW_TOKEN="your-secure-token"
./claw -port 8080
```

### Start Caddy

```bash
# In another terminal, start Caddy (requires sudo for port 80/443)
sudo caddy run -config /etc/caddy/Caddyfile
```

### Systemd Integration (Optional)

For persistent daemon setup:

```bash
# Create systemd service for Claw
sudo tee /etc/systemd/system/claw.service > /dev/null <<'EOF'
[Unit]
Description=Claw MCP Server
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/claw -port 8080
Environment="CLAW_TOKEN=your-secure-token"
Restart=unless-stopped

[Install]
WantedBy=multi-user.target
EOF

# Create systemd service for Caddy
sudo tee /etc/systemd/system/caddy.service > /dev/null <<'EOF'
[Unit]
Description=Caddy (for Claw)
After=network.target

[Service]
Type=notify
ExecStart=/usr/bin/caddy run --config /etc/caddy/Caddyfile
Restart=unless-stopped

[Install]
WantedBy=multi-user.target
EOF

# Enable and start services
sudo systemctl daemon-reload
sudo systemctl enable claw caddy
sudo systemctl start claw caddy
sudo systemctl status claw caddy
```

### Verification

```bash
# Check if Caddy is running
sudo systemctl status caddy

# Test HTTPS access
curl https://yourdomain.com/health

# View certificate details
echo | openssl s_client -connect yourdomain.com:443 2>/dev/null | openssl x509 -noout -dates
```

---

## Certificate Renewal

### Docker Compose

Caddy automatically handles renewal:
- Checks 30 days before expiration
- Renews without downtime
- No manual action needed

To verify renewal logs:
```bash
docker logs claw-caddy-prod | grep -i "certificate\|renewal\|expire"
```

### Standalone VM

If using Caddy via systemd:
- Caddy auto-renews via systemd service
- Check logs: `sudo journalctl -u caddy -f`

If running Caddy manually:
- Restart Caddy periodically to apply renewed certificates
- Or run `sudo caddy reload -config /etc/caddy/Caddyfile` for zero-downtime reload

---

## Troubleshooting

### Caddy won't start (port 80/443 in use)

```bash
# Check what's using the ports
sudo lsof -i :80,443
sudo netstat -tuln | grep -E ':80|:443'

# Kill the process or change port in Caddyfile
# For Docker, ensure no other service is using these ports
docker ps | grep -E '80|443'
```

### Certificate not issuing (Let's Encrypt)

**Issue:** Caddy logs show ACME errors

**Diagnosis:**
1. Verify DNS: `nslookup yourdomain.com` should resolve to your server IP
2. Check port 80: `curl -v http://yourdomain.com` (should eventually redirect to HTTPS)
3. Check logs: `docker logs claw-caddy-prod` or `sudo journalctl -u caddy`

**Solution:**
- Wait a few minutes and restart Caddy
- Verify DNS propagation is complete (can take up to 48 hours)
- Check Let's Encrypt rate limits: https://letsencrypt.org/docs/rate-limits/

### HTTPS errors in browser

**Issue:** Browser shows certificate error

**Likely causes:**
- **Development (localhost):** Self-signed certificate is expected. Use browser's "unsafe" or "accept anyway" option.
- **Production:** Certificate not yet issued. Check Caddy logs and wait for ACME validation.
- **Browser cache:** Clear cache and reload (Ctrl+Shift+Delete).

### Claw unreachable through Caddy

**Issue:** Connection refused or timeout

**Diagnosis:**
```bash
# Docker: Check container network
docker network ls
docker network inspect claw-network
docker ps | grep claw

# Standalone: Check Claw is running
ps aux | grep claw
curl http://localhost:8080/health
```

**Solution:**
- Ensure Claw container is running: `docker ps`
- Check Claw logs: `docker logs claw-server` or `journalctl -u claw`
- Verify internal network in docker-compose: containers must be on same network
- Verify Caddyfile reverse proxy target: should be `http://claw:8080` (Docker) or `http://localhost:8080` (Standalone)

### Certificate expiry warnings

**Issue:** Certificate expired or about to expire

**Check expiry:**
```bash
# Docker
docker exec claw-caddy-prod caddy list-certs

# Standalone
openssl s_client -connect yourdomain.com:443 -showcerts 2>/dev/null | grep -A2 "Validity"
```

**Solution:**
- **Docker:** Caddy auto-renews 30 days before expiration (should be automatic)
- **Standalone:** Restart Caddy or run `sudo caddy reload`
- If not renewing: Check logs for ACME errors and DNS configuration

---

## Kubernetes

For Kubernetes deployments, use your Ingress controller to handle TLS/HTTPS:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: claw-ingress
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
  - hosts:
    - claw.example.com
    secretName: claw-cert
  rules:
  - host: claw.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: claw-service
            port:
              number: 8080
```

See the [Kubernetes README](kubernetes/README.md) for complete deployment instructions.

---

## Support

For issues or questions:
- Check Docker logs: `docker-compose logs -f`
- Check systemd logs: `sudo journalctl -u caddy -f`
- Review [Caddy documentation](https://caddyserver.com/docs/)
- See [Let's Encrypt documentation](https://letsencrypt.org/docs/) for certificate issues
