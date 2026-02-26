## Why

Setting up Claw on a VM requires multiple manual steps: installing Caddy, generating password hashes, configuring Caddyfile, setting up systemd services. A guided setup command reduces friction and ensures users correctly configure TLS + basic auth with best practices.

## What Changes

- Add `./mcpclaw setup` command for interactive VM setup
- Embed Caddyfile template in binary
- Extract and customize Caddyfile with user-provided domain and generated password hash
- Configure systemd services for claw + caddy daemons
- Guard against re-running (fail if already initialized)
- Make `CLAW_TOKEN` optional in application (auth delegated to Caddy basic auth)

## Capabilities

### New Capabilities
- `vm-setup-command`: Interactive setup command for configuring Claw on a Linux/macOS VM with TLS and basic auth

### Modified Capabilities
- `bearer-token-auth`: Make authentication optional when empty/unset (allow Caddy to handle auth instead)

## Impact

- **Application code**: Add CLI command to main.go, make CLAW_TOKEN optional in auth middleware
- **Binary size**: Small increase (embed Caddyfile template ~1KB)
- **No breaking changes**: Existing deployments continue to work with CLAW_TOKEN
- **VM users**: Get guided setup instead of manual steps
- **systemd integration**: Claw and Caddy run as daemons on VM
