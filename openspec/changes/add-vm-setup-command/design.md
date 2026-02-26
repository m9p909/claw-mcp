## Context

Users currently must manually configure Claw on VMs: install Caddy, generate bcrypt password hashes, edit Caddyfile, set up systemd services. The proposed `./mcpclaw setup` command automates this workflow. This is a single-user focused setup (expand to multi-user auth later if needed).

## Goals / Non-Goals

**Goals:**
- Provide interactive guided setup for VM deployments
- Automate Caddyfile generation with user's domain + password hash
- Set up systemd daemons for persistent operation
- Prevent re-initialization (idempotent, fail if already set up)
- Keep auth configuration at Caddy layer (users own it, not Claw)

**Non-Goals:**
- Docker setup (separate change, future)
- Multi-user authentication (single admin user for now)
- Password reset/reconfiguration (initial setup only)
- Automatic Caddy installation (users install Caddy first)
- Kubernetes setup guidance (Ingress controllers handle auth)

## Decisions

**1. Caddyfile as embedded template**
- Rationale: Ensures consistency, survives binary upgrades, easy to maintain
- Alternative: Fetch from GitHub (adds network dependency, slower setup)
- Implementation: Use `embed.FS` to include Caddyfile in binary

**2. Extract to `/etc/caddy/Caddyfile` (user's copy)**
- Rationale: Standard location, survives Claw updates, user can customize later
- Alternative: Generate unique filename like `/etc/caddy/Caddyfile.claw` (confusing)
- Implementation: Check if file exists; if so, error (prevents accidental overwrites)

**3. Marker file at `/var/lib/claw/.setup-done`**
- Rationale: Survives reboots, check prevents running setup twice
- Alternative: Check if systemd services exist (fragile, could be manually edited)
- Implementation: Create after all setup completes; check at start of `setup` command

**4. Password hashing in Go (no shell out to `caddy hash-password`)**
- Rationale: Self-contained, no external dependency, works without Caddy binary on $PATH
- Alternative: Shell out to `caddy hash-password` (requires Caddy already installed)
- Implementation: Use `golang.org/x/crypto/bcrypt` for bcrypt hashing

**5. Single username: `admin`**
- Rationale: Matches common defaults, simple for initial setup, easy to document
- Alternative: Prompt for username (adds complexity, no clear benefit yet)
- Implementation: Hard-code in setup command; document in output

**6. CLAW_TOKEN optional (empty = no auth)**
- Rationale: Delegation to Caddy auth is the primary approach; token becomes fallback
- Alternative: Always require CLAW_TOKEN (confusing if Caddy handles it)
- Implementation: Auth middleware checks if token is set before validating

**7. Systemd services for claw + caddy**
- Rationale: Persistent daemon management, standard Linux practice
- Alternative: systemd user services (requires user-level setup, less common)
- Implementation: Create `/etc/systemd/system/claw.service` and `caddy.service` templates

## Risks / Trade-offs

**[Risk: `/var/lib/claw` permissions]** → Setup command requires write access; document that `sudo ./mcpclaw setup` is required, or setup creates dir with user's permissions if possible.

**[Risk: Existing Caddyfile]** → If user already has `/etc/caddy/Caddyfile`, setup errors out. Mitigation: Document workaround (backup old file, re-run setup).

**[Risk: systemd service failures]** → If `systemctl enable/start` fails (permissions, systemd not available), setup errors. Mitigation: Provide clear error message with fallback manual steps.

**[Risk: Password visibility in logs]** → If user pastes command with password, it could be logged. Mitigation: Prompt interactively (stdin) instead of CLI arg, document not to use `--password` flag.

**[Risk: Single admin user]** → No multi-user auth initially. Mitigation: Document that users can edit Caddyfile later to add more users.

**[Risk: Re-running `setup` after changes]** → User can't re-run if they mess up (marker file prevents it). Mitigation: Document manual recovery steps or provide `./mcpclaw setup --reset` later.
