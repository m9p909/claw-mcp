## ADDED Requirements

### Requirement: Docker Compose TLS deployment guide
The system SHALL provide clear documentation for deploying Claw with TLS using Docker Compose in both development and production environments.

#### Scenario: Developer deploys to localhost with self-signed TLS
- **WHEN** developer follows the development Docker Compose guide
- **THEN** they can start Claw with TLS in under 3 commands
- **AND** guide includes how to trust/bypass self-signed certificate warnings

#### Scenario: Operator deploys to production with Let's Encrypt
- **WHEN** operator follows the production Docker Compose guide
- **THEN** they can deploy with automatic Let's Encrypt certificates by setting two environment variables (DOMAIN, CLAW_TOKEN)
- **AND** guide explains DNS requirements and validation process
- **AND** guide includes certificate renewal and troubleshooting steps

### Requirement: Standalone VM TLS setup guide
The system SHALL provide documentation for deploying Claw with TLS on standalone virtual machines using Caddy.

#### Scenario: User installs Caddy on a VM
- **WHEN** user follows the "Caddy installation" section of the guide
- **THEN** guide provides OS-specific installation commands (Linux/macOS)
- **AND** user can verify Caddy is installed and running

#### Scenario: User configures Caddy as reverse proxy to Claw
- **WHEN** user follows the "Configure Caddyfile" section
- **THEN** guide provides template Caddyfile with example domain substitution
- **AND** user can start Claw and Caddy independently (separate processes/ports)
- **AND** Claw listens on localhost:8080
- **AND** Caddy listens on 0.0.0.0:80 and :443

#### Scenario: User enables Let's Encrypt for production domain
- **WHEN** user updates Caddyfile with production domain and restarts Caddy
- **THEN** guide explains DNS A record requirements
- **AND** guide explains automatic certificate renewal behavior
- **AND** user can verify certificate validity

### Requirement: TLS setup guide contents
The system SHALL provide comprehensive guide content covering Docker Compose and standalone VM deployments.

#### Scenario: Guide includes prerequisites section
- **WHEN** user reads the guide introduction
- **THEN** guide lists prerequisites (Docker/Caddy version, domain name for prod, etc.)

#### Scenario: Guide includes troubleshooting section
- **WHEN** user encounters certificate or connection issues
- **THEN** guide provides diagnostic steps (certificate validation, port checking, Caddy logs)
- **AND** guide includes common error messages and solutions

#### Scenario: Guide includes security considerations
- **WHEN** user reads the guide
- **THEN** guide mentions certificate storage and renewal processes
- **AND** guide explains that Kubernetes users should use Ingress instead

### Requirement: README references TLS setup guide
The system SHALL update the main README to reference the TLS setup guide for production users.

#### Scenario: README security section links to guide
- **WHEN** user reads README.md security section
- **THEN** TLS guide is prominently linked
- **AND** current README statement about reverse proxy is updated to say "See TLS_SETUP.md"
