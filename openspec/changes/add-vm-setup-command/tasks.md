## 1. Application Architecture

- [x] 1.1 Make CLAW_TOKEN optional in auth middleware (allow empty/unset without error)
- [x] 1.2 Skip Bearer token validation when CLAW_TOKEN is empty or unset
- [x] 1.3 Verify /health endpoint remains unauthenticated regardless of CLAW_TOKEN

## 2. Caddyfile Embedding

- [x] 2.1 Add embed package import to main.go
- [x] 2.2 Create //go:embed directive to embed Caddyfile into binary
- [x] 2.3 Verify embedded Caddyfile is accessible at runtime
- [x] 2.4 Test that binary size includes embedded template (~1KB)

## 3. Setup Command Infrastructure

- [x] 3.1 Add `setup` subcommand to CLI (use flag.NewFlagSet or similar)
- [x] 3.2 Create setup.go module with main Setup() function
- [x] 3.3 Add interactive prompt helpers (getInput, getPassword functions)
- [x] 3.4 Implement password input without terminal echo (getPassword using terminal.ReadPassword or similar)

## 4. Setup Flow - Domain and Password Collection

- [x] 4.1 Prompt for domain name with validation (non-empty, basic DNS format check)
- [x] 4.2 Re-prompt if domain is invalid
- [x] 4.3 Prompt for password with confirmation (getPassword twice, compare)
- [x] 4.4 Re-prompt if passwords don't match or is empty

## 5. Bcrypt Hash Generation

- [x] 5.1 Import golang.org/x/crypto/bcrypt
- [x] 5.2 Generate bcrypt hash of password using bcrypt.GenerateFromPassword()
- [x] 5.3 Handle bcrypt errors gracefully (e.g., if cost is invalid)
- [x] 5.4 Verify hash format matches Caddy expectations ($2a$14$...)

## 6. Caddyfile Extraction and Customization

- [x] 6.1 Extract embedded Caddyfile from embedded FS
- [x] 6.2 Replace {$DOMAIN} placeholder with user-provided domain
- [x] 6.3 Add basicauth directive: `basicauth /mcp/* { admin $bcrypt_hash }` (Caddy v2.6.2 and earlier)
- [x] 6.4 Verify customized Caddyfile contains all expected directives
- [x] 6.5 Check if /etc/caddy/Caddyfile already exists (error if it does, prevent overwrite)
- [x] 6.6 Write customized Caddyfile to /etc/caddy/Caddyfile with appropriate permissions

## 7. Marker File Management

- [x] 7.1 Create /var/lib/claw directory if it doesn't exist
- [x] 7.2 Check for /var/lib/claw/.setup-done marker at start of setup command
- [x] 7.3 If marker exists, error with helpful message (already initialized)
- [x] 7.4 Create marker file after setup completes successfully

## 8. Systemd Service Configuration - Claw

- [x] 8.1 Create Claw systemd service template (ExecStart, Environment, Restart, etc.)
- [x] 8.2 Write template to /etc/systemd/system/claw.service
- [x] 8.3 Include CLAW_TOKEN=empty in service (or rely on env var not set)
- [x] 8.4 Set appropriate permissions on service file

## 9. Systemd Service Configuration - Caddy

- [x] 9.1 Create Caddy systemd service template (ExecStart, After, Type=notify, etc.)
- [x] 9.2 Write template to /etc/systemd/system/caddy.service
- [x] 9.3 Set Caddy to start after Claw service
- [x] 9.4 Set appropriate permissions on service file

## 10. Systemd Service Lifecycle

- [x] 10.1 Run `systemctl daemon-reload` after creating service files
- [x] 10.2 Run `systemctl enable claw caddy` to enable auto-start
- [x] 10.3 Run `systemctl start claw caddy` to start services immediately
- [x] 10.4 Verify services are running with `systemctl status`
- [x] 10.5 Handle errors from systemctl commands (provide recovery guidance)

## 11. Error Handling and User Feedback

- [x] 11.1 Validate all inputs (domain, password, filesystem permissions)
- [x] 11.2 Provide clear error messages for each failure point
- [x] 11.3 Suggest recovery steps if setup fails (e.g., "Run with sudo for /etc/caddy access")
- [x] 11.4 Display success message with verification commands at end

## 12. Testing and Verification

- [ ] 12.1 Test setup on Linux VM (manual or test environment)
- [ ] 12.2 Verify /etc/caddy/Caddyfile contains correct domain and bcrypt hash
- [ ] 12.3 Verify systemd services are enabled and running
- [ ] 12.4 Test that second run of setup errors correctly
- [ ] 12.5 Verify HTTPS works on configured domain
- [ ] 12.6 Test basic auth login with username `admin` and user's password
- [ ] 12.7 Verify /mcp endpoint requires basic auth (Caddy enforces it)
- [ ] 12.8 Verify /health endpoint is accessible without auth

**Note:** Testing tasks require a VM deployment environment and manual verification. See implementation notes below.
