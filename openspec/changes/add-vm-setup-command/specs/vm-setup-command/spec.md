## ADDED Requirements

### Requirement: Interactive setup command exists
The Claw binary SHALL expose a `setup` subcommand that guides users through configuring Claw on a Linux/macOS VM with TLS and basic authentication.

#### Scenario: User runs setup command
- **WHEN** user executes `./mcpclaw setup`
- **THEN** system assumes VM deployment and prompts for domain name

### Requirement: User provides domain name
The setup command SHALL prompt for a domain name (for the VM's public address).

#### Scenario: User enters valid domain
- **WHEN** user is prompted for domain and enters "claw.example.com"
- **THEN** system accepts the domain and proceeds to password prompt

#### Scenario: User enters invalid domain
- **WHEN** user enters empty string or invalid format
- **THEN** system re-prompts for a valid domain

### Requirement: User provides password
The setup command SHALL prompt user to enter a password for basic auth (username is always `admin`).

#### Scenario: User enters password interactively
- **WHEN** user is prompted for password
- **THEN** system reads password from stdin without echoing to terminal

#### Scenario: Password is empty
- **WHEN** user provides empty password
- **THEN** system re-prompts for a non-empty password

### Requirement: Password hash is generated
The system SHALL generate a bcrypt hash of the user's password and use it in Caddyfile configuration.

#### Scenario: Bcrypt hash creation
- **WHEN** user provides password "mysecret"
- **THEN** system generates a bcrypt hash ($2a$14$...) and uses it in Caddyfile

#### Scenario: Hash is not stored in plaintext
- **WHEN** setup completes
- **THEN** plaintext password is not written to disk; only the bcrypt hash is used

### Requirement: Caddyfile is extracted and customized
The system SHALL extract the embedded Caddyfile template, customize it with the user's domain and password hash, and write it to `/etc/caddy/Caddyfile`.

#### Scenario: Caddyfile extraction
- **WHEN** setup completes successfully
- **THEN** `/etc/caddy/Caddyfile` exists with user's domain and bcrypt hash embedded

#### Scenario: Caddyfile includes basic auth directive
- **WHEN** Caddyfile is created
- **THEN** it includes `basic_auth /mcp/* { admin $bcrypt_hash }`

#### Scenario: Caddyfile handles both self-signed and Let's Encrypt
- **WHEN** Caddyfile is generated for a custom domain (not localhost)
- **THEN** Caddy will use Let's Encrypt via ACME for certificate provisioning

### Requirement: Systemd services are configured
The system SHALL create systemd service files for both Claw and Caddy, enable them, and start them.

#### Scenario: Claw systemd service created
- **WHEN** setup completes
- **THEN** `/etc/systemd/system/claw.service` exists and is enabled

#### Scenario: Caddy systemd service created
- **WHEN** setup completes
- **THEN** `/etc/systemd/system/caddy.service` exists and is enabled

#### Scenario: Services are started
- **WHEN** setup completes
- **THEN** both claw and caddy services are started and running

#### Scenario: Services survive reboot
- **WHEN** VM is rebooted
- **THEN** claw and caddy services start automatically

### Requirement: Setup cannot be run twice
The system SHALL prevent running `./mcpclaw setup` if it has already been initialized on the system.

#### Scenario: First run succeeds
- **WHEN** user runs `./mcpclaw setup` for the first time
- **THEN** setup completes and creates marker file at `/var/lib/claw/.setup-done`

#### Scenario: Second run fails
- **WHEN** user runs `./mcpclaw setup` again
- **THEN** system detects marker file and exits with error "Setup already completed on this system"

#### Scenario: Error message is helpful
- **WHEN** user attempts to re-run setup
- **THEN** system displays instructions for manual recovery or resetting setup

### Requirement: CLAW_TOKEN is optional
The application SHALL start without requiring `CLAW_TOKEN` environment variable. When unset or empty, the `/mcp` endpoint has no Bearer token validation.

#### Scenario: Claw starts without CLAW_TOKEN
- **WHEN** `CLAW_TOKEN` is unset or empty string
- **THEN** Claw starts successfully and `/mcp` endpoint does not validate Bearer token

#### Scenario: Claw starts with CLAW_TOKEN
- **WHEN** `CLAW_TOKEN` is set to a value
- **THEN** Claw requires Bearer token on `/mcp` endpoint (existing behavior)

#### Scenario: Auth is delegated to Caddy
- **WHEN** Caddy is running with basic auth and CLAW_TOKEN is empty
- **THEN** `/mcp` endpoint is protected by Caddy's basic auth, not Claw's token validation

### Requirement: Embedded Caddyfile is pristine
The Caddyfile in the repository SHALL be the source template. When compiled, it is embedded in the binary as-is and extracted during setup.

#### Scenario: Repository Caddyfile is unchanged
- **WHEN** `./mcpclaw setup` runs
- **THEN** the extracted Caddyfile matches the one in the repository (before customization)

#### Scenario: Setup customizes the extracted copy
- **WHEN** setup extracts and customizes the Caddyfile
- **THEN** only the copy at `/etc/caddy/Caddyfile` is modified; repository Caddyfile remains unchanged
